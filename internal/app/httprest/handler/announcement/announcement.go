package announcement

import (
	// helper "be-idx-tsg/internal/app/helper"

	"be-idx-tsg/internal/app/httprest/model"
	AnnouncementUsecase "be-idx-tsg/internal/app/httprest/usecase/announcement"
	"be-idx-tsg/internal/pkg/httpresponse"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler interface {
	GetAllAnnouncement(c *gin.Context)
	DetailCode(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByCode(c *gin.Context)
	GetByIDandType(c *gin.Context)
	GetAllMin(c *gin.Context)
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
	var query = c.Query("keyword")

	data, err := m.an.GetAllANWithSearch(query)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
func (m *handler) DetailCode(c *gin.Context) {
	ID := c.Query("id")
	data, err := m.an.DetailCode(ID)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) GetAllAnnouncement(c *gin.Context) {
	data, err := m.an.GetAllAnnouncement()
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) GetAllMin(c *gin.Context) {

	data, err := m.an.GetAllMin()
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (m *handler) GetByCode(c *gin.Context) {
	var (
		request struct {
			Code string `json:"code" binding:"required"`
		}
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.an.GetByCode(request.Code)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

//

func (m *handler) Create(c *gin.Context) {
	var (
		request model.CreateAnnouncement
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}

	data, err := m.an.Create(request)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil))
	} else {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, nil))
	}
}
func (m *handler) Update(c *gin.Context) {
	var (
		request model.UpdateAnnouncement
	)
	// Update(id int, Type string, updated_by int) (int64, error)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.an.Update(request)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil))
	} else {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, nil))
	}
}

func (m *handler) Delete(c *gin.Context) {
	var (
		request struct {
			Id        string `json:"id" binding:"required"`
			DeletedBy string `json:"deleted_by" binding:"required"`
		}
	)
	// Delete(id int) (int64, error)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.an.Delete(request.Id, request.DeletedBy)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}
	if data == 1 {
		c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil))
	} else {
		c.JSON(httpresponse.Format(httpresponse.DELETEFAILED_400, nil))
	}
}

func (m *handler) GetByIDandType(c *gin.Context) {
	var (
		request struct {
			ID    string `json:"id" binding:"required"`
			Types string `json:"types" binding:"required"`
		}
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(httpresponse.Format(httpresponse.ERR_REQUESTBODY_400, err))
		return
	}
	data, err := m.an.GetByIDandType(request.ID, request.Types)
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
