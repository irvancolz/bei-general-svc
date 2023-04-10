package helper

import (
	"fmt"
	"strings"
)

type InsertManyTableUserFormRole struct {
	UserId       string
	FormRole     string
	Type         string
	ExternalType string
	Division     string
	CreatedAt    string
	CreatedBy    string
}

func InsertManyTableUserFormRolev1(tableName string, columns []string, records []InsertManyTableUserFormRole) string {
	var values []string
	for _, record := range records {
		value := fmt.Sprintf("('%s', '%s','%s','%s', '%s', '%s', '%s')", record.UserId, record.FormRole, record.Type, record.Division, record.CreatedBy, record.CreatedAt, record.ExternalType)
		values = append(values, value)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}

func InsertManyTableGroupUserFormRolev1(tableName string, columns []string, UserFormRoleId string, groupId []string, CreatedBy string, CreatedAt string) string {
	var values []string
	for _, record := range groupId {
		value := fmt.Sprintf("('%s', '%s','%s','%s')", UserFormRoleId, record, CreatedBy, CreatedAt)
		values = append(values, value)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}

type InsertManyTableGroupMember struct {
	GroupId   string
	CompanyId string
	CreatedAt string
	CreatedBy string
}

func InsertManyTableGroupMemberv1(tableName string, columns []string, records []string, GroupId string, createdAt string, createdBy string) string {
	var values []string
	for _, record := range records {
		value := fmt.Sprintf("('%s', '%s','%s','%s')", GroupId, record, createdBy, createdAt)
		values = append(values, value)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}
