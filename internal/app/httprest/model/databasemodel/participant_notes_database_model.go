package databasemodel

import "time"

type Notes struct {
	ID                string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" json:"id"`
	ParticipantID     string    `gorm:"type:text" json:"participant_id"`
	ReportDescription string    `json:"report_description"`
	ReferenceNo       string    `json:"reference_no"`
	UploadDate        string    `json:"upload_date"`
	ParticipantCode   string    `json:"participant_code"`
	ParticipantName   string    `json:"participant_name"`
	EventDate         string    `json:"event_date"`
	Category          string    `json:"category"`
	Action            string    `json:"action"`
	BursaUser         string    `json:"bursa_user"`
	Description       string    `json:"information"`
	CreatedAt         time.Time `gorm:"<-:false" json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedBy         string    `gorm:"<-:false" json:"updated_by"`
	UpdatedAt         time.Time `gorm:"<-:false" json:"updated_at"`
	DeletedBy         string    `gorm:"<-:false" json:"deleted_by"`
	DeletedAt         time.Time `gorm:"<-:false" json:"deleted_at"`
}

func (Notes) TableName() string {
	return "notes"
}
