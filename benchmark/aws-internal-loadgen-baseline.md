# AWS Internal Load Generator Baseline

Date: 2026-10-06

## Goal

Measure the `2 EC2 + ALB` API deployment using a load generator running inside AWS.

This removes most of the laptop-to-AWS public internet bottleneck from the benchmark path.

Test path:

```text
AWS load generator EC2 -> public ALB -> 2 ECS API tasks -> sidecar Valkey/Postgres
```

## Environment

- Region: `ap-south-1`
- API compute: ECS on 2 EC2 instances
- API entrypoint: Application Load Balancer on port `80`
- Load generator: separate Amazon Linux EC2 instance
- Load generator access: AWS Systems Manager Session/Run Command
- Benchmark tool: `alpine/bombardier`
- Target URL: `<ALB_API_URL>/v1/items/1`
- Duration: `30s`

## Commands

The benchmark was started from the local terminal using SSM:

```powershell
aws ssm send-command `
  --region ap-south-1 `
  --instance-ids <LOAD_GENERATOR_INSTANCE_ID> `
  --document-name "AWS-RunShellScript" `
  --parameters commands='["docker run --rm alpine/bombardier -c 100 -d 30s -l <ALB_API_URL>/v1/items/1"]' `
  --query "Command.CommandId" `
  --output text
```

The same command was repeated with `-c 250` and `-c 500`.

Results were fetched with:

```powershell
$env:AWS_PAGER=""
aws ssm get-command-invocation `
  --region ap-south-1 `
  --command-id <COMMAND_ID> `
  --instance-id <LOAD_GENERATOR_INSTANCE_ID> `
  --query "StandardOutputContent" `
  --output text > aws-benchmark-100.txt

Get-Content .\aws-benchmark-100.txt -Tail 30
```

## Results

| Endpoint | Source | Concurrency | Duration | Avg req/sec | Total 2xx | Errors | Avg latency | p50 | p95 | p99 | Max latency | Throughput | Result |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `/v1/items/1` | AWS loadgen | 100 | 30s | 2,891.92 | 86,466 | 0 | 34.70ms | 25.23ms | 88.19ms | 96.93ms | 186.20ms | 1.25MB/s | Stable |
| `/v1/items/1` | AWS loadgen | 250 | 30s | 3,000.94 | 89,996 | 0 | 83.49ms | 88.92ms | 185.45ms | 197.67ms | 290.23ms | 1.30MB/s | Stable, saturated |
| `/v1/items/1` | AWS loadgen | 500 | 30s | 3,031.14 | 90,793 | 0 | 165.55ms | 173.38ms | 297.95ms | 369.74ms | 471.66ms | 1.31MB/s | Stable, saturated |

## Comparison With Laptop-To-AWS Benchmark

| Concurrency | Laptop-to-AWS RPS | Laptop-to-AWS p99 | AWS Loadgen RPS | AWS Loadgen p99 | Result |
| ---: | ---: | ---: | ---: | ---: | --- |
| 100 | 1,218.68 | 140.26ms | 2,891.92 | 96.93ms | AWS-side loadgen is much faster. |
| 250 | 2,339.33 | 186.91ms | 3,000.94 | 197.67ms | AWS-side loadgen improves throughput. |
| 500 | 2,833.32 | 344.11ms | 3,031.14 | 369.74ms | Cluster is near saturation. |

## Interpretation

Moving the load generator into AWS improved the measured result from about `2.83k RPS` to about `3.03k RPS`.

The more important finding is the plateau:

```text
100 concurrency -> 2,891.92 req/sec, p99 96.93ms
250 concurrency -> 3,000.94 req/sec, p99 197.67ms
500 concurrency -> 3,031.14 req/sec, p99 369.74ms
```

After `100` concurrency, throughput barely increases, but latency rises sharply. That means the current `2 EC2 + ALB` deployment is saturating around `3k RPS`.

## Current Proven AWS Capacity

```text
Architecture: 2 ECS EC2 API tasks behind ALB
Benchmark source: EC2 load generator in the same AWS region
Stable throughput: about 3k requests/sec
Errors: 0
Best balanced result: -c 100, about 2.89k RPS, p99 about 97ms
Highest measured result: -c 500, about 3.03k RPS, p99 about 370ms
```
