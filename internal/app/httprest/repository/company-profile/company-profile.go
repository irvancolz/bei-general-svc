package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type GetProfilePjsppaResponse struct {
	Id                        string `json:"id"`
	Code                      string `json:"code"`
	Name                      string `json:"name"`
	Operational_Status        string `json:"operational_status"`
	Permission_Type           string `json:"permission_type"`
	Registration_Json         string `json:"registration_json"`
	Revocation_Json           string `json:"revocation_json"`
	Name_Json                 string `json:"name_json"`
	Code_Json                 string `json:"code_json"`
	Sppj_Sppa_Json            string `json:"sppj_sppa_json"`
	Permission_Type_Json      string `json:"permission_type_json"`
	Operational_Status_Json   string `json:"operator_status_json"`
	Address_Json              string `json:"address_json"`
	Structure_Management_Json string `json:"structure_management_json"`
	Ownership_Status_Json     string `json:"ownership_status_json"`
	Company_Status_Json       string `json:"company_status_json"`
	Billing_Address_Json      string `json:"billing_address_json"`
	Tax_Payer_Id_Json         string `json:"tax_payer_id_json"`
	Logo_Json                 string `json:"logo_json"`
}

type GetDuProfileResponse struct {
	Id                        string `json:"id"`
	Code                      string `json:"code"`
	Name                      string `json:"name"`
	Operational_Status        string `json:"operational_status"`
	Permission_Type           string `json:"permission_type"`
	Name_Json                 string `json:"name_json"`
	Structure_Management_Json string `json:"structure_management_json"`
	Address_Json              string `json:"address_json"`
	Billing_Address_Json      string `json:"billing_address_json"`
	Tax_Payer_Id_Json         string `json:"tax_payer_id_json"`
	Logo_Json                 string `json:"logo_json"`
	Registration_Json         string `json:"registration_json"`
	Revocation_Json           string `json:"revocation_json"`
}

type GetProfileParticipantResponse struct {
	Id                          string `json:"id"`
	Code                        string `json:"code"`
	Name                        string `json:"name"`
	Operational_Status          string `json:"operational_status"`
	Permission_Type             string `json:"permission_type"`
	Name_Json                   string `json:"name_json"`
	Code_Json                   string `json:"code_json"`
	Permission_Type_Json        string `json:"permission_type_json"`
	Operational_Status_Json     string `json:"operator_status_json"`
	Address_Json                string `json:"address_json"`
	Billing_Address_Json        string `json:"billing_address_json"`
	Structure_Management_Json   string `json:"structure_management_json"`
	Ownership_Status_Json       string `json:"ownership_status_json"`
	Company_Status_Json         string `json:"company_status_json"`
	Installed_Screen_Total_Json string `json:"installed_screen_total_json"`
	Logo_Json                   string `json:"logo_json"`
	Registration_Json           string `json:"registration_json"`
	Revocation_Json             string `json:"revocation_json"`
}

func GetCompanyProfile(extType, id string) (interface{}, error) {

	if extType == "" {
		return nil, errors.New("failed to get company profile : please specify the company external type")
	}

	if id == "" {
		return nil, errors.New("failed to get company profile : please specify the company id")
	}

	query := generateGetCompanyProfileQuery(extType)

	dbConn, errInitDb := helper.InitDBConn(extType)
	if errInitDb != nil {
		return nil, errInitDb
	}

	rowResult := dbConn.QueryRowx(query, id)
	result, errResult := getResultType(extType, rowResult)
	if errResult != nil {
		return nil, errResult
	}

	return result, nil
}

func generateGetCompanyProfileQuery(extType string) string {
	if strings.EqualFold(extType, "pjsppa") {
		return `	SELECT 
				id,
				code,
				name,
				coalesce(operational_status, '') as operational_status,
				coalesce(permission_type, '') as permission_type,
				coalesce(registration_json::text, '') as registration_json,
				coalesce(address_json::text, '') as address_json ,
				coalesce(revocation_json::text, '') as revocation_json,
				coalesce(name_json::text, '') as name_json,
				coalesce(code_json::text, '') as code_json,
				coalesce(sppj_sppa_json::text, '') as sppj_sppa_json,
				coalesce(permission_type_json::text, '') as permission_type_json,
				coalesce(operational_status_json::text, '') as operational_status_json,
				coalesce(address_json::text, '') as address_json,
				coalesce(structure_management_json::text, '') as structure_management_json,
				coalesce(ownership_status_json::text, '') as ownership_status_json,
				coalesce(company_status_json::text, '') as company_status_json,
				coalesce(billing_address_json::text, '') as billing_address_json,
				coalesce(tax_payer_id_json::text, '') as tax_payer_id_json,
				coalesce(logo_json::text, '') as logo_json
				FROM pjsppa
				WHERE id = $1 
			AND deleted_at IS NULL
			AND deleted_by IS NUll`
	}
	if strings.EqualFold(extType, "du") {
		return `		SELECT 
			id,
			code,
			name,
			coalesce(permission_type,'') as permission_type,
			coalesce(operational_status,'') as operational_status,
			coalesce(name_json::text,'') as name_json,
			coalesce(address_json::text,'') as address_json,
			coalesce(structure_management_json::text,'') as structure_management_json,
			coalesce(billing_address_json::text,'') as billing_address_json,
			coalesce(tax_payer_id_json::text,'') as tax_payer_id_json,
			coalesce(logo_json::text,'') as logo_json,
			coalesce(registration_json::text,'') as registration_json,
			coalesce(revocation_json::text,'') as revocation_json
		FROM dealer_utama
		WHERE id = $1
		AND deleted_at IS NULL
		AND deleted_by IS NUll`
	}
	if strings.EqualFold(extType, "participant") {
		return `	SELECT 
			id,
			code,
			name,
			coalesce(permission_type,'') as permission_type,
			coalesce(operational_status,'') as operational_status,
			coalesce(name_json::text,'') as name_json,
			coalesce(code_json::text,'') as code_json,
			coalesce(permission_type_json::text,'') as permission_type_json,
			coalesce(operational_status_json::text,'') as operational_status_json,
			coalesce(address_json::text,'') as address_json,
			coalesce(billing_address_json::text,'') as billing_address_json,
			coalesce(structure_management_json::text,'') as structure_management_json,
			coalesce(ownership_status_json::text,'') as ownership_status_json,
			coalesce(company_status_json::text,'') as company_status_json,
			coalesce(installed_screen_total_json::text,'') as installed_screen_total_json,
			coalesce(logo_json::text,'') as logo_json,
			coalesce(registration_json::text,'') as registration_json,
			coalesce(revocation_json::text,'') as revocation_json
		FROM participant
		WHERE id = $1 
		AND deleted_at IS NULL
		AND deleted_by IS NUll`
	}
	if strings.EqualFold(extType, "ab") {
		return `SELECT 
			COALESCE(id::text,'') as id,
			COALESCE(name::text,'') as name,
			COALESCE(code::text,'') as code,
			COALESCE(permit_bursa::text,'') as permit_bursa,
			COALESCE(company_status::text,'') as company_status,
			COALESCE(operational_status::text,'') as operational_status,
			COALESCE(registration_json::text,'') as registration_json,
			COALESCE(revocation_json::text,'') as revocation_json,
			COALESCE(name_json::text,'') as name_json,
			COALESCE(code_json::text,'') as code_json,
			COALESCE(spab_json::text,'') as spab_json,
			COALESCE(operational_status_json::text,'') as operational_status_json,
			COALESCE(ownership_and_company_status_json::text,'') as ownership_and_company_status_json,
			COALESCE(address_json::text,'') as address_json,
			COALESCE(structure_management_json::text,'') as structure_management_json,
			COALESCE(shareholder_json::text,'') as shareholder_json,
			COALESCE(capital_json::text,'') as capital_json,
			COALESCE(bussiness_permit_ojk_json::text,'') as bussiness_permit_ojk_json,
			COALESCE(permit_bursa_json::text,'') as permit_bursa_json,
			COALESCE(other_business_permit_ojk_json::text,'') as other_business_permit_ojk_json,
			COALESCE(amount_of_customers_and_employees_json::text,'') as amount_of_customers_and_employees_json,
			COALESCE(tax_payer_id_json::text,'') as tax_payer_id_json,
			COALESCE(incorporation_deed_json::text,'') as incorporation_deed_json,
			COALESCE(logo_json::text,'') as logo_json
		FROM anggota_bursa
		WHERE id = $1 AND deleted_by is null`
	}
	return ""
}

func getResultType(extType string, rowResult *sqlx.Row) (interface{}, error) {

	if strings.EqualFold(extType, "participant") {
		result := GetProfileParticipantResponse{}
		errScan := rowResult.StructScan(&result)
		if errScan != nil {
			return nil, errScan
		}
		return result, nil
	}
	if strings.EqualFold(extType, "du") {
		result := GetDuProfileResponse{}
		errScan := rowResult.StructScan(&result)
		if errScan != nil {
			return nil, errScan
		}
		return result, nil
	}
	if strings.EqualFold(extType, "pjsppa") {
		result := GetProfilePjsppaResponse{}
		errScan := rowResult.StructScan(&result)
		if errScan != nil {
			return nil, errScan
		}
		return result, nil
	}

	var result map[string]interface{}
	var scannedResult interface{}
	errScan := rowResult.MapScan(result)
	if errScan != nil {
		log.Println(extType)
		return nil, errScan
	}
	bytesResult, errMarshall := json.Marshal(&result)
	if errMarshall != nil {
		log.Println(errMarshall)
		return nil, errMarshall
	}
	errUnmarshall := json.Unmarshal(bytesResult, &scannedResult)
	if errUnmarshall != nil {
		log.Println(errUnmarshall)
		return nil, errUnmarshall
	}
	return scannedResult, nil
}
