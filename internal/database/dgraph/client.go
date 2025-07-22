package dgraph

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Config holds the configuration for Dgraph client
type Config struct {
	Host             string        `json:"host" yaml:"host"`
	Port             int           `json:"port" yaml:"port"`
	MaxRetries       int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay       time.Duration `json:"retry_delay" yaml:"retry_delay"`
	ConnectTimeout   time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	RequestTimeout   time.Duration `json:"request_timeout" yaml:"request_timeout"`
	MaxConnections   int           `json:"max_connections" yaml:"max_connections"`
	KeepAlive        time.Duration `json:"keep_alive" yaml:"keep_alive"`
	KeepAliveTime    time.Duration `json:"keep_alive_time" yaml:"keep_alive_time"`
	KeepAliveTimeout time.Duration `json:"keep_alive_timeout" yaml:"keep_alive_timeout"`
}

// DefaultConfig returns a default configuration for Dgraph
func DefaultConfig() *Config {
	return &Config{
		Host:             "localhost",
		Port:             9080,
		MaxRetries:       3,
		RetryDelay:       time.Second,
		ConnectTimeout:   10 * time.Second,
		RequestTimeout:   30 * time.Second,
		MaxConnections:   10,
		KeepAlive:        30 * time.Second,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 5 * time.Second,
	}
}

// Client wraps the Dgraph client with additional functionality
type Client struct {
	config *Config
	conn   *grpc.ClientConn
	client *dgo.Dgraph
}

// NewClient creates a new Dgraph client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := &Client{
		config: config,
	}

	if err := client.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to Dgraph: %w", err)
	}

	return client, nil
}

// connect establishes a connection to Dgraph with retry logic
func (c *Client) connect() error {
	var err error

	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying Dgraph connection (attempt %d/%d)", attempt+1, c.config.MaxRetries)
			time.Sleep(c.config.RetryDelay)
		}

		ctx, cancel := context.WithTimeout(context.Background(), c.config.ConnectTimeout)

		// Create gRPC connection with optimized settings
		c.conn, err = grpc.DialContext(ctx,
			fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                c.config.KeepAliveTime,
				Timeout:             c.config.KeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		)

		cancel()

		if err == nil {
			// Create Dgraph client
			c.client = dgo.NewDgraphClient(api.NewDgraphClient(c.conn))
			log.Printf("Successfully connected to Dgraph at %s:%d", c.config.Host, c.config.Port)
			return nil
		}

		log.Printf("Failed to connect to Dgraph: %v", err)
	}

	return fmt.Errorf("failed to connect to Dgraph after %d attempts: %w", c.config.MaxRetries, err)
}

// HealthCheck verifies the connection to Dgraph is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("Dgraph client is not initialized")
	}

	// Simple query to test connectivity
	query := `{
		health(func: has(dgraph.type)) {
			count(uid)
		}
	}`

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	_, err := c.client.NewTxn().Query(ctx, query)
	if err != nil {
		return fmt.Errorf("Dgraph health check failed: %w", err)
	}

	return nil
}

// ApplySchema applies the given schema to Dgraph
func (c *Client) ApplySchema(ctx context.Context, schema string) error {
	if c.client == nil {
		return fmt.Errorf("Dgraph client is not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	op := &api.Operation{Schema: schema}
	if err := c.client.Alter(ctx, op); err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	log.Println("Successfully applied Dgraph schema")
	return nil
}

// DropAll drops all data and schema from Dgraph (use with caution)
func (c *Client) DropAll(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("Dgraph client is not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	op := &api.Operation{DropAll: true}
	if err := c.client.Alter(ctx, op); err != nil {
		return fmt.Errorf("failed to drop all data: %w", err)
	}

	log.Println("Successfully dropped all data from Dgraph")
	return nil
}

// NewTransaction creates a new transaction
func (c *Client) NewTransaction() *dgo.Txn {
	return c.client.NewTxn()
}

// NewReadOnlyTransaction creates a new read-only transaction
func (c *Client) NewReadOnlyTransaction() *dgo.Txn {
	return c.client.NewReadOnlyTxn()
}

// Query executes a query and returns the result
func (c *Client) Query(ctx context.Context, query string) (*api.Response, error) {
	if c.client == nil {
		return nil, fmt.Errorf("Dgraph client is not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	txn := c.client.NewReadOnlyTxn()
	defer txn.Discard(ctx)

	return txn.Query(ctx, query)
}

// QueryWithVars executes a query with variables and returns the result
func (c *Client) QueryWithVars(ctx context.Context, query string, vars map[string]string) (*api.Response, error) {
	if c.client == nil {
		return nil, fmt.Errorf("Dgraph client is not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	txn := c.client.NewReadOnlyTxn()
	defer txn.Discard(ctx)

	return txn.QueryWithVars(ctx, query, vars)
}

// Mutate executes a mutation and returns the result
func (c *Client) Mutate(ctx context.Context, mu *api.Mutation) (*api.Response, error) {
	if c.client == nil {
		return nil, fmt.Errorf("Dgraph client is not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.RequestTimeout)
	defer cancel()

	txn := c.client.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	// Only commit manually if CommitNow is false
	// If CommitNow is true, the transaction is already committed
	if !mu.CommitNow {
		if err := txn.Commit(ctx); err != nil {
			return nil, err
		}
	}

	return assigned, nil
}

// Close closes the connection to Dgraph
func (c *Client) Close() error {
	if c.conn != nil {
		log.Println("Closing Dgraph connection")
		return c.conn.Close()
	}
	return nil
}

// GetClient returns the underlying Dgraph client for advanced operations
func (c *Client) GetClient() *dgo.Dgraph {
	return c.client
}
