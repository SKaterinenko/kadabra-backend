package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost               string
	DBPort               int
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	ServerPort           string
	S3Bucket             string
	S3AccessKey          string
	S3SecretKey          string
	JWTSecret            string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
}

func LoadConfig() (*Config, error) {
	// загружаем .env
	_ = godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               port,
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		DBSSLMode:            os.Getenv("DB_SSLMODE"),
		ServerPort:           os.Getenv("SERVER_PORT"),
		S3Bucket:             os.Getenv("S3_BUCKET"),
		S3AccessKey:          os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:          os.Getenv("S3_SECRET_KEY"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		JWTAccessExpiration:  6 * time.Minute,
		JWTRefreshExpiration: 2 * time.Minute,
	}, nil
}
