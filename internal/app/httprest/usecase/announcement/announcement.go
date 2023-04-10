package announcement

import (
	"be-idx-tsg/internal/app/httprest/model"
	announcement "be-idx-tsg/internal/app/httprest/repository/announcement"
)

type Usecase interface {
	GetAllAnnouncement() ([]*model.Announcement, error)
}

type usecase struct {
	announcementRepo announcement.Repository
}

func NewUsecase() Usecase {
	return &usecase{
		announcement.NewRepository(),
	}
}
func (m *usecase) GetAllAnnouncement() ([]*model.Announcement, error) {
	return m.announcementRepo.GetAllAnnouncement()
}
