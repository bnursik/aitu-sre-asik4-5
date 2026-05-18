# End Term Assignment Gap Analysis

This document compares the current repository with the requirements from `End Term Project.pdf`.

## Verdict

The repository is a suitable starting point. It already has a working microservices demo, Docker Compose deployment, Prometheus scraping, Grafana provisioning, health checks, restart policies, an incident simulation path, and Terraform VM provisioning.

Do not start a completely new repository unless you want to change the application domain. The faster and safer path is to keep this repo and add the missing SRE deliverables around it.

However, the repo is not yet enough for the final assignment submission. The biggest missing items are Kubernetes manifests, Ansible playbooks, a sixth business microservice, explicit SLI/SLO documentation, a full incident postmortem, and capacity planning evidence.

## Assignment Requirements

From the PDF, the project should demonstrate:

1. A distributed microservices architecture with 6 or more services.
2. Docker environment setup, containerized services, and Docker Compose orchestration.
3. Docker Swarm deployment.
4. Kubernetes deployment with Deployments, Services, ConfigMaps, self-healing, and scaling.
5. Terraform infrastructure provisioning, especially VM creation.
6. Ansible configuration management and deployment automation.
7. SLI/SLO definitions for availability, latency, error rate, and request success rate.
8. Prometheus metrics collection.
9. Grafana dashboards.
10. Alert configuration.
11. Simulated incident, especially an Order Service failure caused by bad database configuration.
12. Root cause analysis, recovery process, and postmortem.
13. Automation, health checks, restart policies, and deployment scripts.
14. Capacity planning with load/resource analysis and scaling strategy.
15. Final PDF evidence with Git link, screenshots, and demo proof.

## Current Repository Coverage

| Requirement | Current status | Notes |
| --- | --- | --- |
| Microservices | Present | Current backend business services are `auth-service`, `user-service`, `product-service`, `order-service`, `payment-service`, and `notification-service`. |
| Frontend | Present | Static Nginx frontend exists in `frontend/`. |
| API Gateway | Present | `api-gateway` proxies to backend services. |
| Database | Present | PostgreSQL is used by `order-service`. |
| Dockerfiles | Present | Each Go service and frontend has a Dockerfile. |
| Docker Compose | Present | `docker-compose.yml` starts the app, database, Prometheus, and Grafana. |
| Docker Swarm | Partial | README mentions `docker stack deploy`, but there is no Swarm-specific documentation or `deploy.replicas` configuration. |
| Kubernetes | Missing | No Kubernetes manifests were found. |
| Terraform | Present | `terraform/` provisions an Azure resource group, network, public IP, NSG, NIC, and Ubuntu VM. |
| Ansible | Missing | No playbooks, inventory, or roles were found. |
| Prometheus | Present | `prometheus/prometheus.yml` scrapes all Go services. |
| Grafana | Present | Dashboard and datasource provisioning exist. |
| Alerts | Partial | Prometheus alert rules exist for service down and order-service down. |
| SLI/SLO docs | Missing | No explicit SLI/SLO document was found. |
| Application metrics | Present | All Go services expose request count and latency metrics with `service`, `method`, `path`, and `status` labels where applicable. |
| Health checks | Present | Compose health checks and service `/health` endpoints exist. |
| Restart policies | Present | Compose uses `restart: unless-stopped`. |
| Incident simulation | Partial | README documents a bad `DB_HOST` simulation and there is a log-check script, but no formal incident report or postmortem. |
| Capacity planning | Missing | No load test results, resource usage analysis, or scaling plan document was found. |
| Final evidence | Missing | Screenshots and final PDF report are not in the repo. |

## Suitability As A Starting Point

Use this repository as the base because it already matches the assignment's intended architecture:

- Computer shop domain is simple and understandable.
- There are multiple independent Go services.
- Services are containerized.
- The system has a frontend, API gateway, PostgreSQL, Prometheus, and Grafana.
- The incident scenario in the assignment is already aligned with `order-service` database configuration failure.
- Terraform already provisions a cloud VM, which is one of the major assignment requirements.

Starting from scratch would mostly repeat existing work. The only reason to start over would be if the current implementation cannot be demonstrated on your target environment. Based on the files present, improving this repo is the better option.

## What Needs To Be Added

### 1. Add A Sixth Business Microservice

Status: completed with `notification-service`.

The added service provides:

- `GET /health`
- `GET /metrics`
- `POST /notifications`
- A fake notification result such as `{"status":"sent","channel":"email","message_id":"demo-notification"}`
- Docker Compose, API gateway, Prometheus, README, and frontend wiring.

### 2. Add Real Request Metrics

Status: completed across the API gateway and all backend services.

The services expose middleware metrics for latency, error rate, request rate, and request success rate:

- `http_requests_total{service,method,path,status}`
- `http_request_duration_seconds_bucket{service,method,path}`
- `http_request_duration_seconds_sum`
- `http_request_duration_seconds_count`

Prometheus queries for SLI/SLO evidence:

- Availability: `avg_over_time(up[5m])`
- Latency p95: `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))`
- Error rate: `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`
- Success rate: `sum(rate(http_requests_total{status!~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`

### 3. Create SLI/SLO Documentation

Add a document such as `docs/sli-slo.md` with:

- Availability SLO: at least 99%.
- Latency SLO: p95 less than or equal to 200 ms.
- Error rate SLO: less than or equal to 1%.
- Request success rate target.
- Measurement source: Prometheus.
- Exact PromQL queries.
- Burn-rate or alerting explanation if you want stronger SRE evidence.

### 4. Add Kubernetes Manifests

Create a `k8s/` directory with manifests for:

- Namespace.
- ConfigMaps for service URLs and database config.
- Secret for PostgreSQL password.
- Deployments for frontend, API gateway, all microservices, PostgreSQL, Prometheus, and Grafana.
- Services for each component.
- Optional Ingress for frontend/API access.
- HorizontalPodAutoscaler for at least `api-gateway`, `order-service`, and `payment-service`.

Minimum acceptable structure:

```text
k8s/
  namespace.yaml
  configmap.yaml
  secret.yaml
  postgres.yaml
  api-gateway.yaml
  auth-service.yaml
  user-service.yaml
  product-service.yaml
  order-service.yaml
  payment-service.yaml
  notification-service.yaml
  frontend.yaml
  prometheus.yaml
  grafana.yaml
  hpa.yaml
```

### 5. Strengthen Docker Swarm Evidence

The current Compose file can be a base, but the assignment asks for Swarm as a separate orchestration approach.

Add either:

- A `docker-stack.yml` file with `deploy.replicas`, resource limits, update config, and restart policy.
- Or a `docs/docker-swarm.md` guide showing `docker swarm init`, `docker stack deploy`, scaling commands, and screenshots.

Recommended replicas:

- `api-gateway`: 2
- `product-service`: 2
- `order-service`: 2
- `payment-service`: 2
- `notification-service`: 2

### 6. Add Ansible Automation

Create an `ansible/` directory with:

```text
ansible/
  inventory.ini.example
  playbook.yml
  roles/
    docker/
    deploy/
    monitoring/
```

Minimum tasks:

- Install Docker and Docker Compose plugin.
- Install Git.
- Clone or update the repository on the VM.
- Start the app with Docker Compose or Docker Swarm.
- Optionally install Kubernetes tools if demonstrating Kubernetes on the VM.
- Verify health endpoints after deployment.

This closes the current gap where Terraform creates the VM but setup and deployment are manual.

### 7. Write Incident Report And Postmortem

Add `docs/incident-postmortem.md`.

It should include:

- Incident title.
- Date and duration.
- Impact: order creation unavailable.
- Detection: Prometheus alert or failed health check.
- Root cause: incorrect `DB_HOST` for `order-service`.
- Timeline: when failure was introduced, detected, diagnosed, fixed, and verified.
- Recovery: fix config and restart service.
- What went well.
- What went poorly.
- Action items: config validation, alert improvement, runbook, automated checks.

The existing README incident simulation can become the runbook, but it is not a full postmortem.

### 8. Add Capacity Planning Evidence

Add `docs/capacity-planning.md`.

Include:

- Simple load test method, for example `hey`, `wrk`, or `k6`.
- Request rate tested.
- Observed latency and error rate.
- CPU/memory observations from Docker stats, Grafana, or Kubernetes metrics.
- Bottleneck discussion: PostgreSQL and order/payment paths.
- Scaling plan: horizontal service replicas, database tuning, resource requests/limits, HPA thresholds.

Even a small local test is better than only theoretical text.

### 9. Improve Alerts

Current alerts only detect scrape failure. Add alerts for SLO-related symptoms:

- High 5xx error rate.
- High p95 latency.
- Low request success rate.
- PostgreSQL or order-service health failure.

This requires the request metrics described above.

### 10. Add Demo Evidence

The final PDF should include screenshots of:

- Docker Compose running containers.
- Frontend working.
- API Gateway requests working.
- Prometheus targets page.
- Grafana dashboard.
- Alert firing during the incident.
- Terraform apply/output or Azure VM.
- Ansible playbook run.
- Kubernetes pods/services.
- Docker Swarm stack/services.
- Incident simulation before and after recovery.

## Suggested Implementation Order

1. Add `notification-service` and wire it into Compose, gateway, Prometheus, and frontend if needed.
2. Add request metrics middleware to API gateway and services.
3. Add SLI/SLO document and update Prometheus alerts.
4. Add `k8s/` manifests.
5. Add `docker-stack.yml` and Swarm documentation.
6. Add Ansible playbook for VM setup and deployment.
7. Add incident postmortem and capacity planning documents.
8. Run the app locally, collect screenshots, and build the final PDF report.

## Files To Be Careful With

The repo currently has ignored Terraform local files in the workspace, including `terraform.tfvars` and Terraform state files. They are correctly covered by `.gitignore`, but do not include them in the final Git submission or PDF screenshots if they expose subscription, IP, or SSH-related data.

## Final Recommendation

Keep this repository. It is already about 50-60% aligned with the assignment as a technical base. To become submission-ready, it needs the missing orchestration, automation, documentation, and evidence layers listed above.
