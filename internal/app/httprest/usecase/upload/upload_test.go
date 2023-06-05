package upload

import (
	"path/filepath"
	"testing"
	"time"
)

func TestCheckExtensions(t *testing.T) {
	config := UploadFileConfig{
		Extensions: []string{"a", "b", "c"},
	}

	result := config.CheckFileExt("a")
	resultd := config.CheckFileExt("d")

	if !result {
		t.Error("the file extention should be allowed")
	}
	if resultd {
		t.Error("the file extention should not be allowed")

	}
}

func TestGenerateFilename(t *testing.T) {
	config := UploadFileConfig{}
	name := "test file.pdf"
	mockTime, _ := time.Parse("2006-01-02", "2006-01-02")
	result := config.GenerateFilename(name, mockTime)
	if result != "2006-01-02_00-00-00_test_file.pdf" {
		t.Log(result)
		t.Error("the file name generated didnt fullfilled the requirements")
	}

}

func TestGetFilePath(t *testing.T) {
	path := "http://localhost:9011/api/uploaded/test/2023-05-07_15-20-21_dashboard.pdf"
	expected := "uploaded/test/2023-05-07_15-20-21_dashboard.pdf"
	result := GetFilePath(path)
	if result != filepath.FromSlash(expected) {
		t.Log(result)
		t.Log(filepath.FromSlash(expected))
		t.Error("the resulted string does not as expected")
	}
}
