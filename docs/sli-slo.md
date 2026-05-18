# SLI/SLO Design

This document defines the service level indicators and service level objectives for the Computer Shop microservices project.

Prometheus is the source of truth for measurement. The API Gateway is the main user-facing SLO boundary because all frontend and client traffic flows through it. Backend service SLIs are still tracked to support troubleshooting and dependency health analysis.

## SLO Targets

| SLI | Target | Source |
| --- | --- | --- |
| Availability | `>= 99%` | Prometheus `up` metric |
| Latency | p95 `<= 200 ms` | `http_request_duration_seconds_bucket` |
| Error rate | `<= 1%` | `http_requests_total{status=~"5.."}` |
| Request success rate | `>= 99%` | `http_requests_total{status!~"5.."}` |

For demo evidence, use 5-minute PromQL windows so changes are visible quickly during local runs and incident simulations. For production-style reporting, the same SLIs should be evaluated over a 30-day rolling window.

## SLIs And Queries

### Availability

Availability measures whether Prometheus can scrape each service.

Per-service availability over a 5-minute demo window:

```promql
avg_over_time(up[5m]) * 100
```

API Gateway availability:

```promql
avg_over_time(up{job="api-gateway"}[5m]) * 100
```

Production-style 30-day availability:

```promql
avg_over_time(up{job="api-gateway"}[30d]) * 100
```

### Request Rate

Request rate is not an SLO by itself, but it provides traffic context for latency and error-rate evaluation.

Per-service request rate:

```promql
sum(rate(http_requests_total[5m])) by (service)
```

API Gateway request rate by route:

```promql
sum(rate(http_requests_total{service="api-gateway"}[5m])) by (method, path)
```

### Latency

Latency measures p95 HTTP request duration. The SLO target is p95 less than or equal to 200 ms.

Per-service p95 latency:

```promql
histogram_quantile(
  0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service)
)
```

API Gateway p95 latency:

```promql
histogram_quantile(
  0.95,
  sum(rate(http_request_duration_seconds_bucket{service="api-gateway"}[5m])) by (le)
)
```

API Gateway p95 latency by route:

```promql
histogram_quantile(
  0.95,
  sum(rate(http_request_duration_seconds_bucket{service="api-gateway"}[5m])) by (le, method, path)
)
```

### Error Rate

Error rate measures the share of requests that return HTTP 5xx responses. The SLO target is less than or equal to 1%.

Global error rate:

```promql
sum(rate(http_requests_total{status=~"5.."}[5m]))
/
sum(rate(http_requests_total[5m]))
```

Per-service error rate:

```promql
sum(rate(http_requests_total{status=~"5.."}[5m])) by (service)
/
sum(rate(http_requests_total[5m])) by (service)
```

API Gateway error rate:

```promql
sum(rate(http_requests_total{service="api-gateway",status=~"5.."}[5m]))
/
sum(rate(http_requests_total{service="api-gateway"}[5m]))
```

### Request Success Rate

Request success rate measures the share of requests that do not return HTTP 5xx responses. The SLO target is at least 99%.

Global success rate:

```promql
sum(rate(http_requests_total{status!~"5.."}[5m]))
/
sum(rate(http_requests_total[5m]))
```

Per-service success rate:

```promql
sum(rate(http_requests_total{status!~"5.."}[5m])) by (service)
/
sum(rate(http_requests_total[5m])) by (service)
```

API Gateway success rate:

```promql
sum(rate(http_requests_total{service="api-gateway",status!~"5.."}[5m]))
/
sum(rate(http_requests_total{service="api-gateway"}[5m]))
```

## Alerting Guidance

Final Prometheus alert-rule implementation belongs to a later alerting step. The intended alerting model is:

- Page when any required service is down for more than 30 seconds.
- Page when API Gateway error rate is greater than 1% for 5 minutes.
- Warn when API Gateway p95 latency is greater than 200 ms for 5 minutes.
- Warn when backend service error rate is greater than 1% for 5 minutes.

Burn-rate style interpretation:

- Fast burn: a high error rate or outage during a short 5-minute window should trigger immediate investigation during the demo.
- Slow burn: repeated smaller SLO misses over longer windows should become capacity planning or reliability action items.

## Evidence To Capture

For the final project report, capture screenshots or query output for:

- Prometheus targets with all services up.
- `http_requests_total` showing request counts for API Gateway and backend services.
- p95 latency query result.
- error-rate or success-rate query result.
- incident simulation showing `order-service` availability or request success degradation and recovery.
