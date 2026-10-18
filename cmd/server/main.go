package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/database"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/handlers"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/middleware"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/observability"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/repositories"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/routes"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/services"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/validators"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {

	// Load configuration
	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := observability.NewLogger(cfg.Telemetry.LogLevel, cfg.Telemetry.ServiceName, cfg.Telemetry.Environment)

	shutdownTracing, err := observability.InitTracing(
		context.Background(),
		cfg.Telemetry.ServiceName,
		cfg.Telemetry.Environment,
		cfg.Telemetry.OTLPEndpoint,
	)
	if err != nil {
		// tracing is a debugging aid, it should never stop the api from serving
		slog.Error("tracing disabled, failed to init", "error", err)
		shutdownTracing = func(context.Context) error { return nil }
	}

	// Connect database
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		fatal("failed to connect database", err)
	}

	//redis connection
	redisClient := database.ConnectRedis(cfg)

	// Run migrations
	err = db.AutoMigrate(&models.Lead{})
	if err != nil {
		fatal("migration failed", err)
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
	router := gin.New()
	router.Use(
		gin.Recovery(),
		otelgin.Middleware(cfg.Telemetry.ServiceName),
		middleware.RequestID(),
		middleware.Metrics(),
		middleware.Logger(logger),
	)

	registerOpsRoutes(router, db, redisClient)

	// Register routes
	routes.RegisterLeadRoutes(
		router,
		leadHandler,
	)

	srv := &http.Server{
		Addr:              ":" + cfg.Port.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server
	go func() {
		slog.Info("server starting", "port", cfg.Port.Port, "version", version.Version, "commit", version.Commit)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fatal("server failed", err)
		}
	}()

	// kubernetes sends SIGTERM before killing the pod. Stop taking new
	// requests, let in-flight ones finish, then close connections.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	if err := shutdownTracing(ctx); err != nil {
		slog.Error("tracing shutdown failed", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	if redisClient != nil {
		redisClient.Close()
	}

	slog.Info("shutdown complete")
}

// probes, metrics and version info
func registerOpsRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {

	// liveness: the process is up and serving. It does not check postgres on
	// purpose, a db outage should not make kubernetes restart every pod.
	liveness := func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "healthy",
		})
	}
	router.GET("/healthz", liveness)
	// old endpoint, kept so existing clients do not break
	router.GET("/health", liveness)

	// readiness: can this pod serve traffic right now. Postgres is required.
	// Redis is only reported, the app works without it (just slower).
	router.GET("/readyz", func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()

		checks := gin.H{}
		ready := true

		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.PingContext(c)
		}
		if err != nil {
			ready = false
			checks["postgres"] = "down: " + err.Error()
			observability.DependencyUp.WithLabelValues("postgres").Set(0)
		} else {
			checks["postgres"] = "up"
			observability.DependencyUp.WithLabelValues("postgres").Set(1)
		}

		if redisClient != nil {
			if err := redisClient.Ping(c).Err(); err != nil {
				checks["redis"] = "degraded: " + err.Error()
				observability.DependencyUp.WithLabelValues("redis").Set(0)
			} else {
				checks["redis"] = "up"
				observability.DependencyUp.WithLabelValues("redis").Set(1)
			}
		}

		status := http.StatusOK
		if !ready {
			status = http.StatusServiceUnavailable
		}

		ctx.JSON(status, gin.H{
			"success": ready,
			"checks":  checks,
		})
	})

	router.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, version.Get())
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
