package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func SetupValidator() error {

	validate = validator.New()
	if err := validate.RegisterValidation("username", usernameValidator); err != nil {
		return err
	}
	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		return err
	}
	return nil
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func IsValidEmail(email string) bool {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func usernameValidator(fl validator.FieldLevel) bool {

	username := fl.Field().String()
	username = strings.TrimSpace(username)
	if len(username) < 5 || len(username) > 25 {
		return false
	}
	return true
}

func passwordValidator(fl validator.FieldLevel) bool {

	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return false
	}
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLowercase {
		return false
	}
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false
	}
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	return hasSpecial
}
