package dto

import (
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/google/uuid"
)

type LeadResponse struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Phone     *string           `json:"phone,omitempty"`
	Source    *string           `json:"source,omitempty"`
	Status    models.LeadStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
