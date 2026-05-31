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
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/validators"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	//redis connection
	redisClient := database.ConnectRedis(cfg)

	// Run migrations
	err = db.AutoMigrate(&models.Lead{})
	if err != nil {
		log.Fatal(err)
	}


	// Initialize dependencies
	leadRepository := repositories.NewLeadRepository(db)

	leadService := services.NewLeadService(
		leadRepository,
		redisClient,
	)

	leadHandler := handlers.NewLeadHandler(
		leadService,
	)

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phone", validators.ValidatePhone)
		v.RegisterValidation("fullname", validators.ValidateFullName)
	}

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
