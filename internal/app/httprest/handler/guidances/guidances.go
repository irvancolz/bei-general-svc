package guidances

import (
	usecase "be-idx-tsg/internal/app/httprest/usecase/guidances"
	"be-idx-tsg/internal/pkg/httpresponse"
	"errors"

	"github.com/gin-gonic/gin"
)

type GuidanceHandler interface {
	CreateNewGuidance(c *gin.Context)
	UpdateExistingGuidance(c *gin.Context)
	GetAllGuidanceBasedOnType(C *gin.Context)
	DeleteGuidances(c *gin.Context)
	CreateNewFiles(c *gin.Context)
	UpdateExistingFiles(c *gin.Context)
	CreateNewRegulation(c *gin.Context)
	UpdateExistingRegulation(c *gin.Context)
}

type guidancehandler struct {
	usecase usecase.GuidancesRegulationAndFileUsecaseInterface
}

func NewGuidanceHandler() GuidanceHandler {
	return &guidancehandler{
		usecase: usecase.NewGuidanceUsecase(),
	}
}

func (h *guidancehandler) CreateNewGuidance(c *gin.Context) {
	var request usecase.CreateNewGuidanceProps
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_params))
		return
	}
	result, error_result := h.usecase.CreateNewGuidance(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, result))
}
func (h *guidancehandler) UpdateExistingGuidance(c *gin.Context) {
	var request usecase.UpdateExsistingGuidances
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_params))
		return
	}
	error_result := h.usecase.UpdateExistingGuidance(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil))
}
func (h *guidancehandler) GetAllGuidanceBasedOnType(c *gin.Context) {
	types := c.Query("type")

	switch types {
	case "Guidebook":
		{
			result, error_result := h.usecase.GetAllGuidanceBasedOnType(c, types)
			if error_result != nil {
				c.JSON(httpresponse.Format(httpresponse.READFAILED_400, error_result))
				return
			}
			if len(result) == 0 {
				c.JSON(httpresponse.Format(httpresponse.CONTENTNOTFOUND_404, error_result))
				return
			}
			c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))

			break
		}
	case "File":
		{
			result, error_result := h.usecase.GetAllFilesOnType(c, types)
			if error_result != nil {
				c.JSON(httpresponse.Format(httpresponse.READFAILED_400, error_result))
				return
			}
			if len(result) == 0 {
				c.JSON(httpresponse.Format(httpresponse.CONTENTNOTFOUND_404, error_result))
				return
			}
			c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))

			break
		}

	case "Regulation":
		{
			result, error_result := h.usecase.GetAllRegulationsBasedOnType(c, types)
			if error_result != nil {
				c.JSON(httpresponse.Format(httpresponse.READFAILED_400, error_result))
				return
			}
			if len(result) == 0 {
				c.JSON(httpresponse.Format(httpresponse.CONTENTNOTFOUND_404, error_result))
				return
			}
			c.JSON(httpresponse.Format(httpresponse.READSUCCESS_200, nil, result))

			break
		}
	default:
		{
			c.JSON(httpresponse.Format(httpresponse.READFAILED_400, errors.New("type is unexpected. please change your type to Guidebook || File || Regulation")))
		}

	}
}
func (h *guidancehandler) CreateNewFiles(c *gin.Context) {
	var request usecase.CreateNewGuidanceProps
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_params))
		return
	}
	result, error_result := h.usecase.CreateNewFiles(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, result))
}
func (h *guidancehandler) UpdateExistingFiles(c *gin.Context) {
	var request usecase.UpdateExsistingGuidances
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_params))
		return
	}
	error_result := h.usecase.UpdateExistingFiles(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil))
}
func (h *guidancehandler) CreateNewRegulation(c *gin.Context) {
	var request usecase.CreateNewRegulationsProps
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_params))
		return
	}
	result, error_result := h.usecase.CreateNewRegulations(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.CREATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.CREATESUCCESS_200, nil, result))
}
func (h *guidancehandler) UpdateExistingRegulation(c *gin.Context) {
	var request usecase.UpdateExistingRegulationsProps
	if error_params := c.ShouldBindJSON(&request); error_params != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_params))
		return
	}
	error_result := h.usecase.UpdateExistingRegulations(c, request)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.UPDATEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.UPDATESUCCESS_200, nil))
}
func (h *guidancehandler) DeleteGuidances(c *gin.Context) {
	id := c.Query("id")
	error_result := h.usecase.DeleteGuidances(c, id)
	if error_result != nil {
		c.JSON(httpresponse.Format(httpresponse.DELETEFAILED_400, error_result))
		return
	}
	c.JSON(httpresponse.Format(httpresponse.DELETESUCCESS_200, nil))
}
