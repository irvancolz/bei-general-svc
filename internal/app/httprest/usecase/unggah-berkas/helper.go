package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/usecase/upload"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func formatDataToSlice(data [][]string, columnType []string) [][]interface{} {
	var result [][]interface{}
	if len(columnType) <= 0 {
		for _, item := range data {
			var row []interface{}
			for _, column := range item {
				row = append(row, column)
			}

			result = append(result, row)
		}
		return result
	}

	for _, columns := range data {
		formattedCol := []interface{}{}
		columnLength := len(columns)

		for i, types := range columnType {
			content := func() interface{} {
				if i >= columnLength && strings.EqualFold(types, "number") {
					return 0
				}
				if i >= columnLength && !strings.EqualFold(types, "number") {
					return ""
				}
				if strings.EqualFold(types, "number") {
					formattedNumber, _ := strconv.ParseInt(columns[i], 10, 64)
					return formattedNumber
				}
				return columns[i]
			}()

			formattedCol = append(formattedCol, content)
		}

		result = append(result, formattedCol)
	}
	return result
}

func getUnggahBerkasFile(c *gin.Context, pathFile string) string {
	fileLocation := upload.GetFilePath(pathFile)
	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return ""
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:    upload.GetFilesBucket(pathFile),
		FileSavedName: fileLocation,
	}

	errGetFromMinio := helper.GetFileFromMinio(minioClient, c, minioSaveConfig)
	if errGetFromMinio != nil {
		return ""
	}

	return fileLocation
}

func getDbSvcName(reportType string) string {
	if strings.EqualFold(reportType, "pjsppa") {
		return "pjsppa"
	}
	if strings.EqualFold(reportType, "bulanan ab") {
		return "ab"
	}

	return "participant"
}

func uploadReportToDb(c *gin.Context, pathFile, reportType, referenceNumber string) {
	// download uploaded files
	fileLocation := getUnggahBerkasFile(c, pathFile)
	removeFile := func() {
		err := os.Remove(fileLocation)
		if err != nil {
			log.Println("failed to delete file ", fileLocation, " ", err)
		}
	}

	// read the files
	uploadedData := helper.ReadFileExcel(fileLocation)
	svcName := getDbSvcName(reportType)

	headerHeight := func() int {
		if strings.EqualFold("bulanan", reportType) {
			return 3
		}
		if strings.EqualFold("bulanan ab", reportType) {
			return 1
		}
		if strings.EqualFold("pjsppa", reportType) {
			return 2
		}
		return 4
	}()

	uploadReportQuery := generateUploadReportQuery(reportType)
	sliceFormat := generateSliceFormatter(reportType)
	formattedData := formatDataToSlice(uploadedData[headerHeight:], sliceFormat)

	if reportType == "catatan" && svcName == "participant" {
		uploadParticipantNoteToDb(c, referenceNumber, svcName, uploadedData, removeFile)
	} else {
		DbConn, errCreateConn := helper.InitDBConn(svcName)
		if errCreateConn != nil {
			removeFile()
			log.Println("failed create connection to upload data :", errCreateConn)
			return
		}

		uploadStmt, errorCrateStmt := DbConn.Preparex(uploadReportQuery)
		if errorCrateStmt != nil {
			removeFile()
			log.Println("failed create statement :", errorCrateStmt)
			return
		}

		for _, row := range formattedData {
			row = append(row, referenceNumber, "unggah berkas", helper.GetWIBLocalTime(nil))
			errorInsert := sendRecordToDB(uploadStmt, reportType, row)
			if errorInsert != nil {
				removeFile()
				log.Println("failed to upload data report to database :", errorInsert)
			}
		}

		// delete the saved files
		removeFile()
	}

}

func generateUploadReportQuery(reportType string) string {
	if strings.EqualFold(reportType, "bulanan") {
		return `INSERT
		INTO activity_transaction_report (
			participant_code,
			participant_name,
			period,
			value,
			volume,
			frequency,
			number_of_late_reporting,
			double_report,
			trade_cancel,
			wrong_input,
			no_ref,
			created_by,
			created_at ) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
			)`
	}
	if strings.EqualFold(reportType, "bulanan ab") {
		return `INSERT
		INTO monthly_report (
			company_code,
			report_date,
			company_name,
			total_b,
			total_j,
			no_ref,
			created_by,
			created_at ) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8
			)`
	}
	if strings.EqualFold(reportType, "pjsppa") {
		return `INSERT INTO activity_transaction_report(
			pjsppa_code,
			delivery_date,
			trade_id,
			date,
			instrument,
			trader,
			initiator,
			initiator_trader,
			direction,
			dealer_trader,
			dealer,
			ticker,
			size,
			price,
			settle_date,
			type,
			status,
			value,
			yield,
			accured,
			total_value,
			plte_number,
			cp_switching,
			iceberg,
			currency,
			settlement_currency,
			user_bursa,
			description,
			no_ref,
			created_by,
			created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28, $29, $30, $31
		)`
	}
	return `INSERT INTO visits (
		delivery_date,
		participant_code,
		participant_name,
		visit_date,
		infrastructure_condition,
		connection_and_performance,
		user_obedience_level,
		report_delay_cause,
		suggestions_and_feedback,
		description,
		no_ref,
		created_by,
		created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
}

// todo wait for document to be updated
func generateSliceFormatter(reportType string) []string {
	if strings.EqualFold(reportType, "bulanan") {
		return []string{"number", "string", "string", "string", "number", "number", "number", "number", "number", "number", "number"}
	}
	if strings.EqualFold(reportType, "bulanan ab") {
		return []string{"s", "s", "s", "number", "number"}
	}
	if strings.EqualFold(reportType, "pjsppa") {
		return []string{"s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s"}
	}

	return []string{"s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s", "s"}
}

// todo wait for document to be updated
func sendRecordToDB(stmt *sqlx.Stmt, reportType string, row []interface{}) error {
	var (
		errorInsert       error
		insertQueryResult sql.Result
	)

	if strings.EqualFold(reportType, "bulanan") {
		//  bulanan order index 1 - 13
		insertQueryResult, errorInsert = stmt.Exec(row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13])
	} else if strings.EqualFold(reportType, "pjsppa") {
		insertQueryResult, errorInsert = stmt.Exec(row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15], row[16], row[17], row[18], row[19], row[20], row[21], row[22], row[23], row[24], row[25], row[26], row[27], row[28], row[29], row[30], row[31], row[32], row[33])
	} else if strings.EqualFold(reportType, "bulanan ab") {
		// insert All in report
		insertQueryResult, errorInsert = stmt.Exec(row...)
	} else {
		//  kunjungan order index 3 - 15
		insertQueryResult, errorInsert = stmt.Exec(row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15])
	}

	if errorInsert != nil {
		log.Println("failed to upload report data to database :", errorInsert)
		return errorInsert
	}

	if i, _ := insertQueryResult.RowsAffected(); i == 0 {
		log.Printf("failed to upload data : %v to database", row)
	}

	return nil
}

func buildNoReference(reportType string, timeProvider time.Time, order int) string {

	reportName := func() string {
		if strings.EqualFold(reportType, "bulanan") {
			return "RATPAR"
		}
		if strings.EqualFold(reportType, "bulanan ab") {
			return "RTBAB"
		}
		if strings.EqualFold(reportType, "pjsppa") {
			return "RATPJ"
		}
		if strings.EqualFold(reportType, "catatan") {
			return "NOTEPAR"
		}
		return "VISITPAR"
	}()
	currReportOrder := fmt.Sprintf("%03d", order)
	currDate := func() string {
		currTime := strings.Split(timeProvider.Format(time.DateOnly), "-")

		// forgot bout how its works, should refactor later
		for i, j := 0, len(currTime)-1; i < j; i, j = i+1, j-1 {
			currTime[i], currTime[j] = currTime[j], currTime[i]
		}

		return strings.Join(currTime, "")
	}()

	return reportName + currReportOrder + currDate

}
