package model

import (
	"database/sql"
	"time"
)

type DivisionNameResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Is_Default bool   `json:"is_default"`
}

type InstitutionDivisionResponse struct {
	Id           string                       `json:"id"`
	Name         string                       `json:"name"`
	Company_name string                       `json:"company_name"`
	Company_code string                       `json:"company_code"`
	Company_id   string                       `json:"company_id"`
	Members      []InstitutionMembersResponse `json:"members"`
	Created_at   time.Time                    `json:"created_at"`
	Created_by   string                       `json:"created_by"`
	Updated_at   time.Time                    `json:"updated_at"`
	Updated_by   string                       `json:"updated_by"`
}

type InstitutionDivisionByCodeResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Company_name string `json:"company_name"`
	Company_code string `json:"company_code"`
	Company_id   string `json:"company_id"`
}

type InstitutionDivisionResultSetResponse struct {
	Id           string                       `json:"id"`
	Name         string                       `json:"name"`
	Company_name string                       `json:"company_name"`
	Members      []InstitutionMembersResponse `json:"members"`
	Created_at   sql.NullTime                 `json:"created_at"`
	Created_by   sql.NullString               `json:"created_by"`
	Updated_at   sql.NullTime                 `json:"updated_at"`
	Updated_by   sql.NullString               `json:"updated_by"`
}

type InstitutionProfileDetailResponse struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Code       string    `json:"code"`
	Status     string    `json:"status"`
	Created_at time.Time `json:"created_at"`
	Created_by string    `json:"created_by"`
	Updated_at time.Time `json:"updated_at"`
	Updated_by string    `json:"updated_by"`
}

type InstitutionProfileDetaiResultSetResponse struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Type       sql.NullString `json:"type"`
	Code       string         `json:"code"`
	Status     string         `json:"status"`
	Created_at time.Time      `json:"created_at"`
	Created_by sql.NullString `json:"created_by"`
	Updated_at sql.NullTime   `json:"updated_at"`
	Updated_by sql.NullString `json:"updated_by"`
}

type InstitutionMembersByCompanyResultSetResponse struct {
	Institute_id     string         `json:"institute_id"`
	Institute_status string         `json:"institute_status"`
	Institute_type   sql.NullString `json:"institute_type"`
	Company_code     string         `json:"company_code"`
	Company_name     string         `json:"company_name"`
	Id               string         `json:"id"`
	Name             string         `json:"name"`
	Email            string         `json:"email"`
	Phone            string         `json:"phone"`
	Telephone        string         `json:"telephone"`
	Division         sql.NullString `json:"division"`
	Position         string         `json:"position"`
}

type InstitutionMembersByCompanyResponse struct {
	Institute_id     string `json:"institute_id"`
	Institute_status string `json:"institute_status"`
	Institute_type   string `json:"institute_type"`
	Id               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Telephone        string `json:"telephone"`
	Division         string `json:"division"`
	Position         string `json:"position"`
	Company_code     string `json:"company_code"`
	Company_name     string `json:"company_name"`
}

type InstitutionMembersResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Telephone    string `json:"telephone"`
	Division     string `json:"division"`
	Position     string `json:"position"`
	Company_code string `json:"company_code"`
	Company_name string `json:"company_name"`
}

type InstitutionMembersDetailResponse struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Institution_id string    `json:"company_id"`
	Division_Id    string    `json:"division_id"`
	Division       string    `json:"division"`
	Position       string    `json:"position"`
	Phone          string    `json:"phone"`
	Telephone      string    `json:"telephone"`
	Email          string    `json:"email"`
	Created_at     time.Time `json:"created_at"`
	Created_by     string    `json:"created_by"`
	Updated_at     time.Time `json:"updated_at"`
	Updated_by     string    `json:"updated_by"`
}

type InstitutionMembersDetailResultSetResponse struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	Institution_id string         `json:"company_id"`
	Division_Id    string         `json:"division_id"`
	Division       string         `json:"division"`
	Position       string         `json:"position"`
	Phone          string         `json:"phone"`
	Telephone      string         `json:"telephone"`
	Email          string         `json:"email"`
	Created_at     sql.NullTime   `json:"created_at"`
	Created_by     sql.NullString `json:"created_by"`
	Updated_at     sql.NullTime   `json:"updated_at"`
	Updated_by     sql.NullString `json:"updated_by"`
}

type InstitutionResponse struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Is_deleted bool   `json:"is_deleted"`
}

type InstitutionResultSetResponse struct {
	Id     string         `json:"id"`
	Type   sql.NullString `json:"type"`
	Code   string         `json:"code"`
	Name   string         `json:"name"`
	Status string         `json:"status"`
}

type ContactPersonSyncCompaniesResource struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Is_deleted bool   `json:"is_deleted"`
}
