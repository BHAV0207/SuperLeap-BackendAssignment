package config

import "testing"

func setRequired(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "app")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "leads")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
}

func TestNewAppConfigFromEnv(t *testing.T) {
	setRequired(t)

	cfg, err := NewAppConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.SSLMode != "disable" {
		t.Errorf("default sslmode = %q, want disable", cfg.Database.SSLMode)
	}
	if cfg.Telemetry.ServiceName != "lead-crm" {
		t.Errorf("default service name = %q, want lead-crm", cfg.Telemetry.ServiceName)
	}
	if cfg.Telemetry.OTLPEndpoint != "" {
		t.Errorf("tracing should be off by default, got endpoint %q", cfg.Telemetry.OTLPEndpoint)
	}
}

func TestNewAppConfigMissingDB(t *testing.T) {
	setRequired(t)
	t.Setenv("DB_PASSWORD", "")

	if _, err := NewAppConfig(); err == nil {
		t.Fatal("expected error when DB_PASSWORD is missing")
	}
}
