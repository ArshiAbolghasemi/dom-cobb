package api

import "net/http"

type APIError struct {
	StatusCode int
	Error      string
	Message    string
}

func BadRequestError(code, message string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Error:      code,
		Message:    message,
	}
}

func ConflictError(code, message string) *APIError {
	return &APIError{
		StatusCode: http.StatusConflict,
		Error:      code,
		Message:    message,
	}
}

func NotFoundError(code, message string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Error:      code,
		Message:    message,
	}
}

func InternalServerError(code, message string) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Error:      code,
		Message:    message,
	}
}

func OKError(code, message string) *APIError {
	return &APIError{
		StatusCode: http.StatusOK,
		Error:      code,
		Message:    message,
	}
}
