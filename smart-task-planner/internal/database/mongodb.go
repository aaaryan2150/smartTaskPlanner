package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart-task-planner/config" 
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("🔄 Connecting to MongoDB...")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}

	Client = client
	DB = client.Database(config.AppConfig.MongoDB)

	log.Println("✅ MongoDB connected!")
	log.Printf("   Database: %s", config.AppConfig.MongoDB)

	return nil
}

func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		return err
	}

	log.Println("👋 MongoDB disconnected")
	return nil
}