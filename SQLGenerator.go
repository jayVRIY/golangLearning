package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	var (
		fileName      string
		selectedSheet int
		//sqlTemplate        string
		sqlTemplateAddress string
	)
	fmt.Println("  ________________________________________________________________________________________________________________\n |      _____  ____    __       ______ ______ _   __ ______ ____   ___   ______ ____   ____         ___   ____    |\n |     / ___/ / __ \\  / /      / ____// ____// | / // ____// __ \\ /   | /_  __// __ \\ / __ \\       <  /  / __ \\   |\n |     \\__ \\ / / / / / /      / / __ / __/  /  |/ // __/  / /_/ // /| |  / /  / / / // /_/ /       / /  / / / /   |\n |    ___/ // /_/ / / /___   / /_/ // /___ / /|  // /___ / _, _// ___ | / /  / /_/ // _, _/       / /_ / /_/ /    |\n |   /____/ \\___\\_\\/_____/   \\____//_____//_/ |_//_____//_/ |_|/_/  |_|/_/   \\____//_/ |_|       /_/(_)\\____/     |\n  ________________________________________________________________________________________________________________\n\n  ")
	println("请输入SQL模板地址：")
	fmt.Scanln(&sqlTemplateAddress)
	templete, points := sqlFormat(" INSERT INTO \"CIF_RES_FILE_BATCH\" (\"CLIENT_NO\",\"CH_CLIENT_NAME\",\"DOCUMENT_TYPE\",\"DOCUMENT_ID\",\"MSG_TYPE_EXT\",\"STATUS\",\"ETL_DATE\") VALUES ('${CLIENT_NO}','${CH_CLIENT_NAME}','${DOCUMENT_TYPE}','${DOCUMENT_ID}','${MSG_TYPE_EXT}','0','20220522');")
	if points != nil {
		println("检测到如下模板")
		for index, point := range points {
			fmt.Printf("%d、%s\t", index+11, templete[point])
		}
	}
	println("请输入文件地址")
	fmt.Scanln(&fileName)
	file, error := excelize.OpenFile(fileName)
	if error != nil {
		fmt.Println(error)
	}
	sheepNames := file.GetSheetMap()
	println("请选择要生成sql的表单")
	for i, s := range sheepNames {
		fmt.Printf("表单%d:%s\t", i, s)
	}
	println()
	fmt.Scanln(&selectedSheet)
	rows, error := file.GetRows(sheepNames[selectedSheet])
	if error != nil {
		println(error)
	}
	rowOne := rows[0]
	for i, s := range rowOne {
		print(i+1, "、", s, "\t")
	}
}
