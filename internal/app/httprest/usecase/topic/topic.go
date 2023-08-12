package topic

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	tp "be-idx-tsg/internal/app/httprest/repository/topic"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(c *gin.Context) ([]*model.Topic, error)
	GetTotal(c *gin.Context) (int, int, error)
	GetByID(topicID, keyword string) (*model.Topic, error)
	UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error)
	UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error)
	CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error)
	CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error)
	DeleteTopic(topicID string, c *gin.Context) (int64, error)
	ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error)
	ExportTopic(c *gin.Context) error
}

type usecase struct {
	tpRepo tp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		tp.NewRepository(),
	}
}

func (m *usecase) GetAll(c *gin.Context) ([]*model.Topic, error) {
	return m.tpRepo.GetAll(c)
}

func (m *usecase) GetTotal(c *gin.Context) (int, int, error) {
	return m.tpRepo.GetTotal(c)
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

func (m *usecase) ExportTopic(c *gin.Context) error {
	dataStruct, errorData := m.tpRepo.GetAll(c)
	if errorData != nil {
		return errorData
	}

	var dataToConverted []map[string]interface{}

	for _, data := range dataStruct {
		topic := map[string]interface{}{
			"name":    data.UserFullName,
			"company": data.CompanyName,
			"message": data.Message,
			"date":    data.CreatedAt.Format("2 Jan 2006 - 15:04"),
			"status":  string(data.Status),
		}

		dataToConverted = append(dataToConverted, topic)
	}

	columnHeaders := []string{"Nama", "Nama Perusahaan", "Pertanyaan", "Waktu Pertanyaan", "Status"}
	columnWidth := []float64{30, 30, 60, 40, 20}

	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	exportedFields := []string{"name", "company", "message", "date", "status"}
	var exportedData [][]string

	for _, content := range dataToConverted {
		var item []string
		item = append(item, helper.MapToArray(content, exportedFields)...)

		exportedData = append(exportedData, item)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename: "Pertanyaan Jawaban",
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Pertanyaan Jawaban"},
		},
		PdfConfig: &helper.PdfTableOptions{
			HeaderRows: helper.GenerateTableHeaders(columnHeaders, columnWidth),
		},
		Data:        exportedData,
		Headers:     tablesColumns,
		ColumnWidth: columnWidthInt,
	}

	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return errorExport
	}

	return nil
}
