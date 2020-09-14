provider "azurerm" {
  version = "~>2.5.0"
  features {} // azurerm will err if this is not included

# Authentication uses defaults from azure cli (az account list); 
# these fields can override those defaults
#   subscription_id = "00000000-0000-0000-0000-000000000000"
#   tenant_id       = "11111111-1111-1111-1111-111111111111"

}
resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-${var.resource_group_name}"
  location = var.location
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-${var.cluster_name}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "${var.prefix}"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  addon_profile {
    aci_connector_linux {
      enabled = false
    }

    azure_policy {
      enabled = false
    }

    http_application_routing {
      enabled = false
    }

    kube_dashboard {
      enabled = true
    }

    oms_agent {
      enabled = false
    }
  }
}

output "client_certificate" {
  value = azurerm_kubernetes_cluster.example.kube_config.0.client_certificate
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.example.kube_config_raw
}
