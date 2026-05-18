# Ansible Deployment

This directory automates VM setup and application deployment for the Terraform-created Ubuntu host.

## Prerequisites

- Terraform has created the VM.
- Your SSH key can connect to the VM.
- Ansible is installed on your local machine.
- Your project is available from a Git repository URL.

## Configure Inventory

Copy the example inventory:

```bash
cp ansible/inventory.ini.example ansible/inventory.ini
```

Edit `ansible/inventory.ini`:

```ini
[app]
computer-shop-vm ansible_host=YOUR_VM_PUBLIC_IP ansible_user=azureuser

[app:vars]
ansible_repo_url=https://github.com/<your-user>/<your-repo>.git
ansible_app_dir=/opt/computer-shop-microservices
```

Optional: set `ansible_repo_version` in `[app:vars]` to deploy a branch, tag, or commit other than `main`.

## Run Deployment

```bash
ansible-playbook -i ansible/inventory.ini ansible/playbook.yml
```

The playbook installs Docker, the Docker Compose plugin, and Git, clones or updates the repository, writes a VM-only Compose file, starts the stack, and checks health endpoints.

## Expected URLs

```text
App:        http://<public_ip>
Grafana:    http://<public_ip>:3000
Prometheus: http://<public_ip>:9090
```

Default Grafana login:

```text
admin / admin
```

## Evidence Commands

Use these for screenshots or final report proof:

```bash
curl http://<public_ip>/health
curl http://<public_ip>:9090/-/healthy
curl http://<public_ip>:3000/api/health
ssh azureuser@<public_ip> "cd /opt/computer-shop-microservices && sudo docker compose -f docker-compose.vm.yml ps"
```

Prometheus targets:

```text
http://<public_ip>:9090/targets
```

Grafana dashboard:

```text
http://<public_ip>:3000/d/computer-shop-overview/computer-shop-overview
```
