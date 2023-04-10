package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	MySql            *sqlx.DB
	AuthServiceMySql *sqlx.DB
}
