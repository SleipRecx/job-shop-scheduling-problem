package graph

import (
	"../io"
	"fmt"
)

type Node struct {
	Job, Machine, Time int
}

type Arc struct {
	From, To Node
	Weight     int
	Disjunct bool
	Pheromone float64
}

type Graph struct {
	Edges []Arc
	Nodes []Node
	NeighbourList map[Node][]Node
}


func MakeGraph(problemFormulation io.ProblemFormulation) Graph {
	source := Node{-1, -1, 0}
	sink := Node{-1, -1, 0}
	nodes := []Node{source}
	arcs := make([]Arc, 0)
	machineToNodesMap := make(map[int][]Node)

	// Create conjunctive arcs (technological order)
	for jobId, requirements := range problemFormulation.Sequences {
		previous := source
		for _, requirement := range requirements {
			node := Node{jobId, requirement.Machine, requirement.Time}
			machineToNodesMap[requirement.Machine] = append(machineToNodesMap[requirement.Machine], node)
			arcs = append(arcs, Arc{previous,
				node,
				previous.Time,
				false,
				0.0 })
			nodes = append(nodes, node)
			previous = node
		}
		arcs = append(arcs, Arc{previous,
			sink,
			previous.Time,
			false,
			0.0})
	}
	nodes = append(nodes, sink)

	// Create disjunctive arcs (jobs belonging to machines)
	for _, nodePtrs := range machineToNodesMap {
		for i := range nodePtrs {
			for j := i; j < len(nodePtrs); j++ {
				arcs = append(arcs, Arc{nodePtrs[i],
					nodePtrs[j],
					nodePtrs[i].Time,
					true,
					0.0})
			}
		}
	}

	neighbours := make(map[Node][]Node)
	for _, edge := range arcs {
		neighbours[edge.From] = append(neighbours[edge.From], edge.To)
	}

	return Graph{arcs, nodes, neighbours}
}
