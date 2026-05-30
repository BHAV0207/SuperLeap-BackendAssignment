package dto

import "github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"


type UpdateLeadStatusRequest struct {
	Status models.LeadStatus `json:"status" binding:"required"`
}
