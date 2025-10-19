package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	LLMApiKey    string
	LLMBaseURL   string
	JWTSecret    string
	JWTLifetime  time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		LLMApiKey:    getEnv("LLM_API_KEY", ""),
		LLMBaseURL:   getEnv("LLM_BASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", "default_secret"),
		JWTLifetime:  time.Hour * 72,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("ERROR: Environment variable %s is not set", key)
	}
	return fallback
}