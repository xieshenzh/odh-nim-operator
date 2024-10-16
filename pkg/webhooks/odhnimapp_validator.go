// Copyright (c) 2024 Red Hat, Inc.

package webhooks

import (
	"context"
	"fmt"
	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-nim-opendatahub-io-v1alpha1-odhnimapp,mutating=false,failurePolicy=fail,groups=nim.opendatahub.io,resources=odhnimapps,versions=v1alpha1,name=validate.nim.opendatahub.io.v1alpha1.odhnimapp,sideEffects=None,admissionReviewVersions=v1

type OdhNimAppValidator struct {
	client.Client
}

// SetupWithManager is used for setting up the webhook with a manager (check the init function)
func (w *OdhNimAppValidator) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&v1alpha1.OdhNimApp{}).WithValidator(w).Complete()
}

func (w *OdhNimAppValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	// TODO Only the ODH NIM Operator is allowed to Create or Delete OdhNimApp
	return w.verifyOnlyOneInNamespace(ctx, obj)
}

func (w *OdhNimAppValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	// TODO Users can only Update the OdhNimApp.Spec{.ApiKey.Validate | .Content.Update } keys triggering validation or
	// TODO content fetch, any other spec keys can only be updated by the ODH NIM Operator
	return nil
}

func (w *OdhNimAppValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	// TODO Only the ODH NIM Operator is allowed to Create or Delete OdhNimApp
	return nil
}

func (w *OdhNimAppValidator) verifyOnlyOneInNamespace(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx).WithName("odhnimapp-validator-webhook")

	ns := obj.(*v1alpha1.OdhNimApp).Namespace

	ocfgs := &metav1.PartialObjectMetadataList{}
	ocfgs.SetGroupVersionKind(v1alpha1.GroupVersion.WithKind("OdhNimAppList"))
	if err := w.Client.List(ctx, ocfgs, client.InNamespace(ns)); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if len(ocfgs.Items) > 0 {
		err := fmt.Errorf("an OdhNimApp instance already exists in %s", ns)
		logger.V(1).Error(err, "only one OdhNimApp is allowed in a namespace")
		return err
	}
	logger.V(1).Info(fmt.Sprintf("no OdhNimApp instances found in %s", ns))
	return nil
}
