package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type CostumeSQLSpot struct {
	value []string
}

func main() {
	var (
		SQLTemplate               string
		SQLTemplateAddressPointer string
		SQLS                      []string
	)
	SQLChosen := make(map[string]interface{})
	showBanner()                                          //打印banner
	SQLTemplate = getTemplate(&SQLTemplateAddressPointer) //获取SQL模板
	templateSlices, points := sqlFormat(SQLTemplate)      //处理SQL模板
	for _, point := range points {
		SQLChosen[templateSlices[point]] = nil
	}
	fmt.Printf("检测到%d个模板插槽\n", len(points)) //展示SQL参数
	rows, _ := getSelectedSheet()           //获取表单
	printSheetTitle(rows)                   //获取选中表单标题
	for len(points) != 0 {
		var (
			selectedSQLSpot int
			selectedRow     int
			chosenType      int
			timeOperate     int
		)
		printTemplates(templateSlices, points)
		println("选择要填入SQL的值:")
		fmt.Scanln(&selectedSQLSpot)
		println(getChonsenTemplate(templateSlices, selectedSQLSpot))
		println("选择:1、读取表单列		2、输入自定义值 :")
		fmt.Scanln(&chosenType)
		if chosenType == 1 {
			println("是否操作类型为日期：1、是		2、否")
			var isTime int
			fmt.Scanln(&isTime)
			if isTime == 2 {
				printSheetTitle(rows)
				println("选择对应的值对应的表单列")
				fmt.Scanln(&selectedRow)
				SQLChosen[templateSlices[selectedSQLSpot]] = selectedRow - 1
				points = decrease(points, selectedSQLSpot)
			} else if isTime == 1 {
				c := new(CostumeSQLSpot)
				println("选择对应的值对应的表单列")
				fmt.Scanln(&selectedRow)
				c.value = append(c.value, strconv.Itoa(selectedRow-1))
				println("输入需要添加或减少的月份")
				fmt.Scanln(&timeOperate)
				c.value = append(c.value, "timeOperate:"+strconv.Itoa(timeOperate))
				SQLChosen[templateSlices[selectedSQLSpot]] = c
				points = decrease(points, selectedSQLSpot)
			}
		} else if chosenType == 2 {
			println("输入自定义模板,", getChonsenTemplate(templateSlices, selectedSQLSpot), ":")
			var cosSpot = new(CostumeSQLSpot)
			var costumStr string
			fmt.Scanln(&costumStr)
			cosSpot.value = append(cosSpot.value, costumStr)
			SQLChosen[templateSlices[selectedSQLSpot]] = cosSpot
			points = decrease(points, selectedSQLSpot)
		}

	}
	if len(points) == 0 {
		println("模板填充完成，开始生成sql........")
		for i, row := range rows {
			if i != 0 {
				copyData := make([]string, len(templateSlices))
				copy(copyData, templateSlices)
				for key, value := range SQLChosen {
					for sliceIndex, slice := range templateSlices {
						if slice == key {
							switch value.(type) {
							case *CostumeSQLSpot:
								costumeSQLSpotValue := value.(*CostumeSQLSpot)
								if len(costumeSQLSpotValue.value) == 2 {
									timeRowStr := costumeSQLSpotValue.value[0]
									timeRow, _ := strconv.Atoi(timeRowStr)
									timeValue := row[timeRow]
									timeOperate := costumeSQLSpotValue.value[1]
									if strings.HasPrefix(timeOperate, "timeOperate:") {
										dayOperatorStr := strings.TrimLeft(timeOperate, "timeOperate:")
										dayOperator, _ := strconv.Atoi(dayOperatorStr)
										cosTime, err := time.Parse("2006-01-02 15:04:05", timeValue)
										if err != nil {
											fmt.Println(err)
											return
										}
										cosTime.AddDate(0, dayOperator, 0)
										copyData[sliceIndex] = cosTime.Format("2006-01-02 15:04:05")
										break
									}
								} else {
									copyData[sliceIndex] = costumeSQLSpotValue.value[0]
									break
								}
							case int:
								var rowNum = value.(int)
								rowValue := row[rowNum]
								copyData[sliceIndex] = rowValue
								break
							}
						}
					}
				}
				SQLS = append(SQLS, strings.Join(copyData, ""))
			}
		}
	}
}

func decrease(points []int, selectedSQLSpot int) []int {
	for i, point := range points {
		if point == selectedSQLSpot {
			return append(points[:i], points[i+1:]...)
		}
	}
	return points
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
	println("请输入表单文件地址:")
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
