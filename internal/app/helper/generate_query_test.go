package helper

import (
	"testing"
	"time"
)

func TestSearchQueryGenerator(t *testing.T) {
	config := SearchQueryGenerator{
		ColumnScanned: []string{"a", "b"},
		TableName:     "test",
	}
	expected := "\nAND(\nLOWER( a ) LIKE LOWER('%c%')\nOR LOWER( b ) LIKE LOWER('%c%'))\nAND id IN(\nSELECT id FROM test WHERE (\nLOWER( a ) LIKE LOWER('%d%')\nOR LOWER( b ) LIKE LOWER('%d%')))\nAND id IN(\nSELECT id FROM test WHERE (\nLOWER( a ) LIKE LOWER('%e%')\nOR LOWER( b ) LIKE LOWER('%e%')))"
	result := config.GenerateSearchQuery([]string{"c", "d", "e"}, "")
	if result != expected {
		t.Log("expected query :", expected)
		t.Error("the query search resulted is not matched with the expected : ", result)
	}
}

func TestSearchQueryGeneratorWithoutColumn(t *testing.T) {
	config := SearchQueryGenerator{
		ColumnScanned: []string{"a", "b"},
	}
	expected := ""
	result := config.GenerateSearchQuery([]string{"c"}, "")

	if result != expected {
		t.Error("the query search resulted is not matched with the expected : ", result)
	}

}

func TestSearchQueryGeneratorWithoutValue(t *testing.T) {
	config := SearchQueryGenerator{
		ColumnScanned: []string{},
	}
	expected := ""
	result := config.GenerateSearchQuery([]string{"d"}, "")

	if result != expected {
		t.Error("the query search resulted is not matched with the expected : ", result)
	}

}

func TestTimeConverter(t *testing.T) {
	// Membuat waktu dengan tanggal 26 Mei 2023 jam 00.00
	t1 := time.Date(2023, time.May, 26, 0, 0, 0, 0, time.Local)
	unixTimeStart := t1.Unix()

	// Membuat waktu dengan tanggal 26 Mei 2023 jam 23.59
	t2 := time.Date(2023, time.May, 26, 23, 59, 0, 0, time.Local)
	unixTimeEnd := t2.Unix()

	timeStart := ConvertUnixToDateString(unixTimeStart, "")
	timeEnd := ConvertUnixToDateString(unixTimeEnd, "")

	if timeStart != timeEnd {
		t.Error("the time results is not same")
	}
}
