package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	databaseURL  string
	redisURL     string
	jwtSecret    string
	sessionTTL   time.Duration
	cookieSecure bool
	frontendURL  string
	serverAddr   string
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string
	smtpFromName string
}

func loadConfig() (config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return config{}, fmt.Errorf("load .env: %w", err)
	}

	hours, err := positiveIntEnv("JWT_EXPIRY_HOURS", 24)
	if err != nil {
		return config{}, err
	}
	port, err := positiveIntEnv("SMTP_PORT", 0)
	if err != nil || port > 65535 {
		return config{}, errors.New("SMTP_PORT must be a valid port number")
	}

	cfg := config{
		databaseURL:  strings.TrimSpace(os.Getenv("DATABASE_URL")),
		redisURL:     envOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		jwtSecret:    os.Getenv("JWT_SECRET"),
		sessionTTL:   time.Duration(hours) * time.Hour,
		cookieSecure: strings.EqualFold(strings.TrimSpace(os.Getenv("COOKIE_SECURE")), "true"),
		frontendURL:  envOrDefault("FRONTEND_URL", "http://localhost:5173"),
		serverAddr:   envOrDefault("SERVER_ADDRESS", ":8080"),
		smtpHost:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		smtpPort:     port,
		smtpUsername: strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		smtpPassword: strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		smtpFrom:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
		smtpFromName: strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}

	if cfg.databaseURL == "" {
		return config{}, errors.New("DATABASE_URL must be configured")
	}
	if len(cfg.jwtSecret) < 32 {
		return config{}, errors.New("JWT_SECRET must be at least 32 characters")
	}
	for name, value := range map[string]string{
		"SMTP_HOST": cfg.smtpHost, "SMTP_USERNAME": cfg.smtpUsername,
		"SMTP_PASSWORD": cfg.smtpPassword, "SMTP_FROM": cfg.smtpFrom,
		"SMTP_FROM_NAME": cfg.smtpFromName,
	} {
		if value == "" {
			return config{}, fmt.Errorf("%s must be configured", name)
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
