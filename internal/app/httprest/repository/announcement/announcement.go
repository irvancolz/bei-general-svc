package announcement

import (
	// "be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"
	"time"

	// "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByIDCode(id string) (*model.AnnouncementCode, error)
	GetAllAnnouncement() ([]*model.Announcement, error)
	Create(ab model.CreateAnnouncement) (int64, error)
	Update(ab model.UpdateAnnouncement) (int64, error)
	Delete(id string, deleted_by string) (int64, error)
	GetByCode(id string) ([]model.Announcement, error)
	GetByIDandType(id string, types string) (*model.Announcement, error)
	GetAllMin() (*[]model.GetAllAnnouncement, error)
	GetAllANWithFilter(keyword []string) ([]*model.Announcement, error)
	GetAllANWithSearch(keyword string) ([]*model.Announcement, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (m *repository) GetAllANWithSearch(keyword string) ([]*model.Announcement, error) {
	var querySelect = ` 
	select z.* from (select an.id, an.code, an.type, an.status, an.license, from announcement an where an.is_deleted = false)  as z
	where z.code ilike '%` + keyword + `%' 
	or z.name ilike '%` + keyword + `%'`
	var listData = []*model.Announcement{}
	// log.Println(querySelect)
	selDB, err := m.DB.Query(querySelect)
	if err != nil {
		return nil, errors.New("not found")
	}
	for selDB.Next() {
		an := model.Announcement{}
		err = selDB.Scan(
			&an.ID,
			&an.Code,
			&an.Type,
			&an.Status,
			&an.License,
		)
		if err != nil {
			return nil, errors.New("not found")
		}
		listData = append(listData, &an)
	}
	return listData, nil

}

func (m *repository) GetAllANWithFilter(keyword []string) ([]*model.Announcement, error) {

	var querySelect = " select an.id, an.code, an.type, an.status, an.license, an.operational_status  from announcement an where an.is_deleted = false "

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
		// log.Println(query)
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
	// log.Println(query)
	selDB, err := m.DB.Query(query)
	if err != nil {
		return nil, errors.New("not found")
	}
	for selDB.Next() {
		an := model.Announcement{}
		err = selDB.Scan(
			&an.ID,
			&an.Code,
			&an.Type,
			&an.Status,
			&an.License,
			&an.OperationalStatus,
		)
		if err != nil {
			return nil, errors.New("not found")
		}
		listData = append(listData, &an)
	}
	return listData, nil
}

func (m *repository) GetAllAnnouncement() ([]*model.Announcement, error) {
	result := []*model.Announcement{}
	query := `
	SELECT 
	id
   FROM announcements;`
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println("[AQI-debug] [err] [repository] [Annoucement] [sqlQuery] [GetAllAnnouncement] ", err)
		return nil, errors.New("list announcement not found")
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Announcement
		err := rows.Scan(
			&item.ID,
		)

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [Annoucement] [getQueryData] [GetAllAnnouncement] ", err)
			return nil, errors.New("failed when retrieving data")
		}
		result = append(result, &item)
	}

	return result, nil
}

func (m *repository) GetByIDCode(id string) (*model.AnnouncementCode, error) {
	query := `
		SELECT 
			code,
			type
		FROM announcement
		WHERE id = $1 AND deleted_by IS NULL`
	data := &model.AnnouncementCode{}

	if err := m.DB.QueryRow(query, id).Scan(
		&data.Code,
		&data.Type,
	); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *repository) GetByCode(code string) ([]model.Announcement, error) {
	query := `
		SELECT 
			id,
			code,
			type,
			created_at,
			created_by,
			updated_at,
			updated_by,
			is_deleted
		FROM anggota_bursa
		WHERE code  LIKE $1 AND is_deleted = false`
	listData := []model.Announcement{}
	selDB, err := m.DB.Query(query, "%"+code+"%")
	if err != nil {
		return nil, err
	}
	for selDB.Next() {
		an := model.Announcement{}
		err = selDB.Scan(
			&an.ID,
			&an.Code,
			&an.Type,
			&an.Status,
		)
		if err != nil {
			return nil, err
		}
		listData = append(listData, an)
	}
	return listData, nil
}

func (m *repository) GetAllMin() (*[]model.GetAllAnnouncement, error) {
	query := `
		SELECT 
			id,
			code,
			type
		FROM announcement
		WHERE is_deleted = false
		ORDER BY id ASC`
	// log.Println(id)
	listData := []model.GetAllAnnouncement{}
	selDB, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for selDB.Next() {
		an := model.GetAllAnnouncement{}
		err = selDB.Scan(
			&an.ID,
			&an.Code,
			&an.Type,
		)
		if err != nil {
			return nil, err
		}
		listData = append(listData, an)
	}
	return &listData, nil
}

func (m *repository) Create(an model.CreateAnnouncement) (int64, error) {
	query := `
	INSERT INTO 
	announcement (
		code,
		type,
		role_id,
		effective_date,
		regarding
		status,
		status,
		created_at, 
		created_by,
		is_deleted
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 ,$10, $11, $12, $13, false);`
	// log.Println(id)
	// data := &model.Announcement{}
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	selDB, err := m.DB.Exec(
		query,
		an.Code,
		an.Type,
		an.RoleId,
		an.EffectiveDate,
		an.Regarding,
		an.Status,
		an.FileURL,
		createdAt,
		an.CreatedBy)
	if err != nil {
		return 0, err
	}

	// selDB.LastInsertId()
	LastInsertId, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return LastInsertId, nil
}

func (m *repository) Update(an model.UpdateAnnouncement) (int64, error) {
	query := `
	UPDATE
	anggota_bursa SET
		code = $2,
		type = $3,
		role_id = $4,
		effective_date = $5,
		regarding = $6,
		status = $7,
		file_url = $8,
		updated_at = $13, 
		updated_by = $14
	WHERE id = $1 AND is_deleted = false;`
	// log.Println(id)
	// data := &model.Announcement{}
	updated_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	selDB, err := m.DB.Exec(
		query,
		an.Id,
		an.Code,
		an.Type,
		an.RoleId,
		an.EffectiveDate,
		an.Regarding,
		an.Status,
		an.FileURL,
		updated_at,
		an.UpdatedBy,
	)
	if err != nil {
		return 0, err
	}

	// selDB.LastInsertId()
	RowsAffected, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return RowsAffected, nil
}

func (m *repository) Delete(id string, deleted_by string) (int64, error) {
	query := `
	UPDATE
	announcement SET
		is_deleted = true,
		deleted_by = $2,
		deleted_at = $3,
		updated_by = $4, 
		updated_at = $5
	WHERE 
		id = $1 
		AND
		is_deleted = false;`
	deleted_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	selDB, err := m.DB.Exec(query, id, deleted_by, deleted_at, deleted_by, deleted_at)
	if err != nil {
		return 0, err
	}

	// selDB.LastInsertId()
	RowsAffected, err := selDB.RowsAffected()
	if err != nil {
		return 0, err
	}

	return RowsAffected, nil
}

func (m *repository) GetByIDandType(id string, types string) (*model.Announcement, error) {
	query := `
		SELECT 
			id,
			code,
			type,
			role_id,
			effective_date,
			regarding,
			status,
			file_url,
		FROM announcement
		WHERE id = $1 AND type = $2 AND is_deleted = false`
	log.Println(id)
	data := &model.Announcement{}

	if err := m.DB.QueryRow(query, id, types).Scan(
		&data.ID,
		&data.Code,
		&data.Type,
		&data.RoleId,
		&data.EffectiveDate,
		&data.Regarding,
		&data.Status,
	); err != nil {
		return nil, err
	}

	return data, nil
}
