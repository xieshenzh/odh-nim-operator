// Copyright (c) 2024 Red Hat, Inc.

package v1alpha1

import (
	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	OdhNIMRuntimeSpec struct {
		InferenceServiceSpec kservev1beta1.InferenceServiceSpec `json:"inferenceServiceSpec"`
		PvcSpec              corev1.PersistentVolumeClaimSpec   `json:"pvcSpec"`
		OdhNIMAccountRef     corev1.ObjectReference             `json:"odhNIMAccountRef"`
	}

	OdhNIMRuntimeStatus struct {
		PvcRef              corev1.ObjectReference `json:"pvcRef,omitempty"`
		ImagePullSecretRef  corev1.ObjectReference `json:"imagePullSecretRef,omitempty"`
		NimSecretRef        corev1.ObjectReference `json:"nimSecretRef,omitempty"`
		ServingRuntimeRef   corev1.ObjectReference `json:"servingRuntimeRef,omitempty"`
		InferenceServiceRef corev1.ObjectReference `json:"inferenceServiceRef,omitempty"`
		// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
		Conditions []metav1.Condition `json:"conditions,omitempty"`
	}

	// OdhNIMRuntime is used for creating NIM Kserve-based deployment for Open Data Hub.
	//
	// +kubebuilder:object:root=true
	// +kubebuilder:resource:shortName=onimr
	// +kubebuilder:subresource:status
	// +kubebuilder:printcolumn:name="PVC",type="string",JSONPath=".status.conditions[?(@.type==\"PVCReconciled\")].status",description="The status of PVC reconciliation"
	// +kubebuilder:printcolumn:name="Pull Secret",type="string",JSONPath=".status.conditions[?(@.type==\"ImagePullSecretReconciled\")].status",description="The status of Pull Secret reconciliation"
	// +kubebuilder:printcolumn:name="NIM Secret",type="string",JSONPath=".status.conditions[?(@.type==\"NimSecretReconciled\")].status",description="The status of NIM Secret reconciliation"
	// +kubebuilder:printcolumn:name="ServingRuntime",type="string",JSONPath=".status.conditions[?(@.type==\"ServingRuntimeReconciled\")].status",description="The status of ServingRuntime reconciliation"
	// +kubebuilder:printcolumn:name="InferenceService",type="string",JSONPath=".status.conditions[?(@.type==\"InferenceReconciledCreation\")].status",description="The status of InferenceService reconciliation"
	// +operator-sdk:csv:customresourcedefinitions:displayName="ODH NIM Runtime"
	// +operator-sdk:csv:customresourcedefinitions:resources={{PersistentVolumeClaim,v1,nim-pvc},{Secret,v1,ngc-secret},{Secret,v1,nvidia-nim-secrets},{ConfigMap,v1,nvidia-nim-data},{ServingRuntime,serving.kserve.io/v1alpha1},{InferenceService,serving.kserve.io/v1beta1}}
	OdhNIMRuntime struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
		Spec              OdhNIMRuntimeSpec `json:"spec"`
		// +kubebuilder:validation:Optional
		Status OdhNIMRuntimeStatus `json:"status"`
	}

	// OdhNIMRuntimeList is used for encapsulating OdhNIMRuntime items.
	//
	// +kubebuilder:object:root=true
	OdhNIMRuntimeList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
		Items           []OdhNIMRuntime `json:"items"`
	}
)

func init() {
	SchemeBuilder.Register(&OdhNIMRuntime{}, &OdhNIMRuntimeList{})
}
