package companyprofile

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"be-idx-tsg/internal/app/httprest/model/requestmodel"
	"log"

	"errors"

	"github.com/beevik/etree"
)

const (
	EXTERNAL_TYPE_LIST_AB = "AbList"
	EXTERNAL_TYPE_LIST_PARTICIPANT = "ParticipantList"
	EXTERNAL_TYPE_LIST_PJSPPA = "PjsppaList"
	EXTERNAL_TYPE_LIST_DU = "DuList"
)

func generateProfileXml(id string, registrationJSON []byte, name string) ([]byte, error) {
	// Parse the JSON data into a Registration struct

	file, err := helper.JsonToXml(registrationJSON, name, id)

	if err != nil {
		return nil, err
	}
	// Create a new XML document
	bodyDocument := etree.NewDocument()
	err = bodyDocument.ReadFromBytes(file)

	if err != nil {
		return nil, err
	}
	// Create an "ID" element and set its text to the provided 'id'
	itemXml, err := bodyDocument.WriteToBytes()
		
	if err != nil {
		return nil, err
	}

	return itemXml, nil
}



func GetCompanyProfileAb(request requestmodel.CompanyProfileXml) ([]byte, error)  {
	listData := []byte{}
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	existingDoc := etree.NewDocument()
	element := existingDoc.CreateElement(EXTERNAL_TYPE_LIST_AB)

	dbConn = dbConn.Model(databasemodel.AngggotaBursa{})

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	rows, err := dbConn.Rows()

	if err != nil {
		return listData, errors.New("Failed to get company ab: "+err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.AngggotaBursa
		err := dbConn.ScanRows(rows, &item)

		if err != nil {
			return listData, errors.New("Failed to get company ab - scanrows: "+err.Error())
		}
	

		id := item.ID
		if item.RegistrationJson == nil {
			log.Println(item.ID)
			return listData, errors.New("Failed to get company ab - registrationJson: is null")
		}

		itemXml, err := generateProfileXml(id, item.RegistrationJson, "Ab")
		
		if err != nil {
			return listData, errors.New("Failed to get company ab - generateProfileXml: "+err.Error())
		}


		itemXmlDocument := etree.NewDocument()
		if err := itemXmlDocument.ReadFromBytes(itemXml); err != nil {
			return nil, err
		}

		element.AddChild(itemXmlDocument.Root())

	}
	
	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company ab - dbConn.Error: "+dbConn.Error.Error())
	}
	
	listData, err = existingDoc.WriteToBytes()

	if err != nil {
		return listData, err
	}

	return listData, nil

}


func GetCompanyProfileParticipant(request requestmodel.CompanyProfileXml) ([]byte, error)  {
	listData := []byte{}
	existingDoc := etree.NewDocument()
	element := existingDoc.CreateElement(EXTERNAL_TYPE_LIST_PARTICIPANT)
	// Create a ew XML element from the new XML bytes
	
	
	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	dbConn = dbConn.Model(databasemodel.Participant{})

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	rows, err := dbConn.Rows()

	if err != nil {
		return listData, errors.New("Failed to get company participant -dbConn.Rows(): "+err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		var item databasemodel.Participant
		err := dbConn.ScanRows(rows, &item)
		
		if err != nil {
			return listData, errors.New("Failed to get company participant- dbConn.ScanRows: "+err.Error())
		}

		id := item.ID

		
		itemXml, err := generateProfileXml(id, item.RegistrationJson, "Participant")

		if err != nil {
			return listData, errors.New("Failed to get company participant - generateProfileXml: "+err.Error())
		}

		itemXmlDocument := etree.NewDocument()
		if err := itemXmlDocument.ReadFromBytes(itemXml); err != nil {
			return nil, err
		}


		element.AddChild(itemXmlDocument.Root())

	}

	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company participant - dbConn.Err: "+dbConn.Error.Error())
	}

	listData, err = existingDoc.WriteToBytes()

	if err != nil {
		return listData, err
	}

	return listData, nil
}

func GetCompanyProfilePJSPPA(request requestmodel.CompanyProfileXml) ([]byte, error)  {
	listData := []byte{}
	existingDoc := etree.NewDocument()
	element := existingDoc.CreateElement(EXTERNAL_TYPE_LIST_PJSPPA)

	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}

	dbConn = dbConn.Model(databasemodel.Pjsppa{})

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	rows, err := dbConn.Rows()

	if err != nil {
		return listData, errors.New("Failed to get company pjsppa -dbConn.Rows(): "+err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.Pjsppa
		err := dbConn.ScanRows(rows, &item)
		
		if err != nil {
			return listData, errors.New("Failed to get company pjsppa- dbConn.ScanRows: "+err.Error())
		}

		id := item.ID
		
		itemXml, err := generateProfileXml(id, item.RegistrationJson, "Pjsppa")
		
		if err != nil {
			return listData, errors.New("Failed to get company pjsppa - generateProfileXml: "+err.Error())
		}

		itemXmlDocument := etree.NewDocument()
		if err := itemXmlDocument.ReadFromBytes(itemXml); err != nil {
			return nil, err
		}


		element.AddChild(itemXmlDocument.Root())
	}


	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company pjsppa: "+dbConn.Error.Error())
	}

	listData, err = existingDoc.WriteToBytes()

	if err != nil {
		return listData, err
	}

	return listData, nil
}

func GetCompanyDealerUtama(request requestmodel.CompanyProfileXml) ([]byte, error)  {
	listData := []byte{}
	existingDoc := etree.NewDocument()
	element := existingDoc.CreateElement(EXTERNAL_TYPE_LIST_DU)

	dbConn, errInitDb := helper.InitDBConnGorm(request.ExternalType)
	if errInitDb != nil {
		return listData, errInitDb
	}
	
	dbConn = dbConn.Model(databasemodel.DealerUtama{})

	if len(request.CompanyCode) > 0 {
		dbConn = dbConn.Where("code = ?", request.CompanyCode)
	} 

	rows, err := dbConn.Rows()

	if err != nil {
		return listData, errors.New("Failed to get company dealerutama -dbConn.Rows(): "+err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var item databasemodel.DealerUtama
		err := dbConn.ScanRows(rows, &item)
		
		if err != nil {
			return listData, errors.New("Failed to get company dealerutama- dbConn.ScanRows: "+err.Error())
		}

		id := item.ID
		
		itemXml, err := generateProfileXml(id, item.RegistrationJson, "Du")
		
		if err != nil {
			return listData, errors.New("Failed to get company dealerutama - generateProfileXml: "+err.Error())
		}

		itemXmlDocument := etree.NewDocument()
		if err := itemXmlDocument.ReadFromBytes(itemXml); err != nil {
			return nil, err
		}

		element.AddChild(itemXmlDocument.Root())
	}


	if dbConn.Error != nil {
		return listData, errors.New("Failed to get company dealer utama: "+dbConn.Error.Error())
	}

	listData, err = existingDoc.WriteToBytes()

	if err != nil {
		return listData, err
	}

	return listData, nil
}