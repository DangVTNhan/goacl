package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DangVTNhan/goacl/internal/database/dgraph"
	"github.com/DangVTNhan/goacl/internal/database/redis"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	GRPC   GRPCConfig
	HTTP   HTTPConfig
	Dgraph *dgraph.Config
	Redis  *redis.Config
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
// It automatically loads .env files in the following order:
// 1. .env (if exists)
// 2. .env.local (if exists)
// Environment variables take precedence over .env file values
func Load() *Config {
	// Load .env files automatically
	loadEnvFiles()

	return &Config{
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Dgraph: loadDgraphConfig(),
		Redis:  loadRedisConfig(),
	}
}

// loadEnvFiles loads .env files in order of precedence
func loadEnvFiles() {
	// Try to load .env files in order of precedence
	envFiles := []string{".env.local", ".env"}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Loaded environment from %s", file)
			break // Stop at first successful load
		}
	}
}

// loadDgraphConfig loads Dgraph configuration from environment variables
func loadDgraphConfig() *dgraph.Config {
	config := dgraph.DefaultConfig()

	if host := getEnv("DGRAPH_HOST", ""); host != "" {
		config.Host = host
	}

	if port := getEnvInt("DGRAPH_PORT", 0); port != 0 {
		config.Port = port
	}

	if maxRetries := getEnvInt("DGRAPH_MAX_RETRIES", 0); maxRetries != 0 {
		config.MaxRetries = maxRetries
	}

	if retryDelay := getEnvDuration("DGRAPH_RETRY_DELAY", 0); retryDelay != 0 {
		config.RetryDelay = retryDelay
	}

	if connectTimeout := getEnvDuration("DGRAPH_CONNECT_TIMEOUT", 0); connectTimeout != 0 {
		config.ConnectTimeout = connectTimeout
	}

	if requestTimeout := getEnvDuration("DGRAPH_REQUEST_TIMEOUT", 0); requestTimeout != 0 {
		config.RequestTimeout = requestTimeout
	}

	return config
}

// loadRedisConfig loads Redis configuration from environment variables
func loadRedisConfig() *redis.Config {
	config := redis.DefaultConfig()

	if host := getEnv("REDIS_HOST", ""); host != "" {
		config.Host = host
	}

	if port := getEnvInt("REDIS_PORT", 0); port != 0 {
		config.Port = port
	}

	if password := getEnv("REDIS_PASSWORD", ""); password != "" {
		config.Password = password
	}

	if db := getEnvInt("REDIS_DB", -1); db != -1 {
		config.DB = db
	}

	// Cluster mode
	if getEnvBool("REDIS_CLUSTER_MODE", false) {
		config.ClusterMode = true
		if addrs := getEnv("REDIS_CLUSTER_ADDRS", ""); addrs != "" {
			config.Addrs = strings.Split(addrs, ",")
		}
	}

	// Sentinel mode
	if getEnvBool("REDIS_SENTINEL_MODE", false) {
		config.SentinelMode = true
		if addrs := getEnv("REDIS_SENTINEL_ADDRS", ""); addrs != "" {
			config.SentinelAddrs = strings.Split(addrs, ",")
		}
		if masterName := getEnv("REDIS_MASTER_NAME", ""); masterName != "" {
			config.MasterName = masterName
		}
		if sentinelPassword := getEnv("REDIS_SENTINEL_PASSWORD", ""); sentinelPassword != "" {
			config.SentinelPassword = sentinelPassword
		}
	}

	// Pool settings
	if poolSize := getEnvInt("REDIS_POOL_SIZE", 0); poolSize != 0 {
		config.PoolSize = poolSize
	}

	if minIdleConns := getEnvInt("REDIS_MIN_IDLE_CONNS", 0); minIdleConns != 0 {
		config.MinIdleConns = minIdleConns
	}

	// Timeouts
	if dialTimeout := getEnvDuration("REDIS_DIAL_TIMEOUT", 0); dialTimeout != 0 {
		config.DialTimeout = dialTimeout
	}

	if readTimeout := getEnvDuration("REDIS_READ_TIMEOUT", 0); readTimeout != 0 {
		config.ReadTimeout = readTimeout
	}

	if writeTimeout := getEnvDuration("REDIS_WRITE_TIMEOUT", 0); writeTimeout != 0 {
		config.WriteTimeout = writeTimeout
	}

	// Hostname mapping for Docker environments
	if hostMap := getEnv("REDIS_SENTINEL_HOST_MAP", ""); hostMap != "" {
		config.SentinelHostMap = parseHostMap(hostMap)
	}

	return config
}

// parseHostMap parses a comma-separated list of host mappings
// Format: "internal1:external1,internal2:external2"
func parseHostMap(hostMap string) map[string]string {
	result := make(map[string]string)
	if hostMap == "" {
		return result
	}

	pairs := strings.Split(hostMap, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) == 2 {
			internal := strings.TrimSpace(parts[0])
			external := strings.TrimSpace(parts[1])
			if internal != "" && external != "" {
				result[internal] = external
			}
		}
	}
	return result
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets an environment variable as a boolean
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvDuration gets an environment variable as a duration
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
