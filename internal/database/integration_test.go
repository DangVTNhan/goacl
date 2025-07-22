//go:build integration

package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DangVTNhan/goacl/internal/database/dgraph"
	"github.com/DangVTNhan/goacl/internal/database/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDatabaseIntegration tests the complete database integration
func TestDatabaseIntegration(t *testing.T) {
	ctx := context.Background()

	// Start test containers
	dgraphContainer, redisContainer, err := setupTestContainers(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test containers: %v", err)
	}
	defer func() {
		if err := dgraphContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Dgraph container: %v", err)
		}
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	}()

	// Get container connection details
	dgraphHost, err := dgraphContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Dgraph host: %v", err)
	}
	dgraphPort, err := dgraphContainer.MappedPort(ctx, "9080")
	if err != nil {
		t.Fatalf("Failed to get Dgraph port: %v", err)
	}

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get Redis port: %v", err)
	}

	// Create database manager
	dgraphConfig := &dgraph.Config{
		Host:           dgraphHost,
		Port:           dgraphPort.Int(),
		MaxRetries:     3,
		RetryDelay:     time.Second,
		ConnectTimeout: 30 * time.Second,
		RequestTimeout: 10 * time.Second,
	}

	redisConfig := &redis.Config{
		Host:        redisHost,
		Port:        redisPort.Int(),
		Password:    "",
		DB:          0,
		PoolSize:    5,
		DialTimeout: 10 * time.Second,
	}

	manager, err := NewManager(dgraphConfig, redisConfig)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}
	defer manager.Close()

	// Run integration tests
	t.Run("HealthCheck", func(t *testing.T) {
		testHealthCheck(t, manager)
	})

	t.Run("SchemaInitialization", func(t *testing.T) {
		testSchemaInitialization(t, manager)
	})

	t.Run("NamespaceOperations", func(t *testing.T) {
		testNamespaceOperations(t, manager)
	})

	t.Run("RelationTupleOperations", func(t *testing.T) {
		testRelationTupleOperations(t, manager)
	})

	t.Run("CacheOperations", func(t *testing.T) {
		testCacheOperations(t, manager)
	})
}

// setupTestContainers starts Dgraph and Redis containers for testing
func setupTestContainers(ctx context.Context) (testcontainers.Container, testcontainers.Container, error) {
	// Start Dgraph Zero
	dgraphZeroReq := testcontainers.ContainerRequest{
		Image:        "dgraph/dgraph:v24.0.2",
		ExposedPorts: []string{"5080/tcp", "6080/tcp"},
		Cmd:          []string{"dgraph", "zero", "--my=dgraph-zero:5080", "--replicas=1"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("6080/tcp").WithStartupTimeout(60 * time.Second),
	}

	dgraphZero, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: dgraphZeroReq,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start Dgraph Zero: %w", err)
	}

	// Get Zero's internal IP
	zeroIP, err := dgraphZero.ContainerIP(ctx)
	if err != nil {
		dgraphZero.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get Dgraph Zero IP: %w", err)
	}

	// Wait for Zero to be ready
	time.Sleep(10 * time.Second)

	// Start Dgraph Alpha
	dgraphAlphaReq := testcontainers.ContainerRequest{
		Image:        "dgraph/dgraph:v24.0.2",
		ExposedPorts: []string{"8080/tcp", "9080/tcp"},
		Cmd:          []string{"dgraph", "alpha", "--my=dgraph-alpha:7080", fmt.Sprintf("--zero=%s:5080", zeroIP), "--security=whitelist=0.0.0.0/0"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080/tcp").WithStartupTimeout(120 * time.Second),
	}

	dgraphAlpha, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: dgraphAlphaReq,
		Started:          true,
	})
	if err != nil {
		dgraphZero.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to start Dgraph Alpha: %w", err)
	}

	// Start Redis
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:8.0-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	if err != nil {
		dgraphZero.Terminate(ctx)
		dgraphAlpha.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to start Redis: %w", err)
	}

	return dgraphAlpha, redisContainer, nil
}

// testHealthCheck tests database health checks
func testHealthCheck(t *testing.T, manager *Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := manager.HealthCheck(ctx); err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

// testSchemaInitialization tests schema application and initialization
func testSchemaInitialization(t *testing.T, manager *Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := manager.Initialize(ctx); err != nil {
		t.Errorf("Schema initialization failed: %v", err)
	}

	// Verify namespaces were created
	for _, nsData := range dgraph.InitialNamespaces {
		ns, err := manager.GetNamespaceConfig(ctx, nsData.Name)
		if err != nil {
			t.Errorf("Failed to get namespace %s: %v", nsData.Name, err)
			continue
		}

		if ns.Name != nsData.Name {
			t.Errorf("Expected namespace name %s, got %s", nsData.Name, ns.Name)
		}

		if len(ns.Relations) != len(nsData.Relations) {
			t.Errorf("Expected %d relations for namespace %s, got %d",
				len(nsData.Relations), nsData.Name, len(ns.Relations))
		}
	}
}

// testNamespaceOperations tests namespace CRUD operations
func testNamespaceOperations(t *testing.T, manager *Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test getting existing namespace
	ns, err := manager.GetNamespaceConfig(ctx, "documents")
	if err != nil {
		t.Errorf("Failed to get documents namespace: %v", err)
		return
	}

	if ns.Name != "documents" {
		t.Errorf("Expected namespace name 'documents', got %s", ns.Name)
	}

	// Verify relations exist
	expectedRelations := []string{"owner", "editor", "viewer", "parent"}
	if len(ns.Relations) != len(expectedRelations) {
		t.Errorf("Expected %d relations, got %d", len(expectedRelations), len(ns.Relations))
	}

	relationNames := make(map[string]bool)
	for _, rel := range ns.Relations {
		relationNames[rel.Name] = true
	}

	for _, expected := range expectedRelations {
		if !relationNames[expected] {
			t.Errorf("Expected relation %s not found", expected)
		}
	}
}

// testRelationTupleOperations tests relation tuple CRUD operations
func testRelationTupleOperations(t *testing.T, manager *Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create test relation tuple
	tuple := &RelationTuple{
		Namespace: "documents",
		ObjectID:  "doc123",
		Relation:  "owner",
		UserID:    "user456",
	}

	if err := manager.CreateRelationTuple(ctx, tuple); err != nil {
		t.Errorf("Failed to create relation tuple: %v", err)
		return
	}

	// Create another tuple for the same document
	tuple2 := &RelationTuple{
		Namespace: "documents",
		ObjectID:  "doc123",
		Relation:  "viewer",
		UserID:    "user789",
	}

	if err := manager.CreateRelationTuple(ctx, tuple2); err != nil {
		t.Errorf("Failed to create second relation tuple: %v", err)
	}

	// Test duplicate creation (should update)
	if err := manager.CreateRelationTuple(ctx, tuple); err != nil {
		t.Errorf("Failed to handle duplicate relation tuple: %v", err)
	}
}

// testCacheOperations tests Redis caching functionality
func testCacheOperations(t *testing.T, manager *Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic cache operations
	key := "test:cache:key"
	value := "test_value"

	// Set value
	if err := manager.Redis.Set(ctx, key, value, time.Minute); err != nil {
		t.Errorf("Failed to set cache value: %v", err)
		return
	}

	// Get value
	retrieved, err := manager.Redis.Get(ctx, key)
	if err != nil {
		t.Errorf("Failed to get cache value: %v", err)
		return
	}

	if retrieved != value {
		t.Errorf("Expected cache value %s, got %s", value, retrieved)
	}

	// Test hash operations
	hashKey := "test:hash"
	if err := manager.Redis.HSet(ctx, hashKey, "field1", "value1", "field2", "value2"); err != nil {
		t.Errorf("Failed to set hash values: %v", err)
		return
	}

	hashValue, err := manager.Redis.HGet(ctx, hashKey, "field1")
	if err != nil {
		t.Errorf("Failed to get hash value: %v", err)
		return
	}

	if hashValue != "value1" {
		t.Errorf("Expected hash value 'value1', got %s", hashValue)
	}

	// Test set operations
	setKey := "test:set"
	if err := manager.Redis.SAdd(ctx, setKey, "member1", "member2", "member3"); err != nil {
		t.Errorf("Failed to add set members: %v", err)
		return
	}

	isMember, err := manager.Redis.SIsMember(ctx, setKey, "member2")
	if err != nil {
		t.Errorf("Failed to check set membership: %v", err)
		return
	}

	if !isMember {
		t.Error("Expected member2 to be in set")
	}

	// Cleanup
	if err := manager.Redis.Del(ctx, key, hashKey, setKey); err != nil {
		t.Errorf("Failed to cleanup cache keys: %v", err)
	}
}
