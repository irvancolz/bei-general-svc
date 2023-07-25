package log_system

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/app/httprest/repository/log_system"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAll() ([]*model.LogSystem, error)
	CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error)
	ExportLogSystem(c *gin.Context) error
}

type usecase struct {
	logSystemRepo log_system.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		log_system.NewRepository(),
	}
}

func (m *usecase) GetAll() ([]*model.LogSystem, error) {
	return m.logSystemRepo.GetAll()
}

func (m *usecase) CreateLogSystem(log model.CreateLogSystem, c *gin.Context) (int64, error) {
	return m.logSystemRepo.CreateLogSystem(log, c)
}

func (m *usecase) ExportLogSystem(c *gin.Context) error {
	exportedField := []string{
		"modul",
		"sub",
		"action",
		"detail",
		"user",
		"ip",
		"date",
	}

	tableHeader := []string{
		"Modul",
		"Sub Modul",
		"Aksi",
		"Detail",
		"User",
		"IP",
		"Tanggal",
	}

	var dataToExported [][]string
	var logSystemList []*model.LogSystem
	dataToExported = append(dataToExported, tableHeader)

	logSystemList, _ = m.GetAll()

	for _, data := range logSystemList {
		log := model.LogSystemExport{
			Modul:  data.Modul,
			Sub:    data.SubModul,
			Action: data.Action,
			Detail: data.Detail,
			User:   data.UserName,
			IP:     data.IP,
			Date:   data.CreatedAt.Format("2 Jan 2006 - 15:04"),
		}

		fmt.Println(data.UserName)

		var logData []string
		logData = append(logData, helper.StructToArray(log, exportedField)...)

		dataToExported = append(dataToExported, logData)
	}

	excelConfig := helper.ExportToExcelConfig{
		CollumnStart: "b",
	}
	pdfConfig := helper.PdfTableOptions{
		HeaderTitle: "Log System",
	}
	errorCreateFile := helper.ExportTableToFile(c, helper.ExportTableToFileProps{
		Filename:    "log_system",
		Data:        dataToExported,
		ExcelConfig: &excelConfig,
		PdfConfig:   &pdfConfig,
	})
	if errorCreateFile != nil {
		return errorCreateFile
	}

	return nil
}
