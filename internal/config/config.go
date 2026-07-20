package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Dsn              string
	JwtAccessSecret  string
	JwtRefreshSecret string
}

func LoadEnv() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:             os.Getenv("PORT"),
		Dsn:              os.Getenv("DSN"),
		JwtAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JwtRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
	}
}
