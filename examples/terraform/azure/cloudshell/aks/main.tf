provider "azurerm" {
  version = "~>2.35.0"
  features {} // azurerm will err if this is not included
  subscription_id = var.azure_subscription
}



resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-rg"
  location = var.location
}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "${var.prefix}-${var.cluster_name}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "${var.prefix}-dns"
  tags = {
    aadpodidentity : "enabled",
    policies : "all",
    project : "automation demo"
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    //vnet_subnet_id = element(tolist(azurerm_virtual_network.vnet.subnet), 0).id // subnet object contains one value
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "azure"
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
      enabled = false
      //log_analytics_workspace_id = azurerm_log_analytics_workspace.awp.id
    }
  }
}


resource "azurerm_container_registry" "acr" {
  name                     = var.acr_name
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  sku                      = "Standard"
  admin_enabled            = false

}

resource null_resource "probrimage" {
  provisioner "local-exec" {
      command = "az acr import -n ${var.acr_name} --source docker.io/citihub/probr-probe"
      
  }
}

resource "null_resource" "kubectl" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command     = "echo '${azurerm_kubernetes_cluster.cluster.kube_config_raw}' > ${var.kube_config_filepath}"
    interpreter = ["/bin/bash", "-c"]
  }
}
