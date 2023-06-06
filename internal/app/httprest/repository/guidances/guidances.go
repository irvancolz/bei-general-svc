package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
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
type DeleteExistingDataProps struct {
	Id         string
	Deleted_by string
	Deleted_at time.Time
}

type GuidancesRepoInterface interface {
	CreateNewData(props CreateNewDataProps) (int64, error)
	GetAllData(c *gin.Context) ([]model.GuidanceFileAndRegulationsJSONResponse, error)
	UpdateExistingData(params UpdateExistingDataProps) error
	DeleteExistingData(params DeleteExistingDataProps) error
	CheckIsOrderFilled(order int) bool
	UpdateOrder(order int) error
}

func NewGuidancesRepository() GuidancesRepoInterface {
	return &guidancesRepository{
		DB: *database.Init().MySql,
	}
}

func (r *guidancesRepository) CreateNewData(props CreateNewDataProps) (int64, error) {
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
		},
	}
	query := serchQueryConfig.GenerateGetAllDataQuerry(c, getAllDataQuerry)

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

func (r *guidancesRepository) UpdateExistingData(params UpdateExistingDataProps) error {
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

func (r *guidancesRepository) CheckIsOrderFilled(order int) bool {
	var result int
	rowResult := r.DB.QueryRowx(checkIsOrderFilledQuery, order)
	errorGetRows := rowResult.Scan(&result)
	if errorGetRows != nil {
		log.Println("failed to check guidances order avaliability :", errorGetRows)
		return true
	}
	return result > 0
}

func (r *guidancesRepository) UpdateOrder(order int) error {
	execResult, errorExec := r.DB.Exec(updateOrderQuery, order)
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
