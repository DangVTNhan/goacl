# GoACL - Relationship-Based Access Control (ReBAC) System

A high-performance, scalable authorization system implementing Relationship-Based Access Control (ReBAC) patterns inspired by Google Zanzibar. Built with Go, gRPC, HTTP, Dgraph, and Redis.

## Features

- **Relationship-Based Access Control (ReBAC)**: Flexible authorization model supporting complex permission relationships
- **Google Zanzibar-Inspired**: Implements core concepts like relation tuples, namespace configurations, and consistency tokens
- **High Performance**: Multi-level caching with Redis, optimized database queries, and connection pooling
- **Scalable Architecture**: Horizontal scaling support with stateless design
- **Graph Database**: Dgraph for efficient relationship modeling and traversal
- **Dual Protocol Support**: Both gRPC and HTTP/REST APIs via gRPC-Gateway
- **Production Ready**: Comprehensive logging, health checks, and graceful shutdown

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry points
│       └── main.go      # Server main function
├── internal/            # Private application code
│   ├── app/            # Application orchestration
│   ├── config/         # Configuration management
│   ├── database/       # Database clients and managers
│   │   ├── dgraph/     # Dgraph client and schema
│   │   ├── redis/      # Redis client and operations
│   │   └── manager.go  # Database manager
│   ├── handler/        # gRPC handlers (private)
│   ├── service/        # Business logic services (private)
│   └── server/         # Server setup and management
├── api/                # Generated protobuf files (OpenAPI/gRPC definitions)
├── proto/              # Protocol buffer definitions
├── config/             # Configuration files
│   └── redis/          # Redis configuration
├── docker-compose.yml  # Development environment
└── vendor/             # Vendored dependencies
```

## Quick Start

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- buf (for protocol buffer generation)

### Development Environment Setup

1. **Clone the repository**:
```bash
git clone https://github.com/DangVTNhan/goacl.git
cd goacl
```

2. **Install dependencies**:
```bash
make deps
make install-tools
```

3. **Configure environment**:
```bash
cp .env.example .env
```

4. **Start the development environment**:
```bash
make dev-up
```

This will start:
- Dgraph cluster (Alpha + Zero nodes)
- Redis cluster with sentinel
- Web UIs for development

5. **Generate protobuf files**:
```bash
make generate
```

6. **Build and run the server**:
```bash
make build
make run
```

> **📖 For detailed environment configuration, see [Environment Setup Guide](docs/ENVIRONMENT_SETUP.md)**

### Alternative: Database Only

If you only need the databases without the web UIs:

```bash
make db-up    # Start databases
make db-down  # Stop databases
```

## Configuration

The server automatically loads configuration from `.env` files. Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

The system loads environment files in this order:
1. `.env.local` (if exists) - for local overrides
2. `.env` (if exists) - main configuration
3. System environment variables (highest priority)

### Core Configuration

- `GRPC_PORT`: gRPC server port (default: 50051)
- `HTTP_PORT`: HTTP server port (default: 8080)

### Database Configuration

#### Dgraph
- `DGRAPH_HOST`: Dgraph host (default: localhost)
- `DGRAPH_PORT`: Dgraph gRPC port (default: 9080)
- `DGRAPH_MAX_RETRIES`: Connection retry attempts (default: 3)
- `DGRAPH_CONNECT_TIMEOUT`: Connection timeout (default: 10s)

#### Redis
- `REDIS_HOST`: Redis host (default: localhost)
- `REDIS_PORT`: Redis port (default: 6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number (default: 0)

#### Redis Cluster Mode
- `REDIS_CLUSTER_MODE`: Enable cluster mode (default: false)
- `REDIS_CLUSTER_ADDRS`: Comma-separated cluster addresses

#### Redis Sentinel Mode
- `REDIS_SENTINEL_MODE`: Enable sentinel mode (default: false)
- `REDIS_SENTINEL_ADDRS`: Comma-separated sentinel addresses
- `REDIS_MASTER_NAME`: Master name (default: goacl-master)

## Development

### Available Make Targets

Run `make help` to see all available targets:

#### Development Environment
- `make dev-up` - Start full development environment
- `make dev-down` - Stop development environment
- `make generate` - Generate Go code from proto files
- `make build` - Build the server binary
- `make run` - Run the server directly
- `make dev` - Generate and run (development workflow)

#### Database Management
- `make db-up` - Start database services only
- `make db-down` - Stop database services
- `make db-logs` - Show database logs
- `make db-status` - Show database service status
- `make db-clean` - Clean database volumes

#### Testing
- `make test` - Run unit tests
- `make test-integration` - Run integration tests (requires Docker)
- `make test-all` - Run all tests

### Development Services

When running `make dev-up`, the following services are available:

- **Dgraph Alpha (HTTP)**: http://localhost:8080
- **Dgraph Alpha (gRPC)**: localhost:9080
- **Dgraph Ratel (UI)**: http://localhost:8000
- **Redis**: localhost:6379
- **Redis Commander (UI)**: http://localhost:8081 (admin/admin)

## API Endpoints

### gRPC Services
- **Authorization Service**: Check permissions, expand usersets
- **Relationship Service**: Manage relation tuples
- **Configuration Service**: Manage namespace configurations
- Server runs on `:50051` by default

### HTTP (gRPC-Gateway)
- Server runs on `:8080` by default
- REST endpoints are automatically generated from gRPC definitions
- OpenAPI/Swagger documentation available

### Example API Calls

#### Check Authorization
```bash
curl -X POST http://localhost:8080/v1/check \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "documents",
    "object_id": "doc123",
    "relation": "viewer",
    "user_id": "user456"
  }'
```

#### Create Relation Tuple
```bash
curl -X POST http://localhost:8080/v1/relations \
  -H "Content-Type: application/json" \
  -d '{
    "tuples": [{
      "namespace": "documents",
      "object_id": "doc123",
      "relation": "owner",
      "user_id": "user456"
    }]
  }'
```

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
Integration tests require Docker to be running:

```bash
make test-integration
```

These tests will:
- Start Dgraph and Redis containers
- Test database connectivity and schema application
- Verify CRUD operations for relation tuples
- Test Redis caching functionality
- Clean up containers after completion

### Running All Tests
```bash
make test-all
```

## Architecture

### Database Layer
- **Dgraph**: Graph database for storing relation tuples and metadata
- **Redis**: Multi-level caching for performance optimization

### Application Layer
- **gRPC Server**: High-performance API server
- **HTTP Gateway**: REST API via gRPC-Gateway
- **Database Manager**: Unified database connection management

### Key Components
- **Relation Tuples**: Core authorization data (`<object>#<relation>@<user>`)
- **Namespace Configurations**: Schema definitions for different object types
- **Consistency Tokens**: Ensure causal consistency (Zanzibar's "zookie" concept)
- **Multi-level Caching**: L1 (in-memory), L2 (Redis), L3 (Dgraph cache)

## Production Deployment

### Environment Variables
See `.env.example` for all available configuration options.

### Health Checks
The server provides health check endpoints for monitoring:
- Database connectivity verification
- Service readiness checks

### Monitoring
- Structured logging with configurable levels
- Metrics collection ready
- Graceful shutdown handling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make test-all` to ensure all tests pass
6. Submit a pull request
