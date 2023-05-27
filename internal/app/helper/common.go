package helper

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
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
