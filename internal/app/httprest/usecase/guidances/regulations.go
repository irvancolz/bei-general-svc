package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateNewRegulationsAndFileProps struct {
	Name      string `json:"name" binding:"required"`
	File_name string `json:"file_name" binding:"required"`
	File_size int64  `json:"file_size" binding:"required"`
	File_path string `json:"file_path" binding:"required"`
}

type UpdateExistingRegulationsAndFileProps struct {
	Id        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	File_name string `json:"file_name" binding:"required"`
	File_size int64  `json:"file_size" binding:"required"`
	File_path string `json:"file_path" binding:"required"`
}

type RegulationUsecaseInterface interface {
	CreateNewRegulations(c *gin.Context, props CreateNewRegulationsAndFileProps) (int64, error)
	UpdateExistingRegulations(c *gin.Context, props UpdateExistingRegulationsAndFileProps) error
	GetAllRegulationsBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error)
}

func (r *guidancesUsecase) CreateNewRegulations(c *gin.Context, props CreateNewRegulationsAndFileProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createDataArgs := repo.CreateNewDataProps{
		Category:   "Regulation",
		Name:       props.Name,
		File:       props.File_name,
		File_size:  props.File_size,
		File_path:  props.File_path,
		Created_by: name_user.(string),
		Created_at: time.Now(),
	}
	result, error_result := r.Repository.CreateNewData(createDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}
func (r *guidancesUsecase) UpdateExistingRegulations(c *gin.Context, props UpdateExistingRegulationsAndFileProps) error {
	name_user, _ := c.Get("name_user")

	updateDataArgs := repo.UpdateExistingDataProps{
		Category:   "Regulation",
		Name:       props.Name,
		File:       props.File_name,
		File_size:  props.File_size,
		File_path:  props.File_path,
		Updated_by: name_user.(string),
		Updated_at: time.Now(),
		Id:         props.Id,
	}
	error_result := r.Repository.UpdateExistingData(updateDataArgs)
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
		if item.Category == types {
			result := model.RegulationJSONResponse{
				Id:         item.Id,
				Category:   item.Category,
				Name:       item.Name,
				Created_by: item.Created_by,
				File:       item.File,
				File_size:  item.File_size,
				File_path:  item.File_path,
				Version:    item.Version,
				Created_at: item.Created_at,
				Updated_by: item.Updated_by,
				Updated_at: item.Updated_at,
			}
			results = append(results, result)
		}
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})
	paginatedData := helper.HandleDataPagination(c, filteredData)
	return &paginatedData, nil
}
