package model

type Announcement struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	RoleId        string `json:"role_id"`
	EffectiveDate string `json:"effective_date"`
	Regarding     string `json:"regarding"`
	Status        bool   `json:"status"`
}
