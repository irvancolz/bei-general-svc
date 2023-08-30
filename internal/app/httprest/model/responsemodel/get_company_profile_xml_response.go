package responsemodel

import "be-idx-tsg/internal/app/httprest/model/databasemodel"

type CompanyProfileResponseXml struct {
	AnggotaBursaList 	[]databasemodel.AngggotaBursa
	ParticipantList 	[]databasemodel.Participant
	PjsppaList 			[]databasemodel.Pjsppa
	DealerUtamaList 	[]databasemodel.DealerUtama
}

