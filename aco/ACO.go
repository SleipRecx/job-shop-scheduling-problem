package aco

import (
	"../graph"
	"../constants"
	"../util"
	"../gantt"
	"fmt"
	"math"
	"math/rand"
	"go/constant"
)

func removeFromList(list []graph.Node, element graph.Node) []graph.Node {
	newList := make([]graph.Node, 0)
	for _, node := range list {
		if !(node.Job == element.Job && node.Machine == element.Machine) {
			newList = append(newList, node)
		}
	}
	return newList
}

type Solution struct {
	StartTimeMap map[graph.Node] int
	Nodes []graph.Node
}

func listScheduler(problemGraph graph.Graph) Solution{
	startTimeMap := make(map[graph.Node]int)
	partialSolution := make([]graph.Node, 0)
	unvisited := make([]graph.Node, len(problemGraph.Nodes))
	copy(unvisited, problemGraph.Nodes)
	for range problemGraph.Nodes {
		unvisitedPrime := restrict(partialSolution, unvisited)
		nodeStar := chooseRandom(unvisitedPrime)
		startTimeMap[nodeStar] = earliestStartTime(nodeStar, partialSolution)
		partialSolution = append(partialSolution, nodeStar)
		unvisited = removeFromList(unvisited, nodeStar)
	}
	return Solution{startTimeMap,partialSolution}
}

func calculateMakespan(solution Solution) int {
	max := 0
	for x := range solution.Nodes {
		if time := solution.StartTimeMap[solution.Nodes[x]] + solution.Nodes[x].Time; time > max {
			max = time
		}
	}
	return max
}


func earliestCompletionTime(node graph.Node, ps Solution) int {
	if len(ps.Nodes) == 0 {
		return node.Time
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range ps.Nodes {
		if node.Job == ps.Nodes[x].Job && ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time > JobTimer {
			JobTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
		if node.Machine == ps.Nodes[x].Machine && ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time > MachineTimer {
			MachineTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
	}
	earliestComp := util.Max(JobTimer, MachineTimer + node.Time)
	if earliestComp == -1 {
		return node.Time
	}
	return earliestComp
}

func earliestStartTime(node graph.Node, ps Solution) int {
	if len(ps.Nodes) == 0 {
		return 0
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range ps.Nodes {
		if node.Job == ps.Nodes[x].Job && ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time > JobTimer {
			JobTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
		if node.Machine == ps.Nodes[x].Machine && ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time > MachineTimer {
			MachineTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
	}
	earliestComp := util.Max(JobTimer, MachineTimer)
	return util.Max(earliestComp, 0)
}

func restrict(partialSolution []graph.Node, unVisited []graph.Node) []graph.Node {
	tStar := math.MaxInt32
	restrictedSet := make([]graph.Node, 0)
	for x := range unVisited {
		if t := earliestCompletionTime(unVisited[x], partialSolution); t  < tStar && preStepsExecuted(unVisited[x], partialSolution){
			tStar = t
		}
	}
	for x := range unVisited {
		if earliestStartTime(unVisited[x], partialSolution) <= tStar && preStepsExecuted(unVisited[x],partialSolution) {
			restrictedSet = append(restrictedSet, unVisited[x])
		}
	}
	return restrictedSet
}

func preStepsExecuted(node graph.Node, partialSolution []graph.Node) bool {
	counter := 0
	for x := range partialSolution {
		if partialSolution[x].Job == node.Job {
			counter++
		}
	}
	if counter == node.TechStep {
		return true
	}
	return false
}

func chooseRandom(candidates []graph.Node) graph.Node {
	return candidates[rand.Intn(len(candidates))]
}

func numberOfAnts(problemGraph graph.Graph) int {
	return util.Max(10, len(problemGraph.Nodes) / 10)
}

func InitializePheromoneValues(graph graph.Graph) map[graph.Arc]float64{
	mapping := make(map[graph.Arc]float64)
	for _, arc := range graph.Edges {
		mapping[arc] = constants.InitialPheromone
	}
	return mapping
}

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	arcPheroMap := InitializePheromoneValues(problemGraph)
	var iterationBest []graph.Node		//Sib
	var bestSoFar []graph.Node			//Sbs
	var restartBest []graph.Node		//Srb
	convergenceFactor := 0.0			// cf
	bsUpdate := false
	numberOfAnts := numberOfAnts(problemGraph)
	for true {
		solutions := make([]Solution, 0)
		for i := 0; i < numberOfAnts; i++ {
			solutions = append(solutions, listScheduler(problemGraph))
		}
		ApplyLocalSearch(solutions)
		iterationBest = SolutionWithMinimalMakeSpan(solutions)
		EliteAction(iterationBest)
		Update(iterationBest, restartBest, bestSoFar)
		ApplyPheromoneUpdate(convergenceFactor, bsUpdate, T, restartBest, bestSoFar)
		if convergenceFactor > 0.99 {
			if bsUpdate {
				ResetPheromoneValues(T)
				restartBest = nil
				bsUpdate = false
			} else {
				bsUpdate = true
			}
		}
	}
	fmt.Println("Done!")
	fmt.Println("Best makespan:", calculateMakespan(bestSoFar))
	orders := graph.NodeListToOrderList(bestSoFar)
	gantt.CreateChart("03 - Program Outputs/Chart.xlsx", orders)


}
