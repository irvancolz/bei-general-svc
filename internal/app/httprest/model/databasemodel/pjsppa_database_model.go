package databasemodel

import (
	"time"

	"gorm.io/datatypes"
)

type Pjsppa struct {
	ID                      string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                    string `json:"kode" gorm:"column:code;unique;<-:create" copier:"-"`
	KodeAb                  string `json:"kode_ab" gorm:"column:ab_code"`
	Nama                    string `json:"nama" gorm:"column:name"`
	Tipe                    string `json:"tipe" gorm:"column:permission_type"`
	OperationalStatus       string `json:"status_operasional" gorm:"column:operational_status"`
	RegistrationJson        datatypes.JSON
	RevocationJson          datatypes.JSON
	NameJson                datatypes.JSON
	CodeJson                datatypes.JSON
	SppjSppaJson            datatypes.JSON
	PermissionTypeJson      datatypes.JSON
	OperationalStatusJson   datatypes.JSON
	AddressJSON             datatypes.JSON
	StructureManagementJson datatypes.JSON
	OwnershipStatusJson     datatypes.JSON
	CompanyStatusJson       datatypes.JSON
	BillingAddressJson      datatypes.JSON
	TaxPayerIDJSON          datatypes.JSON
	LogoJSON                datatypes.JSON
	CreatedAt               time.Time `gorm:"<-:create"`
	CreatedBy               string    `gorm:"<-:create"`
	UpdatedAt               time.Time
	UpdatedBy               string
}

func (Pjsppa) TableName() string {
	return "pjsppa"
}
