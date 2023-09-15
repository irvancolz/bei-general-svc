package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"be-idx-tsg/internal/app/httprest/model/responsemodel"
	"encoding/json"
	"log"
)

func GetCompanyProfileParticipantLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*databasemodel.Participant, int, string) {

	latestProfileList := []*databasemodel.Participant{}
	dbConn, errInitDb := helper.InitDBConnGorm(authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company participant latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	}

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company participant latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)
	dbConn = dbConn.Find(&latestProfileList)
	count := helper.GetMaxPage(dbConn, databasemodel.Participant{}, filterqueryparameter.Limit)

	if dbConn.Error != nil {
		return latestProfileList, 0, "Failed to get company participant:Rows() " + dbConn.Error.Error()
	}

	return latestProfileList, count, errorStr
}

//todo
func GetCompanyProfileAbLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*responsemodel.AngggotaBursa, int, string) {

	latestProfileList := []*responsemodel.AngggotaBursa{}
	dbConn, errInitDb := helper.InitDBConnGorm(authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company ab latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	}

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company ab latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	dbConn = dbConn.Model(databasemodel.AngggotaBursa{})
	count := helper.GetMaxPage(dbConn, databasemodel.AngggotaBursa{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)

	rows, err := dbConn.Rows()

	if err != nil {
		return latestProfileList, 0, "Failed to get company ab:Rows() " + err.Error()
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.AngggotaBursa
		err := dbConn.ScanRows(rows, &item)

		if err != nil {
			return latestProfileList, 0, "Failed to get company ab - scanrows: " + err.Error()
		}

		responseModel := responsemodel.AngggotaBursa{}
		helper.Copy(&responseModel, item)
		if item.RegistrationJson != nil {
			var registrationJson interface{}

			log.Println(string(item.RegistrationJson))
			err := json.Unmarshal(item.RegistrationJson, &registrationJson)

			if err != nil {
				return latestProfileList, 0, "Failed to get company ab - unmarshal: " + err.Error()
			}

			responseModel.RegistrationJson = registrationJson

		}

		latestProfileList = append(latestProfileList, &responseModel)
	}

	return latestProfileList, count, errorStr
}

func GetCompanyProfilePjsppaLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*responsemodel.Pjsppa, int, string) {

	latestProfileList := []*responsemodel.Pjsppa{}
	dbConn, errInitDb := helper.InitDBConnGorm(authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	}

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	dbConn = dbConn.Model(databasemodel.Pjsppa{})
	count := helper.GetMaxPage(dbConn, databasemodel.Pjsppa{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)

	rows, err := dbConn.Rows()

	if err != nil {
		return latestProfileList, 0, "Failed to get company pjsppa:Rows() " + err.Error()
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.Pjsppa
		err := dbConn.ScanRows(rows, &item)

		if err != nil {
			return latestProfileList, 0, "Failed to get company pjsppa - scanrows: " + err.Error()
		}

		responseModel := responsemodel.Pjsppa{}
		helper.Copy(&responseModel, item)
		if item.RegistrationJson != nil {
			var registrationJson interface{}

			log.Println(string(item.RegistrationJson))
			err := json.Unmarshal(item.RegistrationJson, &registrationJson)

			if err != nil {
				return latestProfileList, 0, "Failed to get company pjsppa - unmarshal: " + err.Error()
			}

			responseModel.RegistrationJson = registrationJson

		}

		latestProfileList = append(latestProfileList, &responseModel)
	}

	return latestProfileList, count, errorStr
}

func GetCompanyProfileDuLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*responsemodel.DealerUtama, int, string) {

	latestProfileList := []*responsemodel.DealerUtama{}
	dbConn, errInitDb := helper.InitDBConnGorm(authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company du latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	}

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company du latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	dbConn = dbConn.Model(databasemodel.DealerUtama{})
	count := helper.GetMaxPage(dbConn, databasemodel.DealerUtama{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)

	rows, err := dbConn.Rows()

	if err != nil {
		return latestProfileList, 0, "Failed to get company du:Rows() " + err.Error()
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.DealerUtama
		err := dbConn.ScanRows(rows, &item)

		if err != nil {
			return latestProfileList, 0, "Failed to get company du - scanrows: " + err.Error()
		}

		responseModel := responsemodel.DealerUtama{}
		helper.Copy(&responseModel, item)
		if item.RegistrationJson != nil {
			var registrationJson interface{}

			log.Println(string(item.RegistrationJson))
			err := json.Unmarshal(item.RegistrationJson, &registrationJson)

			if err != nil {
				return latestProfileList, 0, "Failed to get company du - unmarshal: " + err.Error()
			}

			responseModel.RegistrationJson = registrationJson

		}

		latestProfileList = append(latestProfileList, &responseModel)
	}

	return latestProfileList, count, errorStr
}
