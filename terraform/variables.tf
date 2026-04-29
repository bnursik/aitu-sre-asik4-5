variable "project_name" {
  description = "Name used for Azure resource names and tags."
  type        = string
  default     = "computer-shop-microservices"
}

variable "location" {
  description = "Azure region where resources will be created."
  type        = string
  default     = "North Europe"
}

variable "vm_size" {
  description = "Azure VM size."
  type        = string
  default     = "Standard_D2s_v3"
}

variable "admin_username" {
  description = "Admin username for SSH access to the Linux VM."
  type        = string
  default     = "azureuser"
}

variable "ssh_public_key_path" {
  description = "Path to the SSH public key file used for VM access."
  type        = string
  default     = "~/.ssh/azure_rsa.pub"
}

variable "ssh_public_key" {
  description = "SSH public key content. If empty, Terraform reads ssh_public_key_path."
  type        = string
  default     = ""
  sensitive   = true
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to SSH into the VM. Use your public IP with /32."
  type        = string
}

variable "allowed_web_cidr" {
  description = "CIDR allowed to access HTTP, Grafana, and Prometheus."
  type        = string
  default     = "0.0.0.0/0"
}
