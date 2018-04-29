package graph

import (
	"../io"
	"../constants"
	"math"
)

type Node struct {
	Job, Machine, Time, StartTime int
}

type Arc struct {
	From, To Node
	Weight   int
	Disjunct bool
	Pheromone float64
}

type Graph struct {
	Edges []Arc
	Nodes []Node
	NeighbourList map[Node][]Node
}

// Probability of moving from node i to node j
func (g *Graph) P(i, j Node) float64 {
	if !g.isNeighbour(i,j) {
		return 0.0
	}
	edgeBetween := g.findEdge(i,j)
	over := math.Pow(edgeBetween.Pheromone, constants.PheromoneFactor) * math.Pow(float64(edgeBetween.Weight), constants.WeightFactor)

	under := 0.0
	for x := range g.NeighbourList[i] {
		tempEdge := g.findEdge(i, g.NeighbourList[i][x])
		under += math.Pow(tempEdge.Pheromone, constants.PheromoneFactor) * math.Pow(float64(tempEdge.Weight), constants.WeightFactor)
	}
	return over / under
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
	source := Node{-1, -1, 0, 0}
	sink := Node{-1, -1, 0, 0}
	nodes := []Node{source}
	arcs := make([]Arc, 0)
	machineToNodesMap := make(map[int][]Node)

	// Create conjunctive arcs (technological order)
	for jobId, requirements := range problemFormulation.Sequences {
		previous := source
		for _, requirement := range requirements {
			node := Node{jobId, requirement.Machine, requirement.Time, 0}
			machineToNodesMap[requirement.Machine] = append(machineToNodesMap[requirement.Machine], node)
			arcs = append(arcs, Arc{previous,
				node,
				previous.Time,
				false,
				constants.InitialPheromone})
			nodes = append(nodes, node)
			previous = node
		}
		arcs = append(arcs, Arc{previous,
			sink,
			previous.Time,
			false,
			constants.InitialPheromone})
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
					constants.InitialPheromone})
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
