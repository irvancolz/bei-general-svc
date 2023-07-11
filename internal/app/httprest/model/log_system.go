package model

import "time"

type LogSystem struct {
	ID                 string    `json:"id" db:"id"`
	Menu               string    `json:"menu" db:"menu"`
	Action             string    `json:"action" db:"action"`
	UserName           string    `json:"user_name" db:"user_name"`
	IP                 string    `json:"ip" db:"ip"`
	Browser            string    `json:"browser" db:"browser"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at"`
}

type CreateLogSystem struct {
	ID        string `db:"id"`
	Menu      string `json:"menu" db:"menu" binding:"required"`
	Action    string `json:"action" db:"action" binding:"required"`
	UserName  string `db:"user_name"`
	IP        string `db:"ip"`
	Browser   string `db:"browser"`
	CreatedBy string `db:"created_by"`
	CreatedAt string `db:"created_at"`
}
