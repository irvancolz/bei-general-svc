package announcement

import (
	// helper "be-idx-tsg/internal/app/helper"

	AnnouncementUsecase "be-idx-tsg/internal/app/httprest/usecase/announcement"
	"be-idx-tsg/internal/pkg/httpresponse"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler interface {
	GetAllAnnouncement(c *gin.Context)
}

type handler struct {
	announcement AnnouncementUsecase.Usecase
}

func NewHandler() Handler {
	return &handler{
		AnnouncementUsecase.NewUsecase(),
	}
}

func (m *handler) GetAllAnnouncement(c *gin.Context) {
	data, err := m.announcement.GetAllAnnouncement()
	if err != nil {
		c.JSON(httpresponse.Format(httpresponse.READFAILED_400, err))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}
