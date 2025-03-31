package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	MongoURI        string
	MongoDB         string
	MongoCollection string
	ServerPort      string
	BaseURL         string
	URLExpiration   time.Duration
}

// Load returns the application configuration
func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Set default values
	config := &Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "url_shortener"),
		MongoCollection: getEnv("MONGO_COLLECTION", "urls"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		BaseURL:         getEnv("BASE_URL", "http://localhost:8080"),
		URLExpiration:   time.Duration(getEnvAsInt("URL_EXPIRATION_HOURS", 24)) * time.Hour,
	}

	return config
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return int(intValue.Hours())
}
