package databasemodel

import (
	"time"

	"gorm.io/datatypes"
)

type AngggotaBursa struct {
	ID                                string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                              string `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	Nama                              string `json:"nama" gorm:"column:name"`
	PermitBursa                       string `json:"permit_bursa" gorm:"column:permit_bursa"`
	CompanyStatus                     string `json:"company_status" gorm:"column:company_status"`
	OperationalStatus                 string `json:"status_operasional" gorm:"column:operational_status"`
	RegistrationJson                  datatypes.JSON
	NameJson                          datatypes.JSON
	CodeJson                          datatypes.JSON //todo dismiss
	SpabJson                          datatypes.JSON //todo
	OperationalStatusJson             datatypes.JSON //todo
	AddressJson                       datatypes.JSON
	StructureManagementJson           datatypes.JSON
	ShareholderJson                   datatypes.JSON
	BussinessPermitOjkJson            datatypes.JSON
	PermitBursaJson                   datatypes.JSON //todo
	OtherBusinessPermitOjkJson        datatypes.JSON  //todo is already inside businesspermitojk
	AmountOfCustomersAndEmployeesJson datatypes.JSON  //todo
	TaxPayerIdJson                    datatypes.JSON 
	IncorporatioincnDeedJson          datatypes.JSON 
	OwnershipAndCompanyStatusJson     datatypes.JSON
	LogoJson                          datatypes.JSON
	RevocationJson                    datatypes.JSON //todo
	CreatedAt                         time.Time `gorm:"<-:create"`
	CreatedBy                         string    `gorm:"<-:create"`
	UpdatedAt                         time.Time
	UpdatedBy                         string
}

type Tabler interface {
	TableName() string
}

func (AngggotaBursa) TableName() string {
	return "anggota_bursa"
}
