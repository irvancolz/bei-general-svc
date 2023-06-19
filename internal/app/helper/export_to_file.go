package helper

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExportTableToFile(c *gin.Context, filename string, data [][]string, excelConfig ExportToExcelConfig) (string, error) {
	fileType := c.Query("file_type")
	if strings.EqualFold("xlsx", fileType) {
		return excelConfig.ExportTableToExcel(filename, data)
	}
	if strings.EqualFold("csv", fileType) {
		return ExportTableToCsv(filename, data)
	}
	if strings.EqualFold("txt", fileType) {
		return ExportTableToTxt(filename, data)
	}
	return "", errors.New("unsupported file type")
}
