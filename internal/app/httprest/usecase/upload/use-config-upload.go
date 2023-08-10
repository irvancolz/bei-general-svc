package upload

import (
	"be-idx-tsg/internal/app/helper"
	"fmt"
	"strings"
	"time"
	"github.com/google/uuid"
)

type UploadFileConfig struct {
	Extensions []string
	MaxSize    int64
	Directory  string
	Host       string
}

func (c *UploadFileConfig) CheckFileExt(fileext string) bool {
	var result bool
	if len(c.Extensions) <= 0 {
		return true
	}
	for _, ext := range c.Extensions {
		result = strings.EqualFold(ext, fileext)
		if result {
			return result
		}
	}
	return result
}

func (c *UploadFileConfig) GenerateFilename(filename string, date time.Time) string {
	nameSlice := strings.Split(filename, " ")
	t := helper.GetWIBLocalTime(&date)
	currentTimestr := t.Format("2006-01-02 15-04-05.000000000")
	uuid := uuid.New().String()
	nameSlice = append([]string{currentTimestr, uuid}, nameSlice...)
	return strings.Join(nameSlice, "_")  
}

func (c *UploadFileConfig) CheckFileSize(size int64) bool {
	return size <= c.MaxSize
}

func (c *UploadFileConfig) GenerateFilePath(filename string) string {
	return fmt.Sprintf("http://%s/%s/%s", c.Host, strings.ToLower(c.Directory), filename)
}
