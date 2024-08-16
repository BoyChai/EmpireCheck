package control

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateExec() {
	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 设置单元格的值

	data := EsFind()
	for i, v := range data {
		f.SetCellValue("Sheet1", fmt.Sprint("A", i+3), v.Name)
		f.SetCellValue("Sheet1", fmt.Sprint("B", i+3), v.Career)
		f.SetCellValue("Sheet1", fmt.Sprint("C", i+3), v.Position)

		// 总繁荣
		prosperous, err := strconv.Atoi(v.Prosperous)
		if err != nil {
			fmt.Println(err)
			continue
		}
		f.SetCellValue("Sheet1", fmt.Sprint("D", i+3), prosperous)
		// 周武勋
		weekMilitaryExploit, err := strconv.Atoi(v.WeekMilitaryExploit)
		if err != nil {
			fmt.Println(err)
			continue
		}
		f.SetCellValue("Sheet1", fmt.Sprint("E", i+3), weekMilitaryExploit)

		// 周贡献
		weekContribute, err := strconv.Atoi(v.WeekContribute)
		if err != nil {
			fmt.Println(err)
			continue
		}
		f.SetCellValue("Sheet1", fmt.Sprint("F", i+3), weekContribute)
		// 今天的繁荣比
		if floatValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(weekMilitaryExploit)/float64(prosperous)), 64); err != nil {
			fmt.Println(err)
			continue
		} else {
			f.SetCellValue("Sheet1", fmt.Sprint("H", i+3), floatValue)

		}

	}
	// 设置工作簿的默认工作表
	// 根据指定路径保存文件
	if err := f.SaveAs("work.xlsx"); err != nil {
		fmt.Println(err)
	}
}
