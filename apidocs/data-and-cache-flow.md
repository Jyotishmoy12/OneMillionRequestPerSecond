# Data And Cache Flow

The project uses Postgres as the source of truth and Valkey as a hot-read cache.

## Database Schema

The `items` table contains:

```text
id
name
description
price_cents
created_at
```

Local Docker Compose initializes this table from `migrations/001_create_items.sql`.

The Go API also has startup migration logic so cloud sidecar Postgres containers can initialize themselves during the AWS learning deployment.

## Cache Key

Item cache keys use this shape:

```text
items:{id}
```

Example:

```text
items:1
```

## Cache Hit Flow

```mermaid
flowchart TD
    A[GET item request] --> B[Parse item ID]
    B --> C[Check Valkey]
    C -->|hit| D[Record cache hit metric]
    D --> E[Return cached JSON]
```

## Cache Miss Flow

```mermaid
flowchart TD
    A[GET item request] --> B[Parse item ID]
    B --> C[Check Valkey]
    C -->|miss| D[Record cache miss metric]
    D --> E[Read item from Postgres]
    E --> F[Write item to Valkey]
    F --> G[Return JSON]
```

## Cache Bypass Flow

```mermaid
flowchart TD
    A[GET item with cache=false] --> B[Parse item ID]
    B --> C[Skip Valkey]
    C --> D[Read item from Postgres]
    D --> E[Return JSON]
```

## Update And Invalidation Flow

```mermaid
flowchart TD
    A[PUT item request] --> B[Validate JSON]
    B --> C[Update Postgres]
    C --> D[Delete Valkey key]
    D --> E[Return updated item]
```

## Why Cache Invalidation Is Delete-Based

The update path deletes the cache key instead of rewriting it directly.

Reason:

- Postgres remains the source of truth.
- The next read repopulates the cache.
- This keeps update behavior simple and predictable.

## Benchmark Meaning

Normal endpoint:

```text
GET /v1/items/1
```

This mostly measures hot-cache read performance after warm-up.

No-cache endpoint:

```text
GET /v1/items/1?cache=false
```

This measures Postgres and database pool pressure.
