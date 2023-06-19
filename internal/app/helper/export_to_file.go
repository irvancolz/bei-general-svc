package helper

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type ExportTableToFileProps struct {
	Filename    string
	Data        [][]string
	Headers     []string
	ExcelConfig *ExportToExcelConfig
	PdfConfig   *PdfTableOptions
}

func ExportTableToFile(c *gin.Context, props ExportTableToFileProps) (string, error) {
	fileType := c.Query("file_type")

	if strings.EqualFold("xlsx", fileType) && props.ExcelConfig != nil {
		return props.ExcelConfig.ExportTableToExcel(props.Filename, props.Data)
	}
	if strings.EqualFold("pdf", fileType) && props.PdfConfig != nil {
		return ExportTableToPDF(c, props.Data, props.Filename, *props.PdfConfig)
	}
	if strings.EqualFold("csv", fileType) {
		return ExportTableToCsv(props.Filename, props.Data)
	}
	if strings.EqualFold("txt", fileType) {
		return ExportTableToTxt(props.Filename, props.Data)
	}
	return "", errors.New("unsupported file type")
}

// func ExportTableToFile(c *gin.Context, filename string, data [][]string, excelConfig ExportToExcelConfig) (string, error) {
// 	fileType := c.Query("file_type")
// 	if strings.EqualFold("xlsx", fileType) {
// 		return excelConfig.ExportTableToExcel(filename, data)
// 	}
// 	if strings.EqualFold("csv", fileType) {
// 		return ExportTableToCsv(filename, data)
// 	}
// 	if strings.EqualFold("txt", fileType) {
// 		return ExportTableToTxt(filename, data)
// 	}
// 	return "", errors.New("unsupported file type")
// }
