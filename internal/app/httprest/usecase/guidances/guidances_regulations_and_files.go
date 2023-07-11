package guidances

import (
	"be-idx-tsg/internal/app/helper"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"strconv"

	"github.com/gin-gonic/gin"
)

type guidancesUsecase struct {
	Repository repo.GuidancesRepoInterface
}

type GuidancesRegulationAndFileUsecaseInterface interface {
	GuidancesUsecaseInterface
	RegulationUsecaseInterface
	FilesUsecaseInterface
	GetAllData(c *gin.Context) (*helper.PaginationResponse, error)
}

func NewGuidanceUsecase() GuidancesRegulationAndFileUsecaseInterface {
	return &guidancesUsecase{
		Repository: repo.NewGuidancesRepository(),
	}
}

func (u *guidancesUsecase) GetAllData(c *gin.Context) (*helper.PaginationResponse, error) {
	databaseResult, errorResult := u.Repository.GetAllData(c)
	if errorResult != nil {
		return nil, errorResult
	}

	var dataToConverted []interface{}
	for _, item := range databaseResult {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})
	sortedData := helper.HandleDataSorting(c, filteredData)
	exportedFields := []string{
		"category",
		"name",
		"description",
		"file",
		"version",
		"order",
		"created_by",
		"link",
	}
	columnHeaders := []string{
		"No",
		"Kategori",
		"Daftar Berkas",
		"Deskripsi",
		"Nama Berkas",
		"Versi",
		"Urutan",
		"User",
		"Link",
	}

	var exportedData [][]string
	exportedData = append(exportedData, columnHeaders)
	for i, item := range sortedData {
		var exportedRows []string
		exportedRows = append(exportedRows, strconv.Itoa(i+1))
		exportedRows = append(exportedRows, helper.MapToArray(item, exportedFields)...)

		exportedData = append(exportedData, exportedRows)
	}
	exportTableProps := helper.ExportTableToFileProps{
		Filename: "Management Berkas",
		Data:     exportedData,
		Headers:  columnHeaders,
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Management Berkas "},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:     300,
			Papperheight:    450,
			PageOrientation: "l",
		},
	}
	errorExport := helper.ExportTableToFile(c, exportTableProps)
	if errorExport != nil {
		return nil, errorExport
	}

	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)
	return &paginatedData, nil

}
