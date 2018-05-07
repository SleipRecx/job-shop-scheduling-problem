package aco

import (
	"../constants"
	"../graph"
	"../util"
	"fmt"
	"math"
	"math/rand"
	"errors"
	"../gantt"
)

func removeFromList(list []graph.Node, element graph.Node) []graph.Node {
	newList := make([]graph.Node, 0)
	for _, node := range list {
		if !(node.Job == element.Job && node.Machine == element.Machine) {
			newList = append(newList, node)
		}
	}
	return newList
}

type Solution struct {
	StartTimeMap map[graph.Node]int
	Nodes        []graph.Node
}

func calculateMakespan(solution Solution) int {
	max := 0
	for x := range solution.Nodes {
		if time := solution.StartTimeMap[solution.Nodes[x]] + solution.Nodes[x].Time; time > max {
			max = time
		}
	}
	return max
}

func numberOfAnts(problemGraph graph.Graph) int {
	return util.Max(10, len(problemGraph.Nodes)/10)
}

func InitializePheromoneValues(problemGraph graph.Graph) map[graph.Arc]float64 {
	mapping := make(map[graph.Arc]float64)
	for _, arc := range problemGraph.Edges {
		mapping[arc] = constants.InitialPheromone
	}
	return mapping
}

func SolutionWithMinimalMakeSpan(solutions []Solution) Solution {
	minMakespan := math.MaxInt32
	var bestSolution Solution
	for _, sol := range solutions {
		if ms := calculateMakespan(sol); ms < minMakespan {
			minMakespan = ms
			bestSolution = sol
		}
	}
	return bestSolution
}

func findCriticalPath(nodes []graph.Node) []graph.Node{
	incomingNeighbourList := make(map[graph.Node][]graph.Node)

	// Construct solution graph
	for i := 0; i < len(nodes); i++ {
		for j := i; j < len(nodes); j++ {
			if nodes[i].Job == nodes[j].Job && nodes[i].TechStep == nodes[j].TechStep-1 {
				incomingNeighbourList[nodes[j]] = append(incomingNeighbourList[nodes[j]], nodes[i])
			} else if nodes[i].Machine == nodes[j].Machine && nodes[i] != nodes[j] {
				incomingNeighbourList[nodes[j]] = append(incomingNeighbourList[nodes[j]], nodes[i])
			}
		}
	}

	// Calculate weight at nodes
	nodeWeight := make(map[graph.Node]int)
	parent := make(map[graph.Node]graph.Node)

	for _, node := range nodes {
		max := 0
		var maxNode graph.Node
		for _, incoming := range incomingNeighbourList[node] {
			if value := nodeWeight[incoming] + incoming.Time; value > max {
				max = value
				maxNode = incoming
			}
		}
		nodeWeight[node] = max
		parent[node] = maxNode
	}

	// Reconstruct critical path
	var node graph.Node
	max := 0
	for _, n := range nodes {
		if w := nodeWeight[n]; w > max {
			max = w
			node = n
		}
	}

	criticalPath := []graph.Node{node}
	nilNode := graph.Node{0,0,0,0}
	for true {
		node = parent[node]
		if node == nilNode {
			break
		}
		criticalPath = append(criticalPath, node)
	}
	return criticalPath
}

func ApplyLocalSearch(solution Solution) Solution{
	criticalPath := findCriticalPath(solution.Nodes)
	blocks := make([][]graph.Node, 0)
	currentBlock := []graph.Node{criticalPath[0]}
	for i := 1; i < len(criticalPath); i++ {
		if criticalPath[i].Machine == currentBlock[len(currentBlock)-1].Machine {
			currentBlock = append(currentBlock, criticalPath[i])
		} else {
			blocks = append(blocks, currentBlock)
			currentBlock = []graph.Node{criticalPath[i]}
		}
	}
	blocks = append(blocks, currentBlock)

	// Create new possible solution by copying old one
	newNodes := make([]graph.Node, 0)
	newStartTimeMap := make(map[graph.Node]int)
	for _, n := range solution.Nodes {
		newNodes = append(newNodes, n)
	}
	for k, v := range solution.StartTimeMap {
		newStartTimeMap[k] = v
	}
	// Try to swap something on the critical path
	for _, block := range blocks {
		if len(block) >= 2 {
			index := rand.Intn(len(block)-1)
			startTime0 := newStartTimeMap[block[index]] + block[index+1].Time
			startTime1 := newStartTimeMap[block[index]]
			newStartTimeMap[block[index]] = startTime0
			newStartTimeMap[block[index+1]] = startTime1
		}
	}
	if calculateMakespan(Solution{newStartTimeMap, newNodes}) < calculateMakespan(solution) {
		return Solution{newStartTimeMap, newNodes}
	}
	return solution
}

func isDone(solution Solution) bool {
	if len(solution.Nodes) == 0 {
		return false
	}
	return float64(calculateMakespan(solution)) < 1.1 * float64(constants.TargetValues[constants.ProblemNumber])
}

func Update(original, candidate Solution) Solution {
	if len(original.Nodes) == 0 {
		return  candidate
	} else if calculateMakespan(original) < calculateMakespan(candidate) {
		return original
	}
	return candidate
}

func index(n graph.Node, nodes []graph.Node) (int, error) {
	for i := range nodes {
		if nodes[i] == n {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func δ(n1, n2 graph.Node, solution Solution) float64{
	n1Index, _ := index(n1, solution.Nodes)
	n2Index, _ := index(n2, solution.Nodes)
	if n1Index < n2Index {
		return 1
	}
	return 0
}

func fmmaas(x float64) float64 {
	if x < constants.TMin {
		return constants.TMin
	}
	if x > constants.TMax {
		return constants.TMax
	}
	return x
}

func ApplyPheromoneUpdate(cf float64, bsUpdate bool, arcPheroMap map[graph.Arc]float64, restartBest, bestSoFar Solution) (map[graph.Arc]float64, float64) {
	solutionToUse := restartBest
	if bsUpdate {
		solutionToUse = bestSoFar
	}
	cfUpper := 0.0
	cfLower := float64(len(arcPheroMap)) * (constants.TMax - constants.TMin)
	for arc, tij := range arcPheroMap {
		arcPheroMap[arc] = fmmaas(tij + constants.EvaporationRate * (δ(arc.From, arc.To, solutionToUse)-tij))
		cfUpper += math.Max(constants.TMax - tij, tij - constants.TMin)
	}
	return arcPheroMap, 2 * ((cfUpper / cfLower) - 0.5)
}

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	arcPheroMap := InitializePheromoneValues(problemGraph)
	var iterationBest Solution		//Sib
	var bestSoFar Solution			//Sbs
	var restartBest Solution		//Srb
	convergenceFactor := 0.0			// cf
	bsUpdate := false
	numberOfAnts := numberOfAnts(problemGraph)
	counter := 0
	for !isDone(bestSoFar) {
		counter += 1
		if counter % 10 == 0 {
			fmt.Println("Iteration", counter, "Makespan", calculateMakespan(bestSoFar))
		}
		solutions := make([]Solution, 0)
		for i := 0; i < numberOfAnts; i++ {
			solutions = append(solutions, listScheduler(problemGraph, arcPheroMap))
		}
		for i := range solutions {
			solutions[i] = ApplyLocalSearch(solutions[i])
		}
		iterationBest = SolutionWithMinimalMakeSpan(solutions)
		//EliteAction(iterationBest)
		bestSoFar = Update(bestSoFar, iterationBest)
		restartBest = Update(restartBest, iterationBest)

		arcPheroMap, convergenceFactor = ApplyPheromoneUpdate(convergenceFactor, bsUpdate, arcPheroMap, restartBest, bestSoFar)
		if convergenceFactor > 0.99 {
			if bsUpdate {
				arcPheroMap = InitializePheromoneValues(problemGraph)
				restartBest = Solution{}
				bsUpdate = false
			} else {
				bsUpdate = true
			}
		}
	}

	fmt.Println("Done!")
	fmt.Println("Makespan:", calculateMakespan(bestSoFar))
	orders := graph.NodeListToOrderList(bestSoFar.Nodes, bestSoFar.StartTimeMap)
	gantt.CreateChart("03 - Program Outputs/Chart.xlsx", orders)

}
