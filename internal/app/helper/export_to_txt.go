package helper

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type ExportTableToTxtProps struct {
	Filename    string
	Data        [][]string
	Header      [][]string
	ColumnWidth []int
}

func ExportTableToTxt(props ExportTableToTxtProps) (string, error) {

	fileName := props.Filename
	data := props.Header
	data = append(data, props.Data...)
	columnWidth := props.ColumnWidth

	fileresultName := generateFileNames(fileName, "_", time.Now())
	file, errorCreate := os.Create(fileresultName + ".txt")
	if errorCreate != nil {
		log.Println("failed to create txt files :", errorCreate)
		return "", errorCreate
	}
	defer file.Close()

	beautifiedData := formatRowsData(data, columnWidth)
	withStylingData := drawTxtTable(beautifiedData, columnWidth)

	txtFile := bufio.NewWriter(file)
	defer txtFile.Flush()

	_, errorResult := txtFile.WriteString(strings.Join(withStylingData, "\n"))
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

func drawTxtTable(data [][]string, columnWidth []int) []string {

	var result []string
	for j, item := range data {
		var beautifiedRows strings.Builder
		for i, text := range item {
			beautifiedRows.WriteString("| " + text)
			for space := len(text); space <= columnWidth[i]; space++ {
				beautifiedRows.WriteString("\u0020")
			}
			beautifiedRows.WriteString("|")
		}

		// header border bottom
		if j == 0 {
			beautifiedRows.WriteString("\n")
			for i := 0; i < len(data[j]); i++ {
				for char := 0; char <= columnWidth[i]; char++ {
					beautifiedRows.WriteString("-")
				}
			}
		}

		result = append(result, beautifiedRows.String())
	}
	return result
}

func checkOverlappedText(data [][]string, widths []int) bool {
	var result bool

	for _, line := range data {

		for i, word := range line {
			maxWordLen := func() int {
				if i >= len(widths) {
					return 20
				}

				return widths[i]
			}()
			if len(word) > maxWordLen {
				return true
			}
		}
	}

	return result
}

func getRowMaxContent(data []string) int {
	var result int

	for _, word := range data {
		if len(word) >= result {
			result = len(word)
		}
	}

	return result
}

func formatRowsData(data [][]string, maxWidths []int) [][]string {

	var result [][]string
	for line := 0; line < len(data); line++ {
		copyOfLine := data[line]
		newLineMock := make([]string, len(copyOfLine))
		for wordIdx, word := range data[line] {
			maxWordLen := func() int {
				if wordIdx >= len(maxWidths) {
					return 20
				}

				return maxWidths[wordIdx]
			}()
			if len(word) > maxWordLen {
				slicedWord := string([]byte(word)[:maxWordLen])
				remainingWord := string([]byte(word)[maxWordLen:])
				copyOfLine[wordIdx] = slicedWord
				newLineMock[wordIdx] = remainingWord
			}
		}
		result = append(result, copyOfLine)
		if getRowMaxContent(newLineMock) != 0 {
			result = append(result, newLineMock)
		}
	}

	if !checkOverlappedText(result, maxWidths) {
		return result
	}

	return formatRowsData(result, maxWidths)

}
