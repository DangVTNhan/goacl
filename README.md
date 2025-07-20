# GoACL

A gRPC and HTTP server implementation with Protocol Buffers.

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry points
│       └── main.go      # Server main function
├── internal/            # Private application code
│   ├── app/            # Application orchestration
│   ├── config/         # Configuration management
│   ├── handler/        # gRPC handlers (private)
│   ├── service/        # Business logic services (private)
│   └── server/         # Server setup and management
├── api/                # Generated protobuf files (OpenAPI/gRPC definitions)
├── proto/              # Protocol buffer definitions
└── vendor/             # Vendored dependencies
```

## Getting Started

### Prerequisites

- Go 1.24+
- buf (for protocol buffer generation)

### Installation

1. Install required tools:
```bash
make install-tools
```

2. Generate protobuf files:
```bash
make generate
```

3. Build the server:
```bash
make build
```

### Running the Server

#### Development Mode
```bash
make dev
```

#### Production Mode
```bash
make build
./bin/server
```

#### Direct Run
```bash
make run
```

### Configuration

The server can be configured using environment variables:

- `GRPC_PORT`: gRPC server port (default: 50051)
- `HTTP_PORT`: HTTP server port (default: 8080)

### Available Make Targets

- `make generate` - Generate Go code from proto files
- `make build` - Build the server binary
- `make run` - Run the server directly
- `make dev` - Generate and run (development workflow)
- `make clean` - Remove generated files
- `make lint` - Lint proto files
- `make format` - Format proto files
- `make help` - Show available targets

## API Endpoints

### gRPC
- Server runs on `:50051` by default
- Available services: Ping

### HTTP (gRPC-Gateway)
- Server runs on `:8080` by default
- REST endpoints are automatically generated from gRPC definitions