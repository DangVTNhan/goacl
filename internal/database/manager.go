package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DangVTNhan/goacl/internal/database/dgraph"
	"github.com/DangVTNhan/goacl/internal/database/redis"
	"github.com/dgraph-io/dgo/v240/protos/api"
)

// Manager manages both Dgraph and Redis connections
type Manager struct {
	Dgraph *dgraph.Client
	Redis  *redis.Client
}

// NewManager creates a new database manager with the given configurations
func NewManager(dgraphConfig *dgraph.Config, redisConfig *redis.Config) (*Manager, error) {
	manager := &Manager{}

	// Initialize Dgraph client
	dgraphClient, err := dgraph.NewClient(dgraphConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Dgraph client: %w", err)
	}
	manager.Dgraph = dgraphClient

	// Initialize Redis client
	redisClient, err := redis.NewClient(redisConfig)
	if err != nil {
		// Close Dgraph client if Redis fails
		dgraphClient.Close()
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}
	manager.Redis = redisClient

	log.Println("Database manager initialized successfully")
	return manager, nil
}

// Initialize sets up the database schema and initial data
func (m *Manager) Initialize(ctx context.Context) error {
	log.Println("Initializing database schema and data...")

	// Apply Dgraph schema
	if err := m.Dgraph.ApplySchema(ctx, dgraph.Schema); err != nil {
		return fmt.Errorf("failed to apply Dgraph schema: %w", err)
	}

	// Create initial namespace configurations
	if err := m.createInitialNamespaces(ctx); err != nil {
		return fmt.Errorf("failed to create initial namespaces: %w", err)
	}

	// Test Redis connectivity
	if err := m.Redis.Set(ctx, "goacl:init", "success", time.Minute); err != nil {
		return fmt.Errorf("failed to test Redis connectivity: %w", err)
	}

	// Clean up test key
	if err := m.Redis.Del(ctx, "goacl:init"); err != nil {
		log.Printf("Warning: failed to clean up test key: %v", err)
	}

	log.Println("Database initialization completed successfully")
	return nil
}

// HealthCheck verifies both database connections are healthy
func (m *Manager) HealthCheck(ctx context.Context) error {
	// Check Dgraph health
	if err := m.Dgraph.HealthCheck(ctx); err != nil {
		return fmt.Errorf("Dgraph health check failed: %w", err)
	}

	// Check Redis health
	if err := m.Redis.HealthCheck(ctx); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// Close closes both database connections
func (m *Manager) Close() error {
	var errors []error

	// Close Redis connection
	if err := m.Redis.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close Redis connection: %w", err))
	}

	// Close Dgraph connection
	if err := m.Dgraph.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close Dgraph connection: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing database connections: %v", errors)
	}

	log.Println("Database connections closed successfully")
	return nil
}

// createInitialNamespaces creates the initial namespace configurations in Dgraph
func (m *Manager) createInitialNamespaces(ctx context.Context) error {
	log.Println("Creating initial namespace configurations...")

	for _, nsData := range dgraph.InitialNamespaces {
		// Check if namespace already exists
		query := `query checkNamespace($name: string) {
			namespace(func: eq(name, $name)) @filter(type(NamespaceConfig)) {
				uid
				name
			}
		}`

		vars := map[string]string{"$name": nsData.Name}
		resp, err := m.Dgraph.QueryWithVars(ctx, query, vars)
		if err != nil {
			return fmt.Errorf("failed to check namespace %s: %w", nsData.Name, err)
		}

		var result struct {
			Namespace []struct {
				UID  string `json:"uid"`
				Name string `json:"name"`
			} `json:"namespace"`
		}

		if err := json.Unmarshal(resp.Json, &result); err != nil {
			return fmt.Errorf("failed to unmarshal namespace check result: %w", err)
		}

		// Skip if namespace already exists
		if len(result.Namespace) > 0 {
			log.Printf("Namespace %s already exists, skipping", nsData.Name)
			continue
		}

		// Create namespace
		if err := m.createNamespace(ctx, nsData); err != nil {
			return fmt.Errorf("failed to create namespace %s: %w", nsData.Name, err)
		}

		log.Printf("Created namespace: %s", nsData.Name)
	}

	return nil
}

// createNamespace creates a single namespace with its relations
func (m *Manager) createNamespace(ctx context.Context, nsData dgraph.NamespaceConfigData) error {
	// Create namespace mutation
	now := time.Now().Format(time.RFC3339)

	// Build the mutation JSON
	namespaceUID := "_:namespace"
	mutation := map[string]interface{}{
		"uid":         namespaceUID,
		"dgraph.type": "NamespaceConfig",
		"name":        nsData.Name,
		"created_at":  now,
		"updated_at":  now,
		"relations":   []interface{}{},
	}

	// Add relations
	relations := make([]interface{}, len(nsData.Relations))
	for i, relData := range nsData.Relations {
		relationUID := fmt.Sprintf("_:relation_%d", i)
		relations[i] = map[string]interface{}{
			"uid":           relationUID,
			"dgraph.type":   "RelationConfig",
			"name":          relData.Name,
			"rewrite_rules": relData.RewriteRules,
			"namespace":     map[string]interface{}{"uid": namespaceUID},
		}
	}
	mutation["relations"] = relations

	// Convert to JSON
	mutationJSON, err := json.Marshal(mutation)
	if err != nil {
		return fmt.Errorf("failed to marshal namespace mutation: %w", err)
	}

	// Execute mutation
	mu := &api.Mutation{
		SetJson:   mutationJSON,
		CommitNow: true,
	}

	_, err = m.Dgraph.Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to execute namespace mutation: %w", err)
	}

	return nil
}

// GetNamespaceConfig retrieves a namespace configuration by name
func (m *Manager) GetNamespaceConfig(ctx context.Context, name string) (*NamespaceConfig, error) {
	query := `query getNamespace($name: string) {
		namespace(func: eq(name, $name)) @filter(type(NamespaceConfig)) {
			uid
			name
			created_at
			updated_at
			relations {
				uid
				name
				rewrite_rules
			}
		}
	}`

	vars := map[string]string{"$name": name}
	resp, err := m.Dgraph.QueryWithVars(ctx, query, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to query namespace %s: %w", name, err)
	}

	var result struct {
		Namespace []NamespaceConfig `json:"namespace"`
	}

	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal namespace result: %w", err)
	}

	if len(result.Namespace) == 0 {
		return nil, fmt.Errorf("namespace %s not found", name)
	}

	return &result.Namespace[0], nil
}

// CreateRelationTuple creates a new relation tuple
func (m *Manager) CreateRelationTuple(ctx context.Context, tuple *RelationTuple) error {
	// Check if tuple already exists
	query := `query checkTuple($namespace: string, $object_id: string, $relation: string, $user_id: string) {
		tuple(func: eq(namespace, $namespace)) @filter(type(RelationTuple) AND eq(object_id, $object_id) AND eq(relation, $relation) AND eq(user_id, $user_id)) {
			uid
		}
	}`

	vars := map[string]string{
		"$namespace": tuple.Namespace,
		"$object_id": tuple.ObjectID,
		"$relation":  tuple.Relation,
		"$user_id":   tuple.UserID,
	}

	resp, err := m.Dgraph.QueryWithVars(ctx, query, vars)
	if err != nil {
		return fmt.Errorf("failed to check existing tuple: %w", err)
	}

	var result struct {
		Tuple []struct {
			UID string `json:"uid"`
		} `json:"tuple"`
	}

	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return fmt.Errorf("failed to unmarshal tuple check result: %w", err)
	}

	// If tuple exists, update it
	var uid string
	if len(result.Tuple) > 0 {
		uid = result.Tuple[0].UID
	} else {
		uid = "_:tuple"
	}

	// Create mutation
	now := time.Now().Format(time.RFC3339)
	mutation := map[string]interface{}{
		"uid":         uid,
		"dgraph.type": "RelationTuple",
		"namespace":   tuple.Namespace,
		"object_id":   tuple.ObjectID,
		"relation":    tuple.Relation,
		"user_id":     tuple.UserID,
		"updated_at":  now,
	}

	if tuple.Userset != "" {
		mutation["userset"] = tuple.Userset
	}

	if uid == "_:tuple" {
		mutation["created_at"] = now
	}

	// Convert to JSON
	mutationJSON, err := json.Marshal(mutation)
	if err != nil {
		return fmt.Errorf("failed to marshal tuple mutation: %w", err)
	}

	// Execute mutation
	mu := &api.Mutation{
		SetJson:   mutationJSON,
		CommitNow: true,
	}

	_, err = m.Dgraph.Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to execute tuple mutation: %w", err)
	}

	// Invalidate related cache entries
	cacheKey := fmt.Sprintf("tuple:%s:%s:%s", tuple.Namespace, tuple.ObjectID, tuple.Relation)
	if err := m.Redis.Del(ctx, cacheKey); err != nil {
		log.Printf("Warning: failed to invalidate cache for key %s: %v", cacheKey, err)
	}

	return nil
}

// NamespaceConfig represents a namespace configuration
type NamespaceConfig struct {
	UID       string           `json:"uid"`
	Name      string           `json:"name"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Relations []RelationConfig `json:"relations"`
}

// RelationConfig represents a relation configuration
type RelationConfig struct {
	UID          string `json:"uid"`
	Name         string `json:"name"`
	RewriteRules string `json:"rewrite_rules"`
}

// RelationTuple represents a relation tuple
type RelationTuple struct {
	UID       string `json:"uid,omitempty"`
	Namespace string `json:"namespace"`
	ObjectID  string `json:"object_id"`
	Relation  string `json:"relation"`
	UserID    string `json:"user_id"`
	Userset   string `json:"userset,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
