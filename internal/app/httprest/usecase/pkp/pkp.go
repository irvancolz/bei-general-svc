package pkp

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	pkp "be-idx-tsg/internal/app/httprest/repository/pkp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAllPKuser(c *gin.Context) (*helper.PaginationResponse, error)
	CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error)
	UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
}

type usecase struct {
	pkpRepo pkp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		pkp.NewRepository(),
	}
}

func (uc *usecase) GetAllPKuser(c *gin.Context) (*helper.PaginationResponse, error) {
	results, errorResults := uc.pkpRepo.GetAllPKuser(c)
	if errorResults != nil {
		return nil, errorResults
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"createdat", "updatedat", "questiondate", "answersat", "deletedat"})
	sortedData := helper.HandleDataSorting(c, filteredData)
	exportedFields := []string{
		"stakeholders",
		"code",
		"name",
		"questiondate",
		"question",
		"answersat",
		"answers",
		"topic",
		"answersby",
		"filename",
		"createby",
		"additionalinfo",
	}
	columnHeaders := []string{
		"No",
		"Identitas Stakeholder",
		"Kode Perusahaan",
		"Nama Perusahaan",
		"Waktu Pertanyaan / Keluhan",
		"Pertanyaan / Keluhan",
		"Waktu Jawaban / Respon",
		"Jawaban / Respon",
		"Topik",
		"Personel Follow Up",
		"Lampiran",
		"User",
		"Sumber Informasi Tambahan",
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	columnWidth := []float64{20, 50, 40, 40, 30, 40, 30, 40, 40, 50, 40, 30, 40, 60}
	tableHeaders := helper.GenerateTableHeaders(columnHeaders, columnWidth)

	var exportedData [][]string

	for i, item := range sortedData {
		var exportedRows []string
		exportedRows = append(exportedRows, strconv.Itoa(i+1))
		exportedRows = append(exportedRows, helper.MapToArray(item, exportedFields)...)
		for i, content := range exportedRows {
			if helper.IsContains([]int{4, 6, 11}, i) {
				unixTime, _ := strconv.Atoi(content)
				dateTime := time.Unix(int64(unixTime), 0)
				dateToFormat := helper.GetWIBLocalTime(&dateTime)
				exportedRows[i] = helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthShortNameInIndo) + " (" + helper.GetTimeAndMinuteOnly(dateToFormat) + ")"
			}
		}

		exportedData = append(exportedData, exportedRows)
	}
	exportTableProps := helper.ExportTableToFileProps{
		Filename: "PKP",
		Data:     exportedData,
		Headers:  tablesColumns,
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Pertanyaan Keluhan Pelanggan"},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:  630,
			Papperheight: 300,
			HeaderRows:   tableHeaders,
		},
		ColumnWidth: []int{20, 50, 40, 40, 50, 40, 50, 40, 40, 50, 40, 40, 40, 60},
	}
	errorExport := helper.ExportTableToFile(c, exportTableProps)
	if errorExport != nil {
		return nil, errorExport
	}

	paginatedData := helper.HandleDataPagination(c, sortedData, filterParameter)
	return &paginatedData, nil
}

func (uc *usecase) CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error) {
	return uc.pkpRepo.CreatePKuser(pkp, c)
}

func (uc *usecase) UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error) {
	return uc.pkpRepo.UpdatePKuser(pkp, c)
}

func (uc *usecase) Delete(id string, c *gin.Context) (int64, error) {
	return uc.pkpRepo.Delete(id, c)
}
