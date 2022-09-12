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
	ArgoReadyCondition                      ConditionType = "ArgoReady"
	TraefikReadyCondition                   ConditionType = "TraefikReady"
	GatewayReadyCondition                   ConditionType = "GatewayReady"
	AuthClientReadyCondition                ConditionType = "AuthClientReady"
)

type ConditionReason string

const (
	K8sInstallingStartedReason                 ConditionReason = "K8sInstallingStarted"
	K8sInstallingReconciliationFailedReason    ConditionReason = "K8sInstallReconciliationFailed"
	K8sInstallingReconciliationSucceededReason ConditionReason = "K8sInstallReconciliationSucceeded"
)

const (
	InfrastructureProvisioningStartedReason                 ConditionReason = "InfrastructureProvisioningStarted"
	InfrastructureProvisioningReconciliationFailedReason    ConditionReason = "InfrastructureProvisioningReconciliationFailed"
	InfrastructureProvisioningReconciliationSucceededReason ConditionReason = "InfrastructureProvisioningReconciliationSucceeded"
)

const (
	KubeconfigCreatingStartedReason                 ConditionReason = "KubeconfigCreatingStarted"
	KubeconfigCreatingReconciliationFailedReason    ConditionReason = "KubeconfigCreatingReconciliationFailed"
	KubeconfigCreatingReconciliationSucceededReason ConditionReason = "KubeconfigCreatingReconciliationSucceeded"
)

const (
	VolumeSettingStartedReason                 ConditionReason = "VolumeSettingStarted"
	VolumeSettingReconciliationFailedReason    ConditionReason = "VolumeSettingReconciliationFailed"
	VolumeSettingReconciliationSucceededReason ConditionReason = "VolumeSettingReconciliationSucceeded"
)

const (
	ControlplaneReadyReason    ConditionReason = "ControlplaneReady"
	ControlplaneNotReadyReason ConditionReason = "ControlplaneNotReady"
)

const (
	ArgoReadyReason    ConditionReason = "ArgoReady"
	ArgoNotReadyReason ConditionReason = "ArgoNotReady"
)

const (
	TraefikReadyReason    ConditionReason = "TraefikReady"
	TraefikNotReadyReason ConditionReason = "TraefikNotReady"
)

const (
	GatewayReadyReason    ConditionReason = "GatewayReady"
	GatewayNotReadyReason ConditionReason = "GatewayNotReady"
)

const (
	AuthClientReadyReason    ConditionReason = "AuthClientReady"
	AuthClientNotReadyReason ConditionReason = "AuthClientNotReady"
)
