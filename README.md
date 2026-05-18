# Computer Shop Microservices

Minimal REST-based online shop for computer peripherals. The project is intentionally simple for demonstrating Docker Compose deployment, Prometheus monitoring, and an incident simulation with a broken database hostname.

## Services

- `api-gateway`: routes HTTP requests to backend services
- `frontend`: static Nginx page for demonstrating the shop flow
- `auth-service`: fake login
- `user-service`: hardcoded demo user
- `product-service`: hardcoded computer peripherals
- `order-service`: stores orders in PostgreSQL
- `payment-service`: fake payment response
- `notification-service`: fake email notification response
- `postgres`: database used only by order-service
- `prometheus`: metrics scraping
- `grafana`: dashboard visualization for Prometheus metrics

## Endpoints

All services expose:

- `GET /health`
- `GET /metrics`

Use the API Gateway on `http://localhost:8080`:

- `POST /api/auth/login`
- `GET /api/users/me`
- `GET /api/products`
- `GET /api/products/{id}`
- `POST /api/orders`
- `GET /api/orders`
- `POST /api/payments`
- `POST /api/notifications`

## Run

```bash
docker compose up --build
```

API Gateway:

```text
http://localhost:8080
```

Prometheus:

```text
http://localhost:9090
```

Grafana:

```text
http://localhost:3001
```

Default Grafana login:

```text
admin / admin
```

Frontend:

```text
http://localhost:3000
```

Prometheus targets:

```text
http://localhost:9090/targets
```

Grafana dashboard:

```text
http://localhost:3001/d/computer-shop-overview/computer-shop-overview
```

## Test With curl

Frontend health:

```bash
curl http://localhost:3000/health
```

Health:

```bash
curl http://localhost:8080/health
```

Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo"}'
```

Current user:

```bash
curl http://localhost:8080/api/users/me
```

Products:

```bash
curl http://localhost:8080/api/products
curl http://localhost:8080/api/products/1
```

Orders:

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product_id":1,"quantity":1}'

curl http://localhost:8080/api/orders
```

Payments:

```bash
curl -X POST http://localhost:8080/api/payments
```

Notifications:

```bash
curl -X POST http://localhost:8080/api/notifications
```

Metrics:

```bash
curl http://localhost:8080/metrics
```

Application metrics exposed by every Go service:

- `http_requests_total{service,method,path,status}`
- `http_request_duration_seconds_bucket{service,method,path}`
- `http_request_duration_seconds_sum`
- `http_request_duration_seconds_count`

SLI/SLO definitions and PromQL queries are documented in [`docs/sli-slo.md`](docs/sli-slo.md).

Example PromQL:

```promql
sum(rate(http_requests_total[5m])) by (service)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
sum(rate(http_requests_total{status!~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

Grafana health:

```bash
curl http://localhost:3001/api/health
```

## Test In Browser

Open:

```text
http://localhost:3000
```

Use the page buttons to:

- log in with the demo credentials
- load the demo user
- load products
- create an order
- list orders
- make a fake payment
- send a fake notification
- check frontend and gateway health

## Incident Simulation

The incident is a broken database hostname for `order-service`.

1. Open `docker-compose.yml`.
2. Change `DB_HOST` for `order-service` from `postgres` to `wrong-postgres`.
3. Recreate the order service:

   ```bash
   docker compose up -d --force-recreate order-service
   ```

4. Check logs:

   ```bash
   docker compose logs order-service
   ```

5. Try creating an order and confirm it fails:

   ```bash
   curl -X POST http://localhost:8080/api/orders \
     -H "Content-Type: application/json" \
     -d '{"user_id":1,"product_id":1,"quantity":1}'
   ```

6. Change `DB_HOST` back to `postgres`.
7. Restart order-service:

   ```bash
   docker compose up -d --force-recreate order-service
   ```

8. Create an order again and confirm it works.
