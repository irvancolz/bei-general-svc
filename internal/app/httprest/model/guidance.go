package model

import (
	"database/sql"
	"time"
)

// raw data obtained from database and need to be filtered in usecase
type GuidanceFileAndRegulationsJSONResponse struct {
	Id          string `json:"id"`
	Category    string `json:"category" binding:"required" validate:"oneof=Guidebook File Regulation"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Link        string `json:"link"`
	File        string `json:"file"`
	File_size   int64  `json:"file_size"`
	File_path   string `json:"file_path"`
	Version     string `json:"version"`
	File_Group  string `json:"file_group"`
	File_Owner  string `json:"file_owner"`
	Is_Visible  bool   `json:"is_visible"`
	Created_by  string `json:"created_by"`
	Created_at  int64  `json:"created_at"`
	Updated_by  string `json:"updated_by"`
	Updated_at  int64  `json:"updated_at"`
}

// struct to handling null result set from database
type GuidanceFileAndRegulationsResultSetResponse struct {
	Id          string         `json:"id"`
	Category    string         `json:"category" binding:"required" validate:"oneof=Guidebook File Regulation"`
	Name        string         `json:"name" binding:"required"`
	Description sql.NullString `json:"description"`
	Link        sql.NullString `json:"link"`
	File        string         `json:"file"`
	File_size   int64          `json:"file_size"`
	File_path   string         `json:"file_path"`
	File_Group  sql.NullString `json:"file_group"`
	File_Owner  sql.NullString `json:"file_owner"`
	Is_Visible  sql.NullBool   `json:"is_visible"`
	Version     sql.NullString `json:"version"`
	Created_by  string         `json:"created_by"`
	Created_at  time.Time      `json:"created_at"`
	Updated_by  sql.NullString `json:"updated_by"`
	Updated_at  sql.NullTime   `json:"updated_at"`
}

// actual result data structure given to user

type GuidanceJSONResponse struct {
	Id          string `json:"id"`
	Category    string `json:"category"`
	File_Group  string `json:"file_group"`
	Description string `json:"description"`
	Name        string `json:"name" binding:"required"`
	File_size   int64  `json:"file_size"`
	File_path   string `json:"file_path"`
	Version     string `json:"version"`
	Owner       string `json:"owner"`
	File        string `json:"file"`
	Link        string `json:"link"`
	Is_Visible  bool   `json:"is_visible"`
	Created_by  string `json:"created_by"`
	Created_at  int64  `json:"created_at"`
	Updated_by  string `json:"updated_by"`
	Updated_at  int64  `json:"updated_at"`
}
type RegulationJSONResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name" binding:"required"`
	Category   string `json:"category"`
	File       string `json:"file"`
	File_size  int64  `json:"file_size"`
	File_path  string `json:"file_path"`
	Version    string `json:"version"`
	Created_by string `json:"created_by"`
	Created_at int64  `json:"created_at"`
	Updated_by string `json:"updated_by"`
	Updated_at int64  `json:"updated_at"`
}

type GuidanceFilesJSONResponse struct {
	Id         string `json:"id"`
	Category   string `json:"category"`
	Name       string `json:"name" binding:"required"`
	File       string `json:"file"`
	File_size  int64  `json:"file_size"`
	File_path  string `json:"file_path"`
	Created_by string `json:"created_by"`
	Created_at int64  `json:"created_at"`
	Updated_by string `json:"updated_by"`
	Updated_at int64  `json:"updated_at"`
}
