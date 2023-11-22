package log_system

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/log_system"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(c *gin.Context) (*helper.PaginationResponse, error)
	CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error)
	ExportLogSystem(c *gin.Context) error
}

type usecase struct {
	logSystemRepo log_system.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		log_system.NewRepository(),
	}
}

func (m *usecase) GetAll(c *gin.Context) (*helper.PaginationResponse, error) {
	paginationData, err := m.logSystemRepo.GetAllWithFilterPagination(c)
	if err != nil {
		return nil, err
	}

	columnHeaders := []string{"Modul", "Sub Modul", "Aksi", "Detail", "User", "IP", "Waktu"}
	columnWidth := []float64{35, 40, 30, 40, 45, 35, 40}

	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	exportedFields := []string{"modul", "sub_modul", "action", "detail", "username", "ip", "created_at"}
	var exportedData [][]string

	for _, content := range paginationData.Data {
		var item []string
		item = append(item, helper.MapToArray(content, exportedFields)...)

		for i, content := range item {
			if i == 6 {
				date, _ := time.Parse("2006-01-02 15:04:05", content[0:19])

				item[i] = date.Format("2 Jan 2006 - 15:04")
			}
		}

		exportedData = append(exportedData, item)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename: "Log-System",
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Log-System"},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:  400,
			Papperheight: 210,
			HeaderRows:   helper.GenerateTableHeaders(columnHeaders, columnWidth),
		},
		Data:        exportedData,
		Headers:     tablesColumns,
		ColumnWidth: columnWidthInt,
	}

	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return nil, errorExport
	}

	return paginationData, nil
}

func (m *usecase) CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error) {
	return m.logSystemRepo.CreateLogSystem(log, c)
}

func (m *usecase) ExportLogSystem(c *gin.Context) error {
	results, errorResults := m.logSystemRepo.GetAll(c)
	if errorResults != nil {
		return errorResults
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, _ := helper.HandleDataFiltering(c, dataToConverted, nil)

	columnHeaders := []string{"Modul", "Sub Modul", "Aksi", "Detail", "User", "IP", "Waktu"}
	columnWidth := []float64{35, 40, 30, 40, 45, 35, 40}

	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	exportedFields := []string{"modul", "submodul", "action", "detail", "username", "ip", "createdat"}
	var exportedData [][]string

	for _, content := range filteredData {
		var item []string
		item = append(item, helper.MapToArray(content, exportedFields)...)

		for i, content := range item {
			if i == 6 {
				date, _ := time.Parse("2006-01-02 15:04:05", content[0:19])

				item[i] = date.Format("2 Jan 2006 - 15:04")
			}
		}

		exportedData = append(exportedData, item)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename: "Log System",
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Log System"},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:  297,
			Papperheight: 210,
			HeaderRows:   helper.GenerateTableHeaders(columnHeaders, columnWidth),
		},
		Data:        exportedData,
		Headers:     tablesColumns,
		ColumnWidth: columnWidthInt,
	}

	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return errorExport
	}

	return nil
}
