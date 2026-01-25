package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
* Config holds the application configuration values
 */
type Config struct {
	DbUrl     string
	Port      string
	JwtSecret string
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

	return Config{
		DbUrl:     dbUrl,
		Port:      port,
		JwtSecret: jwtSecret,
	}
}
