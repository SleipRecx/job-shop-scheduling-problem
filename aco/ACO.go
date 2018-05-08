package aco

import (
	"../constants"
	"../graph"
	"../util"
	"fmt"
	"math"
	"../jssp"
	"errors"
	"../gantt"
	"strconv"
)




func numberOfAnts(problemGraph graph.Graph) int {
	return util.Max(10, len(problemGraph.Nodes)/10)
}

func InitializePheromoneValues(problemGraph graph.Graph) map[graph.Arc]float64 {
	mapping := make(map[graph.Arc]float64)
	for _, arc := range problemGraph.Edges {
		mapping[arc] = constants.InitialPheromone
	}
	return mapping
}


func Update(original, candidate jssp.Solution) jssp.Solution {
	if len(original.Nodes) == 0 {
		return  candidate
	} else if jssp.CalculateMakespan(original) < jssp.CalculateMakespan(candidate) {
		return original
	}
	return candidate
}

func index(n graph.Node, nodes []graph.Node) (int, error) {
	for i := range nodes {
		if nodes[i] == n {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func δ(n1, n2 graph.Node, solution jssp.Solution) float64{
	n1Index, _ := index(n1, solution.Nodes)
	n2Index, _ := index(n2, solution.Nodes)
	if n1Index < n2Index {
		return 1
	}
	return 0
}

func fmmaas(x float64) float64 {
	if x < constants.TMin {
		return constants.TMin
	}
	if x > constants.TMax {
		return constants.TMax
	}
	return x
}

func ApplyPheromoneUpdate(cf float64, bsUpdate bool, arcPheroMap map[graph.Arc]float64, restartBest, bestSoFar jssp.Solution) (map[graph.Arc]float64, float64) {
	solutionToUse := restartBest
	if bsUpdate {
		solutionToUse = bestSoFar
	}
	cfUpper := 0.0
	cfLower := float64(len(arcPheroMap)) * (constants.TMax - constants.TMin)
	for arc, tij := range arcPheroMap {
		arcPheroMap[arc] = fmmaas(tij + constants.EvaporationRate * (δ(arc.From, arc.To, solutionToUse)-tij))
		cfUpper += math.Max(constants.TMax - tij, tij - constants.TMin)
	}
	return arcPheroMap, 2 * ((cfUpper / cfLower) - 0.5)
}

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	arcPheroMap := InitializePheromoneValues(problemGraph)
	var iterationBest jssp.Solution		//Sib
	var bestSoFar jssp.Solution			//Sbs
	var restartBest jssp.Solution		//Srb
	convergenceFactor := 0.0			// cf
	bsUpdate := false
	numberOfAnts := numberOfAnts(problemGraph)
	iterationCount := 0
	for !jssp.IsDone(bestSoFar) {
		iterationCount += 1
		fmt.Println("Iteration", iterationCount, "Best makespan", jssp.CalculateMakespan(bestSoFar))

		solutions := make([]jssp.Solution, 0)
		for i := 0; i < numberOfAnts; i++ {
			solutions = append(solutions, jssp.ListScheduler(problemGraph, arcPheroMap, false))
		}
		for i := range solutions {
			solutions[i] = jssp.ApplyLocalSearch(solutions[i])
		}
		iterationBest = jssp.SolutionWithMinimalMakeSpan(solutions)
		//EliteAction(iterationBest)
		bestSoFar = Update(bestSoFar, iterationBest)
		restartBest = Update(restartBest, iterationBest)

		arcPheroMap, convergenceFactor = ApplyPheromoneUpdate(convergenceFactor, bsUpdate, arcPheroMap, restartBest, bestSoFar)
		if convergenceFactor > 0.99 {
			if bsUpdate {
				arcPheroMap = InitializePheromoneValues(problemGraph)
				restartBest = jssp.Solution{}
				bsUpdate = false
			} else {
				bsUpdate = true
			}
		}
	}

	fmt.Println("Done!")
	fmt.Println("Makespan:", jssp.CalculateMakespan(bestSoFar))
	orders := graph.NodeListToOrderList(bestSoFar.Nodes, bestSoFar.StartTimeMap)
	gantt.CreateChart("03 - Program Outputs/ACO_Chart_" + strconv.Itoa(constants.ProblemNumber) + ".xlsx", orders)

}
