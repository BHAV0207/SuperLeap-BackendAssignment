package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/services"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/utils"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/validators"
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
		utils.ValidationError(c, utils.ParseValidationError(err))
		return
	}

	lead, err := h.service.CreateLead(&req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusCreated, response)
}

func (h *LeadHandler) GetLeadByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationError(c, "invalid lead ID: must be a valid UUID (e.g. 123e4567-e89b-12d3-a456-426614174000)")
		return
	}

	lead, err := h.service.GetLeadByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Error(c, http.StatusNotFound, fmt.Sprintf("lead with ID '%s' not found", c.Param("id")))
		return
	}

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) GetAllLeads(c *gin.Context) {
	statusQuery := c.Query("status")
	emailQuery := c.Query("email")

	validStatuses := []string{"NEW", "CONTACTED", "QUALIFIED", "CONVERTED", "LOST"}

	var status *models.LeadStatus
	var email *string

	if statusQuery != "" {
		s := models.LeadStatus(statusQuery)

		if !validators.IsValidStatus(string(s)) {
			utils.ValidationError(c, fmt.Sprintf(
				"invalid status value '%s': must be one of %s",
				statusQuery,
				strings.Join(validStatuses, ", "),
			))
			return
		}
		status = &s
	}

	if emailQuery != "" {
		email = &emailQuery
	}

	leads, err := h.service.GetAllLeads(
		status,
		email,
	)

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response []dto.LeadResponse

	for _, lead := range leads {
		response = append(
			response,
			utils.ToLeadResponse(&lead),
		)
	}

	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) UpdateLead(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationError(c, "invalid lead ID: must be a valid UUID (e.g. 123e4567-e89b-12d3-a456-426614174000)")
		return
	}

	var req dto.UpdateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationError(err))
		return
	}

	lead, err := h.service.UpdateLead(id, req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Error(c, http.StatusNotFound, fmt.Sprintf("lead with ID '%s' not found", idParam))
		return
	} else if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusOK, response)

}

func (h *LeadHandler) UpdateLeadStatus(c *gin.Context) {

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationError(c, "invalid lead ID: must be a valid UUID (e.g. 123e4567-e89b-12d3-a456-426614174000)")
		return
	}

	var req dto.UpdateLeadStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationError(err))
		return
	}

	lead, err := h.service.UpdateLeadStatus(id, req.Status)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Error(c, http.StatusNotFound, fmt.Sprintf("lead with ID '%s' not found", idParam))
		return
	}

	if errors.Is(err, services.ErrInvalidStatusTransition) {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusOK, response)
}

func (h *LeadHandler) DeleteLead(c *gin.Context) {

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationError(c, "invalid lead ID: must be a valid UUID (e.g. 123e4567-e89b-12d3-a456-426614174000)")
		return
	}

	err = h.service.DeleteLead(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Error(c, http.StatusNotFound, fmt.Sprintf("lead with ID '%s' not found", idParam))
		return
	}

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "lead deleted successfully")
}

// bulk operations
func (h *LeadHandler) BulkCreateLeads(c *gin.Context) {
	var req dto.BulkCreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationError(err))
		return
	}

	if len(req.Leads) == 0 {
		utils.ValidationError(c, "leads array cannot be empty")
		return
	}

	response := h.service.BulkCreateLeads(req)
	utils.Success(c, http.StatusCreated, response)
}

func (h *LeadHandler) BulkUpdateLeads(c *gin.Context) {

	var req dto.BulkUpdateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationError(err))
		return
	}

	if len(req.Leads) == 0 {
		utils.ValidationError(c, "leads array cannot be empty")
		return
	}

	response := h.service.BulkUpdateLeads(req)

	utils.Success(c, http.StatusOK, response)
}
