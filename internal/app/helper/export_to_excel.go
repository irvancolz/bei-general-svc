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
	HeaderText      []string
	headerStartRow  int
	headerMarginBot int
	currentSheet    string
	data            [][]string
}

func (c *ExportToExcelConfig) getTitleHeight() int {
	return len(c.HeaderText)
}

func (c *ExportToExcelConfig) getTitle() []string {
	if len(c.HeaderText) == 0 {
		return []string{"Bursa Efek Indonesia"}
	}
	return c.HeaderText
}

func (c *ExportToExcelConfig) getTitleRowFrom() int {
	if c.headerStartRow <= 0 {
		return 1
	}
	return c.headerStartRow
}

func (c *ExportToExcelConfig) ExportTableToExcel(filenames string, data [][]string) (string, error) {
	excelFile := excelize.NewFile()
	c.currentSheet = "Sheet1"
	c.data = data
	c.headerStartRow = 2
	c.headerMarginBot = 1

	if len(data) <= 0 {
		log.Println("failed to create excel file: try create excel from empty array")
		return "", errors.New("failed to create excel file: try create excel from empty array")
	}

	if c.CollumnStart == "" {
		c.CollumnStart = "b"
	}

	filename := filenames

	if filename == "" {
		filename = "BEI_Report"
	}

	errDrawTitle := c.generateBasicExcelTitle(excelFile)
	if errDrawTitle != nil {
		return "", errDrawTitle
	}

	errorAddTable := c.Addtable(excelFile)
	if errorAddTable != nil {
		return "", errorAddTable
	}

	result := generateFileNames(filename, "_", time.Now()) + ".xlsx"
	errSave := excelFile.SaveAs(result)
	if errSave != nil {
		log.Println("failed to save excel file:", errSave)
		return "", errSave
	}
	return result, nil
}

func (c *ExportToExcelConfig) generateBasicExcelTitle(excelFile *excelize.File) error {
	if c.CollumnStart == "" {
		c.CollumnStart = "b"
	}

	collumnStart := strings.ToUpper(string(c.CollumnStart[len(c.CollumnStart)-1]))
	headerText := c.getTitle()
	headerStartRow := c.getTitleRowFrom()

	for i, text := range headerText {

		currHeaderCell := fmt.Sprintf("%s%v", collumnStart, headerStartRow+i)

		headerStyleId, errorCreateHeaderStyle := excelFile.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size: func() float64 {
					if i == 0 {
						return 12
					}
					return 11
				}(),
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				WrapText: false,
				Vertical: "center",
			},
		})
		if errorCreateHeaderStyle != nil {
			log.Println("failed to create header cell styles :", errorCreateHeaderStyle)
			return errorCreateHeaderStyle
		}

		errorAddHeaderStyle := excelFile.SetCellStyle(c.currentSheet, currHeaderCell, currHeaderCell, headerStyleId)
		if errorAddHeaderStyle != nil {
			log.Println("failed to add styles to header :", errorAddHeaderStyle)
			return errorAddHeaderStyle
		}

		errHeaderTxt := excelFile.SetCellValue(c.currentSheet, currHeaderCell, text)
		if errHeaderTxt != nil {
			log.Println("failed to write header text :", errHeaderTxt)
			return errHeaderTxt
		}
	}

	return nil
}

func (c *ExportToExcelConfig) Addtable(excelFile *excelize.File) error {
	data := c.data
	collumnStart := strings.ToUpper(string(c.CollumnStart[len(c.CollumnStart)-1]))
	headerStartRow := c.getTitleRowFrom()
	headerHeight := c.getTitleHeight()
	headerMarginBot := c.headerMarginBot
	headerEndRow := headerHeight + headerStartRow

	tableStartRow := headerEndRow + headerMarginBot

	if len(data) <= 0 {
		log.Println("failed to create excel file: try create excel from empty array")
		return errors.New("failed to create excel file: try create excel from empty array")
	}

	headerStyleId, errorCreateHeaderStyle := excelFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:      true,
			Color:     "#ffffff",
			VertAlign: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#9f0e0f"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			createBorder("bottom", "#000000", 1),
			createBorder("top", "#000000", 1),
		},
	})

	contentStyleId, errorCreateContentStyle := excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			createBorder("bottom", "#000000", 1),
		},
	})

	if errorCreateContentStyle != nil {
		return errorCreateContentStyle
	}

	if errorCreateHeaderStyle != nil {
		return errorCreateHeaderStyle
	}

	maxColWdth := getColumnMaxWidth(data)

	for rowsIndex, rows := range data {
		for columnIndex, value := range rows {
			currentCol := string([]byte(collumnStart)[0] + byte(columnIndex))
			cellName := currentCol + strconv.Itoa(rowsIndex+tableStartRow) // A1, B1, dst.

			if rowsIndex == 0 {
				errAddStyle := excelFile.SetCellStyle(c.currentSheet, cellName, cellName, headerStyleId)
				if errAddStyle != nil {
					return errAddStyle
				}
			} else {
				errAddStyle := excelFile.SetCellStyle(c.currentSheet, cellName, cellName, contentStyleId)
				if errAddStyle != nil {
					return errAddStyle
				}
			}

			errorSetWidth := excelFile.SetColWidth(c.currentSheet, currentCol, currentCol, float64(maxColWdth[columnIndex]+4))
			if errorSetWidth != nil {
				log.Println("failed to set collumn width : ", errorSetWidth)
			}

			err := excelFile.SetCellValue("Sheet1", cellName, value)
			if err != nil {
				log.Println("error add data to table :", err)
			}
		}
	}

	return nil
}

func createBorder(borderType, borderColor string, borderStyle int) excelize.Border {
	return excelize.Border{
		Type:  borderType,
		Color: borderColor,
		Style: borderStyle,
	}
}
