provider "kubernetes" {
  host                   = azurerm_kubernetes_cluster.cluster.kube_config.0.host
  username               = azurerm_kubernetes_cluster.cluster.kube_config.0.username
  password               = azurerm_kubernetes_cluster.cluster.kube_config.0.password
  client_certificate     = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
  client_key             = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)
  load_config_file       = "false"
}

resource "kubernetes_namespace" "probr" {
  count = length(var.namespaces)

  metadata {
    name = var.namespaces[count.index]

    labels = {
      app = "probr"
    }
  }
}

resource "kubernetes_network_policy" "deny_egress" {
  depends_on = [kubernetes_namespace.probr]
  count      = length(var.namespaces)

  metadata {
    name      = "probr-deny-egress"
    namespace = var.namespaces[count.index]
  }

  spec {
    pod_selector {}
    policy_types = ["Egress"]
  }
}
