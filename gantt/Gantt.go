package gantt

import (
	"../constants"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lucasb-eyer/go-colorful"
	"strconv"
)

// TODO: Move and change to represent actual job etc.
type Order struct {
	JobID     int
	MachineID int
	Time      int
	Duration  int
}

func CreateChart(path string, orders []Order) {
	xlsx := excelize.NewFile()
	xlsx.SetColWidth("Sheet1", "B", "ZZZ", 3.5)
	xlsx.SetColWidth("Sheet1", "A", "A", 5)

	machineStyle, err := xlsx.NewStyle(`{"alignment":{"horizontal": "center"}, 
                                               "fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1}}`)
	if err != nil {
		fmt.Println(err)
	}
	// Machine labels
	for i := 0; i < constants.NMachines; i++ {
		label := "M" + strconv.Itoa(i)
		cell := "A" + strconv.Itoa(i+1)
		xlsx.SetCellValue("Sheet1", cell, label)
		xlsx.SetCellStyle("Sheet1", cell, cell, machineStyle)
	}
	// Time numbering
	for i := 0; i < 2000; i++ {
		row := strconv.Itoa(constants.NMachines + 1)
		col := timeToExcelCol(i)
		if i%5 == 0 {
			xlsx.SetCellValue("Sheet1", col+row, i)
		}
	}
	// Colors
	colorMap := make(map[int]string)
	palette := colorful.FastHappyPalette(constants.NJobs)
	for i, color := range palette {
		colorMap[i] = color.Hex()
	}
	// Color Legend
	legendRow := constants.NMachines + 2
	for i := 0; i < constants.NJobs; i++ {
		row := strconv.Itoa(legendRow + i)
		col := "A"
		cellStyle, err := xlsx.NewStyle(generateBGStyle(colorMap[i]))
		if err == nil {
			xlsx.SetCellValue("Sheet1", col+row, "Job "+strconv.Itoa(i))
			xlsx.SetCellStyle("Sheet1", col+row, col+row, cellStyle)
		}
	}

	// Plot all orders
	for _, order := range orders {
		addJobToExcel(xlsx, colorMap, order)
	}
	err = xlsx.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}

}

func addJobToExcel(xlsx *excelize.File, colorMap map[int]string, order Order) {
	row := strconv.Itoa(order.MachineID + 1)
	cellStyle, err := xlsx.NewStyle(generateBGStyle(colorMap[order.JobID]))
	shouldAddLabel := true
	if err == nil {
		for i := 0; i < order.Duration; i++ {
			col := timeToExcelCol(order.Time + i)
			if shouldAddLabel {
				xlsx.SetCellValue("Sheet1", col+row, order.Duration)
				shouldAddLabel = false
			}
			xlsx.SetCellStyle("Sheet1", col+row, col+row, cellStyle)
		}
	}
}

func generateBGStyle(color string) string {
	return `{"fill":{"type":"pattern","color":["` + color + `"],"pattern":1}}`
}

func timeToExcelCol(time int) string {
	dividend := time + 2 // +2 because time starts at 0 and col A is machine label
	col := ""
	modulo := 0
	for dividend > 0 {
		modulo = (dividend - 1) % 26
		col = string(rune(65+modulo)) + col
		dividend = (int)((dividend - modulo) / 26)
	}
	return col
}
