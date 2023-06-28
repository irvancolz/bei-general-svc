package topic

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	tp "be-idx-tsg/internal/app/httprest/repository/topic"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(keyword, status, name, companyName, startDate, endDate, userId string, page, limit int) ([]*model.Topic, error)
	GetTotal(keyword, status, name, companyName, startDate, endDate, userId string, page, limit int) (int, int, error)
	GetByID(topicID, keyword string) (*model.Topic, error)
	UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error)
	UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error)
	CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error)
	CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error)
	DeleteTopic(topicID string, c *gin.Context) (int64, error)
	ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error)
	ExportTopic(c *gin.Context, keyword, status, name, companyName, startDate, userId string) error
}

type usecase struct {
	tpRepo tp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		tp.NewRepository(),
	}
}

func (m *usecase) GetAll(keyword, status, name, companyName, startDate, endDate, userId string, page, limit int) ([]*model.Topic, error) {
	return m.tpRepo.GetAll(keyword, status, name, companyName, startDate, endDate, userId, page, limit)
}

func (m *usecase) GetTotal(keyword, status, name, companyName, startDate, endDate, userId string, page, limit int) (int, int, error) {
	return m.tpRepo.GetTotal(keyword, status, name, companyName, startDate, endDate, userId, page, limit)
}

func (m *usecase) GetByID(topicID, keyword string) (*model.Topic, error) {
	return m.tpRepo.GetByID(topicID, keyword)
}

func (m *usecase) UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error) {
	return m.tpRepo.UpdateHandler(topic, c)
}

func (m *usecase) UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error) {
	return m.tpRepo.UpdateStatus(topic, c)
}

func (m *usecase) CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error) {
	return m.tpRepo.CreateTopicWithMessage(topic, c, isDraft)
}

func (m *usecase) CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error) {
	return m.tpRepo.CreateMessage(message, c)
}

func (m *usecase) DeleteTopic(topicID string, c *gin.Context) (int64, error) {
	return m.tpRepo.DeleteTopic(topicID, c)
}

func (m *usecase) ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error) {
	return m.tpRepo.ArchiveTopicToFAQ(topic, c)
}

func (m *usecase) ExportTopic(c *gin.Context, keyword, status, name, companyName, startDate, userId string) error {
	exportedField := []string{
		"name",
		"company",
		"message",
		"date",
		"status",
	}

	tableHeader := []string{
		"Nama",
		"Nama Perusahaan",
		"Pertanyaan",
		"Waktu Pertanyaan",
		"Status",
	}

	var dataToExported [][]string
	var topicList []*model.Topic
	dataToExported = append(dataToExported, tableHeader)

	topicList, _ = m.GetAll(keyword, status, name, companyName, startDate, "", userId, 0, 0)

	for _, data := range topicList {
		topic := model.TopicExport{
			Name:    data.UserFullName,
			Company: data.CompanyName,
			Message: data.Message,
			Date:    data.CreatedAt.Format("2 Jan 2006 - 15:04"),
			Status:  string(data.Status),
		}

		var topicData []string
		topicData = append(topicData, helper.StructToArray(topic, exportedField)...)

		dataToExported = append(dataToExported, topicData)
	}

	excelConfig := helper.ExportToExcelConfig{
		CollumnStart: "b",
	}
	pdfConfig := helper.PdfTableOptions{
		HeaderTitle: "Pertanyaan Jawaban",
	}
	errorCreateFile := helper.ExportTableToFile(c, helper.ExportTableToFileProps{
		Filename:    "pertanyaan_jawaban",
		Data:        dataToExported,
		ExcelConfig: &excelConfig,
		PdfConfig:   &pdfConfig,
	})
	if errorCreateFile != nil {
		return errorCreateFile
	}

	return nil
}
