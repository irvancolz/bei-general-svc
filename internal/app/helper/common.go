package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func ConvertArrayStringtoString(array []string) string {
	result := make([]string, len(array))
	for i, v := range array {
		result[i] = fmt.Sprintf("'%s'", v)
	}
	finalResult := strings.Join(result, ",")
	return finalResult
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_+=")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}
