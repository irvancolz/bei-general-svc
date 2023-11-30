package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbConn = &DB{}
)

func Init() *DB {
	var once sync.Once

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PSWD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	once.Do(func() {
		db, err := sqlx.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatalln(err)
		}

		db.SetMaxOpenConns(7)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(time.Duration(300 * time.Second))

		if err := db.Ping(); err != nil {
			log.Fatalln(err)
		}

		dbConn.MySql = db
	})

	return dbConn
}
