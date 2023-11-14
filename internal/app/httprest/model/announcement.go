package model

import "time"

type Announcement struct {
	ID               string `json:"id"`
	Information_Type string `json:"information_type"`
	Effective_Date   int64  `json:"effective_date"`
	Regarding        string `json:"regarding"`
	// todo rm in the future
	Type          string `json:"type"`
	Form_Value_Id string `json:"form_value_id"`
	Creator       string `json:"creator"`
	Creator_Name  string `json:"creator_name"`
}

type CreateAnnouncement struct {
	Information_Type string    `json:"information_type" binding:"required,oneof='INTERNAL BURSA' 'AB' 'PARTICIPANT' 'PJSPPA' 'SEMUA' 'DU'"`
	Effective_Date   time.Time `json:"effective_date" binding:"required"`
	Regarding        string    `json:"regarding" binding:"required"`
	// todo rm in the future
	Type          string `json:"type"`
	Form_Value_Id string `json:"form_value_id"`
}

type UpdateAnnouncement struct {
	ID               string    `json:"id" binding:"required"`
	Information_Type string    `json:"information_type" binding:"required,oneof='INTERNAL BURSA' 'AB' 'PARTICIPANT' 'PJSPPA' 'SEMUA' 'DU'"`
	Effective_Date   time.Time `json:"effective_date" binding:"required"`
	Regarding        string    `json:"regarding" binding:"required"`
	// todo rm in the future
	Type          string `json:"type"`
	Form_Value_Id string `json:"form_value_id"`
}
