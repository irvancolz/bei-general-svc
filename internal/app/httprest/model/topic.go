package model

import "time"

type TopicStatus string

const (
	NotAnsweredTopic TopicStatus = "BELUM TERJAWAB"
	AnsweredTopic    TopicStatus = "SUDAH TERJAWAB"
	DoneTopic        TopicStatus = "SELESAI TERJAWAB"
	DraftTopic       TopicStatus = "DRAFT"
)

type Topic struct {
	ID                    string         `json:"id" db:"id"`
	Created_By            string         `json:"created_by" db:"created_by"`
	Time_Created_At       time.Time      `json:"-" db:"created_at"`
	Created_At            string         `json:"created_at"`
	Time_Updated_At       time.Time      `json:"-" db:"updated_at"`
	Updated_At            string         `json:"updated_at"`
	User_Full_Name        string         `json:"user_full_name,omitempty" db:"user_full_name"`
	Company_Code          string         `json:"company_code,omitempty" db:"company_code"`
	Company_Name          string         `json:"company_name,omitempty" db:"company_name"`
	Status                string         `json:"status" db:"status"`
	Handler_ID            string         `json:"handler_id" db:"handler_id"`
	Handler_Name          *string        `json:"handler_name" db:"handler_name"`
	Handler_User_Type     string         `json:"handler_user_type,omitempty" db:"handler_user_type"`
	Creator_User_Type     string         `json:"creator_user_type,omitempty" db:"creator_user_type"`
	Creator_External_Type string         `json:"creator_external_type,omitempty" db:"creator_external_type"`
	Messages              []TopicMessage `json:"messages"`
	Message               string         `json:"message" db:"message"`
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
	UserType           string    `json:"user_type" db:"user_type"`
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
	UserType     string      `db:"user_type"`
	ExternalType *string     `db:"external_type"`
	Message      string      `db:"message" json:"message" binding:"required"`
}

type UpdateTopicHandler struct {
	TopicID     string      `db:"topic_id" json:"topic_id"`
	Status      TopicStatus `db:"status"`
	HandlerID   string      `db:"handler_id"`
	HandlerName string      `db:"handler_name"`
	UpdatedBy   string      `db:"updated_by"`
	UpdatedAt   string      `db:"updated_at"`
}

type UpdateTopicStatus struct {
	TopicID   string      `db:"topic_id" json:"topic_id"`
	Status    TopicStatus `db:"status"`
	Message   string      `db:"message" json:"message"`
	UpdatedBy string      `db:"updated_by"`
	UpdatedAt string      `db:"updated_at"`
}

type UpdateTopic struct {
	ID        string `db:"id"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt string `db:"updated_at"`
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
