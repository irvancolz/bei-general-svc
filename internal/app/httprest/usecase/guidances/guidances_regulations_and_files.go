package guidances

import (
	"be-idx-tsg/internal/app/helper"
	repo "be-idx-tsg/internal/app/httprest/repository/guidances"

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
	databaseResult, errorResult := u.Repository.GetAllData()
	if errorResult != nil {
		return nil, errorResult
	}

	var dataToConverted []interface{}
	for _, item := range databaseResult {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})
	paginatedData := helper.HandleDataPagination(c, filteredData)
	return &paginatedData, nil

}
