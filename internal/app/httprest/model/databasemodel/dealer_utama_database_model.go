package databasemodel

import (
	"time"

	"gorm.io/datatypes"
)

type DealerUtama struct {
	ID                      string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                    string `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	ParticipantID           string `gorm:"type:text" json:"participant_id"`
	Nama                    string `json:"nama" gorm:"column:name"`
	Tipe                    string `json:"tipe" gorm:"column:permission_type"`
	OperationalStatus       string `json:"status_operasional" gorm:"column:operational_status"`
	NameJson                datatypes.JSON
	StructureManagementJson datatypes.JSON
	AddressJson             datatypes.JSON
	BillingAddressJson      datatypes.JSON
	TaxPayerIdJson          datatypes.JSON
	LogoJson                datatypes.JSON
	RegistrationJson        datatypes.JSON
	RevocationJson          datatypes.JSON    //todo
	CreatedAt               time.Time `gorm:"<-:create"`
	CreatedBy               string    `gorm:"<-:create"`
	UpdatedAt               time.Time
	UpdatedBy               string
}

// TableName specifies the table name for the DealerUtama model.
func (DealerUtama) TableName() string {
	return "dealer_utama"
}
