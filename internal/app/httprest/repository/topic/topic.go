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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(keyword, status, name, company_name, start_date, end_date string, page, limit int) ([]*model.Topic, error)
	GetTotal(keyword, status, name, company_name, start_date, end_date string, page, limit int) (int, int, error)
	GetByID(topicID, keyword string) (*model.Topic, error)
	UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error)
	CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context) (int64, error)
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

func (m *repository) GetAll(keyword, status, name, company_name, start_date, end_date string, page, limit int) ([]*model.Topic, error) {
	var listData = []*model.Topic{}

	query := `SELECT 
	t.id, t.created_by, t.created_at, t.status, COALESCE(t.handler_id, uuid_nil()) AS handler_id, 
	tp.user_full_name, tp.company_name, tp.message
	FROM topics t
	INNER JOIN topic_messages tp ON tp.id = (
		SELECT id FROM topic_messages tp2 WHERE tp2.topic_id = t.id ORDER BY created_at LIMIT 1
	) WHERE t.is_deleted = false`

	if keyword != "" {
		query += ` AND (tp.message ILIKE '%` + keyword + `%' OR tp.company_name ILIKE '%` + keyword + `%'
		OR tp.user_full_name ILIKE '%` + keyword + `%' OR t.status ILIKE '%` + keyword + `%'
		OR t.created_at::text ILIKE '%` + keyword + `%')`
	}

	log.Println(status)

	if status == "BELUM TERJAWAB" || status == "SUDAH TERJAWAB" {
		query += ` AND t.status = '` + status + `'`
	}

	if name != "" {
		query += ` AND tp.user_full_name = '` + name + `'`
	}

	if company_name != "" {
		query += ` AND tp.company_name = '` + company_name + `'`
	}

	if start_date != "" && end_date != "" {
		start_date = parseTime(start_date)
		end_date = parseTime(end_date)

		query += ` AND (tp.created_at BETWEEN '` + start_date + `' AND '` + end_date + `')`
	}

	query += ` ORDER BY created_at DESC`

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit

		query += ` OFFSET ` + strconv.Itoa(offset) + ` LIMIT ` + strconv.Itoa(limit)
	}

	err := m.DB.Select(&listData, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [GetAll] ", err)
		return listData, err
	}

	for i, data := range listData {
		if data.HandlerID == "00000000-0000-0000-0000-000000000000" {
			listData[i].HandlerID = ""
		}

		t, _ := helper.TimeIn(listData[i].CreatedAt, "Asia/Jakarta")
		listData[i].FormattedCreatedAt = t.Format("2006-01-02 15:04")
	}

	return listData, nil
}

func (m *repository) GetTotal(keyword, status, name, company_name, start_date, end_date string, page, limit int) (int, int, error) {
	var totalData int

	query := `SELECT COUNT(t.id)
	FROM topics t
	INNER JOIN topic_messages tp ON tp.id = (
		SELECT id FROM topic_messages tp2 WHERE tp2.topic_id = t.id ORDER BY created_at LIMIT 1
	) WHERE t.is_deleted = false`

	if keyword != "" {
		query += ` AND (tp.message ILIKE '%` + keyword + `%' OR tp.company_name ILIKE '%` + keyword + `%'
		OR tp.user_full_name ILIKE '%` + keyword + `%' OR t.status ILIKE '%` + keyword + `%'
		OR t.created_at::text ILIKE '%` + keyword + `%')`
	}

	if status == "BELUM TERJAWAB" || status == "SUDAH TERJAWAB" {
		query += ` AND t.status = '` + status + `'`
	}

	if name != "" {
		query += ` AND tp.user_full_name = '` + name + `'`
	}

	if company_name != "" {
		query += ` AND tp.company_name = '` + company_name + `'`
	}

	if start_date != "" && end_date != "" {
		start_date = parseTime(start_date)
		end_date = parseTime(end_date)

		query += ` AND (tp.created_at BETWEEN '` + start_date + `' AND '` + end_date + `')`
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

	query := fmt.Sprintf(`SELECT id, created_by, created_at, status, handler_id FROM topics WHERE id = %s AND is_deleted = false`, topicID)
	err := m.DB.Get(&data, query)
	if err != nil {
		return &data, errors.New("not found")
	}

	t, _ := helper.TimeIn(data.CreatedAt, "Asia/Jakarta")
	data.FormattedCreatedAt = t.Format("2006-01-02 15:04")

	query = fmt.Sprintf(`SELECT id, created_by, message, company_id, company_name, user_full_name, created_at FROM topic_messages WHERE topic_id = %s`, topicID)

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

		t, _ := helper.TimeIn(data.Messages[i].CreatedAt, "Asia/Jakarta")
		data.Messages[i].FormattedCreatedAt = t.Format("2006-01-02 15:04")
	}

	return &data, nil
}

func (m *repository) UpdateHandler(topic model.UpdateTopicHandler, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topic.UpdatedAt = t.Format("2006-01-02 15:04:05")

	handlerId, _ := c.Get("user_id")
	topic.HandlerID = handlerId.(string)
	topic.UpdatedBy = handlerId.(string)

	var count int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM topics WHERE id = %s AND handler_id != NULL`, topic.TopicID)

	err := m.DB.Get(&count, query)
	if err != nil || count != 0 {
		return 0, errors.New("topik sudah dihandle")
	}

	query = `UPDATE topics SET handler_id = :handler_id, updated_at = :updated_at, updated_by = :updated_by WHERE id = :topic_id`

	result, err := m.DB.NamedExec(query, &topic)

	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [UpdateHandler] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) CreateTopicWithMessage(topic model.CreateTopicWithMessage, c *gin.Context) (int64, error) {
	topic.Status = model.NotAnswered

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	topic.CreatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	topic.CreatedBy = userId.(string)

	companyId, _ := c.Get("company_id")
	topic.CompanyID = companyId.(string)

	companyName, _ := c.Get("company_name")
	topic.CompanyName = companyName.(string)

	name, _ := c.Get("name")
	topic.UserFullName = name.(string)

	query := `INSERT INTO topics (status, created_by, created_at) VALUES (:status, :created_by, :created_at) RETURNING id AS topic_id`
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

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = %s`, message.TopicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")
	message.CreatedBy = userId.(string)

	if userId.(string) != data.CreatedBy && userId.(string) != data.HandlerID {
		return 0, errors.New("forbidden")
	}

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

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) DeleteTopic(topicID string, c *gin.Context) (int64, error) {
	var data model.Topic

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = %s`, topicID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.CreatedBy && userId.(string) != data.HandlerID {
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

	query := fmt.Sprintf(`SELECT created_by, COALESCE(handler_id, uuid_nil()) AS handler_id FROM topics WHERE id = %s`, topicFAQ.ID)

	err := m.DB.Get(&data, query)
	if err != nil {
		return 0, errors.New("not found")
	}

	userId, _ := c.Get("user_id")

	if userId.(string) != data.CreatedBy && userId.(string) != data.HandlerID {
		return 0, errors.New("forbidden")
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")

	topicFAQ.Status = model.Answered

	topicFAQ.UpdatedAt = t.Format("2006-01-02 15:04:05")

	topicFAQ.UpdatedBy = userId.(string)

	query = `UPDATE topics SET status = :status, updated_by = :updated_by, updated_at = :updated_at WHERE id = :id`

	_, err = m.DB.NamedExec(query, &topicFAQ)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Topic] [sqlQuery] [ArchiveTopicToFAQ] ", err)
		return 0, err
	}

	topicFAQ.CreatedAt = t.Format("2006-01-02 15:04:05")

	topicFAQ.CreatedBy = userId.(string)

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
	output := t.Format("2006-01-02 15:04:05.999 -0700")

	return output
}