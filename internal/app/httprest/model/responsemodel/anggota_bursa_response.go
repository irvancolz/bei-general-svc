package responsemodel

import "time"

type AngggotaBursa struct {
	ID                string      `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode              string      `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	Nama              string      `json:"nama" gorm:"column:name"`
	PermitBursa       string      `json:"permit_bursa" gorm:"column:permit_bursa"`
	CompanyStatus     string      `json:"company_status" gorm:"column:company_status"`
	OperationalStatus string      `json:"status_operasional" gorm:"column:operational_status"`
	RegistrationJson  interface{} `json:"registration_json"`
	CreatedAt         time.Time   `gorm:"<-:create"`
	CreatedBy         string      `gorm:"<-:create"`
	UpdatedAt         time.Time
	UpdatedBy         string
}
