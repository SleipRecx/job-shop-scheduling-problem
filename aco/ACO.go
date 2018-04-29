package aco

import (
	"../graph"
	_ "../constants"
	"fmt"
	"math"
)

func removeFromList(list []graph.Node, element graph.Node) []graph.Node {
	newList := make([]graph.Node, 0)
	for _, node := range list {
		if node != element {
			newList = append(newList, node)
		}
	}
	return newList
}

func evaluate(permutation []graph.Node) {

}


func listScheduler(problemGraph graph.Graph) {
	partialSolution := make([]graph.Node, 0)
	unvisited := make([]graph.Node, len(problemGraph.Nodes))
	copy(unvisited, problemGraph.Nodes)
	for range problemGraph.Nodes {
		unvisitedPrime := restrict(partialSolution, unvisited)
		nodeStar := choose(unvisitedPrime)
		partialSolution = append(partialSolution, nodeStar)
		unvisited = removeFromList(unvisited, nodeStar)
	}
}

func choose(nodes []graph.Node) graph.Node {

	return graph.Node{}
}

func earliestCompletionTime(node graph.Node, partialSolution []graph.Node) int {
	if len(partialSolution) == 0 {
		return node.Time
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range partialSolution {
		if node.Job == partialSolution[x].Job && partialSolution[x].StartTime > JobTimer {
			JobTimer = partialSolution[x].StartTime
		}
		if node.Machine == partialSolution[x].Machine && partialSolution[x].StartTime > MachineTimer {
			MachineTimer = partialSolution[x].StartTime
		}
	}
	earliestComp := math.Max(float64(JobTimer),float64(MachineTimer)) + float64(node.Time)
	return int(earliestComp)
}

func restrict(partialSolution []graph.Node, unVisited []graph.Node) []graph.Node {

}

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	listScheduler(problemGraph)
}
