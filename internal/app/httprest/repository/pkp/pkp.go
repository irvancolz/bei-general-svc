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
	GetAllPKuser(c *gin.Context) ([]*model.PKuser, error)
	CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error)
	UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
	GetAllWithFilter(keyword []string) ([]*model.PKuser, error)
	GetAllWithSearch(Code string, Name string, QuestionDate time.Time, Question string, Answers string, answered_by string, AnsweredAt time.Time) ([]*model.PKuser, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAllWithSearch(Code string, Name string, QuestionDate time.Time, Question string, Answers string, answered_by string, AnsweredAt time.Time) ([]*model.PKuser, error) {
	var querySelect = `SELECT id, stakeholders, code, name, question_date, question, answers, answered_by, answered_at, topic
						FROM questions
						WHERE code = $1 
						AND name ILIKE '%' || $2 || '%' 
						AND question_date BETWEEN $3 AND $4 
						AND (answered_at BETWEEN $5 AND $6 OR answered_at IS NULL)
						AND deleted_at IS NULL`

	var listData = []*model.PKuser{}
	selDB, err := m.DB.Query(querySelect, Code, Name, QuestionDate, Question, Answers, answered_by, AnsweredAt)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [questions] [sqlQuery] [GetAllWithSearch] ", err)
		return nil, err
	}
	for selDB.Next() {
		pk := model.PKuser{}
		err = selDB.Scan(
			&pk.ID,
			&pk.Stakeholders,
			&pk.Code,
			&pk.Name,
			&pk.QuestionDate,
			&pk.Question,
			&pk.Answers,
			&pk.AnswersBy,
			&pk.AnswersAt,
			&pk.Topic,
		)
		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [questions] [sqlQuery] [GetAllWithSearch] ", err)
			return nil, err
		}
		listData = append(listData, &pk)
	}
	return listData, nil
}

func (m *repository) GetAllWithFilter(keyword []string) ([]*model.PKuser, error) {

	var querySelect = "select name, code, question_date, answers_by, topic, from PKusers where is_deleted = false"

	if len(keyword) > 3 {
		return nil, errors.New("keyword more than three")
	}

	var query, oneKeyword, twoKeywords, threeKeyword string

	if len(keyword) == 1 {
		oneKeyword = `
		SELECT * FROM (
			` + querySelect + `
		) AS a 
		WHERE a.name ILIKE '%` + keyword[0] + `%' 
			OR a.code ILIKE '%` + keyword[0] + `%' 
			OR a.question_date ILIKE '%` + keyword[0] + `%'  
			OR a.answers_by ILIKE '%` + keyword[0] + `%'
			OR a.topic ILIKE '%` + keyword[0] + `%'`
		query = oneKeyword
	} else if len(keyword) == 2 {
		oneKeyword = `
		SELECT * FROM (
			` + querySelect + `
		) AS a 
		WHERE a.name ILIKE '%` + keyword[0] + `%' 
			OR a.code ILIKE '%` + keyword[0] + `%' 
			OR a.question_date ILIKE '%` + keyword[0] + `%'  
			OR a.answers_by ILIKE '%` + keyword[0] + `%'
			OR a.topic ILIKE '%` + keyword[0] + `%'`

		twoKeywords = `
		SELECT * FROM (
			` + oneKeyword + `
		) AS b
		WHERE b.name ILIKE '%` + keyword[1] + `%'
			OR b.code ILIKE '%` + keyword[1] + `%'
			OR b.question_date ILIKE '%` + keyword[1] + `%'
			OR b.answers_by ILIKE '%` + keyword[1] + `%'
			OR b.topic ILIKE '%` + keyword[1] + `%'`
		query = twoKeywords
	} else if len(keyword) == 3 {
		oneKeyword = `
		SELECT * FROM (
			` + querySelect + `
		) AS a 
		WHERE a.name ILIKE '%` + keyword[0] + `%' 
			OR a.code ILIKE '%` + keyword[0] + `%' 
			OR a.question_date ILIKE '%` + keyword[0] + `%'  
			OR a.answers_by ILIKE '%` + keyword[0] + `%'
			OR a.topic ILIKE '%` + keyword[0] + `%'`

		twoKeywords = `
		SELECT * FROM (
			` + oneKeyword + `
		) AS b
		WHERE b.name ILIKE '%` + keyword[1] + `%'
			OR b.code ILIKE '%` + keyword[1] + `%'
			OR b.question_date ILIKE '%` + keyword[1] + `%'
			OR b.answers_by ILIKE '%` + keyword[1] + `%'
			OR b.topic ILIKE '%` + keyword[1] + `%'`

		threeKeyword = `
		SELECT * FROM (
			` + twoKeywords + `
		) AS c
		WHERE c.name ILIKE '%` + keyword[2] + `%'
			OR c.code ILIKE '%` + keyword[2] + `%'
			OR c.question_date ILIKE '%` + keyword[2] + `%'
			OR c.answers_by ILIKE '%` + keyword[2] + `%'
			OR c.topic ILIKE '%` + keyword[2] + `%'`
		query = threeKeyword
	} else {
		query = querySelect
	}
	var listData = []*model.PKuser{}
	selDB, err := m.DB.Query(query)
	if err != nil {
		return nil, errors.New("not found")
	}
	for selDB.Next() {
		pkp := model.PKuser{}
		err = selDB.Scan(
			&pkp.ID,
			&pkp.Name,
			&pkp.Code,
			&pkp.QuestionDate,
			&pkp.AnswersBy,
			&pkp.Topic,
		)
		if err != nil {
			return nil, errors.New("not found")
		}
		listData = append(listData, &pkp)
	}
	return listData, nil
}

func (m *repository) GetAllPKuser(c *gin.Context) ([]*model.PKuser, error) {
	result := []*model.PKuser{}
	query := `
	SELECT 
		id,
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
		created_at
	FROM PKuser;`
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKuser] [sqlQuery] [GetAllPKuser] ", err)
		return nil, errors.New("list complaints not found")
	}
	defer rows.Close()
	for rows.Next() {
		var item model.PKuser

		err := rows.Scan(
			&item.ID,
			&item.Stakeholders,
			&item.Code,
			&item.Name,
			&item.QuestionDate,
			&item.Question,
			&item.Answers,
			&item.AnswersAt,
			&item.Topic,
			&item.FileName,
			&item.FilePath,
		)

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [pkp] [getQueryData] [GetAllComplaints] ", err)
			return nil, errors.New("failed when retrieving data")
		}
		result = append(result, &item)
	}

	return result, nil
}

func (m *repository) CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error) {
	UserId, _ := c.Get("id")
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
			create_by,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`

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
		UserId,
		t,
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
	userId, _ := c.Get("id")
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
			created_by = $13,
			created_at = $14,
			updated_by = $15,
			updated_at = $16, 
		WHERE id = $1 
			AND deleted_by IS NULL
			AND is_deleted = false;
	`

	UpdatedAt := time.Now().UTC()
	CreatedAt := time.Now().UTC()

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
		pkp.Topic,
		pkp.FileName,
		pkp.FilePath,
		pkp.CreateBy,
		CreatedAt,
		pkp.UpdatedBy,
		UpdatedAt,
		userId,
	)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [PKP] [updateData] ", err)
		return 0, errors.New("update pkp failed")
	}

	rowsAffected, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (m *repository) Delete(id string, c *gin.Context) (int64, error) {
	userId, _ := c.Get("user_id")
	deleted_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	query := `
	UPDATE
	pkp SET
		is_deleted = true,
		deleted_by = $2,
		deleted_at = $3
	WHERE 
		id = $1 
		AND
		is_deleted = false;`
	selDB, err := m.DB.Exec(query, id, userId, deleted_at)
	if err != nil {
		return 0, err
	}

	RowsAffected, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return RowsAffected, nil
}
