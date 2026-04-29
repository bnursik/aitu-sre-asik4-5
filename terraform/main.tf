terraform {
  required_version = ">= 1.5.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.110"
    }
  }
}

provider "azurerm" {
  features {}
}

locals {
  common_tags = {
    Project   = var.project_name
    ManagedBy = "terraform"
  }

  ssh_public_key = var.ssh_public_key != "" ? var.ssh_public_key : file(pathexpand(var.ssh_public_key_path))
}

resource "azurerm_resource_group" "app" {
  name     = "${var.project_name}-rg"
  location = var.location
  tags     = local.common_tags
}

resource "azurerm_virtual_network" "app" {
  name                = "${var.project_name}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.app.location
  resource_group_name = azurerm_resource_group.app.name
  tags                = local.common_tags

  subnet {
    name           = "${var.project_name}-subnet"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_public_ip" "app" {
  name                = "${var.project_name}-public-ip"
  location            = azurerm_resource_group.app.location
  resource_group_name = azurerm_resource_group.app.name
  allocation_method   = "Static"
  sku                 = "Standard"
  tags                = local.common_tags
}

resource "azurerm_network_security_group" "app" {
  name                = "${var.project_name}-nsg"
  location            = azurerm_resource_group.app.location
  resource_group_name = azurerm_resource_group.app.name
  tags                = local.common_tags

  security_rule {
    name                       = "Allow-SSH"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = var.allowed_ssh_cidr
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Allow-HTTP"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = var.allowed_web_cidr
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Allow-Grafana"
    priority                   = 120
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3000"
    source_address_prefix      = var.allowed_web_cidr
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Allow-Prometheus"
    priority                   = 130
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "9090"
    source_address_prefix      = var.allowed_web_cidr
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface" "app" {
  name                = "${var.project_name}-nic"
  location            = azurerm_resource_group.app.location
  resource_group_name = azurerm_resource_group.app.name
  tags                = local.common_tags

  ip_configuration {
    name                          = "internal"
    subnet_id                     = one(azurerm_virtual_network.app.subnet).id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.app.id
  }
}

resource "azurerm_network_interface_security_group_association" "app" {
  network_interface_id      = azurerm_network_interface.app.id
  network_security_group_id = azurerm_network_security_group.app.id
}

resource "azurerm_linux_virtual_machine" "app" {
  name                  = "${var.project_name}-vm"
  location              = azurerm_resource_group.app.location
  resource_group_name   = azurerm_resource_group.app.name
  size                  = var.vm_size
  admin_username        = var.admin_username
  network_interface_ids = [azurerm_network_interface.app.id]
  tags                  = local.common_tags

  disable_password_authentication = true

  admin_ssh_key {
    username   = var.admin_username
    public_key = local.ssh_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }
}
