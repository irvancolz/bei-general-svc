package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DateStr string

const (
	QUERY_START_DATE = "start_date"
	QUERY_END_DATE   = "end_date"
    DEFAULT_ORDER_BY = "CASE WHEN updated_at IS NOT NULL THEN updated_at ELSE created_at END DESC"
)

func getCaseInsensitiveQuery(key string) string {
	return "LOWER(CAST(" + key + " AS VARCHAR)) = LOWER(?)"
}

func getCaseInsensitveContainsQuery(key string) string {
	return "LOWER(CAST(" + key + " AS VARCHAR)) LIKE LOWER(?)"
}

func (dateStr DateStr) ToTimeFormat() string {
	return string(dateStr + " 23:59:59")
}

func GetDefaultLimit(c *gin.Context) int {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 50
	}
	return limit
}

func GetDefaultPage(c *gin.Context) int {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 1
	}
	return page
}

func GetDefaultOffset(page int, limit int) int {
	return (page - 1) * limit
}

func GetDefaultOrder(c *gin.Context) string {
	order := c.DefaultQuery("order", "desc")
	return order
}

func GetDefaultOrderBy(c *gin.Context) string {
	orderBy := c.DefaultQuery("orderBy", DEFAULT_ORDER_BY)
	return orderBy
}

func isColumnQuery(key string) bool {
	key = strings.TrimSpace(key)
	return !(strings.EqualFold(key, "limit") || strings.EqualFold(key, "order") || strings.EqualFold(key, "orderBy") || strings.EqualFold(key, "page") || strings.EqualFold(key, "offset"))
}

func isDateTimeQuery(key string) bool {
	key = strings.TrimSpace(key)
	return strings.EqualFold(key, "created_at") || strings.EqualFold(key, "updated_at") || strings.EqualFold(key, "deleted_at")
}

func GetDefaultEndDate(c *gin.Context, now time.Time) time.Time {
	endDate := DateStr(c.DefaultQuery(QUERY_END_DATE, now.Format(time.DateOnly)))
	defaultEndDate, err := time.Parse(time.DateTime, endDate.ToTimeFormat())

	if err != nil {
		defaultEndDate = now
	}

	return defaultEndDate
}

func MatchingData(query, filter, data string) string {
	if data != "" {
		if filter != "" {
			filter = filter + "or"
		} else {
			filter = filter + "where"
		}
		filter = filter + query
	}

	return filter
}

func GetMaxPage(db *gorm.DB, model interface{}, limit int) int {
	var count int64
	countDb := db
	countDb.Model(&model).Count(&count)

	totalPage := float64(float64(count) / float64(limit))

	return int(math.Ceil(totalPage))
}

func GetGormQueryFilter(db *gorm.DB, query url.Values, endDate time.Time, now time.Time) (*gorm.DB, string) {
	var errorStringBuilder strings.Builder

	for key, values := range query {
		//end date udah diolah, skip
		if key == QUERY_END_DATE {
			continue
		}

		if strings.Contains(key, "contains.") {
			//handle keyword
			for _, value := range values {

				decodeValue, err := url.QueryUnescape(value)

				if err != nil {
					println(err.Error())
					continue
				}

				actualKey := strings.Split(key, ".")[1]

				db = db.Or(getCaseInsensitveContainsQuery(actualKey), "%"+decodeValue+"%")
			}
		} else if strings.Contains(key, "data.") {
			//handle keyword

			for _, value := range values {
				decodeValue, err := url.QueryUnescape(value)

				if err != nil {
					log.Println(err)
					errorStringBuilder.WriteString(err.Error())
					continue
				}
				actualKey := strings.Split(key, "_")[1]
				db = db.Where("data LIKE ?", "%"+actualKey+"%").Where("data LIKE ?", "%"+decodeValue+"%")
			}
		} else if key == QUERY_START_DATE {

			startDate, error := time.Parse(time.DateOnly, values[0])

			if error != nil {
				startDate = now.AddDate(-99, 0, 0)
			}

			db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		} else if len(values) > 1 {
			db = db.Where(getCaseInsensitiveQuery(key), values[0])

			for _, value := range values[1:] {
				db = db.Or(getCaseInsensitiveQuery(key), value)
			}
		} else if isDateTimeQuery(key) {

			decodeValue, err := url.QueryUnescape(values[0])
			if err != nil {
				log.Println(err)
				errorStringBuilder.WriteString(err.Error())
				continue
			}
			db = db.Where("TO_CHAR(created_at, 'yyyy-mm-dd') LIKE ?", "%"+decodeValue+"%")

		} else if isColumnQuery(key) {
			decodeValue, err := url.QueryUnescape(values[0])
			if err != nil {
				log.Println(err)
				errorStringBuilder.WriteString(err.Error())
				continue
			}
			db = db.Where(getCaseInsensitiveQuery(key), decodeValue)
		}
	}

	return db, errorStringBuilder.String()
}

func GetFilterQueryParameter(c *gin.Context) (model.FilterQueryParameter) {
	order := GetDefaultOrder(c)
	orderBy := GetDefaultOrderBy(c)
	limit := GetDefaultLimit(c)
	page := GetDefaultPage(c)
	offset := GetDefaultOffset(page, limit)
	now := time.Now()
	endDate := GetDefaultEndDate(c, now)

	filterqueryparameter := model.FilterQueryParameter{
		QueryList: c.Request.URL.Query(),
		Order:     order,
		OrderBy:   orderBy,
		Limit:     limit,
		Offset:    offset,
		EndDate:   endDate,
	}

	return filterqueryparameter
}

