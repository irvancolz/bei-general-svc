package guidances

import (
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"time"

	"github.com/gin-gonic/gin"
)

type FilesUsecaseInterface interface {
	GetAllFilesOnType(c *gin.Context, types string) ([]*model.GuidanceFilesJSONResponse, error)
	UpdateExistingFiles(c *gin.Context, props UpdateExistingRegulationsAndFileProps) error
	CreateNewFiles(c *gin.Context, props CreateNewRegulationsAndFileProps) (int64, error)
}

func (u *guidancesUsecase) GetAllFilesOnType(c *gin.Context, types string) ([]*model.GuidanceFilesJSONResponse, error) {
	var results []*model.GuidanceFilesJSONResponse
	raw_result, error_result := u.Repository.GetAllData()
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		if item.Category == types {
			result := model.GuidanceFilesJSONResponse{
				Id:         item.Id,
				Name:       item.Name,
				Category:   item.Category,
				File_size:  item.File_size,
				File:       item.File,
				Created_by: item.Created_by,
				Created_at: item.Created_at,
				Updated_by: item.Updated_by,
				Updated_at: item.Updated_at,
			}
			results = append(results, &result)
		}
	}
	return results, nil
}

func (u *guidancesUsecase) UpdateExistingFiles(c *gin.Context, props UpdateExistingRegulationsAndFileProps) error {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.UpdateExistingDataProps{
		Id:         props.Id,
		Category:   "File",
		Name:       props.Name,
		File:       props.File_name,
		File_size:  props.File_size,
		Updated_at: time.Now(),
		Updated_by: name_user.(string),
	}
	error_result := u.Repository.UpdateExistingData(createNewDataArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}

func (u *guidancesUsecase) CreateNewFiles(c *gin.Context, props CreateNewRegulationsAndFileProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createNewDataArgs := repo.CreateNewDataProps{
		Category:   "File",
		Name:       props.Name,
		File:       props.File_name,
		File_size:  props.File_size,
		Created_at: time.Now(),
		Created_by: name_user.(string),
	}
	result, error_result := u.Repository.CreateNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}
