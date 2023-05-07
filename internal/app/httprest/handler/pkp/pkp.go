package pkp

import (
	"be-idx-tsg/internal/app/httprest/model"
	PkpUsecase "be-idx-tsg/internal/app/httprest/usecase/pkp"
	"be-idx-tsg/internal/pkg/httpresponse"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler interface {
	GetAllPKuser(c *gin.Context)
	CreatePKuser(c *gin.Context)
	UpdatePKuser(c *gin.Context)
	Delete(c *gin.Context)
	GetAllWithFilter(c *gin.Context)
	GetAllWithSearch(c *gin.Context)
}

type handler struct {
	pkp PkpUsecase.Usecase
}

func NewHandler() Handler {
	return &handler{
		PkpUsecase.DetailUseCase(),
	}
}

func (m *handler) GetAllWithFilter(c *gin.Context) {
	var query = c.Query("keyword")

	var list = strings.Split(query, ",")

	data, err := m.pkp.GetAllWithFilter(list)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
func (m *handler) GetAllWithSearch(c *gin.Context) {
	var (
		request struct {
			Code         string    `json:"code"`
			Name         string    `json:"name"`
			QuestionDate time.Time `json:"question_date"`
			Question     string    `json:"question"`
			Answers      string    `json:"answers"`
			AnsweredBy   string    `json:"answered_by"`
			AnsweredAt   time.Time `json:"answered_at"`
		}
	)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.pkp.GetAllWithSearch(request.Code, request.Name, request.QuestionDate, request.Question, request.Answers, request.AnsweredBy, request.AnsweredAt)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) GetAllPKuser(c *gin.Context) {
	data, err := m.pkp.GetAllPKuser(c)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) CreatePKuser(c *gin.Context) {
	var (
		request model.CreatePKuser
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}

	data, err := m.pkp.CreatePKuser(request, c)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
	} else {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, nil, data))
	}
}

func (m *handler) UpdatePKuser(c *gin.Context) {
	var (
		request model.UpdatePKuser
	)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.pkp.UpdatePKuser(request, c)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))

}

func (m *handler) Delete(c *gin.Context) {
	ID := c.Query("id")
	data, err := m.pkp.Delete(ID, c)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, data))

}
