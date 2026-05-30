package main

import (
	"net/http"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		panic(err)
	}

	_ = db

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "healthy",
		})
	})

	if err := router.Run(":" + cfg.Port.Port); err != nil {
		panic(err)
	}
}
