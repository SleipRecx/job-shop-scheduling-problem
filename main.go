package main

import (
	"./constants"
	"./graph"
	"./io"
	"./ba"
	"./aco"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("Job Shop Scheduling Problem")
	constants.InitialPheromone = 0.5
	constants.Ants = 2
	constants.WeightFactor = 1.0
	constants.PheromoneFactor = 1.0
	constants.EvaporationRate = 0.01
	constants.TMax = 0.999
	constants.TMin = 0.001
	constants.Beta = 10
	constants.Scouts = 10
	constants.NBest = 4
	constants.NRB = 1
	constants.NElites = 2
	constants.NRE = 2
	constants.TargetValues = map[int]int{
		1: 56,
		2: 1059,
		3: 1276,
		4: 1130,
		5: 1451,
		6: 979,
	}
	constants.Iterations = 50
	constants.ProblemNumber = 3
	problemFormulation := io.ReadProblem(constants.ProblemNumber)
	constants.NMachines = problemFormulation.NMachines
	constants.NJobs = problemFormulation.NJobs
	problemGraph := graph.MakeGraph(problemFormulation)
	//aco.ACO(problemGraph)
	start := time.Now()
	ba.BA(problemGraph)
	fmt.Println("Time:", time.Now().Sub(start))
	start = time.Now()
	aco.ACO(problemGraph)
	fmt.Println("Time:", time.Now().Sub(start))
}
