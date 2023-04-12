package guidances

import (
	"be-idx-tsg/internal/app/httprest/model"

	"github.com/gin-gonic/gin"
)

type FilesUsecaseInterface interface {
	GetAllGuidanceBasedOnType(c *gin.Context, types string) ([]*model.GuidanceFilesJSONResponse, error)
}
