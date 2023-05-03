package guidances

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"
	"time"

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
	Version     float64
	Order       int64
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
	Version     float64
	Order       int64
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
	GetAllDataBasedOnCategory(category_type string) ([]*model.GuidanceFileAndRegulationsJSONResponse, error)
	UpdateExistingData(params UpdateExistingDataProps) error
	DeleteExistingData(params DeleteExistingDataProps) error
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

	createNewDataQuerry := `INSERT INTO public.guidance_file_and_regulation(
		category,
		name,
		description,
		link,
		file,
		file_size,
		version,
		"order",
		created_by,
		created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	insert_res, error_insert := r.DB.Exec(createNewDataQuerry,
		props.Category,
		props.Name,
		props.Description,
		props.Link,
		props.File,
		props.File_size,
		props.Version,
		props.Order,
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

func (r *guidancesRepository) GetAllDataBasedOnCategory(category_type string) ([]*model.GuidanceFileAndRegulationsJSONResponse, error) {
	var results []*model.GuidanceFileAndRegulationsJSONResponse
	getAllDataQuerry := `SELECT 
	id, 
	category, 
	name, 
	description, 
	link,
	file,
	file_size,
	version,
	"order",
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM public.guidance_file_and_regulation
	WHERE category = $1
	AND deleted_by IS NUll`

	result_rows, error_rows := r.DB.Queryx(getAllDataQuerry, category_type)
	if error_rows != nil {
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
			Id:          result_set.Id,
			Category:    result_set.Category,
			Name:        result_set.Name,
			Description: result_set.Description,
			Version:     result_set.Version,
			Order:       result_set.Order,
			Created_by:  result_set.Created_by,
			Created_at:  result_set.Created_at,
		}
		result.File = result_set.File.String
		result.File_size = result_set.File_size.String
		result.Link = result_set.Link.String
		result.Updated_at = result_set.Updated_at.Time
		result.Updated_by = result_set.Updated_by.String
		results = append(results, &result)
	}
	return results, nil
}

func (r *guidancesRepository) UpdateExistingData(params UpdateExistingDataProps) error {
	error_validate := helper.Validator().Struct(params)
	if error_validate != nil {
		log.Println("data is not passed validation ", error_validate)
		return error_validate
	}

	querryUpdate := `UPDATE public.guidance_file_and_regulation 
	SET category  = $1,
	name = $2,
	description = $3,
	link = $4,
	file = $5,
	version = $6,
	"order" = $7,
	updated_by = $8,
	updated_at = $9,
	file_size = $10
	WHERE id = $11
	AND category = $1`
	updated_rows, error_update := r.DB.Exec(querryUpdate,
		params.Category,
		params.Name,
		params.Description,
		params.Link,
		params.File,
		params.Version,
		params.Order,
		params.Updated_by,
		params.Updated_at,
		params.File_size,
		params.Id)
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
	querryDelete := `UPDATE public.guidance_file_and_regulation 
	SET deleted_at  = $1,
	deleted_by = $2
	WHERE id = $3`
	updated_rows, error_update := r.DB.Exec(querryDelete,
		time.Now(),
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
