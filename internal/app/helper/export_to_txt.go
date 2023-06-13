package helper

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

func ExportTableToTxt(fileName string, data [][]string) (string, error) {
	fileresultName := generateFileNames(fileName, "_", time.Now())
	collumnWidth := getCollumnMaxWidth(data)
	file, errorCreate := os.Create(fileresultName + ".txt")
	if errorCreate != nil {
		log.Println("failed to create txt files :", errorCreate)
		return "", errorCreate
	}
	defer file.Close()
	var beautifiedData []string

	for _, item := range data {
		var beautifiedRows strings.Builder
		for i, text := range item {
			beautifiedRows.WriteString(text)
			for space := len(text); space <= collumnWidth[i]; space++ {
				beautifiedRows.WriteString("\u0020")
			}
			beautifiedRows.WriteString("\u0020")
		}
		beautifiedData = append(beautifiedData, beautifiedRows.String())
	}

	txtFile := bufio.NewWriter(file)
	defer txtFile.Flush()

	_, errorResult := txtFile.WriteString(strings.Join(beautifiedData, "\n"))
	if errorResult != nil {
		log.Println("failed to write data to txt :", errorResult)
		return "", errorResult
	}
	return fileresultName + ".txt", nil
}

func getCollumnMaxWidth(data [][]string) []int {
	var result []int
	var tableCollumnContent [][]string

	for i := 0; i < len(data[0]); i++ {
		tableCollumnContent = append(tableCollumnContent, []string{})
	}

	for _, rows := range data {
		for d, text := range rows {
			tableCollumnContent[d] = append(tableCollumnContent[d], text)
		}
	}

	for _, rows := range tableCollumnContent {
		var max int
		for _, text := range rows {
			if len(text) > max {
				max = len(text)
			}
		}
		result = append(result, max)
	}
	return result
}
