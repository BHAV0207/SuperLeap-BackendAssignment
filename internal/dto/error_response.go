package dto

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}