package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	SERVER_PORT   string
	S3_BUCKET     string
	S3_ACCESS_KEY string
	S3_SECRET_KEY string
}

func LoadConfig() (*Config, error) {
	// загружаем .env
	_ = godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        port,
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		SERVER_PORT:   os.Getenv("SERVER_PORT"),
		S3_BUCKET:     os.Getenv("S3_BUCKET"),
		S3_ACCESS_KEY: os.Getenv("S3_ACCESS_KEY"),
		S3_SECRET_KEY: os.Getenv("S3_SECRET_KEY"),
	}, nil
}
