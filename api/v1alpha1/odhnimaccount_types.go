// Copyright (c) 2024 Red Hat, Inc.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	// OdhNIMAccountSpec defines the desired state of an OdhNIMAccount object.
	OdhNIMAccountSpec struct {
		// A reference to the Secret containing the NGC API Key.
		SecretRef corev1.ObjectReference `json:"secretRef"`
	}

	// OdhNIMAccountStatus defines the observed state of an OdhNIMAccount object.
	OdhNIMAccountStatus struct {
		// A reference to the Template for NIM ServingRuntime.
		TemplateRef *corev1.ObjectReference `json:"templateRef,omitempty"`
		// A reference to the ConfigMap with data for NIM deployment.
		ConfigMapRef *corev1.ObjectReference `json:"configMapRef,omitempty"`

		// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
		Conditions []metav1.Condition `json:"conditions,omitempty"`
	}

	// OdhNIMAccount is used for adopting a NIM Account for Open Data Hub.
	//
	// +kubebuilder:object:root=true
	// +kubebuilder:resource:shortName=onima
	// +kubebuilder:subresource:status
	//
	// +kubebuilder:printcolumn:name="Template",type="string",JSONPath=".status.templateRef.name",description="The name of the Template"
	// +kubebuilder:printcolumn:name="ConfigMap",type="string",JSONPath=".status.configMapRef.name",description="The name of the ConfigMap"
	//
	// +operator-sdk:csv:customresourcedefinitions:displayName="ODH NIM Account"
	// +operator-sdk:csv:customresourcedefinitions:resources={{ConfigMap,v1},{Template,template.openshift.io/v1}}
	OdhNIMAccount struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Spec   OdhNIMAccountSpec   `json:"spec,omitempty"`
		Status OdhNIMAccountStatus `json:"status,omitempty"`
	}

	// OdhNIMAccountList is used for encapsulating OdhNIMAccount items.
	//
	// +kubebuilder:object:root=true
	OdhNIMAccountList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
		Items           []OdhNIMAccount `json:"items"`
	}
)

func init() {
	SchemeBuilder.Register(&OdhNIMAccount{}, &OdhNIMAccountList{})
}
