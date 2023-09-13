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
	result, errorResult := companyprofile.GetCompanyProfile(c)
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
