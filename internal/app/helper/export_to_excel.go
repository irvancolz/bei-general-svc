package helper

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportToExcelConfig struct {
	CollumnStart    string
	HeaderStartRow  int
	HeaderHeight    int
	HeaderMarginBot int
	HeaderText      string
}

func (c *ExportToExcelConfig) ExportTableToExcel(filenames string, data [][]string) (string, error) {
	excelFile := excelize.NewFile()
	currentSheet := "Sheet1"

	if len(data) <= 0 {
		log.Println("failed to create excel file: try create excel from empty array")
		return "", errors.New("failed to create excel file: try create excel from empty array")
	}

	collumnStart := strings.ToUpper(string(c.CollumnStart[len(c.CollumnStart)-1]))
	tableEndCol := string([]byte(collumnStart)[len(collumnStart)-1] + byte(len(data[0])-1))
	headerStartRow := c.HeaderStartRow
	if headerStartRow <= 0 {
		headerStartRow = 1
	}
	headerHeight := c.HeaderHeight
	if headerHeight <= 0 {
		headerHeight = 1
	}
	headerMarginBot := c.HeaderMarginBot
	headerEndRow := headerHeight + headerStartRow
	headerStartCell := fmt.Sprintf("%s%v", collumnStart, headerStartRow)
	headerEndCell := fmt.Sprintf("%s%v", tableEndCol, headerEndRow)
	tableStartRow := headerEndRow + headerMarginBot
	tableEndRow := tableStartRow + len(data)
	// cant handle collumn with more than 1 letter e.g AA1
	tableRange := collumnStart + strconv.Itoa(tableStartRow) + ":" + tableEndCol + strconv.Itoa(tableEndRow)
	headerText := c.HeaderText
	filename := filenames

	if headerText == "" {
		headerText = "Bursa Efek Indonesia"
	}

	if filename == "" {
		filename = "BEI_Report"
	}

	errorMerge := excelFile.MergeCell(currentSheet, headerStartCell, headerEndCell)
	if errorMerge != nil {
		log.Println("failed to merge header :", errorMerge)
		return "", errorMerge
	}

	headerStyleId, errorCreateHeaderStyle := excelFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 24,
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if errorCreateHeaderStyle != nil {
		log.Println("failed to create header cell styles :", errorCreateHeaderStyle)
		return "", errorCreateHeaderStyle
	}

	errorAddHeaderStyle := excelFile.SetCellStyle(currentSheet, headerStartCell, headerStartCell, headerStyleId)
	if errorAddHeaderStyle != nil {
		log.Println("failed to add styles to header :", errorAddHeaderStyle)
		return "", errorAddHeaderStyle
	}

	errHeaderTxt := excelFile.SetCellValue(currentSheet, headerStartCell, headerText)
	if errHeaderTxt != nil {
		log.Println("failed to write header text :", errHeaderTxt)
		return "", errHeaderTxt
	}

	errorTable := excelFile.AddTable("Sheet1", &excelize.Table{
		Range:     tableRange,
		Name:      "data",
		StyleName: "TableStyleLight21",
	})
	if errorTable != nil {
		log.Println("failed to create table :", errorTable)
		return "", errorTable
	}

	for rowsIndex, rows := range data {
		for valueNumber, value := range rows {
			// 			set collumn with current data position and collumn start	set rows position with start rows + 1 for header
			cellName := string([]byte(collumnStart)[0]+byte(valueNumber)) + strconv.Itoa(rowsIndex+tableStartRow+1) // A1, B1, dst.
			err := excelFile.SetCellValue("Sheet1", cellName, value)
			if err != nil {
				log.Println("error add data to table :", err)
			}
		}
	}

	result := generateFileNames(filename, "_", time.Now()) + ".xlsx"
	errSave := excelFile.SaveAs(result)
	if errSave != nil {
		log.Println("failed to save excel file:", errSave)
		return "", errSave
	}
	return result, nil
}
