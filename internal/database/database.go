package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	connectAttempts = 10
	connectBackoff  = 3 * time.Second
)

func ConnectDatabase(cfg *config.AppConfig) (*gorm.DB, error) {
	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// on kubernetes the app pod can come up before postgres is ready, so
	// retry for ~30s instead of crashing straight away
	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= connectAttempts; attempt++ {
		db, err = open(dsn)
		if err == nil {
			break
		}
		slog.Warn("database not reachable yet", "attempt", attempt, "error", err)
		time.Sleep(connectBackoff)
	}
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// adds a span for every query so we can see db time inside a request trace
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return nil, fmt.Errorf("failed to enable db tracing: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// keep the pool small, postgres has a low connection limit and we may run several replicas
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
