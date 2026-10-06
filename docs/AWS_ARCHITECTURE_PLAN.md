# AWS Architecture Plan

This is the pre-Terraform plan. It documents what we intend to build on AWS, but it intentionally does not include Terraform code.

## Goal

Move the local high-throughput API stack toward an AWS deployment that can be benchmarked step by step.

The long-term target is:

```text
1 million requests per second
very low latency
measured and documented at each stage
```

## First AWS Version

Use this architecture first:

```text
Internet
  -> Network Load Balancer
  -> ECS on EC2
  -> Go API containers
  -> ElastiCache Valkey
  -> RDS Postgres
```

We will not use Kubernetes/EKS for the first version because the goal is controlled AWS scaling without EKS control-plane cost.

## AWS Services

Planned services:

- VPC with public and private subnets
- Internet Gateway
- Security groups
- Network Load Balancer
- ECS cluster using EC2 capacity
- EC2 Auto Scaling Group for ECS capacity
- ECR for API images
- RDS Postgres
- ElastiCache Valkey
- CloudWatch logs
- Benchmark worker EC2 instances

Optional later:

- Managed Prometheus/Grafana or CloudWatch dashboards
- Multiple Availability Zones for production-style resilience
- Route 53 and TLS certificate
- WAF or rate limiting

## Deployment Phases

### Phase 1: Container Registry

Create ECR repository and push the API image.

Success criteria:

- API image builds locally.
- API image pushes to ECR.
- Image tag is documented.

### Phase 2: Network Foundation

Create VPC, subnets, routes, and security groups.

Success criteria:

- ECS instances can reach the internet for image pulls.
- API tasks can reach RDS and Valkey.
- Public traffic only enters through the load balancer.

### Phase 3: Data Layer

Create RDS Postgres and ElastiCache Valkey.

Success criteria:

- API can connect to Postgres.
- API can connect to Valkey.
- Migrations or seed data are handled deliberately.

### Phase 4: API Runtime

Run API on ECS using EC2 capacity.

Success criteria:

- `/healthz` works through the load balancer.
- `/readyz` works through the load balancer.
- `/v1/items/1` works through the load balancer.
- Logs appear in CloudWatch.

### Phase 5: AWS Benchmarking

Run load from AWS benchmark workers, not from the local laptop.

Success criteria:

- Benchmark traffic stays in AWS region.
- Results include RPS, p95, p99, errors, instance count, and approximate cost.
- Benchmark results are recorded under `benchmark/`.

## Benchmark Strategy

Start small and scale gradually:

```text
10k RPS
25k RPS
50k RPS
100k RPS
250k RPS
500k RPS
1M RPS
```

For each step, record:

- API task count
- ECS instance type and count
- DB configuration
- Valkey configuration
- Benchmark worker count
- Average RPS
- p50, p95, p99 latency
- Error count
- CPU, memory, network, and DB pool behavior
- Approximate cost

## Cost Guard

Before running larger tests:

- Confirm instance types and counts.
- Confirm expected hourly cost.
- Set a short benchmark duration first.
- Stop benchmark workers immediately after each run.
- Destroy temporary infrastructure after testing.

Default benchmark posture:

```text
guarded scale-up
```

Do not jump directly to 1M RPS.

## Current Local Evidence

Current local cache hot path:

```text
12,770.87 requests/sec
250 concurrency
p99 56.28ms
0 errors
```

Current local no-cache DB path:

```text
10,345.19 requests/sec
250 concurrency
p99 59.84ms
0 errors
DB_MAX_CONNS=50
DB_MIN_CONNS=10
```

Current local break behavior:

```text
500 concurrency causes connection refused errors locally.
Likely bottleneck: Docker Desktop networking or local connection handling.
```

## Terraform Boundary

Terraform will be added in the next phase.

When we start Terraform, we will build it step by step and validate each root before moving on. No Terraform code is included in this plan.
