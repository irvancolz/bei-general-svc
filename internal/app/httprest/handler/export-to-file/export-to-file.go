package exporttofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"

	"github.com/gin-gonic/gin"
)

type ExportToFileinterface interface {
	ExportTableToFile(c *gin.Context)
}

type handler struct{}

func NewHandler() ExportToFileinterface {
	return &handler{}
}

func (h *handler) ExportTableToFile(c *gin.Context) {

	var request helper.ExportTableToFileProps

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	errExport := helper.ExportTableToFile(c, request)

	if errExport != nil {
		model.GenerateReadErrorResponse(c, errExport)
		return
	}

}
