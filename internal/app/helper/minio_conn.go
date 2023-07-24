package helper

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() (*minio.Client, error) {

	errGetEnv := godotenv.Load()
	if errGetEnv != nil {
		log.Println("failed to read .env file :", errGetEnv)
		return nil, errGetEnv
	}

	minoEndpointPath := os.Getenv("MINIO_ENDPOINT")
	minioAcceskeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	minioSecretAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioUseSSLMode, _ := strconv.ParseBool(os.Getenv("MINIO_SSL_MODE"))

	// Initialize minio client object.
	minioClient, err := minio.New(minoEndpointPath, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAcceskeyID, minioSecretAccessKey, ""),
		Secure: minioUseSSLMode,
	})

	if err != nil {
		log.Println("failed to connect to minio server :", err)
		return nil, err
	}

	return minioClient, nil
}

type UploadToMinioProps struct {
	BucketName     string
	FileOriginName string
	FileSavedName  string
}

func UploadToMinio(client *minio.Client, c context.Context, props UploadToMinioProps) (*minio.UploadInfo, error) {

	isBucketExists, errorCheckBucketExists := client.BucketExists(c, props.BucketName)
	if errorCheckBucketExists != nil {
		log.Println("failed to check the specified bucket location: ", errorCheckBucketExists)
		return nil, errorCheckBucketExists
	}

	if !isBucketExists {
		errCreateBucket := client.MakeBucket(c, props.BucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
		if errCreateBucket != nil {
			log.Println("failed to create new bucket :", errCreateBucket)
			return nil, errCreateBucket
		}
	}

	result, errorUpload := client.FPutObject(c, props.BucketName, props.FileSavedName, props.FileOriginName, minio.PutObjectOptions{})
	if errorUpload != nil {
		log.Println("failed to upload file to minio :", errorUpload)
		return nil, errorUpload
	}

	return &result, nil
}

func GetFileFromMinio(client *minio.Client, c context.Context, props UploadToMinioProps) error {
	errorGetObj := client.FGetObject(c, props.BucketName, props.FileSavedName, props.FileSavedName, minio.GetObjectOptions{})
	if errorGetObj != nil {
		log.Println("failed to get file from minio :", errorGetObj)
		return errorGetObj
	}
	return nil
}

func DeleteFileInMinio(client *minio.Client, c context.Context, props UploadToMinioProps) error {
	deleteProps := minio.RemoveObjectOptions{
		ForceDelete: false,
	}

	errGetFile := GetFileFromMinio(client, c, props)
	if errGetFile != nil {
		return errGetFile
	}

	errorDelete := client.RemoveObject(c, props.BucketName, props.FileSavedName, deleteProps)

	if errorDelete != nil {
		log.Println("failed to delete file in minio :", errorDelete)
		return errorDelete
	}

	return nil
}
