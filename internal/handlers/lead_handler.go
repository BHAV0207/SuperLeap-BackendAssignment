package handlers

import (
	"errors"
	"net/http"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/services"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeadHandler struct {
	service *services.LeadService
}

func NewLeadHandler(service *services.LeadService) *LeadHandler {
	return &LeadHandler{
		service: service,
	}
}

func (h *LeadHandler) CreateLead(c *gin.Context) {
	var req dto.CreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		utils.ValidationError(c, err.Error())
		return
	}

	lead, err := h.service.CreateLead(&req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.LeadResponse{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Phone:     lead.Phone,
		Source:    lead.Source,
		Status:    lead.Status,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: lead.UpdatedAt,
	}

	utils.Success(c, http.StatusCreated, response)
}

func (h *LeadHandler) GetLeadByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {

		utils.ValidationError(c, "invalid UUID format")
		return
	}

	lead, err := h.service.GetLeadByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "lead not found",
		})
		return
	}

	response := dto.LeadResponse{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Phone:     lead.Phone,
		Source:    lead.Source,
		Status:    lead.Status,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: lead.UpdatedAt,
	}

	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) GetAllLeads(c *gin.Context) {
	statusQuery := c.Query("status")

	var status *models.LeadStatus

	if statusQuery != "" {
		s := models.LeadStatus(statusQuery)
		status = &s
	}

	leads, err := h.service.GetAllLeads(status)
	if err != nil {

		utils.Error(c, http.StatusInternalServerError, err.Error())

		return
	}

	var response []dto.LeadResponse

	for _, lead := range leads {
		response = append(response, dto.LeadResponse{
			ID:        lead.ID,
			Name:      lead.Name,
			Email:     lead.Email,
			Phone:     lead.Phone,
			Source:    lead.Source,
			Status:    lead.Status,
			CreatedAt: lead.CreatedAt,
			UpdatedAt: lead.UpdatedAt,
		})
	}

	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) UpdateLead(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {

		utils.ValidationError(c, "invalid UUID format")
		return
	}

	var req dto.UpdateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
	
		utils.ValidationError(c, err.Error())
		return
	}

	lead, err := h.service.UpdateLead(id, req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "lead not found",
		})
		return
	} else if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.LeadResponse{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Phone:     lead.Phone,
		Source:    lead.Source,
		Status:    lead.Status,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: lead.UpdatedAt,
	}

	utils.Success(c, http.StatusOK, response)

}

func (h *LeadHandler) UpdateLeadStatus(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
	
		utils.ValidationError(c, "invalid UUID format")
		return
	}

	var req dto.UpdateLeadStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		utils.ValidationError(c, err.Error())
		return
	}

	lead, err := h.service.UpdateLeadStatus(id, req.Status)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "lead not found",
		})
		return
	} else if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.LeadResponse{
		ID:        lead.ID,
		Name:      lead.Name,
		Email:     lead.Email,
		Phone:     lead.Phone,
		Source:    lead.Source,
		Status:    lead.Status,
		CreatedAt: lead.CreatedAt,
		UpdatedAt: lead.UpdatedAt,
	}
	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) DeleteLead(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationError(c, "invalid UUID format")
		return
	}

	err = h.service.DeleteLead(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "lead not found",
		})
		return
	} else if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "lead deleted successfully")
}
