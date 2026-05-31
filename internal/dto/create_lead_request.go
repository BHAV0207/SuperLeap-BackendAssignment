package dto

type CreateLeadRequest struct {
	Name   string  `json:"name" binding:"required,fullname,min=2,max=100"`
	Email  string  `json:"email" binding:"required,email"`
	Phone  *string `json:"phone,omitempty" binding:"omitempty,phone"`
	Source *string `json:"source,omitempty"`
}