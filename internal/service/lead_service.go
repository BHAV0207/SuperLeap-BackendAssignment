package service

import (
	"errors"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/repositories"
	"github.com/google/uuid"
)

type LeadService struct {
	repo *repositories.LeadRepository
}

func NewLeadService(repo *repositories.LeadRepository) *LeadService {
	return &LeadService{
		repo: repo,
	}
}

func (s *LeadService) CreateLead(req *dto.CreateLeadRequest) (*models.Lead, error) {
	lead := models.Lead{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Source: req.Source,
	}

	err := s.repo.Create(&lead)
	if err != nil {
		return nil, err
	}

	return &lead, nil
}

func (s *LeadService) GetLeadByID(id uuid.UUID) (*models.Lead, error) {
	return s.repo.GetByID(id)
}

func (s *LeadService) GetAllLeads(status *models.LeadStatus) ([]models.Lead, error) {
	return s.repo.GetAll(status)
}

func (s *LeadService) UpdateLead(id uuid.UUID, req dto.UpdateLeadRequest) (*models.Lead, error) {

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		lead.Name = *req.Name
	}

	if req.Email != nil {
		lead.Email = *req.Email
	}

	if req.Phone != nil {
		lead.Phone = req.Phone
	}

	if req.Source != nil {
		lead.Source = req.Source
	}

	err = s.repo.Update(lead)
	if err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *LeadService) UpdateLeadStatus(
	id uuid.UUID,
	newStatus models.LeadStatus,
) (*models.Lead, error) {

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !isValidTransition(lead.Status, newStatus) {
		return nil, errors.New(
			"invalid status transition",
		)
	}

	lead.Status = newStatus

	err = s.repo.Update(lead)
	if err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *LeadService) DeleteLead(id uuid.UUID) error {
	return s.repo.Delete(id)
}