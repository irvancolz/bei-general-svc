package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func InitDBConn(externalType string) (*sqlx.DB, error) {
	var (
		dbUrl  string
		host   string
		port   string
		user   string
		pass   string
		dbname string
	)

	if strings.EqualFold(externalType, "ab") {
		host = os.Getenv("DB_AB_HOST")
		port = os.Getenv("DB_AB_PORT")
		user = os.Getenv("DB_AB_USER")
		pass = os.Getenv("DB_AB_PSWD")
		dbname = os.Getenv("DB_AB_NAME")
	} else if strings.EqualFold(externalType, "du") {
		host = os.Getenv("DB_DU_HOST")
		port = os.Getenv("DB_DU_PORT")
		user = os.Getenv("DB_DU_USER")
		pass = os.Getenv("DB_DU_PSWD")
		dbname = os.Getenv("DB_DU_NAME")
	} else if strings.EqualFold(externalType, "pjsppa") {
		host = os.Getenv("DB_PJSPPA_HOST")
		port = os.Getenv("DB_PJSPPA_PORT")
		user = os.Getenv("DB_PJSPPA_USER")
		pass = os.Getenv("DB_PJSPPA_PSWD")
		dbname = os.Getenv("DB_PJSPPA_NAME")
	} else if strings.EqualFold(externalType, "participant") {
		host = os.Getenv("DB_PARTICIPANT_HOST")
		port = os.Getenv("DB_PARTICIPANT_PORT")
		user = os.Getenv("DB_PARTICIPANT_USER")
		pass = os.Getenv("DB_PARTICIPANT_PSWD")
		dbname = os.Getenv("DB_PARTICIPANT_NAME")
	} else if strings.EqualFold(externalType, "auth") {
		host = os.Getenv("DB_AUTH_HOST")
		port = os.Getenv("DB_AUTH_PORT")
		user = os.Getenv("DB_AUTH_USER")
		pass = os.Getenv("DB_AUTH_PSWD")
		dbname = os.Getenv("DB_AUTH_NAME")
	}
	dbUrl = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	return sqlx.Connect("postgres", dbUrl)
}
