package contactperson

import (
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ContactPersonRepositoryInterface interface {
	SynchronizeInstitutionProfile(data []model.ContactPersonSyncCompaniesResource, company_type string) error
	GetAllCompany(keyword string) ([]*model.InstitutionResponse, error)
	GetCompanyDetail(id string) (*model.InstitutionResponse, error)
	GetAllDivision() ([]model.DivisionNameResponse, error)
	GetAllDivisionByCompany(company_id string) ([]*model.InstitutionDivisionResponse, error)
	AddDivision(c *gin.Context, props AddDivisionprops) (int64, error)
	EditDivision(c *gin.Context, props EditDivisionprops) (int64, error)
	AddMember(c *gin.Context, props AddMemberProps) (int64, error)
	EditMember(c *gin.Context, props EditMemberProps) (int64, error)
	GetMemberByDivision(division_id []string, company_type string) ([]model.InstitutionMembersResponse, error)
	GetMemberByDivisionAndCompanyId(division_id, company_id []string) ([]model.InstitutionMembersResponse, error)
	GetMemberByCompany() ([]*model.InstitutionMembersByCompanyResponse, error)
	GetMemberByID(member_id string) (*model.InstitutionMembersDetailResponse, error)
	DeleteMemberByID(props DeleteDataProps) (int64, error)
	DeleteDivisionByID(props DeleteDataProps) (int64, error)
	CheckMemberAvailability(id string) bool
	CheckDivisionEditAvailability(id string) bool
	CheckDivisionDeleteAvailability(division_id string) bool
	CheckDivisionViewAvailability(division_id string) bool
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository() ContactPersonRepositoryInterface {
	return &repository{
		DB: database.Init().MySql,
	}
}

func (r *repository) GetAllDivision() ([]model.DivisionNameResponse, error) {
	var results []model.DivisionNameResponse
	row_result, error_row := r.DB.Queryx(getAllDivisionQuerry)
	if error_row != nil {
		log.Println("failed to get division data from databases : ", error_row)
		return nil, error_row
	}
	defer row_result.Close()
	for row_result.Next() {
		var result model.DivisionNameResponse

		err_scan := row_result.StructScan(
			&result,
		)
		result.Is_Default = !r.CheckDivisionEditAvailability(result.Id)

		if err_scan != nil {
			log.Println("failed to copy data from database : ", err_scan)
			return nil, err_scan
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *repository) CheckMemberAvailability(id string) bool {
	result := r.DB.QueryRowx(checkMemberViewAvailabilityQuerry, id)
	var total int64
	error_result := result.Scan(&total)
	if error_result != nil {
		log.Println("division not found : ", error_result)
		return false
	}

	return total > 0
}

func (r *repository) CheckDivisionDeleteAvailability(id string) bool {
	result := r.DB.QueryRowx(checkDivisionDeleteAvailabilityQuerry, id)

	var total int64
	error_result := result.Scan(&total)
	if error_result != nil {
		log.Println("division not found : ", error_result)
		return false
	}

	return total <= 0
}

func (r *repository) CheckDivisionEditAvailability(division_id string) bool {
	var result int64
	row_result := r.DB.QueryRowx(checkDivisionEditAvailabilityQuerry, division_id)

	err_scan := row_result.Scan(&result)
	if err_scan != nil {
		log.Println("failed to get division delete Availability : ", err_scan)
		return false
	}

	return result > 0
}

func (r *repository) CheckDivisionViewAvailability(division_id string) bool {
	var result int64
	row_result := r.DB.QueryRowx(checkDivisionViewAvailabilityQuerry, division_id)

	err_scan := row_result.Scan(&result)
	if err_scan != nil {
		log.Println("failed to get division delete Availability : ", err_scan)
		return false
	}

	return result > 0
}

func (r *repository) SynchronizeInstitutionProfile(data []model.ContactPersonSyncCompaniesResource, types string) error {
	statement, errStatement := r.DB.Preparex(syncContactPersonCompaniesQuery)
	if errStatement != nil {
		log.Println("failed to prepara statement on sync contact person :", errStatement)
		return errStatement
	}

	for _, companies := range data {
		execResult, errorExec := statement.Exec(
			companies.Code,
			companies.Name,
			companies.Status,
			types,
			companies.Is_deleted,
		)
		if errorExec != nil {
			log.Println("failed to add latest data to contact person :", errorExec)
			return errorExec
		}
		_, errorResult := execResult.RowsAffected()
		if errorResult != nil {
			log.Println("failed to get edited companies in contact person :", errorResult)
			return errorResult
		}
	}

	return nil
}

func (r *repository) GetAllDivisionByCompany(company_id string) ([]*model.InstitutionDivisionResponse, error) {
	var results []*model.InstitutionDivisionResponse
	row_result, error_row := r.DB.Queryx(getAllDivisionByCompanyQuerry)
	if error_row != nil {
		log.Println("failed to get division data from databases : ", error_row)
		return nil, error_row
	}
	defer row_result.Close()
	for row_result.Next() {
		var result_mock model.InstitutionDivisionResultSetResponse

		err_scan := row_result.StructScan(&result_mock)
		if err_scan != nil {
			log.Println("failed to copy data from database : ", err_scan)
			return nil, err_scan
		}

		result := model.InstitutionDivisionResponse{
			Id:         result_mock.Id,
			Name:       result_mock.Name,
			Company_id: company_id,

			// handling null response from database
			Created_at: result_mock.Created_at.Time,
			Created_by: result_mock.Created_by.String,
			Updated_at: result_mock.Updated_at.Time,
			Updated_by: result_mock.Updated_by.String,
		}

		company_detail, error_company := r.GetCompanyDetail(company_id)
		if error_company != nil {
			return nil, error_company
		}

		result.Company_name = company_detail.Name
		result.Company_code = company_detail.Code

		division_members, error_get_members := r.GetMemberByDivisionAndCompanyId([]string{result.Id}, []string{company_id})
		if error_get_members != nil {
			return nil, error_get_members
		}
		result.Members = division_members
		if len(division_members) == 0 {
			result.Members = make([]model.InstitutionMembersResponse, 0)
		}

		results = append(results, &result)
	}
	sort.SliceStable(results, func(current, before int) bool {
		return len(results[current].Members) > len(results[before].Members)
	})
	return results, nil
}

type AddDivisionprops struct {
	Name       string
	IsDefault  bool
	Creator    string
	Created_at time.Time
}

func (r *repository) AddDivision(c *gin.Context, props AddDivisionprops) (int64, error) {
	rows_result, error_sync := r.DB.Exec(addDivisionQuerry,
		props.IsDefault,
		props.Name,
		props.Created_at,
		props.Creator)
	if error_sync != nil {
		log.Println("failed to add division : ", error_sync)
		return 0, error_sync
	}

	result, error_result := rows_result.RowsAffected()
	if error_result != nil {
		log.Println("failed to get added division total : ", error_result)
		return 0, error_result
	}

	if result == 0 {
		log.Println("the division is failed to created, please check your querry and try again")
		return 0, errors.New("the division is failed to created, please check your querry and try again")
	}

	return result, nil
}

type EditDivisionprops struct {
	Name        string
	Division_id string
	Creator     string
	Created_at  time.Time
}

func (r *repository) EditDivision(c *gin.Context, props EditDivisionprops) (int64, error) {

	isEditable := r.CheckDivisionEditAvailability(props.Division_id)
	if !isEditable {
		return 0, errors.New("cannot edit the default division, please select the non default division")
	}

	edit_result, error_edit := r.DB.Exec(editDivisionQuerry, props.Division_id, props.Name, props.Created_at, props.Creator)
	if error_edit != nil {
		log.Println("failed to edit divison in database : ", error_edit)
		return 0, error_edit
	}

	result, error_rsult := edit_result.RowsAffected()
	if error_rsult != nil {
		log.Println("failed to get division updated detail : ", error_rsult)
		return 0, error_rsult
	}

	if result == 0 {
		log.Println("failed to edit division, please check your id and try again")
		return 0, errors.New("failed to edit division, please check your id and try again")
	}
	return result, nil
}

type AddMemberProps struct {
	Name           string
	Institution_id string
	Division_id    string
	Position       string
	Email          string
	Phone          string
	Telephone      string
	Created_at     time.Time
	Creator        string
}

type EditMemberProps struct {
	ID          string
	Name        string
	Division_id string
	Position    string
	Email       string
	Phone       string
	Telephone   string
	Created_at  time.Time
	Creator     string
}

func (r *repository) AddMember(c *gin.Context, props AddMemberProps) (int64, error) {
	rows_result, error_row_result := r.DB.Exec(crateMemberQuerry, props.Institution_id, props.Division_id, props.Name, props.Phone, props.Telephone, props.Email, props.Position, props.Created_at, props.Creator)
	if error_row_result != nil {
		log.Println("failed to add new member to database : ", error_row_result)
		return 0, error_row_result
	}

	result, error_result := rows_result.RowsAffected()
	if error_result != nil {
		log.Println("failed to get added member : ", error_result)
		return 0, error_result
	}

	if result == 0 {
		log.Println("the member failed to be added to database, please try again")
		return 0, errors.New("the member failed to be added to database, please try again")
	}

	return result, nil
}

func (r *repository) EditMember(c *gin.Context, props EditMemberProps) (int64, error) {
	edit_result, error_edit := r.DB.Exec(editMemberQuerry,
		props.ID,
		props.Name,
		props.Position,
		props.Email,
		props.Phone,
		props.Telephone,
		props.Created_at,
		props.Creator,
		props.Division_id)
	if error_edit != nil {
		log.Println("failed to edit members data in database : ", error_edit)
		return 0, error_edit
	}

	result, error_result := edit_result.RowsAffected()
	if error_result != nil {
		log.Println("failed to get edited member total : ", error_result)
		return 0, error_result
	}

	if result == 0 {
		log.Println("the member is failed to edited, please check your querry / id and try again")
		return 0, errors.New("the member is failed to edited, please check your querry / id and try again")
	}
	return result, nil
}

type GetMemberByDivisionProps struct {
	Company_id  string
	Division_id string
}

func SearchMemberByDivisionAndCompanyIdQuerry(divisons, companies []string) string {

	var sb strings.Builder
	sb.WriteString(getMemberDetailByDivisionAndCompanyCodeBaseQuery)

	var formatedDivisions, formattedCompanies []string

	for _, id := range divisons {
		formatedDivisions = append(formatedDivisions, "'"+id+"'")
	}
	for _, company := range companies {
		formattedCompanies = append(formattedCompanies, "'"+company+"'")
	}
	sb.WriteString(fmt.Sprintf("\nAND d.deleted_at IS NULL \nAND d.deleted_by IS NULL \nAND d.id IN (%s) \nAND m.institution_id IN (%s) ORDER BY m.name ASC", strings.Join(formatedDivisions, ","), strings.Join(formattedCompanies, ",")))

	return sb.String()
}

func SearchMemberByDivisionAndCompanyIdProps(division_id, company_id []string) []GetMemberByDivisionProps {
	var results []GetMemberByDivisionProps

	divisionMemberGap := len(company_id) - len(division_id)
	if len(division_id) < len(company_id) {
		for i := 0; i < divisionMemberGap; i++ {
			division_id = append(division_id, "")
		}
	}

	for i := 0; i < len(company_id); i++ {
		result := GetMemberByDivisionProps{
			Company_id:  company_id[i],
			Division_id: division_id[i],
		}

		results = append(results, result)
	}

	return results
}

func (r *repository) GetMemberByDivisionAndCompanyId(division_id, company_id []string) ([]model.InstitutionMembersResponse, error) {
	var results []model.InstitutionMembersResponse

	if len(division_id) <= 0 {
		return nil, errors.New("cannot get the division list : no division_id is given, if you want to search the whole company member use get-all-company-member instead")
	}

	if len(company_id) <= 0 {
		return nil, errors.New("please specify the company you want to search")
	}

	if len(company_id) == 1 && !r.CheckDivisionViewAvailability(division_id[0]) {
		return nil, errors.New("cannot find the specified division")
	}

	querry := SearchMemberByDivisionAndCompanyIdQuerry(division_id, company_id)
	row_result, error_row := r.DB.Queryx(querry)
	if error_row != nil {
		log.Print("failed to get member by the division given from database: ", error_row)
		return nil, error_row
	}
	defer row_result.Close()

	for row_result.Next() {
		var result model.InstitutionMembersResponse
		err_scan := row_result.StructScan(&result)
		if err_scan != nil {
			log.Println("failed to store data from database :  ", err_scan)
			return nil, err_scan
		}

		results = append(results, result)
	}
	return results, nil
}

func generateGetMemberDivisionQuerry(division_id []string, company_type string) string {
	var sb strings.Builder
	var formated_ids []string
	for _, id := range division_id {
		formated_ids = append(formated_ids, "'"+id+"'")
	}
	sb.WriteString(fmt.Sprintf(getMemberDetailByDivisionBaseQuery, strings.Join(formated_ids, ",")))

	if company_type != "" {
		sb.WriteString(fmt.Sprintf(" \nAND i.Type = '%s'", company_type))
	}

	return sb.String()
}

func (r *repository) GetMemberByDivision(division_id []string, company_type string) ([]model.InstitutionMembersResponse, error) {
	var results []model.InstitutionMembersResponse

	if len(division_id) <= 0 {
		return nil, errors.New("cannot get the division list : no division_id is given, if you want to search the whole company member use get-all-company-member instead")
	}

	querry := generateGetMemberDivisionQuerry(division_id, company_type)
	row_result, error_row := r.DB.Queryx(querry)
	if error_row != nil {
		log.Print("failed to get member by the division given from database: ", error_row)
		return nil, error_row
	}
	defer row_result.Close()

	for row_result.Next() {
		var result model.InstitutionMembersResponse
		err_scan := row_result.StructScan(&result)
		if err_scan != nil {
			log.Println("failed to store data from database :  ", err_scan)
			return nil, err_scan
		}

		results = append(results, result)
	}
	return results, nil
}

func (r *repository) GetMemberByCompany() ([]*model.InstitutionMembersByCompanyResponse, error) {
	var results []*model.InstitutionMembersByCompanyResponse

	row_result, error_row := r.DB.Queryx(getMemberByCompanyQuerry)
	if error_row != nil {
		log.Print("failed to get member by the company given from database: ", error_row)
		return nil, error_row
	}

	defer row_result.Close()
	for row_result.Next() {
		var mock model.InstitutionMembersByCompanyResultSetResponse
		err_scan := row_result.StructScan(&mock)
		if err_scan != nil {
			log.Println("failed to store data from database :  ", err_scan)
			return nil, err_scan
		}

		result := model.InstitutionMembersByCompanyResponse{
			Institute_id:     mock.Institute_id,
			Institute_status: mock.Institute_status,
			Id:               mock.Id,
			Name:             mock.Name,
			Email:            mock.Email,
			Phone:            mock.Phone,
			Telephone:        mock.Telephone,
			Division:         mock.Division.String,
			Position:         mock.Position,
			Company_code:     mock.Company_code,
			Company_name:     mock.Company_name,

			// handling error null value posibilitites
			Institute_type: mock.Institute_type.String,
		}

		results = append(results, &result)
	}

	return results, nil
}

func (r *repository) GetAllCompany(keyword string) ([]*model.InstitutionResponse, error) {
	var results []*model.InstitutionResponse
	var row_result *sqlx.Rows
	var error_row error

	if len(keyword) <= 0 {
		row_result, error_row = r.DB.Queryx(getAllCompanyQuerry)
	} else {
		row_result, error_row = r.DB.Queryx(getAllCompanyQuerryWithKeyword, keyword)
	}
	if error_row != nil {
		log.Println("failed to get company list : ", error_row)
		return nil, error_row
	}
	defer row_result.Close()

	for row_result.Next() {
		var mock model.InstitutionResultSetResponse

		error_scan := row_result.StructScan(&mock)
		if error_scan != nil {
			log.Println("failed to save data to struct : ", error_scan)
			return nil, error_scan
		}
		result := model.InstitutionResponse{
			Id:     mock.Id,
			Code:   mock.Code,
			Name:   mock.Name,
			Status: mock.Status,

			Type: mock.Type.String,
		}
		results = append(results, &result)
	}

	return results, nil
}

func (r *repository) GetMemberByID(member_id string) (*model.InstitutionMembersDetailResponse, error) {
	var mock model.InstitutionMembersDetailResultSetResponse

	isMemberAvailable := r.CheckMemberAvailability(member_id)
	if !isMemberAvailable {
		return nil, errors.New("cannot find the specified members")
	}

	members := r.DB.QueryRowx(getMembersDetailQuerry, member_id)
	error_scan := members.StructScan(&mock)
	if error_scan != nil {
		return nil, error_scan
	}

	result := model.InstitutionMembersDetailResponse{
		Id:             mock.Id,
		Name:           mock.Name,
		Institution_id: mock.Institution_id,
		Position:       mock.Position,
		Phone:          mock.Phone,
		Telephone:      mock.Telephone,
		Email:          mock.Email,
		Division:       mock.Division,
		Division_Id:    mock.Division_Id,

		// handling null results
		Created_at: mock.Created_at.Time,
		Created_by: mock.Created_by.String,
		Updated_at: mock.Created_at.Time,
		Updated_by: mock.Updated_by.String,
	}

	return &result, nil
}

type DeleteDataProps struct {
	Id         string
	Deleted_by string
	Deleted_at time.Time
}

func (r *repository) DeleteMemberByID(props DeleteDataProps) (int64, error) {
	delete_result, error_delete := r.DB.Exec(deleteCompanyMemberQuerry, props.Id, props.Deleted_by, props.Deleted_at)
	if error_delete != nil {
		log.Println("failed to delete members data from database : ", error_delete)
		return 0, error_delete
	}

	result, error_result := delete_result.RowsAffected()
	if error_result != nil {
		log.Println("failed to get deleted member total : ", error_result)
		return 0, error_result
	}

	if result == 0 {
		log.Println("the member is failed to deleted, please check your querry / id and try again")
		return 0, errors.New("the member is failed to deleted, please check your querry / id and try again")
	}
	return result, nil
}

func (r *repository) DeleteDivisionByID(props DeleteDataProps) (int64, error) {
	log.Println(props.Id)
	delete_result, error_delete := r.DB.Exec(deleteCompanyDivisionQuerry, props.Id, props.Deleted_by, props.Deleted_at)
	if error_delete != nil {
		log.Println("failed to delete division data from database : ", error_delete)
		return 0, error_delete
	}

	result, error_result := delete_result.RowsAffected()
	if error_result != nil {
		log.Println("failed to get deleted division total : ", error_result)
		return 0, error_result
	}

	if result == 0 {
		log.Println("the division is failed to deleted, please check your querry / id and try again")
		return 0, errors.New("the division is failed to deleted, please check your querry / id and try again")
	}
	return result, nil
}

func (r *repository) GetCompanyDetail(id string) (*model.InstitutionResponse, error) {
	var result model.InstitutionResponse

	if id == "" {
		return nil, errors.New("cannot find the company : please specify the company_id")
	}

	rows := r.DB.QueryRowx(getCompanyDetailQuery, id)
	errorScan := rows.StructScan(&result)
	if errorScan != nil {
		log.Println("failed to copy company detail from database : ", errorScan)
		return nil, errorScan
	}

	return &result, nil
}
