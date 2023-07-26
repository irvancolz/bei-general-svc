package unggahberkas

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UnggahBerkasRepoInterface interface {
	UploadNew(props UploadNewFilesProps) (int64, error)
	GetUploadedFiles(c *gin.Context) ([]model.UploadedFilesMenuResponse, error)
	GetUploadedFilesPath(c *gin.Context, id string) (string, error)
	DeleteUploadedFiles(props DeleteUploadedFilesProps) error
	CheckFileAvaliability(id string) bool
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() UnggahBerkasRepoInterface {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (r *repository) CheckFileAvaliability(id string) bool {
	var total int64
	result := r.DB.QueryRow(checkDataAvaliabilityQuery, id)
	errorScan := result.Scan(&total)
	if errorScan != nil {
		log.Println("failed to get file avaliability from database :", errorScan)
		return false
	}
	return total > 0
}

type UploadNewFilesProps struct {
	Type         string `validate:"oneof:catatan kunjungan bulanan pjsppa"`
	Company_code string
	Company_name string
	Company_id   string
	Is_Uploaded  bool
	File_Name    string
	File_Path    string
	File_Size    int64
	Created_by   string
	Created_at   int64
}

func (r *repository) UploadNew(props UploadNewFilesProps) (int64, error) {
	execResult, errorExec := r.DB.Exec(uploadNewFilesQuery,
		props.Type,
		props.Company_code,
		props.Company_name,
		props.Company_id,
		props.Is_Uploaded,
		props.File_Name,
		props.File_Path,
		props.File_Size,
		props.Created_by,
		props.Created_at)
	if errorExec != nil {
		log.Println("failed when try record uploaded files to database :", errorExec)
		return 0, errorExec
	}

	result, errorResult := execResult.RowsAffected()
	if errorResult != nil {
		log.Println("failed when try to check uploaded files in databse :", errorResult)
		return 0, errorResult
	}

	if result == 0 {
		log.Println("the files is failed to uploaded to database, please try again")
		return 0, errors.New("the files is failed to uploaded to database, please try again")
	}

	return result, nil
}

func (r *repository) GetUploadedFiles(c *gin.Context) ([]model.UploadedFilesMenuResponse, error) {
	var results []model.UploadedFilesMenuResponse
	serchQueryConfig := helper.SearchQueryGenerator{
		TableName: "uploaded_files",
		ColumnScanned: []string{
			"company_code",
			"company_name",
			"file_name",
			"id",
		},
	}

	query := serchQueryConfig.GenerateGetAllDataQuerry(c, getUploadedFilesQuery)
	rowResults, errorRows := r.DB.Queryx(query)
	if errorRows != nil {
		log.Println("failed to get uploaded files from databases :", errorRows)
		return nil, errorRows
	}
	defer rowResults.Close()

	for rowResults.Next() {
		var mock model.UploadedFilesMenuResultSet
		errorScan := rowResults.StructScan(&mock)
		if errorScan != nil {
			log.Println("failed when try to parsing data from database :", errorScan)
			return nil, errorScan
		}
		result := model.UploadedFilesMenuResponse{
			Id:           mock.Id,
			Type:         mock.Type,
			Company_code: mock.Company_code,
			Company_name: mock.Company_name,
			Company_id:   mock.Company_id,
			Is_Uploaded:  mock.Is_Uploaded,
			Created_By:   mock.Created_By,
			Created_At:   mock.Created_At,

			Updated_By: mock.Updated_By.String,
			Updated_At: mock.Updated_At.Int64,
			File_Size:  mock.File_Size.Int64,
			File_Path:  mock.File_Path.String,
			File_Name:  mock.File_Name.String,
		}
		results = append(results, result)
	}

	return results, nil
}

type DeleteUploadedFilesProps struct {
	Id         string
	Deleted_at int64
	Deleted_by string
}

func (r *repository) DeleteUploadedFiles(props DeleteUploadedFilesProps) error {
	txResult, errorTx := r.DB.Exec(deleteUploadedFilesQuery, props.Id, props.Deleted_by, props.Deleted_at)
	if errorTx != nil {
		log.Println("failed to deleted the uploaded files from database :", errorTx)
		return errorTx
	}
	result, errorResult := txResult.RowsAffected()
	if errorResult != nil {
		log.Println("failed when try to check deleted files from databse :", errorResult)
		return errorResult
	}

	if result == 0 {
		log.Println("the files is failed to be deleted to database, please try again")
		return errors.New("the files is failed to be deleted to database, please try again")
	}

	return nil
}

func (r *repository) GetUploadedFilesPath(c *gin.Context, id string) (string, error) {
	var result string
	rowResults := r.DB.QueryRowx(getUploadedFilesPathQuery, id)
	errorScan := rowResults.Scan(&result)
	if errorScan != nil {
		log.Println("failed to scan file path from database :", errorScan)
		return "", errorScan
	}

	return result, nil
}
