package model

import (
	"database/sql"
	"time"
)

// default struct to add data to db
type GuidanceFileAndRegulationsDBStructure struct {
	Id          string    `json:"id"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	File        string    `json:"file"`
	File_size   int64     `json:"file_size"`
	Version     int64     `json:"version"`
	Order       bool      `json:"order"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
	Updated_by  string    `json:"updated_by"`
	Updated_at  time.Time `json:"updated_at"`
	Deleted_by  string    `json:"deleted_by"`
	Deleted_at  time.Time `json:"deleted_at"`
}

// raw data obtained from database and need to be filtered in usecase
type GuidanceFileAndRegulationsJSONResponse struct {
	Id          string    `json:"id"`
	Category    string    `json:"category" binding:"required" validate:"oneof=Guidebook File Regulation"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	File        string    `json:"file"`
	File_size   int64     `json:"file_size"`
	Version     float64   `json:"version"`
	Order       int64     `json:"order"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
	Updated_by  string    `json:"updated_by"`
	Updated_at  time.Time `json:"updated_at"`
}

// struct to handling null result set from database
type GuidanceFileAndRegulationsResultSetResponse struct {
	Id          string         `json:"id"`
	Category    string         `json:"category" binding:"required" validate:"oneof=Guidebook File Regulation"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Link        sql.NullString `json:"link"`
	File        sql.NullString `json:"file"`
	File_size   sql.NullInt64  `json:"file_size"`
	Version     float64        `json:"version"`
	Order       int64          `json:"order"`
	Created_by  string         `json:"created_by"`
	Created_at  time.Time      `json:"created_at"`
	Updated_by  sql.NullString `json:"updated_by"`
	Updated_at  sql.NullTime   `json:"updated_at"`
}

// actual result data structure given to user

type GuidanceJSONResponse struct {
	Id          string    `json:"id"`
	Category    string    `json:"category" binding:"required" validate:"oneof=Guidebook"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	File        string    `json:"file"`
	File_size   int64     `json:"file_size"`
	Version     float64   `json:"version"`
	Order       int64     `json:"order"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
	Updated_by  string    `json:"updated_by"`
	Updated_at  time.Time `json:"updated_at"`
}
type RegulationJSONResponse struct {
	Id         string    `json:"id"`
	Category   string    `json:"category" binding:"required" validate:"oneof=File"`
	Name       string    `json:"name" binding:"required"`
	Link       string    `json:"link"`
	Order      int64     `json:"order"`
	Created_by string    `json:"created_by"`
	Created_at time.Time `json:"created_at"`
	Updated_by string    `json:"updated_by"`
	Updated_at time.Time `json:"updated_at"`
}

type GuidanceFilesJSONResponse struct {
	Id          string    `json:"id"`
	Category    string    `json:"category" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Version     float64   `json:"version"`
	Order       int64     `json:"order"`
	File        string    `json:"file"`
	File_size   int64     `json:"file_size"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
	Updated_by  string    `json:"updated_by"`
	Updated_at  time.Time `json:"updated_at"`
}
