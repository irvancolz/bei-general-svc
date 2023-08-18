package topic

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	tp "be-idx-tsg/internal/app/httprest/repository/topic"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(c *gin.Context) (*helper.PaginationResponse, error)
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

func (m *usecase) GetAll(c *gin.Context) (*helper.PaginationResponse, error) {
	results, errorResults := m.tpRepo.GetAll(c)
	if errorResults != nil {
		return nil, errorResults
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, nil)

	startDate := c.Query("start_date")

	if startDate != "" {
		var temp []map[string]interface{}

		for _, data := range filteredData {
			if parseTime(startDate) == data["time_created_at"].(time.Time).Format("2006-01-02") {
				temp = append(temp, data)
			}
		}

		filteredData = temp
	}

	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)

	return &paginatedData, nil
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
	results, errorResults := m.tpRepo.GetAll(c)
	if errorResults != nil {
		return errorResults
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, _ := helper.HandleDataFiltering(c, dataToConverted, nil)

	startDate := c.Query("start_date")

	if startDate != "" {
		var temp []map[string]interface{}

		for _, data := range filteredData {
			if parseTime(startDate) == data["time_created_at"].(time.Time).Format("2006-01-02") {
				temp = append(temp, data)
			}
		}

		filteredData = temp
	}

	columnHeaders := []string{"Nama", "Nama Perusahaan", "Pertanyaan", "Waktu Pertanyaan", "Status"}
	columnWidth := []float64{30, 30, 60, 40, 20}

	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	exportedFields := []string{"user_full_name", "company_name", "message", "created_at", "status"}
	var exportedData [][]string

	for _, content := range filteredData {
		var item []string
		item = append(item, helper.MapToArray(content, exportedFields)...)

		for i, content := range item {
			if i == 3 {
				date, _ := time.Parse("2006-01-02 15:04:05", content[0:19])

				item[i] = date.Format("2 Jan 2006 - 15:04")
			}
		}

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

func parseTime(input string) string {
	// parse input string menjadi time.Time object
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		log.Println("error parsing time:", err)
		return ""
	}

	// set timezone yang diinginkan
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("error loading location:", err)
		return ""
	}

	// konversi time.Time object ke timezone yang diinginkan
	t = t.In(location)

	// format output string
	output := t.Format("2006-01-02")

	return output
}
