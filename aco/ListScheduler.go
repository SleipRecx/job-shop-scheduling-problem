package aco

import (
	"../graph"
	"../util"
	"math"
	"math/rand"
	"../constants"
)

func listScheduler(problemGraph graph.Graph, arcPheroMap map[graph.Arc]float64) Solution {
	startTimeMap := make(map[graph.Node]int)
	partialSolution := make([]graph.Node, 0)
	unvisited := make([]graph.Node, len(problemGraph.Nodes))
	copy(unvisited, problemGraph.Nodes)
	for range problemGraph.Nodes {
		unvisitedPrime := restrict(Solution{startTimeMap, partialSolution}, unvisited)
		nodeStar := choose(unvisitedPrime, unvisited, arcPheroMap, Solution{startTimeMap, partialSolution})
		startTimeMap[nodeStar] = earliestStartTime(nodeStar, Solution{startTimeMap, partialSolution})
		partialSolution = append(partialSolution, nodeStar)
		unvisited = removeFromList(unvisited, nodeStar)
	}
	return Solution{startTimeMap, partialSolution}
}

func earliestCompletionTime(node graph.Node, ps Solution) int {
	if len(ps.Nodes) == 0 {
		return node.Time
	}
	MachineTimer := -1
	JobTimer := -1
	for x := range ps.Nodes {
		if node.Job == ps.Nodes[x].Job && ps.StartTimeMap[ps.Nodes[x]]+ps.Nodes[x].Time > JobTimer {
			JobTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
		if node.Machine == ps.Nodes[x].Machine && ps.StartTimeMap[ps.Nodes[x]]+ps.Nodes[x].Time > MachineTimer {
			MachineTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
	}
	earliestComp := util.Max(JobTimer, MachineTimer+node.Time)
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
		if node.Job == ps.Nodes[x].Job && ps.StartTimeMap[ps.Nodes[x]]+ps.Nodes[x].Time > JobTimer {
			JobTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
		if node.Machine == ps.Nodes[x].Machine && ps.StartTimeMap[ps.Nodes[x]]+ps.Nodes[x].Time > MachineTimer {
			MachineTimer = ps.StartTimeMap[ps.Nodes[x]] + ps.Nodes[x].Time
		}
	}
	earliestComp := util.Max(JobTimer, MachineTimer)
	return util.Max(earliestComp, 0)
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


func restrict(partialSolution Solution, unVisited []graph.Node) []graph.Node {
	tStar := math.MaxInt32
	restrictedSet := make([]graph.Node, 0)
	for x := range unVisited {
		if t := earliestCompletionTime(unVisited[x], partialSolution); t < tStar && preStepsExecuted(unVisited[x], partialSolution.Nodes) {
			tStar = t
		}
	}
	for x := range unVisited {
		if earliestStartTime(unVisited[x], partialSolution) <= tStar && preStepsExecuted(unVisited[x], partialSolution.Nodes) {
			restrictedSet = append(restrictedSet, unVisited[x])
		}
	}
	return restrictedSet
}

func chooseRandom(candidates []graph.Node) graph.Node {
	return candidates[rand.Intn(len(candidates))]
}

// Tror denne er good
func eta(n graph.Node, ps Solution, candidates []graph.Node) float64 {
	numerator := 1.0 / float64(earliestStartTime(n, ps)) + 1.0
	denominator := 0.0
	for _, k := range candidates {
		denominator += 1.0 / (float64(earliestStartTime(k, ps)) + 1.0)
	}
	return numerator / denominator
}

func choose(candidates, unvisited []graph.Node, arcPheroMap map[graph.Arc]float64, ps Solution) graph.Node {
	probabilities := make(map[graph.Node]float64)

	// both 0 and +Inf
	denominator := 0.0
	for _, n := range candidates {
		min := math.MaxFloat64
		for _, u := range unvisited {
			if n.Machine == u.Machine && u != n {
				v := arcPheroMap[graph.Arc{n, u}] * math.Pow(eta(n, ps, candidates), constants.Beta)
				if v < min {
					min = v
				}
			}
		}
		if min != math.MaxFloat64 {
			denominator += min
		}
	}
	for _, n := range candidates {
		intersection := make([]graph.Node, 0)
		for _, u := range unvisited {
			if n.Machine == u.Machine && n != u{
				intersection = append(intersection, u)
			}
		}
		numerator := math.MaxFloat64
		for j := range intersection {
			v := arcPheroMap[graph.Arc{n, intersection[j]}] * math.Pow(eta(n, ps, candidates), constants.Beta)
			if v < numerator {
				numerator = v
			}
		}
		if numerator == math.MaxFloat64 {
			numerator = 0
		}
		if denominator == 0 {
			probabilities[n] = 0
		} else {
			probabilities[n] = numerator / denominator
		}
	}
	//fmt.Println(probabilities)
	return chooseRandom(candidates)
}