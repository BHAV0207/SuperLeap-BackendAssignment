package utils

import (
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
)

func ToLeadResponse(lead *models.Lead) dto.LeadResponse {
	return dto.LeadResponse{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Phone:     lead.Phone,
		Source:    lead.Source,
		Status:    lead.Status,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: lead.UpdatedAt,
	}
}