package config

import "os"

// Config holds all configuration for the application
type Config struct {
	GRPC GRPCConfig
	HTTP HTTPConfig
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port string
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
