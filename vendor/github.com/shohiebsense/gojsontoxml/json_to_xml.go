package gojsontoxml

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)


//and other special characters that might break it, a bit lazy to handle the edge cases :)
func getSafeKey(key string) string {
	return strings.Replace(key, ":", "", -1)
}

func parseXml(doc *etree.Document, obj *etree.Element, data map[string]interface{}) {

	for k, element := range data {
		k = getSafeKey(k)
		switch vv := element.(type) {
		case string:
			obj.CreateElement(k).SetText(vv)
		case float64:
			floatstr := fmt.Sprintf("%f", vv)
			obj.CreateElement(k).SetText(floatstr)
		case []string:
			newElement := obj.CreateElement(k)
			parseSliceOfString(newElement, k, vv)
		case []interface{}:
			newElement := obj.CreateElement(k)
			parseSliceOfInterface(doc, newElement, k, vv)
		case map[string]interface{}:
			newElement := obj.CreateElement(k)
			parseXml(doc, newElement, vv)
		case bool:
			obj.CreateElement(k).SetText(strconv.FormatBool(vv))
		default:
			xType := reflect.TypeOf(element)
			fmt.Println("this guy has not been handled: ", xType)
		}
	}
}

func parseSliceOfInterface(doc *etree.Document, obj *etree.Element, key string, data []interface{}) {
	defaultItemTagName := "item"
	for _, element := range data {
		switch vv := element.(type) {
		case string:
			obj.CreateElement(defaultItemTagName).SetText(vv)
		case int:
			intStr := strconv.Itoa(vv)
			obj.CreateElement(defaultItemTagName).SetText(intStr)
		case float64:
			floatstr := fmt.Sprintf("%f", vv)
			obj.CreateElement(defaultItemTagName).SetText(floatstr)
		case bool:
			obj.CreateElement(defaultItemTagName).SetText(strconv.FormatBool(vv))
		case []interface{}:
			newElement := obj.CreateElement(key)
			parseSliceOfInterface(doc, newElement, key, vv)
		case interface{}:
			newElement := obj.CreateElement(key)
			parseInterface(doc, newElement, key, vv)
		default:
			xType := reflect.TypeOf(element)
			fmt.Println("this guy has not been handled: ", xType)

		}
	}
}

func parseInterface(doc *etree.Document, obj *etree.Element, key string, data interface{}) {
	switch v := data.(type) {
	case string:
		obj.CreateElement(key).SetText(v)
	case map[string]interface{}:
		parseXml(doc, obj, v)
	default:
		xType := reflect.TypeOf(v)
		fmt.Println("this guy has not been handled: ", xType)
	}
}

func parseSliceOfString(obj *etree.Element, key string, data []string) {
	for _, element := range data {
		obj.CreateAttr(key, element)
	}
}

func JsonToXml(data map[string]interface{}) ([]byte, error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("Object")
	parseXml(doc, element, data)
	doc.Indent(2)
	
	dataBytes, err := doc.WriteToBytes()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return dataBytes, nil
}