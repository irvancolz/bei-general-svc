package announcement

import (
	"be-idx-tsg/internal/app/httprest/model"
	AnnouncementUsecase "be-idx-tsg/internal/app/httprest/usecase/announcement"
	"be-idx-tsg/internal/pkg/httpresponse"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler interface {
	GetAllAnnouncement(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAllANWithFilter(c *gin.Context)
	GetAllANWithSearch(c *gin.Context)
}

type handler struct {
	an AnnouncementUsecase.Usecase
}

func NewHandler() Handler {
	return &handler{
		AnnouncementUsecase.DetailUseCase(),
	}
}
func (m *handler) GetAllANWithFilter(c *gin.Context) {
	var query = c.Query("keyword")

	var list = strings.Split(query, ",")

	data, err := m.an.GetAllANWithFilter(list)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
func (m *handler) GetAllANWithSearch(c *gin.Context) {
	var (
		request struct {
			Keyword         string `json:"keyword"`
			InformationType string `json:"information_type"`
			StartDate       string `json:"start_date"`
			EndDate         string `json:"end_date"`
		}
	)
	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}
	data, err := m.an.GetAllANWithSearch(request.Keyword, request.InformationType, request.StartDate, request.EndDate)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
func (m *handler) GetById(c *gin.Context) {
	ID := c.Query("id")
	data, err := m.an.Detail(ID, c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) GetAllAnnouncement(c *gin.Context) {
	data, err := m.an.GetAllAnnouncement(c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) Create(c *gin.Context) {
	var (
		request model.CreateAnnouncement
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	data, err := m.an.Create(request, c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
	} else {
		model.GenerateInsertErrorResponse(c, err)
	}
}
func (m *handler) Update(c *gin.Context) {
	var (
		request model.UpdateAnnouncement
	)
	if err := c.ShouldBindJSON(&request); err != nil {
		model.GenerateInvalidJsonResponse(c, err) 
		return
	}
	data, err := m.an.Update(request, c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil))
	} else {
		model.GenerateUpdateErrorResponse(c, err)
	}
}

func (m *handler) Delete(c *gin.Context) {
	ID := c.Query("id")
	data, err := m.an.Delete(ID, c)
	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil))
	} else {
		model.GenerateDeleteErrorResponse(c, err)
	}
}
