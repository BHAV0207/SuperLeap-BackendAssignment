package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port     PortConfig
	Database DatabaseConfig
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
}

func NewAppConfig() (*AppConfig, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	portConfig, err := newPortConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading port configuration: %w", err)
	}

	dbConfig, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading database configuration: %w", err)
	}

	return &AppConfig{
		Port:     *portConfig,
		Database: *dbConfig,
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
	}

	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Name == "" {
		return nil, fmt.Errorf("all database configuration fields are required")
	}

	return dbConfig, nil
}

//generally the helper functions job is to pass the error to the caller and then the caller decides how to handel THE ERROR .
