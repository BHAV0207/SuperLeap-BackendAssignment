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

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusCreated, response)
}

func (h *LeadHandler) GetLeadByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {

		utils.ValidationError(c, "invalid UUID format")
		return
	}

	lead, err := h.service.GetLeadByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {

		utils.Error(c, http.StatusNotFound, "lead not found")

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
		utils.Error(c, http.StatusNotFound, "lead not found")
		return
	} else if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.ToLeadResponse(lead)

	utils.Success(c, http.StatusOK, response)

}

func (h *LeadHandler) UpdateLeadStatus(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
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

		utils.Error(c, http.StatusNotFound, "lead not found")
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationError(c, "invalid UUID format")
		return
	}

	err = h.service.DeleteLead(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Error(c, http.StatusNotFound, "lead not found")
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

	if c.Request.Body == nil {
		utils.ValidationError(c, "request body cannot be empty")
		return
	}

	var req dto.BulkCreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
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

		utils.ValidationError(c, err.Error())
		return
	}

	if len(req.Leads) == 0 {
		utils.ValidationError(c, "leads array cannot be empty")
		return
	}

	response := h.service.BulkUpdateLeads(req)

	utils.Success(c, http.StatusOK, response)
}
