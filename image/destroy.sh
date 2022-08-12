#!/bin/sh
terraform -chdir=$PWD/contrib/terraform/aws/ destroy -auto-approve --state=/context/terraform.tfstate
