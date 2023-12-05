package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	repo "be-idx-tsg/internal/app/httprest/repository/unggah-berkas"
	"be-idx-tsg/internal/app/httprest/usecase/upload"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UnggahBerkasUsecaseInterface interface {
	UploadNew(c *gin.Context, props UploadNewFilesProps) (int64, error)
	GetUploadedFiles(c *gin.Context) (*helper.PaginationResponse, error)
	DeleteUploadedFiles(c *gin.Context, id string) error
}

type usecase struct {
	Repo repo.UnggahBerkasRepoInterface
}

type UploadNewFilesProps struct {
	Type      string `json:"type" binding:"oneof='catatan' 'kunjungan' 'bulanan' 'pjsppa' 'bulanan ab',required"`
	File_Name string `json:"file_name"`
	File_Path string `json:"file_path"`
	File_Size int64  `json:"file_size"`
	Periode   int64  `json:"periode" binding:"required"`
}

func NewUsecase() UnggahBerkasUsecaseInterface {
	return &usecase{
		Repo: repo.NewRepository(),
	}
}

func (u *usecase) UploadNew(c *gin.Context, props UploadNewFilesProps) (int64, error) {
	userName, _ := c.Get("name_user")
	userId, _ := c.Get("user_id")
	companyName, _ := c.Get("company_name")
	companyCode, _ := c.Get("company_code")
	companyId, _ := c.Get("company_id")
	userType, _ := c.Get("type")

	createNewArgs := repo.UploadNewFilesProps{
		Type: props.Type,
		Company_code: func() string {
			if strings.EqualFold(userType.(string), "internal") {
				return "BEI"
			}
			return companyCode.(string)
		}(),
		Company_name: func() string {
			if strings.EqualFold(userType.(string), "internal") {
				return "BURSA EFEK INDONESIA"
			}
			return companyName.(string)
		}(),
		Company_id:  companyId.(string),
		File_Name:   props.File_Name,
		File_Path:   props.File_Path,
		File_Size:   props.File_Size,
		Periode:     props.Periode,
		Is_Uploaded: true,
		Created_by:  userName.(string),
		Created_at:  time.Now().Unix(),
	}

	if props.File_Size <= 0 || props.File_Path == "" {
		createNewArgs.Is_Uploaded = false
	}

	referenceNumber := buildNoReference(props.Type, time.Now(), u.Repo.CurrentFileUploadedOrderToday(props.Type))

	go uploadReportToDb(c, props.File_Path, props.Type, referenceNumber)

	uploadRes, errUpload := u.Repo.UploadNew(createNewArgs)
	if errUpload != nil {
		return 0, errUpload
	}

	moduleName := func() string {
		if strings.EqualFold(props.Type, "bulanan") {
			return "Laporan Rekapitulasi Aktivitas Transaksi Partisipan"
		}
		if strings.EqualFold(props.Type, "pjsppa") {
			return "Laporan Rekapitulasi Aktivitas Transaksi PJSPPA"
		}
		if strings.EqualFold(props.Type, "catatan") {
			return "Laporan Rekapitulasi Catatan"
		}
		return "Laporan Historis Kunjungan Partisipan"
	}()

	notifMsg := "Laporan Baru Telah Berhasil Di Upload"
	emailMsg := fmt.Sprintf("%s telah melakukan Upload data pada Modul Unggah Data transaksi terkait %s", userName.(string), moduleName)
	notifType := "Unggah Data"

	utilities.CreateNotif(c, userId.(string), notifType, "Laporan Berhasil Diupload")
	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)

	return uploadRes, nil
}

func (u *usecase) GetUploadedFiles(c *gin.Context) (*helper.PaginationResponse, error) {

	dataStruct, errorData := u.Repo.GetUploadedFiles(c)
	if errorData != nil {
		return nil, errorData
	}
	var dataToConverted []interface{}
	for _, item := range dataStruct {
		dataToConverted = append(dataToConverted, item)
	}
	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at", "periode"})

	collumnName := []string{"Periode", "Diupload Oleh", "Tanggal Upload", "Ukuran"}
	columnWidthFloat := []float64{40, 50, 60, 40}
	var columnWidthInt []int

	var txtFileHeaders [][]string
	txtFileHeaders = append(txtFileHeaders, collumnName)

	for _, width := range columnWidthFloat {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	dataOrder := []string{"periode", "created_by", "created_at", "file_size"}
	var exportedData [][]string

	for _, content := range filteredData {
		var item []string
		item = append(item, helper.MapToArray(content, dataOrder)...)

		for i, content := range item {
			if i == 2 {
				unixTime, _ := strconv.Atoi(content)
				dateTime := time.Unix(int64(unixTime), 0)
				dateToFormat := helper.GetWIBLocalTime(&dateTime)
				item[i] = helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthFullNameInIndo) + " - " + helper.GetTimeAndMinuteOnly(dateToFormat)
			}
			if i == 0 {
				unixTime, _ := strconv.Atoi(content)
				dateTime := time.Unix(int64(unixTime), 0)
				dateToFormat := helper.GetWIBLocalTime(&dateTime)
				item[i] = helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthFullNameInIndo)
			}
			if i == len(item)-1 {
				fileSizeInInt, _ := strconv.Atoi(content)
				filesSizeInMb := fileSizeInInt / 1000
				item[i] = fmt.Sprintf("%v MB", filesSizeInMb)
			}
		}

		exportedData = append(exportedData, item)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename:    "Unggah Data",
		Headers:     txtFileHeaders,
		ColumnWidth: columnWidthInt,
		Data:        exportedData,
		ExcelConfig: &helper.ExportToExcelConfig{},
		PdfConfig: &helper.PdfTableOptions{
			HeaderRows: helper.GenerateTableHeaders(collumnName, columnWidthFloat),
		},
	}

	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return nil, errorExport
	}

	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)
	return &paginatedData, nil
}

func (u *usecase) DeleteUploadedFiles(c *gin.Context, id string) error {
	userName, _ := c.Get("name_user")

	isFileAvaliable := u.Repo.CheckFileAvaliability(id)
	if !isFileAvaliable {
		return errors.New("failed to delete files, files cannot found in database")
	}

	filePath, errorPath := u.Repo.GetUploadedFilesPath(c, id)
	if errorPath != nil {
		return errorPath
	}
	removeFileFromDiskArgs := upload.UploadFileConfig{}
	removeFileFromDiskErr := upload.NewUsecase().DeleteFile(c, removeFileFromDiskArgs, filePath)
	if removeFileFromDiskErr != nil {
		return removeFileFromDiskErr
	}

	deleteFileArgs := repo.DeleteUploadedFilesProps{
		Id:         id,
		Deleted_at: time.Now().Unix(),
		Deleted_by: userName.(string),
	}

	return u.Repo.DeleteUploadedFiles(deleteFileArgs)
}
