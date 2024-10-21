// Copyright (c) 2024 Red Hat, Inc.

package controllers

import (
	"context"

	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=nim.opendatahub.io,resources=odhnimaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nim.opendatahub.io,resources=odhnimaccounts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nim.opendatahub.io,resources=odhnimaccounts/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="template.openshift.io",resources=templates,verbs=get;list;watch;create;update;patch;delete

type OdhNIMAccountController struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager is used for setting up the controller with a manager (check the init function)
func (r *OdhNIMAccountController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("odh-nim-account-controller").
		For(&v1alpha1.OdhNIMAccount{}).
		Complete(r)
}

func (r *OdhNIMAccountController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}
