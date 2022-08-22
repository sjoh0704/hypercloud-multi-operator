#!/bin/sh
terraform -chdir=/kubespray/contrib/terraform/aws/ apply -auto-approve --state=/context/terraform.tfstate
