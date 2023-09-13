package responsemodel

import (
	"time"
)

type Pjsppa struct {
	ID                string      `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode              string      `json:"kode" gorm:"column:code;unique;<-:create" copier:"-"`
	KodeAb            string      `json:"kode_ab" gorm:"column:ab_code"`
	Nama              string      `json:"nama" gorm:"column:name"`
	Tipe              string      `json:"tipe" gorm:"column:permission_type"`
	OperationalStatus string      `json:"status_operasional" gorm:"column:operational_status"`
	RegistrationJson  interface{} `json:"registration_json"`
	CreatedAt         time.Time   `gorm:"<-:create"`
	CreatedBy         string      `gorm:"<-:create"`
	UpdatedAt         time.Time
	UpdatedBy         string
}

func (Pjsppa) TableName() string {
	return "pjsppa"
}
