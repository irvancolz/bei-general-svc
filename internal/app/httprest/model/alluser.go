package model

type AllUserRequest struct {
	OrderBy    string                   `json:"orderby"`
	SortBy     string                   `json:"sortby"`
	Search     string                   `json:"search"`
	Pagination AllUserPaginationRequest `json:"pagination" binding:"required"`
}

type AllUserPaginationRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type AllUserResponse struct {
	ID         string `json:"id"`
	UserRoleID string `json:"user_role_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Fullname   string `json:"fullname"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
