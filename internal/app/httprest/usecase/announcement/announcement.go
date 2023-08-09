package announcement

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	an "be-idx-tsg/internal/app/httprest/repository/announcement"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetAllAnnouncement(c *gin.Context) (*helper.PaginationResponse, error)
	Detail(id string, c *gin.Context) (*model.Announcement, error)
	Create(ab model.CreateAnnouncement, c *gin.Context) (int64, error)
	Update(ab model.UpdateAnnouncement, c *gin.Context) (int64, error)
	Delete(id string, c *gin.Context) (int64, error)
	Export(c *gin.Context, id string) error
	GetAllANWithFilter(keyword []string) ([]*model.Announcement, error)
	GetAllANWithSearch(keyword string, InformationType string, startDate string, endDate string) ([]*model.Announcement, error)
}

type usecase struct {
	anRepo an.Repository
}

func DetailUseCase() Usecase {
	return &usecase{
		an.NewRepository(),
	}
}
func (m *usecase) Detail(id string, c *gin.Context) (*model.Announcement, error) {
	return m.anRepo.GetByID(id, c)
}
func (m *usecase) GetAllAnnouncement(c *gin.Context) (*helper.PaginationResponse, error) {
	dataStruct, errorData := m.anRepo.GetAllAnnouncement(c)
	if errorData != nil {
		return nil, errorData
	}
	var dataToConverted []interface{}
	for _, item := range dataStruct {
		dataToConverted = append(dataToConverted, item)
	}
	filteredData, filterParameter := helper.HandleDataFiltering(c, dataToConverted, []string{"effective_date"})

	columnHeaders := []string{"Jenis Information", "Perihal", "Waktu", "Tanggal Efektif"}
	columnWidth := []float64{30, 60, 30, 50}
	var columnWidthInt []int

	for _, width := range columnWidth {
		columnWidthInt = append(columnWidthInt, int(width))
	}

	var tablesColumns [][]string
	tablesColumns = append(tablesColumns, columnHeaders)

	dataOrder := []string{"information_type", "regarding", "effective_date"}
	var exportedData [][]string

	for _, content := range filteredData {
		var item []string
		item = append(item, helper.MapToArray(content, dataOrder)...)

		for i, content := range item {
			if i == 2 {
				unixTime, _ := strconv.Atoi(content)
				dateToFormat := time.Unix(int64(unixTime), 0)
				item[i] = helper.GetTimeAndMinuteOnly(dateToFormat)
				item = append(item, helper.ConvertTimeToHumanDateOnly(dateToFormat, helper.MonthFullNameInIndo))
			}
		}

		exportedData = append(exportedData, item)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename:    "Announcements",
		ExcelConfig: &helper.ExportToExcelConfig{},
		PdfConfig: &helper.PdfTableOptions{
			HeaderRows: helper.GenerateTableHeaders(columnHeaders, columnWidth),
		},
		Data:        exportedData,
		Headers:     tablesColumns,
		ColumnWidth: columnWidthInt,
	}
	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return nil, errorExport
	}

	paginatedData := helper.HandleDataPagination(c, filteredData, filterParameter)
	return &paginatedData, nil
}
func (m *usecase) Create(an model.CreateAnnouncement, c *gin.Context) (int64, error) {
	return m.anRepo.Create(an, c)
}
func (m *usecase) Update(an model.UpdateAnnouncement, c *gin.Context) (int64, error) {
	return m.anRepo.Update(an, c)
}
func (m *usecase) Delete(id string, c *gin.Context) (int64, error) {
	return m.anRepo.Delete(id, c)
}
func (m *usecase) GetAllANWithFilter(keyword []string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithFilter(keyword)
}
func (m *usecase) GetAllANWithSearch(keyword string, InformationType string, startDate string, endDate string) ([]*model.Announcement, error) {
	return m.anRepo.GetAllANWithSearch(InformationType, keyword, startDate, endDate)
}

func (m *usecase) Export(c *gin.Context, id string) error {
	var exportedData *model.Announcement

	exportedData, errorGetData := m.Detail(id, c)
	if errorGetData != nil {
		return errorGetData
	}

	exportConfig := helper.ExportAnnouncementsToFileProps{
		Filename: "Announcements",
		Data:     *exportedData,
		PdfConfig: helper.PdfTableOptions{
			HeaderTitle: "Pengumuman",
		},
	}
	errExport := helper.ExportAnnouncementsToFile(c, exportConfig)
	if errExport != nil {
		return errExport
	}

	return nil
}
