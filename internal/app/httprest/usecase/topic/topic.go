package topic

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	tp "be-idx-tsg/internal/app/httprest/repository/topic"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll(c *gin.Context) (*helper.PaginationResponse, error)
	GetByID(c *gin.Context, topicID, keyword string) (*model.Topic, error)
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

	createdAtFrom := c.Query("created_at_from")
	createdAtEnd := c.Query("created_at_end")

	if len(createdAtEnd) > 0 && len(createdAtFrom) > 0 {
		from := helper.ConvertUnixStrToTime(createdAtFrom)
		end := helper.ConvertUnixStrToTime(createdAtEnd)

		var temp []map[string]interface{}

		for _, data := range filteredData {
			createdAt := data["time_created_at"].(time.Time)

			if createdAt.After(from) && createdAt.Before(end) {
				temp = append(temp, data)
			}
		}

		filteredData = temp
	}

	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)

	return &paginatedData, nil
}

func (m *usecase) GetByID(c *gin.Context, topicID, keyword string) (*model.Topic, error) {
	return m.tpRepo.GetByID(c, topicID, keyword)
}

func (m *usecase) UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error) {
	return m.tpRepo.UpdateHandler(topic, c)
}

func (m *usecase) UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error) {
	updater, _ := c.Get("name_user")
	updateStatRes, errUpdate := m.tpRepo.UpdateStatus(topic, c)
	if errUpdate != nil {
		return 0, errUpdate
	}

	topicCreatorId := m.tpRepo.GetCreator(c, topic.TopicID)
	topicCreator, errGetTopicCreator := email.GetUser(c, topicCreatorId)
	if errGetTopicCreator != nil {
		return 0, errGetTopicCreator
	}

	go utilities.CreateNotif(c, topicCreatorId, "Pertanyaan", "Pertanyaan telah ditandai sebagai terjawab")
	go email.SendEmailNotification(*topicCreator, "Update Status pada pertanyaan anda", "Pertanyaan anda berhasil terjawab")

	internalBursaUser := email.GetAllUserInternalBursa(c)
	var internalBursaUserId []string
	for _, user := range internalBursaUser {
		internalBursaUserId = append(internalBursaUserId, user.Id)
	}
	go utilities.CreateGroupNotif(c, internalBursaUserId, "Pertanyaan", "Pertanyaan telah ditandai sebagai terjawab")

	for _, user := range internalBursaUser {
		go email.SendEmailNotification(user, "Aktivitas Baru Di Menu Pertanyaan", fmt.Sprintf("user %s menandai pertanyaan sebagai sudah terjawab", updater.(string)))
	}

	return updateStatRes, nil
}

func (m *usecase) CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error) {
	notifCreatorId, _ := c.Get("user_id")
	notifCreatorUserName, _ := c.Get("name_user")
	notifCreatorEmail, _ := c.Get("email")

	createTopic, errCreateTopic := m.tpRepo.CreateTopicWithMessage(topic, c, isDraft)
	if errCreateTopic != nil {
		return 0, errCreateTopic
	}

	go utilities.CreateNotif(c, notifCreatorId.(string), "Pertanyaan", "Pertanyaan Berhasil Dibuat")
	go email.SendEmailNotification(model.UsersIdWithEmail{Id: notifCreatorId.(string), Username: notifCreatorUserName.(string), Email: notifCreatorEmail.(string)}, "Pertanyaan Berhasil Dibuat", "Pertanyaan Berhasil Dibuat")

	internalBursaUser := email.GetAllUserInternalBursa(c)
	var internalBursaUserId []string
	for _, user := range internalBursaUser {
		internalBursaUserId = append(internalBursaUserId, user.Id)
	}
	go utilities.CreateGroupNotif(c, internalBursaUserId, "Pertanyaan", fmt.Sprintf("user %s menambahkan pertanyaan baru", notifCreatorUserName.(string)))

	for _, user := range internalBursaUser {
		go email.SendEmailNotification(user, "Pertanyaan Berhasil Dibuat", fmt.Sprintf("user %s menambahkan pertanyaan baru", notifCreatorUserName.(string)))
	}

	return createTopic, nil
}

func (m *usecase) CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error) {
	createmsgRes, errCreatemsg := m.tpRepo.CreateMessage(message, c)
	if errCreatemsg != nil {
		return 0, errCreatemsg
	}

	topicCreatorId := m.tpRepo.GetCreator(c, message.TopicID)
	topicCreator, errGetTopicCreator := email.GetUser(c, topicCreatorId)
	if errGetTopicCreator != nil {
		return 0, errGetTopicCreator
	}

	go utilities.CreateNotif(c, topicCreatorId, "Pertanyaan", "Aktivitas baru di pertanyaan anda")
	go email.SendEmailNotification(*topicCreator, "Pertanyaan Anda Di Response", "Aktivitas baru di pertanyaan anda")

	internalBursaUser := email.GetAllUserInternalBursa(c)
	var internalBursaUserId []string
	for _, user := range internalBursaUser {
		internalBursaUserId = append(internalBursaUserId, user.Id)
	}
	go utilities.CreateGroupNotif(c, internalBursaUserId, "Pertanyaan", "aktivitas baru di menu pertanyaan")

	for _, user := range internalBursaUser {
		go email.SendEmailNotification(user, "Aktivitas Baru Di Menu Pertanyaan", fmt.Sprintf("user %s merespons pada menu pertanyaan ", topicCreator.Username))
	}

	return createmsgRes, nil
}

func (m *usecase) DeleteTopic(topicID string, c *gin.Context) (int64, error) {
	return m.tpRepo.DeleteTopic(topicID, c)
}

func (m *usecase) ArchiveTopicToFAQ(topic model.ArchiveTopicToFAQ, c *gin.Context) (int64, error) {
	data, err := m.tpRepo.ArchiveTopicToFAQ(topic, c)
	if err != nil {
		return 0, err
	}

	internalBursaUser := email.GetAllUserInternalBursa(c)
	var internalBursaUserId []string
	for _, user := range internalBursaUser {
		internalBursaUserId = append(internalBursaUserId, user.Id)
	}
	go utilities.CreateGroupNotif(c, internalBursaUserId, "FAQ", fmt.Sprintf("User %s menambahkan FAQ baru", c.GetString("name_user")))

	for _, user := range internalBursaUser {
		go email.SendEmailNotification(user, "Aktivitas Baru Di Menu FAQ", fmt.Sprintf("User %s menambahkan FAQ baru", c.GetString("name_user")))
	}

	return data, nil

	return data, nil
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
	columnWidth := []float64{30, 30, 175, 40, 30}

	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	exportedFields := []string{"user_full_name", "company_name", "message", "time_created_at", "status"}
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
			PapperWidth:  335,
			Papperheight: 475,
			HeaderRows:   helper.GenerateTableHeaders(columnHeaders, columnWidth),
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
