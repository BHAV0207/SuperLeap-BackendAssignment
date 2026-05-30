package routes

import (
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterLeadRoutes(
	router *gin.Engine,
	handler *handlers.LeadHandler,
) {

	leads := router.Group("/leads")

	{
		leads.POST("", handler.CreateLead)
		leads.GET("", handler.GetAllLeads)
		leads.GET("/:id", handler.GetLeadByID)
		leads.PUT("/:id", handler.UpdateLead)
		leads.DELETE("/:id", handler.DeleteLead)
		leads.PATCH("/:id/status", handler.UpdateLeadStatus)
		// Bulk operations
		leads.POST("/bulk", handler.BulkCreateLeads)
		leads.PUT("/bulk", handler.BulkUpdateLeads)
	}
}
