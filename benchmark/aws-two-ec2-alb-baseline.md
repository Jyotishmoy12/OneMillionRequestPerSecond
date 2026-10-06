# AWS Two EC2 ALB Baseline

Date: 2026-10-06

## Goal

Measure the first horizontally scaled AWS deployment:

```text
local benchmark client -> public ALB -> 2 ECS tasks -> API + sidecar Postgres + sidecar Valkey
```

This benchmark checks whether adding a second API task behind a load balancer improves throughput and latency compared with the single-EC2 direct-public-IP baseline.

## Environment

- Region: `ap-south-1`
- Compute: ECS on 2 EC2 instances
- API: 2 Go API containers, one task per EC2 instance
- Database: Postgres sidecar container per task
- Cache: Valkey sidecar container per task
- Load balancer: public Application Load Balancer on port `80`
- Benchmark client: local machine using `alpine/bombardier` through Docker
- Target URL: `<ALB_API_URL>/v1/items/1`
- Duration: `30s`

Important: these numbers still include public internet latency between the local benchmark client and AWS Mumbai.

## Health Checks

Before benchmarking, both ECS service and ALB target health were confirmed:

```text
ECS desired tasks: 2
ECS running tasks: 2
ECS pending tasks: 0
ALB healthy targets: 2
```

## Commands

```powershell
docker run --rm alpine/bombardier -c 100 -d 30s -l <ALB_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 250 -d 30s -l <ALB_API_URL>/v1/items/1
docker run --rm alpine/bombardier -c 500 -d 30s -l <ALB_API_URL>/v1/items/1
```

## Results

| Endpoint | Concurrency | Duration | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Throughput | Result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/v1/items/1` | 100 | 30s | 1,218.68 | 36,588 | 0 | 82.03ms | 79.69ms | 100.44ms | 140.26ms | 325.49ms | 541.95KB/s | Stable |
| `/v1/items/1` | 250 | 30s | 2,339.33 | 70,253 | 0 | 106.92ms | 101.06ms | 159.29ms | 186.91ms | 453.00ms | 1.01MB/s | Stable |
| `/v1/items/1` | 500 | 30s | 2,833.32 | 85,340 | 0 | 176.20ms | 172.42ms | 277.03ms | 344.11ms | 793.43ms | 1.23MB/s | Stable |

## Comparison With Single EC2 Baseline

| Concurrency | Single EC2 RPS | Single EC2 p99 | 2 EC2 + ALB RPS | 2 EC2 + ALB p99 | Result |
| ---: | ---: | ---: | ---: | ---: | --- |
| 100 | 1,160.67 | 128.60ms | 1,218.68 | 140.26ms | Similar throughput; ALB adds small overhead. |
| 250 | 1,486.16 | 434.60ms | 2,339.33 | 186.91ms | Clear improvement. |
| 500 | 1,659.17 | 471.12ms | 2,833.32 | 344.11ms | Clear improvement. |

## Interpretation

Adding a second EC2-backed ECS task behind the ALB improved throughput and tail latency under higher concurrency.

The most important comparison is at `250` concurrency:

```text
Single EC2: 1,486.16 req/sec, p99 434.60ms
2 EC2 ALB: 2,339.33 req/sec, p99 186.91ms
```

This means horizontal scaling is working.

Current proven AWS capacity from local-to-AWS testing:

```text
Stable concurrency: 500
Average throughput: about 2.83k requests/sec
p99 latency: about 344ms
Errors: 0
```

## Limitations

This benchmark is still not measuring pure AWS-side capacity because the load generator is outside AWS.

Current bottlenecks may include:

- local Docker Desktop load generation
- home/ISP network path
- public internet latency to `ap-south-1`
- ALB public endpoint latency
- small EC2 instance size

The next accurate scaling benchmark should run the load generator inside AWS in the same region.
