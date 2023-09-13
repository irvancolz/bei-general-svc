package helper

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	QUERY_FORM_CODE             = "form_code"
	CONST_AUTH_NAME             = "name"
	CONST_AUTH_KEY_USER_ID      = "user_id"
	CONST_AUTH_KEY_NAME_USER    = "name_user"
	CONST_AUTH_KEY_USER_ROLE    = "user_role"
	CONST_AUTH_KEY_USER_ROLE_ID = "user_role_id"
	CONST_AUTH_KEY_COMPANY_ID   = "company_id"
	CONST_AUTH_KEY_COMPANY_NAME = "company_name"
	CONST_AUTH_KEY_COMPANY_CODE = "company_code"
)

func GetSafeString(data interface{}) string {
	if data == nil {
		return ""
	}
	return data.(string)
}

func GetAuthValue(c *gin.Context, key string) string {
	authValue, _ := c.Get(key)
	return GetSafeString(authValue)
}

func GetUserNameInfo(c *gin.Context) string {
	userName := GetAuthValue(c, CONST_AUTH_KEY_NAME_USER)
	if len(userName) > 0 {
		return "User " + userName + ". "
	}
	return userName
}

func GetFormCodeInfo(c *gin.Context) string {
	formCode := GetFormCode(c)
	if len(formCode) > 0 {
		return "Form Code " + formCode + ". "
	}
	return formCode
}

func GetRequestInfo(c *gin.Context) string {
	userNameInfo := GetUserNameInfo(c)
	messageInfo := "."
	if len(userNameInfo) > 0 {
		messageInfo = " " + userNameInfo
	}
	return messageInfo + GetFormCodeInfo(c)
}

func GetFormCode(c *gin.Context) string {
	return strings.ToUpper(c.Query(QUERY_FORM_CODE))
}

func GetRequestQueryValueUserId(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_KEY_USER_ID)
}

func GetRequestQueryValueUserName(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_KEY_NAME_USER)
}

func GetRequestQueryValueCompanyCode(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_KEY_COMPANY_CODE)
}

func GetRequestQueryValueCompanyId(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_KEY_COMPANY_ID)
}

func GetRequestQueryValueCompanyName(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_KEY_COMPANY_NAME)
}

func GetRequestQueryValueNameStr(c *gin.Context) string {
	return c.DefaultQuery(CONST_AUTH_KEY_USER_ID, "")
}

func GetRequestQueryValueName(c *gin.Context) string {
	return GetAuthValue(c, CONST_AUTH_NAME)
}

func ValidateUserFormRole(inType, inExternalType, inFormRole, inDivision string) (int, error)  {
	if len(inType) == 0 {
		return 0, errors.New("didn't have type")
	}

	if len(inExternalType) == 0 {
		return 0, errors.New("didn't have external type)")
	}

	if len(inFormRole) == 0 {
		return 0, errors.New("didn't have form role")
	}

	if len(inDivision) == 0 {
		return 0, errors.New("didn't have division")
	}
	return 0, nil
}