package main

import (
	"./io"
	"./graph"
	"fmt"
	"./constants"
)

func main() {
	fmt.Println("Job Shop Scheduling Problem")
	constants.Ants = 2
	constants.WeightFactor = 1.0
	constants.PheromoneFactor = 1.0
	constants.EvaporationRate = 0.01
	constants.InitialPheromone = 1.0
	graph := graph.MakeGraph(io.ReadProblem(9))

	n1, n2 := graph.Nodes[2], graph.NeighbourList[graph.Nodes[2]][0]
	fmt.Println(graph.P(n1,n2))

}
