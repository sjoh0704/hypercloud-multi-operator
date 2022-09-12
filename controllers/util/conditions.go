package util

import (
	clusterV1alpha1 "github.com/tmax-cloud/hypercloud-multi-operator/apis/cluster/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CheckConditionExist(conditions []metav1.Condition, conditionType clusterV1alpha1.ConditionType) bool {
	if meta.FindStatusCondition(conditions, string(conditionType)) != nil {
		return true
	}
	return false
}

// condition이 있는지 먼저 체크 후, condtion 값이 false라면 true를 반환한다.
func CheckConditionExistAndConditionFalse(conditions []metav1.Condition, conditionType clusterV1alpha1.ConditionType) bool {
	if !CheckConditionExist(conditions, conditionType) {
		return false
	}

	if meta.IsStatusConditionFalse(conditions, string(conditionType)) {
		return true
	}
	return false
}

// condition이 있는지 먼저 체크 후, condtion 값이 true라면 true를 반환한다.
func CheckConditionExistAndConditionTrue(conditions []metav1.Condition, conditionType clusterV1alpha1.ConditionType) bool {
	return meta.IsStatusConditionTrue(conditions, string(conditionType))
}
