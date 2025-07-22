# GoACL Implementation Plan: Relationship-Based Access Control (ReBAC) System

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
