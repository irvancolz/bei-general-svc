package model

import "time"

type Status string

const (
	NotAnswered Status = "BELUM TERJAWAB"
	Answered    Status = "SUDAH TERJAWAB"
)

type Topic struct {
	ID                 string         `json:"id" db:"id"`
	CreatedBy          string         `json:"created_by" db:"created_by"`
	CreatedAt          time.Time      `json:"-" db:"created_at"`
	FormattedCreatedAt string         `json:"created_at"`
	UserFullName       string         `json:"user_full_name,omitempty" db:"user_full_name"`
	CompanyName        string         `json:"company_name,omitempty" db:"company_name"`
	Status             Status         `json:"status" db:"status"`
	HandlerID          string         `json:"handler_id" db:"handler_id"`
	Message            string         `json:"message,omitempty" db:"message"`
	Messages           []TopicMessage `json:"messages,omitempty"`
}

type TopicMessage struct {
	ID                 string    `json:"id" db:"id"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	Message            string    `json:"message" db:"message"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	CompanyName        string    `json:"company_name" db:"company_name"`
	UserFullName       string    `json:"user_full_name" db:"user_full_name"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at"`
}

type CreateTopicWithMessage struct {
	Status       Status `json:"status" db:"status"`
	CreatedBy    string `db:"created_by"`
	CreatedAt    string `db:"created_at"`
	CompanyID    string `db:"company_id"`
	CompanyName  string `db:"company_name"`
	TopicID      string `db:"topic_id"`
	UserFullName string `db:"user_full_name"`
	Message      string `db:"message" json:"message" binding:"required"`
}

type UpdateTopicHandler struct {
	TopicID   string `db:"topic_id" json:"topic_id"`
	HandlerID string `db:"handler_id"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt string `db:"updated_at"`
}

type DeleteTopic struct {
	ID        string `db:"id"`
	DeletedBy string `db:"deleted_by"`
	DeletedAt string `db:"deleted_at"`
}

type ArchiveTopicToFAQ struct {
	ID        string `json:"id" db:"id"`
	Status    Status `db:"status"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt string `db:"updated_at"`
	CreatedBy string `db:"created_by"`
	CreatedAt string `db:"created_at"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
}

type CreateMessage struct {
	CreatedBy    string `db:"created_by"`
	CreatedAt    string `db:"created_at"`
	CompanyID    string `db:"company_id"`
	CompanyName  string `db:"company_name"`
	UserFullName string `db:"user_full_name"`
	TopicID      string `json:"topic_id" db:"topic_id" binding:"required"`
	Message      string `json:"message" binding:"required"`
}