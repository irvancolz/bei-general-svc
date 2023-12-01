package upload

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/utilities"
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
	DownloadFromLocal(c *gin.Context) error
	DeleteFile(c *gin.Context, props UploadFileConfig, slug string) error
	IsFileExists(c *gin.Context, slug string) error
}

type usecase struct{}

func NewUsecase() UploadFileUsecaseInterface {
	return &usecase{}
}

func (u *usecase) Upload(c *gin.Context, props UploadFileConfig) (*model.UploadFileResponse, error) {

	result := &model.UploadFileResponse{}
	cleanDisk := func(file string) {
		errRemoveFile := os.Remove(file)
		if errRemoveFile != nil {
			log.Println("failed to cleanup directory : ", errRemoveFile)
			return
		}
	}

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

	if len(newFileName) > 240 {
		errorText := "failed to save file to server : file name is too long"
		log.Println(errorText)
		return nil, errors.New(errorText)
	}
	// save file temporary to local server before uploaded to minio
	errSave := c.SaveUploadedFile(file, file.Filename)
	if errSave != nil {
		log.Println("failed to save file to server : ", errSave)
		return nil, errSave
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:     strings.ToLower(props.Directory),
		FileOriginName: file.Filename,
		FileSavedName:  newFileName,
	}

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		cleanDisk(file.Filename)
		return nil, errorCreateMinionConn
	}
	_, errorUploadToMinio := helper.UploadToMinio(minioClient, c, minioSaveConfig)
	if errorUploadToMinio != nil {
		cleanDisk(file.Filename)
		return nil, errorUploadToMinio
	}

	result.FileName = newFileName
	result.FileSize = file.Size
	result.OgFileName = file.Filename
	result.Filepath = props.GenerateFilePath(newFileName)

	// cleanup
	cleanDisk(file.Filename)
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
	bucketName := GetFilesBucket(slug)
	prohibitedExt := []string{".go", ".env", ".dev", ".yml", ".sql"}

	// get file extension from filename
	ext := filepath.Ext(slug)

	isExtAvailable := props.CheckFileExt(ext)
	if !isExtAvailable || helper.IsContains(prohibitedExt, ext) {
		return errors.New("the file ext does not supported")
	}

	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:    bucketName,
		FileSavedName: fileLocation,
	}

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return errorCreateMinionConn
	}

	err := helper.DeleteFileInMinio(minioClient, c, minioSaveConfig)
	if err != nil {
		log.Println("failed to remove file : ", err)
		return err
	}

	if strings.EqualFold(bucketName, "form") {
		utilities.UpdateFormAttachmentFileStatus(c, fileLocation)
	}

	return nil
}

func (u *usecase) IsFileExists(c *gin.Context, slug string) error {

	fileLocation := GetFilePath(slug)
	minioSaveConfig := helper.UploadToMinioProps{
		BucketName:    GetFilesBucket(slug),
		FileSavedName: fileLocation,
	}

	minioClient, errorCreateMinionConn := helper.InitMinio()
	if errorCreateMinionConn != nil {
		return errorCreateMinionConn
	}

	err := helper.CheckIsObjExists(minioClient, c, minioSaveConfig)
	if err != nil {
		log.Println("failed to remove file : ", err)
		return err
	}

	return nil
}

func (u usecase) DownloadFromLocal(c *gin.Context) error {
	reportType := c.DefaultQuery("report", "")

	fileName := func() string {
		if strings.EqualFold(reportType, "lraktp") {
			return "static/Laporan Rekapitulasi Aktivitas Transaksi Partisipan.xlsx"
		}
		if strings.EqualFold(reportType, "ltbab") {
			return "Laporan Transaksi Bulanan Anggota Bursa.xlsx"
		}
		if strings.EqualFold(reportType, "lraktpjsppa") {
			return "static/Laporan Rekapitulasi Aktivitas Transaksi PJSPPA.xlsx"
		}
		if strings.EqualFold(reportType, "lhkp") {
			return "static/Laporan Historis Kunjungan Partisipan.xlsx"
		}
		if strings.EqualFold(reportType, "notepar") {
			return "static/Laporan Rekapitulasi Catatan.xlsx"
		}
		return ""

	}()

	if fileName == "" {
		return errors.New("cannot find the report template, please select the correct type : lraktp || lraktpjsppa || lhkp")
	}

	c.File(fileName)

	return nil
}
