package main

import ("fmt"
"./gantt")

func main() {

	fmt.Println("Job Shop Scheduling Problem")
	gantt.CreateChart("./03 - Program Outputs/Chart.xlsx", 9)

}
