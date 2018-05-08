package main

import (
	"./constants"
	"./graph"
	"./io"
	"./ba"
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
	constants.Scouts = 50
	constants.NBest = 5
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

	constants.ProblemNumber = 3
	problemFormulation := io.ReadProblem(constants.ProblemNumber)
	constants.NMachines = problemFormulation.NMachines
	constants.NJobs = problemFormulation.NJobs
	problemGraph := graph.MakeGraph(problemFormulation)
	//aco.ACO(problemGraph)
	ba.BA(problemGraph)
}
