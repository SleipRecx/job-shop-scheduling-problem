package gantt

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"fmt"
	"strconv"
)

func CreateChart(path string, nMachines int) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "B2", 100)
	xlsx.SetColWidth("Sheet1", "B", "ZZ", 3)
	xlsx.SetColWidth("Sheet1", "A", "A", 4)


	for i := 0; i < nMachines; i++ {
		label := "M" + strconv.Itoa(i+1)
		axis := "A" + strconv.Itoa(i+1)
		xlsx.SetCellValue("Sheet1", axis, label)
	}
	err := xlsx.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}

}
