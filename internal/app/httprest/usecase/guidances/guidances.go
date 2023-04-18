package guidances

import (
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"time"

	"github.com/gin-gonic/gin"
)

type GuidancesUsecaseInterface interface {
	CreateNewGuidance(c *gin.Context, props CreateNewGuidanceProps) (int64, error)
	UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidances) error
	GetAllGuidanceBasedOnType(c *gin.Context, types string) ([]*model.GuidanceJSONResponse, error)
	DeleteGuidances(c *gin.Context, id string) error
}

type CreateNewGuidanceProps struct {
	Description string  `json:"description"`
	Name        string  `json:"name"`
	File        string  `json:"file" binding:"required"`
	File_size   string  `json:"file_size" binding:"required"`
	Version     float64 `json:"version"`
	Order       int64   `json:"order"`
}
type UpdateExsistingGuidances struct {
	Id          string  `json:"id" binding:"required"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	File        string  `json:"file" binding:"required"`
	File_size   string  `json:"file_size" binding:"required"`
	Version     float64 `json:"version"`
	Order       int64   `json:"order"`
}

func (u *guidancesUsecase) UpdateExistingGuidances(c *gin.Context, props UpdateExsistingGuidances) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:          props.Id,
		Category:    "Guidebook",
		Description: props.Description,
		Name:        props.Name,
		File:        props.File,
		File_size:   props.File_size,
		Version:     props.Version,
		Order:       props.Order,
		Updated_at:  time.Now(),
		Updated_by:  name_user.(string),
	}
	error_result := u.Repository.UpdateExistingData(createNewDataArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}
func (u *guidancesUsecase) CreateNewGuidance(c *gin.Context, props CreateNewGuidanceProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:    "Guidebook",
		Description: props.Description,
		Name:        props.Name,
		File:        props.File,
		File_size:   props.File_size,
		Version:     props.Version,
		Order:       props.Order,
		Created_at:  time.Now(),
		Created_by:  name_user.(string),
	}
	result, error_result := u.Repository.CreateNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}

func (u *guidancesUsecase) GetAllGuidanceBasedOnType(c *gin.Context, types string) ([]*model.GuidanceJSONResponse, error) {
	var results []*model.GuidanceJSONResponse
	raw_result, error_result := u.Repository.GetAllDataBasedOnCategory(types)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		result := model.GuidanceJSONResponse{
			Id:          item.Id,
			Category:    item.Category,
			Name:        item.Name,
			Description: item.Description,
			Version:     item.Version,
			Order:       item.Order,
			File:        item.File,
			File_size:   item.File_size,
			Created_by:  item.Created_by,
			Created_at:  item.Created_at,
			Updated_by:  item.Updated_by,
			Updated_at:  item.Updated_at,
		}
		results = append(results, &result)
	}
	return results, nil
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
