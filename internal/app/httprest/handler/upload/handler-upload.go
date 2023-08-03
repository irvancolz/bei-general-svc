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
	UploadParameterAdminFile(c *gin.Context)
	UploadParameterAdminImage(c *gin.Context)
	IsFileExists(c *gin.Context)
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
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "form",
		MaxSize:   10240000, // 10mb max
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
		Host:       os.Getenv("MINIO_ENDPOINT"),
		Directory:  "user",
		Extensions: []string{".jpg", ".png", ".jpeg"},
		MaxSize:    10240000, // 10mb max
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
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "report",
		MaxSize:   10240000, // 10mb max
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
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "admin",
		MaxSize:   10240000, // 10mb max
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
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "pkp",
		MaxSize:   10240000, // 10mb max
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
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "guidebook",
		MaxSize:   10240000, // 10mb max
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}
func (h *handler) UploadParameterAdminImage(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "ParameterAdmin",
		MaxSize:   10240000, // 10mb max
	}
	result, error_result := h.Usecase.Upload(c, config)
	if error_result != nil {
		model.GenerateUploadErrorResponse(c, error_result)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPLOADSUCCESS_200, nil, result))
}

func (h *handler) UploadParameterAdminFile(c *gin.Context) {
	config := usecase.UploadFileConfig{
		Host:      os.Getenv("MINIO_ENDPOINT"),
		Directory: "ParameterAdmin",
		MaxSize:   10240000, // 10mb max
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
}

func (h *handler) Remove(c *gin.Context) {
	slug := c.Query("path")

	config := usecase.UploadFileConfig{
		Host:       os.Getenv("MINIO_ENDPOINT"),
		Directory:  "test",
		Extensions: []string{},
	}
	errorResult := h.Usecase.DeleteFile(c, config, slug)
	if errorResult != nil {
		model.GenerateDeleteErrorResponse(c, errorResult)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, "berhasil menghapus file"))
}

func (h *handler) IsFileExists(c *gin.Context) {
	slug := c.Query("path")

	errorResult := h.Usecase.IsFileExists(c, slug)
	if errorResult != nil {
		model.GenerateReadErrorResponse(c, errorResult)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, "file ada dalam penyimpanan"))
}
