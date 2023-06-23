package helper

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ExportTableToFileProps struct {
	Filename string
	// the table data
	Data [][]string
	// data showed on the top table collumn
	Headers []string
	// if not specified meant the file format is unsupported
	ExcelConfig *ExportToExcelConfig
	// if not specified meant the file format is unsupported
	PdfConfig *PdfTableOptions
}

// export file to desired file extensions
// make sure there is checking request.ContentLength to be not more than 1, if you use this export feature in get-all func
// if you dont check the content length, it will raise panic if there is more than 1 response written into the response body
func ExportTableToFile(c *gin.Context, props ExportTableToFileProps) error {
	fileType := c.Query("export")
	var filePath string
	var errorPath error

	// skip export if the file type is not specified
	if fileType == "" {
		return nil
	}

	if strings.EqualFold("xlsx", fileType) && props.ExcelConfig != nil {
		filePath, errorPath = props.ExcelConfig.ExportTableToExcel(props.Filename, props.Data)
	} else if strings.EqualFold("pdf", fileType) && props.PdfConfig != nil {
		filePath, errorPath = ExportTableToPDF(c, props.Data, props.Filename, *props.PdfConfig)
	} else if strings.EqualFold("csv", fileType) {
		filePath, errorPath = ExportTableToCsv(props.Filename, props.Data)
	} else if strings.EqualFold("txt", fileType) {
		filePath, errorPath = ExportTableToTxt(props.Filename, props.Data)
	} else {
		return errors.New("unsupported file type")
	}
	if errorPath != nil {
		return errorPath
	}

	c.File(filePath)
	c.Abort()
	errRemoveFile := os.Remove(filePath)
	if errRemoveFile != nil {
		log.Println("failed to clean server disk after create files :", errRemoveFile)
	}
	return nil
}
