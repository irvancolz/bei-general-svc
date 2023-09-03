package companyprofile

import (
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"be-idx-tsg/internal/app/httprest/model/requestmodel"
	"be-idx-tsg/internal/app/httprest/model/responsemodel"
	companyprofilerepository "be-idx-tsg/internal/app/httprest/repository/company-profile"
	"strings"

	"github.com/gin-gonic/gin"
)

type onGetPjsppaList = func([]databasemodel.Pjsppa)
type onGetParticipantList = func([]databasemodel.Participant)
type onGetAbList = func([]databasemodel.AngggotaBursa)
type onGetDuList = func([]databasemodel.DealerUtama)

func GetCompanyProfileXml(c *gin.Context, request requestmodel.CompanyProfileXml) (responsemodel.CompanyProfileResponseXml, error) {

	companyProfileXml := responsemodel.CompanyProfileResponseXml{}

	err := handleCompanyType(request,
		func(abList []databasemodel.AngggotaBursa) {
			companyProfileXml.AnggotaBursaList = abList
		},
		func(participantList []databasemodel.Participant) {
			companyProfileXml.ParticipantList = participantList
		},
		func(pjsppaList []databasemodel.Pjsppa) {
			companyProfileXml.PjsppaList = pjsppaList
		},
		func(duList []databasemodel.DealerUtama) {
			companyProfileXml.DealerUtamaList = duList
		},
	)

	if err != nil {
		return companyProfileXml, err
	}

	return companyProfileXml, nil
}

func handleCompanyType(request requestmodel.CompanyProfileXml,
	onGetAbList onGetAbList,
	onGetParticipantList onGetParticipantList,
	onGetPjsppaList onGetPjsppaList,
	onGetDuList onGetDuList,
) error {

	if len(request.ExternalType) == 0 {
		anggotaBursaList, err := companyprofilerepository.GetCompanyProfileAb(request)

		if err != nil {
			return err
		}

		participantList, err := companyprofilerepository.GetCompanyProfileParticipant(request)

		if err != nil {
			return err
		}

		pjsppaList, err := companyprofilerepository.GetCompanyProfilePJSPPA(request)

		if err != nil {
			return err
		}

		onGetPjsppaList(pjsppaList)

		dealerUtamaList, err := companyprofilerepository.GetCompanyDealerUtama(request)

		if err != nil {
			return err
		}

		onGetAbList(anggotaBursaList)
		onGetParticipantList(participantList)
		onGetPjsppaList(pjsppaList)
		onGetDuList(dealerUtamaList)
	} else if strings.EqualFold(request.ExternalType, "ab") {
		anggotaBursaList, err := companyprofilerepository.GetCompanyProfileAb(request)

		if err != nil {
			return err
		}
		onGetAbList(anggotaBursaList)

	} else if strings.EqualFold(request.ExternalType, "participant") {
		participantList, err := companyprofilerepository.GetCompanyProfileParticipant(request)

		if err != nil {
			return err
		}
		onGetParticipantList(participantList)
	} else if strings.EqualFold(request.ExternalType, "pjsppa") {
		pjsppaList, err := companyprofilerepository.GetCompanyProfilePJSPPA(request)

		if err != nil {
			return err
		}

		onGetPjsppaList(pjsppaList)

	} else if strings.EqualFold(request.ExternalType, "du") {
		dealerUtamaList, err := companyprofilerepository.GetCompanyDealerUtama(request)

		if err != nil {
			return err
		}

		onGetDuList(dealerUtamaList)
	}

	return nil
}
