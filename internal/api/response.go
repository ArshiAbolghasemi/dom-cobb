package api

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func RespondAPIError(c *gin.Context, apiErr *APIError) {
	c.JSON(apiErr.StatusCode, ErrorResponse{
		Error:   apiErr.Error,
		Message: apiErr.Message,
	})
}

func RespondSuccess(c *gin.Context, statusCode int, msg string, data any) {
	c.JSON(statusCode, SuccessResponse{
		Message: msg,
		Data:    data,
	})
}

type PaginationResponse struct {
	Page       uint `json:"page"`
	Size       uint `json:"size"`
	Total      uint `json:"total"`
	TotalPages uint `json:"total_pages"`
}
