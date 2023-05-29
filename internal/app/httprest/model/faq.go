package model

import "time"

type FAQ struct {
	ID                 string    `json:"id"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at"`
	Question           string    `json:"question"`
	Answer             string    `json:"answer"`
}

type CreateFAQ struct {
	CreatedBy string `db:"created_by"`
	CreatedAt string `db:"created_at"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
}

type DeleteFAQ struct {
	ID        string `db:"id"`
	DeletedBy string `db:"deleted_by"`
	DeletedAt string `db:"deleted_at"`
}
