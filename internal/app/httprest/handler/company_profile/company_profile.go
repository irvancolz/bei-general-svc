package companyprofile

import (
	"be-idx-tsg/internal/app/httprest/model"
	companyprofile "be-idx-tsg/internal/app/httprest/usecase/company_profile"
	"be-idx-tsg/internal/pkg/httpresponse"

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
