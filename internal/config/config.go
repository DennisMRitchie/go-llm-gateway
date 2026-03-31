package config

import (
	"os"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port           string
	OllamaURL      string
	OpenAIAPIKey   string
	RedisAddr      string
	DefaultModel   string
	RateLimitRPS   int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8082"),
		OllamaURL:    getEnv("OLLAMA_URL", "http://ollama:11434"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		RedisAddr:    getEnv("REDIS_ADDR", "redis:6379"),
		DefaultModel: getEnv("DEFAULT_MODEL", "qwen3"),
		RateLimitRPS: 100,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
