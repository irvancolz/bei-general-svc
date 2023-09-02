package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"be-idx-tsg/internal/app/httprest/model/requestmodel"
	"errors"
)


func GetCompanyProfileAb(request requestmodel.CompanyProfileXml) ([]databasemodel.AngggotaBursa, error)  {
	listData := []databasemodel.AngggotaBursa{}
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	dbConn.Find(&listData)
	
	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company ab: "+dbConn.Error.Error())
	}
	
	return listData, nil
}


func GetCompanyProfileParticipant(request requestmodel.CompanyProfileXml) ([]databasemodel.Participant, error)  {
	listData := []databasemodel.Participant{}
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	dbConn.Find(&listData)

	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company participant: "+dbConn.Error.Error())
	}

	return listData, nil
}

func GetCompanyProfilePJSPPA(request requestmodel.CompanyProfileXml) ([]databasemodel.Pjsppa, error)  {
	listData := []databasemodel.Pjsppa{}
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	dbConn.Find(&listData)

	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company pjsppa: "+dbConn.Error.Error())
	}

	return listData, nil
}

func GetCompanyDealerUtama(request requestmodel.CompanyProfileXml) ([]databasemodel.DealerUtama, error)  {
	listData := []databasemodel.DealerUtama{}
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	dbConn.Find(&listData)

	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company dealer utama: "+dbConn.Error.Error())
	}

	return listData, nil
}