variable "azure_subscription"{
  description = "Azure subscription to use"
}

variable "prefix" {
  description = "Value that will be prepended to most others"
}

variable "location" {
  description = "Location display name. az account list-locations --output table"
}

variable "cluster_name" {
  description = "K8s cluster. Should recieve prefix."
}

variable "kube_config_filepath" {
  description = "Filepath for kube config to be written to"
}

variable "demo_acr" {
  default = "automation"
}

variable "acr_name" {
  
}

variable "probr_probe_msi_name" {

}

variable "psp_policy_name" {
  description = "Restricted PSP policy"
}

variable "restrict_registry_policy_name" {
  description = "Restrict container registry"
}


variable "namespaces" {
  type = list

  default = [
    "probr-container-access-test-ns",
    "probr-general-test-ns",
    "probr-network-access-test-ns",
    "probr-pod-security-test-ns",
    "probr-rbac-test-ns"
  ]

}
