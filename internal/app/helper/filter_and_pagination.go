package helper

import (
	"fmt"
	"math"
	"sort"
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

func HandleDataSorting(c *gin.Context, data []map[string]interface{}) []map[string]interface{} {
	sortedField := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	results := data

	if sortedField == "" {
		return data
	}

	if strings.EqualFold(sortOrder, "desc") {
		sort.SliceStable(results, func(current, before int) bool {
			return fmt.Sprintf("%v", results[current][sortedField]) > fmt.Sprintf("%v", results[before][sortedField])
		})
		return results
	}

	sort.SliceStable(results, func(current, before int) bool {
		return strings.ToLower(fmt.Sprintf("%v", results[current][sortedField])) < strings.ToLower(fmt.Sprintf("%v", results[before][sortedField]))
	})
	return results
}

// this func will get all unique values from exported data to users, this will only retrieve
// data thats saved in string type, this will not return another data types.
func generateFilterParameter(data []map[string]interface{}) map[string][]string {
	results := make(map[string][]string)
	if len(data) <= 0 {
		return results
	}
	mapKeys := GetMapKeys(data[0])

	for _, items := range data {
		for _, keys := range mapKeys {
			if IsString(items[keys]) && !IsContains(results[keys], items[keys].(string)) && items[keys].(string) != "" {
				results[keys] = append(results[keys], items[keys].(string))
			}
		}
	}

	return results
}

// filtering data obtained from database with comparing value on response object and params given by user
// the time field used specially if we want to search data created/updated on a single day
// there is excluded properties that will not filtered = "page", "limit", "search", "export", "orientation"
func HandleDataFiltering(c *gin.Context, data []interface{}, timeField []string) (filteredData []map[string]interface{}, filterParameter map[string][]string) {
	querries := c.Request.URL.Query()
	rangeTimeParams := generateTimeRangeParamList(timeField)
	results := ConvertToMap(data)
	filteredParameterResults := generateFilterParameter(results)
	if len(querries) <= 0 {
		return results, filteredParameterResults
	}

	var filteredResults []map[string]interface{}
	for _, maps := range results {
		mapKeys := GetMapKeys(maps)
		var isMatched []bool
		for key := range querries {
			if IsContains(timeField, key) {
				isMatched = append(isMatched, ConvertUnixStrToDateString(querries[key][0], "") == ConvertUnixToDateString(maps[key].(int64), ""))
			} else if IsContains(rangeTimeParams, key) {
				isMatched = append(isMatched, CheckIsOnSpecifiedTimeRange(c, key, maps[getBaseTimeRange(key)].(int64)))
			} else if IsContains(mapKeys, key) && IsString(maps[key]) {
				// convert current obj props value to double quoted str then remove the quote to be compared with current params given
				isMatched = append(isMatched, IsContains(querries[key], strings.ReplaceAll(strconv.Quote(maps[key].(string)), `"`, "")))
			} else if IsContains(mapKeys, key) {
				isMatched = append(isMatched, IsContains(querries[key], fmt.Sprintf("%v", maps[key])))
			} else {
				isMatched = append(isMatched, true)
			}
		}
		if !IsContains(isMatched, false) {
			filteredResults = append(filteredResults, maps)
		}
	}

	sortedResults := HandleDataSorting(c, filteredResults)

	return sortedResults, filteredParameterResults
}

type PaginationResponse struct {
	TotalPage       int                      `json:"total_page"`
	TotalData       int                      `json:"total_data"`
	Data            []map[string]interface{} `json:"data"`
	FilterParameter map[string][]string      `json:"filter_parameter"`
	Next            bool                     `json:"next"`
	Previous        bool                     `json:"previous"`
	CurrentPage     int                      `json:"current_page"`
	Limit           int                      `json:"limit"`
}

func HandleDataPagination(c *gin.Context, data []map[string]interface{}, filterParameter map[string][]string) PaginationResponse {
	var result PaginationResponse
	totalData := len(data)
	pageCount, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if pageCount == 0 {
		pageCount = 1
	}
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if pageLimit == 0 {
		pageLimit = 5
	}
	showedDatafrom := (pageCount - 1) * pageLimit
	pageTotal := float64(totalData) / float64(pageLimit)

	result.FilterParameter = filterParameter
	result.TotalPage = int(math.Ceil(pageTotal))
	result.Limit = pageLimit
	result.CurrentPage = pageCount
	result.Next = true
	result.Previous = true
	result.TotalData = totalData

	if pageCount <= 1 {
		result.Previous = false
	}

	if pageCount*pageLimit > totalData {
		result.Next = false
	}

	if showedDatafrom >= 0 && showedDatafrom < len(data) {
		showedDataEnd := showedDatafrom + pageLimit
		if showedDataEnd > totalData {
			showedDataEnd = totalData
		}
		result.Data = data[showedDatafrom:showedDataEnd]
	}

	if showedDatafrom > totalData || totalData <= 0 {
		result.Data = make([]map[string]interface{}, 0)
	}

	return result
}
