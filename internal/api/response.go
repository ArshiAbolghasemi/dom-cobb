package api

import (
	"net/http"

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

func RespondInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "Internal Server Error",
		Message: err.Error(),
	})
}

func RespondSuccess(c *gin.Context, statusCode int, msg string, data any) {
	c.JSON(statusCode, SuccessResponse{
		Message: msg,
		Data: data,
	})
}
