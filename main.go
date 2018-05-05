package main

import (
	"./aco"
	"./constants"
	"./graph"
	"./io"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("Job Shop Scheduling Problem")
	problemFormulation := io.ReadProblem(9)
	constants.InitialPheromone = 0.5
	constants.Ants = 2
	constants.WeightFactor = 1.0
	constants.PheromoneFactor = 1.0
	constants.EvaporationRate = 0.01
	constants.NMachines = problemFormulation.NMachines
	constants.NJobs = problemFormulation.NJobs
	problemGraph := graph.MakeGraph(problemFormulation)
	aco.ACO(problemGraph)

}
