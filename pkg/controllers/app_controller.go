// Copyright (c) 2024 Red Hat, Inc.

package controllers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type AppController struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager is used for setting up the controller with a manager (check the init function)
func (r *AppController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("odh-nim-app-controller").
		Complete(r)
}

// rbac markers are in controllers.go

func (r *AppController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("app-controller")
	ctx = log.IntoContext(ctx, logger)
	// all funcs we invoke in this context should use 'logger := log.FromContext(ctx)' to get the correct logger
	logger.V(1).Info(fmt.Sprintf("got request for Secret %s", req.NamespacedName))

	// TODO write code
	// 1. Fetch OdhNimApp
	// 2. If OdhNimApp NOT found (deleted):
	//		- Break reconciliation, we use the finalizer mechanism to cleanups
	//
	// 3. If in deletion process !OdhNimApp.DeletionTimestamp.IsZero() (note the !):
	//		3.1 If has our finalizer "nim.opendatahub.io/cleanup_finalizer" (const in controllers.go):
	//			- Tear down (if we want)
	//			- Remove our finalizer
	//		3.2 Break reconciliation
	//
	// 4. If doesn't have our finalizer "nim.opendatahub.io/cleanup_finalizer":
	//		- Add the finalizer
	//
	// 5. If OdhNimApp.Spec.ApiKey.Validate is True (defaults to true):
	//		5.1 Validate the API Key!!
	// 		5.2 Patch OdhNimApp.Status.Condition[Type=ApiKeyValidated] to True/False based on the validation (consts in controllers.go)
	//		5.3 Patch OdhNimApp.Spec.ApiKey.Validate to False (but store the original value for step 7)
	//		5.4 If the validation NOT successful:
	//			- Do we want to tear down?
	//			- Break reconciliation
	//
	// 6. Reconcile the Template, OdhNimApp.Spec.TemplateRef, if empty, create and patch the reference
	//
	// 7. If OdhNimApp.Spec.ApiKey.Validate WAS True (before 5.3) OR
	//	OdhNimApp.Spec.Content.Update IS True OR
	//	OdhNimApp.Spec.Content.ConfigMapRef is empty (we haven't created it yet):
	//		7.1 Fetch NIM Images and models
	//		7.2 Patch OdhNimApp.Status.Condition[Type=ContentUpdated] to True/False based on the fetch status (consts in controllers.go)
	//		7.3 Patch OdhNimApp.Spec.Content.Update to False (if wasn't false to begin with)
	//		7.4 If the fetching was successful:
	//			- Reconcile the content ConfigMap with the updated data
	//			- Patch the OdhNimApp.Spec.Content.ConfigMapRef with the reference to the ConfigMap if empty
	//
	// 8.  Reconcile a daily recurring Cron Job owned by this OdhNimApp, patching OdhNimApp.Spec.Content.Update to True

	return ctrl.Result{}, nil
}

// init is used for registering the odh-nim-app controller for loading
func init() {
	controllerSetups = append(controllerSetups, func(opts ControllerOptions) error {
		return (&SecretController{
			opts.Manager.GetClient(),
			opts.Manager.GetScheme(),
		}).SetupWithManager(opts.Manager)
	})
}
