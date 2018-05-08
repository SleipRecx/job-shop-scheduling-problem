package jssp

import (
	"../graph"
	"../constants"
	"math"
	"math/rand"
)

type Solution struct {
	StartTimeMap map[graph.Node]int
	Nodes        []graph.Node
}

func CalculateMakespan(solution Solution) int {
	max := 0
	for x := range solution.Nodes {
		if time := solution.StartTimeMap[solution.Nodes[x]] + solution.Nodes[x].Time; time > max {
			max = time
		}
	}
	return max
}

func MakeSpan2(solution Solution) int {
	ms := CalculateMakespan(solution)
	if ms == 0 {
		return math.MaxInt32
	}
	return ms
}

func SolutionWithMinimalMakeSpan(solutions []Solution) Solution {
	minMakespan := math.MaxInt32
	var bestSolution Solution
	for _, sol := range solutions {
		if ms := CalculateMakespan(sol); ms < minMakespan {
			minMakespan = ms
			bestSolution = sol
		}
	}
	return bestSolution
}

func FindCriticalPath(nodes []graph.Node) []graph.Node{
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
	criticalPath := FindCriticalPath(solution.Nodes)
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
	if CalculateMakespan(Solution{newStartTimeMap, newNodes}) < CalculateMakespan(solution) {
		return Solution{newStartTimeMap, newNodes}
	}
	return solution
}

func IsDone(solution Solution) bool {
	if len(solution.Nodes) == 0 {
		return false
	}
	return float64(CalculateMakespan(solution)) < 1.1 * float64(constants.TargetValues[constants.ProblemNumber])
}

func GenerateNeighbourHood(solution Solution) []Solution {
	criticalPath := FindCriticalPath(solution.Nodes)
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

	neighbours := []Solution{solution}
	for _, block := range blocks {
		if len(block) >= 2 {
			for i := 0; i < len(block) - 1; i++ {
				// Create new possible solution by copying old one
				newNodes := make([]graph.Node, 0)
				newStartTimeMap := make(map[graph.Node]int)
				for _, n := range solution.Nodes {
					newNodes = append(newNodes, n)
				}
				for k, v := range solution.StartTimeMap {
					newStartTimeMap[k] = v
				}
				startTime0 := newStartTimeMap[block[i]] + block[i+1].Time
				startTime1 := newStartTimeMap[block[i]]
				newStartTimeMap[block[i]] = startTime0
				newStartTimeMap[block[i+1]] = startTime1
				neighbours = append(neighbours, Solution{newStartTimeMap, newNodes})
			}
		}
	}
	return neighbours
}