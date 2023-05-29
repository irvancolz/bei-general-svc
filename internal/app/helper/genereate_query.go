package helper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchQueryGenerator struct {
	ColumnScanned []string
	TableName     string
}

func (m *SearchQueryGenerator) GenerateSearchQuery(valueExpected []string, constraint string) string {

	var columnFlag string
	if constraint == "" {
		columnFlag = "id"
	} else {
		columnFlag = constraint
	}

	if len(m.ColumnScanned) <= 0 {
		log.Println("failed to create search query : please specify the collumn name you want to search")
		return ""
	}

	if len(valueExpected) <= 0 {
		log.Println("failed to create search query : there is no data to search for")
		return ""
	}

	if len(m.TableName) <= 0 {
		log.Println("failed to create search query : there is no table to search for")
		return ""
	}

	var sb strings.Builder
	for i, keyword := range valueExpected {
		sb.WriteString("\nAND")
		if i != 0 {
			sb.WriteString(fmt.Sprintf(" %s IN", columnFlag))
		}
		sb.WriteString("(\n")

		if i != 0 {
			sb.WriteString(fmt.Sprintf("SELECT %s FROM %s WHERE (\n", columnFlag, m.TableName))
		}

		for j, column := range m.ColumnScanned {
			if j != 0 {
				sb.WriteString("\nOR ")
			}
			sb.WriteString(fmt.Sprintf(`LOWER( %s ) LIKE LOWER('%%%s%%')`, column, keyword))

		}

		if i != 0 {
			sb.WriteString(")")
		}
		sb.WriteString(")")
	}
	return sb.String()
}

func (s *SearchQueryGenerator) GenerateGetAllDataQuerry(c *gin.Context, baseQuery string) string {
	queryParams := c.Request.URL.Query()

	for key, values := range queryParams {
		if key == "search" {
			return baseQuery + s.GenerateSearchQuery(values, "")
		}
	}
	return baseQuery
}
