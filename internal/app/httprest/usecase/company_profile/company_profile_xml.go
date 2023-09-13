package companyprofile

import (
	"be-idx-tsg/internal/app/httprest/model/requestmodel"
	companyprofilerepository "be-idx-tsg/internal/app/httprest/repository/company-profile"
	"strings"

	"github.com/beevik/etree"
	"github.com/gin-gonic/gin"
)

type onGetPjsppaList = func([]byte)
type onGetParticipantList = func([]byte)
type onGetAbList = func([]byte)
type onGetDuList = func([]byte)
type onGetAllList = func([]byte)

const (
	EXTERNAL_TYPE_LIST_AB = "AbList"
	EXTERNAL_TYPE_LIST_PARTICIPANT = "ParticipantList"
	EXTERNAL_TYPE_LIST_PJSPPA = "PjsppaList"
	EXTERNAL_TYPE_LIST_DU = "DuList"

	REQUEST_EXTERNAL_TYPE_AB = "ab"
	REQUEST_EXTERNAL_TYPE_PARTICIPANT = "participant"
	REQUEST_EXTERNAL_TYPE_PJSPPA = "pjsppa"
	REQUEST_EXTERNAL_TYPE_DU = "du"
)

func byteToDocument(data []byte) (*etree.Document, error) {
	// Create a ew XML element from the new XML bytes
	document := etree.NewDocument()
	if err := document.ReadFromBytes(data); err != nil {
		return nil, err
	}

	// Add the new element to the existing document


	//log.Println("hey")
	//data, err := abCompanyList.WriteToBytes()

	//log.Println(string(data))



	return document, nil
}

func GetCompanyProfileXml(c *gin.Context, request requestmodel.CompanyProfileXml) ([]byte, error) {

	companyProfileXmlBytes := []byte{}

	err := handleCompanyType(request,
		func(abList []byte) {
			companyProfileXmlBytes = append(companyProfileXmlBytes, abList...) 
		},
		func(participantList []byte) {
			companyProfileXmlBytes = append(companyProfileXmlBytes, participantList...) 
		},
		func(pjsppaList []byte) {
			companyProfileXmlBytes = append(companyProfileXmlBytes, pjsppaList...) 
		},
		func(duList []byte) {
			companyProfileXmlBytes = append(companyProfileXmlBytes, duList...) 
		},
		func(allList []byte) {
			companyProfileXmlBytes = append(companyProfileXmlBytes, allList...) 
		},
	)

	if err != nil {
		return companyProfileXmlBytes, err
	}

	return companyProfileXmlBytes, nil
}

func handleCompanyType(request requestmodel.CompanyProfileXml,
	onGetAbList onGetAbList,
	onGetParticipantList onGetParticipantList,
	onGetPjsppaList onGetPjsppaList,
	onGetDuList onGetDuList,
	onGetAllList onGetAllList,
) error {


	companyList := etree.NewDocument()
	companyList.CreateElement("CompanyProfile")


	if len(request.ExternalType) == 0 {
		var allList[]byte  
		request.ExternalType = REQUEST_EXTERNAL_TYPE_AB
		anggotaBursaList, err := companyprofilerepository.GetCompanyProfileAb(request)

		if err != nil {
			return err
		}

		_, err = combineXml(companyList,  anggotaBursaList)

		if err != nil {
			return err
		}
		

		request.ExternalType = REQUEST_EXTERNAL_TYPE_PARTICIPANT
		participantList, err := companyprofilerepository.GetCompanyProfileParticipant(request)

		if err != nil {
			return err
		}

		_, err = combineXml(companyList,  participantList)

		if err != nil {
			return err
		}

		request.ExternalType = REQUEST_EXTERNAL_TYPE_PJSPPA
		pjsppaList, err := companyprofilerepository.GetCompanyProfilePJSPPA(request)

		if err != nil {
			return err
		}

		_, err = combineXml(companyList, pjsppaList)

		if err != nil {
			return err
		}

		request.ExternalType = REQUEST_EXTERNAL_TYPE_DU
		dealerUtamaList, err := companyprofilerepository.GetCompanyDealerUtama(request)

		if err != nil {
			return err
		}

		allList, err = combineXml(companyList, dealerUtamaList)

		if err != nil {
			return err
		}


		onGetAllList(allList)
	} else if strings.EqualFold(request.ExternalType, REQUEST_EXTERNAL_TYPE_AB) {
		anggotaBursaList, err := companyprofilerepository.GetCompanyProfileAb(request)

		if err != nil {
			return err
		}

		anggotaBursaList, err = combineXml(companyList, anggotaBursaList)

		if err != nil {
			return err
		}

		onGetAbList(anggotaBursaList)

	} else if strings.EqualFold(request.ExternalType, REQUEST_EXTERNAL_TYPE_PARTICIPANT) {
		participantList, err := companyprofilerepository.GetCompanyProfileParticipant(request)

		if err != nil {
			return err
		}

		participantList, err = combineXml(companyList, participantList)

		if err != nil {
			return err
		}

		onGetParticipantList(participantList)
	} else if strings.EqualFold(request.ExternalType, REQUEST_EXTERNAL_TYPE_PJSPPA) {
		pjsppaList, err := companyprofilerepository.GetCompanyProfilePJSPPA(request)

		if err != nil {
			return err
		}

		pjsppaList, err = combineXml(companyList, pjsppaList)

		if err != nil {
			return err
		}

		onGetPjsppaList(pjsppaList)

	} else if strings.EqualFold(request.ExternalType, REQUEST_EXTERNAL_TYPE_DU) {
		dealerUtamaList, err := companyprofilerepository.GetCompanyDealerUtama(request)

		if err != nil {
			return err
		}

		dealerUtamaList, err = combineXml(companyList, dealerUtamaList)

		if err != nil {
			return err
		}

		onGetDuList(dealerUtamaList)
	}

	return nil
}



func combineXml(companyList *etree.Document, data []byte)  ([]byte, error) {

	dataDocument, err :=  byteToDocument(data)

	if err != nil {
		return nil, err
	}

	companyList.Root().AddChild(dataDocument.Root())
	data, err = companyList.WriteToBytes()

	if err != nil {
		return nil, err
	}

	return data, nil
}