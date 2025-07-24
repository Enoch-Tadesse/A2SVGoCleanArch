package infrastructure

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// initMongo creates a new MongoDB client and verifies the connection
func InitMongo(uri string) *mongo.Database {
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
