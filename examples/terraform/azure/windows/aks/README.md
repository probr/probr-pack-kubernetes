 # Azure Terraform

This Terraform build is the create a Windows version of the AKS for development purposes.

#Prequistes for development on Windows

- Install Terraform on your local machine, follow the instructions here https://www.terraform.io/downloads.html

- Install kubectl on you your local machine, follow the instructions here
https://kubernetes.io/docs/tasks/tools/install-kubectl/

- Download code to your local machine. 

Running Terraform modules
- Navigate to probr/examples/terraform/azure/windows/aks/

- Open the terraform.tfvars configuration file

- Remove example Azure subscription entries for your own Azure subcription entries.

- Make sure you are in the probr/examples/terraform/azure/windows/aks/ directory.

- Run the command terraform plan, review plan

- Run the command 
 ~~~
 terraform apply -var-file="terraform.tfvars" -auto-approve
 ~~~
 Make sure there are no errors.
 
- Run
 ~~~
 az aks get-credentials --resource-group <ResourceGroup Name> --name  <AKS Cluster Name> --overwrite-existing
 ~~~ 
 this will allow to login to your AKS cluster.
 ~~~
 Example: az aks get-credentials --resource-group probr-automation-rg --name probr-automation-cluster --overwrite-existing
 ~~~
  
- Once connected to your AKS cluster, run the command:
 ~~~
 kubectl get pods --all-namespaces, review pods
 ~~~








