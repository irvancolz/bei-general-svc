package model

type GetAllAnnouncement struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Type string `json:"type"`
}

type Announcement struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Type              string `json:"type"`
	RoleId            string `json:"role_id"`
	Role              string `json:"role"`
	EffectiveDate     string `json:"effective_date"`
	Regarding         string `json:"regarding"`
	Status            bool   `json:"status"`
}

type CreateAnnouncement struct {
	Code          string `json:"code"`
	Type          string `json:"type" binding:"required"`
	RoleId        string `json:"role_id" binding:"required"`
	EffectiveDate string `json:"effective_date" binding:"required"`
	Regarding     string `json:"regarding"`
	Status        bool   `json:"status"`
	FileURL       string `json:"file_url"`
}

type UpdateAnnouncement struct {
	Id            string `json:"id" binding:"required"`
	Code          string `json:"code"`
	Type          string `json:"type" binding:"required"`
	RoleId        string `json:"role_id" binding:"required"`
	EffectiveDate string `json:"effective_date" binding:"required"`
	Regarding     string `json:"regarding"`
	Status        bool   `json:"status" binding:"required"`
	FileURL       string `json:"file_url" binding:"required"`
	UpdatedBy     string `json:"updated_by" binding:"required"`
}

type AnnouncementCode struct {
	Code string `json:"code"`
	Type string `json:"type" binding:"required"`
}
