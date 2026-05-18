# Docker Swarm Deployment

This document shows how to deploy the Computer Shop stack with Docker Swarm.

The Swarm stack uses `docker-stack.yml`. Unlike Docker Compose, Swarm does not build images during `docker stack deploy`, so build the local images first.

## Build Images

From the repository root:

```bash
docker compose build
```

The stack file uses these local image names:

- `endterm-frontend`
- `endterm-api-gateway`
- `endterm-auth-service`
- `endterm-user-service`
- `endterm-product-service`
- `endterm-order-service`
- `endterm-payment-service`
- `endterm-notification-service`

## Start Swarm

If Swarm is not already enabled:

```bash
docker swarm init
```

If Docker Compose is already running, stop it first so ports `3000`, `8080`, `9090`, and `3001` are free:

```bash
docker compose down
```

## Deploy Stack

```bash
docker stack deploy -c docker-stack.yml computer-shop
```

Check service state:

```bash
docker stack services computer-shop
docker service ls
docker service ps computer-shop_api-gateway
```

Expected replicated services:

- `computer-shop_api-gateway`: 2 replicas
- `computer-shop_product-service`: 2 replicas
- `computer-shop_order-service`: 2 replicas
- `computer-shop_payment-service`: 2 replicas
- `computer-shop_notification-service`: 2 replicas

## Access

```text
Frontend:    http://localhost:3000
API Gateway: http://localhost:8080
Prometheus:  http://localhost:9090
Grafana:     http://localhost:3001
```

Default Grafana login:

```text
admin / admin
```

## Smoke Tests

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/products
curl -X POST http://localhost:8080/api/notifications
curl http://localhost:9090/-/healthy
```

Prometheus targets:

```text
http://localhost:9090/targets
```

## Scaling Demo

Scale a service manually:

```bash
docker service scale computer-shop_product-service=3
docker service ps computer-shop_product-service
```

Scale it back to the stack baseline:

```bash
docker service scale computer-shop_product-service=2
```

## Rolling Update Demo

Force a rolling update without changing the image:

```bash
docker service update --force computer-shop_api-gateway
docker service ps computer-shop_api-gateway
```

The stack uses `update_config` with one task at a time for replicated application services.

## Evidence To Capture

For the final report, capture screenshots or command output for:

- `docker node ls`
- `docker stack services computer-shop`
- `docker service ps computer-shop_api-gateway`
- `docker service ps computer-shop_product-service`
- frontend/API smoke test output
- Prometheus targets page
- manual scaling command and resulting replica count

## Cleanup

Remove the stack:

```bash
docker stack rm computer-shop
```

Wait until services are removed:

```bash
docker stack services computer-shop
```

Leave Swarm mode only if you no longer need it:

```bash
docker swarm leave --force
```
