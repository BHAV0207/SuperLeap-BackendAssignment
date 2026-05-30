package dto

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
