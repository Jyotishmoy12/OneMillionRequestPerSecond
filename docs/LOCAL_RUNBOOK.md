# Local Runbook

This runbook explains how to run, test, observe, and benchmark the local stack before moving to AWS.

## Stack

- Go API on port `8080`
- Postgres on port `5432`
- Valkey cache on port `6379`
- Prometheus on port `9090`
- Grafana on port `3000`

## Start

```powershell
docker compose up --build -d
```

Check containers:

```powershell
docker compose ps
```

Check API logs:

```powershell
docker logs one-million-rps-api --tail 30
```

Expected startup log includes:

```text
database pool configured
api server listening
```

## Test API

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/healthz
Invoke-WebRequest -UseBasicParsing http://localhost:8080/readyz
Invoke-WebRequest -UseBasicParsing http://localhost:8080/v1/items/1
```

Expected:

- `/healthz` returns `{"status":"ok"}`
- `/readyz` returns `{"status":"ready"}`
- `/v1/items/1` returns item JSON

## Test Cache

Clear cache:

```powershell
docker exec -it one-million-rps-redis valkey-cli FLUSHALL
```

First request should miss cache and refill from Postgres:

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/v1/items/1
```

Second request should hit cache:

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/v1/items/1
```

Check cache metrics:

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/metrics | Select-String "item_cache_requests_total"
```

Check raw cache value:

```powershell
docker exec -it one-million-rps-redis valkey-cli GET items:1
```

## Test No-Cache DB Path

Use this only for DB tuning and benchmarking:

```powershell
Invoke-WebRequest -UseBasicParsing "http://localhost:8080/v1/items/1?cache=false"
```

This bypasses Valkey for that request only.

## Test Cache Invalidation

Warm cache:

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/v1/items/1
docker exec -it one-million-rps-redis valkey-cli GET items:1
```

Update item:

```powershell
Invoke-WebRequest `
  -UseBasicParsing `
  -Method PUT `
  -Uri http://localhost:8080/v1/items/1 `
  -ContentType "application/json" `
  -Body '{"name":"Mechanical Keyboard Pro","description":"Updated low-latency keyboard.","price_cents":9999}'
```

Cache should be deleted:

```powershell
docker exec -it one-million-rps-redis valkey-cli GET items:1
```

Expected:

```text
(nil)
```

## Observability

Prometheus targets:

```text
http://localhost:9090/targets
```

Grafana dashboard:

```text
http://localhost:3000/d/one-million-rps-api/one-million-rps-api
```

Grafana login:

```text
username: admin
password: admin
```

Important panels:

- Requests Per Second
- HTTP Latency
- Cache Hits / Misses
- DB Pool Connections
- DB Pool Wait Count
- DB Pool Wait Duration
- Go Goroutines
- Go Memory Allocated

## Benchmarks

Normal cache path:

```powershell
.\benchmark\run-local.ps1 -Url "http://host.docker.internal:8080/v1/items/1" -ConcurrencyLevels 100,250,500
```

Postgres no-cache path:

```powershell
.\benchmark\run-local.ps1 -Url "http://host.docker.internal:8080/v1/items/1?cache=false" -ConcurrencyLevels 100,250,500
```

Short burst test:

```powershell
.\benchmark\run-local.ps1 -Duration "1s" -ConcurrencyLevels 100,250,500
```

Record new results with:

```text
benchmark/RESULT_TEMPLATE.md
```

## Current Local Baselines

Cache hot path:

```text
12,770.87 requests/sec
250 concurrency
p99 56.28ms
0 errors
```

No-cache DB path with tuned DB pool:

```text
10,345.19 requests/sec
250 concurrency
p99 59.84ms
0 errors
DB_MAX_CONNS=50
DB_MIN_CONNS=10
```

## Stop

Stop containers:

```powershell
docker compose down
```

Stop and remove volumes only when you want a fresh database and Grafana state:

```powershell
docker compose down -v
```
