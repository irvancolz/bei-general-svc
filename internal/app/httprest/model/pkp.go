package model

import (
	"database/sql"
	"time"
)

type PKuser struct {
	ID             string `json:"id"`
	Stakeholders   string `json:"stakeholders"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	QuestionDate   int64  `json:"question_date"`
	Question       string `json:"question"`
	Answers        string `json:"answers"`
	AnswersBy      string `json:"answers_by"`
	AnswersAt      int64  `json:"answers_at"`
	Topic          string `json:"topic"`
	ExternalType   string `json:"external_type"`
	AdditionalInfo string `json:"additional_info"`
	FileName       string `json:"file_name"`
	FilePath       string `json:"file_path"`
	CreateBy       string `json:"create_by"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedBy      string `json:"update_by"`
	UpdatedAt      int64  `json:"updated_at"`
	DeletedBy      string `json:"deleted_by"`
	DeletedAt      int64  `json:"deleted_at"`
}

type PKuserResultSet struct {
	ID             string         `json:"id"`
	Stakeholders   string         `json:"stakeholders"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	QuestionDate   time.Time      `json:"question_date"`
	Question       string         `json:"question"`
	Answers        string         `json:"answers"`
	AnswersBy      string         `json:"answers_by"`
	AnswersAt      time.Time      `json:"answers_at"`
	Topic          string         `json:"topic"`
	FileName       string         `json:"file_name"`
	FilePath       string         `json:"file_path"`
	CreateBy       string         `json:"create_by"`
	ExternalType   string         `json:"external_type"`
	AdditionalInfo sql.NullString `json:"additional_info" `
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedBy      sql.NullString `json:"update_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	DeletedBy      sql.NullString `json:"deleted_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
}

type CreatePKuser struct {
	Stakeholders   string    `json:"stakeholders" binding:"required"`
	Code           string    `json:"code" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	QuestionDate   time.Time `json:"question_date" binding:"required"`
	Question       string    `json:"question" binding:"required"`
	Answers        string    `json:"answers" binding:"required"`
	AnswersBy      string    `json:"answers_by" binding:"required"`
	AnswersAt      time.Time `json:"answers_at" binding:"required"`
	Topic          string    `json:"topic" binding:"required"`
	FileName       string    `json:"file_name" binding:"required"`
	FilePath       string    `json:"file_path" binding:"required"`
	ExternalType   string    `json:"external_type" binding:"required"`
	AdditionalInfo string    `json:"additional_info" binding:"required"`
}

type UpdatePKuser struct {
	ID             string    `json:"id" binding:"required"`
	ExternalType   string    `json:"external_type" binding:"required"`
	AdditionalInfo string    `json:"additional_info" binding:"required"`
	Stakeholders   string    `json:"stakeholders"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	QuestionDate   time.Time `json:"question_date"`
	Question       string    `json:"question"`
	Answers        string    `json:"answers"`
	AnswersBy      string    `json:"answers_by"`
	AnswersAt      time.Time `json:"answers_at"`
	Topic          string    `json:"topic"`
	FileName       string    `json:"file_name"`
	FilePath       string    `json:"file_path"`
	UpdatedBy      string    `json:"update_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}
