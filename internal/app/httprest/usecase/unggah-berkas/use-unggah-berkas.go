package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	repo "be-idx-tsg/internal/app/httprest/repository/unggah-berkas"
	"errors"
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
	Type        string `json:"type" binding:"oneof=catatan kunjungan bulanan pjsppa,required"`
	Report_Code string `json:"report_code"`
	Report_Name string `json:"report_name"`
	File_Name   string `json:"file_name"`
	File_Path   string `json:"file_path"`
	File_Size   int64  `json:"file_size"`
}

func NewUsecase() UnggahBerkasUsecaseInterface {
	return &usecase{
		Repo: repo.NewRepository(),
	}
}

func (u *usecase) UploadNew(c *gin.Context, props UploadNewFilesProps) (int64, error) {
	userName, _ := c.Get("name_user")
	createNewArgs := repo.UploadNewFilesProps{
		Type:        props.Type,
		Report_Code: props.Report_Code,
		Report_Name: props.Report_Name,
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
	filteredData := helper.HandleDataFiltering(c, dataToConverted, []string{"created_at", "updated_at"})
	paginatedData := helper.HandleDataPagination(c, filteredData)
	return &paginatedData, nil
}

func (u *usecase) DeleteUploadedFiles(c *gin.Context, id string) error {
	userName, _ := c.Get("name_user")

	isFileAvaliable := u.Repo.CheckFileAvaliability(id)
	if !isFileAvaliable {
		return errors.New("failed to delete files, files cannot found in database")
	}

	deleteFileArgs := repo.DeleteUploadedFilesProps{
		Id:         id,
		Deleted_at: time.Now().Unix(),
		Deleted_by: userName.(string),
	}

	return u.Repo.DeleteUploadedFiles(deleteFileArgs)
}
