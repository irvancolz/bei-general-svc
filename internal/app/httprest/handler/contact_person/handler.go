package contactperson

import (
	"be-idx-tsg/internal/app/httprest/model"
	usecase "be-idx-tsg/internal/app/httprest/usecase/contact_person"
	"be-idx-tsg/internal/pkg/httpresponse"

	"github.com/gin-gonic/gin"
)

type ContactPersonHandlerInterface interface {
	SynchronizeInstitutionProfile(c *gin.Context)
	AddDivision(c *gin.Context)
	EditDivision(c *gin.Context)
	GetAllDivision(c *gin.Context)
	GetAllDivisionByCompany(c *gin.Context)
	AddMember(c *gin.Context)
	EditMember(c *gin.Context)
	GetMemberByDivision(c *gin.Context)
	GetMemberByDivisionAndCompanyID(c *gin.Context)
	GetMemberByCompanyID(c *gin.Context)
	GetAllCompanyByType(c *gin.Context)
	GetMemberByID(c *gin.Context)
	DeleteDivisionByID(c *gin.Context)
	DeleteMemberByID(c *gin.Context)
	GetMemberByCompanyType(c *gin.Context)
	SearchCompany(c *gin.Context)
	ExportMember(c *gin.Context)
	GetAllMemberEmail(c *gin.Context)
}

type handler struct {
	Usecase usecase.ContactPersonUsecaseInterface
}

func NewHandler() ContactPersonHandlerInterface {
	return &handler{
		Usecase: usecase.NewUsecase(),
	}
}

func (h *handler) GetAllMemberEmail(c *gin.Context) {
	result, errResult := h.Usecase.GetAllMembersEmail(c)
	if errResult != nil {
		model.GenerateReadErrorResponse(c, errResult)
		return
	}
	if !c.Writer.Written() {

		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))
	}

}

func (h *handler) SynchronizeInstitutionProfile(c *gin.Context) {
	companyType := c.Query("company_type")
	data, error_data := h.Usecase.SynchronizeInstitutionProfile(c, companyType)
	if error_data != nil {
		model.GenerateUpdateErrorResponse(c, error_data)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

type AddDivisionProps struct {
	Name string `json:"name" binding:"required"`
}

func (h *handler) AddDivision(c *gin.Context) {
	var request AddDivisionProps

	if error_param := c.ShouldBindJSON(&request); error_param != nil {
		model.GenerateInsertErrorResponse(c, error_param)
		return
	}
	data, error_data := h.Usecase.AddDivision(c, request.Name)
	if error_data != nil {
		model.GenerateInsertErrorResponse(c, error_data)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (h *handler) AddMember(c *gin.Context) {
	var request usecase.AddMemberProps

	if error_param := c.ShouldBindJSON(&request); error_param != nil {
		model.GenerateInsertErrorResponse(c, error_param)
		return
	}

	data, error_data := h.Usecase.AddMember(c, request)
	if error_data != nil {
		model.GenerateInsertErrorResponse(c, error_data)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, data))
}

func (h *handler) GetMemberByDivision(c *gin.Context) {

	divisions := c.QueryArray("division")
	types := c.Query("company_type")

	data, error_data := h.Usecase.GetMemberByDivision(types, divisions)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) GetMemberByDivisionAndCompanyID(c *gin.Context) {

	divisions := c.QueryArray("division")
	company := c.QueryArray("company_id")

	data, error_data := h.Usecase.GetMemberByDivisionAndCompanyID(divisions, company)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) GetAllDivision(c *gin.Context) {
	data, error_data := h.Usecase.GetAllDivision(c)

	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if !c.Writer.Written() {
		if len(data) == 0 {
			c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
			return
		}
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
	}

}

func (h *handler) GetAllDivisionByCompany(c *gin.Context) {
	company_id := c.Query("company_id")
	data, error_data := h.Usecase.GetAllDivisionByCompany(company_id)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) GetMemberByCompanyID(c *gin.Context) {
	company_id := c.QueryArray("id")
	data, error_data := h.Usecase.GetMemberByCompanyID(company_id)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) GetAllCompanyByType(c *gin.Context) {
	company_type := c.Query("type")
	data, error_data := h.Usecase.GetAllCompanyByType(company_type)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) SearchCompany(c *gin.Context) {
	keyword := c.Query("search")
	data, error_data := h.Usecase.SearchCompany(keyword)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	if len(data) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) GetMemberByID(c *gin.Context) {
	member_id := c.Query("id")
	data, error_data := h.Usecase.GetMemberByID(member_id)
	if error_data != nil {
		model.GenerateReadErrorResponse(c, error_data)
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, data))
}

func (h *handler) EditMember(c *gin.Context) {
	var request usecase.EditMemberProps

	if error_param := c.ShouldBindJSON(&request); error_param != nil {
		model.GenerateUpdateErrorResponse(c, error_param)
		return
	}

	data, error_data := h.Usecase.EditMember(c, request)
	if error_data != nil {
		model.GenerateUpdateErrorResponse(c, error_data)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

func (h *handler) EditDivision(c *gin.Context) {
	var request usecase.EditDivisionprops

	if error_param := c.ShouldBindJSON(&request); error_param != nil {
		model.GenerateUpdateErrorResponse(c, error_param)
		return
	}

	data, error_data := h.Usecase.EditDivision(c, request)
	if error_data != nil {
		model.GenerateUpdateErrorResponse(c, error_data)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil, data))
}

func (h *handler) DeleteDivisionByID(c *gin.Context) {
	division_id := c.Query("id")

	result, error_result := h.Usecase.DeleteDivisionByID(c, division_id)
	if error_result != nil {
		model.GenerateDeleteErrorResponse(c, error_result)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, result))
}

func (h *handler) DeleteMemberByID(c *gin.Context) {
	division_id := c.Query("id")

	result, error_result := h.Usecase.DeleteMemberByID(c, division_id)
	if error_result != nil {
		model.GenerateDeleteErrorResponse(c, error_result)
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil, result))
}

func (h *handler) GetMemberByCompanyType(c *gin.Context) {
	company_type := c.Query("type")

	result, error_result := h.Usecase.GetMemberByCompanyType(company_type)
	if error_result != nil {
		model.GenerateReadErrorResponse(c, error_result)
		return
	}

	if len(result) == 0 {
		c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, make([]map[string]interface{}, 0)))
		return
	}

	c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))
}

func (h *handler) ExportMember(c *gin.Context) {
	companyType := c.Query("company_type")
	companyId := c.Query("company_id")
	divisionId := c.Query("division_id")

	errorExport := h.Usecase.ExportMember(c, companyType, companyId, divisionId)
	if errorExport != nil {
		model.GenerateReadErrorResponse(c, errorExport)
		return
	}
}
