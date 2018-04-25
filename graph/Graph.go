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
}

type Graph struct {
	Edges []Arc
}

func MakeGraph(problemFormulation io.ProblemFormulation) {
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
			arcs = append(arcs, Arc{previous, node, previous.Time, false})
			nodes = append(nodes, node)
			previous = node
		}
		arcs = append(arcs, Arc{previous, sink, previous.Time, false})
	}
	nodes = append(nodes, sink)

	// Create disjunctive arcs (jobs belonging to machines)
	for _, nodePtrs := range machineToNodesMap {
		for i := range nodePtrs {
			for j := i; j < len(nodePtrs); j++ {
				arcs = append(arcs, Arc{nodePtrs[i],
					nodePtrs[j],
					nodePtrs[i].Time,
					true})
			}
		}
	}

	neighbours := make(map[Node][]Node)
	for _, edge := range conjunctiveArcs {
		neighbours[edge.From] = append(neighbours[edge.From], edge.To)
	}
	for _, edge := range disjunctiveArcs {
		neighbours[edge.From] = append(neighbours[edge.From], edge.To)
	}
}
