package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port      PortConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Telemetry TelemetryConfig
}

type PortConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host string
	Port string
}

type TelemetryConfig struct {
	ServiceName string
	Environment string
	// empty means tracing is disabled
	OTLPEndpoint string
	LogLevel     string
}

func NewAppConfig() (*AppConfig, error) {

	// .env is only for local runs. In docker / kubernetes the values come
	// from real environment variables, so a missing file is not an error.
	_ = godotenv.Load()

	portConfig, err := newPortConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading port configuration: %w", err)
	}

	dbConfig, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading database configuration: %w", err)
	}
	redisConfig, err := newRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading Redis configuration: %w", err)
	}

	return &AppConfig{
		Port:      *portConfig,
		Database:  *dbConfig,
		Redis:     *redisConfig,
		Telemetry: newTelemetryConfig(),
	}, nil
}

func newPortConfig() (*PortConfig, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable is required")
	}

	return &PortConfig{Port: port}, nil
}

func newDatabaseConfig() (*DatabaseConfig, error) {

	dbConfig := &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Name == "" {
		return nil, fmt.Errorf("all database configuration fields are required")
	}

	return dbConfig, nil
}

func newRedisConfig() (*RedisConfig, error) {

	redisConfig := &RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}

	if redisConfig.Host == "" || redisConfig.Port == "" {
		return nil, fmt.Errorf("all Redis configuration fields are required")
	}

	return redisConfig, nil
}

func newTelemetryConfig() TelemetryConfig {
	return TelemetryConfig{
		ServiceName:  getEnv("OTEL_SERVICE_NAME", "lead-crm"),
		Environment:  getEnv("APP_ENV", "local"),
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

//generally the helper functions job is to pass the error to the caller and then the caller decides how to handel THE ERROR .
