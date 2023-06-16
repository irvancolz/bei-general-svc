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
	Question           string    `json:"question"`
	Answer             string    `json:"answer"`
	Status             FAQStatus `json:"status"`
}

type CreateFAQ struct {
	CreatedBy string    `db:"created_by"`
	CreatedAt string    `db:"created_at"`
	Status    FAQStatus `db:"status"`
	Question  string    `json:"question" db:"question"`
	Answer    string    `json:"answer" db:"answer"`
}

type UpdateFAQStatus struct {
	ID        string    `db:"id" json:"id"`
	Status    FAQStatus `db:"status"`
	UpdatedBy string    `db:"updated_by"`
	UpdatedAt string    `db:"updated_at"`
}

type DeleteFAQ struct {
	ID        string `db:"id"`
	DeletedBy string `db:"deleted_by"`
	DeletedAt string `db:"deleted_at"`
}
