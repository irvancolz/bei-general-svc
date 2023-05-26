package unggahberkas

import (
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/unggah-berkas"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UnggahBerkasUsecaseInterface interface {
	UploadNew(c *gin.Context, props UploadNewFilesProps) (int64, error)
	GetUploadedFiles(c *gin.Context, filetypes string) ([]model.UploadedFilesMenuResponse, error)
	DeleteUploadedFiles(c *gin.Context, id string) error
}

type usecase struct {
	Repo repo.UnggahBerkasRepoInterface
}

type UploadNewFilesProps struct {
	Id          string
	Type        string `validate:"oneof:catatan kunjungan bulanan pjsppa"`
	Report_Code string
	Report_Name string
	File_Name   string
	File_Path   string
	File_Size   int64
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

func (u *usecase) GetUploadedFiles(c *gin.Context, filetypes string) ([]model.UploadedFilesMenuResponse, error) {

	results, errorResults := u.Repo.GetUploadedFiles()
	if errorResults != nil {
		return nil, errorResults
	}

	if filetypes != "" {
		var resultByType []model.UploadedFilesMenuResponse
		for _, item := range results {
			if strings.EqualFold(filetypes, item.Report_Type) {
				resultByType = append(resultByType, item)
			}
		}

		return resultByType, nil
	}

	return results, nil
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
