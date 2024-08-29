package application

import (
	"fmt"
	"os"
)

// NOTE: for now database options are just pointers (nil if option is not used).
// For a more interchangable config it could be a DatabaseConfig struct/interface.
type Config struct {
	ServerPort  string
	PostgresUrl *string
	// RedisUrl    *string
}

func LoadConfig() (Config, error) {
	// Preload default values
	cfg := Config{ServerPort: "8080"}

	if port, exists := os.LookupEnv("PORT"); exists {
		cfg.ServerPort = port
	}

	// NOTE: For now just checks if JWT_SECRET exists to be used with os.Getenv in JWT middleware
	// preferably it should be added to Config
	if _, exists := os.LookupEnv("JWT_SECRET"); !exists {
		return Config{}, fmt.Errorf("Environment variable JWT_SECRET must be set")
	}

	// For now it's required since it's the only database supported
	// but this config gives the option to add more databases
	if dbUrl, exists := os.LookupEnv("DB_URL"); exists {
		cfg.PostgresUrl = &dbUrl
	} else {
		return Config{}, fmt.Errorf("Environment variable DB_URL must be set")
	}

	return cfg, nil
}
