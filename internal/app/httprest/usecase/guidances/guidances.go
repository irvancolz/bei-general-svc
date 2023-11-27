package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GuidancesUsecaseInterface interface {
	CreateNewGuidance(c *gin.Context, props CreateNewGuidanceAndFilesProps) (int64, error)
	UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidancesAndFilesProps) error
	GetAllGuidanceBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error)
	DeleteGuidances(c *gin.Context, id string) error
}

type CreateNewGuidanceAndFilesProps struct {
	Owner       string `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	File        string `json:"file" binding:"required"`
	File_size   int64  `json:"file_size" binding:"required"`
	File_path   string `json:"file_path" binding:"required"`
	Version     string `json:"version" binding:"required,numeric"`
	Order       int    `json:"order" binding:"min=1,required"`
}

type UpdateExsistingGuidancesAndFilesProps struct {
	Id          string `json:"id" binding:"required"`
	Owner       string `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	File        string `json:"file" binding:"required"`
	File_path   string `json:"file_path" binding:"required"`
	File_size   int64  `json:"file_size" binding:"required"`
	Version     string `json:"version" binding:"required,numeric"`
	Order       int    `json:"order" binding:"min=1,required"`
}

const BukuPetunjuk = "Buku Petunjuk"

func (u *guidancesUsecase) UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidancesAndFilesProps) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:          props.Id,
		Category:    BukuPetunjuk,
		File_Owner:  props.Owner,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_size:   props.File_size,
		File_path:   props.File_path,
		Version:     props.Version,
		Order:       props.Order,
		Updated_at:  time.Now(),
		Updated_by:  name_user.(string),
	}

	isOrderFilled := u.Repository.CheckIsOrderFilled(createNewDataArgs.Order, BukuPetunjuk)
	if isOrderFilled {
		errorSetOrder := u.Repository.UpdateOrder(createNewDataArgs.Order, BukuPetunjuk)
		if errorSetOrder != nil {
			return errorSetOrder
		}
	}

	error_result := u.Repository.UpdateExistingData(c, createNewDataArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}

func (u *guidancesUsecase) CreateNewGuidance(c *gin.Context, props CreateNewGuidanceAndFilesProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:    BukuPetunjuk,
		File_Owner:  props.Owner,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_size:   props.File_size,
		Version:     props.Version,
		Order:       props.Order,
		File_path:   props.File_path,
		Created_at:  time.Now(),
		Created_by:  name_user.(string),
	}

	isOrderFilled := u.Repository.CheckIsOrderFilled(createNewDataArgs.Order, BukuPetunjuk)
	if isOrderFilled {
		errorSetOrder := u.Repository.UpdateOrder(createNewDataArgs.Order, BukuPetunjuk)
		if errorSetOrder != nil {
			return 0, errorSetOrder
		}
	}

	result, error_result := u.Repository.CreateNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}

	utilities.CreateNotifForAdminApp(c, "management berkas", fmt.Sprintf("%s menambahkan berkas baru", name_user.(string)))
	email.SendEmailForUserAdminApp(c, "Penambahan Berkas Baru", fmt.Sprintf("%s menambahkan berkas baru", name_user.(string)))

	return result, nil
}

func (u *guidancesUsecase) GetAllGuidanceBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error) {
	var results []model.GuidanceJSONResponse
	raw_result, error_result := u.Repository.GetAllData(c)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		if strings.EqualFold(item.Category, types) {
			result := model.GuidanceJSONResponse{
				Id:          item.Id,
				Name:        item.Name,
				Category:    item.Category,
				Description: item.Description,
				Version:     item.Version,
				File:        item.File,
				File_size:   item.File_size,
				File_path:   item.File_path,
				File_Group:  item.File_Group,
				Owner:       item.File_Owner,
				Link:        item.Link,
				Order:       item.Order,
				Created_by:  item.Created_by,
				Created_at:  item.Created_at,
				Updated_by:  item.Updated_by,
				Updated_at:  item.Updated_at,
			}
			results = append(results, result)
		}
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}
	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})

	columnHeaders := []string{"No", "Nama berkas", "Deskripsi", "Versi", "File Lampiran", "Ukuran File"}
	columnWidth := []float64{20, 50, 60, 20, 50, 30}

	tableHeaders := helper.GenerateTableHeaders(columnHeaders, columnWidth)

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	dataOrder := []string{"name", "description", "version", "file", "file_size"}
	var exportedData [][]string

	for i, content := range filteredData {
		var item []string
		item = append(item, fmt.Sprintf("%v", i+1))
		item = append(item, helper.MapToArray(content, dataOrder)...)

		exportedData = append(exportedData, item)
	}

	var columnWidthInINT []int
	for _, width := range columnWidth {
		columnWidthInINT = append(columnWidthInINT, int(width))
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename:    "guidances",
		Data:        exportedData,
		Headers:     tablesColumns,
		ColumnWidth: columnWidthInINT,
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Buku Petunjuk"},
		},
		PdfConfig: &helper.PdfTableOptions{
			HeaderRows:      tableHeaders,
			PageOrientation: "l",
		},
	}

	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return nil, errorExport
	}
	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)
	return &paginatedData, nil
}

func (u *guidancesUsecase) DeleteGuidances(c *gin.Context, id string) error {
	user_id, _ := c.Get("name_user")
	deleteGuidancesArgs := repo.DeleteExistingDataProps{
		Deleted_at: time.Now(),
		Deleted_by: user_id.(string),
		Id:         id,
	}
	error_result := u.Repository.DeleteExistingData(deleteGuidancesArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}
