package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBType      string
	PostgresURL string
	MongoDBURL  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	return &Config{
		DBType:      os.Getenv("DB_TYPE"),
		PostgresURL: os.Getenv("POSTGRESQL_URL"),
		MongoDBURL:  os.Getenv("MONGODB_URL"),
	}
}
