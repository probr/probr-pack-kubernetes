provider "helm" {
  debug = true
  kubernetes {
    host     = azurerm_kubernetes_cluster.cluster.kube_config.0.host
    username = azurerm_kubernetes_cluster.cluster.kube_config.0.username
    password = azurerm_kubernetes_cluster.cluster.kube_config.0.password

    client_certificate     = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
    client_key             = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)

    load_config_file = false
  }
}

resource "helm_release" "aad-pod-identity" {
  namespace  = "kube-system"
  name       = "aad-pod-identity"
  repository = "https://raw.githubusercontent.com/Azure/aad-pod-identity/master/charts"
  chart      = "aad-pod-identity"
  version    = "2.0.2"
}

resource "azurerm_user_assigned_identity" "probe_msi" {
  resource_group_name = azurerm_kubernetes_cluster.cluster.node_resource_group
  location            = azurerm_resource_group.rg.location

  name = var.probr_probe_msi_name
}

data "azurerm_subscription" "this" {}

data "template_file" "azureidentity" {
  template = file("./azure-identity.yaml")
  vars = {
    subscription_id = data.azurerm_subscription.this.id
    node_rg_name    = azurerm_kubernetes_cluster.cluster.node_resource_group
    msi_name        = var.probr_probe_msi_name
    msi_object_id   = azurerm_user_assigned_identity.probe_msi.principal_id
  }
}

data "template_file" "azureidentitybinding" {
  template = file("./azure-identity-binding.yaml")
}

resource "null_resource" "azureidentity_apply" {
  depends_on = [null_resource.kubectl, helm_release.aad-pod-identity]
  provisioner "local-exec" {
    command     = "cat <<EOL | kubectl apply -n default --kubeconfig=${var.kube_config_filepath} -f - \n${data.template_file.azureidentity.rendered}\nEOL"
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "null_resource" "azureidentitybinding_apply" {
  depends_on = [null_resource.kubectl, helm_release.aad-pod-identity]

  provisioner "local-exec" {
    command     = "cat <<EOL | kubectl apply -n default --kubeconfig=${var.kube_config_filepath} -f - \n${data.template_file.azureidentitybinding.rendered}\nEOL"
    interpreter = ["/bin/bash", "-c"]
  }
}
