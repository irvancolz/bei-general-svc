package model

import "time"

type AuthenticationResponse struct {
	ID              string          `json:"id"`
	Name            *string         `json:"name"`
	Email           *string         `json:"email"`
	UserName        *string         `json:"user_name"`
	UserRole        *string         `json:"user_role"`
	UserRoleID      *string         `json:"user_role_id"`
	UserType        *string         `json:"user_type"`
	FirstLogin      bool            `json:"first_login"`
	PasswordExpired *string         `json:"password_expired"`
	CompanyName     *string         `json:"company_name"`
	CompanyCode     *string         `json:"company_code"`
	CompanyId       *string         `json:"company_id"`
	GroupType       *string         `json:"group_type"`
	ExternalType    *string         `json:"external_type"`
	UserRoleForm    []*UserRoleForm `json:"user_form_role"`
}

type AuthUserDetail struct {
	ID              string          `json:"id"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	UserName        string         `json:"user_name"`
	UserRole        string         `json:"user_role"`
	UserRoleID      string         `json:"user_role_id"`
	UserType        string         `json:"user_type"`
	FirstLogin      bool            `json:"first_login"`
	PasswordExpired string         `json:"password_expired"`
	CompanyName     string         `json:"company_name"`
	CompanyCode     string         `json:"company_code"`
	CompanyId       string         `json:"company_id"`
	GroupType       string         `json:"group_type"`
	ExternalType    string         `json:"external_type"`
	UserRoleForm    []UserRoleForm `json:"user_form_role"`
}
type UserRoleForm struct {
	Id           string   `json:"id"`
	UserId       string   `json:"user_id"`
	Type         string   `json:"type"`
	ExternalType string   `json:"external_type"`
	Division     string   `json:"division"`
	GroupId      []string `json:"group_id"`
	FormRole     string   `json:"form_role"`
}

type AuthConfirmPasswordResponse struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	UserName            string     `json:"user_name"`
	Password            string     `json:"password"`
	PasswordResetToken  string     `json:"password_reset_token"`
	CountResetPassword  int        `json:"count_reset_password"`
	CountForgotPassword int    `json:"count_forgot_password"`
	DateForgotPassword  *string    `json:"date_forgot_password"`
	PasswordResetAt     *time.Time `json:"password_reset_at"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type GetUserName struct {
	UserName string `json:"user_name"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type EmailData struct {
	URL     string
	Name    string
	Subject string
}

type SiteVerifyResponse struct {
	Success     *bool     `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type GetAuthResponses struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}