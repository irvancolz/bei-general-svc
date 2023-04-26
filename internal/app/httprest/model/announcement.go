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
	EffectiveDate     string `json:"effective_date"`
	Regarding         string `json:"regarding"`
	Status            bool   `json:"status"`
	License           string `json:"license"`
	OperationalStatus string `json:"operational_status"`
}

type CreateAnnouncement struct {
	Code          string `json:"code"`
	Type          string `json:"type" binding:"required"`
	RoleId        string `json:"role_id" binding:"required"`
	EffectiveDate string `json:"effective_date" binding:"required"`
	Regarding     string `json:"regarding"`
	Status        bool   `json:"status" binding:"required"`
	FileURL       string `json:"file_url" binding:"required"`
	CreatedBy     string `json:"created_by" binding:"required"`
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
