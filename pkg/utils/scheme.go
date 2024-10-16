// Copyright (c) 2024 Red Hat, Inc.

package utils

import (
	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// InstallTypes is used for installing our required types with a given scheme.
func InstallTypes(scheme *runtime.Scheme) error {
	installs := []func(*runtime.Scheme) error{
		v1alpha1.Install,   // our own api
		corev1.AddToScheme, // ConfigMaps and Secrets
	}

	for _, install := range installs {
		if err := install(scheme); err != nil {
			return err
		}
	}

	return nil
}
