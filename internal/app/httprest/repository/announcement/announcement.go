package announcement

import (
	// "be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"log"

	// "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAllAnnouncement() ([]*model.Announcement, error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	return &repository{
		DB: database.Init().MySql,
	}
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
		return nil, errors.New("List User Not Found")
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Announcement
		err := rows.Scan(
			&item.ID,
		)

		if err != nil {
			log.Println("[AQI-debug] [err] [repository] [Annoucement] [getQueryData] [GetAllAnnouncement] ", err)
			return nil, errors.New("Failed When Retrieving data")
		}
		result = append(result, &item)
	}

	return result, nil
}
