package pkp

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAllPKuser(c *gin.Context) ([]model.PKuser, error)
	CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error)
	UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAllPKuser(c *gin.Context) ([]model.PKuser, error) {
	result := []model.PKuser{}
	query := `
	SELECT 
		id,
		stakeholders,
		code,
		name,
		question_date AS QuestionDate,
		question,
		answers,
		answers_by AS AnswersBy,
		answers_at AS AnswersAt,
		topic,
		external_type AS ExternalType,
		additional_information AS AdditionalInfo,
		file_name AS filename,
		file_path AS filepath,
		created_by AS CreateBy,
		created_at AS CreatedAt,
		updated_by AS UpdatedBy,
		updated_at AS UpdatedAt,
		deleted_by AS DeletedBy,
		deleted_at AS DeletedAt
	FROM pkp 
	WHERE deleted_by IS NULL
	AND deleted_at IS NULL`

	orderQuery := ` ORDER BY
	CASE
		WHEN updated_at IS NOT NULL 
			THEN updated_at
		ELSE created_at
	END DESC`

	searchQueryConfig := helper.SearchQueryGenerator{
		ColumnScanned: []string{
			"name",
			"code",
			"question",
			"answers_by",
			"topic",
			"answers",
			"stakeholders",
			"file_name",
			"created_by",
			"additional_information",
		},
		TableName: "pkp",
	}

	getAllQuery := searchQueryConfig.GenerateGetAllDataQuerry(c, query)
	getAllQuery += orderQuery

	authDb, errGetAuthCon := helper.InitDBConn("auth")
	if errGetAuthCon != nil {
		log.Println("failed inisiate connection to get pkp user :", errGetAuthCon)
		return nil, errGetAuthCon
	}
	defer authDb.Close()

	rows, err := m.DB.Queryx(getAllQuery)
	if err != nil {
		log.Println(getAllQuery)
		log.Println("[AQI-debug] [err] [repository] [PKuser] [sqlQuery] [GetAllPKuser] ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item model.PKuserResultSet

		err := rows.StructScan(
			&item,
		)

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [pkp] [getQueryData] [GetAllComplaints] ", err)
			return nil, err
		}

		data := model.PKuser{
			ID:             item.ID,
			Stakeholders:   item.Stakeholders,
			Code:           item.Code,
			Name:           item.Name,
			QuestionDate:   item.QuestionDate.Unix(),
			Question:       item.Question,
			Answers:        item.Answers,
			AnswersBy:      item.AnswersBy,
			AnswersAt:      item.AnswersAt.Unix(),
			Topic:          item.Topic,
			FileName:       item.FileName,
			FilePath:       item.FilePath,
			ExternalType:   item.ExternalType,
			AdditionalInfo: item.AdditionalInfo.String,
			CreatedAt:      item.CreatedAt.Unix(),
			CreateBy:       item.CreateBy,

			UpdatedBy: item.UpdatedBy.String,
			DeletedBy: item.DeletedBy.String,
			UpdatedAt: func() int64 {
				if item.UpdatedBy.Valid {
					return item.UpdatedAt.Time.Unix()
				}
				return 0
			}(),
			DeletedAt: func() int64 {
				if item.DeletedBy.Valid {
					return item.DeletedAt.Time.Unix()
				}
				return 0
			}(),
		}

		result = append(result, data)
	}

	return result, nil
}

func (m *repository) CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error) {
	userName, _ := c.Get("name_user")
	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")

	QuestionDate := pkp.QuestionDate
	AnswersAt := pkp.AnswersAt

	query := `
		INSERT INTO 
		pkp (
			stakeholders,
			code,
			name,
			question_date,
			question,
			answers,
			answers_by,
			answers_at,
			topic,
			file_name,
			file_path,
			created_by,
			created_at,
			external_type,
			additional_information
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`

	result, err := m.DB.Exec(
		query,
		pkp.Stakeholders,
		pkp.Code,
		pkp.Name,
		QuestionDate,
		pkp.Question,
		pkp.Answers,
		pkp.AnswersBy,
		AnswersAt,
		pkp.Topic,
		pkp.FileName,
		pkp.FilePath,
		userName.(string),
		t,
		pkp.ExternalType,
		pkp.AdditionalInfo,
	)

	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (m *repository) UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error) {
	userName, _ := c.Get("name_user")
	query := `
		UPDATE
			pkp SET
			stakeholders = $2,
			code = $3,
			name = $4,
			question_date = $5,
			question = $6,
			answers = $7,
			answers_by = $8,
			answers_at = $9,
			topic = $10,
			file_name = $11,
			file_path = $12,
			updated_by = $13,
			updated_at = $14,
			external_type = $15,
			additional_information = $16
		WHERE id = $1 
			AND deleted_by IS NULL
			AND  deleted_at IS NULL;
	`

	UpdatedAt := time.Now().UTC()

	selDB, err := m.DB.Exec(
		query,
		pkp.ID,
		pkp.Stakeholders,
		pkp.Code,
		pkp.Name,
		pkp.QuestionDate,
		pkp.Question,
		pkp.Answers,
		pkp.AnswersBy,
		pkp.AnswersAt,
		pkp.Topic,
		pkp.FileName,
		pkp.FilePath,
		userName.(string),
		UpdatedAt,
		pkp.ExternalType,
		pkp.AdditionalInfo,
	)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKP] [updateData] failed to update PKP data: ", err)
		return 0, err
	}

	rowsAffected, err := selDB.RowsAffected()
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKP] [updateData] failed to retrieve rows affected: ", err)
		return 0, err
	}

	if rowsAffected == 0 {
		err = errors.New("pkp failed to updated, please check the id and try again")
		log.Println("[AQI-debug] [err] [repository] [PKP] [updateData] ", err)
		return 0, err
	}

	log.Printf("[AQI-debug] [info] [repository] [PKP] [updateData] %d rows affected\n", rowsAffected)
	return rowsAffected, nil
}

func (m *repository) Delete(id string, c *gin.Context) (int64, error) {
	userName, _ := c.Get("name_user")
	deleted_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	query := `
	UPDATE
	pkp SET
		deleted_by = $2,
		deleted_at = $3
	WHERE 
		id = $1 
		AND deleted_by IS NULL
		AND deleted_at IS NULL;`
	selDB, err := m.DB.Exec(query, id, userName.(string), deleted_at)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKP] [Delete] ", err)
		return 0, err
	}

	rowsAffected, err := selDB.RowsAffected()
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKP] [Delete] ", err)
		return 0, err
	}

	if rowsAffected == 0 {
		err = errors.New("failed to delete pkp, please specify the id and try again")
		log.Println("[AQI-debug] [warn] [repository] [PKP] [Delete] ", err)
		return 0, err
	}

	log.Println("[AQI-debug] [info] [repository] [PKP] [Delete] PKP deleted successfully")
	return rowsAffected, nil
}
