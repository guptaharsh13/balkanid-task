package utils

import (
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func SetupValidator() error {
	validate = validator.New()
	return nil
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
