package config

import (
	"os"
	"strings"
)

type Config struct {
	AppPort            string
	CORSAllowedOrigins []string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	GooseDriver       string
	GooseDBString     string
	GooseMigrationDir string

	ResendAPIKey string
	JWTSecretKey string
}

func Load() Config {
	return Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		CORSAllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),

		GooseDriver:       os.Getenv("GOOSE_DRIVER"),
		GooseDBString:     os.Getenv("GOOSE_DBSTRING"),
		GooseMigrationDir: os.Getenv("GOOSE_MIGRATION_DIR"),

		ResendAPIKey: os.Getenv("RESEND_API_KEY"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvList(key string) []string {
	value := os.Getenv(key)

	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")

	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			result = append(result, origin)
		}
	}

	return result
}
