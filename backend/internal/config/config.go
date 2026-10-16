package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Telegram
	TelegramToken string

	// LLM
	LLMBaseURL string
	LLMAPIKey  string
	LLMModel   string

	// Server
	ServerPort string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "expenses"),
		DBPassword: getEnv("DB_PASSWORD", "expenses"),
		DBName:     getEnv("DB_NAME", "expenses"),

		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),

		LLMBaseURL: os.Getenv("LLM_BASE_URL"),
		LLMAPIKey:  os.Getenv("LLM_API_KEY"),
		LLMModel:   os.Getenv("LLM_MODEL"),

		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.LLMBaseURL == "" {
		return nil, fmt.Errorf("LLM_BASE_URL is required")
	}
	if cfg.LLMAPIKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY is required")
	}
	if cfg.LLMModel == "" {
		return nil, fmt.Errorf("LLM_MODEL is required")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
