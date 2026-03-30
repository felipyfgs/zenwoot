package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	ServerHost     string
	APIKey         string
	JWTSecret      string
	LogLevel       string
	Environment    string
	AllowedOrigins string

	DatabaseURL string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool

	NatsURL string

	GlobalWebhookURL string

	// wzap external provider
	WzapBaseURL  string
	WzapAdminKey string
}

func Load() *Config {
	_ = godotenv.Load()

	env := getEnv("ENVIRONMENT", "development")
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		if env == "production" {
			panic("JWT_SECRET is required in production")
		}
		jwtSecret = "zenwoot-dev-secret-change-in-production"
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		ServerHost:     getEnv("SERVER_HOST", "0.0.0.0"),
		APIKey:         getEnv("API_KEY", ""),
		JWTSecret:      jwtSecret,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		Environment:    env,
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", ""),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://wzap:wzap123@localhost:5435/wzap?sslmode=disable"),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9010"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "admin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "admin123"),
		MinioBucket:    getEnv("MINIO_BUCKET", "wzap-media"),
		MinioUseSSL:    getEnvAsBool("MINIO_USE_SSL", false),

		NatsURL: getEnv("NATS_URL", "nats://localhost:4222"),

		GlobalWebhookURL: getEnv("GLOBAL_WEBHOOK_URL", ""),

		WzapBaseURL:  getEnv("WZAP_BASE_URL", "http://localhost:8080"),
		WzapAdminKey: getEnv("WZAP_ADMIN_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
