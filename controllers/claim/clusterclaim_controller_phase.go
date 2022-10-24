package controllers

import (
	"context"
	"os"

	claimV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/claim/v1alpha1"
	clusterV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/cluster/v1alpha1"
	"github.com/tmax-cloud/hypercloud-multi-operator/controllers/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ClusterClaimReconciler) CreateClusterManager(ctx context.Context, cc *claimV1alpha1.ClusterClaim) (ctrl.Result, error) {
	if cc.Status.Phase != claimV1alpha1.ClusterClaimPhaseApproved {
		return ctrl.Result{}, nil
	}

	log := r.Log.WithValues("clusterclaim", cc.GetNamespacedName())
	log.Info("Start to reconcile phase for CreateClusterManager")

	key := types.NamespacedName{
		Name:      cc.Spec.ClusterName,
		Namespace: cc.Namespace,
	}

	if err := r.CheckSshKeySecretExist(cc); err != nil {
		log.Error(err, "Check if ssh key secret exist")
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
	}

	if err := r.Get(context.TODO(), key, &clusterV1alpha1.ClusterManager{}); errors.IsNotFound(err) {
		newClusterManager := &clusterV1alpha1.ClusterManager{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      cc.Spec.ClusterName,
				Namespace: cc.Namespace,
				Labels: map[string]string{
					clusterV1alpha1.LabelKeyClmClusterType: clusterV1alpha1.ClusterTypeCreated,
					clusterV1alpha1.LabelKeyClcName:        cc.Name,
				},
				Annotations: map[string]string{
					"owner":                                cc.Annotations[util.AnnotationKeyCreator],
					"creator":                              cc.Annotations[util.AnnotationKeyCreator],
					clusterV1alpha1.AnnotationKeyClmSuffix: util.CreateSuffixString(),
					clusterV1alpha1.AnnotationKeyClmDomain: os.Getenv(util.HC_DOMAIN),
				},
			},
			Spec: clusterV1alpha1.ClusterManagerSpec{
				Provider:  cc.Spec.Provider,
				Version:   cc.Spec.Version,
				MasterNum: cc.Spec.MasterNum,
				WorkerNum: cc.Spec.WorkerNum,
			},
			AwsSpec: clusterV1alpha1.ProviderAwsSpec{
				Region:         cc.Spec.ProviderAwsSpec.Region,
				SshKeyName:     cc.Spec.ProviderAwsSpec.SshKeyName,
				Bastion:        clusterV1alpha1.Instance(cc.Spec.ProviderAwsSpec.Bastion),
				Master:         clusterV1alpha1.Instance(cc.Spec.ProviderAwsSpec.Master),
				Worker:         clusterV1alpha1.Instance(cc.Spec.ProviderAwsSpec.Worker),
				HostOS:         cc.Spec.ProviderAwsSpec.HostOS,
				NetworkSpec:    clusterV1alpha1.NetworkSpec(cc.Spec.ProviderAwsSpec.NetworkSpec),
				AdditionalTags: cc.Spec.ProviderAwsSpec.AdditionalTags,
			},
			VsphereSpec: clusterV1alpha1.ProviderVsphereSpec{
				PodCidr:             cc.Spec.ProviderVsphereSpec.PodCidr,
				VcenterIp:           cc.Spec.ProviderVsphereSpec.VcenterIp,
				VcenterId:           cc.Spec.ProviderVsphereSpec.VcenterId,
				VcenterPassword:     cc.Spec.ProviderVsphereSpec.VcenterPassword,
				VcenterThumbprint:   cc.Spec.ProviderVsphereSpec.VcenterThumbprint,
				VcenterNetwork:      cc.Spec.ProviderVsphereSpec.VcenterNetwork,
				VcenterDataCenter:   cc.Spec.ProviderVsphereSpec.VcenterDataCenter,
				VcenterDataStore:    cc.Spec.ProviderVsphereSpec.VcenterDataStore,
				VcenterFolder:       cc.Spec.ProviderVsphereSpec.VcenterFolder,
				VcenterResourcePool: cc.Spec.ProviderVsphereSpec.VcenterResourcePool,
				VcenterKcpIp:        cc.Spec.ProviderVsphereSpec.VcenterKcpIp,
				VcenterCpuNum:       cc.Spec.ProviderVsphereSpec.VcenterCpuNum,
				VcenterMemSize:      cc.Spec.ProviderVsphereSpec.VcenterMemSize,
				VcenterDiskSize:     cc.Spec.ProviderVsphereSpec.VcenterDiskSize,
				VcenterTemplate:     cc.Spec.ProviderVsphereSpec.VcenterTemplate,
			},
		}

		if err := r.Create(context.TODO(), newClusterManager); err != nil {
			return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
		}

	} else if err != nil {
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
	}
	return ctrl.Result{}, nil
}
