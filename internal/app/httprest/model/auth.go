package model

type User struct {
	ID           string          `json:"id" binding:"required"`
	Name         string          `json:"name"`
	Username     string          `json:"username" binding:"required"`
	Email        string          `json:"email" binding:"required"`
	Fullname     *string         `json:"fullname" binding:"required"`
	Address      *string         `json:"address"`
	Phone        *string         `json:"phone"`
	Position     *string         `json:"position"`
	RoleStatus   *string         `json:"role_status"`
	Telephone    *string         `json:"telephone"`
	Type         string          `json:"type" binding:"required"`
	CompanyId    *string         `json:"company_id" `
	CompanyName  string          `json:"company_name"`
	Status       *string         `json:"status"`
	Password     *string         `json:"password"`
	RoleId       string          `json:"role_id" `
	UserRoleForm []*UserRoleForm `json:"user_role_form" binding:"required"`
	ExternalType string          `json:"external_type"`
	StatusUser   *bool           `json:"status_user"`
	PhotoUrl     *string         `json:"photo_url"`
}
