package unggahberkas

import (
	"testing"
	"time"
)

func TestBuildReference(t *testing.T) {
	mockTime := time.Date(2022, 1, 24, 10, 0, 0, 0, time.UTC)
	result := buildNoReference("l", mockTime, 1)
	t.Log(result)
}
