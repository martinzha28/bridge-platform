package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
	Env         string
}

func Load() Config {
	env := getEnv("ENV", "development")

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		if env == "production" {
			log.Fatal("JWT_SECRET must be set when ENV=production")
		}
		jwtSecret = "dev-secret-change-in-production"
	}

	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bridge?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   jwtSecret,
		Env:         env,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
