package graph

import (
	"../io"
	"fmt"
)

type Node struct {
	Job, Machine, Time int
}

type ConjunctiveArc struct {
	From, To Node
	Time     int
}

type DisjunctiveArc struct {
	From, To Node
	Time     int
}

func MakeGraph(problemFormulation io.ProblemFormulation) {
	source := Node{-1, -1, 0}
	sink := Node{-1, -1, 0}
	nodes := []Node{source}
	conjunctiveArcs := make([]ConjunctiveArc, 0)
	disjunctiveArcs := make([]DisjunctiveArc, 0)
	machineToNodesMap := make(map[int][]Node)

	// Create conjunctive arcs (technological order)
	for jobId, requirements := range problemFormulation.Sequences {
		previous := source
		for _, requirement := range requirements {
			node := Node{jobId, requirement.Machine, requirement.Time}
			machineToNodesMap[requirement.Machine] = append(machineToNodesMap[requirement.Machine], node)
			conjunctiveArcs = append(conjunctiveArcs, ConjunctiveArc{previous, node, previous.Time})
			nodes = append(nodes, node)
			previous = node
		}
		conjunctiveArcs = append(conjunctiveArcs, ConjunctiveArc{previous, sink, previous.Time})
	}
	nodes = append(nodes, sink)

	// Create disjunctive arcs (jobs belonging to machines)
	for _, nodePtrs := range machineToNodesMap {
		for i := range nodePtrs {
			for j := i; j < len(nodePtrs); j++ {
				disjunctiveArcs = append(disjunctiveArcs, DisjunctiveArc{nodePtrs[i],
					nodePtrs[j],
					nodePtrs[i].Time})
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
