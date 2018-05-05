package aco

import (
	"../constants"
	graph "../graph"
	"../util"
	"fmt"
	"math"
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

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	//arcPheroMap := InitializePheromoneValues(problemGraph)
	//var iterationBest Solution		//Sib
	//var bestSoFar []graph.Node			//Sbs
	//var restartBest []graph.Node		//Srb
	//convergenceFactor := 0.0			// cf
	//bsUpdate := false
	//numberOfAnts := numberOfAnts(problemGraph)
	//for true {
	//	solutions := make([]Solution, 0)
	//	for i := 0; i < numberOfAnts; i++ { cx /*<dfgdzg<drg<orgsrijgsirgjrigjsigj<rgjr<j<rgjrjgorjggp<erjoprgjropg*/
	//		solutions = append(solutions, listScheduler(problemGraph))
	//	}
	//	ApplyLocalSearch(solutions)
	//	iterationBest = SolutionWithMinimalMakeSpan(solutions)
	//	EliteAction(iterationBest)
	//	Update(iterationBest, restartBest, bestSoFar)
	//	ApplyPheromoneUpdate(convergenceFactor, bsUpdate, T, restartBest, bestSoFar)
	//	if convergenceFactor > 0.99 {
	//		if bsUpdate {
	//			ResetPheromoneValues(T)
	//			restartBest = nil
	//			bsUpdate = false
	//		} else {
	//			bsUpdate = true
	//		}
	//	}
	//}


	//bestSoFar := listScheduler(problemGraph)
	//fmt.Println(bestSoFar.Nodes)
	//fmt.Println("Done!")
	//fmt.Println("Best makespan:", calculateMakespan(bestSoFar))
	//orders := graph.NodeListToOrderList(bestSoFar.Nodes, bestSoFar.StartTimeMap)
	//gantt.CreateChart("03 - Program Outputs/Chart.xlsx", orders)

}
