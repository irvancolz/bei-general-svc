package helper

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ConvertUnixStrToDateString(unix string, format string) string {
	convertFormat := "2006-01-02"
	if format != "" {
		convertFormat = format
	}

	unixDate, _ := strconv.Atoi(unix)
	result := time.Unix(int64(unixDate), 0).Format(convertFormat)

	return result
}

func ConvertUnixToDateString(unix int64, format string) string {
	convertFormat := "2006-01-02"
	if format != "" {
		convertFormat = format
	}

	result := time.Unix(int64(unix), 0).Format(convertFormat)

	return result
}

// filtering data obtained from database with comparing value on response object and params given by user
// the time field used specially if we want to search data created/updated on a single day
func HandleDataFiltering(c *gin.Context, data []interface{}, timeField []string) []map[string]interface{} {
	querries := c.Request.URL.Query()
	results := ConvertToMap(data)
	if len(querries) <= 0 {
		return results
	}

	var filteredResults []map[string]interface{}
	for _, maps := range results {
		mapKeys := GetMapKeys(maps)
		isMatched := make([]bool, len(querries))
		for key := range querries {
			if !IsContains(mapKeys, key) {
				isMatched = append(isMatched, true)
			} else if !IsContains(timeField, key) {
				isMatched = append(isMatched, IsContains(querries[key], fmt.Sprintf("%v", maps[key])))
			} else {
				isMatched = append(isMatched, ConvertUnixStrToDateString(querries[key][0], "") == ConvertUnixToDateString(maps[key].(int64), ""))
			}

		}
		if !IsContains(isMatched, false) {
			filteredResults = append(results, maps)
		}
	}
	return filteredResults
}

type PaginationResponse struct {
	TotalPage   int                      `json:"total_page"`
	Data        []map[string]interface{} `json:"data"`
	Next        bool                     `json:"next"`
	Previous    bool                     `json:"previous"`
	CurrentPage int                      `json:"current_page"`
	Limit       int                      `json:"limit"`
}

func HandleDataPagination(c *gin.Context, data []map[string]interface{}) PaginationResponse {
	var result PaginationResponse

	pageCount, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if pageCount == 0 {
		pageCount = 1
	}
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if pageLimit == 0 {
		pageLimit = 5
	}
	showedDatafrom := (pageCount - 1) * pageLimit
	pageTotal := float64(len(data)) / float64(pageLimit)

	result.TotalPage = int(math.Ceil(pageTotal))
	result.Limit = pageLimit
	result.CurrentPage = pageCount
	result.Next = true
	result.Previous = true

	if pageCount <= 1 {
		result.Previous = false
	}

	if pageCount*pageLimit > len(data) {
		result.Next = false
	}

	if showedDatafrom >= 0 && showedDatafrom < len(data) {
		showedDataEnd := showedDatafrom + pageLimit
		if showedDataEnd > len(data) {
			showedDataEnd = len(data)
		}
		result.Data = data[showedDatafrom:showedDataEnd]
	}

	if showedDatafrom > len(data) || len(data) <= 0 {
		result.Data = make([]map[string]interface{}, 0)
	}

	return result
}
