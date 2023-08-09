package announcement

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByID(id string, c *gin.Context) (*model.Announcement, error)
	GetAllAnnouncement(c *gin.Context) ([]model.Announcement, error)
	Create(an model.CreateAnnouncement, c *gin.Context) (int64, error)
	Update(ab model.UpdateAnnouncement, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
	GetAllANWithFilter(keyword []string) ([]*model.Announcement, error)
	GetAllANWithSearch(Information_Type string, keyword string, startDate string, endDate string) ([]*model.Announcement, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAllANWithSearch(Information_Type string, keyword string, startDate string, endDate string) ([]*model.Announcement, error) {
	var querySelect = ` 
		select a.id, 
		a.information_type, 
		a.effective_date, 
		a.regarding,
		a.type 
		from 
		announcements a 
		where 
		a.information_type = $1
		and a.regarding ilike '%` + keyword + `%'
		and (a.effective_date between  $2 AND $3 ) 
		and a.deleted_by is null
	`

	var listData = []*model.Announcement{}
	selDB, err := m.DB.Query(querySelect, Information_Type, parseTime(startDate), parseTime(endDate))
	if err != nil {
		log.Println("time ", startDate, endDate)
		log.Println("[AQI-debug] [err] [repository] [Annoucement] [sqlQuery] [GetAllANWithSearch] ", err)
		return nil, errors.New("not found")
	}
	for selDB.Next() {
		an := model.Announcement{}
		err = selDB.Scan(
			&an.ID,
			&an.Information_Type,
			&an.Effective_Date,
			&an.Regarding,
			&an.Type,
		)
		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [Annoucement] [sqlQuery] [GetAllANWithSearch] ", err)
			return nil, errors.New("not found")
		}
		listData = append(listData, &an)
	}
	return listData, nil

}

func (m *repository) GetAllANWithFilter(keyword []string) ([]*model.Announcement, error) {

	var querySelect = " select an.id, an.code, an.type, an.status, an.license, an.operational_status  from announcements an where an.is_deleted = false "

	if len(keyword) > 3 {
		return nil, errors.New("keyword more than three ")
	}
	var query, oneKeyword, twoKeywords, threeKeyword string
	if len(keyword) == 1 {
		oneKeyword = `
	select z.* from (
		` + querySelect + `
		) as z 
		where
		z.code ilike '%` + keyword[0] + `%' 
		or z.name ilike '%` + keyword[0] + `%' 
		or z.status ilike '%` + keyword[0] + `%'  
		or z.license ilike '%` + keyword[0] + `%'  
		or z.operational_status ilike '%` + keyword[0] + `%' `
		query = oneKeyword
	} else if len(keyword) == 2 {
		oneKeyword = `
	select z.* from (
		` + querySelect + `
		) as z 
		where
		z.code ilike '%` + keyword[0] + `%' 
		or z.name ilike '%` + keyword[0] + `%' 
		or z.status ilike '%` + keyword[0] + `%'  
		or z.license ilike '%` + keyword[0] + `%'  
		or z.operational_status ilike '%` + keyword[0] + `%' `

		twoKeywords = `
		select x.* from (
			` + oneKeyword + `
		) as x 
		where 
		x.code ilike '%` + keyword[1] + `%'
		or x.name ilike '%` + keyword[1] + `%'
		or x.status ilike '%` + keyword[1] + `%' 
		or x.license ilike '%` + keyword[1] + `%' 
		or x.operational_status ilike '%` + keyword[1] + `%' `
		query = twoKeywords
	} else if len(keyword) == 3 {
		oneKeyword = `
	select z.* from (
		` + querySelect + `
		) as z 
		where
		z.code ilike '%` + keyword[0] + `%' 
		or z.name ilike '%` + keyword[0] + `%' 
		or z.status ilike '%` + keyword[0] + `%'  
		or z.license ilike '%` + keyword[0] + `%'  
		or z.operational_status ilike '%` + keyword[0] + `%' `

		twoKeywords = `
		select x.* from (
			` + oneKeyword + `
		) as x 
		where 
		x.code ilike '%` + keyword[1] + `%'
		or x.name ilike '%` + keyword[1] + `%'
		or x.status ilike '%` + keyword[1] + `%' 
		or x.license ilike '%` + keyword[1] + `%' 
		or x.operational_status ilike '%` + keyword[1] + `%' `

		threeKeyword = `
			select y.* from (
				` + twoKeywords + `) as y
			where
				y.code ilike '%` + keyword[2] + `%'
				or y.name ilike '%` + keyword[2] + `%'
				or y.status ilike '%` + keyword[2] + `%'
				or y.license ilike '%` + keyword[2] + `%'
				or y.operational_status ilike '%` + keyword[2] + `%'
		`
		query = threeKeyword
	} else {
		query = querySelect
	}
	var listData = []*model.Announcement{}
	selDB, err := m.DB.Queryx(query)
	if err != nil {
		return nil, errors.New("not found")
	}
	for selDB.Next() {
		an := model.UpdateAnnouncement{}
		err = selDB.StructScan(&an)
		if err != nil {
			return nil, errors.New("not found")
		}

		result := model.Announcement{}
		listData = append(listData, &result)
	}
	return listData, nil
}

func (m *repository) GetAllAnnouncement(c *gin.Context) ([]model.Announcement, error) {
	userType, _ := c.Get("type")
	ExternalType, _ := c.Get("external_type")
	filterQuery := ""
	if strings.ToLower(userType.(string)) == "internal" {
		filterQuery = "where deleted_by IS NULL"
	} else {
		str, _ := ExternalType.(string)
		if str != "" {
			filterQuery = "where lower(information_type) in (lower('SEMUA'), lower(" + str + ")) and deleted_by IS NULL "
		} else {
			filterQuery = "where information_type in ('SEMUA') and deleted_by IS NULL "
		}
	}

	result := []model.Announcement{}
	getAllQuery := `
	SELECT 
	id,
	information_type,
	effective_date,
	regarding,
	created_by,
	type
   FROM announcements
	` + filterQuery

	serchQueryConfig := helper.SearchQueryGenerator{
		TableName: "announcements",
		ColumnScanned: []string{
			"regarding",
		},
	}
	query := serchQueryConfig.GenerateGetAllDataQuerry(c, getAllQuery)

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println(query)
		log.Println("[AQI-debug] [err] [repository] [Annoucement] [sqlQuery] [GetAllAnnouncement] ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item model.UpdateAnnouncement
		var userId string
		var anEffectiveDate time.Time
		err := rows.Scan(
			&item.ID,
			&item.Information_Type,
			&anEffectiveDate,
			&item.Regarding,
			&userId,
			&item.Type,
		)

		announcement := model.Announcement{
			ID:               item.ID,
			Information_Type: item.Information_Type,
			Regarding:        item.Regarding,
			Type:             item.Type,
			Creator:          userId,
		}
		announcement.Effective_Date = anEffectiveDate.Unix()

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [Annoucement] [getQueryData] [GetAllAnnouncement] ", err)
			return nil, errors.New("failed when retrieving data")
		}

		result = append(result, announcement)
	}

	return result, nil
}

func (m *repository) GetByID(id string, c *gin.Context) (*model.Announcement, error) {
	var creatorId string
	var efectiveDateTime time.Time
	query := `
		SELECT 
		id,
		information_type,
		effective_date,
		created_by,
		type,
		regarding
		FROM announcements
		WHERE id = $1 
		AND deleted_by IS NULL
		ORDER BY effective_date DESC
		`
	item := &model.Announcement{}

	if err := m.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.Information_Type,
		&efectiveDateTime,
		&creatorId,
		&item.Type,
		&item.Regarding,
	); err != nil {
		log.Println(query)
		log.Println("[AQI-debug] [err] [repository] [Annoucement] [getQueryData] [GetByID] ", err)
		return nil, errors.New("announcement Not Found")
	}

	item.Creator = utilities.GetUserNameByID(c, creatorId)
	item.Effective_Date = efectiveDateTime.Unix()

	return item, nil
}

func IsExternalUser(c *gin.Context) bool {
	userType, _ := c.Get("type")
	return strings.EqualFold(userType.(string), "external")
}

func (m *repository) Create(an model.CreateAnnouncement, c *gin.Context) (int64, error) {
	userId, _ := c.Get("user_id")

	if !IsExternalUser(c) {
		return 0, errors.New("you dont have permission to create announcement")
	}

	t, _ := helper.TimeIn(time.Now(), "Asia/Jakarta")
	CreatedAt := t.Format("2006-01-02 15:04:05")
	query := `
	INSERT INTO 
	announcements (
		information_type,
		effective_date,
		regarding,
		created_at, 
		created_by,
		is_deleted,
		type
	)
	VALUES ($1, $2, $3, $4, $5, false, $6);`
	Effective_Date := an.Effective_Date
	Effective_DateParse, _ := time.Parse(time.RFC3339, Effective_Date)
	selDB, err := m.DB.Exec(
		query,
		an.Information_Type,
		Effective_DateParse,
		an.Regarding,
		CreatedAt,
		userId,
		an.Type)
	if err != nil {
		return 0, err
	}

	LastInsertId, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return LastInsertId, nil
}

func (m *repository) Update(an model.UpdateAnnouncement, c *gin.Context) (int64, error) {
	userId, _ := c.Get("user_id")

	if !IsExternalUser(c) {
		return 0, errors.New("you dont have permission to edit announcement")
	}

	query := `
	UPDATE
		announcements SET
		information_type = $2,
		effective_date = $3,
		regarding = $4,
		updated_at = $5, 
		updated_by = $6,
		type = $7
	WHERE id = $1 AND is_deleted = false;`
	updated_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	selDB, err := m.DB.Exec(
		query,
		an.ID,
		an.Information_Type,
		an.Effective_Date,
		an.Regarding,
		updated_at,
		userId,
		an.Type,
	)
	if err != nil {
		return 0, err
	}

	RowsAffected, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return RowsAffected, nil
}

func (m *repository) Delete(id string, c *gin.Context) (int64, error) {
	userId, _ := c.Get("user_id")
	deleted_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	query := `
	UPDATE
	announcements SET
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

func parseTime(input string) string {
	// parse input string menjadi time.Time object
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		log.Println("error parsing time:", err)
		return ""
	}

	// set timezone yang diinginkan
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("error loading location:", err)
		return ""
	}

	// konversi time.Time object ke timezone yang diinginkan
	t = t.In(location)

	// format output string
	output := t.Format("2006-01-02 15:04:05.999 -0700")

	return output
}
