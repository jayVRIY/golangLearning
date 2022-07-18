package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
)

type CostumeSQLSpot struct {
	value string
}

func main() {
	var (
		SQLTemplate               string
		SQLTemplateAddressPointer string
		//SQLS                      []string
	)
	SQLChosen := make(map[string]any)
	showBanner()                                          //打印banner
	SQLTemplate = getTemplate(&SQLTemplateAddressPointer) //获取SQL模板
	templateSlices, points := sqlFormat(SQLTemplate)      //处理SQL模板
	for _, point := range points {
		SQLChosen[templateSlices[point]] = nil
	}
	printTemplates(templateSlices, points) //展示SQL参数
	rows, _ := getSelectedSheet()          //获取表单
	printSheetTitle(rows)                  //获取选中表单标题
	for len(points) != 0 {
		var (
			selectedSQLSpot int
			selectedRow     int
			chosenType      int
		)
		printTemplates(templateSlices, points)
		println("选择要填入SQL的值:")
		fmt.Scanln(&selectedSQLSpot)
		println(getChonsenTemplate(templateSlices, selectedSQLSpot))
		println("选择输入自定义值或读取表单列:")
		fmt.Scanln(&chosenType)
		if chosenType == 2 {
			printSheetTitle(rows)
			println("选择对应的值对应的表单列")
			fmt.Scanln(&selectedRow)
			SQLChosen[templateSlices[selectedSQLSpot]] = selectedRow
			for i, point := range points {
				if point == selectedSQLSpot {
					points = append(points[:i], points[i+1:]...)
				}
			}
		} else if chosenType == 1 {
			println("输入自定义模板,", getChonsenTemplate(templateSlices, selectedSQLSpot), ":")
			var cosSpot = new(CostumeSQLSpot)
			fmt.Scanln(&cosSpot.value)
			SQLChosen[templateSlices[selectedSQLSpot]] = cosSpot
			for i, point := range points {
				if point == selectedSQLSpot {
					points = append(points[:i], points[i+1:]...)
				}
			}
		}

	}
}

func getChonsenTemplate(slices []string, point int) string {
	return "已选中模板：" + slices[point]
}

func printSheetTitle(rows [][]string) {
	rowOne := rows[0]
	for i, s := range rowOne {
		print(i+1, "、", s, "   ")
	}
	println()
}

func getSelectedSheet() ([][]string, [][]string) {
	var (
		fileName      string
		selectedSheet int
	)
	println("请输入表单文件地址")
	fmt.Scanln(&fileName)
	file, error := excelize.OpenFile(fileName)
	sheepNames := file.GetSheetMap()
	println("请选择要生成sql的表单")
	for i, s := range sheepNames {
		fmt.Printf("表单%d:%s   ", i, s)
	}
	println()
	fmt.Scanln(&selectedSheet)
	rows, error := file.GetRows(sheepNames[selectedSheet])
	cols, error := file.GetCols(sheepNames[selectedSheet])
	if error != nil {
		fmt.Println(error)
	}
	return rows, cols
}

func printTemplates(template []string, points []int) {
	if points != nil {
		println("检测到如下模板")
		for _, point := range points {
			fmt.Printf("%d、%s   ", point, template[point])
		}
	}
	println()
}

func getTemplate(SQLTemplateAddressPointer *string) string {
	println("请输入SQL模板地址：")
	fmt.Scanln(SQLTemplateAddressPointer)
	SQLBytes, err := ioutil.ReadFile(*SQLTemplateAddressPointer)
	if err != nil {
		println(err)
	}
	return string(SQLBytes)
}

func showBanner() {
	fmt.Println("  ________________________________________________________________________________________________________________\n |      _____  ____    __       ______ ______ _   __ ______ ____   ___   ______ ____   ____         ___   ____    |\n |     / ___/ / __ \\  / /      / ____// ____// | / // ____// __ \\ /   | /_  __// __ \\ / __ \\       <  /  / __ \\   |\n |     \\__ \\ / / / / / /      / / __ / __/  /  |/ // __/  / /_/ // /| |  / /  / / / // /_/ /       / /  / / / /   |\n |    ___/ // /_/ / / /___   / /_/ // /___ / /|  // /___ / _, _// ___ | / /  / /_/ // _, _/       / /_ / /_/ /    |\n |   /____/ \\___\\_\\/_____/   \\____//_____//_/ |_//_____//_/ |_|/_/  |_|/_/   \\____//_/ |_|       /_/(_)\\____/     |\n  ________________________________________________________________________________________________________________\n\n  ")
}
