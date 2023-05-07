package upload

import (
	"be-idx-tsg/internal/app/httprest/model"
	usecase "be-idx-tsg/internal/app/httprest/usecase/upload"
	"be-idx-tsg/internal/pkg/httpresponse"
	"os"

	"github.com/gin-gonic/gin"
)

type UploadFileHandlreInterface interface {
	UploadForm(c *gin.Context)
	UploadReport(c *gin.Context)
	UploadAdmin(c *gin.Context)
	UploadUser(c *gin.Context)
	UploadPkp(c *gin.Context)
	UploadGuidebook(c *gin.Context)
	Download(c *gin.Context)
	Remove(c *gin.Context)
}

type handler struct {
	Usecase usecase.UploadFileUsecaseInterface
}

func NewHandler() UploadFileHandlreInterface {
	return &handler{
		Usecase: usecase.NewUsecase(),
	}
}

func (h *handler) UploadForm(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "form",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadUser(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "user",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadReport(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "report",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadAdmin(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "admin",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadPkp(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "pkp",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadGuidebook(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("DIR_HOST"),
		Directory: "guidebook",
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) Download(c *gin.Context) {
	filepth := c.Query("path")
	errorResult := h.Usecase.Download(c, filepth)
	if errorResult != nil {
		model.GenerateIFileNotFoundErrorResponse(c, errorResult)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DOWNLOADSUCCESS_200, nil, "file berhasil diunduh"))
}

func (h *handler) Remove(c *gin.Context) {
	slug := c.Query("path")

	config := usecase.UploadFileConfig{
		Host:       os.Getenv("DIR_HOST"),
		Directory:  "test",
		Extensions: []string{".pdf"},
	}
	errorResult := h.Usecase.DeleteFile(c, config, slug)
	if errorResult != nil {
		model.GenerateDeleteErrorResponse(c, errorResult)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, "berhasil menghapus file"))
}
