package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/model/requestmodel"
	companyprofile "be-idx-tsg/internal/app/httprest/usecase/company_profile"
	"be-idx-tsg/internal/pkg/httpresponse"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetCompanyProfile(c *gin.Context) {
	extType := c.Query("external_type")

	result, errorResult := companyprofile.GetCompanyProfile(c, extType)
	if errorResult != nil {
		model.GenerateReadErrorResponse(c, errorResult)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, result))
}

func GetCompanyProfileSingleLatest(c *gin.Context) {
	filterQueryParameter := helper.GetFilterQueryParameter(c)
	filterQueryParameter.Limit = 1

	responseData, maxPage, errorStr := companyprofile.GetCompanyProfileLatest(c, filterQueryParameter)
	if len(errorStr) > 0 {
		model.GenerateReadErrorResponse(c, errors.New(errorStr))
		return
	}

	model.GenerateQueryListResponse(c, responseData, maxPage)
}

func GetCompanyProfileXml(c *gin.Context) {
	var companyProfileXml requestmodel.CompanyProfileXml

	if err := c.ShouldBindXML(&companyProfileXml); err != nil {
		model.GenerateReadErrorResponseXml(c, err)
		return
	}

	if len(companyProfileXml.CompanyCode) > 0 && len(companyProfileXml.ExternalType) == 0 {
		model.GenerateReadErrorResponseXml(c, errors.New("Invalid request, please supply the external tyoe"))
		return
	}

	companyProfileList, errorResult := companyprofile.GetCompanyProfileXml(c, companyProfileXml)
	if errorResult != nil {
		model.GenerateReadErrorResponseXml(c, errorResult)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/xml")
	_, err := c.Writer.Write(companyProfileList)

	if err != nil {
		model.GenerateReadErrorResponseXml(c, err)
		return
	}

}

func GetCompanyProfileJSON(c *gin.Context) {
	var listAllCompany []map[string]any

	listCompanyAB, _ := companyprofile.GetCompanyProfile(c, "ab")

	for _, ab := range listCompanyAB {
		kontak, _ := ab["registration_json"].(map[string]interface{})["kontak"]

		company := map[string]any{
			"company_id":   ab["id"],
			"company_code": ab["code"],
			"name":         ab["name"],
			"phone":        kontak.(map[string]interface{})["telepon"],
			"fax":          kontak.(map[string]interface{})["faksimili"],
			"address":      kontak.(map[string]interface{})["alamat"],
			"status":       ab["operational_status"],
		}

		listAllCompany = append(listAllCompany, company)
	}

	listCompanyParticipant, _ := companyprofile.GetCompanyProfile(c, "participant")

	for _, p := range listCompanyParticipant {
		company := map[string]any{
			"company_id":   p["id"],
			"company_code": p["code"],
			"company_name": p["name"],
			"phone":        p["registration_json"].(map[string]interface{})["phone"],
			"fax":          p["registration_json"].(map[string]interface{})["fax"],
			"address":      p["registration_json"].(map[string]interface{})["address"],
			"status":       p["operational_status"],
		}

		listAllCompany = append(listAllCompany, company)
	}

	listCompanyPjsppa, _ := companyprofile.GetCompanyProfile(c, "pjsppa")

	for _, pj := range listCompanyPjsppa {
		kontak, _ := pj["registration_json"].(map[string]interface{})["kontak"]

		company := map[string]any{
			"company_id":   pj["id"],
			"company_code": pj["code"],
			"company_name": pj["name"],
			"phone":        kontak.(map[string]interface{})["telepon"],
			"fax":          kontak.(map[string]interface{})["faksimili"],
			"address":      kontak.(map[string]interface{})["alamat"],
			"status":       pj["operational_status"],
		}

		listAllCompany = append(listAllCompany, company)
	}

	listCompanyDu, _ := companyprofile.GetCompanyProfile(c, "du")

	for _, du := range listCompanyDu {
		kontak, _ := du["registration_json"].(map[string]interface{})["kontak"]

		company := map[string]any{
			"company_id":   du["id"],
			"company_code": du["code"],
			"company_name": du["name"],
			"phone":        kontak.(map[string]interface{})["telepon"],
			"fax":          kontak.(map[string]interface{})["faksimili"],
			"address":      kontak.(map[string]interface{})["alamat"],
			"status":       du["operational_status"],
		}

		listAllCompany = append(listAllCompany, company)
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, listAllCompany))
}
