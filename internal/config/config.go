package config

import (
	"os"
)

type Config struct {
	ServerPort     string
	GitHubSecret   string
	LarkWebhookURL string
}

func New() *Config {
	return &Config{
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		GitHubSecret:   getEnvOrDefault("GITHUB_SECRET", ""),
		LarkWebhookURL: getEnvOrDefault("LARK_WEBHOOK_URL", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
