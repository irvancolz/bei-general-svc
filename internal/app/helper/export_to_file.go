package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
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
	// data showed on the top table collumn for txt files
	Headers [][]string
	// specify each co width
	ColumnWidth []int
	// if not specified meant the file format is unsupported
	ExcelConfig *ExportToExcelConfig
	// if not specified meant the file format is unsupported
	PdfConfig *PdfTableOptions
	// if not specified meant the file format is unsupported
	TxtConfig bool
	// if not specified meant the file format is unsupported
	CsvConfig bool
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
		var data [][]string
		data = append(data, props.Headers...)
		data = append(data, props.Data...)
		filePath, errorPath = props.ExcelConfig.ExportTableToExcel(props.Filename, data)
	} else if strings.EqualFold("pdf", fileType) && props.PdfConfig != nil {
		filePath, errorPath = ExportTableToPDF(c, props.Data, props.Filename, props.PdfConfig)
	} else if strings.EqualFold("csv", fileType) {
		var data [][]string
		data = append(data, props.Headers...)
		data = append(data, props.Data...)
		filePath, errorPath = ExportTableToCsv(props.Filename, data)
	} else if strings.EqualFold("txt", fileType) {
		txtConfig := ExportTableToTxtProps{
			Filename:    props.Filename,
			Data:        props.Data,
			Header:      props.Headers,
			ColumnWidth: props.ColumnWidth,
		}
		filePath, errorPath = ExportTableToTxt(txtConfig)
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

type ExportAnnouncementsToFileProps struct {
	Filename string
	// the table data
	Data model.Announcement
	// if not specified meant the file format is unsupported
	ExcelConfig *ExportToExcelConfig
	// if not specified meant the file format is unsupported
	PdfConfig PdfTableOptions
}

// export file to desired file extensions
// make sure there is checking request.ContentLength to be not more than 1, if you use this export feature in get-all func
// if you dont check the content length, it will raise panic if there is more than 1 response written into the response body
func ExportAnnouncementsToFile(c *gin.Context, props ExportAnnouncementsToFileProps) error {
	fileType := c.Query("export")
	var filePath string
	var errorPath error

	// skip export if the file type is not specified
	if fileType == "" {
		return nil
	}

	if strings.EqualFold("pdf", fileType) {
		filePath, errorPath = ExportAnnouncementToPdf(c, props.Data, props.PdfConfig, props.Filename)
	} else if strings.EqualFold("xlsx", fileType) {
		filePath, errorPath = ExportAnnouncementsToExcel(props.Filename, props.Data)
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
