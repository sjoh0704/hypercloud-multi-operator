/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ClusterClaimPhase string

const (
	// 클러스터 클레임이 생성되고, 관리자의 승인/거절을 기다리는 상태
	ClusterClaimPhaseAwaiting = ClusterClaimPhase("Awaiting")
	// 관리자에 의해 클레임이 승인된 상태
	ClusterClaimPhaseApproved = ClusterClaimPhase("Approved")
	// 관리자에 의헤 클레임이 거절된 상태
	ClusterClaimPhaseRejected = ClusterClaimPhase("Rejected")
	// 클러스터가 삭제된 상태
	ClusterClaimPhaseClusterDeleted = ClusterClaimPhase("Cluster Deleted")
	// 클러스터 생성과정에서 에러가 발생한 상태
	ClusterClaimPhaseError = ClusterClaimPhase("Error")
)

const (
	ClusterClaimDeprecatedPhaseClusterDeleted = ClusterClaimPhase("ClusterDeleted")
)

// ClusterClaimSpec defines the desired state of ClusterClaim
type ClusterClaimSpec struct {
	// +kubebuilder:validation:Required
	// The name of the cluster to be created
	ClusterName string `json:"clusterName"`
	// +kubebuilder:validation:Required
	// The version of kubernetes
	Version string `json:"version"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=AWS;vSphere
	// The type of provider
	Provider string `json:"provider"`
	// +kubebuilder:validation:Required
	// The number of master node
	MasterNum int `json:"masterNum"`
	// +kubebuilder:validation:Required
	// The number of worker node
	WorkerNum int `json:"workerNum"`
	// Provider Aws Spec
	ProviderAwsSpec AwsClaimSpec `json:"providerAwsSpec,omitempty"`
	// Provider vSphere Spec
	ProviderVsphereSpec VsphereClaimSpec `json:"providerVsphereSpec,omitempty"`
}

type AwsClaimSpec struct {
	// The ssh key secret name to access VM. Ssh key name and secret name must be same
	SshKeyName string `json:"sshKeyName,omitempty"`
	// +kubebuilder:validation:Enum:=ap-northeast-1;ap-northeast-2;ap-south-1;ap-southeast-1;ap-northeast-2;ca-central-1;eu-central-1;eu-west-1;eu-west-2;eu-west-3;sa-east-1;us-east-1;us-east-2;us-west-1;us-west-2;
	// The region where VM is working
	Region string `json:"region,omitempty"`
	// The info of bastion instance
	Bastion Instance `json:"bastion,omitempty"`
	// The info of master instance
	Master Instance `json:"master,omitempty"`
	// The info of worker instance
	Worker Instance `json:"worker,omitempty"`
	// +kubebuilder:validation:Enum:=rhel;ubuntu;
	// The type of OS that instances(master, worker, bastion) use
	HostOS string `json:"hostOs,omitempty"`
	// The network spec that cluster uses
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`
	// The additional tag attached to aws resources
	AdditionalTags map[string]string `json:"additionalTags,omitempty"`
}

type Instance struct {
	// Indicates the size of the instance
	Type string `json:"type,omitempty"`
	// The number of instance
	Num int `json:"num,omitempty"`
	// The disk size of the instance
	DiskSize int `json:"diskSize,omitempty"`
}

type NetworkSpec struct {
	// The size of the vpc to which the cluster will be deployed
	VpcCidrBlock string `json:"vpcCidrBlock,omitempty"`
	// The size of private subnet belonging to vpc cidr
	PrivateSubnetCidrBlock []string `json:"privateSubnetCidrBlock,omitempty"`
	// The size of public subnet belonging to vpc cidr
	PublicSubnetCidrBlock []string `json:"publicSubnetCidrBlock,omitempty"`
}

type VsphereClaimSpec struct {
	// The internal IP address cider block for pods
	PodCidr string `json:"podCidr,omitempty"`
	// The IP address of vCenter Server Application(VCSA)
	VcenterIp string `json:"vcenterIp,omitempty"`
	// The user id of VCSA
	VcenterId string `json:"vcenterId,omitempty"`
	// The password of VCSA
	VcenterPassword string `json:"vcenterPassword,omitempty"`
	// The TLS thumbprint of machine certificate
	VcenterThumbprint string `json:"vcenterThumbprint,omitempty"`
	// The name of network
	VcenterNetwork string `json:"vcenterNetwork,omitempty"`
	// The name of data center
	VcenterDataCenter string `json:"vcenterDataCenter,omitempty"`
	// The name of data store
	VcenterDataStore string `json:"vcenterDataStore,omitempty"`
	// The name of folder
	VcenterFolder string `json:"vcenterFolder,omitempty"`
	// The name of resource pool
	VcenterResourcePool string `json:"vcenterResourcePool,omitempty"`
	// The IP address of control plane for remote cluster(vip)
	VcenterKcpIp string `json:"vcenterKcpIp,omitempty"`
	// The number of cpus for vm
	VcenterCpuNum int `json:"vcenterCpuNum,omitempty"`
	// The memory size for vm, write as MB without unit. Example: 8192
	VcenterMemSize int `json:"vcenterMemSize,omitempty"`
	// The disk size for vm, write as GB without unit. Example: 25
	VcenterDiskSize int `json:"vcenterDiskSize,omitempty"`
	// The template name for cloud init
	VcenterTemplate string `json:"vcenterTemplate,omitempty"`
}

// ClusterClaimStatus defines the observed state of ClusterClaim
type ClusterClaimStatus struct {
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	Reason  string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`

	// +kubebuilder:validation:Enum=Awaiting;Admitted;Approved;Rejected;Error;ClusterDeleted;Cluster Deleted;
	Phase ClusterClaimPhase `json:"phase,omitempty" protobuf:"bytes,4,opt,name=phase"`
}

func (c *ClusterClaimStatus) SetTypedPhase(p ClusterClaimPhase) {
	c.Phase = p
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusterclaims,shortName=cc,scope=Namespaced
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.reason`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// ClusterClaim is the Schema for the clusterclaims API
type ClusterClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterClaimSpec   `json:"spec"`
	Status ClusterClaimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ClusterClaimList contains a list of ClusterClaim
type ClusterClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterClaim{}, &ClusterClaimList{})
}

func (cc *ClusterClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      cc.Name,
		Namespace: cc.Namespace,
	}
}
