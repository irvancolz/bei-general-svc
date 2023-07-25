package log_system

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/usecase/log_system"
	"be-idx-tsg/internal/pkg/httpresponse"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAll(c *gin.Context)
	CreateLogSystem(c *gin.Context)
	ExportLogSystem(c *gin.Context)
}

type handler struct {
	log log_system.Usecase
}

func NewHandler() Handler {
	return &handler{
		log_system.DetailUseCase(),
	}
}

func (m *handler) GetAll(c *gin.Context) {
	data, err := m.log.GetAll()
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) CreateLogSystem(c *gin.Context) {
	var (
		request model.CreateLogSystem
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.log.CreateLogSystem(request, c)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateInsertErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (m *handler) ExportLogSystem(c *gin.Context) {
	err := m.log.ExportLogSystem(c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
}
