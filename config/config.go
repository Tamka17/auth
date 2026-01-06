package configs

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	Secret        string
	TokenDuration int // в минутах
}

func LoadConfig() (*Config, error) {
	port := getEnv("PORT", "8080")
	secret := getEnv("SECRET_KEY", "your-secret-key")

	tokenDurationStr := getEnv("TOKEN_DURATION", "60")
	tokenDuration, err := strconv.Atoi(tokenDurationStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          port,
		Secret:        secret,
		TokenDuration: tokenDuration,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
