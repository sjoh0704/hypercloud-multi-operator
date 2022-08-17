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

package controllers

import (
	"context"

	claimV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/claim/v1alpha1"
	clusterV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/cluster/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ClusterClaimReconciler) requeueClusterClaimsForClusterManager(o client.Object) []ctrl.Request {
	clm := o.DeepCopyObject().(*clusterV1alpha1.ClusterManager)
	log := r.Log.WithValues("objectMapper", "clusterManagerToClusterClaim", "clusterManager", clm.Name)
	log.Info("Start to clusterManagerToClusterClaim mapping...")

	//get clusterManager
	cc := &claimV1alpha1.ClusterClaim{}
	key := types.NamespacedName{
		Name:      clm.Labels[clusterV1alpha1.LabelKeyClcName],
		Namespace: clm.Namespace,
	}
	if err := r.Get(context.TODO(), key, cc); errors.IsNotFound(err) {
		log.Info("ClusterClaim resource not found. Ignoring since object must be deleted")
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get ClusterClaim")
		return nil
	}

	if cc.Status.Phase != "Approved" {
		log.Info("ClusterClaims for ClusterManager [" + cc.Spec.ClusterName + "] is already delete... Do not update cc status to delete ")
		return nil
	}

	cc.Status.Phase = "ClusterDeleted"
	cc.Status.Reason = "cluster is deleted"
	err := r.Status().Update(context.TODO(), cc)
	if err != nil {
		log.Error(err, "Failed to update ClusterClaim status")
		return nil //??
	}
	return nil
}
