package helper

import (
	"strconv"
	"testing"
	"time"
)

func TestIsArray(t *testing.T) {
	char := []string{"a", "b"}
	result := isArray(char)
	if !result {
		t.Error("the args is an array")
	}
}

func TestIsNotAnArray(t *testing.T) {
	char := "a"
	result := isArray(char)
	if result {
		t.Error("the args is not an array")
	}
}

type convertToMapMockStruct struct {
	Id       int64
	Name     string
	Is_owned bool
}

func areSlicesEqual(slice1, slice2 interface{}) bool {
	if len(slice1.([]map[string]interface{})) != len(slice2.([]map[string]interface{})) {
		return false
	}
	for i := range slice1.([]map[string]interface{}) {
		if areMapsEqual(slice1.([]map[string]interface{})[i], slice2.([]map[string]interface{})[i]) {
			return false
		}
	}
	return true
}

func areMapsEqual(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		value2, exists := map2[key]
		if !exists || value1 != value2 {
			return false
		}
	}

	return true
}

func TestConvertToMap(t *testing.T) {
	var datas []interface{}
	count := 3
	for i := 0; i < count; i++ {
		data := convertToMapMockStruct{
			Is_owned: true,
			Id:       int64(i) + 1,
			Name:     "data ke : " + strconv.Itoa(i+1),
		}
		datas = append(datas, data)
	}

	var expected []map[string]interface{}

	for i := 0; i < count; i++ {
		expect := make(map[string]interface{})
		expect["id"] = i + 1
		expect["is_owned"] = true
		expect["name"] = "data ke : " + strconv.Itoa(i+1)

		expected = append(expected, expect)
	}

	results := ConvertToMap(datas)
	if !areSlicesEqual(results, expected) {
		t.Error("the resulted map should be same as expected")
	}
}

func TestGenerateParamsTimeName(t *testing.T) {
	paramsList := []string{"lorem", "ipsum", "dolor"}
	result := generateTimeRangeParamList(paramsList)
	expected := []string{"lorem_from", "ipsum_from", "dolor_from", "lorem_end", "ipsum_end", "dolor_end"}
	for _, data := range result {
		if !IsContains(expected, data) {
			t.Log(data)
			t.Error()
		}
	}
}

func TestCheckStarterTimeParam(t *testing.T) {
	param := "lorem_from"
	paramEnd := "lorem_end"
	if !isStarterTimeParam(param) {
		t.Error()
	}
	if isStarterTimeParam(paramEnd) {
		t.Error()
	}
}

func TestGetTimeBaseName(t *testing.T) {
	param := "lorem_ipsum_from"
	paramEnd := "lorem_ipsum_end"
	result := getBaseTimeRange(param)
	resultend := getBaseTimeRange(paramEnd)
	if result != "lorem_ipsum" {
		t.Error()
	}
	if resultend != "lorem_ipsum" {
		t.Error()
	}
}

func TestGenerateParams(t *testing.T) {
	var exportedData []map[string]interface{}
	for i := 0; i < 10; i++ {
		data := make(map[string]interface{})
		data["count"] = strconv.Itoa(i + 1)
		exportedData = append(exportedData, data)
	}

	filteredParameter := generateFilterParameter(exportedData)
	t.Log(filteredParameter)
}

func TestGetLocalTime(t *testing.T) {
	utc := time.Date(2022, 01, 02, 11, 00, 00, 00, time.UTC)
	jktTime := GetWIBLocalTime(&utc)
	t.Log(jktTime)
}

func TestGenerateHumanDate(t *testing.T) {
	utc := time.Date(2022, 01, 02, 11, 20, 00, 00, time.UTC)
	t.Log(ConvertTimeToHumanDate(utc))
}
