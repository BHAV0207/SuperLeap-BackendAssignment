package dto

type UpdateLeadRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone  *string `json:"phone,omitempty"`
	Source *string `json:"source,omitempty"`
}

