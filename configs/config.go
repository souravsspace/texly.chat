package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

/*
* Config holds the application configuration values
 */
type Config struct {
	DatabaseURL          string
	DatabaseMaxConns     int
	DatabaseMaxIdleConns int
	Port                 string
	JWTSecret            string
	OpenAIAPIKey         string
	EmbeddingModel       string
	EmbeddingDimension   int
	ChatModel            string
	ChatTemperature      float64
	MaxContextChunks     int
	// MinIO Configuration
	MinIOEndpoint   string
	MinIOAccessKey  string
	MinIOSecretKey  string
	MinIOBucket     string
	MinIOUseSSL     bool
	MaxUploadSizeMB int
	// Redis Configuration
	RedisURL string

	// Google OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
}

/*
* Load initializes the configuration from environment variables
* It will panic if required environment variables are missing
 */
func Load() Config {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, using environment variables")
	}

	return Config{
		DatabaseURL:          getEnv("DATABASE_URL", true),
		DatabaseMaxConns:     getEnvAsInt("DATABASE_MAX_CONNS", 25),
		DatabaseMaxIdleConns: getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 5),
		Port:                 getEnv("PORT", false, "8080"),
		JWTSecret:            getEnv("JWT_SECRET", true),
		OpenAIAPIKey:         getEnv("OPENAI_API_KEY", true),
		EmbeddingModel:       getEnv("EMBEDDING_MODEL", false, "text-embedding-3-small"),
		EmbeddingDimension:   getEnvAsInt("EMBEDDING_DIMENSION", 1536),
		ChatModel:            getEnv("OPENAI_CHAT_MODEL", false, "gpt-4o-mini"),
		ChatTemperature:      getEnvAsFloat("CHAT_TEMPERATURE", 0.7),
		MaxContextChunks:     getEnvAsInt("MAX_CONTEXT_CHUNKS", 5),
		MinIOEndpoint:        getEnv("MINIO_ENDPOINT", true),
		MinIOAccessKey:       getEnv("MINIO_ACCESS_KEY", true),
		MinIOSecretKey:       getEnv("MINIO_SECRET_KEY", true),
		MinIOBucket:          getEnv("MINIO_BUCKET", false, "texly-uploads"),
		MinIOUseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		MaxUploadSizeMB:      getEnvAsInt("MAX_UPLOAD_SIZE_MB", 100),
		RedisURL:             getEnv("REDIS_URL", true),
		GoogleClientID:       getEnv("GOOGLE_CLIENT_ID", false),
		GoogleClientSecret:   getEnv("GOOGLE_CLIENT_SECRET", false),
		GoogleRedirectURL:    getEnv("GOOGLE_REDIRECT_URL", false),
		FrontendURL:          getEnv("FRONTEND_URL", false, "http://localhost:5173"),
	}
}

func getEnv(key string, required bool, fallback ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if required {
			log.Fatalf("Fatal: Environment variable %s is required but not set", key)
		}
		if len(fallback) > 0 {
			return fallback[0]
		}
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer for %s: %s. Using default: %d", key, valueStr, fallback)
		return fallback
	}
	return value
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Printf("Warning: Invalid float for %s: %s. Using default: %f", key, valueStr, fallback)
		return fallback
	}
	return value
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	return valueStr == "true"
}
