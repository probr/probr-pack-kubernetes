provider "azurerm" {
  version = "~>2.5.0"
  features {} // azurerm will err if this is not included

# Authentication uses defaults from azure cli (az account list); 
# these fields can override those defaults
#   subscription_id = "00000000-0000-0000-0000-000000000000"
#   tenant_id       = "11111111-1111-1111-1111-111111111111"

}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "${var.prefix}-${var.cluster_name}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "${var.prefix}-dns"
  tags = {
    aadpodidentity: "enabled",
    policies: "all",
    project: "probr"
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    vnet_subnet_id = element(tolist(azurerm_virtual_network.vnet.subnet), 0).id // subnet object contains one value
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
  }

  addon_profile {
    azure_policy {
      enabled = true
    }

    http_application_routing {
      enabled = false
    }

    kube_dashboard {
      enabled = false
    }
    
    oms_agent {
      enabled                    = true
      log_analytics_workspace_id = azurerm_log_analytics_workspace.awp.id
    }
  }
}

resource "null_resource" "kubectl" {
  provisioner "local-exec" {
    command = "echo '${azurerm_kubernetes_cluster.cluster.kube_config_raw}' > .kubeconfig"
    interpreter = ["/bin/bash", "-c"]
  }
  provisioner "local-exec" {
    command = "kubectl apply --kubeconfig=${var.kube_config} -f https://raw.githubusercontent.com/Azure/aad-pod-identity/master/deploy/infra/deployment-rbac.yaml"
    interpreter = ["/bin/bash", "-c"]
  }
}

output "client_certificate" {
  value = azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.cluster.kube_config_raw
}
