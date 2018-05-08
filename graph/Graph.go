package graph

import (
	"../gantt"
	"../io"
)

type Node struct {
	Job, Machine, Time, TechStep int
}

type Arc struct {
	From, To Node
}

type Graph struct {
	Edges         []Arc
	Nodes         []Node
	NeighbourList map[Node][]Node
}

func NodeListToOrderList(nodes []Node, startTimeMap map[Node]int) []gantt.Order {
	orders := make([]gantt.Order, 0)
	for _, n := range nodes {
		orders = append(orders, gantt.Order{n.Job, n.Machine, startTimeMap[n], n.Time})
	}
	return orders
}

func (g *Graph) isNeighbour(i, j Node) bool {
	for x := range g.NeighbourList[i] {
		if g.NeighbourList[i][x] == j {
			return true
		}
	}
	return false
}

func (g *Graph) findEdge(from, to Node) Arc {
	for x := range g.Edges {
		if g.Edges[x].From == from && g.Edges[x].To == to {
			return g.Edges[x]
		}
	}
	return Arc{}
}

func MakeGraph(problemFormulation io.ProblemFormulation) Graph {
	source := Node{-1, -1, 0, -1}
	sink := Node{-1, -1, 0, -1}
	nodes := []Node{}
	arcs := make([]Arc, 0)
	machineToNodesMap := make(map[int][]Node)

	// Create conjunctive arcs (technological order)
	for jobId, requirements := range problemFormulation.Sequences {
		previous := source
		for index, requirement := range requirements {
			node := Node{jobId, requirement.Machine, requirement.Time, index}
			machineToNodesMap[requirement.Machine] = append(machineToNodesMap[requirement.Machine], node)
			arcs = append(arcs, Arc{previous,
				node})
			nodes = append(nodes, node)
			previous = node
		}
		arcs = append(arcs, Arc{previous,
			sink})
	}

	// Create disjunctive arcs (jobs belonging to machines)
	for _, nodePtrs := range machineToNodesMap {
		for i := range nodePtrs {
			for j := i; j < len(nodePtrs); j++ {
				arcs = append(arcs, Arc{nodePtrs[i],
					nodePtrs[j]})
			}
		}
	}

	neighbours := make(map[Node][]Node)
	for _, edge := range arcs {
		if edge.From != edge.To {
			neighbours[edge.From] = append(neighbours[edge.From], edge.To)
		}
	}

	return Graph{arcs, nodes, neighbours}
}
