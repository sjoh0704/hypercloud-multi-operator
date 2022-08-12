#!/bin/sh
HOST_FILE=/kubespray/inventory/tmaxcloud/hosts

# setting
./setting.sh

# 특이사항 
adduser $USER --disabled-password

# ansible-playbook
ansible-playbook -i $HOST_FILE ./cluster.yml -e ansible_user=$USER -e bootstrap_os=$HOST_OS -e ansible_ssh_private_key_file=/kubespray/key.pem -e cloud_provider=aws -b --become-user=root --flush-cache -v

# admin.conf 파일 저장
cp /kubespray/inventory/tmaxcloud/artifacts/admin.conf /context/admin.conf

# admin.conf 파일 update
API_SERVER_LB_DNS=`grep apiserver_loadbalancer_domain_name $HOST_FILE | cut -d "=" -f 2`
API_SERVER_LB_DNS="${API_SERVER_LB_DNS%\"}"
API_SERVER_LB_DNS="${API_SERVER_LB_DNS#\"}"
sed -i "s/^.*server:.*$/    server: https:\/\/"${API_SERVER_LB_DNS}":6443/g" /context/admin.conf
chmod 777 /context/admin.conf