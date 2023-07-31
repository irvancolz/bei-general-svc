package model

import "time"

type LogSystem struct {
	ID                 string    `json:"id" db:"id"`
	Modul              string    `json:"modul" db:"modul" binding:"required"`
	SubModul           string    `json:"sub_modul" db:"sub_modul"`
	Action             string    `json:"action" db:"action"`
	Detail             string    `json:"detail" db:"detail"`
	UserName           string    `json:"user_name" db:"user_name"`
	IP                 string    `json:"ip" db:"ip"`
	Browser            string    `json:"browser" db:"browser"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	CreatedAt          time.Time `json:"-" db:"created_at"`
	FormattedCreatedAt string    `json:"created_at"`
}

type CreateLogSystem struct {
	ID        string `db:"id"`
	Modul     string `json:"modul" db:"modul" binding:"required"`
	SubModul  string `json:"sub_modul" db:"sub_modul"`
	Action    string `json:"action" db:"action"`
	Detail    string `json:"detail" db:"detail"`
	UserName  string `db:"user_name"`
	IP        string `db:"ip"`
	Browser   string `db:"browser"`
	CreatedBy string `db:"created_by"`
	CreatedAt string `db:"created_at"`
}

type LogSystemExport struct {
	Modul  string `json:"modul"`
	Sub    string `json:"sub"`
	Action string `json:"action"`
	Detail string `json:"detail"`
	User   string `json:"user"`
	IP     string `json:"ip"`
	Date   string `json:"date"`
}
