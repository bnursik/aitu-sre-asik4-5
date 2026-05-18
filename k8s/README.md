# Kubernetes Deployment

This directory contains local Kubernetes manifests for the Computer Shop microservices stack.

The manifests use local image names such as `endterm-api-gateway` and `imagePullPolicy: IfNotPresent`. This is intended for Minikube or kind after loading locally built images into the cluster.

## Build Images

From the repository root:

```bash
docker compose build
```

## Load Images Into A Local Cluster

For Minikube:

```bash
minikube image load endterm-api-gateway
minikube image load endterm-auth-service
minikube image load endterm-user-service
minikube image load endterm-product-service
minikube image load endterm-order-service
minikube image load endterm-payment-service
minikube image load endterm-notification-service
minikube image load endterm-frontend
```

For kind, replace `computer-shop` with your kind cluster name if different:

```bash
kind load docker-image endterm-api-gateway --name computer-shop
kind load docker-image endterm-auth-service --name computer-shop
kind load docker-image endterm-user-service --name computer-shop
kind load docker-image endterm-product-service --name computer-shop
kind load docker-image endterm-order-service --name computer-shop
kind load docker-image endterm-payment-service --name computer-shop
kind load docker-image endterm-notification-service --name computer-shop
kind load docker-image endterm-frontend --name computer-shop
```

## Deploy

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/
```

Check rollout state:

```bash
kubectl get pods,svc,hpa -n computer-shop
kubectl rollout status deployment/api-gateway -n computer-shop
kubectl rollout status deployment/frontend -n computer-shop
```

## Access

NodePort endpoints:

```text
Frontend:    http://<node-ip>:30080
API Gateway: http://<node-ip>:30081
Prometheus:  http://<node-ip>:30090
Grafana:     http://<node-ip>:30300
```

For Minikube, use:

```bash
minikube service frontend -n computer-shop
minikube service api-gateway -n computer-shop
minikube service prometheus -n computer-shop
minikube service grafana -n computer-shop
```

Default Grafana login:

```text
admin / admin
```

## Smoke Tests

```bash
kubectl get pods -n computer-shop
kubectl port-forward svc/api-gateway 8080:8080 -n computer-shop
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/api/products
curl -X POST http://127.0.0.1:8080/api/notifications
```

Prometheus should show scrape targets for:

- `api-gateway`
- `auth-service`
- `user-service`
- `product-service`
- `order-service`
- `payment-service`
- `notification-service`

## Notes

- HPA manifests require metrics-server in the cluster.
- PostgreSQL uses a 1 Gi persistent volume claim.
- No Ingress controller is required; NodePort services are used for the demo.

## Cleanup

```bash
kubectl delete namespace computer-shop
```
