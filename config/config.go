package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL  string
	APIKey   string
	MongoURI string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		BaseURL:  os.Getenv("BASE_URL"),
		APIKey:   os.Getenv("API_KEY"),
		MongoURI: os.Getenv("MONGO_URI"),
	}
}
