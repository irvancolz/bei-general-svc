package databasemodel

import (
	"time"
)

type AngggotaBursa struct {
	ID                                string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                              string `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	Nama                              string `json:"nama" gorm:"column:name"`
	PermitBursa                       string `json:"permit_bursa" gorm:"column:permit_bursa"`
	CompanyStatus                     string `json:"company_status" gorm:"column:company_status"`
	OperationalStatus                 string `json:"status_operasional" gorm:"column:operational_status"`
	RegistrationJson                  []byte
	NameJson                          []byte
	CodeJson                          []byte //todo dismiss
	SpabJson                          []byte //todo
	OperationalStatusJson             []byte //todo
	AddressJson                       []byte
	StructureManagementJson           []byte
	ShareholderJson                   []byte
	BussinessPermitOjkJson            []byte
	PermitBursaJson                   []byte //todo
	OtherBusinessPermitOjkJson        []byte  //todo is already inside businesspermitojk
	AmountOfCustomersAndEmployeesJson []byte  //todo
	TaxPayerIdJson                    []byte 
	IncorporatioincnDeedJson          []byte 
	OwnershipAndCompanyStatusJson     []byte
	LogoJson                          []byte
	RevocationJson                    []byte //todo
	CreatedAt                         time.Time `gorm:"<-:create"`
	CreatedBy                         string    `gorm:"<-:create"`
	UpdatedAt                         time.Time
	UpdatedBy                         string
}

type Ab struct {
	ID                                string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                              string `json:"kode" gorm:"column:code;unique;<-:create" copier:"CompanyCode"`
	Nama                              string `json:"nama" gorm:"column:name"`
	PermitBursa                       string `json:"permit_bursa" gorm:"column:permit_bursa"`
	CompanyStatus                     string `json:"company_status" gorm:"column:company_status"`
	OperationalStatus                 string `json:"status_operasional" gorm:"column:operational_status"`
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
