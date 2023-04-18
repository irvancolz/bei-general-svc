package guidances

import (
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateNewRegulationsProps struct {
	Name  string `json:"name"`
	Link  string `json:"link"`
	Order int64  `json:"order"`
}

type UpdateExistingRegulationsProps struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Link  string `json:"link"`
	Order int64  `json:"order"`
}

type RegulationUsecaseInterface interface {
	CreateNewRegulations(c *gin.Context, props CreateNewRegulationsProps) (int64, error)
	UpdateExistingRegulations(c *gin.Context, props UpdateExistingRegulationsProps) error
	GetAllRegulationsBasedOnType(c *gin.Context, types string) ([]*model.RegulationJSONResponse, error)
}

func (r *guidancesUsecase) CreateNewRegulations(c *gin.Context, props CreateNewRegulationsProps) (int64, error) {
	name_user, _ := c.Get("name_user")

	createDataArgs := repo.CreateNewDataProps{
		Category:   "Regulation",
		Name:       props.Name,
		Link:       props.Link,
		Order:      props.Order,
		Created_by: name_user.(string),
		Created_at: time.Now(),
		Version:    1.0, //hardcoded, because no params given
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
		Category:   "Regulation",
		Name:       props.Name,
		Link:       props.Link,
		Order:      props.Order,
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
func (r *guidancesUsecase) GetAllRegulationsBasedOnType(c *gin.Context, types string) ([]*model.RegulationJSONResponse, error) {
	var results []*model.RegulationJSONResponse
	raw_result, error_result := r.Repository.GetAllDataBasedOnCategory(types)
	if error_result != nil {
		return nil, error_result
	}
	for _, item := range raw_result {
		result := model.RegulationJSONResponse{
			Id:         item.Id,
			Category:   item.Category,
			Name:       item.Name,
			Link:       item.Link,
			Version:    item.Version,
			File:       item.File,
			File_size:  item.File_size,
			Order:      item.Order,
			Created_by: item.Created_by,
			Created_at: item.Created_at,
			Updated_by: item.Updated_by,
			Updated_at: item.Updated_at,
		}
		results = append(results, &result)
	}
	return results, nil
}
