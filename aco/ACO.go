package aco

import (
	"../graph"
	_ "../constants"
	"../gantt"
	"../util"
	"fmt"
	"math"
	"math/rand"
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


func listScheduler(problemGraph graph.Graph) []graph.Node{
	partialSolution := make([]graph.Node, 0)
	unvisited := make([]graph.Node, len(problemGraph.Nodes))
	copy(unvisited, problemGraph.Nodes)
	for range problemGraph.Nodes {
		unvisitedPrime := restrict(partialSolution, unvisited)
		nodeStar := chooseRandom(unvisitedPrime)
		nodeStar.StartTime = earliestStartTime(nodeStar, partialSolution)
		partialSolution = append(partialSolution, nodeStar)
		unvisited = removeFromList(unvisited, nodeStar)
	}
	return partialSolution
}

func calculateMakespan(solution []graph.Node) int {
	max := 0
	for x := range solution {
		if time := solution[x].StartTime + solution[x].Time; time > max {
			max = time
		}
	}
	return max
}


func earliestCompletionTime(node graph.Node, partialSolution []graph.Node) int {
	if len(partialSolution) == 0 {
		return node.Time
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range partialSolution {
		if node.Job == partialSolution[x].Job && partialSolution[x].StartTime + partialSolution[x].Time > JobTimer {
			JobTimer = partialSolution[x].StartTime + partialSolution[x].Time
		}
		if node.Machine == partialSolution[x].Machine && partialSolution[x].StartTime + partialSolution[x].Time > MachineTimer {
			MachineTimer = partialSolution[x].StartTime + partialSolution[x].Time
		}
	}
	earliestComp := math.Max(float64(JobTimer),float64(MachineTimer)) + float64(node.Time)
	if earliestComp == -1 {
		return node.Time
	}
	return int(earliestComp)
}

func earliestStartTime(node graph.Node, partialSolution []graph.Node) int {
	if len(partialSolution) == 0 {
		return 0
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range partialSolution {
		if node.Job == partialSolution[x].Job && partialSolution[x].StartTime + partialSolution[x].Time > JobTimer {
			JobTimer = partialSolution[x].StartTime + partialSolution[x].Time
		}
		if node.Machine == partialSolution[x].Machine && partialSolution[x].StartTime + partialSolution[x].Time > MachineTimer {
			MachineTimer = partialSolution[x].StartTime + partialSolution[x].Time
		}
	}
	earliestComp := math.Max(float64(JobTimer),float64(MachineTimer))
	return int(math.Max(earliestComp, 0))
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

func ACO(problemGraph graph.Graph) {
	fmt.Println("Running ACO")
	best := math.MaxInt32
	var bestSol []graph.Node
	for i := 0; i < 10000; i++ {
		solution := listScheduler(problemGraph)
		//fmt.Println("Solution: ", solution)
		makeSpan := calculateMakespan(solution)
		//fmt.Println("Makespan: ", makeSpan)
		if makeSpan < best {
			best = makeSpan
			bestSol = solution
		}

	}
	fmt.Println("-------")
	fmt.Println("Best makespan:", best)
	orders := graph.NodeListToOrderList(bestSol)
	gantt.CreateChart("03 - Program Outputs/Chart.xlsx", orders)
}
