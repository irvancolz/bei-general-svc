package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
)

func GetCompanyProfileParticipantLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*databasemodel.Participant, int, string)  {
	
	latestProfileList := []*databasemodel.Participant{}
	dbConn, errInitDb := helper.InitDBConnGorm(*authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company participant latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(*authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	} 

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company participant latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	count := helper.GetMaxPage(dbConn, databasemodel.Participant{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)
	dbConn = dbConn.Find(&latestProfileList)

	if dbConn.Error != nil {
		dbConn.Order(filterqueryparameter.OrderBy + " " + filterqueryparameter.Order).Offset(filterqueryparameter.Offset).Find(&latestProfileList)
	}

	return latestProfileList, count, errorStr
}


func GetCompanyProfileAbLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter)([]*databasemodel.AngggotaBursa, int, string)  {
	
	latestProfileList := []*databasemodel.AngggotaBursa{}
	dbConn, errInitDb := helper.InitDBConnGorm(*authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company ab latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(*authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	} 

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company ab latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	count := helper.GetMaxPage(dbConn, databasemodel.AngggotaBursa{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)
	dbConn = dbConn.Find(&latestProfileList)

	if dbConn.Error != nil {
		dbConn.Order(filterqueryparameter.OrderBy).Offset(filterqueryparameter.Offset).Find(&latestProfileList)
	}

	return latestProfileList, count, errorStr
}

func GetCompanyProfilePjsppaLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*databasemodel.Pjsppa, int, string)  {
	
	latestProfileList := []*databasemodel.Pjsppa{}
	dbConn, errInitDb := helper.InitDBConnGorm(*authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(*authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	} 

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	count := helper.GetMaxPage(dbConn, databasemodel.Pjsppa{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)
	dbConn = dbConn.Find(&latestProfileList)

	if dbConn.Error != nil {
		dbConn.Order(filterqueryparameter.OrderBy).Offset(filterqueryparameter.Offset).Find(&latestProfileList)
	}

	return latestProfileList, count, errorStr
}


func GetCompanyProfileDuLatest(authUserDetail model.AuthUserDetail, filterqueryparameter model.FilterQueryParameter) ([]*databasemodel.DealerUtama, int, string)  {
	
	latestProfileList := []*databasemodel.DealerUtama{}
	dbConn, errInitDb := helper.InitDBConnGorm(*authUserDetail.ExternalType)
	if errInitDb != nil {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.InitDBConnGorm(): " + errInitDb.Error()
	}

	if len(*authUserDetail.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", authUserDetail.CompanyCode)
	} 

	dbConn, errorStr := helper.GetGormQueryFilter(dbConn, filterqueryparameter.QueryList, filterqueryparameter.EndDate, filterqueryparameter.Now)

	if len(errorStr) > 0 {
		return latestProfileList, 0, "Failed to get company pjsppa latest -dbConn.GetGormQueryFilter(): " + errorStr
	}

	count := helper.GetMaxPage(dbConn, databasemodel.DealerUtama{}, filterqueryparameter.Limit)
	dbConn = dbConn.Order(helper.DEFAULT_ORDER_BY).Limit(filterqueryparameter.Limit).Offset(filterqueryparameter.Offset)
	dbConn = dbConn.Find(&latestProfileList)

	if dbConn.Error != nil {
		dbConn.Order(filterqueryparameter.OrderBy).Offset(filterqueryparameter.Offset).Find(&latestProfileList)
	}

	return latestProfileList, count, errorStr
}
