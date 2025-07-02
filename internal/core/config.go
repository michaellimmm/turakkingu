package core

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBUri  string
	MongoDBName string
	HttpPort    string
	WebPort     string
	Domain      string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		MongoDBUri:  os.Getenv("MONGODB_URI"),
		MongoDBName: os.Getenv("MONGODB_NAME"),
		HttpPort:    os.Getenv("HTTP_PORT"),
		WebPort:     os.Getenv("WEB_PORT"),
		Domain:      os.Getenv("DOMAIN"),
	}

	return config, nil
}
