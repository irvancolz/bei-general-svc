package pkp

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	pkp "be-idx-tsg/internal/app/httprest/repository/pkp"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAllPKuser(c *gin.Context) (*helper.PaginationResponse, error)
	CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error)
	UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
}

type usecase struct {
	pkpRepo pkp.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		pkp.NewRepository(),
	}
}

func (uc *usecase) GetAllPKuser(c *gin.Context) (*helper.PaginationResponse, error) {
	results, errorResults := uc.pkpRepo.GetAllPKuser(c)
	if errorResults != nil {
		return nil, errorResults
	}

	var dataToConverted []interface{}
	for _, item := range results {
		dataToConverted = append(dataToConverted, item)
	}

	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"createdat", "updatedat", "questiondate", "answersat", "deletedat"})
	sortedData := helper.HandleDataSorting(c, filteredData)
	exportedFields := []string{
		"stakeholders",
		"code",
		"name",
		"questiondate",
		"question",
		"answersat",
		"answers",
		"topic",
		"answersby",
		"filename",
		"additionalinfo",
		"createdat",
		"createby",
	}
	columnHeaders := []string{
		"No",
		"Identitas Stakeholder",
		"Kode Perusahaan",
		"Nama Perusahaan",
		"Waktu Pertanyaan / Keluhan",
		"Pertanyaan / Keluhan",
		"Waktu Jawaban / Respon",
		"Jawaban / Respon",
		"Topik",
		"Personel Follow Up",
		"Lampiran",
		"Sumber Informasi Tambahan",
		"Waktu Buat",
		"Dibuat Oleh",
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	columnWidth := []float64{20, 50, 40, 40, 30, 40, 30, 40, 40, 50, 40, 30, 40, 40}
	tableHeaders := helper.GenerateTableHeaders(columnHeaders, columnWidth)

	var exportedData [][]string

	for i, item := range sortedData {
		var exportedRows []string
		exportedRows = append(exportedRows, strconv.Itoa(i+1))
		exportedRows = append(exportedRows, helper.MapToArray(item, exportedFields)...)
		for i, content := range exportedRows {
			if helper.IsContains([]int{4, 6, 12}, i) {
				unixTime, _ := strconv.Atoi(content)
				dateTime := time.Unix(int64(unixTime), 0)
				dateToFormat := helper.GetWIBLocalTime(&dateTime)
				exportedRows[i] = helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthShortNameInIndo) + " (" + helper.GetTimeAndMinuteOnly(dateToFormat) + ")"
			}
			if helper.IsContains([]int{10}, i) {
				attachmentPath := strings.Split(content, "_")
				if len(attachmentPath) > 2 {
					exportedRows[i] = strings.Join(attachmentPath[2:], "_")
				}
			}
		}

		exportedData = append(exportedData, exportedRows)
	}
	exportTableProps := helper.ExportTableToFileProps{
		Filename: "PKP",
		Data:     exportedData,
		Headers:  tablesColumns,
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Pertanyaan Keluhan Pelanggan"},
		},
		PdfConfig: &helper.PdfTableOptions{
			PapperWidth:  630,
			Papperheight: 300,
			HeaderRows:   tableHeaders,
		},
		ColumnWidth: []int{4, 25, 20, 25, 30, 25, 25, 25, 25, 20, 30, 15, 20, 20},
	}
	errorExport := helper.ExportTableToFile(c, exportTableProps)
	if errorExport != nil {
		return nil, errorExport
	}

	paginatedData := helper.HandleDataPagination(c, sortedData, filterParameter)
	followUpPersonel := utilities.GetFollowUpPersonel(c)
	paginatedData.FilterParameter["answersby"] = followUpPersonel
	return &paginatedData, nil
}

func (uc *usecase) CreatePKuser(pkp model.CreatePKuser, c *gin.Context) (int64, error) {
	data, err := uc.pkpRepo.CreatePKuser(pkp, c)
	if err != nil {
		return 0, err
	}

	notifMsg := "PKP Baru Telah Berhasil Dibuat"
	emailMsg := fmt.Sprintf("%s telah melakukan penambahan data pada Modul PKP", c.GetString("name_user"))
	notifType := "PKP"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return data, nil
}

func (uc *usecase) UpdatePKuser(pkp model.UpdatePKuser, c *gin.Context) (int64, error) {
	data, err := uc.pkpRepo.UpdatePKuser(pkp, c)
	if err != nil {
		return 0, err
	}

	notifMsg := "PKP Telah Berhasil Diubah"
	emailMsg := fmt.Sprintf("%s telah melakukan perubahan data pada Modul PKP", c.GetString("name_user"))
	notifType := "PKP"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)

	return data, nil
}

func (uc *usecase) Delete(id string, c *gin.Context) (int64, error) {
	data, err := uc.pkpRepo.Delete(id, c)
	if err != nil {
		return 0, err
	}

	notifMsg := "PKP Telah Berhasil Dihapus"
	emailMsg := fmt.Sprintf("%s telah melakukan penghapusan data pada Modul PKP", c.GetString("name_user"))
	notifType := "PKP"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)

	return data, nil
}
