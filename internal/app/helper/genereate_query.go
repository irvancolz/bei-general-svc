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
		sb.WriteString(` AND`)
		if i != 0 {
			sb.WriteString(fmt.Sprintf(` %s IN`, columnFlag))
		}
		sb.WriteString(`( `)

		if i != 0 {
			sb.WriteString(fmt.Sprintf(`SELECT %s FROM %s WHERE (`, columnFlag, m.TableName))
		}

		for j, column := range m.ColumnScanned {
			if j != 0 {
				sb.WriteString(` OR `)
			}
			sb.WriteString(fmt.Sprintf(` CAST ( %s AS TEXT ) LIKE '%%%s%%'`, column, keyword))
		}

		if i != 0 {
			sb.WriteString(`)`)
		}
		sb.WriteString(`)`)
	}
	return sb.String()
}

func (s *SearchQueryGenerator) GenerateGetAllDataQuerry(c *gin.Context, baseQuery string) string {
	queryParams := c.QueryArray("search")
	if len(queryParams) > 0 {
		return baseQuery + s.GenerateSearchQuery(queryParams, "")
	}
	return baseQuery
}

func (s *SearchQueryGenerator) GenerateGetAllDataByQueryKeyword(queryParams []string, baseQuery string) string {
	if len(queryParams) > 0 {
		return baseQuery + s.GenerateSearchQuery(queryParams, "")
	}
	return baseQuery
}
