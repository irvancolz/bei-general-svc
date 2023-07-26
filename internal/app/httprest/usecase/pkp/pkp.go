package pkp

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	pkp "be-idx-tsg/internal/app/httprest/repository/pkp"
	"strconv"

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
		"createdat",
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
		"Waktu",
		"User",
		"Sumber Informasi Tambahan",
	}

	columnWidth := []float64{20, 50, 40, 40, 50, 40, 50, 40, 40, 50, 40, 40, 40, 60}

	var tableHeaders []helper.TableHeader

	for i, col := range columnHeaders {
		header := helper.TableHeader{
			Title: col,
			Width: columnWidth[i],
		}
		tableHeaders = append(tableHeaders, header)
	}

	var exportedData [][]string
	// exportedData = append(exportedData, columnHeaders)
	for i, item := range sortedData {
		var exportedRows []string
		exportedRows = append(exportedRows, strconv.Itoa(i+1))
		exportedRows = append(exportedRows, helper.MapToArray(item, exportedFields)...)

		exportedData = append(exportedData, exportedRows)
	}
	exportTableProps := helper.ExportTableToFileProps{
		Filename: "PKP",
		Data:     exportedData,
		Headers:  columnHeaders,
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Pertanyaan Keluhan Pelanggan"},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:  630,
			Papperheight: 300,
			HeaderRows:   tableHeaders,
		},
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
