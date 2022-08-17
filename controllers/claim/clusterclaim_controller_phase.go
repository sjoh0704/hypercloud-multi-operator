package controllers

import (
	"context"
	"fmt"
	"os"

	claimV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/claim/v1alpha1"
	clusterV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/cluster/v1alpha1"
	"github.com/tmax-cloud/hypercloud-multi-operator/controllers/util"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ClusterClaimReconciler) CreateClusterManager(ctx context.Context, cc *claimV1alpha1.ClusterClaim) (ctrl.Result, error) {
	if cc.Status.Phase != "Approved" {
		return ctrl.Result{}, nil
	}

	log := r.Log.WithValues("clusterclaim", cc.GetNamespacedName())
	log.Info("Start to reconcile phase for CreateClusterManagerPhase")

	key := types.NamespacedName{
		Name:      cc.Spec.ClusterName,
		Namespace: cc.Namespace,
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
					clusterV1alpha1.AnnotationKeyClmDomain: os.Getenv("HC_DOMAIN"),
					clusterV1alpha1.AnnotationKeyClmSuffix: util.CreateSuffixString(),
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

func (r *ClusterClaimReconciler) CreatePersistentVolumeClaim(ctx context.Context, cc *claimV1alpha1.ClusterClaim) (ctrl.Result, error) {
	if cc.Status.Phase != "Approved" {
		return ctrl.Result{}, nil
	}
	storageClassName := ""

	log := r.Log.WithValues("clusterclaim", cc.GetNamespacedName())
	log.Info("Start to reconcile phase for CreatePersistentVolumeClaimPhase")

	key := types.NamespacedName{
		Name:      cc.Spec.ClusterName,
		Namespace: cc.Namespace,
	}

	// clm의 suffix값을 이용해서 pvc name을 지정해준다.
	clm := &clusterV1alpha1.ClusterManager{}
	err := r.Get(context.TODO(), key, clm)
	if errors.IsNotFound(err) {
		log.Info("Wait for creating cluster manager")
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, nil
	} else if err != nil {
		log.Error(err, "Fail to get cluster manager")
		return ctrl.Result{}, nil
	}

	key = types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", cc.Spec.ClusterName, clm.Annotations[clusterV1alpha1.AnnotationKeyClmSuffix]),
		Namespace: cc.Namespace,
	}

	if err := r.Get(context.TODO(), key, &coreV1.PersistentVolumeClaim{}); errors.IsNotFound(err) {
		pvc := &coreV1.PersistentVolumeClaim{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", cc.Spec.ClusterName, clm.Annotations[clusterV1alpha1.AnnotationKeyClmSuffix]),
				Namespace: cc.Namespace,
				Labels: map[string]string{
					clusterV1alpha1.LabelKeyClmName: cc.Spec.ClusterName,
				},
			},
			Spec: coreV1.PersistentVolumeClaimSpec{
				StorageClassName: &storageClassName,
				AccessModes: []coreV1.PersistentVolumeAccessMode{
					coreV1.ReadWriteOnce,
				},
				Resources: coreV1.ResourceRequirements{
					Limits: coreV1.ResourceList{},
					Requests: coreV1.ResourceList{
						coreV1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}

		if err := r.Create(context.TODO(), pvc); err != nil {
			return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
		}

	} else if err != nil {
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
	}

	return ctrl.Result{}, nil
}
func (r *ClusterClaimReconciler) CreatePersistentVolume(ctx context.Context, cc *claimV1alpha1.ClusterClaim) (ctrl.Result, error) {
	if cc.Status.Phase != "Approved" {
		return ctrl.Result{}, nil
	}

	log := r.Log.WithValues("clusterclaim", cc.GetNamespacedName())
	log.Info("Start to reconcile phase for CreatePersistentVolumePhase")

	key := types.NamespacedName{
		Name:      cc.Spec.ClusterName,
		Namespace: cc.Namespace,
	}

	// clm의 suffix값을 이용해서 pvc name을 지정해준다.
	clm := &clusterV1alpha1.ClusterManager{}
	err := r.Get(context.TODO(), key, clm)
	if errors.IsNotFound(err) {
		log.Info("Wait for creating cluster manager")
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, nil
	} else if err != nil {
		log.Error(err, "Fail to get cluster manager")
		return ctrl.Result{}, nil
	}

	key = types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", cc.Spec.ClusterName, clm.Annotations[clusterV1alpha1.AnnotationKeyClmSuffix]),
		Namespace: cc.Namespace,
	}

	if err := r.Get(context.TODO(), key, &coreV1.PersistentVolume{}); errors.IsNotFound(err) {
		pvc := &coreV1.PersistentVolume{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", cc.Spec.ClusterName, clm.Annotations[clusterV1alpha1.AnnotationKeyClmSuffix]),
				Namespace: cc.Namespace,
				Labels: map[string]string{
					clusterV1alpha1.LabelKeyClmName: cc.Spec.ClusterName,
				},
			},
			Spec: coreV1.PersistentVolumeSpec{
				Capacity: coreV1.ResourceList{
					coreV1.ResourceStorage: resource.MustParse("1Gi"),
				},
				AccessModes: []coreV1.PersistentVolumeAccessMode{
					coreV1.ReadWriteOnce,
				},
				ClaimRef: &coreV1.ObjectReference{
					Name:      fmt.Sprintf("%s-%s", cc.Spec.ClusterName, clm.Annotations[clusterV1alpha1.AnnotationKeyClmSuffix]),
					Namespace: cc.Namespace,
				},
				PersistentVolumeReclaimPolicy: coreV1.PersistentVolumeReclaimRetain,
			},
		}

		if err := r.Create(context.TODO(), pvc); err != nil {
			return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
		}

	} else if err != nil {
		return ctrl.Result{RequeueAfter: requeueAfter10Second}, err
	}
	return ctrl.Result{}, nil
}
