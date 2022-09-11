/*
Copyright 2020 The Kubernetes Authors.

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

type ConditionType string

const (
	K8sInstalledReadyCondition              ConditionType = "K8sInstalledReady"
	InfrastructureProvisionedReadyCondition ConditionType = "InfrastructureProvisionedReady"
	KubeconfigCreatedReadyCondition         ConditionType = "KubeconfigCreatedReady"
	VolumeReadyCondition                    ConditionType = "VolumeReady"
	ControlplaneReadyCondition              ConditionType = "ControlplaneReady"
)

const (
	K8sInstallingStartedReason                 = "K8sInstallingStarted"
	K8sInstallingReconciliationFailedReason    = "K8sInstallReconciliationFailed"
	K8sInstallingReconciliationSucceededReason = "K8sInstallReconciliationSucceeded"
)

const (
	InfrastructureProvisioningStartedReason                 = "InfrastructureProvisioningStarted"
	InfrastructureProvisioningReconciliationFailedReason    = "InfrastructureProvisioningReconciliationFailed"
	InfrastructureProvisioningReconciliationSucceededReason = "InfrastructureProvisioningReconciliationSucceeded"
)

const (
	KubeconfigCreatingStartedReason                 = "KubeconfigCreatingStarted"
	KubeconfigCreatingReconciliationFailedReason    = "KubeconfigCreatingReconciliationFailed"
	KubeconfigCreatingReconciliationSucceededReason = "KubeconfigCreatingReconciliationSucceeded"
)

const (
	VolumeSettingStartedReason                 = "VolumeSettingStarted"
	VolumeSettingReconciliationFailedReason    = "VolumeSettingReconciliationFailed"
	VolumeSettingReconciliationSucceededReason = "VolumeSettingReconciliationSucceeded"
)

const (
	ControlplaneReady                          = "ControlplaneReady"
	ControlplaneNotReady                       = "ControlplaneNotReady"
)

