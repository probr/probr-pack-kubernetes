resource "azurerm_policy_assignment" "psp" {
  name                 = var.psp_policy_name
  scope                = azurerm_resource_group.rg.id
  policy_definition_id = "/providers/Microsoft.Authorization/policySetDefinitions/42b8ef37-b724-4e24-bbc8-7a7708edfe00"
  description          = "Restricted PSP policy"
  display_name         = "Restricted PSP policy"

  parameters = <<PARAMETERS
{
  "effect": {
    "value": "deny"
  }
}
PARAMETERS

}

resource "azurerm_policy_assignment" "restrict_container_registries" {
  name                 = var.restrict_registry_policy_name
  scope                = azurerm_resource_group.rg.id
  policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/febd0533-8e55-448f-b837-bd0e06f16469"
  description          = "Restrict container registry"
  display_name         = "Restrict container registry"

  parameters = <<PARAMETERS
{
  "allowedContainerImagesRegex": {
    "value": "^.+.azurecr.io/.+$|mcr.microsoft.com/oss/azure/aad-pod-identity"
  },
  "effect": {
    "value": "deny"
  }
}
PARAMETERS

}
