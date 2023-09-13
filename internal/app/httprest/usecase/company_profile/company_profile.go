package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	companyprofile "be-idx-tsg/internal/app/httprest/repository/company-profile"

	"github.com/gin-gonic/gin"
)

func GetCompanyProfile(c *gin.Context) (interface{}, error) {
	extType := c.Query("external_type")
	id := c.Query("id")

	return companyprofile.GetCompanyProfile(extType, id)
}

func GetCompanyProfileLatest(c *gin.Context, filterQueryParameter model.FilterQueryParameter) (interface{}, int, string) {
	
	authUserDetail, err := helper.GetAuthUserDetail(c)

	if err != nil {
		return nil, 0, "Failed to get auth user detail: " + err.Error()
	}

	switch authUserDetail.ExternalType {
	case REQUEST_EXTERNAL_TYPE_PARTICIPANT:
		{
			return companyprofile.GetCompanyProfileParticipantLatest(authUserDetail, filterQueryParameter)
		}
	case REQUEST_EXTERNAL_TYPE_AB:
		{
			return companyprofile.GetCompanyProfileAbLatest(authUserDetail, filterQueryParameter)
		}
	case REQUEST_EXTERNAL_TYPE_PJSPPA:
		{
			return companyprofile.GetCompanyProfilePjsppaLatest(authUserDetail, filterQueryParameter)
		}
	case REQUEST_EXTERNAL_TYPE_DU:
		{
			return companyprofile.GetCompanyProfileDuLatest(authUserDetail, filterQueryParameter)
		}
	default:
		{
			return companyprofile.GetCompanyProfileAbLatest(authUserDetail, filterQueryParameter)
		}
	}

}
