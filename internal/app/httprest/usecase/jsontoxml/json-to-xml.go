package jsontoxml

import (
	"log"
	"github.com/shohiebsense/gojsontoxml"
)

type Usecase interface {
	ToXml(data map[string]interface{})([]byte, error) 
}

type usecase struct {
}

// ToXml implements Usecase
func ToXmlUseCase() Usecase {
	return &usecase{}
}


func (*usecase) ToXml(data map[string]interface{}) ([]byte, error) {

	dataBytes, err := gojsontoxml.JsonToXml(data)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return dataBytes, nil
}
