package upload

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	usecase "be-idx-tsg/internal/app/httprest/usecase/upload"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/httpresponse"
	"errors"
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
		Extensions: GetExtension(c),
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
		Extensions: GetExtension(c),
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
		Extensions: GetExtension(c),
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
		Extensions: GetExtension(c),
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
		Extensions: GetExtension(c),
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
		Extensions: GetExtension(c),
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
		Host:       os.Getenv("DIR_HOST"),
		Directory:  "ParameterAdmin",
		Extensions: GetExtension(c),
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
		Host:      os.Getenv("DIR_HOST"),
		Directory: "ParameterAdmin",
		Extensions: GetExtension(c),
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

func GetExtension(c *gin.Context)([]string){
	datas, err := utilities.GetParameterAdminImageExtension(c)
	if err != nil {
		model.GenerateReadErrorResponse(c, errors.New("Fail to Get Data"))
		c.Abort()
	}
	ext := helper.ConvertListInterfaceToListString(datas.Data["value"].([]interface{}))
	return ext
}