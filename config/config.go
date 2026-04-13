package config

import (
	"os"
)

type Config struct {
	JWTSecret     string
	Port          string
	DBPath        string
	RateLimitRate float64
	RateLimitBurst int
}

func LoadConfig() Config {
	return Config{
		JWTSecret:      getEnv("JWT_SECRET", "golearn_secret_key_2024"),
		Port:           getEnv("PORT", "8090"),
		DBPath:         getEnv("DB_PATH", "golearn.db"),
		RateLimitRate:  5.0,
		RateLimitBurst: 10,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
