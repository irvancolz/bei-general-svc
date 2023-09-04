package pkp

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func GetUserFullname(db *sqlx.DB, user_id string) string {
	var result string
	query := `SELECT fullname FROM users WHERE id = $1`
	if err := db.QueryRowx(query, user_id).Scan(&result); err != nil {
		log.Println("failed get user fullname from db :", err)
	}
	return result
}
