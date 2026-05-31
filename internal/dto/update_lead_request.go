package dto

type UpdateLeadRequest struct {
	Name   *string `json:"name,omitempty" binding:"omitempty,fullname,min=2,max=100"`
	Email  *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone  *string `json:"phone,omitempty" binding:"omitempty,phone"`
	Source *string `json:"source,omitempty"`
}