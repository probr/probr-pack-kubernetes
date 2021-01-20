
# This script will add environment variables with global scope in Windows.
# Launch a Terminal and "Run as Administrator".
# Restart VS Code and check values were applied: go run .\examples\getenvvars.go

# Azure credentials
$AzureTenantID = Read-Host -Prompt 'Enter value for AZURE_TENANT_ID'
[System.Environment]::SetEnvironmentVariable('AZURE_TENANT_ID', $AzureTenantID, [System.EnvironmentVariableTarget]::Machine)

$AzureSubscriptionID = Read-Host -Prompt 'Enter value for AZURE_SUBSCRIPTION_ID'
[System.Environment]::SetEnvironmentVariable('AZURE_SUBSCRIPTION_ID', $AzureSubscriptionID, [System.EnvironmentVariableTarget]::Machine)

$AzureClientID = Read-Host -Prompt 'Enter value for AZURE_CLIENT_ID'
[System.Environment]::SetEnvironmentVariable('AZURE_CLIENT_ID', $AzureClientID, [System.EnvironmentVariableTarget]::Machine)

$AzureClientSecret = Read-Host -Prompt 'Enter value for AZURE_CLIENT_SECRET'
[System.Environment]::SetEnvironmentVariable('AZURE_CLIENT_SECRET', $AzureClientSecret, [System.EnvironmentVariableTarget]::Machine)