package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"log"

	"github.com/gin-gonic/gin"
)

func uploadParticipantNoteToDb(c *gin.Context, referenceNumber, svcName string, uploadedData [][]string, removeFile func()) {
	DbConn, errCreateConn := helper.InitDBConnGorm(svcName)
	if errCreateConn != nil {
		removeFile()
		log.Println("failed create connection to upload data :", errCreateConn)
		return
	}

	if uploadedData != nil {
		if len(uploadedData) > 2 {
			catatanParticipantList := []databasemodel.Notes{}
			for i := 3; i < len(uploadedData); i++ {
				catatanParticipant := databasemodel.Notes{
					ReferenceNo:       referenceNumber,
					UploadDate:        safeAccess(uploadedData, i, 3),
					ParticipantCode:   safeAccess(uploadedData, i, 4),
					ParticipantName:   safeAccess(uploadedData, i, 5),
					EventDate:         safeAccess(uploadedData, i, 6),
					Category:          safeAccess(uploadedData, i, 7),
					ReportDescription: safeAccess(uploadedData, i, 8),
					Action:            safeAccess(uploadedData, i, 9),
					BursaUser:         safeAccess(uploadedData, i, 10),
					Description:       safeAccess(uploadedData, i, 11),
					CreatedBy:         "Unggah Berkas",
				}

				catatanParticipantList = append(catatanParticipantList, catatanParticipant)
			}

			DbConn.Create(&catatanParticipantList)
			if DbConn.Error != nil {
				removeFile()
			}
		}
	}
}

func safeAccess(slice [][]string, columnIndex, rowIndex int) string {
	if rowIndex >= 0 && rowIndex < len(slice[columnIndex]) {
		return slice[columnIndex][rowIndex]
	}
	return ""
}
