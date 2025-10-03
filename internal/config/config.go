package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	DatabaseURL      string
	GroqAPIKey       string
	Port             string
	WebhookURL       string
	DefaultTimezone  string
	DefaultLanguage  string
	GroqConfig       GroqConfig
}

type GroqConfig struct {
	Model       string
	Temperature float64
	MaxTokens   int
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable"),
		GroqAPIKey:       getEnv("GROQ_API_KEY", ""),
		Port:             getEnv("PORT", "8080"),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		DefaultTimezone:  getEnv("DEFAULT_TIMEZONE", "Asia/Jakarta"),
		DefaultLanguage:  getEnv("DEFAULT_LANGUAGE", "id"),
		GroqConfig: GroqConfig{
			Model:       getEnv("GROQ_MODEL", "llama-3.1-70b-versatile"),
			Temperature: getFloatEnv("GROQ_TEMPERATURE", 0.3),
			MaxTokens:   getIntEnv("GROQ_MAX_TOKENS", 1000),
		},
	}

	// Validate required environment variables
	if config.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if config.GroqAPIKey == "" {
		log.Fatal("GROQ_API_KEY is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}