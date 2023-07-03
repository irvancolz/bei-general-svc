package pkp

import (
	"be-idx-tsg/internal/app/httprest/model"
	PkpUsecase "be-idx-tsg/internal/app/httprest/usecase/pkp"
	"be-idx-tsg/internal/pkg/httpresponse"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler interface {
	GetAllPKuser(c *gin.Context)
	CreatePKuser(c *gin.Context)
	UpdatePKuser(c *gin.Context)
	Delete(c *gin.Context)
}

type handler struct {
	pkp PkpUsecase.Usecase
}

func NewHandler() Handler {
	return &handler{
		PkpUsecase.DetailUseCase(),
	}
}

func (m *handler) GetAllPKuser(c *gin.Context) {
	data, err := m.pkp.GetAllPKuser(c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	if !c.Writer.Written() {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
		return
	}
}

func (m *handler) CreatePKuser(c *gin.Context) {
	var (
		request model.CreatePKuser
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	data, err := m.pkp.CreatePKuser(request, c)
	if err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))

}

func (m *handler) UpdatePKuser(c *gin.Context) {
	var (
		request model.UpdatePKuser
	)
	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}
	data, err := m.pkp.UpdatePKuser(request, c)
	if err != nil {
		model.GenerateUpdateErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))

}

func (m *handler) Delete(c *gin.Context) {
	ID := c.Query("id")
	data, err := m.pkp.Delete(ID, c)
	if err != nil {
		model.GenerateInsertErrorResponse(c, err)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, data))

}
