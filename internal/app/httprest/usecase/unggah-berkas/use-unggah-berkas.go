package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	repo "be-idx-tsg/internal/app/httprest/repository/unggah-berkas"
	"be-idx-tsg/internal/app/httprest/usecase/upload"
	"errors"
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
	Type      string `json:"type" binding:"oneof=catatan kunjungan bulanan pjsppa,required"`
	File_Name string `json:"file_name"`
	File_Path string `json:"file_path"`
	File_Size int64  `json:"file_size"`
}

func NewUsecase() UnggahBerkasUsecaseInterface {
	return &usecase{
		Repo: repo.NewRepository(),
	}
}

func (u *usecase) UploadNew(c *gin.Context, props UploadNewFilesProps) (int64, error) {
	userName, _ := c.Get("name_user")
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
		Is_Uploaded: true,
		Created_by:  userName.(string),
		Created_at:  time.Now().Unix(),
	}

	if props.File_Size <= 0 || props.File_Path == "" {
		createNewArgs.Is_Uploaded = false
	}

	return u.Repo.UploadNew(createNewArgs)
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
	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})

	collumnName := []string{"Jenis Lampiran", "Kode", "Nama", "Tanggal", "Status"}
	columnWidthFloat := []float64{50, 20, 50, 40, 30}
	var columnWidthInt []int

	var txtFileHeaders [][]string
	txtFileHeaders = append(txtFileHeaders, collumnName)

	for _, width := range columnWidthFloat {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	dataOrder := []string{"type", "company_code", "company_name", "created_at", "is_uploaded"}
	var exportedData [][]string

	for _, content := range filteredData {
		var item []string
		item = append(item, helper.MapToArray(content, dataOrder)...)

		for i, content := range item {
			if helper.IsContains([]int{3}, i) {
				unixTime, _ := strconv.Atoi(content)
				dateTime := time.Unix(int64(unixTime), 0)
				dateToFormat := helper.GetWIBLocalTime(&dateTime)
				item[i] = helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthFullNameInIndo) + " - " + helper.GetTimeAndMinuteOnly(dateToFormat)
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
