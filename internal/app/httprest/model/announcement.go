package model

type Announcement struct {
	ID              string `json:"id"`
	InformationType string `json:"information_type"`
	EffectiveDate   string `json:"effective_date"`
	Regarding       string `json:"regarding"`
	Type            *string `json:"type"`
}

type CreateAnnouncement struct {
	InformationType string `json:"information_type" binding:"required,oneof='INTERNAL BURSA' 'AB' 'PARTICIPANT' 'PJSPPA' 'SEMUA' "`
	EffectiveDate   string `json:"effective_date" binding:"required"`
	Regarding       string `json:"regarding" binding:"required"`
	Type            *string `json:"type" binding:"required"`
}

type UpdateAnnouncement struct {
	ID              string `json:"id" binding:"required"`
	InformationType string `json:"information_type" binding:"required,oneof='Internal Bursa' 'Anggota Bursa' 'Participant' 'PJSPPA' 'Semua' "`
	EffectiveDate   string `json:"effective_date" binding:"required"`
	Regarding       string `json:"regarding" binding:"required"`
	Type           *string `json:"type" binding:"required"`
}
