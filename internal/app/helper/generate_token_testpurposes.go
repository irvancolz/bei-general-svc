package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"gopkg.in/dgrijalva/jwt-go.v3"
)


type JWTClaim2 struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Email           string                `json:"email"`
	UserName        string                `json:"user_name"`
	UserRole        string                `json:"user_role"`
	UserRoleID      string                `json:"user_role_id"`
	UserType        string                `json:"user_type"`
	GroupType       string                `json:"group_type"`
	CompanyId       string                `json:"company_id"`
	CompanyName     string                `json:"company_name"`
	CompanyCode     string                `json:"company_code"`
	FirstLogin      bool                  `json:"first_login"`
	PasswordExpired string               `json:"password_expired"`
	UserFormRole    []*model.UserRoleForm `json:"user_form_role"`
	ExternalType    string               `json:"external_type"`
	jwt.StandardClaims
}
