package faq

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(keyword string) ([]*model.FAQ, error)
	CreateFAQ(faq model.CreateFAQ, c *gin.Context) (int64, error)
	DeleteFAQ(faqID string, c *gin.Context) (int64, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAll(keyword string) ([]*model.FAQ, error) {
	var listData = []*model.FAQ{}

	query := `SELECT id, created_by, created_at, question, answer FROM faqs WHERE is_deleted = false`

	if keyword != "" {
		query += ` AND (question ILIKE '%` + keyword + `%' OR answer ILIKE '%` + keyword + `%')`
	}

	query += ` ORDER BY created_at DESC`

	err := m.DB.Select(&listData, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return listData, err
	}

	for i := range listData {
		listData[i].FormattedCreatedAt = listData[i].CreatedAt.Format("2006-01-02 15:04")
	}

	return listData, nil
}

func (m *repository) CreateFAQ(faq model.CreateFAQ, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	faq.CreatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	faq.CreatedBy = userId.(string)

	query := `INSERT INTO faqs (question, answer, created_by, created_at) VALUES (:question, :answer, :created_by, :created_at)`

	result, err := m.DB.NamedExec(query, &faq)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [CreateFAQ] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) DeleteFAQ(faqID string, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")

	userId, _ := c.Get("user_id")

	topic := model.DeleteFAQ{
		ID:        faqID,
		DeletedAt: t.Format("2006-01-02 15:04:05"),
		DeletedBy: userId.(string),
	}

	query := `UPDATE faqs SET is_deleted = true, deleted_by = :deleted_by, deleted_at = :deleted_at WHERE id = :id`

	result, err := m.DB.NamedExec(query, &topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [DeleteFAQ] ", err)
		return 0, nil
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}
