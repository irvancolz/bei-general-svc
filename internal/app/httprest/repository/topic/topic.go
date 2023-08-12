package topic

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(c *gin.Context) ([]model.Topic, error)
	GetTotal(c *gin.Context) (int, int, error)
	GetByID(topicID, keyword string) (*model.Topic, error)
	UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error)
	UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error)
	CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error)
	CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error)
	DeleteTopic(topicID string, c *gin.Context) (int64, error)
	ArchiveTopicToFAQ(topicFAQ model.ArchiveTopicToFAQ, c *gin.Context) (int64, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAll(c *gin.Context) ([]model.Topic, error) {
	keyword := c.Query("keyword")
	userId, _ := c.Get("user_id")
	userType, _ := c.Get("type")

	var listData = []model.Topic{}

	query := `SELECT 
	t.id, t.created_by, t.created_at, COALESCE(tp3.created_at, COALESCE(t.updated_at, t.created_at)) AS updated_at, t.status, COALESCE(t.handler_id, uuid_nil()) AS handler_id, t.handler_name,
	t.company_code, t.company_name, tp.user_full_name, tp.message,
	CASE WHEN tp.company_id != '00000000-0000-0000-0000-000000000000' THEN 'External' ELSE 'Internal' END AS creator_user_type,
	CASE WHEN tp3.company_id != '00000000-0000-0000-0000-000000000000' THEN 'External' ELSE 'Internal' END AS handler_user_type
	FROM topics t
	LEFT JOIN topic_messages tp ON tp.id = (
		SELECT id FROM topic_messages tp2 WHERE tp2.topic_id = t.id ORDER BY created_at LIMIT 1
	) 
	LEFT JOIN topic_messages tp3 ON tp3.id = (
		SELECT id FROM topic_messages tp4 WHERE tp4.topic_id = t.id AND tp4.created_by = COALESCE(t.handler_id, uuid_nil()) ORDER BY created_at DESC LIMIT 1
	) 
	WHERE t.is_deleted = false AND (t.status IN ('SUDAH TERJAWAB', 'BELUM TERJAWAB') OR (t.status = 'DRAFT' AND t.created_by = '` + userId.(string) + `'))`

	if keyword != "" {
		keywords := strings.Split(keyword, ",")

		var filterQuery []string

		for _, v := range keywords {
			filterQuery = append(filterQuery, `tp.message ILIKE '%`+v+`%' OR tp.company_name ILIKE '%`+v+`%'
			OR tp.user_full_name ILIKE '%`+v+`%' OR t.status ILIKE '%`+v+`%'
			OR t.created_at::text ILIKE '%`+v+`%'`)
		}

		query += `AND (` + strings.Join(filterQuery, " OR ") + ")"
	}

	if userType.(string) == "External" {
		companyID, _ := c.Get("company_id")

		query += ` AND tp.company_id = '` + companyID.(string) + `'`
	}

	query += ` ORDER BY CASE WHEN status = 'DRAFT' THEN 1 WHEN status = 'BELUM TERJAWAB' AND handler_id IS NULL THEN 2 CASE WHEN status = 'SUDAH TERJAWAB' THEN 4 ELSE 3 END, t.created_at DESC`

	err := m.DB.Select(&listData, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [GetAll] ", err)
		return listData, err
	}

	for i, data := range listData {
		if data.Handler_ID == "00000000-0000-0000-0000-000000000000" {
			listData[i].Handler_ID = ""
		}
	}

	return listData, nil
}

func (m *repository) GetTotal(c *gin.Context) (int, int, error) {
	keyword := c.Query("keyword")
	limit, _ := strconv.Atoi(c.Query("limit"))
	status := c.Query("status")
	name := c.Query("name")
	companyName := c.Query("company_name")
	startDate := c.Query("start_date")
	userId, _ := c.Get("user_id")
	userType, _ := c.Get("type")

	var totalData int

	query := `SELECT COUNT(t.id)
	FROM topics t
	INNER JOIN topic_messages tp ON tp.id = (
		SELECT id FROM topic_messages tp2 WHERE tp2.topic_id = t.id ORDER BY created_at LIMIT 1
	) WHERE t.is_deleted = false AND (t.status IN ('SUDAH TERJAWAB', 'BELUM TERJAWAB') OR (t.status = 'DRAFT' AND t.created_by = '` + userId.(string) + `'))`

	if keyword != "" {
		keywords := strings.Split(keyword, ",")

		var filterQuery []string

		for _, v := range keywords {
			filterQuery = append(filterQuery, `tp.message ILIKE '%`+v+`%' OR tp.company_name ILIKE '%`+v+`%'
			OR tp.user_full_name ILIKE '%`+v+`%' OR t.status ILIKE '%`+v+`%'
			OR t.created_at::text ILIKE '%`+v+`%'`)
		}

		query += `AND (` + strings.Join(filterQuery, " OR ") + ")"
	}

	if userType.(string) == "External" {
		companyID, _ := c.Get("company_id")

		query += ` AND tp.company_id = '` + companyID.(string) + `'`
	}

	var queryFilter []string

	if status == "BELUM TERJAWAB" || status == "SUDAH TERJAWAB" || status == "DRAFT" {
		statuses := strings.Split(status, ",")

		queryFilter = append(queryFilter, "t.status IN ('"+strings.Join(statuses, "','")+"')")
	}

	if name != "" {
		names := strings.Split(name, ",")

		queryFilter = append(queryFilter, "tp.user_full_name IN ('"+strings.Join(names, "','")+"')")
	}

	if companyName != "" {
		companyNames := strings.Split(companyName, ",")

		queryFilter = append(queryFilter, "tp.company_name IN ('"+strings.Join(companyNames, "','")+"')")
	}

	if startDate != "" {
		startDate = parseTime(startDate)

		queryFilter = append(queryFilter, `t.created_at::TEXT LIKE '`+startDate+`%'`)
	}

	if len(queryFilter) > 0 {
		query += ` AND (` + strings.Join(queryFilter, " OR ") + ")"
	}

	err := m.DB.Get(&totalData, query)
	if err != nil {
		return 0, 0, err
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	return totalData, totalPage, nil
}

func (m *repository) GetByID(topicID, keyword string) (*model.Topic, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT id, created_by, created_at, status, COALESCE(handler_id, uuid_nil()) AS handler_id, handler_name, company_code, company_name FROM topics WHERE id = '%s' AND is_deleted = false`, topicID)
	err := m.DB.Get(&data, query)
	if err != nil {
		return &data, errors.New("not found")
	}

	if data.Handler_ID == "00000000-0000-0000-0000-000000000000" {
		data.Handler_ID = ""
	}

	query = fmt.Sprintf(`SELECT id, created_by, message, company_id, company_name, user_full_name, created_at, CASE WHEN company_id != '00000000-0000-0000-0000-000000000000' THEN 'External' ELSE 'Internal' END as user_type FROM topic_messages WHERE topic_id = '%s'`, topicID)

	if keyword != "" {
		query += ` AND (message ILIKE '%` + keyword + `%' OR company_name ILIKE '%` + keyword + `%' OR user_full_name ILIKE '%` + keyword + `%' OR created_at::text ILIKE '%` + keyword + `%')`
	}

	query += ` ORDER BY created_at ASC`

	err = m.DB.Select(&data.Messages, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [GetByID] ", err)
		return &data, err
	}

	for i, message := range data.Messages {
		if message.CompanyID == "00000000-0000-0000-0000-000000000000" {
			data.Messages[i].CompanyID = ""
		}

		data.Messages[i].FormattedCreatedAt = data.Messages[i].CreatedAt.Format("2006-01-02 15:04")
	}

	if len(data.Messages) == 0 {
		data.Messages = []model.TopicMessage{}
	}

	return &data, nil
}

func (m *repository) UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = '%s'`, topic.TopicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	if data.Handler_ID != "00000000-0000-0000-0000-000000000000" {
		return 0, errors.New("forbidden")
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topic.UpdatedAt = t.Format("2006-01-02 15:04:05")

	handler_Id, _ := c.Get("user_id")
	topic.HandlerID = handler_Id.(string)
	topic.UpdatedBy = handler_Id.(string)

	name, _ := c.Get("name")
	topic.HandlerName = name.(string)

	query = `UPDATE topics SET handler_id = :handler_id, handler_name = :handler_name, updated_at = :updated_at, updated_by = :updated_by WHERE id = :topic_id`

	result, err := m.DB.NamedExec(query, &topic)

	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [UpdateHandler] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) UpdateStatus(topic model.UpdateTopicStatus, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = '%s'`, topic.TopicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.Created_By && userId.(string) != data.Handler_ID {
		return 0, errors.New("forbidden")
	}

	topic.UpdatedBy = userId.(string)

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topic.UpdatedAt = t.Format("2006-01-02 15:04:05")

	topic.Status = model.AnsweredTopic
	if strings.Contains(c.FullPath(), "publish") {
		topic.Status = model.NotAnsweredTopic
	}

	query = `UPDATE topics SET status = :status, updated_at = :updated_at, updated_by = :updated_by WHERE id = :topic_id`

	result, err := m.DB.NamedExec(query, &topic)

	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [UpdateHandler] ", err)
		return 0, err
	}

	if topic.Message != "" {
		query = `UPDATE topic_messages SET message = :message, updated_at = :updated_at, updated_by = :updated_by WHERE id = (SELECT id FROM topic_messages tp WHERE tp.topic_id = :topic_id ORDER BY created_at LIMIT 1)`

		result, err = m.DB.NamedExec(query, &topic)

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [UpdateHandler] ", err)
			return 0, err
		}
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context, isDraft bool) (int64, error) {
	topic.Status = model.NotAnsweredTopic
	if isDraft {
		topic.Status = model.DraftTopic
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topic.CreatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	topic.CreatedBy = userId.(string)

	companyId, _ := c.Get("company_id")
	topic.CompanyID = companyId.(string)

	companyCode, _ := c.Get("company_code")
	topic.CompanyCode = companyCode.(string)
	if topic.CompanyCode == "" {
		topic.CompanyCode = "BEI"
	}

	companyName, _ := c.Get("company_name")
	topic.CompanyName = companyName.(string)
	if topic.CompanyName == "" {
		topic.CompanyName = "Bursa Efek Indonesia"
	}

	name, _ := c.Get("name")
	topic.UserFullName = name.(string)

	if topic.CompanyID == "" {
		topic.CompanyID = "00000000-0000-0000-0000-000000000000"
	}

	query := `INSERT INTO topics (status, created_by, created_at, company_code, company_name) VALUES (:status, :created_by, :created_at, :company_code, :company_name) RETURNING id AS topic_id`
	row, err := m.DB.NamedQuery(query, topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [CreateTopicWithMessage] ", err)
		return 0, err
	}

	for row.Next() {
		row.StructScan(&topic)
	}

	query = `INSERT INTO topic_messages (topic_id, message, company_id, company_name, user_full_name, created_by, created_at) 
	VALUES (:topic_id, :message, :company_id, :company_name, :user_full_name, :created_by, :created_at)`

	result, err := m.DB.NamedExec(query, &topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [CreateTopicWithMessage] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) CreateMessage(message model.CreateMessage, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = '%s'`, message.TopicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.Created_By && userId.(string) != data.Handler_ID {
		return 0, errors.New("forbidden")
	}

	message.CreatedBy = userId.(string)

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	message.CreatedAt = t.Format("2006-01-02 15:04:05")

	companyId, _ := c.Get("company_id")
	message.CompanyID = companyId.(string)

	companyName, _ := c.Get("company_name")
	message.CompanyName = companyName.(string)

	name, _ := c.Get("name")
	message.UserFullName = name.(string)

	if message.CompanyID == "" {
		message.CompanyID = "00000000-0000-0000-0000-000000000000"
	}

	query = `INSERT INTO topic_messages (topic_id, message, company_id, company_name, user_full_name, created_by, created_at) 
	VALUES (:topic_id, :message, :company_id, :company_name, :user_full_name, :created_by, :created_at)`

	result, err := m.DB.NamedExec(query, &message)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [CreateMessage] ", err)
		return 0, err
	}

	recipientID := data.Created_By
	title := fmt.Sprintf("%s \u0020 telah membalas jawaban pertanyaan anda", name.(string))
	if userId.(string) == data.Created_By {
		recipientID = data.Handler_ID
		title = fmt.Sprintf("%s \u0020 telah membalas pertanyaan anda", name.(string))
	}

	param := helper.CreateSingleNotificationParam{
		UserID: recipientID,
		Data: helper.NotificationData{
			Title: title,
			Date:  helper.GetCurrentTime(),
		},
		Link: message.TopicID,
		Type: "Topic",
	}

	helper.CreateSingleNotification(c, param)

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) DeleteTopic(topicID string, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = '%s'`, topicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, err
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.Created_By && userId.(string) != data.Handler_ID {
		return 0, errors.New("forbidden")
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")

	topic := model.DeleteTopic{
		ID:        topicID,
		DeletedAt: t.Format("2006-01-02 15:04:05"),
		DeletedBy: userId.(string),
	}

	query = `UPDATE topics SET is_deleted = true, deleted_by = :deleted_by, deleted_at = :deleted_at WHERE id = :id`

	_, err = m.DB.NamedExec(query, &topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [DeleteTopic] ", err)
		return 0, nil
	}

	query = `UPDATE topic_messages SET is_deleted = true, deleted_by = :deleted_by, deleted_at = :deleted_at WHERE topic_id = :id`

	result, err := m.DB.NamedExec(query, &topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [DeleteTopic] ", err)
		return 0, nil
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) ArchiveTopicToFAQ(topicFAQ model.ArchiveTopicToFAQ, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = '%s'`, topicFAQ.ID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.Created_By && userId.(string) != data.Handler_ID {
		return 0, errors.New("forbidden")
	}

	topicFAQ.CreatedBy = userId.(string)

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topicFAQ.CreatedAt = t.Format("2006-01-02 15:04:05")

	query = `INSERT INTO faqs (question, answer, created_by, created_at) VALUES (:question, :answer, :created_by, :created_at)`

	result, err := m.DB.NamedExec(query, &topicFAQ)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [ArchiveTopicToFAQ] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
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

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
