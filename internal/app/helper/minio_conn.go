package helper

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() (*minio.Client, error) {

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
		log.Println(minoEndpointPath)
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

		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": "s3:GetObject",
					"Resource": [
						"arn:aws:s3:::` + props.BucketName + `/*"
					]
				}
			]
		}`

		errCreateBucket := client.MakeBucket(c, props.BucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
		if errCreateBucket != nil {
			log.Println("failed to create new bucket :", errCreateBucket)
			return nil, errCreateBucket
		}

		errSetBucketPolicy := client.SetBucketPolicy(c, props.BucketName, policy)
		if errSetBucketPolicy != nil {
			log.Println("failed to set bucket to public access :", errSetBucketPolicy)
			return nil, errSetBucketPolicy
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

	errGetFile := CheckIsObjExists(client, c, props)
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

func CheckIsObjExists(client *minio.Client, c context.Context, props UploadToMinioProps) error {
	_, errGetStat := client.StatObject(c, props.BucketName, props.FileSavedName, minio.StatObjectOptions{})
	if errGetStat != nil {
		log.Println("failed to get obj stat :", errGetStat)
	}
	return errGetStat
}
