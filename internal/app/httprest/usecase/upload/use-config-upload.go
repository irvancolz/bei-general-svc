package upload

import (
	"be-idx-tsg/internal/app/helper"
	"fmt"
	"strings"
	"time"
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
	currentTimestr := t.Format("2006-01-02 15-04")
	nameSlice = append([]string{currentTimestr}, nameSlice...)
	return strings.Join(nameSlice, "_")
}

func (c *UploadFileConfig) CheckFileSize(size int64) bool {
	return size <= c.MaxSize
}

func (c *UploadFileConfig) GenerateFilePath(filename string) string {
	return fmt.Sprintf("%s/uploaded/%s/%s", c.Host, c.Directory, filename)
}
