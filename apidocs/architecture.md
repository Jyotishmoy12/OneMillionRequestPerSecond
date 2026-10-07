# Architecture

This project is a layered Go API designed for high-read traffic experiments.

The architecture is intentionally simple at the application layer and progressively more realistic at the infrastructure layer.

## Local Architecture

```mermaid
flowchart LR
    Client[Benchmark client] --> API[Go API]
    API --> Cache[Valkey cache]
    API --> DB[Postgres]
    API --> Metrics[Prometheus metrics endpoint]
    Prometheus[Prometheus] --> API
    Grafana[Grafana] --> Prometheus
```

Local Docker Compose services:

- `api`
- `postgres`
- `redis` using Valkey
- `prometheus`
- `grafana`

## AWS Learning Architecture

The final AWS learning shape used two ECS tasks behind an Application Load Balancer.

```mermaid
flowchart LR
    LoadGen[AWS load generator EC2] --> ALB[Application Load Balancer]
    ALB --> TaskA[ECS task on EC2 A]
    ALB --> TaskB[ECS task on EC2 B]

    subgraph TaskAGroup[Task A]
        APIA[Go API]
        CacheA[Valkey sidecar]
        DBA[Postgres sidecar]
        APIA --> CacheA
        APIA --> DBA
    end

    subgraph TaskBGroup[Task B]
        APIB[Go API]
        CacheB[Valkey sidecar]
        DBB[Postgres sidecar]
        APIB --> CacheB
        APIB --> DBB
    end
```

This was chosen for learning and cost control. It is not the final architecture for very high RPS production systems.

## Application Layers

```mermaid
flowchart TD
    Main[cmd/api/main.go] --> Config[internal/config]
    Main --> Controllers[internal/controller]
    Controllers --> Service[internal/service]
    Service --> Cache[internal/cache]
    Service --> Repository[internal/repository]
    Repository --> Database[internal/database]
    Main --> Middleware[internal/middleware]
    Main --> Metrics[internal/metrics]
```

Layer responsibilities:

- `cmd/api`: application wiring and server lifecycle.
- `internal/config`: environment configuration.
- `internal/controller`: HTTP request parsing and response writing.
- `internal/service`: business flow and cache orchestration.
- `internal/repository`: Postgres queries.
- `internal/cache`: Valkey read/write/delete logic.
- `internal/database`: Postgres pool setup and migrations.
- `internal/middleware`: request logging, request IDs, metrics middleware.
- `internal/metrics`: Prometheus metric definitions and database pool observers.

## Request Path

Normal cached read:

```mermaid
sequenceDiagram
    participant C as Client
    participant API as Go API
    participant S as Item service
    participant V as Valkey
    participant P as Postgres

    C->>API: GET /v1/items/{id}
    API->>S: GetByID
    S->>V: Get item
    alt cache hit
        V-->>S: item
        S-->>API: item
        API-->>C: 200 JSON
    else cache miss
        V-->>S: miss
        S->>P: SELECT item
        P-->>S: item
        S->>V: Set item
        S-->>API: item
        API-->>C: 200 JSON
    end
```

Update path:

```mermaid
sequenceDiagram
    participant C as Client
    participant API as Go API
    participant S as Item service
    participant P as Postgres
    participant V as Valkey

    C->>API: PUT /v1/items/{id}
    API->>S: Validate and update
    S->>P: UPDATE item
    P-->>S: updated item
    S->>V: Delete cached item
    S-->>API: updated item
    API-->>C: 200 JSON
```

## Important Design Decisions

- Use Go standard `net/http` routing.
- Keep packages small and explicit.
- Use `pgxpool` for Postgres pooling.
- Use Valkey for hot-read caching.
- Use Prometheus metrics for visibility during load testing.
- Use Terraform for repeatable AWS infrastructure.
- Keep AWS tests short to control cost.

## Production Gap

The AWS sidecar database/cache pattern is useful for learning but not ideal for production. A production path would normally move Postgres and Valkey to shared managed services or dedicated clusters.
