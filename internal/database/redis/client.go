package redis

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds the configuration for Redis client
type Config struct {
	// Basic connection settings
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`

	// Deployment modes
	ClusterMode  bool `json:"cluster_mode" yaml:"cluster_mode"`
	SentinelMode bool `json:"sentinel_mode" yaml:"sentinel_mode"`

	// Cluster configuration
	Addrs []string `json:"addrs" yaml:"addrs"`

	// Sentinel configuration
	SentinelAddrs    []string          `json:"sentinel_addrs" yaml:"sentinel_addrs"`
	MasterName       string            `json:"master_name" yaml:"master_name"`
	SentinelPassword string            `json:"sentinel_password" yaml:"sentinel_password"`
	SentinelHostMap  map[string]string `json:"sentinel_host_map" yaml:"sentinel_host_map"`

	// Connection pool settings
	PoolSize        int           `json:"pool_size" yaml:"pool_size"`
	MinIdleConns    int           `json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`

	// Timeouts
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`

	// Retry settings
	MaxRetries      int           `json:"max_retries" yaml:"max_retries"`
	MinRetryBackoff time.Duration `json:"min_retry_backoff" yaml:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff"`
}

// DefaultConfig returns a default configuration for Redis
func DefaultConfig() *Config {
	return &Config{
		// Basic connection
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,

		// Deployment modes
		ClusterMode:  false,
		SentinelMode: false,

		// Sentinel settings
		MasterName: "goacl-master",
		SentinelHostMap: map[string]string{
			"redis-master":  "localhost",
			"redis-replica": "localhost",
		},

		// Connection pool
		PoolSize:        10,
		MinIdleConns:    2,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 30 * time.Minute,
		ConnMaxLifetime: time.Hour,

		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		// Retry settings
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}
}

// Client wraps the Redis client with additional functionality
type Client struct {
	config *Config
	client redis.UniversalClient
}

// NewClient creates a new Redis client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := &Client{
		config: config,
	}

	if err := client.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// createDialer creates a custom dialer that maps hostnames for Docker environments
func (c *Client) createDialer() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		// Map Docker hostnames to localhost if configured
		if mappedHost, exists := c.config.SentinelHostMap[host]; exists {
			addr = net.JoinHostPort(mappedHost, port)
			log.Printf("Redis: Mapped hostname %s to %s", host, mappedHost)
		}

		dialer := &net.Dialer{Timeout: c.config.DialTimeout}
		return dialer.DialContext(ctx, network, addr)
	}
}

// connect establishes a connection to Redis
func (c *Client) connect() error {
	opts := c.buildRedisOptions()
	c.client = redis.NewUniversalClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), c.config.DialTimeout)
	defer cancel()

	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	mode := c.getConnectionMode()
	log.Printf("Successfully connected to Redis in %s mode", mode)
	return nil
}

// buildRedisOptions creates Redis options based on configuration
func (c *Client) buildRedisOptions() *redis.UniversalOptions {
	// Base options common to all modes
	opts := &redis.UniversalOptions{
		Password:        c.config.Password,
		PoolSize:        c.config.PoolSize,
		MinIdleConns:    c.config.MinIdleConns,
		MaxIdleConns:    c.config.MaxIdleConns,
		ConnMaxIdleTime: c.config.ConnMaxIdleTime,
		ConnMaxLifetime: c.config.ConnMaxLifetime,
		DialTimeout:     c.config.DialTimeout,
		ReadTimeout:     c.config.ReadTimeout,
		WriteTimeout:    c.config.WriteTimeout,
		MaxRetries:      c.config.MaxRetries,
		MinRetryBackoff: c.config.MinRetryBackoff,
		MaxRetryBackoff: c.config.MaxRetryBackoff,
	}

	// Configure based on deployment mode
	switch {
	case c.config.ClusterMode:
		opts.Addrs = c.config.Addrs
	case c.config.SentinelMode:
		opts.Addrs = c.config.SentinelAddrs
		opts.MasterName = c.config.MasterName
		opts.SentinelPassword = c.config.SentinelPassword
		opts.DB = c.config.DB
		opts.Dialer = c.createDialer()
	default:
		// Single instance mode
		opts.Addrs = []string{fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)}
		opts.DB = c.config.DB
	}

	return opts
}

// getConnectionMode returns a string describing the current connection mode
func (c *Client) getConnectionMode() string {
	switch {
	case c.config.ClusterMode:
		return "cluster"
	case c.config.SentinelMode:
		return "sentinel"
	default:
		return "single"
	}
}

// HealthCheck verifies the connection to Redis is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result := c.client.Get(ctx, key)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return "", nil // Key doesn't exist
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return result.Val(), nil
}

// Set stores a key-value pair with optional expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := c.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Del deletes one or more keys
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys %v: %w", keys, err)
	}
	return nil
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result := c.client.Exists(ctx, key)
	if err := result.Err(); err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return result.Val() > 0, nil
}

// HSet sets a field in a hash
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	if err := c.client.HSet(ctx, key, values...).Err(); err != nil {
		return fmt.Errorf("failed to hset key %s: %w", key, err)
	}
	return nil
}

// HGet gets a field from a hash
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	result := c.client.HGet(ctx, key, field)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return "", nil // Field doesn't exist
		}
		return "", fmt.Errorf("failed to hget key %s field %s: %w", key, field, err)
	}
	return result.Val(), nil
}

// HGetAll gets all fields from a hash
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result := c.client.HGetAll(ctx, key)
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("failed to hgetall key %s: %w", key, err)
	}
	return result.Val(), nil
}

// SAdd adds members to a set
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	if err := c.client.SAdd(ctx, key, members...).Err(); err != nil {
		return fmt.Errorf("failed to sadd key %s: %w", key, err)
	}
	return nil
}

// SMembers gets all members of a set
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	result := c.client.SMembers(ctx, key)
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("failed to smembers key %s: %w", key, err)
	}
	return result.Val(), nil
}

// SIsMember checks if a value is a member of a set
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	result := c.client.SIsMember(ctx, key, member)
	if err := result.Err(); err != nil {
		return false, fmt.Errorf("failed to sismember key %s: %w", key, err)
	}
	return result.Val(), nil
}

// Expire sets an expiration time for a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to expire key %s: %w", key, err)
	}
	return nil
}

// Pipeline creates a new pipeline for batch operations
func (c *Client) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

// TxPipeline creates a new transaction pipeline
func (c *Client) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	if err := c.client.Publish(ctx, channel, message).Err(); err != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}
	return nil
}

// Subscribe subscribes to channels
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// FlushPattern deletes all keys matching a pattern
func (c *Client) FlushPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		// Delete in batches to avoid memory issues
		if len(keys) >= 1000 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
			keys = keys[:0] // Reset slice
		}
	}
	
	// Delete remaining keys
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete remaining keys: %w", err)
		}
	}
	
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys with pattern %s: %w", pattern, err)
	}
	
	return nil
}

// Close closes the connection to Redis
func (c *Client) Close() error {
	if c.client != nil {
		log.Println("Closing Redis connection")
		return c.client.Close()
	}
	return nil
}

// GetClient returns the underlying Redis client for advanced operations
func (c *Client) GetClient() redis.UniversalClient {
	return c.client
}

// Info returns Redis server information
func (c *Client) Info(ctx context.Context, section ...string) (string, error) {
	result := c.client.Info(ctx, section...)
	if err := result.Err(); err != nil {
		return "", fmt.Errorf("failed to get Redis info: %w", err)
	}
	return result.Val(), nil
}

// GetConnectionInfo returns information about the current connection
func (c *Client) GetConnectionInfo() string {
	var info []string
	
	if c.config.ClusterMode {
		info = append(info, fmt.Sprintf("Mode: Cluster, Addrs: %v", c.config.Addrs))
	} else if c.config.SentinelMode {
		info = append(info, fmt.Sprintf("Mode: Sentinel, Master: %s, Sentinels: %v", 
			c.config.MasterName, c.config.SentinelAddrs))
	} else {
		info = append(info, fmt.Sprintf("Mode: Single, Host: %s:%d", c.config.Host, c.config.Port))
	}
	
	info = append(info, fmt.Sprintf("Pool Size: %d", c.config.PoolSize))
	
	return strings.Join(info, ", ")
}
