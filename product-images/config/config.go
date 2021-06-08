package config

import (
	"os"
	"strconv"
)

type ServerConfig struct {
	BindAddress string
	PORT        int
}
type StorageConfig struct {
	BasePath string
}

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			BindAddress: getEnv("BIND_ADDRESS", "localhost"),
			PORT:        getEnvAsInt("PORT", 9000),
		},
		Storage: StorageConfig{
			BasePath: getEnv("BASE_PATH", "./imagestore"),
		},
	}
}

// Helper function to read an environment or return a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment into integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to read an environment into Boolean or return a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
