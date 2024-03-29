package contactperson

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/contact_person"
	"be-idx-tsg/internal/app/utilities"
	"be-idx-tsg/internal/pkg/email"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactPersonUsecaseInterface interface {
	SynchronizeInstitutionProfile(c *gin.Context, company_type string) ([]*model.InstitutionResponse, error)
	AddDivision(c *gin.Context, name string) (int64, error)
	EditDivision(c *gin.Context, props EditDivisionprops) (int64, error)
	GetAllDivision(c *gin.Context) ([]model.DivisionNameResponse, error)
	GetAllDivisionByCompany(company_id string) ([]*model.InstitutionDivisionResponse, error)
	AddMember(c *gin.Context, props AddMemberProps) (int64, error)
	EditMember(c *gin.Context, props EditMemberProps) (int64, error)
	GetMemberByDivision(company_type string, division_id []string) ([]model.InstitutionMembersResponse, error)
	GetMemberByDivisionAndCompanyID(division_id, company_id []string) ([]model.InstitutionMembersResponse, error)
	GetMemberByCompanyID(company_id []string) ([]model.InstitutionMembersResponse, error)
	GetMemberByCompanyType(company_type string) ([]model.InstitutionMembersResponse, error)
	GetMemberByID(id string) (*model.InstitutionMembersDetailResponse, error)
	GetAllCompanyByType(company_type string) ([]*model.InstitutionResponse, error)
	SearchCompany(keyword string) ([]*model.InstitutionResponse, error)
	DeleteMemberByID(c *gin.Context, member_id string) (int64, error)
	DeleteDivisionByID(c *gin.Context, division_id string) (int64, error)
	ExportMember(c *gin.Context, company_type, company_id, division_id string) error
	GetAllMembersEmail(c *gin.Context) (*helper.PaginationResponse, error)
}

type usecase struct {
	Repository repo.ContactPersonRepositoryInterface
}

func NewUsecase() ContactPersonUsecaseInterface {
	return &usecase{
		Repository: repo.NewRepository(),
	}
}

type AddMemberProps struct {
	Name           string `json:"name" binding:"required"`
	Institution_id string `json:"company_id" binding:"required"`
	Division_id    string `json:"division_id" binding:"required"`
	Position       string `json:"position" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
	Telephone      string `json:"telephone" binding:"required"`
}

type EditMemberProps struct {
	Id          string `json:"member_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Division_id string `json:"division_id" binding:"required"`
	Position    string `json:"position" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Telephone   string `json:"telephone" binding:"required"`
}

func (u *usecase) AddMember(c *gin.Context, props AddMemberProps) (int64, error) {
	user_id, _ := c.Get("user_id")

	addMemberArgs := repo.AddMemberProps{
		Name:           props.Name,
		Institution_id: props.Institution_id,
		Position:       props.Position,
		Email:          props.Email,
		Phone:          props.Phone,
		Telephone:      props.Telephone,
		Creator:        user_id.(string),
		Created_at:     time.Now(),
		Division_id:    props.Division_id,
	}
	result, error_result := u.Repository.AddMember(c, addMemberArgs)
	if error_result != nil {
		return 0, error_result
	}

	notifMsg := "Member Baru Berhasil Ditambahkan"
	emailMsg := fmt.Sprintf("%s telah menambahkan data member baru di contact person", c.GetString("name_user"))
	notifType := "contact person"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return result, nil
}

func (u *usecase) EditMember(c *gin.Context, props EditMemberProps) (int64, error) {
	user_id, _ := c.Get("user_id")

	editMemberArgs := repo.EditMemberProps{
		ID:          props.Id,
		Name:        props.Name,
		Position:    props.Position,
		Email:       props.Email,
		Phone:       props.Phone,
		Telephone:   props.Telephone,
		Division_id: props.Division_id,
		Created_at:  time.Now(),
		Creator:     user_id.(string),
	}
	result, error_result := u.Repository.EditMember(c, editMemberArgs)
	if error_result != nil {
		return 0, error_result
	}

	notifMsg := "Member Telah Berhasil DiPerbaharui"
	emailMsg := fmt.Sprintf("%s telah memperbaharui data member di contact person", c.GetString("name_user"))
	notifType := "contact person"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return result, nil
}

func (u *usecase) AddDivision(c *gin.Context, Name string) (int64, error) {
	user_id, _ := c.Get("user_id")

	createDivisionArgs := repo.AddDivisionprops{
		Name:       Name,
		IsDefault:  false,
		Creator:    user_id.(string),
		Created_at: time.Now(),
	}
	result, error_result := u.Repository.AddDivision(c, createDivisionArgs)
	if error_result != nil {
		return 0, error_result
	}

	notifMsg := "Fungsi Baru Telah Berhasil Dibuat"
	emailMsg := fmt.Sprintf("%s telah melakukan penambahan data pada Modul Fungsi", c.GetString("name_user"))
	notifType := "Fungsi"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return result, nil
}

type EditDivisionprops struct {
	Name        string `json:"name" binding:"required"`
	Division_id string `json:"division_id"`
}

func (u *usecase) EditDivision(c *gin.Context, props EditDivisionprops) (int64, error) {
	user_id, _ := c.Get("user_id")

	if props.Division_id == "" {
		log.Println("can not edit the default division, please specify the division id")
		return 0, errors.New("can not edit the default division, please specify the division id")
	}
	editDivArgs := repo.EditDivisionprops{
		Name:        props.Name,
		Division_id: props.Division_id,
		Creator:     user_id.(string),
		Created_at:  time.Now(),
	}

	notifMsg := "Fungsi Telah Berhasil Diubah"
	emailMsg := fmt.Sprintf("%s telah melakukan perubahan data pada Modul Fungsi", c.GetString("name_user"))
	notifType := "Fungsi"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)

	return u.Repository.EditDivision(c, editDivArgs)
}

func (u *usecase) SynchronizeInstitutionProfile(c *gin.Context, company_type string) ([]*model.InstitutionResponse, error) {

	if company_type == "" {
		log.Println("failed to sync contact person companies : please specify the company type you want to sync with")
		return nil, errors.New("failed to sync contact person companies : please specify the company type you want to sync with")
	}

	latestCompanyList, errorGetCompanies := u.Repository.GetAllCompanyWithExtType(company_type)
	if errorGetCompanies != nil {
		return nil, errorGetCompanies
	}

	errorSync := u.Repository.SynchronizeInstitutionProfile(latestCompanyList, company_type)
	if errorSync != nil {
		return nil, errorSync
	}

	latestCompaniesData, errorGetList := u.GetAllCompanyByType(company_type)
	if errorGetList != nil {
		return nil, errorGetList
	}

	notifMsg := "Sinkronisasi Contact Person Berhasil"
	emailMsg := fmt.Sprintf("%s telah melakukan sinkronisasi data company di contact person", c.GetString("name_user"))
	notifType := "contact person"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return latestCompaniesData, nil
}

func (u *usecase) GetAllDivision(c *gin.Context) ([]model.DivisionNameResponse, error) {
	results, errorResult := u.Repository.GetAllDivision()
	if errorResult != nil {
		return nil, errorResult
	}

	var exportedData []interface{}

	for _, item := range results {
		exportedData = append(exportedData, item)
	}

	tableColumn := []string{"No", "Nama"}
	columnWidth := []float64{20, 80}

	var tableTxtCol [][]string
	var tableTxtWidth []int

	tableTxtCol = append(tableTxtCol, tableColumn)
	for _, width := range columnWidth {
		tableTxtWidth = append(tableTxtWidth, int(width))
	}

	tableHeaders := helper.GenerateTableHeaders(tableColumn, columnWidth)

	var exportedDataStr [][]string

	exportedDataMap := helper.ConvertToMap(exportedData)

	for i, maps := range exportedDataMap {
		var exportedRows []string
		// add number
		exportedRows = append(exportedRows, fmt.Sprintf("%v", i+1))
		exportedRows = append(exportedRows, helper.MapToArray(maps, []string{"name"})...)

		// add rows to data
		exportedDataStr = append(exportedDataStr, exportedRows)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename:    "company division",
		ExcelConfig: &helper.ExportToExcelConfig{},
		PdfConfig: &helper.PdfTableOptions{
			HeaderRows: tableHeaders,
		},
		Data:        exportedDataStr,
		Headers:     tableTxtCol,
		ColumnWidth: tableTxtWidth,
	}
	errorExport := helper.ExportTableToFile(c, exportConfig)
	if errorExport != nil {
		return nil, errorExport
	}

	return results, nil
}

func (u *usecase) GetAllDivisionByCompany(company_id string) ([]*model.InstitutionDivisionResponse, error) {
	return u.Repository.GetAllDivisionByCompany(company_id)
}

func (u *usecase) GetMemberByDivision(company_type string, division_id []string) ([]model.InstitutionMembersResponse, error) {
	return u.Repository.GetMemberByDivision(division_id, company_type)
}

func (u *usecase) GetMemberByDivisionAndCompanyID(division_id, company_id []string) ([]model.InstitutionMembersResponse, error) {
	return u.Repository.GetMemberByDivisionAndCompanyId(division_id, company_id)
}

func (u *usecase) SearchCompany(keyword string) ([]*model.InstitutionResponse, error) {
	return u.Repository.GetAllCompany(keyword)
}

func (u *usecase) GetMemberByCompanyID(company_id []string) ([]model.InstitutionMembersResponse, error) {
	var results []model.InstitutionMembersResponse
	memberList, error_member := u.Repository.GetMemberByCompany()
	if error_member != nil {
		return nil, error_member
	}
	for _, data := range memberList {
		if helper.IsContains(company_id, data.Institute_id) {
			result := model.InstitutionMembersResponse{
				Id:           data.Id,
				Name:         data.Name,
				Email:        data.Email,
				Phone:        data.Phone,
				Telephone:    data.Telephone,
				Position:     data.Position,
				Division:     data.Division,
				Company_code: data.Company_code,
				Company_name: data.Company_name,
			}

			results = append(results, result)
		}
	}

	return results, nil
}

func (u *usecase) GetAllCompanyByType(company_type string) ([]*model.InstitutionResponse, error) {
	var results []*model.InstitutionResponse
	company_list, err_company_list := u.Repository.GetAllCompany("")
	if err_company_list != nil {
		return nil, err_company_list
	}
	for _, data := range company_list {
		if strings.EqualFold(data.Type, company_type) {
			results = append(results, data)
		}
	}

	return results, nil
}

func (u *usecase) GetMemberByID(id string) (*model.InstitutionMembersDetailResponse, error) {
	return u.Repository.GetMemberByID(id)
}

func (u *usecase) DeleteMemberByID(c *gin.Context, member_id string) (int64, error) {
	user_id, _ := c.Get("user_id")

	if !u.Repository.CheckMemberAvailability(member_id) {
		return 0, errors.New("cannot delete members, this member not found on database")
	}

	deleteMemberArgs := repo.DeleteDataProps{
		Id:         member_id,
		Deleted_by: user_id.(string),
		Deleted_at: time.Now(),
	}

	notifMsg := "Member Telah Berhasil Dihapus"
	emailMsg := fmt.Sprintf("%s telah menghapus data member dari contact person", c.GetString("name_user"))
	notifType := "contact person"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)

	return u.Repository.DeleteMemberByID(deleteMemberArgs)
}

func (u *usecase) DeleteDivisionByID(c *gin.Context, division_id string) (int64, error) {
	user_id, _ := c.Get("user_id")
	isNotDefaultDivision := u.Repository.CheckDivisionEditAvailability(division_id)
	if !isNotDefaultDivision {
		return 0, errors.New("divisi merupakan default sistem")
	}

	is_deletadble := u.Repository.CheckDivisionDeleteAvailability(division_id)
	if !is_deletadble {
		return 0, errors.New("divisi memiliki karyawan")
	}

	deleteDivArgs := repo.DeleteDataProps{
		Id:         division_id,
		Deleted_by: user_id.(string),
		Deleted_at: time.Now(),
	}

	notifMsg := "Fungsi Telah Berhasil Dihapus"
	emailMsg := fmt.Sprintf("%s telah melakukan penghapusan data pada Modul Fungsi", c.GetString("name_user"))
	notifType := "Fungsi"

	utilities.CreateNotifForAdminApp(c, notifType, notifMsg)
	email.SendEmailForUserAdminApp(c, notifMsg, emailMsg)
	utilities.CreateNotifForUserAng(c, notifType, notifMsg)
	email.SendEmailForUserAng(c, notifMsg, emailMsg)

	return u.Repository.DeleteDivisionByID(deleteDivArgs)
}

func (u *usecase) GetMemberByCompanyType(company_type string) ([]model.InstitutionMembersResponse, error) {
	var results []model.InstitutionMembersResponse
	memberList, error_member := u.Repository.GetMemberByCompany()
	if error_member != nil {
		return nil, error_member
	}
	for _, data := range memberList {
		if strings.EqualFold(data.Institute_type, company_type) {
			result := model.InstitutionMembersResponse{
				Id:           data.Id,
				Name:         data.Name,
				Email:        data.Email,
				Phone:        data.Phone,
				Telephone:    data.Telephone,
				Position:     data.Position,
				Division:     data.Division,
				Company_code: data.Company_code,
				Company_name: data.Company_name,
			}

			results = append(results, result)
		}
	}

	return results, nil
}

func (u *usecase) ExportMember(c *gin.Context, company_type, company_id, division_id string) error {
	var exportedField, tableHeaderName []string
	var columnWidths []float64

	exportTitle := "CONTACT PERSON "
	exportedField = []string{
		"division",
		"name",
		"position",
		"phone",
		"telephone",
		"email",
	}

	tableHeaderName = []string{
		"No",
		"Fungsi",
		"Nama",
		"Jabatan",
		"No. HP",
		"No Tel. Kantor",
		"Email"}

	columnWidths = []float64{20, 40, 60, 60, 50, 40, 60}

	tableHeaders := helper.GenerateTableHeaders(tableHeaderName, columnWidths)

	excelConfig := helper.ExportToExcelConfig{
		CollumnStart: "b",
	}
	pdfConfig := helper.PdfTableOptions{
		HeaderTitle:  "CONTACT PERSON ANGGOTA BURSA / PARTISIPAN / PJ SPPA / DU",
		HeaderRows:   tableHeaders,
		PapperWidth:  500,
		Papperheight: 300,
	}

	var dataToExported [][]string
	var memberStructList []model.InstitutionMembersResponse

	if len(company_type) <= 0 && len(company_id) <= 0 {
		return errors.New("failed to export data members: please specify which companies member to be exported")
	}

	if len(division_id) > 0 && len(company_id) > 0 {
		memberList, errorGetMember := u.GetMemberByDivisionAndCompanyID([]string{division_id}, []string{company_id})
		if errorGetMember != nil {
			return errorGetMember
		}
		memberStructList = append(memberStructList, memberList...)

		companyCode := func() string {
			if len(memberStructList) <= 0 {
				return ""
			}
			return memberStructList[0].Company_code
		}()

		excelConfig.HeaderText = []string{exportTitle + u.Repository.GetCompanyType(company_id), fmt.Sprintf("Kode :	%s", companyCode), fmt.Sprintf("Nama Perusahaan :	%s", companyCode)}
		pdfConfig.HeaderRows = tableHeaders
	} else if len(company_id) > 0 {
		memberList, errorGetMember := u.GetMemberByCompanyID([]string{company_id})
		if errorGetMember != nil {
			return errorGetMember
		}
		memberStructList = append(memberStructList, memberList...)

		companyCode := func() string {
			if len(memberStructList) <= 0 {
				return ""
			}
			return memberStructList[0].Company_code
		}()

		pdfConfig.HeaderRows = tableHeaders
		excelConfig.HeaderText = []string{exportTitle + u.Repository.GetCompanyType(company_id), fmt.Sprintf("Kode :	%s", companyCode), fmt.Sprintf("Nama Perusahaan :	%s", companyCode)}
	} else if len(company_type) >= 0 {
		exportedField = []string{
			"company_code",
			"company_name",
			"division",
			"name",
			"position",
			"phone",
			"telephone",
			"email",
		}

		tableHeaderName = []string{
			"No",
			"Kode",
			"Nama Perusahaan",
			"Fungsi",
			"Nama",
			"Jabatan",
			"No. HP",
			"No Tel. Kantor",
			"Email"}

		columnWidths = []float64{20, 40, 50, 40, 60, 60, 50, 40, 60}

		tableHeaders := helper.GenerateTableHeaders(tableHeaderName, columnWidths)

		pdfConfig.HeaderRows = tableHeaders
		excelConfig.HeaderText = []string{exportTitle + " " + company_type}
		memberList, errorGetMember := u.GetMemberByCompanyType(company_type)
		if errorGetMember != nil {
			return errorGetMember
		}
		memberStructList = append(memberStructList, memberList...)
	}

	for i, member := range memberStructList {
		var memberData []string
		memberData = append(memberData, strconv.Itoa(i+1))
		memberData = append(memberData, helper.StructToArray(member, exportedField)...)

		dataToExported = append(dataToExported, memberData)
	}

	var tablesColumns [][]string
	var ColumnWidtINT []int

	for _, width := range columnWidths {
		ColumnWidtINT = append(ColumnWidtINT, int(width))
	}

	tablesColumns = append(tablesColumns, tableHeaderName)

	errorCreateFile := helper.ExportTableToFile(c, helper.ExportTableToFileProps{
		Filename:    "contact_person_members",
		Data:        dataToExported,
		Headers:     tablesColumns,
		ExcelConfig: &excelConfig,
		PdfConfig:   &pdfConfig,
		ColumnWidth: ColumnWidtINT,
	})
	if errorCreateFile != nil {
		return errorCreateFile
	}

	return nil
}

func (m *usecase) GetAllMembersEmail(c *gin.Context) (*helper.PaginationResponse, error) {
	membersEmail, errGetEmail := m.Repository.GetAllMembersEmail(c)
	if errGetEmail != nil {
		return nil, errGetEmail
	}
	dataToFilter := []interface{}{}

	for _, item := range membersEmail {
		dataToFilter = append(dataToFilter, item)
	}

	filteredData, filterParams := helper.HandleDataFiltering(c, dataToFilter, []string{})
	column := []string{
		"No",
		"Tipe Perusahaan",
		"Kode Perusahaan",
		"Nama Perusahaan",
		"Fungsi",
		"Nama",
		"Email",
	}
	colWidth := []float64{20, 40, 40, 70, 50, 60, 60}

	pdfConfig := helper.PdfTableOptions{
		HeaderTitle:  "Database Email",
		HeaderRows:   helper.GenerateTableHeaders(column, colWidth),
		PapperWidth:  400,
		Papperheight: 300,
	}

	columnOrder := []string{
		"companyName",
		"companyCode",
		"companyName",
		"division",
		"name",
		"email",
	}
	dataToExported := [][]string{}

	for i, item := range filteredData {
		var data []string
		data = append(data, strconv.Itoa(i+1))
		dataInArr := helper.MapToArray(item, columnOrder)
		data = append(data, dataInArr...)
		dataToExported = append(dataToExported, data)
	}

	exportConfig := helper.ExportTableToFileProps{
		Filename:    "contact_person_member_email",
		PdfConfig:   &pdfConfig,
		Data:        dataToExported,
		Headers:     [][]string{column},
		ColumnWidth: []int{5, 10, 10, 30, 15, 20, 20},
		ExcelConfig: &helper.ExportToExcelConfig{
			HeaderText: []string{"Database Email"},
		},
	}
	errExport := helper.ExportTableToFile(c, exportConfig)
	if errExport != nil {
		return nil, errExport
	}

	paginatedRes := helper.HandleDataPagination(c, filteredData, filterParams)
	return &paginatedRes, nil
}
