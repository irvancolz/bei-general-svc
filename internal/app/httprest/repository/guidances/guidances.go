package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	fileOrganizer "be-idx-tsg/internal/app/httprest/usecase/upload"
	"be-idx-tsg/internal/pkg/database"
	"database/sql"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type guidancesRepository struct {
	DB sqlx.DB
}

type GuidancesRepoInterface interface {
	CreateNewData(props CreateNewDataProps) (int64, error)
	insertNewData(props CreateNewDataProps) (int64, error)
	GetAllData(c *gin.Context) ([]model.GuidanceFileAndRegulationsJSONResponse, error)
	insertUpdatedData(params UpdateExistingDataProps) error
	UpdateExistingData(c *gin.Context, params UpdateExistingDataProps) error
	DeleteExistingData(params DeleteExistingDataProps) error
	CheckIsOrderFilled(order int, category string) bool
	UpdateOrder(order int, category string) error
	GetFileSavedPath(id string) string
}

func NewGuidancesRepository() GuidancesRepoInterface {
	return &guidancesRepository{
		DB: *database.Init().MySql,
	}
}

type CreateNewDataProps struct {
	Category    string `validate:"oneof=Guidebook File Regulation"`
	Description string
	Name        string
	Link        string
	File        string
	File_size   int64
	File_path   string
	File_Group  string
	File_Owner  string
	Order       int
	Version     string
	Created_by  string
	Created_at  time.Time
}

func (u *guidancesRepository) CreateNewData(props CreateNewDataProps) (int64, error) {
	createNewDataArgs := props
	if props.Order <= 0 {
		createNewDataArgs.Order = 1
	}

	isOrderFilled := u.CheckIsOrderFilled(props.Order, props.Category)
	if isOrderFilled {
		errorSetOrder := u.UpdateOrder(props.Order, props.Category)
		if errorSetOrder != nil {
			return 0, errorSetOrder
		}
	}

	result, error_result := u.insertNewData(createNewDataArgs)
	if error_result != nil {
		return 0, error_result
	}
	return result, nil
}

type UpdateExistingDataProps struct {
	Id          string
	Category    string `validate:"oneof=Guidebook File Regulation"`
	Name        string
	Description string
	Link        string
	File        string
	File_size   int64
	File_path   string
	File_Group  string
	File_Owner  string
	Order       int
	Version     string
	Updated_by  string
	Updated_at  time.Time
}

func (u *guidancesRepository) UpdateExistingData(c *gin.Context, props UpdateExistingDataProps) error {
	createNewDataArgs := props
	if props.Order <= 0 {
		createNewDataArgs.Order = 1
	}

	isOrderFilled := u.CheckIsOrderFilled(props.Order, props.Category)
	if isOrderFilled && props.Order != int(u.GetCurrentOrder(props.Id)) {
		errorSetOrder := u.UpdateOrder(props.Order, props.Category)
		if errorSetOrder != nil {
			return errorSetOrder
		}
	}

	// update the filesaved
	savedFile := u.GetFileSavedPath(props.Id)
	if props.File_path != savedFile && savedFile != "" {
		errRemove := fileOrganizer.NewUsecase().DeleteFile(c, fileOrganizer.UploadFileConfig{}, savedFile)
		if errRemove != nil {
			log.Println("failed to delete existing file saved :", errRemove)
			return errRemove
		}
	}

	error_result := u.insertUpdatedData(createNewDataArgs)
	if error_result != nil {
		return error_result
	}
	return nil
}

func (r *guidancesRepository) insertNewData(props CreateNewDataProps) (int64, error) {
	error_validate := helper.Validator().Struct(props)
	if error_validate != nil {
		log.Println("data is not passed validation ", error_validate)
		return 0, error_validate
	}

	insert_res, error_insert := r.DB.Exec(createNewDataQuerry,
		props.Category,
		props.Name,
		props.Description,
		props.Link,
		props.File,
		props.File_size,
		props.File_path,
		props.File_Group,
		props.File_Owner,
		props.Order,
		props.Version,
		props.Created_by,
		props.Created_at)

	if error_insert != nil {
		log.Println("failed to insert data to database : ", error_insert)
		return 0, error_insert
	}
	results, error_results := insert_res.RowsAffected()
	if error_results != nil {
		log.Println("failed to get affected rows after inserting new data : ", error_results)
		return 0, error_results
	}
	if results == 0 {
		log.Println("database is not updated after transactions")
		return 0, errors.New("DATABASE IS NOT UPDATED, PLEASE TRY AGAIN")
	}

	return results, nil
}

func (r *guidancesRepository) GetAllData(c *gin.Context) ([]model.GuidanceFileAndRegulationsJSONResponse, error) {
	var results []model.GuidanceFileAndRegulationsJSONResponse

	serchQueryConfig := helper.SearchQueryGenerator{
		TableName: "public.guidance_file_and_regulation",
		ColumnScanned: []string{
			"category",
			"name",
			"description",
			"link",
			"file",
			"file_group",
			"owner",
		},
	}
	orderQuery := ` ORDER BY
		CASE
			WHEN updated_at IS NOT NULL 
				THEN updated_at
			ELSE created_at
		END DESC`
	query := serchQueryConfig.GenerateGetAllDataQuerry(c, getAllDataQuerry) + orderQuery

	result_rows, error_rows := r.DB.Queryx(query)
	if error_rows != nil {
		log.Println(query)
		log.Println("failed to excecute script : ", error_rows)
		return nil, error_rows
	}
	defer result_rows.Close()
	for result_rows.Next() {
		var result_set model.GuidanceFileAndRegulationsResultSetResponse
		error_scan := result_rows.StructScan(&result_set)
		if error_scan != nil {
			log.Println("failed to convert result set to struct : ", error_scan)
			return nil, error_scan
		}

		result := model.GuidanceFileAndRegulationsJSONResponse{
			Id:         result_set.Id,
			Category:   result_set.Category,
			Name:       result_set.Name,
			File:       result_set.File,
			File_size:  result_set.File_size,
			Created_by: result_set.Created_by,
			Created_at: result_set.Created_at.Unix(),
			File_path:  result_set.File_path,

			Description: result_set.Description.String,
			Link:        result_set.Link.String,
			File_Group:  result_set.File_Group.String,
			File_Owner:  result_set.File_Owner.String,
			Order:       int(result_set.Order.Int32),
			Version:     result_set.Version.String,
			Updated_by:  result_set.Updated_by.String,
			Updated_at:  result_set.Updated_at.Time.Unix(),
		}
		if !result_set.Updated_at.Valid {
			result.Updated_at = 0
		}

		results = append(results, result)
	}

	sort.SliceStable(results, func(current, before int) bool {
		return results[current].Order < results[before].Order
	})

	return results, nil
}

func (r *guidancesRepository) insertUpdatedData(params UpdateExistingDataProps) error {
	error_validate := helper.Validator().Struct(params)
	if error_validate != nil {
		log.Println("data is not passed validation ", error_validate)
		return error_validate
	}

	updated_rows, error_update := r.DB.Exec(querryUpdate,
		params.Id,
		params.Category,
		params.Name,
		params.Description,
		params.Link,
		params.File,
		params.Version,
		params.Updated_by,
		params.Updated_at,
		params.File_size,
		params.File_path,
		params.File_Group,
		params.File_Owner,
		params.Order)
	if error_update != nil {
		log.Println("failed to excecute script to update data : ", error_update)
		return error_update
	}

	results, error_results := updated_rows.RowsAffected()
	if error_results != nil {
		log.Println("failed to get affected rows after inserting new data : ", error_results)
		return error_results
	}
	if results == 0 {
		log.Println("database is not updated after transactions")
		return errors.New("DATABASE IS NOT UPDATED, PLEASE TRY AGAIN")
	}

	return nil
}

type DeleteExistingDataProps struct {
	Id         string
	Deleted_by string
	Deleted_at time.Time
}

func (r *guidancesRepository) DeleteExistingData(params DeleteExistingDataProps) error {
	updated_rows, error_update := r.DB.Exec(querryDelete,
		params.Deleted_at,
		params.Deleted_by,
		params.Id)
	if error_update != nil {
		log.Println("failed to excecute script to update data : ", error_update)
		return error_update
	}

	results, error_results := updated_rows.RowsAffected()
	if error_results != nil {
		log.Println("failed to get affected rows after delete new data : ", error_results)
		return error_results
	}
	if results == 0 {
		log.Println("database is not updated after transactions")
		return errors.New("DATABASE IS NOT UPDATED, PLEASE TRY AGAIN")
	}
	return nil
}

func (r *guidancesRepository) CheckIsOrderFilled(order int, category string) bool {
	var result int
	rowResult := r.DB.QueryRowx(checkIsOrderFilledQuery, order, category)
	errorGetRows := rowResult.Scan(&result)
	if errorGetRows != nil {
		log.Println("failed to check guidances order avaliability :", errorGetRows)
		return true
	}
	return result > 0
}

func (r *guidancesRepository) UpdateOrder(order int, category string) error {
	execResult, errorExec := r.DB.Exec(updateOrderQuery, order, category)
	if errorExec != nil {
		log.Println("failed to update order on guidances :", errorExec)
		return errorExec
	}

	_, errorResult := execResult.RowsAffected()
	if errorResult != nil {
		log.Println("failed to get updated rows after editing order :", errorResult)
	}

	return nil
}

func (r *guidancesRepository) GetFileSavedPath(guidancesid string) string {

	filePath := r.DB.QueryRow(getFileSavedPathQuery, guidancesid)
	var result sql.NullString

	errorPath := filePath.Scan(&result)
	if errorPath != nil {
		log.Println("failed to get filepath saved :", errorPath)
		return ""
	}

	return result.String
}

func (r *guidancesRepository) GetCurrentOrder(id string) int64 {
	var result sql.NullInt64
	checkResult := r.DB.QueryRowx(getCurrentOrderQuery, id)

	errorScan := checkResult.Scan(&result)
	if errorScan != nil {
		log.Println("failed to get current files order :", errorScan)
		return 0
	}

	return result.Int64
}
