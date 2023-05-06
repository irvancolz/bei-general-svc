package model

import "time"

type PKuser struct {
	ID           string    `json:"id"`
	Stakeholders string    `json:"stakeholders"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	QuestionDate time.Time `json:"question_date"`
	Question     string    `json:"question"`
	Answers      string    `json:"answers"`
	AnswersBy    string    `json:"answers_by"`
	AnswersAt    time.Time `json:"answers_at`
	Topic        string    `json:"topic"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	CreateBy     string    `json:"create_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    string    `json:"update_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedBy    *int      `json:"deleted_by"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type CreatePKuser struct {
	Stakeholders string    `json:"stakeholders"`
	Code         string    `json:"code" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	QuestionDate time.Time `json:"question_date"" validate:"required"`
	Question     string    `json:"question"`
	Answers      string    `json:"answers"`
	AnswersBy    string    `json:"answers_by"`
	AnswersAt    string    `json:"answers_at`
	Topic        string    `json:"topic" validate:"required"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
}

type UpdatePKuser struct {
	ID           string    `json:"id"`
	Stakeholders string    `json:"stakeholders"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	QuestionDate time.Time `json:"question_date"`
	Question     string    `json:"question"`
	Answers      string    `json:"answers"`
	AnswersBy    string    `json:"answers_by"`
	AnswersAt    string    `json:"answers_at`
	Topic        string    `json:"topic"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	CreateBy     string    `json:"create_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    string    `json:"update_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EditRequest struct {
	ID           string `json:"id"`
	Stakeholders string `json:"stakeholders"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	AnswersAt    string `json:"answers_at`
	Topic        string `json:"topic"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
}
