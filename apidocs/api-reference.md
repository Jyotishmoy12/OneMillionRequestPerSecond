# API Reference

This document describes the API surface used by the benchmark project.

All examples use placeholders instead of live URLs.

## Health Check

```text
GET /healthz
```

Purpose:

- Confirms the HTTP server is alive.
- Does not check database readiness.

Example response:

```json
{
  "status": "ok"
}
```

## Readiness Check

```text
GET /readyz
```

Purpose:

- Confirms the API can reach Postgres.
- Useful for container orchestration and load balancer health checks.

Example response:

```json
{
  "status": "ready"
}
```

## Get Item

```text
GET /v1/items/{id}
```

Purpose:

- Reads an item by ID.
- Uses Valkey first.
- Falls back to Postgres on cache miss.
- Populates Valkey after a successful database read.

Example response:

```json
{
  "id": 1,
  "name": "Mechanical Keyboard",
  "description": "Low-latency keyboard for serious typing and gaming.",
  "price_cents": 8999,
  "created_at": "timestamp"
}
```

## Get Item Without Cache

```text
GET /v1/items/{id}?cache=false
```

Purpose:

- Bypasses Valkey.
- Always reads from Postgres.
- Used for database pool and Postgres pressure benchmarks.

## Update Item

```text
PUT /v1/items/{id}
```

Request body:

```json
{
  "name": "Mechanical Keyboard Pro",
  "description": "Updated low-latency keyboard.",
  "price_cents": 9999
}
```

Purpose:

- Updates an item in Postgres.
- Deletes the matching Valkey cache entry after a successful update.

## Metrics

```text
GET /metrics
```

Purpose:

- Exposes Prometheus metrics.
- Used by Prometheus and Grafana dashboards.

## Common Error Responses

Invalid item ID:

```json
{
  "error": "invalid item id"
}
```

Missing item:

```json
{
  "error": "item not found"
}
```

Invalid update body:

```json
{
  "error": "invalid json body"
}
```

Internal error:

```json
{
  "error": "internal server error"
}
```
