package databasemodel

import (
	"time"
)

type DealerUtama struct {
	ID                      string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                    string `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	ParticipantID           string `gorm:"type:text" json:"participant_id"`
	Nama                    string `json:"nama" gorm:"column:name"`
	Tipe                    string `json:"tipe" gorm:"column:permission_type"`
	OperationalStatus       string `json:"status_operasional" gorm:"column:operational_status"`
	NameJson                []byte
	StructureManagementJson []byte
	AddressJson             []byte
	BillingAddressJson      []byte
	TaxPayerIdJson          []byte
	LogoJson                []byte
	RegistrationJson        []byte
	RevocationJson          []byte    //todo
	CreatedAt               time.Time `gorm:"<-:create"`
	CreatedBy               string    `gorm:"<-:create"`
	UpdatedAt               time.Time
	UpdatedBy               string
}

// TableName specifies the table name for the DealerUtama model.
func (DealerUtama) TableName() string {
	return "dealer_utama"
}
