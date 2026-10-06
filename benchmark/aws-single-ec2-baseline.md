# AWS Single EC2 Baseline

Date: 2026-10-06

## Goal

Measure the first AWS deployment before adding a load balancer, multiple EC2 instances, managed Postgres, or managed Valkey.

This benchmark proves the API can run end to end on AWS and gives us a cautious starting point for scaling.

## Environment

- Region: `ap-south-1`
- Compute: ECS on one EC2 instance
- API: Go API container
- Database: Postgres sidecar container in the same ECS task
- Cache: Valkey sidecar container in the same ECS task
- Public endpoint: direct EC2 public IP on port `8080`
- Benchmark client: local machine using `alpine/bombardier` through Docker
- Target URL: `<EC2_API_URL>/v1/items/1`
- Duration: `30s`

Important: these numbers include public internet latency between the local benchmark client and AWS Mumbai. This is not the same as running the load generator inside AWS.

## How To Read These Numbers

`Avg req/sec` is the average number of requests handled per second during the 30-second run.

`p50`, `p95`, and `p99` are latency percentiles:

- `p50`: 50% of requests completed within this time.
- `p95`: 95% of requests completed within this time.
- `p99`: 99% of requests completed within this time.

`Errors` should stay at `0` for a result to be considered stable.

## Commands

```powershell
docker run --rm alpine/bombardier -c 25 -d 30s -l <EC2_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 50 -d 30s -l <EC2_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 100 -d 30s -l <EC2_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 250 -d 30s -l <EC2_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 500 -d 30s -l <EC2_API_URL>/v1/items/1
```

## Results

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Throughput | Result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/v1/items/1` | 25 | 30s | 312.75 | 9,409 | 0 | 79.84ms | 77.31ms | 98.64ms | 144.99ms | 369.16ms | 116.08KB/s | Stable |
| `/v1/items/1` | 50 | 30s | 626.69 | 18,853 | 0 | 79.73ms | 78.27ms | 97.97ms | 114.76ms | 224.70ms | 232.60KB/s | Stable |
| `/v1/items/1` | 100 | 30s | 1,160.67 | 34,804 | 0 | 86.30ms | 83.87ms | 105.18ms | 128.60ms | 359.58ms | 429.80KB/s | Stable |
| `/v1/items/1` | 250 | 30s | 1,486.16 | 44,795 | 0 | 167.64ms | 176.57ms | 258.34ms | 434.60ms | 2.28s | 551.81KB/s | Stable, latency rising |
| `/v1/items/1` | 500 | 30s | 1,659.17 | 50,156 | 0 | 299.65ms | 297.75ms | 398.88ms | 471.12ms | 1.00s | 615.41KB/s | Stable, likely saturated |

## Interpretation

The first AWS deployment is healthy:

```text
Health check: 200 OK
Readiness check: 200 OK
Item endpoint: 200 OK
```

Throughput increased as concurrency increased from `25` to `100`, and there were no benchmark errors.

At `250` and `500` concurrency, errors remained at `0`, but latency increased sharply while throughput improved only slightly:

```text
100 concurrency -> 1,160.67 req/sec, p99 128.60ms
250 concurrency -> 1,486.16 req/sec, p99 434.60ms
500 concurrency -> 1,659.17 req/sec, p99 471.12ms
```

That means this single EC2 deployment is still accepting requests, but it is likely saturating around `1.5k` to `1.7k` requests/sec from this local-to-AWS benchmark path.

Current proven AWS capacity from local-to-AWS testing:

```text
Stable concurrency: 500
Average throughput: about 1.66k requests/sec
p99 latency: about 471ms
Errors: 0
```

Best balanced result so far:

```text
Concurrency: 100
Average throughput: about 1.16k requests/sec
p99 latency: about 129ms
Errors: 0
```

This is lower than the local Docker benchmark because the client is outside AWS and every request crosses the public internet. Later, to measure server-side capacity more accurately, run the load generator from AWS in the same region.

## Next Action

Do not keep increasing concurrency on this exact topology as the main path to `1M RPS`. The next meaningful improvement is architectural:

```text
Add multiple API tasks behind a load balancer.
Move the benchmark client into AWS in the same region.
Then compare per-instance and total cluster throughput.
```
