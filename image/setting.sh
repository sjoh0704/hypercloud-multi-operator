#!/bin/sh
echo "setup starts"

HOST_FILE_PATH=/kubespray/inventory/tmaxcloud/hosts
VAR_PATH=/kubespray/inventory/tmaxcloud/group_vars
ELB_FILE_PATH=$VAR_PATH/all/all.yml
EFS_FILE_PATH=$VAR_PATH/k8s_cluster/addons.yml
CALICO_FILE_PATH=$VAR_PATH/k8s_cluster/k8s-net-calico.yml
OFFLINE_FILE_PATH=$VAR_PATH/all/offline.yml
KUBECONFIG_SETTING_FILE_PATH=$VAR_PATH/k8s_cluster/k8s-cluster.yml
CLUSTER_FILE_PATH=/kubespray/cluster.yml
NOT_USED="#NOT_USED#"

echo "copy hosts file from context file to $HOST_FILE_PATH"
cp /context/hosts $HOST_FILE_PATH

echo "fill etcd host in $HOST_FILE_PATH"
ipset=`sed -n "/\[kube_control_plane\]/,/\[kube_node\]/p" $HOST_FILE_PATH | head -n -1 | tail -n +2`
for i in $ipset
do 
   sed -i "/\[etcd\]/a"${i}"" $HOST_FILE_PATH
done

echo "set up loadbalancer in $ELB_FILE_PATH"
apiserver_lb_domain_name=`grep apiserver_loadbalancer_domain_name $HOST_FILE_PATH | cut -d "=" -f 2`
sed -i "s/^apiserver_loadbalancer_domain_name:.*$/apiserver_loadbalancer_domain_name: "${apiserver_lb_domain_name}"/" $ELB_FILE_PATH
sed -i "s/  address:.*$/${NOT_USED} address:/" $ELB_FILE_PATH

echo "set up elastic file system in $EFS_FILE_PATH"
aws_efs_filesystem_id=`grep aws_efs_filesystem_id $HOST_FILE_PATH | cut -d "=" -f 2`
sed -i "s/^default_storageclass_name:.*$/default_storageclass_name: efs-sc/" $EFS_FILE_PATH
sed -i "s/^nfs_external_provisioner_enabled:.*$/nfs_external_provisioner_enabled: false/" $EFS_FILE_PATH
sed -i "s/^aws_efs_csi_enabled:.*$/aws_efs_csi_enabled: true/" $EFS_FILE_PATH
sed -i "s/^aws_efs_filesystem_id:.*$/aws_efs_filesystem_id: "${aws_efs_filesystem_id}"/" $EFS_FILE_PATH

echo "set up calico in $CALICO_FILE_PATH"
sed -i "s/^#.*calico_ipip_mode:.*$/calico_ipip_mode: 'Always'/" $CALICO_FILE_PATH
calico_cidr=$(echo $TF_VAR_aws_cidr_subnets_private | cut -d "=" -f2)
length=${#calico_cidr}
calico_cidr="${calico_cidr:2:length-3}" 
calico_cidr=`echo $calico_cidr | sed "s/\"//g" | sed "s/, /,/g"`  
sed -i "/^calico_ip_auto_method:.*$/ccalico_ip_auto_method: \"cidr="${calico_cidr}"\"" $CALICO_FILE_PATH

echo "set up offline option false in $OFFLINE_FILE_PATH"
sed -i "s/^is_this_offline:.*$/is_this_offline: false/g" $OFFLINE_FILE_PATH
sed -i "s/^registry_host:.*$/registry_host: \"\"/g" $OFFLINE_FILE_PATH
# TODO: 수정 필요
#del_line=`nl $OFFLINE_FILE_PATH | grep files_repo: | awk '{print $1}'`
#echo $del_line
sed -i '1,5!d' $OFFLINE_FILE_PATH

# TODO 검토 필요
echo "delete bootstrap-os in $CLUSTER_FILE_PATH"
sed -i '/bootstrap-os/d' $CLUSTER_FILE_PATH  

echo "setup kubeconfig option true in $KUBECONFIG_SETTING_FILE_PATH"
sed -i "s/^#.*kubeconfig_localhost:.*$/kubeconfig_localhost: true/g" $KUBECONFIG_SETTING_FILE_PATH

echo "setup completes"