package model

import "time"

type FAQStatus string

const (
	PublishedFAQ FAQStatus = "PUBLISHED"
	DraftFAQ     FAQStatus = "DRAFT"
)

type FAQ struct {
	ID                 string    `json:"id"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at"`
	Question           string    `json:"question" db:"question"`
	Answer             string    `json:"answer" db:"answer"`
	Status             FAQStatus `json:"status" db:"status"`
	OrderNum           int       `json:"order_num" db:"order_num"`
}

type CreateFAQ struct {
	Status    FAQStatus `db:"status"`
	Question  string    `json:"question" db:"question" binding:"required"`
	Answer    string    `json:"answer" db:"answer" binding:"required"`
	OrderNum  int       `db:"order_num"`
	CreatedBy string    `db:"created_by"`
	CreatedAt string    `db:"created_at"`
}

type UpdateFAQ struct {
	ID        string `db:"id" json:"id" binding:"required"`
	Question  string `json:"question" db:"question" binding:"required"`
	Answer    string `json:"answer" db:"answer" binding:"required"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt string `db:"updated_at"`
}

type UpdateFAQStatus struct {
	ID        string    `db:"id" json:"id" binding:"required"`
	Status    FAQStatus `db:"status"`
	Question  string    `json:"question" db:"question" binding:"required"`
	Answer    string    `json:"answer" db:"answer" binding:"required"`
	OrderNum  int       `db:"order_num"`
	CreatedBy string    `db:"created_by"`
	CreatedAt string    `db:"created_at"`
	UpdatedBy string    `db:"updated_by"`
	UpdatedAt string    `db:"updated_at"`
}

type UpdateFAQOrder struct {
	ID        string `db:"id" json:"id"`
	OrderNum  int    `db:"order_num" json:"order_num"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt string `db:"updated_at"`
}

type DeleteFAQ struct {
	ID        string `db:"id"`
	DeletedBy string `db:"deleted_by"`
	DeletedAt string `db:"deleted_at"`
}
