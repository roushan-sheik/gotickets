package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Dsn              string
	JwtAccessSecret  string
	JwtRefreshSecret string
	JwtAccessExpiry  string
	JwtRefreshExpiry string
}

func LoadEnv() *Config {
	// Load .env file for local development.
	// In production (Docker), environment variables are injected by the runtime,
	// so the absence of a .env file is expected and not an error.
	_ = godotenv.Load()

	return &Config{
		Port:             os.Getenv("PORT"),
		Dsn:              os.Getenv("DSN"),
		JwtAccessSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		JwtRefreshSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		JwtAccessExpiry:  os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"),
		JwtRefreshExpiry: os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"),
	}
}
