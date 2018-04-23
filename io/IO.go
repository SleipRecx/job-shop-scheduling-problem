package io

import (
	"os"
	"strconv"
	"bufio"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Requirement struct {
	Machine, Time int
}

type ProblemFormulation struct {
	NJobs, NMachines int
	Sequences map[int][]Requirement

}

func ReadProblem(n int) ProblemFormulation{
	infile, err := os.Open("./02 - Test Data/" + strconv.Itoa(n) + ".txt")
	defer infile.Close()
	check(err)
	scanner := bufio.NewScanner(infile)
	problemFormulation := ProblemFormulation{}
	problemFormulation.Sequences = make(map[int][]Requirement)
	jobId := 0
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		if len(data) > 0 {
			if len(data) == 2 {
				problemFormulation.NJobs, err = strconv.Atoi(data[0])
				check(err)
				problemFormulation.NMachines, err = strconv.Atoi(data[1])
				check(err)
			} else {
				requirements := make([]Requirement, 0)
				for i := 0; i < len(data); i+=2 {
					machine, err := strconv.Atoi(data[i])
					check(err)
					time, err := strconv.Atoi(data[i+1])
					check(err)
					requirements = append(requirements, Requirement{machine, time})
				}
				problemFormulation.Sequences[jobId] = requirements
				jobId++
			}
		}

	}
	return problemFormulation
}