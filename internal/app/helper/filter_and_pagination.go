package helper

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ConvertUnixStrToDateString(unix string, format string) string {
	convertFormat := "2006-01-02"
	if format != "" {
		convertFormat = format
	}

	unixDate := ConvertUnixStrToNumber(unix)
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

func ConvertUnixStrToNumber(str string) int64 {
	result, errorResult := strconv.Atoi(str)
	if errorResult != nil {
		return 0
	}
	return int64(result)
}

func generateTimeRangeParamsName(params string) [2]string {
	if params == "" {
		return [2]string{}
	}
	return [2]string{params + "_from", params + "_end"}
}

func getBaseTimeRange(param string) string {
	paramBase := strings.Split(param, "_")
	return strings.Join(paramBase[:len(paramBase)-1], "_")
}

func generateTimeRangeParamList(params []string) []string {
	var result []string
	for _, param := range params {
		timeRangeParam := generateTimeRangeParamsName(param)
		result = append(result, timeRangeParam[0], timeRangeParam[1])
	}
	return result
}

func isStarterTimeParam(key string) bool {
	params := strings.Split(key, "_")
	return params[len(params)-1] == "from"
}

func CheckIsOnSpecifiedTimeRange(c *gin.Context, key string, target int64) bool {
	baseTimeParam := getBaseTimeRange(key)
	rangeTimeParam := generateTimeRangeParamsName(baseTimeParam)
	timeStart := ConvertUnixStrToNumber(c.Query(rangeTimeParam[0]))
	timeEnd := ConvertUnixStrToNumber(c.Query(rangeTimeParam[1]))

	if timeEnd == 0 {
		return target >= timeStart
	}

	return target >= timeStart && target <= timeEnd
}

// filtering data obtained from database with comparing value on response object and params given by user
// the time field used specially if we want to search data created/updated on a single day
// there is excluded properties that will not filtered = "page", "limit", "search", "export", "orientation"
func HandleDataFiltering(c *gin.Context, data []interface{}, timeField []string) []map[string]interface{} {
	querries := c.Request.URL.Query()
	rangeTimeParams := generateTimeRangeParamList(timeField)
	results := ConvertToMap(data)
	if len(querries) <= 0 {
		return results
	}

	var filteredResults []map[string]interface{}
	for _, maps := range results {
		mapKeys := GetMapKeys(maps)
		var isMatched []bool
		for key := range querries {
			if IsContains(mapKeys, key) {
				isMatched = append(isMatched, IsContains(querries[key], fmt.Sprintf("%v", maps[key])))
			} else if IsContains(timeField, key) {
				isMatched = append(isMatched, ConvertUnixStrToDateString(querries[key][0], "") == ConvertUnixToDateString(maps[key].(int64), ""))
			} else if IsContains(rangeTimeParams, key) {
				isMatched = append(isMatched, CheckIsOnSpecifiedTimeRange(c, key, maps[getBaseTimeRange(key)].(int64)))
			} else {
				isMatched = append(isMatched, true)
			}
		}
		if !IsContains(isMatched, false) {
			filteredResults = append(filteredResults, maps)
		}
	}
	return filteredResults
}

type PaginationResponse struct {
	TotalPage   int                      `json:"total_page"`
	TotalData   int                      `json:"total_data"`
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
	result.TotalData = len(data)

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
