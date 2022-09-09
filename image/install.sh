#!/bin/sh
HOST_FILE=/kubespray/inventory/tmaxcloud/hosts
ADMIN_CONF_FILE=/context/admin.conf

check_error () {
    if ! [ $? -eq 0 ]
    then
        echo "error occurs"
        exit 1
    fi
}

# setting
./setting.sh

# 특이사항 
adduser $USER --disabled-password
check_error

# ansible-playbook
ansible-playbook -i $HOST_FILE ./cluster.yml -e ansible_user=$USER -e bootstrap_os=$HOST_OS -e ansible_ssh_private_key_file=/kubespray/key.pem -e cloud_provider=aws -b --become-user=root --flush-cache -v
check_error

# admin.conf 파일 저장
cp /kubespray/inventory/tmaxcloud/artifacts/admin.conf $ADMIN_CONF_FILE
check_error

# admin.conf 파일 update
API_SERVER_LB_DNS=`grep apiserver_loadbalancer_domain_name $HOST_FILE | cut -d "=" -f 2`
API_SERVER_LB_DNS="${API_SERVER_LB_DNS%\"}"
API_SERVER_LB_DNS="${API_SERVER_LB_DNS#\"}"
sed -i "s/^.*server:.*$/    server: https:\/\/"${API_SERVER_LB_DNS}":6443/g" $ADMIN_CONF_FILE
check_error

chmod 777 $ADMIN_CONF_FILE
check_error
