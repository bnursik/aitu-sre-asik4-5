output "public_ip" {
  description = "Public IP address of the Azure VM."
  value       = azurerm_public_ip.app.ip_address
}

output "ssh_command" {
  description = "SSH command for connecting to the Azure VM."
  value       = "ssh ${var.admin_username}@${azurerm_public_ip.app.ip_address}"
}

output "app_url" {
  description = "HTTP application URL after manual deployment."
  value       = "http://${azurerm_public_ip.app.ip_address}"
}

output "prometheus_url" {
  description = "Prometheus URL after manual deployment."
  value       = "http://${azurerm_public_ip.app.ip_address}:9090"
}

output "grafana_url" {
  description = "Grafana URL after manual deployment."
  value       = "http://${azurerm_public_ip.app.ip_address}:3000"
}
