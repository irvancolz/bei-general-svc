package jsontoxml

import (
	"be-idx-tsg/internal/app/httprest/model"
	JsonToXmlUseCase "be-idx-tsg/internal/app/httprest/usecase/jsontoxml"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ToXml(c *gin.Context)
}

type handler struct {
	jsonToXmlUseCase JsonToXmlUseCase.Usecase
}

func NewHandler() Handler {
	return &handler{
		JsonToXmlUseCase.ToXmlUseCase(),
	}
}

func (jsonToXmlHandler *handler) ToXml(c *gin.Context) {
	var data map[string]interface{}

	err := c.ShouldBindJSON(&data)

	if err != nil {
		log.Println(err)
		model.GenerateInvalidJsonResponse(c, err)
		return
	}

	dataBytes, err := jsonToXmlHandler.jsonToXmlUseCase.ToXml(data)

	if err != nil {
		model.GenerateReadErrorResponse(c, err)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/xml")
	_, err = c.Writer.Write(dataBytes)

	if err != nil {
		model.GenerateReadErrorResponse(c, err)
	}
}
