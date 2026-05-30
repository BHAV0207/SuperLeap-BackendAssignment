package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/database"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/handlers"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/repositories"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/routes"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/services"
)

func main() {

	// Load configuration
	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect database
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.Lead{})
	if err != nil {
		log.Fatal(err)
	}

	err = database.SeedData(db)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize dependencies
	leadRepository := repositories.NewLeadRepository(db)

	leadService := services.NewLeadService(
		leadRepository,
	)

	leadHandler := handlers.NewLeadHandler(
		leadService,
	)

	// Create router
	router := gin.Default()

	// Health check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "healthy",
		})
	})

	// Register routes
	routes.RegisterLeadRoutes(
		router,
		leadHandler,
	)

	// Start server
	err = router.Run(":" + cfg.Port.Port)
	if err != nil {
		log.Fatal(err)
	}
}
