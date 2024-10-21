// Copyright (c) 2024 Red Hat, Inc.

package webhooks

import (
	"context"

	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-nim-opendatahub-io-v1alpha1-odhnimaccount,mutating=false,failurePolicy=fail,groups=nim.opendatahub.io,resources=odhnimaccounts,versions=v1alpha1,name=validate.nim.opendatahub.io.v1alpha1.odhnimaccount,sideEffects=None,admissionReviewVersions=v1

type OdhNIMAccountValidator struct {
	client.Client
}

// SetupWithManager is used for setting up the webhook with a manager (check the init function)
func (w *OdhNIMAccountValidator) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.OdhNIMAccount{}).
		WithValidator(w).
		Complete()
}

func (w *OdhNIMAccountValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	// TODO Only the ODH NIM Operator is allowed to Create or Delete OdhNIMAccount
	//return w.verifyOnlyOneInNamespace(ctx, obj)
	return nil
}

func (w *OdhNIMAccountValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	return nil
}

func (w *OdhNIMAccountValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}
