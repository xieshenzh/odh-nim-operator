// Copyright (c) 2024 Red Hat, Inc.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// +groupName=nim.opendatahub.io
// +kubebuilder:object:generate=true
// +kubebuilder:validation:Required

var (
	GroupVersion  = schema.GroupVersion{Group: "nim.opendatahub.io", Version: "v1alpha1"}
	schemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	Install       = schemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}

type (
	OdhNimAppSpecApiKey struct {
		// +kubebuilder:default=true
		// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
		Validate bool `json:"validate"`
		// +kubebuilder:validation:Optional
		// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
		SecretRef *corev1.ObjectReference `json:"secretRef,omitempty"`
	}

	OdhNimAppSpecContent struct {
		// +kubebuilder:default=true
		// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
		Update bool `json:"update"`
		// +kubebuilder:validation:Optional
		// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
		ConfigMapRef *corev1.ObjectReference `json:"configMapRef,omitempty"`
	}

	OdhNimAppSpec struct {
		ApiKey  OdhNimAppSpecApiKey  `json:"apiKey"`
		Content OdhNimAppSpecContent `json:"content"`
		// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
		TemplateRef *corev1.ObjectReference `json:"templateRef"`
	}

	OdhNimAppStatus struct {
		// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
		Conditions []metav1.Condition `json:"conditions,omitempty"`
	}

	// OdhNimApp is used for activating NIM integration reconciliation in Open Data Hub.
	//
	// +kubebuilder:object:root=true
	// +kubebuilder:resource:shortName=ona
	// +kubebuilder:subresource:status
	// +kubebuilder:printcolumn:name="Validated",type="string",JSONPath=".status.conditions[?(@.type==\"ApiKeyValidated\")].status",description="The validation status of the API Key"
	// +kubebuilder:printcolumn:name="Updated",type="string",JSONPath=".status.conditions[?(@.type==\"ContentUpdated\")].status",description="The status of the last content update"
	// +operator-sdk:csv:customresourcedefinitions:displayName="ODH NIM App"
	// +operator-sdk:csv:customresourcedefinitions:resources={{OdhNimApp,nim.opendatahub.io/v1alpha1},{PersistentVolumeClaim,v1,nim-pvc},{Secret,v1,ngc-secret},{Secret,v1,nvidia-nim-secrets},{ConfigMap,v1,odh-nim-app-content},{Template,template.openshift.io/v1,nvidia-nim-serving-template}}
	OdhNimApp struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
		Spec              OdhNimAppSpec `json:"spec"`
		// +kubebuilder:validation:Optional
		Status OdhNimAppStatus `json:"status,omitempty"`
	}

	// OdhNimAppList is used for encapsulating OdhNimApp items.
	//
	// +kubebuilder:object:root=true
	OdhNimAppList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
		Items           []OdhNimApp `json:"items"`
	}
)

func init() {
	schemeBuilder.Register(&OdhNimApp{}, &OdhNimAppList{})
}
