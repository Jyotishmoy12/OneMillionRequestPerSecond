# Limitations And Next Steps

This project is a learning project, not a finished 1M RPS production system.

## Current Limitations

Small AWS instances:

```text
The AWS tests used small EC2 instances, so CPU, memory, and network capacity were limited.
```

Sidecar database and cache:

```text
Each ECS task had its own Postgres and Valkey sidecars.
This helped learning and cost control, but it is not a shared production data layer.
```

Limited load generation:

```text
Only one AWS-side load generator was used.
At higher scale, the load generator can become a bottleneck too.
```

Public ALB path:

```text
The ALB added realistic infrastructure behavior, but also added another network hop.
```

Per-request work:

```text
The API still performs JSON encoding, middleware, metrics, request IDs, and logging.
These are useful but not free under load.
```

No automated autoscaling policy:

```text
The project used manual scaling steps during experiments.
```

No production secrets management:

```text
The learning setup used simple environment variables.
Production should use a proper secrets system.
```

## Current Proven Capacity

```text
Local Docker hot-cache: about 12.7k requests/sec
AWS two-EC2 ALB from AWS load generator: about 3.0k requests/sec
```

This is not a contradiction. These measured different paths and different infrastructure sizes.

## What To Improve Next

Application:

- Reduce or disable per-request logging during benchmark mode.
- Profile CPU usage.
- Profile allocations.
- Review JSON encoding overhead.
- Add pprof for deeper performance analysis.

Infrastructure:

- Test 4 API tasks across 4 EC2 instances.
- Use larger EC2 instances for API tasks.
- Add multiple AWS-side load generators.
- Move Postgres and Valkey out of sidecars.
- Add autoscaling policies.
- Add cost alarms before larger tests.

Benchmarking:

- Run each benchmark multiple times.
- Compare median result, not only one run.
- Track CPU, memory, network, cache, and database metrics during the run.
- Record exact infrastructure size with each report.

## Suggested Next Experiment

```text
4 API tasks
4 EC2 instances
1 or 2 AWS-side load generators
30-second benchmark windows
immediate teardown after testing
```

Expected learning:

- whether RPS scales roughly with task count
- whether ALB or load generator becomes the next bottleneck
- whether cache/database sidecars become wasteful

## Path Toward 1M RPS

Reaching 1M RPS would require more than just increasing concurrency.

Likely requirements:

- many API instances
- multiple load generators
- shared high-performance cache
- careful network design
- larger instances or container capacity
- reduced per-request overhead
- strong observability
- budget controls
- short test windows

The right path is incremental:

```text
3k -> 10k -> 25k -> 50k -> 100k -> beyond
```

Each step should be measured and explained before moving to the next.
