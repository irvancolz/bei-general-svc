package model

import "time"

type LogSystem struct {
	ID           string    `json:"id" db:"id"`
	Modul        string    `json:"modul" db:"modul" binding:"required"`
	Sub_Modul    string    `json:"sub_modul" db:"sub_modul"`
	Action       string    `json:"action" db:"action"`
	Detail       string    `json:"detail" db:"detail"`
	UserName     string    `json:"user_name" db:"user_name"`
	IP           string    `json:"ip" db:"ip"`
	Browser      string    `json:"browser" db:"browser"`
	Created_By   string    `json:"created_by" db:"created_by"`
	Created_At   string    `json:"created_at"`
	T_Created_At time.Time `json:"t_created_at" db:"t_created_at"`
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

type LogSystemFilter struct {
	Modul    string `db:"modul"`
	SubModul string `db:"sub_modul"`
	Detail   string `db:"detail"`
	Action   string `db:"action"`
	User     string `db:"user_name"`
	IP       string `db:"ip"`
}

func IsAllowedAction(action string) bool {
	blackList := []string{
		"List",
		"Filter",
		"Detail",
		"Sort",
	}

	for _, v := range blackList {
		if v == action {
			return false
		}
	}

	return true
}
