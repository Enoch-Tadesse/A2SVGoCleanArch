package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/A2SVTask7/Delivery/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AppConfig holds application configuration values loaded from environment variables
type AppConfig struct {
	MongoURI  string
	DbName    string
	JWTSecret string
	Port      string
	Timeout   time.Duration
}

func main() {
	loadEnv()              // Load environment variables from .env file if present
	config := loadConfig() // Load required config values from environment variables

	db := initMongo(config.MongoURI) // Initialize MongoDB connection
	defer func() {
		_ = db.Client().Disconnect(context.Background()) // Disconnect Mongo client on program exit
	}()

	router := gin.Default()                    // Create a default Gin router with Logger and Recovery middleware
	routers.SetUp(config.Timeout, *db, router) // Setup all routes with middleware and handlers

	log.Printf("🚀 Server running at http://localhost:%s", config.Port)
	if err := router.Run(":" + config.Port); err != nil { // Start the HTTP server
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// loadEnv loads environment variables from a .env file if it exists
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found (may be intentional)")
	}
}

// loadConfig reads configuration from environment variables with validation and defaults
func loadConfig() AppConfig {
	config := AppConfig{
		MongoURI:  os.Getenv("MONGO_URI"),
		DbName:    os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
		Timeout:   5 * time.Second, // Default request timeout duration
	}

	// Validate required environment variables
	if config.MongoURI == "" || config.DbName == "" || config.JWTSecret == "" {
		log.Fatal("❌ Required environment variables are missing")
	}

	// Default to port 8080 if none is set
	if config.Port == "" {
		config.Port = "8080"
	}
	return config
}

// initMongo creates a new MongoDB client and verifies the connection
func initMongo(uri string) *mongo.Database {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	dbName := os.Getenv("DB_NAME") // Get DB name again for clarity, could also pass in
	return client.Database(dbName)
}
