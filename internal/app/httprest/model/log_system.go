package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

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
	Modul     NullString `db:"modul"`
	Sub_Modul NullString `db:"sub_modul"`
	Detail    NullString `db:"detail"`
	Action    NullString `db:"action"`
	User      NullString `db:"user_name"`
	IP        NullString `db:"ip"`
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

type NullString string

// Scan scans the value and assigns it to the NullString receiver.
//
// The value parameter is the value to be scanned. If the value is nil, the NullString receiver is set to an empty string.
// If the value is not a string, an error is returned with the message "Column is not a string".
// The function returns an error if the value parameter is not a string or nil.
func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("Column is not a string")
	}
	*s = NullString(strVal)
	return nil
}

// Value returns the value of the NullString as a driver.Value.
//
// It checks if the length of the NullString is 0 and returns nil, nil if it is.
// Otherwise, it returns the NullString as a string and nil for the error.
func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}
