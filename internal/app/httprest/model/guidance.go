package model

import "time"

type GuidanceFileAndRegulationsDBStructure struct {
	Id          string    `json:"id"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	File        string    `json:"file"`
	Version     float64   `json:"version"`
	Order       int64     `json:"order"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
	Updated_by  string    `json:"updated_by"`
	Updated_at  time.Time `json:"updated_at"`
	Deleted_by  string    `json:"deleted_by"`
	Deleted_at  time.Time `json:"deleted_at"`
}
