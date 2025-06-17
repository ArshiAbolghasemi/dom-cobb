package api

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
