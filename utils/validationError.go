package utils

import "github.com/go-playground/validator"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type validationErrorsResponse struct {
	Success bool              `json:"success"`
	Errors  []ValidationError `json:"errors"`
}

func formValidationMessage(tag string) string {

	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "username":
		return "Username must be between 5 and 25 characters long"
	case "password":
		return "Password must contain at least 8 characters, one uppercase letter, one number, and one special character"
	}
	return ""
}

func ValidationErrorResponse(validationErrors error) validationErrorsResponse {

	errors := []ValidationError{}
	for _, validationError := range validationErrors.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   validationError.Field(),
			Message: formValidationMessage(validationError.Tag()),
		})
	}
	return validationErrorsResponse{
		Success: false,
		Errors:  errors,
	}
}
