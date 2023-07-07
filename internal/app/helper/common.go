package helper

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ConvertArrayStringtoString(array []string) string {
	result := make([]string, len(array))
	for i, v := range array {
		result[i] = fmt.Sprintf("'%s'", v)
	}
	finalResult := strings.Join(result, ",")
	return finalResult
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_+=")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func GetWIBLocalTime(date *time.Time) time.Time {
	var timeToParse time.Time
	if date != nil {
		timeToParse = *date
	} else {
		timeToParse = time.Now()
	}

	t, _ := TimeIn(timeToParse, "Asia/Jakarta")
	return t
}

func ConvertListInterfaceToListString(list []interface{}) []string {
	stringList := make([]string, len(list))

	for i, v := range list {
		if str, ok := v.(string); ok {
			stringList[i] = str
		} else {
			stringList[i] = ""
		}
	}
	return stringList
}

func isArray(arg interface{}) bool {
	argType := reflect.TypeOf(arg)
	return argType.Kind() == reflect.Array || argType.Kind() == reflect.Slice
}

func FilterData(data interface{}) (interface{}, error) {
	if !isArray(data) {
		log.Println("the parameter should be and array")
		return nil, errors.New("the parameter should be and array")
	}

	return "a", nil
}

func ConvertToMap(dataStruct []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, data := range dataStruct {
		mapResult := make(map[string]interface{})
		baseStruct := reflect.ValueOf(data)
		baseStructTotalProps := baseStruct.NumField()

		for i := 0; i < baseStructTotalProps; i++ {
			mapResult[strings.ToLower(baseStruct.Type().Field(i).Name)] = baseStruct.Field(i).Interface()
		}
		result = append(result, mapResult)
	}

	return result
}

func GetMapKeys(data map[string]interface{}) []string {
	var results []string
	for key := range data {
		results = append(results, key)
	}
	return results
}

func IsContainsString(list []string, data string) bool {
	for _, item := range list {
		if item == data {
			return true
		}
	}
	return false
}

func IsContainsBool(list []bool, data bool) bool {
	for _, item := range list {
		if item == data {
			return true
		}
	}
	return false
}

func IsContains[T comparable](list []T, data T) bool {
	for _, item := range list {
		if item == data {
			return true
		}
	}
	return false
}

// expected result = function will return array of rows data from data given so it can be used to provide data for export feature
// mechanic get all values from properties in given params
// order is mandatory to keep data order consistency cannot give 0 order
// make sure there is no pointer used on this func, if there is pointer it will raise panic error
func MapToArray(data map[string]interface{}, order []string) []string {
	var result []string
	if len(order) <= 0 {
		log.Println("failed to convert map to array: please specify at least one array order to prevent unconsistent result")
		return result
	}
	for _, orderValue := range order {
		for key := range data {

			// if key == orderValue {
			if strings.EqualFold(key, orderValue) {
				result = append(result, fmt.Sprintf("%v", data[key]))
			}
		}
	}
	return result
}

func StructToArray(data interface{}, order []string) []string {
	var result []string

	if len(order) <= 0 {
		log.Println("failed to convert struct to array: please specify at least one array order to prevent unconsistent result")
		return result
	}

	dataType := reflect.ValueOf(data)
	dataProps := dataType.NumField()

	for _, arrkey := range order {

		for i := 0; i < dataProps; i++ {
			if strings.EqualFold(arrkey, dataType.Type().Field(i).Name) {
				result = append(result, fmt.Sprintf("%v", dataType.Field(i).Interface()))
			}

		}
	}
	return result
}

func generateFileNames(fileName, separator string, date time.Time) string {
	return fileName + "_" + strconv.Itoa(int(date.UnixNano()))
}

func IsString(val interface{}) bool {
	_, isString := val.(string)
	return isString
}
