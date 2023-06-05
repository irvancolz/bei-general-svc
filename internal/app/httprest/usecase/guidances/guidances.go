package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type GuidancesUsecaseInterface interface {
	CreateNewGuidance(c *gin.Context, props CreateNewGuidanceProps) (int64, error)
	UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidances) error
	GetAllGuidanceBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error)
	DeleteGuidances(c *gin.Context, id string) error
}

type CreateNewGuidanceProps struct {
	Owner       string `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	File        string `json:"file" binding:"required"`
	File_size   int64  `json:"file_size" binding:"required"`
	File_path   string `json:"file_path" binding:"required"`
	Version     string `json:"version" binding:"required,numeric"`
	Is_visible  bool   `json:"visible" binding:"required"`
	Link        string `json:"link" binding:"required"`
}
type UpdateExsistingGuidances struct {
	Id          string `json:"id" binding:"required"`
	Owner       string `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	File        string `json:"file" binding:"required"`
	File_path   string `json:"file_path" binding:"required"`
	File_size   int64  `json:"file_size" binding:"required"`
	Version     string `json:"version" binding:"required,numeric"`
	Is_visible  bool   `json:"visible" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

func (u *guidancesUsecase) UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidances) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:          props.Id,
		Category:    "Guidebook",
		File_Owner:  props.Owner,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_size:   props.File_size,
		File_path:   props.File_path,
		Version:     props.Version,
		Is_Visible:  props.Is_visible,
		Link:        props.Link,
		Updated_at:  time.Now(),
		Updated_by:  name_user.(string),
	}
	error_result := u.Repository.UpdateExistingData(createNewDataArgs)
	if error_result != nil {
		log.Println(error_result)
		return error_result
	}
	return nil
}
func (u *guidancesUsecase) CreateNewGuidance(c *gin.Context, props CreateNewGuidanceProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:    "Guidebook",
		File_Owner:  props.Owner,
		Name:        props.Name,
		Description: props.Description,
		File:        props.File,
		File_size:   props.File_size,
		Version:     props.Version,
		Is_Visible:  props.Is_visible,
		File_path:   props.File_path,
		Link:        props.Link,
		Created_at:  time.Now(),
		Created_by:  name_user.(string),
	}
	result, error_result := u.Repository.CreateNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}

func (u *guidancesUsecase) GetAllGuidanceBasedOnType(c *gin.Context, types string) (*helper.PaginationResponse, error) {
	var results []model.GuidanceJSONResponse
	raw_result, error_result := u.Repository.GetAllData(c)
	if error_result != nil {
		log.Println(error_result)
		return nil, error_result
	}
	for _, item := range raw_result {
		if item.Category == types {
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
				Is_Visible:  item.Is_Visible,
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

	filteredData := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})
	paginatedData := helper.HandleDataPagination(c, filteredData)
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
