package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model/databasemodel"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)






func uploadParticipantNoteToDb(c *gin.Context, pathFile, reportType, referenceNumber string) {
	// download uploaded files

	// read the files
	fileLocation := getUnggahBerkasFile(c, pathFile)

	if len(fileLocation) == 0 {
		return
	}

	removeFile := func() error { return os.Remove(fileLocation) }

	// read the files
	uploadedData := helper.ReadFileExcel(fileLocation)
	svcName := getDbSvcName(reportType)

	if svcName == "participant" {
		DbConn, errCreateConn := helper.InitDBConnGorm(svcName)
		if errCreateConn != nil {
			removeFile()
			log.Println("failed create connection to upload data :", errCreateConn)
			return
		}

		if uploadedData != nil && len(uploadedData) >2 {
			catatanParticipantList := []databasemodel.Notes{}
			for i := 3; i < len(uploadedData); i++ {
				for j := 1; j < len(uploadedData[i]); j++ {
					catatanParticipant := databasemodel.Notes{
						ID:              safeAccess(uploadedData, i, 1),
						ReferenceNo:     safeAccess(uploadedData, i, 2),
						UploadDate:      safeAccess(uploadedData, i, 3),
						ParticipantCode: safeAccess(uploadedData, i, 4),
						ParticipantName: safeAccess(uploadedData, i, 5),
						EventDate:       safeAccess(uploadedData, i, 6),
						Category:        safeAccess(uploadedData, i, 7),
						ReportDescription:     safeAccess(uploadedData, i, 8),
						Action:          safeAccess(uploadedData, i, 9),
						BursaUser:       safeAccess(uploadedData, i, 10),
						Description: 	 safeAccess(uploadedData, i, 11),
					}

					catatanParticipantList = append(catatanParticipantList, catatanParticipant)
				}
			}
		
			DbConn.Save(&catatanParticipantList)
		}

		defer func() {
			dbInstance, _ := DbConn.DB()
			_ = dbInstance.Close()
		}()
	}
}

func safeAccess(slice [][]string, columnIndex, rowIndex int) string {
	if rowIndex >= 0 && rowIndex < len(slice) {
		return slice[columnIndex][rowIndex]
	}
	return ""
}