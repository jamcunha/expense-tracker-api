package application

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	ServerPort  string
	PostgresUrl string
	// RedisUrl    string

	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExp     time.Duration
	JWTRefreshExp    time.Duration
}

func LoadConfig() (Config, error) {
	// Preload default values
	cfg := Config{ServerPort: "8080"}

	if port, exists := os.LookupEnv("PORT"); exists {
		cfg.ServerPort = port
	}

	if secret, exists := os.LookupEnv("JWT_ACCESS_SECRET"); exists {
		cfg.JWTAccessSecret = secret
	} else {
		return Config{}, fmt.Errorf("Environment variable JWT_ACCESS_SECRET must be set")
	}

	if secret, exists := os.LookupEnv("JWT_REFRESH_SECRET"); exists {
		cfg.JWTRefreshSecret = secret
	} else {
		return Config{}, fmt.Errorf("Environment variable JWT_REFRESH_SECRET must be set")
	}

	accessExp, exists := os.LookupEnv("JWT_ACCESS_EXPIRATION")
	if !exists {
		return Config{}, fmt.Errorf("Environment variable JWT_ACCESS_EXPIRATION must be set")
	}

	refreshExp, exists := os.LookupEnv("JWT_REFRESH_EXPIRATION")
	if !exists {
		return Config{}, fmt.Errorf("Environment variable JWT_REFRESH_EXPIRATION must be set")
	}

	accessExpDuration, err := time.ParseDuration(accessExp + "m")
	if err != nil {
		return Config{}, fmt.Errorf("Failed to parse JWT_ACCESS_EXPIRATION: %w", err)
	}

	refreshExpDuration, err := time.ParseDuration(refreshExp + "m")
	if err != nil {
		return Config{}, fmt.Errorf("Failed to parse JWT_REFRESH_EXPIRATION: %w", err)
	}

	cfg.JWTAccessExp = accessExpDuration
	cfg.JWTRefreshExp = refreshExpDuration

	// For now it's required since it's the only database supported
	// but this config gives the option to add more databases
	if dbUrl, exists := os.LookupEnv("DB_URL"); exists {
		cfg.PostgresUrl = dbUrl
	} else {
		return Config{}, fmt.Errorf("Environment variable DB_URL must be set")
	}

	return cfg, nil
}
