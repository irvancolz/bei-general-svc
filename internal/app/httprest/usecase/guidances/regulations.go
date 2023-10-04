package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateNewRegulationsProps struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
	Order int    `json:"order" binding:"min=1,required"`
	Link  string `json:"link" binding:"required"`
}

type UpdateExistingRegulationsProps struct {
	Id    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
	Order int    `json:"order" binding:"min=1,required"`
	Link  string `json:"link" binding:"required"`
}

type RegulationUsecaseInterface interface {
	CreateNewRegulations(c *gin.Context, props CreateNewRegulationsProps) (int64, error)
	UpdateExistingRegulations(c *gin.Context, props UpdateExistingRegulationsProps) error
	GetAllRegulationsBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error)
}

const Peraturan = "Peraturan"

func (r *guidancesUsecase) CreateNewRegulations(c *gin.Context, props CreateNewRegulationsProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createDataArgs := repo.CreateNewDataProps{
		Category:   Peraturan,
		Name:       props.Name,
		File_Owner: props.Owner,
		Link:       props.Link,
		Order:      props.Order,
		Created_by: name_user.(string),
		Created_at: time.Now(),
	}

	result, error_result := r.Repository.CreateNewData(createDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}

func (r *guidancesUsecase) UpdateExistingRegulations(c *gin.Context, props UpdateExistingRegulationsProps) error {
	name_user, _ := c.Get("name_user")

	updateDataArgs := repo.UpdateExistingDataProps{
		Category:   Peraturan,
		Name:       props.Name,
		File_Owner: props.Owner,
		Link:       props.Link,
		Order:      props.Order,
		Updated_by: name_user.(string),
		Updated_at: time.Now(),
		Id:         props.Id,
	}

	if props.Order <= 0 {
		updateDataArgs.Order = 1
	}

	isOrderFilled := r.Repository.CheckIsOrderFilled(updateDataArgs.Order, Peraturan)
	if isOrderFilled {
		errorSetOrder := r.Repository.UpdateOrder(updateDataArgs.Order, Peraturan)
		if errorSetOrder != nil {
			return errorSetOrder
		}
	}

	error_result := r.Repository.UpdateExistingData(c, updateDataArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}
func (r *guidancesUsecase) GetAllRegulationsBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error) {
	var results []model.RegulationJSONResponse
	raw_result, error_result := r.Repository.GetAllData(c)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		if strings.EqualFold(item.Category, types) {
			result := model.RegulationJSONResponse{
				Id:         item.Id,
				Category:   item.Category,
				Name:       item.Name,
				Created_by: item.Created_by,
				Created_at: item.Created_at,
				Updated_by: item.Updated_by,
				Updated_at: item.Updated_at,
				Order:      item.Order,
				Link:       item.Link,
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
