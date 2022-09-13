package util

import (
	"fmt"

	clusterV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/cluster/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
)

func CreateEnvFromClustermanagerSpec(clusterManager *clusterV1alpha1.ClusterManager) ([]coreV1.EnvVar, error) {
	EnvList := []coreV1.EnvVar{}
	AwsSpec := clusterManager.AwsSpec

	// region 종속 파라미터를 가져온다.
	params, err := GetAwsRegionalPreset(AwsSpec.Region, AwsSpec.HostOS)
	if err != nil {
		return nil, err
	}

	// region 별 종속 데이터를 환경변수로 지정한다.
	EnvList = append(EnvList,
		coreV1.EnvVar{
			Name:  "TF_VAR_AWS_DEFAULT_REGION",
			Value: fmt.Sprintf("%s", AwsSpec.Region),
		},
		coreV1.EnvVar{
			Name:  "HOST_OS",
			Value: fmt.Sprintf("%s", AwsSpec.HostOS),
		},
		coreV1.EnvVar{
			Name:  "TF_VAR_AWS_SSH_KEY_NAME",
			Value: fmt.Sprintf("%s", params.SshKeyName),
		},
		coreV1.EnvVar{
			Name:  "TF_VAR_aws_ami_name",
			Value: fmt.Sprintf("[\"%s\"]", params.AmiName),
		},
		coreV1.EnvVar{
			Name:  "TF_VAR_aws_ami_owner",
			Value: fmt.Sprintf("[\"%s\"]", params.AmiOwner),
		},
		coreV1.EnvVar{
			Name:  "USER",
			Value: fmt.Sprintf("%s", params.User),
		},
	)

	// cluster name, master num, worker num
	EnvList = append(EnvList,
		coreV1.EnvVar{
			Name:  "TF_VAR_aws_cluster_name",
			Value: fmt.Sprintf("%s", clusterManager.Name),
		},
		coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_master_num",
			Value: fmt.Sprintf("%d", clusterManager.Spec.MasterNum),
		},
		coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_worker_num",
			Value: fmt.Sprintf("%d", clusterManager.Spec.WorkerNum),
		},
	)

	// bastion // default 1
	if AwsSpec.Bastion.Num > 0 {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_bastion_num",
			Value: fmt.Sprintf("%d", AwsSpec.Bastion.Num),
		})
	}

	if AwsSpec.Bastion.Type != "" {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_bastion_size",
			Value: fmt.Sprintf("%s", AwsSpec.Bastion.Type),
		})
	}

	// master
	if AwsSpec.Master.Type != "" {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_master_size",
			Value: fmt.Sprintf("%s", AwsSpec.Master.Type),
		})
	}

	if AwsSpec.Master.DiskSize != 0 {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_master_disk_size",
			Value: fmt.Sprintf("%d", AwsSpec.Master.DiskSize),
		})
	}

	// worker
	if AwsSpec.Worker.Type != "" {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_worker_size",
			Value: fmt.Sprintf("%s", AwsSpec.Worker.Type),
		})
	}

	if AwsSpec.Worker.DiskSize != 0 {
		EnvList = append(EnvList, coreV1.EnvVar{
			Name:  "TF_VAR_aws_kube_worker_disk_size",
			Value: fmt.Sprintf("%d", AwsSpec.Worker.DiskSize),
		})
	}

	if AwsSpec.NetworkSpec.VpcCidrBlock != "" {
		if len(AwsSpec.NetworkSpec.PrivateSubnetCidrBlock) != len(AwsSpec.NetworkSpec.PublicSubnetCidrBlock) {
			return nil, fmt.Errorf("PrivateSubnetCidrBlock and PublicSubnetCidrBlock must have same length of list")
		}

		publicCidr := ""
		for _, cidr := range AwsSpec.NetworkSpec.PublicSubnetCidrBlock {
			publicCidr += fmt.Sprintf("\"%s\", ", cidr)
		}

		privateCidr := ""
		for _, cidr := range AwsSpec.NetworkSpec.PrivateSubnetCidrBlock {
			privateCidr += fmt.Sprintf("\"%s\", ", cidr)
		}

		EnvList = append(EnvList,
			coreV1.EnvVar{
				Name:  "TF_VAR_aws_vpc_cidr_block",
				Value: fmt.Sprintf("%s", AwsSpec.NetworkSpec.VpcCidrBlock),
			},
			coreV1.EnvVar{
				Name:  "TF_VAR_aws_cidr_subnets_public",
				Value: fmt.Sprintf("[%s]", publicCidr[:len(publicCidr)-2]),
			},
			coreV1.EnvVar{
				Name:  "TF_VAR_aws_cidr_subnets_private",
				Value: fmt.Sprintf("[%s]", privateCidr[:len(privateCidr)-2]),
			})
	}

	// additional tags
	if len(AwsSpec.AdditionalTags) > 0 {
		addtionalTags := ""
		for key, value := range AwsSpec.AdditionalTags {
			addtionalTags += fmt.Sprintf("%s=\"%s\", ", key, value)
		}
		EnvList = append(EnvList,
			coreV1.EnvVar{
				Name:  "TF_VAR_default_tags",
				Value: fmt.Sprintf("{ %s }", addtionalTags[:len(addtionalTags)-2]),
			})
	}

	return EnvList, nil
}

func GetTerminatedPodStateReason(pod *coreV1.Pod) (string, error) {
	if pod == nil {
		return "", fmt.Errorf("pod doesn't exist")
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == clusterV1alpha1.Kubespray {
			if containerStatus.State.Terminated != nil {
				return containerStatus.State.Terminated.Message, nil
			}
			return "", fmt.Errorf("not yet terminated")
		}
	}
	return "", fmt.Errorf("not found kubespray container")
}

func GetAwsRegionalPreset(region string, os string) (*AwsParams, error) {
	if !IsAwsRegion(region) {
		return nil, fmt.Errorf("%s is not existing region", region)
	}
	if !IsAwsSupportOs(os) {
		return nil, fmt.Errorf("%s is not supporting os", region)
	}
	params := &AwsParams{}
	if region == AP_NORTHEAST_2 {
		params.SshKeyName = "ap-northeast-2-mo"
		if os == UBUNTU {
			params.AmiName = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20220609"
			params.AmiOwner = "099720109477"
			params.User = "ubuntu"

		} else if os == RHEL {
			params.AmiName = "RHEL-8.2.0_HVM-20210907-x86_64-0-Hourly2-GP2"
			params.AmiOwner = "309956199498"
			params.User = "ec2-user"
		}

	} else if region == EU_WEST_1 {
		params.SshKeyName = "eu-west-1-mo"
		if os == UBUNTU {
			params.AmiName = "ami-ubuntu-18.04-1.13.0-00-1548773800"
			params.AmiOwner = "258751437250"
			params.User = "ubuntu"
		} else if os == RHEL {
			params.AmiName = "RHEL-8.2.0_HVM-20210907-x86_64-0-Hourly2-GP2"
			params.AmiOwner = "309956199498"
			params.User = "ec2-user"
		}
	}
	return params, nil
}

func IsAwsRegion(s string) bool {
	for _, r := range AwsRegionSet {
		if r == s {
			return true
		}
	}
	return false
}

func IsAwsSupportOs(s string) bool {
	for _, r := range AwsSupportOs {
		if r == s {
			return true
		}
	}
	return false
}
