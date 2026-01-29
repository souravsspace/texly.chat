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
	DbUrl              string
	Port               string
	JwtSecret          string
	OpenAIAPIKey       string
	EmbeddingModel     string
	EmbeddingDimension int
	ChatModel          string
	ChatTemperature    float64
	MaxContextChunks   int
	// MinIO Configuration
	MinIOEndpoint    string
	MinIOAccessKey   string
	MinIOSecretKey   string
	MinIOBucket      string
	MinIOUseSSL      bool
	MaxUploadSizeMB  int
}


/*
* Load initializes the configuration from environment variables or defaults
*/
func Load() Config {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, using environment variables")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "data/dev.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	openAIAPIKey := os.Getenv("OpenAIAPIKey")

	embeddingModel := os.Getenv("EMBEDDING_MODEL")
	if embeddingModel == "" {
		embeddingModel = "text-embedding-3-small"
	}

	embeddingDimension := 1536 // Default for text-embedding-3-small
	if dimStr := os.Getenv("EMBEDDING_DIMENSION"); dimStr != "" {
		if dim, err := strconv.Atoi(dimStr); err == nil {
			embeddingDimension = dim
		}
	}

	chatModel := os.Getenv("OPENAI_CHAT_MODEL")
	if chatModel == "" {
		chatModel = "gpt-4o-mini"
	}

	chatTemperature := 0.7
	if tempStr := os.Getenv("CHAT_TEMPERATURE"); tempStr != "" {
		if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
			chatTemperature = temp
		}
	}

	maxContextChunks := 5
	if maxStr := os.Getenv("MAX_CONTEXT_CHUNKS"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil {
			maxContextChunks = max
		}
	}

	// MinIO Configuration
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}

	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}

	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "texly-uploads"
	}

	minioUseSSL := false
	if useSSLStr := os.Getenv("MINIO_USE_SSL"); useSSLStr == "true" {
		minioUseSSL = true
	}

	maxUploadSizeMB := 100
	if maxUploadStr := os.Getenv("MAX_UPLOAD_SIZE_MB"); maxUploadStr != "" {
		if maxUpload, err := strconv.Atoi(maxUploadStr); err == nil {
			maxUploadSizeMB = maxUpload
		}
	}

	return Config{
		DbUrl:              dbUrl,
		Port:               port,
		JwtSecret:          jwtSecret,
		OpenAIAPIKey:       openAIAPIKey,
		EmbeddingModel:     embeddingModel,
		EmbeddingDimension: embeddingDimension,
		ChatModel:          chatModel,
		ChatTemperature:    chatTemperature,
		MaxContextChunks:   maxContextChunks,
		MinIOEndpoint:      minioEndpoint,
		MinIOAccessKey:     minioAccessKey,
		MinIOSecretKey:     minioSecretKey,
		MinIOBucket:        minioBucket,
		MinIOUseSSL:        minioUseSSL,
		MaxUploadSizeMB:    maxUploadSizeMB,
	}
}
