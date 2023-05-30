package unggahberkas

import (
	"be-idx-tsg/internal/app/httprest/model"
	usecase "be-idx-tsg/internal/app/httprest/usecase/unggah-berkas"
	"be-idx-tsg/internal/pkg/httpresponse"

	"github.com/gin-gonic/gin"
)

type UnggahBerkasHandlerInterface interface {
	UploadNew(c *gin.Context)
	GetUploadedFiles(c *gin.Context)
	DeleteUploadedFiles(c *gin.Context)
}

type handler struct {
	Usecase usecase.UnggahBerkasUsecaseInterface
}

func NewHandler() UnggahBerkasHandlerInterface {
	return &handler{
		Usecase: usecase.NewUsecase(),
	}
}

func (h *handler) UploadNew(c *gin.Context) {
	var request usecase.UploadNewFilesProps
	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	results, errorResults := h.Usecase.UploadNew(c, request)
	if errorResults != nil {
		model.GenerateInsertErrorResponse(c, errorResults)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, results))

}
func (h *handler) GetUploadedFiles(c *gin.Context) {
	results, errorResults := h.Usecase.GetUploadedFiles(c)
	if errorResults != nil {
		model.GenerateReadErrorResponse(c, errorResults)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, results))
}
func (h *handler) DeleteUploadedFiles(c *gin.Context) {
	id := c.Query("id")
	errorResults := h.Usecase.DeleteUploadedFiles(c, id)
	if errorResults != nil {
		model.GenerateDeleteErrorResponse(c, errorResults)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, nil))
}
