package helper

import (
	"encoding/csv"
	"log"
	"os"
	"time"
)

func ExportTableToCsv(fileName string, data [][]string) (string, error) {
	fileNameResult := generateFileNames(fileName, "_", time.Now()) + ".csv"
	file, errorFile := os.Create(fileNameResult)
	if errorFile != nil {
		log.Println("failed to create csv file", errorFile)
		return "", errorFile
	}
	defer file.Close()

	csvFile := csv.NewWriter(file)
	defer csvFile.Flush()

	csvFile.WriteAll(data)

	return fileNameResult, nil
}
