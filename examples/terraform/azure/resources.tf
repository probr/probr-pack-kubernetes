resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-${var.resource_group_name}"
  location = var.location
}

resource "random_id" "law_id" {
    byte_length = 8
    prefix      = "${var.prefix}-log-analytics-"
}

resource "azurerm_log_analytics_workspace" "awp" {
    # The WorkSpace name has to be unique across the whole of azure, not just the current subscription/tenant.
    name                = random_id.law_id.dec
    location            = var.location
    resource_group_name = azurerm_resource_group.rg.name
    sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "as" {
    solution_name         = "ContainerInsights"
    location              = azurerm_log_analytics_workspace.awp.location
    resource_group_name   = azurerm_resource_group.rg.name
    workspace_resource_id = azurerm_log_analytics_workspace.awp.id
    workspace_name        = azurerm_log_analytics_workspace.awp.name

    plan {
        publisher = "Microsoft"
        product   = "OMSGallery/ContainerInsights"
    }
}

resource "azurerm_virtual_network" "vnet" {
  name                = "${var.prefix}-network"
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.location
  address_space       = ["10.1.2.0/24"]

  subnet {
    name           = "${var.prefix}-subnet"
    address_prefix = "10.1.2.0/25"
  }
}