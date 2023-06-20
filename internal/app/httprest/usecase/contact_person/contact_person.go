package contactperson

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	repo "be-idx-tsg/internal/app/httprest/repository/contact_person"
	"be-idx-tsg/internal/app/utilities"
	"errors"
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
	GetAllDivision() ([]model.DivisionNameResponse, error)
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
	return u.Repository.EditDivision(c, editDivArgs)
}

func (u *usecase) SynchronizeInstitutionProfile(c *gin.Context, company_type string) ([]*model.InstitutionResponse, error) {

	var latestCompanyList []model.ContactPersonSyncCompaniesResource
	var errorGetCompanies error

	if company_type == "" {
		log.Println("failed to sync contact person companies : please specify the company type you want to sync with")
		return nil, errors.New("failed to sync contact person companies : please specify the company type you want to sync with")
	}

	if strings.EqualFold(company_type, "AB") {
		latestCompanyList, errorGetCompanies = utilities.GetLatestABCompanies(c)
	}

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

	return latestCompaniesData, nil
}

func (u *usecase) GetAllDivision() ([]model.DivisionNameResponse, error) {
	return u.Repository.GetAllDivision()
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
		if data.Type == company_type {
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
	return u.Repository.DeleteDivisionByID(deleteDivArgs)
}

func (u *usecase) GetMemberByCompanyType(company_type string) ([]model.InstitutionMembersResponse, error) {
	var results []model.InstitutionMembersResponse
	memberList, error_member := u.Repository.GetMemberByCompany()
	if error_member != nil {
		return nil, error_member
	}
	for _, data := range memberList {
		if data.Institute_type == company_type {
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
	exportedField := []string{
		"company_code",
		"company_name",
		"division",
		"name",
		"email",
		"phone",
		"telephone",
		"position"}

	tableHeader := []string{
		"No",
		"Kode Perusahaan",
		"Nama Perusahaan",
		"Divisi",
		"Nama",
		"Email",
		"Nomor HP",
		"No Telepon Kantor",
		"Posisi"}

	var dataToExported [][]string
	var memberStructList []model.InstitutionMembersResponse
	dataToExported = append(dataToExported, tableHeader)

	if len(company_type) <= 0 && len(company_id) <= 0 {
		return errors.New("failed to export data members: please specify which companies member to be exported")
	}

	if len(division_id) > 0 && len(company_id) > 0 {
		memberList, errorGetMember := u.GetMemberByDivisionAndCompanyID([]string{division_id}, []string{company_id})
		if errorGetMember != nil {
			return errorGetMember
		}
		memberStructList = append(memberStructList, memberList...)

	} else if len(company_id) > 0 {
		memberList, errorGetMember := u.GetMemberByCompanyID([]string{company_id})
		if errorGetMember != nil {
			return errorGetMember
		}
		memberStructList = append(memberStructList, memberList...)

	} else if len(company_type) >= 0 {
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

	excelConfig := helper.ExportToExcelConfig{
		CollumnStart: "b",
	}
	pdfConfig := helper.PdfTableOptions{
		HeaderTitle: "Contact Person Member",
	}
	errorCreateFile := helper.ExportTableToFile(c, helper.ExportTableToFileProps{
		Filename:    "contact_person_members",
		Data:        dataToExported,
		ExcelConfig: &excelConfig,
		PdfConfig:   &pdfConfig,
	})
	if errorCreateFile != nil {
		return errorCreateFile
	}

	return nil
}
