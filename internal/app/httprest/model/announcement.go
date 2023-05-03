package model
type Announcement struct {
	ID              string `json:"id"`
	InformationType string `json:"information_type"`
	EffectiveDate   string `json:"effective_date"`
	Regarding       string `json:"regarding"`
}

type CreateAnnouncement struct {
	InformationType string `json:"information_type" binding:"required"`
	EffectiveDate   string `json:"effective_date" binding:"required"`
	Regarding       string `json:"regarding" binding:"required"`
}

type UpdateAnnouncement struct {
	ID              string `json:"id" binding:"required"`
	InformationType string `json:"information_type" binding:"required"`
	EffectiveDate   string `json:"effective_date" binding:"required"`
	Regarding       string `json:"regarding" binding:"required"`
}
