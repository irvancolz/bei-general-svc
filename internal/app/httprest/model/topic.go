package model

import "time"

type TopicStatus string

const (
	NotAnsweredTopic TopicStatus = "BELUM TERJAWAB"
	AnsweredTopic    TopicStatus = "SUDAH TERJAWAB"
	DraftTopic       TopicStatus = "DRAFT"
)

type Topic struct {
	ID                 string         `json:"id" db:"id"`
	CreatedBy          string         `json:"created_by" db:"created_by"`
	CreatedAt          time.Time      `json:"-" db:"created_at"`
	FormattedCreatedAt string         `json:"created_at"`
	UpdatedAt          time.Time      `json:"-" db:"updated_at"`
	FormattedUpdatedAt string         `json:"pdated_at"`
	UserFullName       string         `json:"user_full_name,omitempty" db:"user_full_name"`
	CompanyCode        string         `json:"company_code,omitempty" db:"company_code"`
	CompanyName        string         `json:"company_name,omitempty" db:"company_name"`
	Status             TopicStatus    `json:"status" db:"status"`
	HandlerID          string         `json:"handler_id" db:"handler_id"`
	HandlerName        *string        `json:"handler_name" db:"handler_name"`
	Message            string         `json:"message" db:"message"`
	Messages           []TopicMessage `json:"messages"`
}

type TopicExport struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Message string `json:"message"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}

type TopicMessage struct {
	ID                 string    `json:"id,omitempty" db:"id"`
	CreatedBy          string    `json:"created_by,omitempty" db:"created_by"`
	Message            string    `json:"message,omitempty" db:"message"`
	CompanyID          string    `json:"company_id,omitempty" db:"company_id"`
	CompanyName        string    `json:"company_name,omitempty" db:"company_name"`
	UserFullName       string    `json:"user_full_name,omitempty" db:"user_full_name"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at,omitempty"`
}

type CreateTopicWithMessage struct {
	Status       TopicStatus `json:"status" db:"status"`
	CreatedBy    string      `db:"created_by"`
	CreatedAt    string      `db:"created_at"`
	CompanyID    string      `db:"company_id"`
	CompanyCode  string      `db:"company_code"`
	CompanyName  string      `db:"company_name"`
	TopicID      string      `db:"topic_id"`
	UserFullName string      `db:"user_full_name"`
	Message      string      `db:"message" json:"message" binding:"required"`
}

type UpdateTopicHandler struct {
	TopicID     string `db:"topic_id" json:"topic_id"`
	HandlerID   string `db:"handler_id"`
	HandlerName string `db:"handler_name"`
	UpdatedBy   string `db:"updated_by"`
	UpdatedAt   string `db:"updated_at"`
}

type UpdateTopicStatus struct {
	TopicID   string      `db:"topic_id" json:"topic_id"`
	Status    TopicStatus `db:"status"`
	Message   string      `db:"message" json:"message"`
	UpdatedBy string      `db:"updated_by"`
	UpdatedAt string      `db:"updated_at"`
}

type DeleteTopic struct {
	ID        string `db:"id"`
	DeletedBy string `db:"deleted_by"`
	DeletedAt string `db:"deleted_at"`
}

type ArchiveTopicToFAQ struct {
	ID        string      `json:"id" db:"id"`
	Status    TopicStatus `db:"status"`
	UpdatedBy string      `db:"updated_by"`
	UpdatedAt string      `db:"updated_at"`
	CreatedBy string      `db:"created_by"`
	CreatedAt string      `db:"created_at"`
	Question  string      `json:"question" db:"question"`
	Answer    string      `json:"answer" db:"answer"`
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
