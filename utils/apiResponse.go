package utils

import "net/http"

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type Error struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Success bool `json:"success"`
	Error
}

func SuccessResponse(data interface{}) successResponse {
	return successResponse{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(code uint, message string) errorResponse {
	return errorResponse{
		Success: false,
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}

func BadRequestResponse(message string) errorResponse {
	return ErrorResponse(http.StatusBadRequest, message)
}

func UnauthorizedResponse(message string) errorResponse {
	return ErrorResponse(http.StatusUnauthorized, message)
}

func ForbiddenResponse(message string) errorResponse {
	return ErrorResponse(http.StatusForbidden, message)
}

func NotFoundResponse(message string) errorResponse {
	return ErrorResponse(http.StatusNotFound, message)
}

func ConflictResponse(message string) errorResponse {
	return ErrorResponse(http.StatusConflict, message)
}

func InternalServerErrorResponse(message string) errorResponse {
	return ErrorResponse(http.StatusInternalServerError, message)
}
