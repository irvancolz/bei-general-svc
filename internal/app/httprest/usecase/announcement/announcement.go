package announcement

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	an "be-idx-tsg/internal/app/httprest/repository/announcement"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAllAnnouncement(c *gin.Context) ([]*model.Announcement, error)
	Detail(id string, c *gin.Context) (*model.Announcement, error)
	Create(ab model.CreateAnnouncement, c *gin.Context) (int64, error)
	Update(ab model.UpdateAnnouncement, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
	Export(c *gin.Context, id string) error
	GetAllANWithFilter(keyword []string) ([]*model.Announcement, error)
	GetAllANWithSearch(keyword string, InformationType string, startDate string, endDate string) ([]*model.Announcement, error)
}

type usecase struct {
	anRepo an.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		an.NewRepository(),
	}
}
func (m *usecase) Detail(id string, c *gin.Context) (*model.Announcement, error) {
	return m.anRepo.GetByID(id, c)
}
func (m *usecase) GetAllAnnouncement(c *gin.Context) ([]*model.Announcement, error) {
	return m.anRepo.GetAllAnnouncement(c)
}
func (m *usecase) Create(an model.CreateAnnouncement, c *gin.Context) (int64, error) {
	// ab := model.CreateAnnouncement
	return m.anRepo.Create(an, c)
}
func (m *usecase) Update(an model.UpdateAnnouncement, c *gin.Context) (int64, error) {
	return m.anRepo.Update(an, c)
}
func (m *usecase) Delete(id string, c *gin.Context) (int64, error) {
	return m.anRepo.Delete(id, c)
}
func (m *usecase) GetAllANWithFilter(keyword []string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithFilter(keyword)
}
func (m *usecase) GetAllANWithSearch(keyword string, InformationType string, startDate string, endDate string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithSearch(InformationType, keyword, startDate, endDate)
}

func (m *usecase) Export(c *gin.Context, id string) error {
	var exportedData *model.Announcement

	exportedData, errorGetData := m.Detail(id, c)
	if errorGetData != nil {
		return errorGetData
	}

	exportConfig := helper.ExportAnnouncementsToFileProps{
		Filename: "Announcements",
		Data:     *exportedData,
		PdfConfig: helper.PdfTableOptions{
			HeaderTitle: "Pengumuman",
		},
	}
	errExport := helper.ExportAnnouncementsToFile(c, exportConfig)
	if errExport != nil {
		return errExport
	}

	return nil
}
