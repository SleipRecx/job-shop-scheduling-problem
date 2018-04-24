package main

import (
	"./io"
	"./graph"
	"fmt"
)

func main() {

	fmt.Println("Job Shop Scheduling Problem")
	graph.MakeGraph(io.ReadProblem(9))
}
