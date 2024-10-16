// Copyright (c) 2024 Red Hat, Inc.

package controllers

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type SecretController struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager is used for setting up the controller with a manager (check the init function)
// Note the event filtering, we only watch Secrets with nim.opendatahub.io/nim-app set to true
func (r *SecretController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("odh-nim-secret-controller").
		For(&corev1.Secret{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(createEvent event.CreateEvent) bool {
				value, found := createEvent.Object.GetLabels()[Label_NimApp]
				return found && value == "true"
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				value, found := deleteEvent.Object.GetLabels()[Label_NimApp]
				return found && value == "true"
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				// TODO do we want to tear down when this labels is removed, if not, than we only need the new object
				valueOld, foundOld := updateEvent.ObjectOld.GetLabels()[Label_NimApp]
				valueNew, foundNew := updateEvent.ObjectOld.GetLabels()[Label_NimApp]
				return (foundOld && valueOld == "true") || (foundNew && valueNew == "true")
			},
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				value, found := genericEvent.Object.GetLabels()[Label_NimApp]
				return found && value == "true"
			},
		}).
		Complete(r)
}

// rbac markers are in controllers.go

func (r *SecretController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("secret-controller")
	ctx = log.IntoContext(ctx, logger)
	// all funcs we invoke in this context should use 'logger := log.FromContext(ctx)' to get the correct logger
	logger.V(1).Info(fmt.Sprintf("got request for Secret %s", req.NamespacedName))

	// TODO write code
	// 1. Fetch the Secret
	// 2. Fetch the OdhNimApp referencing this Secret by Name and Namespace (OdhNimApp.Spec.ApiKey.SecretRef)
	// 3. Fetch the Cron Job owned by this secret
	//
	// 4. If Secret NOT found (deleted):
	//		- Delete the OdhNimApp (if found)
	//		- Delete the CronJob (if found)
	//		- Break reconciliation
	//
	// # This is statement depends on comment at the UpdateFunc of the Predicates,
	// # if we only check for the label on the new object, this IF can be removed.
	// 5. If label "nim.opendatahub.io/nim-app" NOT "true":
	//		- Patch OdhNimApp.Status.Condition[Type=ApiKeyValidated] to False
	//		- Do we want to tear down anything at this point?
	//		- Break reconciliation
	//
	// 6. If the OdhNimApp found:
	//		- Patch OdhNimApp.Spec.ApiKey.Validate to True
	//
	// 7. Else (OdhNimApp NOT found):
	//		- Create new OdhNimApp, set:
	//			- OdhNimApp.Spec.ApiKey.Validate to True
	//			- OdhNimApp.Spec.ApiKey.SecretRef to reference this Secret by Name and Namespace
	//
	// 8.  Reconcile a daily recurring Cron Job owned by the OdhNimApp, patching OdhNimApp.Spec.ApiKey.Validate to True

	return ctrl.Result{}, nil
}

// init is used for registering the secret controller for loading
func init() {
	controllerSetups = append(controllerSetups, func(opts ControllerOptions) error {
		return (&SecretController{
			opts.Manager.GetClient(),
			opts.Manager.GetScheme(),
		}).SetupWithManager(opts.Manager)
	})
}
