package faq

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/usecase/faq"
	"be-idx-tsg/internal/pkg/httpresponse"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAll(c *gin.Context)
	CreateFAQ(c *gin.Context)
	DeleteFAQ(c *gin.Context)
	UpdateStatusFAQ(c *gin.Context)
	UpdateFAQ(c *gin.Context)
	UpdateOrderFAQ(c *gin.Context)
}

type handler struct {
	faq faq.Usecase
}

func NewHandler() Handler {
	return &handler{
		faq.DetailUseCase(),
	}
}

func (m *handler) GetAll(c *gin.Context) {
	var keyword = c.Query("keyword")
	userId, _ := c.Get("user_id")

	data, err := m.faq.GetAll(keyword, userId.(string))
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) CreateFAQ(c *gin.Context) {
	var (
		request model.CreateFAQ
		isDraft bool
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	isDraft, _ = strconv.ParseBool(c.DefaultQuery("draft", "0"))

	data, err := m.faq.CreateFAQ(request, c, isDraft)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	if data != 1 {
		model.GenerateInsertErrorResponse(c, err)
	}

	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (m *handler) DeleteFAQ(c *gin.Context) {
	ID := c.Query("id")

	data, err := m.faq.DeleteFAQ(ID, c)
	if err != nil {
		model.GenerateDeleteErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, data))
}

func (m *handler) UpdateStatusFAQ(c *gin.Context) {
	var (
		request model.UpdateFAQStatus
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.faq.UpdateStatusFAQ(request, c)
	if err != nil {
		model.GenerateUpdateErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

func (m *handler) UpdateFAQ(c *gin.Context) {
	var (
		request model.UpdateFAQ
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.faq.UpdateFAQ(request, c)
	if err != nil {
		model.GenerateUpdateErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

func (m *handler) UpdateOrderFAQ(c *gin.Context) {
	var (
		request []model.UpdateFAQOrder
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.faq.UpdateOrderFAQ(request, c)
	if err != nil {
		model.GenerateUpdateErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}
