package companyprofile

import (
	companyprofile "be-idx-tsg/internal/app/httprest/repository/company-profile"

	"github.com/gin-gonic/gin"
)

func GetCompanyProfile(c *gin.Context) (interface{}, error) {
	extType := c.Query("external_type")
	id := c.Query("id")

	return companyprofile.GetCompanyProfile(extType, id)
}

