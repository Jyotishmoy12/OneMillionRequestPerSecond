# DB Pool Tuning Benchmark

Date: 2026-10-06

## Goal

Compare different Postgres connection pool sizes and observe impact on throughput, p99 latency, errors, and DB pool wait metrics.

This benchmark intentionally uses the no-cache path:

```text
GET /v1/items/1?cache=false
```

That makes the API bypass Valkey and hit Postgres directly, so the DB pool metrics are meaningful.

## Environment

- API: Go API running through Docker Compose
- Database: Postgres from Docker Compose
- Cache: Valkey still running, but bypassed with `cache=false`
- Benchmark tool: `alpine/bombardier`
- Duration: `30s`
- Concurrency levels: `100`, `250`, `500`
- Target URL: `<DOCKER_HOST_API_URL>/v1/items/1?cache=false`

## Commands

For each DB pool setting, `docker-compose.yml` was updated, the API was recreated, then this command was run:

```powershell
docker compose up --build -d
.\benchmark\run-local.ps1 -Url "<DOCKER_HOST_API_URL>/v1/items/1?cache=false" -ConcurrencyLevels 100,250,500
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/metrics | Select-String "database_pool"
```

## Benchmark Results

| DB Max Conns | DB Min Conns | Concurrency | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Result |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| 10 | 2 | 100 | 11,679.45 | 350,344 | 0 | 8.56ms | 7.84ms | 13.65ms | 24.58ms | 255.96ms | Stable |
| 10 | 2 | 250 | 9,638.05 | 289,067 | 0 | 25.95ms | 23.20ms | 46.98ms | 79.22ms | 607.25ms | Stable |
| 10 | 2 | 500 | 8,249.94 | 231,329 | 15,360 | 60.83ms | 31.21ms | 308.10ms | 566.93ms | 2.68s | Broke: connection refused |
| 25 | 5 | 100 | 9,645.56 | 289,053 | 0 | 10.37ms | 9.29ms | 18.03ms | 31.70ms | 149.93ms | Stable |
| 25 | 5 | 250 | 9,602.19 | 288,015 | 0 | 26.04ms | 23.28ms | 49.93ms | 72.60ms | 315.01ms | Stable |
| 25 | 5 | 500 | 6,640.71 | 182,827 | 16,018 | 75.47ms | 31.27ms | 451.99ms | 706.04ms | 3.49s | Broke: connection refused |
| 50 | 10 | 100 | 10,055.59 | 301,636 | 0 | 9.94ms | 8.92ms | 16.40ms | 34.91ms | 126.57ms | Stable |
| 50 | 10 | 250 | 10,345.19 | 310,299 | 0 | 24.17ms | 22.60ms | 41.99ms | 59.84ms | 340.53ms | Stable |
| 50 | 10 | 500 | 7,659.40 | 211,525 | 18,460 | 65.26ms | 23.56ms | 463.41ms | 620.43ms | 2.82s | Broke: connection refused |

## DB Pool Metrics After Each Full Batch

These metrics were captured after completing the `100`, `250`, and `500` concurrency runs for each pool setting. They are cumulative for that API process.

| DB Max Conns | DB Min Conns | Total Conns Observed | Idle Conns After Run | Wait Count | Wait Duration | Notes |
| ---: | ---: | ---: | ---: | ---: | ---: | --- |
| 10 | 2 | 10 | 7 | 204,199 | 77.07s | Pool was heavily saturated. |
| 25 | 5 | 25 | 25 | 13,229 | 12.16s | Much less DB waiting than 10. |
| 50 | 10 | 36 | 36 | 611 | 1.75s | Lowest DB wait pressure. |

## Interpretation

Increasing the DB pool size reduced DB connection waiting very clearly:

```text
10 max conns -> 204,199 waits
25 max conns -> 13,229 waits
50 max conns -> 611 waits
```

The best stable no-cache result in this test was:

```text
DB pool: 50 max / 10 min
Concurrency: 250
Throughput: 10,345.19 requests/sec
p99 latency: 59.84ms
Errors: 0
```

At `500` concurrency, every pool setting still broke with `connection refused`. Because this failure remained even when DB wait pressure became very low, the next bottleneck is likely local Docker Desktop networking, TCP connection handling, or single-container server limits rather than the Postgres pool alone.

## Decision

Chosen local DB pool config for the next phase:

```text
DB_MAX_CONNS=50
DB_MIN_CONNS=10
```

Reason:

- It produced the lowest DB wait count and wait duration.
- It gave the best stable `250` concurrency no-cache result.
- It gives more room for Postgres-backed tests before moving to AWS.

## Next Action

Keep `DB_MAX_CONNS=50` and `DB_MIN_CONNS=10` for local no-cache DB testing.

For high-cache-read benchmarks, continue using the normal endpoint:

```text
GET /v1/items/1
```

For DB pool or Postgres pressure benchmarks, use:

```text
GET /v1/items/1?cache=false
```
