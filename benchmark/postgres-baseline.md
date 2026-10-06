# Postgres Baseline Benchmark

Date: 2026-09-29

## Goal

Measure the current API before Redis/Valkey caching so we have a clean baseline to compare against after adding the cache layer.

Current request paths:

- `GET /healthz`: raw HTTP handler, no database call.
- `GET /v1/items/1`: API -> controller -> service -> repository -> Postgres.

## Environment

- App runtime: Go API running through Docker Compose
- Database: Postgres from Docker Compose
- Benchmark tool: `alpine/bombardier` Docker image
- Benchmark client target: `host.docker.internal:8080`
- Concurrency: `100`
- Duration: `30s`

Validation before benchmark:

```powershell
go test ./...
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/healthz
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/v1/items/1
```

## How To Read These Numbers

`Avg req/sec` means average requests per second. If a 30-second test shows `10,675.79 req/sec`, the API sustained about 10.6k requests every second during that run.

`Duration` is set to `30s` to measure sustained capacity, not only a one-second burst. A `1s` test can be noisy because connection setup, Docker networking, CPU scheduling, and short spikes can distort the result.

`p50`, `p95`, and `p99` are latency percentiles:

- `p50`: 50% of requests completed within this time.
- `p95`: 95% of requests completed within this time.
- `p99`: 99% of requests completed within this time.

For example, `p99 latency: 60.81ms` means 99 out of 100 requests finished in 60.81ms or less. The remaining 1% were slower.

`Errors` or `Error count` means requests that failed from the benchmark client's point of view. In these tests, `connection refused` means the local Docker/API networking stack stopped accepting some new TCP connections under pressure.

For capacity claims, use all three together:

```text
Requests/sec = how much traffic it handled
p99 latency = how slow the slowest normal requests were
Errors = whether the result was stable
```

## Commands

Health endpoint:

```powershell
docker run --rm alpine/bombardier -c 100 -d 30s -l <DOCKER_HOST_API_URL>/healthz
```

Postgres-backed item endpoint:

```powershell
docker run --rm alpine/bombardier -c 100 -d 30s -l <DOCKER_HOST_API_URL>/v1/items/1
```

## Results

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Error count | Avg latency | p50 | p95 | p99 | Max latency | Throughput |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `/healthz` | 100 | 30s | 5,815.59 | 174,518 | 0 | 17.19ms | 15.32ms | 33.11ms | 51.29ms | 170.65ms | 1.13MB/s |
| `/v1/items/1` | 100 | 30s | 5,731.95 | 171,966 | 0 | 17.44ms | 15.62ms | 32.93ms | 50.95ms | 282.10ms | 1.97MB/s |

Earlier raw runs without `-l` showed higher reported throughput:

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Error count | Avg latency | Max latency |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `/healthz` | 100 | 30s | 11,892.77 | 356,572 | 0 | 8.41ms | 286.85ms |
| `/v1/items/1` | 100 | 30s | 10,759.62 | 322,832 | 0 | 9.30ms | 105.39ms |

## Step-Up Break Test

These runs increase concurrency on the Postgres-backed endpoint to find where the local setup stops being stable.

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/v1/items/1` | 250 | 30s | 10,675.79 | 320,210 | 0 | 23.43ms | 21.85ms | 40.61ms | 60.81ms | 419.93ms | Stable |
| `/v1/items/1` | 500 | 30s | 8,844.79 | 244,131 | 21,331 | 56.53ms | 15.36ms | 489.93ms | 578.24ms | 2.81s | Broke: connection refused |
| `/v1/items/1` | 750 | 30s | 9,192.25 | 252,592 | 23,635 | 81.50ms | 37.03ms | 517.97ms | 620.34ms | 3.17s | Broke: connection refused |
| `/v1/items/1` | 1000 | 30s | 8,843.52 | 229,684 | 34,905 | 113.49ms | 22.30ms | 717.78ms | 0.88s | 3.03s | Broke: connection refused |
| `/v1/items/1` | 2000 | 30s | 8,883.79 | 233,250 | 33,813 | 232.55ms | 28.90ms | 0.98s | 2.83s | 38.00s | Broke: connection refused |

Highest stable result observed so far:

```text
Concurrency: 250
Throughput: 10,675.79 requests/sec
p99 latency: 60.81ms
Errors: 0
```

First confirmed break point:

```text
Concurrency: 500
Throughput: 8,844.79 requests/sec
p99 latency: 578.24ms
Errors: 21,331 connection refused
```

Interpretation: the current local Postgres-only setup is comfortable around 250 concurrency, but unstable by 500 concurrency. More concurrency does not improve useful throughput because the system starts refusing connections and tail latency becomes too high.
