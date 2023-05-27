package helper

import (
	"strconv"
	"testing"
)

func TestIsArray(t *testing.T) {
	char := []string{"a", "b"}
	result := isArray(char)
	if !result {
		t.Error("the args is an array")
	}
}

func TestIsArray2(t *testing.T) {
	char := "a"
	result := isArray(char)
	if result {
		t.Error("the args is not an array")
	}
}

type dummyStruct struct {
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
		data := dummyStruct{
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
