package companyprofile

import (
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

	result, errorResult := companyprofile.GetCompanyProfileXml(c, companyProfileXml)
	if errorResult != nil {
		model.GenerateReadErrorResponseXml(c, errorResult)
		return
	}
	
	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))
}
