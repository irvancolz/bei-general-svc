package faq

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(keyword, userId string) ([]*model.FAQ, error)
	CreateFAQ(faq model.CreateFAQ, c *gin.Context, isDraft bool) (int64, error)
	DeleteFAQ(faqID string, c *gin.Context) (int64, error)
	UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error)
	UpdateOrderFAQ(faqs []model.UpdateFAQOrder, c *gin.Context) (int64, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAll(keyword, userId string) ([]*model.FAQ, error) {
	var listData = []*model.FAQ{}

	query := `SELECT id, created_by, created_at, question, answer, status, order_num FROM faqs WHERE is_deleted = false AND (status = 'PUBLISHED' OR (status = 'DRAFT' AND created_by = '` + userId + `'))`

	if keyword != "" {
		keywords := strings.Split(keyword, ",")

		var filterQuery []string

		for _, v := range keywords {
			filterQuery = append(filterQuery, `question ILIKE '%`+v+`%' OR answer ILIKE '%`+v+`%'`)
		}

		query += `AND (` + strings.Join(filterQuery, " OR ") + ")"
	}

	query += ` ORDER BY order_num ASC`

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

func (m *repository) CreateFAQ(faq model.CreateFAQ, c *gin.Context, isDraft bool) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	faq.CreatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	faq.CreatedBy = userId.(string)

	faq.Status = model.PublishedFAQ

	if isDraft {
		faq.Status = model.DraftFAQ
	}
	faq.OrderNum = 0

	query := `INSERT INTO faqs (question, answer, status, created_by, created_at, order_num) VALUES (:question, :answer, :status, :created_by, :created_at, :order_num)`

	result, err := m.DB.NamedExec(query, &faq)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [CreateFAQ] ", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) UpdateStatusFAQ(faq model.UpdateFAQStatus, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	faq.UpdatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	faq.UpdatedBy = userId.(string)

	faq.Status = model.PublishedFAQ

	var count int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM faqs WHERE id = '%s' AND created_by = '%s'`, faq.ID, userId)

	err := m.DB.Get(&count, query)
	if err != nil || count == 0 {
		return 0, errors.New("forbidden")
	}

	query = `SELECT MAX(order_num) + 1 AS max FROM faqs`

	err = m.DB.Get(&faq.OrderNum, query)
	if err != nil {
		return 0, err
	}

	query = `UPDATE faqs SET status = :status, question = :question, answer = :answer, updated_by = :updated_by, updated_at = :updated_at, order_num = :order_num WHERE id = :id`

	result, err := m.DB.NamedExec(query, &faq)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [UpdateStatusFAQ] ", err)
		return 0, nil
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) DeleteFAQ(faqID string, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")

	userId, _ := c.Get("user_id")

	var count int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM faqs WHERE id = '%s' AND created_by = '%s'`, faqID, userId)

	err := m.DB.Get(&count, query)
	if err != nil || count == 0 {
		return 0, errors.New("forbidden")
	}

	topic := model.DeleteFAQ{
		ID:        faqID,
		DeletedAt: t.Format("2006-01-02 15:04:05"),
		DeletedBy: userId.(string),
	}

	query = `UPDATE faqs SET is_deleted = true, deleted_by = :deleted_by, deleted_at = :deleted_at WHERE id = :id`

	result, err := m.DB.NamedExec(query, &topic)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [DeleteFAQ] ", err)
		return 0, nil
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (m *repository) UpdateOrderFAQ(faqs []model.UpdateFAQOrder, c *gin.Context) (int64, error) {
	var IDs []string

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	updatedAt := t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	updatedBy := userId.(string)

	query := `UPDATE faqs SET order_num = CASE `
	for _, faq := range faqs {
		query += fmt.Sprintf(`WHEN id = '%s' THEN '%d' `, faq.ID, faq.OrderNum)

		IDs = append(IDs, fmt.Sprintf(`'%s'`, faq.ID))
	}
	query += fmt.Sprintf(` ELSE order_num END, updated_at = '%s', updated_by = '%s' WHERE id IN (%s)`, updatedAt, updatedBy, strings.Join(IDs, ", "))

	result, err := m.DB.Exec(query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [UpdateOrderFAQ] ", err)
		return 0, nil
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}
