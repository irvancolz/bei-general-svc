package log_system

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll() ([]*model.LogSystem, error)
	CreateLogSystem(logSystem model.CreateLogSystem, c *gin.Context) (int64, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAll() ([]*model.LogSystem, error) {
	var listData = []*model.LogSystem{}

	query := `SELECT * FROM log_systems`

	err := m.DB.Select(&listData, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return listData, err
	}

	for i := range listData {
		listData[i].FormattedCreatedAt = listData[i].CreatedAt.Format("2006-01-02 15:04:05")
	}

	return listData, nil
}

func (m *repository) CreateLogSystem(logSystem model.CreateLogSystem, c *gin.Context) (int64, error) {
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	logSystem.CreatedAt = t.Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	logSystem.CreatedBy = userId.(string)

	name, _ := c.Get("name")
	logSystem.UserName = name.(string)

	logSystem.Browser = c.Request.Header.Get("User-Agent")

	logSystem.IP = getClientIpAddress(c.Request)

	query := `INSERT INTO log_systems (menu, action, user_name, ip, browser, created_by, created_at) VALUES (:menu, :action, :user_name, :ip, :browser, :created_by, :created_at)`

	result, err := m.DB.NamedExec(query, &logSystem)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func getClientIpAddress(req *http.Request) string {
	ip := req.Header.Get("X-FORWARDED-FOR")
	if ip != "" {
		return ip
	}
	return req.RemoteAddr
}
