package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	SecretKey   []byte
	AppPassword string
	AppEmail    string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	config := Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		SecretKey:   []byte(os.Getenv("SEKRET_KEY")),
		AppPassword: os.Getenv("APP_PASSWORD"),
		AppEmail: os.Getenv("APP_EMAIL"),
	}
	return config
}
