package topic

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/usecase/topic"
	"be-idx-tsg/internal/pkg/httpresponse"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	CreateTopicWithMessage(c *gin.Context)
	UpdateHandler(c *gin.Context)
	CreateMessage(c *gin.Context)
	DeleteTopic(c *gin.Context)
	ArchiveTopicToFAQ(c *gin.Context)
}

type handler struct {
	tp topic.Usecase
}

func NewHandler() Handler {
	return &handler{
		topic.DetailUseCase(),
	}
}

func (m *handler) GetAll(c *gin.Context) {
	var keyword = c.Query("keyword")
	var page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	var limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	var status = c.Query("status")
	var name = c.Query("name")
	var company_name = c.Query("company_name")
	var start_date = c.Query("start_date")
	var end_date = c.Query("end_date")

	data, err := m.tp.GetAll(keyword, status, name, company_name, start_date, end_date, page, limit)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	totalData, totalPage, err := m.tp.GetTotal(keyword, status, name, company_name, start_date, end_date, page, limit)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	var pagination = []any{data, totalData, page, limit, len(data), totalPage}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, pagination...))
}

func (m *handler) GetById(c *gin.Context) {
	ID := c.Query("id")
	keyword := c.Query("keyword")

	data, err := m.tp.GetByID(ID, keyword)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) CreateTopicWithMessage(c *gin.Context) {
	var (
		request model.CreateTopicWithMessage
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.tp.CreateTopicWithMessage(request, c)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateInsertErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (m *handler) UpdateHandler(c *gin.Context) {
	var (
		request model.UpdateTopicHandler
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.tp.UpdateHandler(request, c)
	if err != nil {
		model.GenerateUpdateErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateUpdateErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

func (m *handler) CreateMessage(c *gin.Context) {
	var (
		request model.CreateMessage
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.tp.CreateMessage(request, c)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateInsertErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (m *handler) DeleteTopic(c *gin.Context) {
	ID := c.Query("id")

	data, err := m.tp.DeleteTopic(ID, c)
	if err != nil {
		model.GenerateDeleteErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, data))
}

func (m *handler) ArchiveTopicToFAQ(c *gin.Context) {
	var (
		request model.ArchiveTopicToFAQ
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.tp.ArchiveTopicToFAQ(request, c)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateInsertErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}