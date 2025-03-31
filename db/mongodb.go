package db

import (
	"context"
	"log"
	"time"

	"urlShortner/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is our MongoDB client
var Client *mongo.Client
var Database *mongo.Database
var URLCollection *mongo.Collection

// Connect establishes a connection to MongoDB
func Connect(cfg *config.Config) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Check the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Initialize database and collection
	Database = Client.Database(cfg.MongoDB)
	URLCollection = Database.Collection(cfg.MongoCollection)
}

// Disconnect closes the MongoDB connection
func Disconnect(ctx context.Context) error {
	err := Client.Disconnect(ctx)
	if err != nil {
		log.Println("Failed to disconnect from MongoDB:", err)
		return err
	}
	log.Println("Disconnected from MongoDB")
	return nil
}
