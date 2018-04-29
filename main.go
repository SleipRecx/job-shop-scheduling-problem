package main

import (
	"./io"
	"./graph"
	"fmt"
	"./constants"
	"./aco"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("Job Shop Scheduling Problem")
	constants.Ants = 2
	constants.WeightFactor = 1.0
	constants.PheromoneFactor = 1.0
	constants.EvaporationRate = 0.01
	constants.InitialPheromone = 1.0
	problemGraph := graph.MakeGraph(io.ReadProblem(1))
	aco.ACO(problemGraph)

}
