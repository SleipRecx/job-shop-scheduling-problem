package gantt

import (
	"fmt"
	"math"
	"../util"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

type DataPoint struct {
	Iteration int
	MakeSpans []int
}

func getMaxMin(nums []int) (int, int) {
	min := math.MaxInt32
	max := math.MinInt32
	for _, n := range nums {
		min = util.Min(min, n)
		max = util.Max(max, n)
	}
	return min, max
}

func mean(nums []int) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += float64(n)
	}
	return sum / float64(len(nums))
}

func std(nums []int) float64 {
	mean := mean(nums)
	stdSum := 0.0
	for _, n := range nums {
		stdSum += math.Pow(float64(n) - mean, 2.0)
	}
	stdSum = stdSum / float64(len(nums))
	return math.Sqrt(stdSum)
}


func CreateSummaryGraph(path string, dataPoints []DataPoint) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A" + strconv.Itoa(1), "Iteration")
	xlsx.SetCellValue("Sheet1", "B" + strconv.Itoa(1), "Min")
	xlsx.SetCellValue("Sheet1", "C" + strconv.Itoa(1), "Max")
	xlsx.SetCellValue("Sheet1", "D" + strconv.Itoa(1), "Mean")
	row := 2
	for _, d := range dataPoints {
		min, max := getMaxMin(d.MakeSpans)
		xlsx.SetCellValue("Sheet1", "A" + strconv.Itoa(row), d.Iteration)
		xlsx.SetCellValue("Sheet1", "B" + strconv.Itoa(row), min)
		xlsx.SetCellValue("Sheet1", "C" + strconv.Itoa(row), max)
		xlsx.SetCellValue("Sheet1", "D" + strconv.Itoa(row), mean(d.MakeSpans))
		row += 1
	}

	err := xlsx.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
}