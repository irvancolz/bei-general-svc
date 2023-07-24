package upload

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadFileUsecaseInterface interface {
	Upload(c *gin.Context, props UploadFileConfig) (*model.UploadFileResponse, error)
	Download(c *gin.Context, path string) error
	DeleteFile(c *gin.Context, props UploadFileConfig, slug string) error
}

type usecase struct{}

func NewUsecase() UploadFileUsecaseInterface {
	return &usecase{}
}

func (u *usecase) Upload(c *gin.Context, props UploadFileConfig) (*model.UploadFileResponse, error) {

	result := &model.UploadFileResponse{}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("failed to get the specified file : ", err)
		return nil, err
	}

	ext := filepath.Ext(file.Filename)
	if !props.CheckFileExt(ext) {
		log.Println("the file extension is not allowed to upload : ", ext)
		return nil, errors.New("the file extension is not allowed to upload")
	}
	if props.MaxSize > 0 {
		if !props.CheckFileSize(file.Size) {
			log.Println("the file size is exceed the maximum size allowed : ", file.Size)
			return nil, errors.New("the file size is exceed the maximum size allowed")
		}
	}

	newFileName := props.GenerateFilename(file.Filename, time.Now())

	// save file temporary to local server before uploaded to minio
	errSave := c.SaveUploadedFile(file, file.Filename)
	if errSave != nil {
		log.Println("failed to save file to server : ", errSave)
		return nil, errSave
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:     props.Directory,
		FileOriginName: file.Filename,
		FileSavedName:  newFileName,
	}

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return nil, errorCreateMinionConn
	}
	_, errorUploadToMinio := helper.UploadToMinio(minioClient, c, minioSaveConfig)
	if errorUploadToMinio != nil {
		return nil, errorUploadToMinio
	}

	result.FileName = newFileName
	result.FileSize = file.Size
	result.Filepath = props.GenerateFilePath(newFileName)

	// cleanup
	errRemoveFile := os.Remove(file.Filename)
	if errRemoveFile != nil {
		log.Println("failed to cleanup directory : ", err)
		return nil, errRemoveFile
	}

	return result, nil
}

func GetFilePath(path string) string {
	pat := filepath.FromSlash(path)
	pathStr := strings.Split(pat, string(os.PathSeparator))
	result := pathStr[len(pathStr)-1]
	return result
}

func GetFilesBucket(path string) string {
	pat := filepath.FromSlash(path)
	pathStr := strings.Split(pat, string(os.PathSeparator))
	result := pathStr[len(pathStr)-2]
	return result
}

func GetFileName(slug string) string {
	pathByOs := filepath.FromSlash(slug)
	fileSlug := strings.Split(pathByOs, string(os.PathSeparator))
	return fileSlug[len(fileSlug)-1]
}

func IsFIleExists(slug string) bool {
	fileLocation := GetFilePath(slug)
	fileName := GetFileName(slug)
	errorPath := filepath.Walk(filepath.FromSlash(fileLocation), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == fileName {
			log.Println(path)
		}
		return nil
	})

	if errorPath != nil {
		log.Println(errorPath)
		return false
	}

	return true
}

func (u *usecase) Download(c *gin.Context, pathFile string) error {
	fileLocation := GetFilePath(pathFile)

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return errorCreateMinionConn
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:    GetFilesBucket(pathFile),
		FileSavedName: fileLocation,
	}

	errGetFromMinio := helper.GetFileFromMinio(minioClient, c, minioSaveConfig)
	if errGetFromMinio != nil {
		return errGetFromMinio
	}

	c.File(fileLocation)

	errorCleanup := os.Remove(fileLocation)
	if errorCleanup != nil {
		log.Println("failed to do cleanup on downloaded files :", fileLocation)
		return nil
	}

	return nil
}

func (u *usecase) DeleteFile(c *gin.Context, props UploadFileConfig, slug string) error {

	if slug == "" {
		return nil
	}

	fileLocation := GetFilePath(slug)
	prohibitedExt := []string{".go", ".env", ".dev", ".yml", ".sql"}

	// get file extension from filename
	ext := filepath.Ext(slug)

	isExtAvailable := props.CheckFileExt(ext)
	if !isExtAvailable || helper.IsContains(prohibitedExt, ext) {
		return errors.New("the file ext does not supported")
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:    GetFilesBucket(slug),
		FileSavedName: fileLocation,
	}

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return errorCreateMinionConn
	}

	// err := os.Remove(fileLocation)
	err := helper.DeleteFileInMinio(minioClient, c, minioSaveConfig)
	if err != nil {
		log.Println("failed to remove file : ", err)
		return err
	}

	return nil
}
