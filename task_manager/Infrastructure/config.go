package infrastructure

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	CollectionTask string
	CollectionUser string
	JWTSecret      string
	DBName         string
	Port           string
	Timeout        time.Duration
}

var AppConfig Config

func LoadConfig() {
	// Load .env file if exists (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading env variables directly")
	}

	AppConfig = Config{
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		CollectionTask: getEnv("COLLECTION_TASK", "tasks"),
		CollectionUser: getEnv("COLLECTION_USER", "users"),
		JWTSecret:      getEnv("JWT_SECRET", "supersecretkey"),
		DBName:         getEnv("DBName", "managers"),
		Port:           getEnv("Port", "8080"),
	}

	// set the timeout
	timeoutStr := getEnv("APP_TIMEOUT", "5s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Printf("Invalid timeout format, defaulting to 5s: %v", err)
		timeout = 5 * time.Second
	}
	AppConfig.Timeout = timeout

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
