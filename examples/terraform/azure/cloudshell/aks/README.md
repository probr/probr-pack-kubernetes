<h1> Azure Terraform </h1>

This Terraform code will run in Azure CLI bash in Cloudshell. This Terraform build is the create a cloudshell version of the AKS for development purposes.

- To install Azure CLI Cloudshell, please follow the directions here https://docs.microsoft.com/en-us/cli/azure/install-azure-cli

- Please navigate your storage account where a fileshare has been created for your.  

- This would be storage account -> fileshare, and you fileshare name.

- On upper bar click Add Directory, please name the directory what you like, this is where you will upload you Terraform code.

- After the directory is created, click Upload, upload the code in cloudshell directory.

- Open cloudshell and navigate to you clouddrive where you will see the your uploaded code. 

- Make sure your Cloudshell is configured for bash

- Run 
  ~~~
  az login
  ~~~
  to login to Azure 

- Run 
  ~~~
  az account list
  ~~~
  to list subscriptions if necessary

- Run 
  ~~~
  az account set -s <YOUR SUBSCRIPTION>
  ~~~
- Run 
  ~~~
  terraform init
  ~~~
- Run 
  ~~~
  terraform plan
  ~~~
- Run
  ~~~
  terraform apply -var-file="variables.tfvars" -auto-approve
  ~~~
- Run 
  ~~~
  az aks get-credentials --resource-group <ResourceGroup Name> --name  <AKS Cluster Name>
  ~~~
- Run 
  ~~~
  kubectl get pods --all-namespaces
  ~~~
