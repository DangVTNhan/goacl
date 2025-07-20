# GoACL Implementation Plan: Relationship-Based Access Control (ReBAC) System

## Executive Summary

This document outlines the implementation plan for transforming the existing GoACL server into a comprehensive Relationship-Based Access Control (ReBAC) system inspired by Google Zanzibar. The system will use Dgraph as the primary graph database for storing permission relationships and Redis v8 as the caching layer, while maintaining the existing gRPC/HTTP infrastructure.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Research Findings](#research-findings)
3. [Database Schema Design](#database-schema-design)
4. [Caching Strategy](#caching-strategy)
5. [API Design](#api-design)
6. [Implementation Phases](#implementation-phases)
7. [Integration Points](#integration-points)
8. [Performance Considerations](#performance-considerations)
9. [Security Considerations](#security-considerations)
10. [References](#references)

## Architecture Overview

### High-Level System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │    │   HTTP Client   │    │  Admin Client   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      gRPC Gateway         │
                    │   (existing HTTP layer)   │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │       GoACL Server       │
                    │  ┌─────────────────────┐  │
                    │  │   Authorization     │  │
                    │  │     Service         │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │   Relationship      │  │
                    │  │     Service         │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │     Cache           │  │
                    │  │     Manager         │  │
                    │  └─────────────────────┘  │
                    └─────────────┬─────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
┌─────────▼─────────┐  ┌─────────▼─────────┐  ┌─────────▼─────────┐
│   Redis Cluster   │  │   Dgraph Cluster  │  │   Configuration   │
│   (Cache Layer)   │  │  (Graph Database) │  │     Storage       │
└───────────────────┘  └───────────────────┘  └───────────────────┘
```

### Component Responsibilities

- **GoACL Server**: Core authorization engine with ReBAC logic
- **Dgraph**: Primary storage for relationship tuples and object metadata
- **Redis v8**: High-performance caching layer for authorization decisions
- **gRPC/HTTP Gateway**: Existing API layer (maintained)

## Research Findings

### Google Zanzibar Key Concepts

Based on the Zanzibar paper analysis, the following core concepts will be implemented:

1. **Relation Tuples**: Basic authorization data structure `<object>#<relation>@<user>`
2. **Namespace Configuration**: Schema definitions for different object types
3. **Userset Rewrites**: Rules for computing effective permissions
4. **Consistency Model**: External consistency with bounded staleness
5. **Zookie Protocol**: Consistency tokens for causally ordered operations

### ReBAC Patterns Identified

1. **Data Ownership**: Users own resources they create
2. **Parent-Child Resources**: Hierarchical permission inheritance
3. **User Groups**: Team-based access control
4. **Recursive Relationships**: Nested group memberships
5. **Role-Based Relations**: Traditional RBAC as a subset of ReBAC

### Dgraph Integration Benefits

- Native graph structure for relationship modeling
- GraphQL-like query language (DQL)
- ACID transactions with distributed consistency
- Built-in indexing for performance
- Horizontal scalability

### Redis v8 Caching Advantages

- Improved memory efficiency
- Enhanced data structures (JSON, streams)
- Better clustering support
- Advanced eviction policies
- Pub/Sub for cache invalidation

## Database Schema Design

### Dgraph Schema

```graphql
# Core Types
type User {
  id: string @id
  email: string @index(exact)
  name: string
  created_at: datetime
  updated_at: datetime
  
  # Relationships
  member_of: [Group] @reverse
  owns: [Resource] @reverse
  has_role: [RoleAssignment] @reverse
}

type Group {
  id: string @id
  name: string @index(exact)
  description: string
  parent_group: Group @reverse
  created_at: datetime
  updated_at: datetime
  
  # Relationships
  members: [User]
  child_groups: [Group] @reverse
  has_role: [RoleAssignment] @reverse
}

type Resource {
  id: string @id
  type: string @index(exact)
  name: string
  parent_resource: Resource @reverse
  owner: User
  created_at: datetime
  updated_at: datetime
  
  # Relationships
  child_resources: [Resource] @reverse
  permissions: [Permission] @reverse
}

type Role {
  id: string @id
  name: string @index(exact)
  description: string
  permissions: [string]
  created_at: datetime
  updated_at: datetime
}

type RoleAssignment {
  id: string @id
  user: User
  group: Group
  role: Role
  resource: Resource
  granted_by: User
  granted_at: datetime
  expires_at: datetime
}

type Permission {
  id: string @id
  action: string @index(exact)
  resource: Resource
  effect: string @index(exact) # ALLOW or DENY
  conditions: string # JSON conditions
  created_at: datetime
}

# Relation Tuples (Zanzibar-style)
type RelationTuple {
  id: string @id
  namespace: string @index(exact)
  object_id: string @index(exact)
  relation: string @index(exact)
  user_id: string @index(exact)
  userset: string # For indirect relationships
  created_at: datetime
  updated_at: datetime
}

# Namespace Configuration
type NamespaceConfig {
  id: string @id
  name: string @index(exact)
  relations: [RelationConfig]
  created_at: datetime
  updated_at: datetime
}

type RelationConfig {
  id: string @id
  name: string @index(exact)
  rewrite_rules: string # JSON configuration
  namespace: NamespaceConfig @reverse
}
```

### Indexing Strategy

```graphql
# Performance-critical indexes
RelationTuple.namespace: @index(exact)
RelationTuple.object_id: @index(exact)
RelationTuple.relation: @index(exact)
RelationTuple.user_id: @index(exact)

# Composite indexes for common queries
RelationTuple: @index(exact) on (namespace, object_id, relation)
RelationTuple: @index(exact) on (namespace, user_id, relation)

# Full-text search capabilities
Resource.name: @index(fulltext)
Group.name: @index(fulltext)
User.name: @index(fulltext)
```

## Caching Strategy

### Redis v8 Cache Architecture

#### Cache Layers

1. **L1 Cache**: In-memory application cache (5-minute TTL)
2. **L2 Cache**: Redis cluster cache (30-minute TTL)
3. **L3 Cache**: Dgraph query result cache (2-hour TTL)

#### Cache Key Patterns

```
# Authorization decisions
authz:check:{namespace}:{object_id}:{relation}:{user_id}:{timestamp}

# Relation tuples
tuple:{namespace}:{object_id}:{relation}

# User permissions
user:perms:{user_id}:{resource_type}

# Group memberships
group:members:{group_id}

# Namespace configurations
ns:config:{namespace}

# Computed usersets
userset:{namespace}:{object_id}:{relation}
```

#### Cache Invalidation Strategy

1. **Time-based**: TTL for all cache entries
2. **Event-based**: Redis Pub/Sub for real-time invalidation
3. **Version-based**: Cache versioning with consistency tokens
4. **Pattern-based**: Wildcard invalidation for related keys

#### Redis Data Structures

```redis
# Hash for relation tuples
HSET tuple:doc:123:viewer user:alice 1
HSET tuple:doc:123:viewer user:bob 1

# Set for group memberships
SADD group:eng:members user:alice user:bob user:charlie

# Sorted set for time-ordered operations
ZADD authz:timeline 1640995200 "tuple:doc:123:viewer@user:alice"

# JSON for complex objects (Redis v8)
JSON.SET ns:config:documents $ '{"relations": {"viewer": {...}}}'

# Stream for audit logs
XADD authz:audit * action check namespace documents object_id 123
```

## API Design

### gRPC Service Definitions

```protobuf
syntax = "proto3";

package goacl.v1;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";

// Authorization Service
service AuthorizationService {
  // Check if user has permission
  rpc Check(CheckRequest) returns (CheckResponse) {
    option (google.api.http) = {
      post: "/v1/check"
      body: "*"
    };
  }
  
  // Expand userset for debugging
  rpc Expand(ExpandRequest) returns (ExpandResponse) {
    option (google.api.http) = {
      post: "/v1/expand"
      body: "*"
    };
  }
  
  // List user permissions
  rpc ListPermissions(ListPermissionsRequest) returns (ListPermissionsResponse) {
    option (google.api.http) = {
      get: "/v1/users/{user_id}/permissions"
    };
  }
}

// Relationship Management Service
service RelationshipService {
  // Write relation tuple
  rpc WriteRelation(WriteRelationRequest) returns (WriteRelationResponse) {
    option (google.api.http) = {
      post: "/v1/relations"
      body: "*"
    };
  }
  
  // Delete relation tuple
  rpc DeleteRelation(DeleteRelationRequest) returns (DeleteRelationResponse) {
    option (google.api.http) = {
      delete: "/v1/relations"
    };
  }
  
  // Read relation tuples
  rpc ReadRelations(ReadRelationsRequest) returns (ReadRelationsResponse) {
    option (google.api.http) = {
      get: "/v1/relations"
    };
  }
  
  // Watch for changes
  rpc WatchRelations(WatchRelationsRequest) returns (stream WatchRelationsResponse);
}

// Configuration Service
service ConfigurationService {
  // Create/Update namespace
  rpc WriteNamespace(WriteNamespaceRequest) returns (WriteNamespaceResponse) {
    option (google.api.http) = {
      post: "/v1/namespaces"
      body: "*"
    };
  }
  
  // Get namespace configuration
  rpc ReadNamespace(ReadNamespaceRequest) returns (ReadNamespaceResponse) {
    option (google.api.http) = {
      get: "/v1/namespaces/{namespace}"
    };
  }
}

// Message Definitions
message CheckRequest {
  string namespace = 1;
  string object_id = 2;
  string relation = 3;
  string user_id = 4;
  string consistency_token = 5; // Zookie equivalent
}

message CheckResponse {
  bool allowed = 1;
  string consistency_token = 2;
  google.protobuf.Timestamp checked_at = 3;
}

message RelationTuple {
  string namespace = 1;
  string object_id = 2;
  string relation = 3;
  string user_id = 4;
  string userset = 5; // For indirect relationships
}

message WriteRelationRequest {
  repeated RelationTuple tuples = 1;
  string consistency_token = 2;
}

message WriteRelationResponse {
  string consistency_token = 1;
  google.protobuf.Timestamp written_at = 2;
}
```

### REST API Endpoints

The gRPC-Gateway will automatically generate REST endpoints:

- `POST /v1/check` - Authorization check
- `POST /v1/expand` - Expand userset
- `GET /v1/users/{user_id}/permissions` - List permissions
- `POST /v1/relations` - Write relation
- `DELETE /v1/relations` - Delete relation
- `GET /v1/relations` - Read relations
- `POST /v1/namespaces` - Create namespace
- `GET /v1/namespaces/{namespace}` - Get namespace

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
- Set up Dgraph cluster and basic schema
- Implement Redis v8 caching layer
- Create basic relation tuple CRUD operations
- Extend existing gRPC services

### Phase 2: Core ReBAC Engine (Weeks 3-4)
- Implement authorization check logic
- Add userset expansion functionality
- Create namespace configuration system
- Implement consistency token (zookie) system

### Phase 3: Advanced Features (Weeks 5-6)
- Add hierarchical relationships
- Implement group membership resolution
- Create audit logging system
- Add performance monitoring

### Phase 4: Optimization & Production (Weeks 7-8)
- Performance tuning and caching optimization
- Load testing and scalability improvements
- Security hardening
- Documentation and deployment guides

## Integration Points

### Existing Infrastructure Integration

1. **gRPC Server**: Extend existing server with new services
2. **HTTP Gateway**: Leverage existing gRPC-Gateway setup
3. **Configuration**: Extend current config system
4. **Logging**: Integrate with existing logging infrastructure

### External System Integration

1. **Identity Providers**: OIDC/SAML integration for user authentication
2. **Audit Systems**: Export authorization events to SIEM
3. **Monitoring**: Prometheus metrics and Grafana dashboards
4. **CI/CD**: Automated testing and deployment pipelines

## Performance Considerations

### Scalability Targets

- **Throughput**: 10,000+ authorization checks per second
- **Latency**: <10ms p95 for cached checks, <50ms for uncached
- **Storage**: Support for millions of relation tuples
- **Concurrent Users**: 10,000+ simultaneous users

### Optimization Strategies

1. **Caching**: Multi-level caching with intelligent invalidation
2. **Batching**: Batch multiple checks in single requests
3. **Indexing**: Optimized Dgraph indexes for common query patterns
4. **Connection Pooling**: Efficient database connection management
5. **Horizontal Scaling**: Stateless service design for easy scaling

## Security Considerations

### Data Protection

1. **Encryption**: TLS for all communications, encryption at rest
2. **Authentication**: Strong authentication for all API access
3. **Authorization**: Self-hosted authorization for admin operations
4. **Audit**: Comprehensive audit logging for compliance

### Threat Mitigation

1. **Rate Limiting**: Prevent abuse and DoS attacks
2. **Input Validation**: Strict validation of all inputs
3. **Injection Prevention**: Parameterized queries and sanitization
4. **Access Control**: Principle of least privilege

## References

### Academic Papers
- [Zanzibar: Google's Consistent, Global Authorization System](https://research.google/pubs/pub48190/)
- [Relationship-Based Access Control: Protection Model and Policy Language](https://dl.acm.org/doi/10.1145/1102120.1102146)

### Implementation References
- [OpenFGA](https://openfga.dev/) - Open source Zanzibar implementation
- [SpiceDB](https://spicedb.dev/) - Production-ready Zanzibar implementation
- [Ory Keto](https://www.ory.sh/keto/) - Cloud-native access control server

### Technical Documentation
- [Dgraph Documentation](https://dgraph.io/docs/)
- [Redis v8 Documentation](https://redis.io/docs/)
- [gRPC-Go Documentation](https://grpc.io/docs/languages/go/)

---

**Next Steps**: Begin Phase 1 implementation with Dgraph setup and basic relation tuple operations.
