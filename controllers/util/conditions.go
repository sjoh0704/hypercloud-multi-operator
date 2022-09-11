package util

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CheckConditionExist(conditions []metav1.Condition, conditionType string) bool {
	if meta.FindStatusCondition(conditions, conditionType) != nil {
		return true
	}
	return false
}

// condition이 있는지 먼저 체크 후, condtion 값이 false라면 true를 반환한다.
func CheckConditionExistAndConditionFalse(conditions []metav1.Condition, conditionType string) bool {

	if !CheckConditionExist(conditions, conditionType) {
		return false
	}

	if meta.IsStatusConditionFalse(conditions, conditionType) {
		return true
	}
	return false
}

// condition이 있는지 먼저 체크 후, condtion 값이 true라면 true를 반환한다.
func CheckConditionExistAndConditionTrue(conditions []metav1.Condition, conditionType string) bool {
	return meta.IsStatusConditionTrue(conditions, conditionType)
}
