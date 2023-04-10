package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbConn = &DB{}
)

// disable this. development only!
// func init() {
// 	if os.Getenv("HTTP_PORT") == "" {
// 		log.Println("executed development code!!! please check inside database package")
// 		if err := godotenv.Load("./configs/.env"); err != nil {
// 			log.Fatalln("cannot load env files", err)
// 		}
// 	}
// }

func Init() *DB {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PSWD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

	return dbConn
}
