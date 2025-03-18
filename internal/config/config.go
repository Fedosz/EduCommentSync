package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TokenFile   string
	ServerPort  string
	DatabaseURL string
	AuthURL     string
	SecretKey   []byte
}

func LoadConfig() Config {
	err := godotenv.Load() // Загружаем переменные из .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		TokenFile:   os.Getenv("TOKEN_FILE"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AuthURL:     os.Getenv("AUTH_URL"),
		SecretKey:   []byte(os.Getenv("SECRET_KEY")),
	}
}
