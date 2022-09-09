#!/bin/sh
# state file 경로 
STATE_FILE_PATH=/context/terraform.tfstate

# terraform 실행 경로 
TF_FILE_PATH=/kubespray/contrib/terraform/aws/

# terraform 실행 
terraform -chdir=$TF_FILE_PATH destroy -auto-approve --state=$STATE_FILE_PATH
