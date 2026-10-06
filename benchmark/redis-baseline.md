# Redis/Valkey Cache Benchmark

Date: 2026-09-30

## Goal

Measure the API after adding Redis/Valkey cache support and compare it with the Postgres-only baseline.

Current hot-cache request path:

- `GET /v1/items/1`: API -> controller -> service -> Valkey -> response.

On a cache miss, the fallback path is:

- API -> controller -> service -> Valkey miss -> Postgres -> Valkey set -> response.

## Environment

- App runtime: Go API running through Docker Compose
- Database: Postgres from Docker Compose
- Cache: Valkey from Docker Compose
- Benchmark tool: `alpine/bombardier` Docker image
- Benchmark client target: `host.docker.internal:8080`
- Duration: `30s`

Cache warm-up before benchmark:

```powershell
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/v1/items/1
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/v1/items/1
```

## How To Read These Numbers

`Avg req/sec` means average requests per second. If a 30-second test shows `12,770.87 req/sec`, the API sustained about 12.7k requests every second during that run.

`Duration` is set to `30s` to measure sustained capacity, not only a one-second burst. A `1s` test can be noisy because connection setup, Docker networking, CPU scheduling, cache behavior, and short spikes can distort the result.

`p50`, `p95`, and `p99` are latency percentiles:

- `p50`: 50% of requests completed within this time.
- `p95`: 95% of requests completed within this time.
- `p99`: 99% of requests completed within this time.

For example, `p99 latency: 56.28ms` means 99 out of 100 requests finished in 56.28ms or less. The remaining 1% were slower.

`Errors` means requests that failed from the benchmark client's point of view. In these tests, `connection refused` means the local Docker/API networking stack stopped accepting some new TCP connections under pressure.

For capacity claims, use all three together:

```text
Requests/sec = how much traffic it handled
p99 latency = how slow the slowest normal requests were
Errors = whether the result was stable
```

## Commands

```powershell
docker run --rm alpine/bombardier -c 100 -d 30s -l <DOCKER_HOST_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 250 -d 30s -l <DOCKER_HOST_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 500 -d 30s -l <DOCKER_HOST_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 750 -d 30s -l <DOCKER_HOST_API_URL>/v1/items/1
```

## Results

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Throughput | Result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/v1/items/1` | 100 | 30s | 13,101.37 | 392,732 | 0 | 7.64ms | 6.87ms | 12.97ms | 24.83ms | 144.04ms | 4.49MB/s | Stable |
| `/v1/items/1` | 250 | 30s | 12,770.87 | 382,624 | 0 | 19.60ms | 17.97ms | 35.51ms | 56.28ms | 238.03ms | 4.38MB/s | Stable |
| `/v1/items/1` | 500 | 30s | 9,839.75 | 276,392 | 18,938 | 50.84ms | 23.85ms | 304.74ms | 516.46ms | 2.67s | 3.15MB/s | Broke: connection refused |
| `/v1/items/1` | 750 | 30s | 9,712.74 | 275,069 | 16,078 | 77.41ms | 26.22ms | 602.65ms | 1.05s | 3.67s | 3.14MB/s | Broke: connection refused |

## Comparison With Postgres Baseline

At 100 concurrency:

```text
Postgres-only: 5,731.95 requests/sec, p99 50.95ms, 0 errors
Valkey cache: 13,101.37 requests/sec, p99 24.83ms, 0 errors
```

At 250 concurrency:

```text
Postgres-only: 10,675.79 requests/sec, p99 60.81ms, 0 errors
Valkey cache: 12,770.87 requests/sec, p99 56.28ms, 0 errors
```

At 500 concurrency:

```text
Postgres-only: 8,844.79 requests/sec, p99 578.24ms, 21,331 errors
Valkey cache: 9,839.75 requests/sec, p99 516.46ms, 18,938 errors
```

## Interpretation

Valkey improves the stable hot-read path clearly at lower concurrency. At 100 concurrency, throughput more than doubled and p99 latency dropped by about half.

At 250 concurrency, the cache path remains stable and reaches the best stable result observed so far:

```text
Concurrency: 250
Throughput: 12,770.87 requests/sec
p99 latency: 56.28ms
Errors: 0
```

At 500 and 750 concurrency, the system still breaks with connection refused errors. Since this happens even after removing most Postgres pressure, the likely local bottleneck is Docker Desktop networking, connection handling, or single-container/server limits rather than Postgres alone.

## Current Proven Capacity

The current local Docker setup has proven:

```text
Stable hot-cache capacity: about 12.7k requests/sec
Stable concurrency: 250
p99 latency: about 56ms
Error count: 0
```

The first confirmed break point remains:

```text
Concurrency: 500
Failure mode: connection refused
```
