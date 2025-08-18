package config

import (
	"log"
	"os"
)

type Config struct {
	DBUser           string
	DBPassword       string
	DBHost           string
	DBPort           string
	DBName           string
	HttpPort         string
	ExternalAPIToken string
	ExternalAPIURL   string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBUser:           getEnv("DB_USER", ""),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "26257"), // Cockroach default port
		DBName:           getEnv("DB_NAME", "defaultdb"),
		HttpPort:         getEnv("HTTP_PORT", "8080"),
		ExternalAPIToken: getEnv("EXTERNAL_API_TOKEN", ""),
		ExternalAPIURL:   getEnv("EXTERNAL_API_URL", ""),
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" {
		log.Fatal("❌ DB_USER y DB_PASSWORD son obligatorios")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
