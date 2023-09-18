package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	companyprofile "be-idx-tsg/internal/app/httprest/repository/company-profile"
	"encoding/json"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCompanyProfile(c *gin.Context) ([]map[string]interface{}, error) {
	extType := c.Query("external_type")

	companyList, errGetList := companyprofile.GetCompanyProfile(extType)
	if errGetList != nil {
		return nil, errGetList
	}

	filteredData, _ := helper.HandleDataFiltering(c, companyList, []string{})

	result := []map[string]interface{}{}
	for _, item := range filteredData {
		for key := range item {
			ogString := item[key]

			if json.Valid([]byte(ogString.(string))) {
				var formattedProps interface{}
				errUnmarshall := json.Unmarshal([]byte(ogString.(string)), &formattedProps)
				if errUnmarshall != nil {
					log.Println("failed to convert data to json :", errUnmarshall)
				}
				item[key] = formattedProps
			} else {
				item[key] = ogString
			}

		}

		result = append(result, item)
	}

	return result, nil
}

func GetCompanyProfileLatest(c *gin.Context, filterQueryParameter model.FilterQueryParameter) (interface{}, int, string) {

	authUserDetail, err := helper.GetAuthUserDetail(c)

	if err != nil {
		return nil, 0, "Failed to get auth user detail: " + err.Error()
	}

	switch strings.ToLower(authUserDetail.ExternalType) {
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
