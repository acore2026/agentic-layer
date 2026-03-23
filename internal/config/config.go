package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LLMProvider  string // "gemini" or "kimi"
	GeminiAPIKey string
	KimiAPIKey   string
	KimiBaseURL  string
	KimiModel    string
	ACRFURL      string
	IGWURL       string
	Port         string
}

// Load loads the configuration from .env file and environment variables.
func Load() (*Config, error) {
	// Try loading .env file but don't fail if it's missing (might be in CI/CD or Docker)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		LLMProvider:  getEnv("AGENTIC_LLM_PROVIDER", "gemini"),
		GeminiAPIKey: os.Getenv("AGENTIC_GEMINI_API_KEY"),
		KimiAPIKey:   os.Getenv("AGENTIC_KIMI_API_KEY"),
		KimiBaseURL:  getEnv("AGENTIC_KIMI_BASE_URL", "https://api.moonshot.cn/v1"),
		KimiModel:    getEnv("AGENTIC_KIMI_MODEL", "kimi-k2.5"),
		ACRFURL:      getEnv("AGENTIC_ACRF_URL", "http://localhost:18080"),
		IGWURL:       getEnv("AGENTIC_IGW_URL", "http://localhost:18081"),
		Port:         getEnv("AGENTIC_AAIHF_PORT", "18082"),
	}

	// Validation
	if cfg.LLMProvider == "gemini" && cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("required environment variable AGENTIC_GEMINI_API_KEY is missing for provider 'gemini'")
	}
	if cfg.LLMProvider == "kimi" && cfg.KimiAPIKey == "" {
		return nil, fmt.Errorf("required environment variable AGENTIC_KIMI_API_KEY is missing for provider 'kimi'")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
