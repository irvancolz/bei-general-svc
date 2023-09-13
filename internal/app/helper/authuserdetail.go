package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"errors"

	"github.com/gin-gonic/gin"
)


func getError(message string) error {
	return errors.New("key " + message + " tidak ada")
}

func GetAuthUserDetail(c *gin.Context) (model.AuthUserDetail, error) {
	authUserDetail  := model.AuthUserDetail{}

	user_id, isUserIdExisting := c.Get("user_id")
	if !isUserIdExisting {
		return authUserDetail, getError("user_id")
	}
	authUserDetail.ID = user_id.(string)

	name_user, isUserNameExisting := c.Get("name_user")
	if !isUserNameExisting {
		return authUserDetail, getError("name_user")
	}
	*authUserDetail.Name = name_user.(string)

	user_role, isUserRoleExisting := c.Get("user_role")
	if !isUserRoleExisting {
		return authUserDetail, getError("user_role")
	}
	*authUserDetail.UserRole = user_role.(string)

	user_role_id, isUserRoleExisting := c.Get("user_role_id")
	if !isUserRoleExisting {
		return authUserDetail, getError("user_role_id")
	}
	*authUserDetail.UserRoleID = user_role_id.(string)

	company_name, _ := c.Get("company_name")
	*authUserDetail.CompanyName = GetSafeString(company_name)

	company_id, _ := c.Get("company_id")
	*authUserDetail.CompanyId = GetSafeString(company_id)

	company_code, _ := c.Get("company_code")
	*authUserDetail.CompanyCode =  GetSafeString(company_code)

	externalType, _ := c.Get("external_type")
	*authUserDetail.ExternalType =  GetSafeString(externalType)


	return authUserDetail, nil
}