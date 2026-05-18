# Azure Terraform Deployment

This folder provisions the infrastructure required by the assignment using AzureRM.

Terraform creates infrastructure only. Docker installation and application deployment are automated with the Ansible workflow in [`../ansible/`](../ansible/).

## What It Creates

- Azure resource group
- Virtual network
- Subnet
- Static public IP
- Network security group
- Network interface
- Ubuntu Linux VM using `azurerm_linux_virtual_machine`

Opened inbound ports:

- SSH: `22`
- HTTP: `80`
- Grafana: `3000`
- Prometheus: `9090`

## Prerequisites

- Terraform installed
- Azure CLI installed
- Azure account/subscription
- SSH public key, for example `~/.ssh/id_ed25519.pub`

Login to Azure:

```bash
az login
```

Select the correct subscription if needed:

```bash
az account set --subscription "<subscription-id>"
```

## Configure Variables

Copy the example file:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:

```hcl
project_name        = "computer-shop-microservices"
location            = "North Europe"
vm_size             = "Standard_D2s_v3"
admin_username      = "azureuser"
ssh_public_key_path = "~/.ssh/azure_rsa.pub"
ssh_public_key      = ""
allowed_ssh_cidr    = "YOUR_PUBLIC_IP/32"
allowed_web_cidr    = "0.0.0.0/0"
```

By default Terraform reads the public key from `ssh_public_key_path`. If you prefer, you can paste the public key directly into `ssh_public_key` instead.

Do not commit `terraform.tfvars`.

## Run Terraform

```bash
terraform init
terraform fmt
terraform validate
terraform plan
terraform apply
terraform output
```

## Outputs

Terraform outputs:

- `public_ip`
- `ssh_command`
- `app_url`
- `prometheus_url`
- `grafana_url`

## Deploy With Ansible After VM Creation

Copy and edit the Ansible inventory:

```bash
cp ../ansible/inventory.ini.example ../ansible/inventory.ini
```

Set `ansible_host` to the Terraform `public_ip` output and set `ansible_repo_url` to your Git repository URL.

Run the playbook:

```bash
ansible-playbook -i ../ansible/inventory.ini ../ansible/playbook.yml
```

The playbook installs Docker, installs Git, clones or updates the repository, writes a VM Compose file with ports matching the Terraform outputs, starts the app with Docker Compose, and verifies the frontend, Prometheus, and Grafana health endpoints.

Open:

```text
App:        http://<public_ip>
Grafana:    http://<public_ip>:3000
Prometheus: http://<public_ip>:9090
```

## Cleanup

```bash
terraform destroy
```
