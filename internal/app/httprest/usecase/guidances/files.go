package guidances

import (
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"time"

	"github.com/gin-gonic/gin"
)

type FilesUsecaseInterface interface {
	GetAllFilesOnType(c *gin.Context, types string) ([]*model.GuidanceFilesJSONResponse, error)
	UpdateExistingFiles(c *gin.Context, props UpdateExsistingGuidances) error
	CreateNewFiles(c *gin.Context, props CreateNewGuidanceProps) (int64, error)
}

func (u *guidancesUsecase) GetAllFilesOnType(c *gin.Context, types string) ([]*model.GuidanceFilesJSONResponse, error) {
	var results []*model.GuidanceFilesJSONResponse
	raw_result, error_result := u.Repository.GetAllDataBasedOnCategory(types)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		result := model.GuidanceFilesJSONResponse{
			Id:         item.Id,
			Category:   item.Category,
			Name:       item.Name,
			File_size:  item.File_size,
			Created_by: item.Created_by,
			Created_at: item.Created_at,
			Updated_by: item.Updated_by,
			Updated_at: item.Updated_at,
		}
		results = append(results, &result)
	}
	return results, nil
}

func (u *guidancesUsecase) UpdateExistingFiles(c *gin.Context, props UpdateExsistingGuidances) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:          props.Id,
		Category:    "File",
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

func (u *guidancesUsecase) CreateNewFiles(c *gin.Context, props CreateNewGuidanceProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:    "File",
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
