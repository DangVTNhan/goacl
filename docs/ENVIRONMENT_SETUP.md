# Environment Configuration Guide

This guide explains how to configure the GoACL environment variables to work with the Docker Compose setup.

## Quick Start

1. **Copy the environment template**:
   ```bash
   cp .env.example .env
   ```

2. **Start the development environment**:
   ```bash
   make dev-up
   ```

3. **Verify services are running**:
   ```bash
   make db-status
   ```

## Environment Files

The system automatically loads environment files in this order:

| File | Purpose | Priority |
|------|---------|----------|
| `.env.example` | Template with all available options | Template only |
| `.env.local` | Local overrides (optional) | High |
| `.env` | Main configuration | Medium |
| System environment | System-level variables | Highest |

**Loading Order**: `.env.local` → `.env` → System environment variables (highest priority)

## Service Endpoints

When using `docker-compose.yml`, these services are available:

### Dgraph Services
- **Dgraph Alpha (HTTP)**: http://localhost:8080
- **Dgraph Alpha (gRPC)**: localhost:9080 ← *Use this for DGRAPH_PORT*
- **Dgraph Zero**: localhost:5080 (cluster coordinator)
- **Dgraph Ratel (UI)**: http://localhost:8000

### Redis Services
- **Redis Master**: localhost:6379 ← *Use this for REDIS_PORT*
- **Redis Replica**: localhost:6380 (read-only)
- **Redis Sentinel**: localhost:26379 (monitoring/failover)
- **Redis Commander (UI)**: http://localhost:8081 (admin/admin)

### Application Services
- **GoACL gRPC**: localhost:50051
- **GoACL HTTP**: localhost:8080 ← *Note: Same as Dgraph Alpha*

## Configuration Modes

### Mode 1: Simple Development (Default in .env.local)
```bash
# Use Redis Sentinel for high availability
REDIS_SENTINEL_MODE=true
REDIS_SENTINEL_ADDRS=localhost:26379
REDIS_MASTER_NAME=goacl-master
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Mode 2: Single Redis Instance
```bash
# Simpler setup, no high availability
REDIS_SENTINEL_MODE=false
REDIS_CLUSTER_MODE=false
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Mode 3: Redis Cluster
```bash
# For distributed Redis setup
REDIS_CLUSTER_MODE=true
REDIS_CLUSTER_ADDRS=localhost:7000,localhost:7001,localhost:7002
REDIS_SENTINEL_MODE=false
```

## Key Configuration Variables

### Dgraph Settings
```bash
DGRAPH_HOST=localhost          # Dgraph server host
DGRAPH_PORT=9080              # Dgraph gRPC port (from docker-compose)
DGRAPH_MAX_RETRIES=3          # Connection retry attempts
DGRAPH_CONNECT_TIMEOUT=10s    # Connection timeout
DGRAPH_REQUEST_TIMEOUT=30s    # Request timeout
```

### Redis Settings
```bash
REDIS_HOST=localhost          # Redis master host
REDIS_PORT=6379              # Redis master port (from docker-compose)
REDIS_PASSWORD=              # Redis password (empty for local dev)
REDIS_DB=0                   # Redis database number
REDIS_POOL_SIZE=10           # Connection pool size
```

### Application Settings
```bash
GRPC_PORT=50051              # GoACL gRPC server port
HTTP_PORT=8080               # GoACL HTTP server port
LOG_LEVEL=debug              # Logging level (debug/info/warn/error)
DEV_MODE=true                # Enable development features
```

## Port Conflicts

**Important**: The GoACL HTTP server (port 8080) conflicts with Dgraph Alpha HTTP (also port 8080).

**Solutions**:

1. **Change GoACL HTTP port** (recommended):
   ```bash
   HTTP_PORT=8090
   ```

2. **Use different Dgraph ports** in docker-compose.yml:
   ```yaml
   ports:
     - "8081:8080"  # Dgraph Alpha HTTP
     - "9080:9080"  # Dgraph Alpha gRPC (keep same)
   ```

## Troubleshooting

### Connection Issues

1. **Verify services are running**:
   ```bash
   make db-status
   curl http://localhost:8080/health  # Dgraph
   redis-cli ping                     # Redis
   ```

2. **Check service logs**:
   ```bash
   make db-logs
   ```

3. **Reset environment**:
   ```bash
   make db-reset
   ```

### Configuration Validation

Test your configuration:

```bash
# Test Dgraph connection
curl http://localhost:9080/health

# Test Redis connection
redis-cli -h localhost -p 6379 ping

# Test Redis Sentinel
redis-cli -h localhost -p 26379 ping
```

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_PORT` | 50051 | GoACL gRPC server port |
| `HTTP_PORT` | 8080 | GoACL HTTP server port |
| `DGRAPH_HOST` | localhost | Dgraph server hostname |
| `DGRAPH_PORT` | 9080 | Dgraph gRPC port |
| `REDIS_HOST` | localhost | Redis server hostname |
| `REDIS_PORT` | 6379 | Redis server port |
| `REDIS_SENTINEL_MODE` | true | Enable Redis Sentinel |
| `REDIS_SENTINEL_ADDRS` | localhost:26379 | Sentinel addresses |
| `REDIS_MASTER_NAME` | goacl-master | Redis master name |
| `LOG_LEVEL` | debug | Logging level |
| `DEV_MODE` | true | Development mode |

## Best Practices

1. **Always use .env files**: Don't hardcode configuration in source code
2. **Keep .env in .gitignore**: Never commit actual environment files
3. **Use .env.local for overrides**: Local-specific settings that shouldn't be shared
4. **Document custom settings**: Add comments for non-standard configurations
5. **Test configuration changes**: Restart services after environment changes
6. **Environment precedence**: System env > .env > .env.local > defaults

## Next Steps

After configuring your environment:

1. **Start services**: `make dev-up`
2. **Build application**: `make build`
3. **Run tests**: `make test-all`
4. **Start development**: `make run`
