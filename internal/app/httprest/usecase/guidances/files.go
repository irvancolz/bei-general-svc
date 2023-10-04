package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FilesUsecaseInterface interface {
	GetAllFilesOnType(c *gin.Context, types string) (*helper.PaginationResponse, error)
	UpdateExistingFiles(c *gin.Context, props UpdateExsistingGuidancesAndFilesProps) error
	CreateNewFiles(c *gin.Context, props CreateNewGuidanceAndFilesProps) (int64, error)
}

const Berkas = "Berkas"

func (u *guidancesUsecase) GetAllFilesOnType(c *gin.Context, types string) (*helper.PaginationResponse, error) {
	var results []model.GuidanceFilesJSONResponse
	raw_result, error_result := u.Repository.GetAllData(c)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		if strings.EqualFold(item.Category, types) {
			result := model.GuidanceFilesJSONResponse{
				Id:          item.Id,
				Name:        item.Name,
				Category:    item.Category,
				Description: item.Description,
				Version:     item.Version,
				File_size:   item.File_size,
				File_path:   item.File_path,
				File:        item.File,
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
	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)
	return &paginatedData, nil

}

func (u *guidancesUsecase) UpdateExistingFiles(c *gin.Context, props UpdateExsistingGuidancesAndFilesProps) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:          props.Id,
		Category:    Berkas,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_path:   props.File_path,
		File_size:   props.File_size,
		Order:       props.Order,
		Version:     props.Version,
		File_Owner:  props.Owner,
		Updated_at:  time.Now(),
		Updated_by:  name_user.(string),
	}

	isOrderFilled := u.Repository.CheckIsOrderFilled(createNewDataArgs.Order, Berkas)
	if isOrderFilled {
		errorSetOrder := u.Repository.UpdateOrder(createNewDataArgs.Order, Berkas)
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

func (u *guidancesUsecase) CreateNewFiles(c *gin.Context, props CreateNewGuidanceAndFilesProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:    Berkas,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_path:   props.File_path,
		File_size:   props.File_size,
		Order:       props.Order,
		Version:     props.Version,
		File_Owner:  props.Owner,
		Created_at:  time.Now(),
		Created_by:  name_user.(string),
	}

	isOrderFilled := u.Repository.CheckIsOrderFilled(createNewDataArgs.Order, Berkas)
	if isOrderFilled {
		errorSetOrder := u.Repository.UpdateOrder(createNewDataArgs.Order, Berkas)
		if errorSetOrder != nil {
			return 0, errorSetOrder
		}
	}

	result, error_result := u.Repository.CreateNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}
