package helper

import (
	"log"

	"github.com/beevik/etree"
)

func GetXmlDocument(dataBytes []byte) (*etree.Document, error) {
	bodyDocument := etree.NewDocument()
	err := bodyDocument.ReadFromBytes(dataBytes)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return bodyDocument, nil
}

func WrapToXml(data []byte, name string) ([]byte, error) {
	// Parse the JSON data into a Registration struct

	// Create a new XML document

	// Create the root element "Participant"
	dataDocument, err := GetXmlDocument(data)

	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	doc.Root().AddChild(dataDocument.Root())

	dataDocument.SetRoot(&dataDocument.Element)
	//bodyDocument.AddChild(&dataDocument.Element)
	
	// Create an "ID" element and set its text to the provided 'id'
	itemXml, err := doc.WriteToBytes()
		
	if err != nil {
		return nil, err
	}

	return itemXml, nil
}
