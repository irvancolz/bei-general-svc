package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		model.GenerateFlowErrorResponse(c, err)
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" && ext != ".csv" && ext != ".xls" && ext != ".xlsx" {
		model.GenerateFlowErrorFromMessageResponse(c, "invalid file type")
		return
	}

	// Create directory if it doesn't exist
	uploadDir := "./uploaded"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		model.GenerateInternalErrorResponse(c, err)
		return
	}

	// get original file name
	filename := file.Filename

	// replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// generate timestamp
	t := time.Now().Format("2006-01-02_15-04-05")

	// Generate a unique filename
	newName := filename[:len(filename)-len(ext)] + "-" + t + ext

	// Save the file to disk
	err = c.SaveUploadedFile(file, filepath.Join(uploadDir, newName))
	if err != nil {
		model.GenerateInternalErrorResponse(c, err)
		return
	}

	host := os.Getenv("DIR_HOST")
	url := fmt.Sprintf("%sapi/uploaded/%s", host, newName)

	c.JSON(http.StatusOK, gin.H{
		"filename": newName,
		"url":      url,
		"size":     strconv.Itoa(int(math.Round(float64(file.Size)/1000))) + " KB",
		"message":  "sukses menambahkan file",
	})
}

func getFileExtension(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) > 1 {
		return "." + split[len(split)-1]
	}
	return ""
}

func DeleteFile(c *gin.Context) {
	fileName := c.Query("filename")

	// get file extension from filename
	ext := getFileExtension(fileName)

	// check file extension and delete file from appropriate directory
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".pdf" || ext == ".csv" || ext == ".xls" || ext == ".xlsx" {
		// file does not exist, return error message
		if _, err := os.Stat("./uploaded/" + filepath.Base(fileName)); os.IsNotExist(err) {
			model.GenerateIFileNotFoundErrorResponse(c, err)
			return
		}
		err := os.Remove("./uploaded/" + filepath.Base(fileName))
		if err != nil {
			model.GenerateRemoveFileErrorResponse(c, err)
			return
		}
	} else {
		// unsupported file extension, return error
		model.GenerateFlowErrorFromMessageResponse(c, "gagal menghapus file. Unsupported file extension")
		return
	}

	// return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "sukses menghapus file",
	})
}

func GetFile(c *gin.Context) {
	filename := c.Param("filename")
	c.File("./uploaded/" + filename)
}
