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
	columnWidth := getColumnMaxWidth(data)
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
			for space := len(text); space <= columnWidth[i]; space++ {
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

func getColumnMaxWidth(data [][]string) []int {
	var result []int
	var tablecolumnContent [][]string

	for i := 0; i < len(data[0]); i++ {
		tablecolumnContent = append(tablecolumnContent, []string{})
	}

	for _, rows := range data {
		for d, text := range rows {
			tablecolumnContent[d] = append(tablecolumnContent[d], text)
		}
	}

	for _, rows := range tablecolumnContent {
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
