package config

import (
	"os"
)

type Config struct {
	NotionAPIKey     string
	NotionAPIVersion string
	NotionAPIURL     string
}

func Load() *Config {
	return &Config{
		NotionAPIKey:     getEnv("NOTION_API_KEY", ""),
		NotionAPIVersion: getEnv("NOTION_VERSION", "2022-06-28"),
		NotionAPIURL:     getEnv("NOTION_API_URL", "https://api.notion.com/v1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
