package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	SessionTTL   time.Duration
	CookieSecure bool
	FrontendURL  string
	ServerAddr   string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	hours, err := positiveIntEnv("JWT_EXPIRY_HOURS", 24)
	if err != nil {
		return Config{}, err
	}
	port, err := positiveIntEnv("SMTP_PORT", 0)
	if err != nil || port > 65535 {
		return Config{}, errors.New("SMTP_PORT must be a valid port number")
	}

	cfg := Config{
		DatabaseURL:  strings.TrimSpace(os.Getenv("DATABASE_URL")),
		RedisURL:     envOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		SessionTTL:   time.Duration(hours) * time.Hour,
		CookieSecure: strings.EqualFold(strings.TrimSpace(os.Getenv("COOKIE_SECURE")), "true"),
		FrontendURL:  envOrDefault("FRONTEND_URL", "http://localhost:5173"),
		ServerAddr:   envOrDefault("SERVER_ADDRESS", ":8080"),
		SMTPHost:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:     port,
		SMTPUsername: strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		SMTPPassword: strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		SMTPFrom:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
		SMTPFromName: strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL must be configured")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must be at least 32 characters")
	}
	for name, value := range map[string]string{
		"SMTP_HOST": cfg.SMTPHost, "SMTP_USERNAME": cfg.SMTPUsername,
		"SMTP_PASSWORD": cfg.SMTPPassword, "SMTP_FROM": cfg.SMTPFrom,
		"SMTP_FROM_NAME": cfg.SMTPFromName,
	} {
		if value == "" {
			return Config{}, fmt.Errorf("%s must be configured", name)
		}
	}
	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func positiveIntEnv(name string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive number", name)
	}
	return parsed, nil
}
