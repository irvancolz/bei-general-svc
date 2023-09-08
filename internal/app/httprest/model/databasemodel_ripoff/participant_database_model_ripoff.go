package databasemodelripoff

import (
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"time"
)

type Participant struct {
	ID                  string                       `gorm:"type:uuid;default:gen_random_uuid();primaryKey;<-:false" copier:"-"`
	Kode                string                       `json:"kode" gorm:"column:code;unique;<-:create" copier:"-"`
	KodeAb              string                       `json:"kode_ab" gorm:"column:ab_code"`
	Nama                string                       `json:"nama" gorm:"column:name"`
	Tipe                string                       `json:"tipe" gorm:"column:permission_type"`
	KodeUnik            string                       `json:"kode_unik" gorm:"-:all"`
	OperationalStatus   string                       `json:"status_operasional" gorm:"column:operational_status"`
	Alamat              string                       `json:"alamat" gorm:"-:all"`
	Telepon             *string                      `json:"telepon" gorm:"-:all"`
	Faks                string                       `json:"faks" gorm:"-:all"`
	Contact             string                       `json:"contact" gorm:"-:all"`
	MainNetwork         string                       `json:"main_network" gorm:"-:all"`
	User                []databasemodel.UserPLTEData `json:"user" gorm:"-:all"`
	TanggalEfektif      time.Time                    `json:"tanggal_efektif" gorm:"-:all"`
	NomorSuratPendukung string                       `json:"nomor_surat_pendukung" gorm:"-:all"`
	Alasan              string                       `json:"alasan" gorm:"-:all"`
	CreatedAt           time.Time                    `gorm:"<-:create"`
	CreatedBy           string                       `gorm:"<-:create"`
	UpdatedAt           time.Time
	UpdatedBy           string
}
