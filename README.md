# One Million RPS API

A learning project for building, observing, benchmarking, and scaling a Go API toward high request-per-second workloads.

The project starts as a local Go API with Postgres and Valkey, then moves into AWS using Docker, Terraform, ECR, ECS on EC2, an Application Load Balancer, and an AWS-side load generator.

The name is intentionally aspirational. The current proven result is not 1 million RPS. The value of the project is the step-by-step path: build, measure, find the bottleneck, improve, and measure again.

## What This Project Does

- Serves a small item API written in Go.
- Reads item data from Postgres.
- Uses Valkey as a cache for hot reads.
- Invalidates cache entries on item update.
- Exposes health, readiness, and Prometheus metrics endpoints.
- Runs locally with Docker Compose.
- Provides Grafana dashboards for API, cache, runtime, and database-pool metrics.
- Includes Terraform for AWS infrastructure experiments.
- Records benchmark results in the `benchmark/` folder.

## Main API Endpoints

```text
GET /healthz
GET /readyz
GET /v1/items/{id}
GET /v1/items/{id}?cache=false
PUT /v1/items/{id}
GET /metrics
```

The `cache=false` query parameter intentionally bypasses Valkey. It is useful when benchmarking Postgres and database pool behavior.

## Local Stack

The local Docker Compose stack runs:

- Go API
- Postgres
- Valkey
- Prometheus
- Grafana

Local startup:

```powershell
docker compose up --build
```

Local verification uses placeholders in this README to avoid publishing live links:

```powershell
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/healthz
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/readyz
Invoke-WebRequest -UseBasicParsing <LOCAL_API_URL>/v1/items/1
```

## AWS Stack

The AWS learning stack evolved through these stages:

1. Remote Terraform state with S3 and DynamoDB locking.
2. ECR repository for the API image.
3. ECS cluster backed by EC2 capacity.
4. Single ECS task running API, Postgres sidecar, and Valkey sidecar.
5. Application Load Balancer.
6. Two ECS API tasks across two EC2 instances.
7. Separate EC2 load generator controlled through AWS Systems Manager.

The AWS stack was destroyed after testing to control cost.

## Current Proven Benchmarks

Local Docker hot-cache result:

```text
About 12.7k requests/sec
250 concurrency
0 errors
p99 around 56ms
```

AWS two-EC2 ALB result from laptop:

```text
About 2.8k requests/sec
500 concurrency
0 errors
p99 around 344ms
```

AWS two-EC2 ALB result from AWS-side load generator:

```text
About 3.0k requests/sec
500 concurrency
0 errors
p99 around 370ms
```

Best balanced AWS-side result:

```text
About 2.89k requests/sec
100 concurrency
0 errors
p99 around 97ms
```

## Why Local RPS Was Higher Than AWS RPS

The local benchmark ran with everything close together:

```text
local load generator -> local Docker network -> API -> local Valkey/Postgres
```

The AWS benchmark used a distributed path:

```text
AWS load generator -> ALB -> ECS task -> API -> sidecar Valkey/Postgres
```

AWS did not make a tiny setup faster by default. AWS gave the project a path to horizontal scaling, load balancing, repeatable infrastructure, and realistic cloud testing.

## Documentation

Detailed documentation lives in `apidocs/`:

- `apidocs/architecture.md`
- `apidocs/api-reference.md`
- `apidocs/data-and-cache-flow.md`
- `apidocs/observability.md`
- `apidocs/benchmarking.md`
- `apidocs/aws-deployment.md`
- `apidocs/limitations-and-next-steps.md`

Benchmark reports live in `benchmark/`.

Operational local notes live in `docs/LOCAL_RUNBOOK.md`.

## Limitations

- The current AWS test used small EC2 instances.
- Postgres and Valkey were sidecars during the AWS learning phase, not managed shared services.
- The current proven AWS capacity is around 3k RPS, not 1M RPS.
- More realistic high-scale testing needs larger instances, more API tasks, multiple load generators, cost alarms, and short controlled test windows.

## Recommended Next Step

When resuming the project:

1. Recreate the AWS dev stack.
2. Scale carefully to 4 API tasks across 4 EC2 instances.
3. Run AWS-side benchmarks.
4. Record cost and performance.
5. Destroy the stack immediately after the test.

Do not leave AWS compute or load balancers running while idle.
