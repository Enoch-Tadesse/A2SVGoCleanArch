package main

import (
	"context"
	"log"

	"github.com/A2SVTask7/Delivery/routers"
	infrastructure "github.com/A2SVTask7/Infrastructure"
	"github.com/gin-gonic/gin"
)

func main() {
	infrastructure.LoadConfig()
	config := infrastructure.AppConfig

	db, err := infrastructure.InitMongo(context.TODO(), config.MongoURI, config.DBName) // Initialize MongoDB connection
	if err != nil {
		log.Fatal("Failed to connect to db: ", err.Error())
	}
	defer func() {
		_ = db.Client().Disconnect(context.Background()) // Disconnect Mongo client on program exit
	}()

	router := gin.Default()                            // Create a default Gin router with Logger and Recovery middleware
	routers.SetUp(config.Timeout, *db, router, config) // Setup all routes with middleware and handlers

	log.Printf("🚀 Server running at http://localhost:%s", config.Port)
	if err := router.Run(":" + config.Port); err != nil { // Start the HTTP server
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
