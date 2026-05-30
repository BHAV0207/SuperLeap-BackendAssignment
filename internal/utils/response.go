package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/dto"
)

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, dto.SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, err interface{}) {
	c.JSON(status, dto.ErrorResponse{
		Success: false,
		Error:   err,
	})
}

func ValidationError(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Success: false,
		Error:   err,
	})
}
