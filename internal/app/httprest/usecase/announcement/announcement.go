package announcement

import (
	"be-idx-tsg/internal/app/httprest/model"
	an "be-idx-tsg/internal/app/httprest/repository/announcement"
)

type Usecase interface {
	GetAllAnnouncement() ([]*model.Announcement, error)
	DetailCode(id string) (*model.AnnouncementCode, error)
	Create(ab model.CreateAnnouncement) (int64, error)
	Update(ab model.UpdateAnnouncement) (int64, error)
	Delete(id string, deleted_by string) (int64, error)
	GetByCode(id string) ([]model.Announcement, error)
	GetByIDandType(id string, types string) (*model.Announcement, error)
	GetAllMin() (*[]model.GetAllAnnouncement, error)
	GetAllANWithFilter(keyword []string) ([]*model.Announcement, error)
	GetAllANWithSearch(keyword string) ([]*model.Announcement, error)
}

type usecase struct {
	anRepo an.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		an.NewRepository(),
	}
}
func (m *usecase) DetailCode(id string) (*model.AnnouncementCode, error) {
	return m.anRepo.GetByIDCode(id)
}
func (m *usecase) GetByCode(id string) ([]model.Announcement, error) {
	return m.anRepo.GetByCode(id)
}
func (m *usecase) GetAllAnnouncement() ([]*model.Announcement, error) {
	return m.anRepo.GetAllAnnouncement()
}
func (m *usecase) Create(an model.CreateAnnouncement) (int64, error) {
	// ab := model.CreateAnnouncement
	return m.anRepo.Create(an)
}
func (m *usecase) Update(an model.UpdateAnnouncement) (int64, error) {
	// ab := model.CreateAnnouncement
	return m.anRepo.Update(an)
}
func (m *usecase) Delete(id string, deleted_by string) (int64, error) {
	return m.anRepo.Delete(id, deleted_by)
}

// Communication Use
func (m *usecase) GetByIDandType(id string, types string) (*model.Announcement, error) {
	return m.anRepo.GetByIDandType(id, types)
}

func (m *usecase) GetAllMin() (*[]model.GetAllAnnouncement, error) {
	return m.anRepo.GetAllMin()
}
func (m *usecase) GetAllANWithFilter(keyword []string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithFilter(keyword)
}
func (m *usecase) GetAllANWithSearch(keyword string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithSearch(keyword)
}
