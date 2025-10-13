package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

var AppConfig *Config

func Load() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppConfig = &Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGODB_DATABASE", "task_planner"),
	}

	log.Println("✅ Configuration loaded")
	log.Printf("   Port: %s", AppConfig.Port)
	log.Printf("   Database: %s", AppConfig.MongoDB)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}