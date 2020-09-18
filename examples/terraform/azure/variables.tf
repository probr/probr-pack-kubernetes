variable "prefix" {
    default     = "example-k8s"
    description = "Value that will be prepended to most others"
}

variable "resource_group_name" {
    default     = "resource-group"
    description = "K8s resource group. Should recieve prefix."
}

variable "location" {
    default     = "East US 2"
    description = "Location display name. az account list-locations --output table"
}

variable "cluster_name" {
    default     = "cluster"
    description = "K8s cluster. Should recieve prefix."
}

variable "kube_config" {
    default     = ""
    description = "Filepath for kube config to be written to"
}
