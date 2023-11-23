package log_system

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(c *gin.Context) ([]model.LogSystem, error)
	GetAllWithFilterPagination(c *gin.Context) (*helper.PaginationResponse, error)
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

func (m *repository) GetAll(c *gin.Context) ([]model.LogSystem, error) {
	var listData = []model.LogSystem{}

	query := `SELECT ls.id, ls.modul, COALESCE(ls.sub_modul, '') AS sub_modul, ls.action, COALESCE(ls.detail, '') AS detail,  ls.user_name, ls.ip, ls.browser, ls.created_by, ls.created_at FROM log_systems ls ORDER BY ls.created_at DESC LIMIT 10`

	err := m.DB.Select(&listData, query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return listData, err
	}

	for i := range listData {
		listData[i].Created_At = listData[i].T_Created_At.Format("2006-01-02 15:04:05")
	}

	return listData, nil
}

func (m *repository) GetAllWithFilterPagination(c *gin.Context) (*helper.PaginationResponse, error) {
	var (
		listData     []model.LogSystem
		listFilter   model.LogSystemFilter
		totalData    int
		filteredData []map[string]interface{}
		conditions   []string
		filters      []string
	)

	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	searches := c.QueryArray("search")
	export := c.Query("export")
	listFilter.Modul = model.NullString(c.Query("modul"))
	listFilter.SubModul = model.NullString(c.Query("sub_modul"))
	listFilter.Action = model.NullString(c.Query("action"))
	listFilter.Detail = model.NullString(c.Query("detail"))
	listFilter.User = model.NullString(c.Query("user_name"))
	listFilter.IP = model.NullString(c.Query("ip"))
	createdAtEnd := c.Query("created_at_end")
	createdAtFrom := c.Query("created_at_from")

	selectQuery := `SELECT ls.id, ls.modul, COALESCE(ls.sub_modul, '') AS sub_modul, ls.action, COALESCE(ls.detail, '') AS detail,  ls.user_name, ls.ip, ls.browser, ls.created_by, ls.created_at AS t_created_at FROM log_systems ls `

	countQuery := `SELECT COUNT(*) FROM log_systems ls `

	filterParameterQuery := `SELECT 
		string_agg(DISTINCT modul::text, ';') AS modul,
		string_agg(DISTINCT sub_modul::text, ';') AS sub_modul,
		string_agg(DISTINCT detail::text, ';') AS detail,
		string_agg(DISTINCT action::text, ';') AS action,
		string_agg(DISTINCT user::text, ';') AS user_name,
		string_agg(DISTINCT ip::text, ';') AS ip
	FROM 
		log_systems ls `

	if len(searches) > 0 {
		for _, search := range searches {
			condition := fmt.Sprintf("(ls.modul ILIKE '%%%s%%' OR ls.sub_modul ILIKE '%%%s%%' OR ls.action ILIKE '%%%s%%' OR ls.detail ILIKE '%%%s%%' OR ls.user_name ILIKE '%%%s%%' OR ls.ip ILIKE '%%%s%%')", search, search, search, search, search, search)

			conditions = append(conditions, condition)
		}

		selectQuery += " WHERE (" + strings.Join(conditions, " AND ") + ")"
		countQuery += " WHERE (" + strings.Join(conditions, " AND ") + ")"
		filterParameterQuery += " WHERE (" + strings.Join(conditions, " AND ") + ")"
	}

	for i := 0; i < reflect.TypeOf(listFilter).NumField(); i++ {
		fieldValue := reflect.ValueOf(listFilter).Field(i).String()

		if len(fieldValue) > 0 {
			var filter string

			tag := reflect.TypeOf(listFilter).Field(i).Tag.Get("db")

			if tag == "ip" || tag == "action" || tag == "detail" {
				filter = tag + " ILIKE '%" + fieldValue + "%'"
			} else {
				filter = tag + " = '" + fieldValue + "'"
			}

			filters = append(filters, filter)
		}
	}

	if len(filters) > 0 {
		if !strings.Contains(selectQuery, "WHERE") && !strings.Contains(countQuery, "WHERE") {
			selectQuery += " WHERE "
			countQuery += " WHERE "
		} else {
			selectQuery += " AND "
			countQuery += " AND "
		}

		selectQuery += strings.Join(filters, " AND ")
		countQuery += strings.Join(filters, " AND ")
	}

	if len(createdAtEnd) > 0 && len(createdAtFrom) > 0 {
		if !strings.Contains(selectQuery, "WHERE") && !strings.Contains(countQuery, "WHERE") {
			selectQuery += " WHERE "
			countQuery += " WHERE "
		} else {
			selectQuery += " AND "
			countQuery += " AND "
		}

		from := helper.ConvertUnixStrToDateString(createdAtFrom, "2006-01-02 15:04:05")
		end := helper.ConvertUnixStrToDateString(createdAtEnd, "2006-01-02 15:04:05")

		selectQuery += " ls.created_at >= '" + from + "' AND ls.created_at <= '" + end + "' "
		countQuery += " ls.created_at >= '" + from + "' AND ls.created_at <= '" + end + "' "
	}

	selectQuery += " ORDER BY ls.created_at DESC"

	if limit != 0 && page != 0 && len(export) == 0 {
		selectQuery += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa((page-1)*limit)
	}

	err := m.DB.Select(&listData, selectQuery)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return nil, err
	}

	err = m.DB.Get(&totalData, countQuery)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return nil, err
	}

	err = m.DB.Get(&listFilter, filterParameterQuery)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [FAQ] [sqlQuery] [GetAll] ", err)
		return nil, err
	}

	var dataToConverted []interface{}
	for i := range listData {
		listData[i].Created_At = listData[i].T_Created_At.Format("2006-01-02 15:04:05")

		dataToConverted = append(dataToConverted, listData[i])
	}

	filterParameter := make(map[string][]string)
	for i := 0; i < reflect.TypeOf(listFilter).NumField(); i++ {
		fieldValue := reflect.ValueOf(listFilter).Field(i).String()

		param := strings.Split(fieldValue, ";")

		filterParameter[strings.ToLower(reflect.TypeOf(listFilter).Field(i).Name)] = param
	}

	filteredData = helper.ConvertToMap(dataToConverted)

	paginatedData := helper.HandleDataPaginationFromDB(c, filteredData, filterParameter, totalData)

	return &paginatedData, nil
}

func (m *repository) CreateLogSystem(logSystem model.CreateLogSystem, c *gin.Context) (int64, error) {
	if allowedAction := model.IsAllowedAction(logSystem.Action); !allowedAction {
		return 1, nil
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	logSystem.CreatedAt = t.Format("2006-01-02 15:04:05")

	logSystem.CreatedBy = c.GetString("user_id")

	logSystem.UserName = c.GetString("name_user")

	logSystem.Browser = c.Request.Header.Get("User-Agent")

	logSystem.IP = getClientIpAddress(c.Request)

	query := `INSERT INTO log_systems (modul, sub_modul, action, detail, user_name, ip, browser, created_by, created_at) VALUES (:modul, :sub_modul, :action, :detail, :user_name, :ip, :browser, :created_by, :created_at)`

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
