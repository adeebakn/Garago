package app

import (
	"testing"
	"time"
)

func TestLoadConfigUsesDefaults(t *testing.T) {
	setRequiredConfig(t)
	t.Setenv("JWT_EXPIRY_HOURS", "")
	t.Setenv("REDIS_URL", "")
	t.Setenv("FRONTEND_URL", "")
	t.Setenv("SERVER_ADDRESS", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg.sessionTTL != 24*time.Hour {
		t.Fatalf("sessionTTL = %v, want 24h", cfg.sessionTTL)
	}
	if cfg.redisURL != "redis://localhost:6379/0" {
		t.Fatalf("redisURL = %q", cfg.redisURL)
	}
	if cfg.serverAddr != ":8080" {
		t.Fatalf("serverAddr = %q", cfg.serverAddr)
	}
}

func TestLoadConfigRejectsShortJWTSecret(t *testing.T) {
	setRequiredConfig(t)
	t.Setenv("JWT_SECRET", "too-short")

	if _, err := loadConfig(); err == nil {
		t.Fatal("loadConfig() expected an error for a short JWT secret")
	}
}

func TestLoadConfigRejectsInvalidSMTPPort(t *testing.T) {
	setRequiredConfig(t)
	t.Setenv("SMTP_PORT", "70000")

	if _, err := loadConfig(); err == nil {
		t.Fatal("loadConfig() expected an error for an invalid SMTP port")
	}
}

func setRequiredConfig(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("SMTP_USERNAME", "user")
	t.Setenv("SMTP_PASSWORD", "password")
	t.Setenv("SMTP_FROM", "sender@example.com")
	t.Setenv("SMTP_FROM_NAME", "Garago")
}
