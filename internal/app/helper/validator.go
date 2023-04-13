package helper

import "github.com/go-playground/validator/v10"

func Validator() *validator.Validate {
	return validator.New()
}
