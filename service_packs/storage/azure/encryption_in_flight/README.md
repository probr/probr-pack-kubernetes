# Encryption in Flight Probe Notes

This directory contains the feature file and code realted to the probing of encryption in flight controls

## Mandatory Azure Configuration Variables

- ***AZURE_SUBSCRIPTION_ID*** - the user supplied azure subscription id
- ***AZURE_TENANT_ID*** - the user supplied azure tenant id
- ***AZURE_CLIENT_ID*** - the user supplied azure client id (will n ormally be a service principal application id)
- ***AZURE_CLIENT_SECRET*** - the secret required for client authentication
- ***AZURE_RESOURCE_GROUP*** - the user supplied resource group for Probr purposes and must exist in the specified subscription
- ***AZURE_LOCATION*** - the azure geo location where test storage account resources may be created

## Azure Policy prerequiste

A policy which denies the creation of storage accounts with non-secure http access enabled, must be assigned to the user's azure subscription or azure management group. The applicable built-in azure policy is:
`Secure transfer to storage accounts should be enabled`. The assignment must set the 'Effect' parameter value to 'Deny', in order to prevent creation of storage accounts with the EnableHTTPSTrafficOnly option not set to true. Note that the default value is 'Audit', which will not prevent non-compliant account creation.

## Preventative scenario outline

Probr attempts to create a storage account for the following scenarios:

- http and https access is switched on - creation should be denied
- only http access is switched on - creation should be denied
- only https access is switched on - creation should be allowed