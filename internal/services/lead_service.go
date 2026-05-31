package services

import (
	"fmt"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/repositories"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LeadService struct {
	repo  *repositories.LeadRepository
	redis *redis.Client
}

func NewLeadService(repo *repositories.LeadRepository, redisClient *redis.Client) *LeadService {
	return &LeadService{
		repo:  repo,
		redis: redisClient,
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

	cachedLead, err := s.getCachedLead(id.String())

	if err == nil && cachedLead != nil {
		return cachedLead, nil
	}

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.setCachedLead(lead)

	return lead, nil
}

func (s *LeadService) GetAllLeads(status *models.LeadStatus,email *string,) ([]models.Lead, error) {
	return s.repo.GetAll(status , email)
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

	s.deleteCachedLead(id.String())

	return lead, nil
}

func (s *LeadService) UpdateLeadStatus(id uuid.UUID, newStatus models.LeadStatus) (*models.Lead, error) {

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !isValidTransition(lead.Status, newStatus) {

		return nil, fmt.Errorf(
			"%w: %s -> %s",
			ErrInvalidStatusTransition,
			lead.Status,
			newStatus,
		)
	}

	lead.Status = newStatus

	err = s.repo.Update(lead)
	if err != nil {
		return nil, err
	}

	s.deleteCachedLead(id.String())

	return lead, nil
}

func (s *LeadService) DeleteLead(id uuid.UUID) error {
	s.deleteCachedLead(id.String())
	return s.repo.Delete(id)
}

// bulk requests
func (s *LeadService) BulkCreateLeads(req dto.BulkCreateLeadRequest) dto.BulkResponse {

	response := dto.BulkResponse{
		Total: len(req.Leads),
	}

	validate := binding.Validator.Engine().(*validator.Validate)

	for index, leadReq := range req.Leads {

		err := validate.Struct(leadReq)
		if err != nil {

			response.Failed++

			response.Results = append(
				response.Results,
				dto.BulkResult{
					Index:   index,
					Success: false,
					Error:   err.Error(),
				},
			)

			continue
		}

		lead, err := s.CreateLead(&leadReq)
		if err != nil {

			response.Failed++

			response.Results = append(
				response.Results,
				dto.BulkResult{
					Index:   index,
					Success: false,
					Error:   err.Error(),
				},
			)

			continue
		}

		leadResponse := dto.LeadResponse{
			ID:        lead.ID,
			Name:      lead.Name,
			Email:     lead.Email,
			Phone:     lead.Phone,
			Source:    lead.Source,
			Status:    lead.Status,
			CreatedAt: lead.CreatedAt,
			UpdatedAt: lead.UpdatedAt,
		}

		response.Successful++

		response.Results = append(
			response.Results,
			dto.BulkResult{
				Index:   index,
				Success: true,
				Lead:    &leadResponse,
			},
		)
	}

	return response
}

func (s *LeadService) BulkUpdateLeads(req dto.BulkUpdateLeadRequest) dto.BulkResponse {

	response := dto.BulkResponse{Total: len(req.Leads)}

	for index, leadReq := range req.Leads {

		id, err := uuid.Parse(leadReq.ID)
		if err != nil {

			response.Failed++

			response.Results = append(
				response.Results,
				dto.BulkResult{
					Index:   index,
					Success: false,
					Error:   fmt.Sprintf("invalid UUID for lead at index %d: must be a valid UUID (e.g. 123e4567-e89b-12d3-a456-426614174000)", index),
				},
			)

			continue
		}

		lead, err := s.UpdateLead(id, leadReq.UpdateLeadRequest)

		if err != nil {

			response.Failed++

			response.Results = append(
				response.Results,
				dto.BulkResult{
					Index:   index,
					Success: false,
					Error:   err.Error(),
				},
			)

			continue
		}

		leadResponse := dto.LeadResponse{
			ID:        lead.ID,
			Name:      lead.Name,
			Email:     lead.Email,
			Phone:     lead.Phone,
			Source:    lead.Source,
			Status:    lead.Status,
			CreatedAt: lead.CreatedAt,
			UpdatedAt: lead.UpdatedAt,
		}

		response.Successful++

		response.Results = append(
			response.Results,
			dto.BulkResult{
				Index:   index,
				Success: true,
				Lead:    &leadResponse,
			},
		)
	}

	return response
}
