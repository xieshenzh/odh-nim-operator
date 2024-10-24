// Copyright (c) 2024 Red Hat, Inc.

package controllers

import (
	"context"

	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"
	dscv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/datasciencecluster/v1"
	tmpv1 "github.com/openshift/api/template/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

const (
	kServeCondition    = "KServeReadiness"
	apiKeyCondition    = "APIKeyValidation"
	templateCondition  = "TemplateUpdate"
	configMapCondition = "ConfigMapUpdate"

	kServeNotReady = "KServeNotReady"
	kServeReady    = "KServeReady"
)

// SetupWithManager is used for setting up the controller with a manager (check the init function)
func (r *OdhNIMAccountController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("odh-nim-account-controller").
		For(&v1alpha1.OdhNIMAccount{}).
		Complete(r)
}

func (r *OdhNIMAccountController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).
		WithName("odh-nim-account-controller").
		WithValues("OdhNIMAccount", req.NamespacedName)

	ctx = log.IntoContext(ctx, logger)

	logger.V(1).Info("got request for OdhNIMAccount")

	var account v1alpha1.OdhNIMAccount
	if err := r.Get(ctx, req.NamespacedName, &account); err != nil {
		if errors.IsNotFound(err) {
			logger.V(1).Info("OdhNIMAccount not found")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "error fetching OdhNIMAccount")
		return ctrl.Result{}, err
	}

	if ready, err := r.checkKServeReady(ctx); err != nil {
		conditions := []metav1.Condition{
			{
				Type:   kServeCondition,
				Status: metav1.ConditionFalse,
				Reason: kServeNotReady,
			},
		}
		_, _ = r.updateStatus(ctx, &account, conditions, account.Status.TemplateRef, account.Status.ConfigMapRef)

		logger.Error(err, "error fetching DataScienceCluster")
		return ctrl.Result{}, err
	} else if !ready {
		if r, e := r.handleKServeNotReady(ctx, &account); e != nil {
			logger.Error(e, "error updating OdhNIMAccount status")
			return ctrl.Result{}, err
		} else if r {
			logger.V(1).Info("OdhNIMAccount modified")
			return ctrl.Result{Requeue: true}, nil
		}
		logger.V(1).Info("KServe not ready")
		return ctrl.Result{Requeue: true}, nil
	}

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

func (r *OdhNIMAccountController) updateStatus(ctx context.Context, account *v1alpha1.OdhNIMAccount,
	conditions []metav1.Condition, template *corev1.ObjectReference, configMap *corev1.ObjectReference) (bool, error) {
	for _, condition := range conditions {
		apimeta.SetStatusCondition(&account.Status.Conditions, condition)
	}

	account.Status.TemplateRef = template
	account.Status.ConfigMapRef = configMap

	if err := r.Status().Update(ctx, account); err != nil {
		if errors.IsConflict(err) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (r *OdhNIMAccountController) checkKServeReady(ctx context.Context) (bool, error) {
	var dscList dscv1.DataScienceClusterList
	if err := r.List(ctx, &dscList); err != nil {
		return false, err
	}

	if len(dscList.Items) > 0 {
		dsc := dscList.Items[0]
		if dsc.Status.Phase == "Ready" {
			ready, ok := dsc.Status.InstalledComponents["kserve"]
			if ok && ready {
				return true, nil
			}
		}
	}

	return false, nil
}

func (r *OdhNIMAccountController) handleKServeNotReady(ctx context.Context, account *v1alpha1.OdhNIMAccount) (bool, error) {
	conditions := []metav1.Condition{
		{
			Type:   kServeCondition,
			Status: metav1.ConditionFalse,
			Reason: kServeNotReady,
		},
		{
			Type:   apiKeyCondition,
			Status: metav1.ConditionUnknown,
			Reason: kServeNotReady,
		},
		{
			Type:   templateCondition,
			Status: metav1.ConditionFalse,
			Reason: kServeNotReady,
		},
		{
			Type:   configMapCondition,
			Status: metav1.ConditionFalse,
			Reason: kServeNotReady,
		},
	}

	rq, err := r.updateStatus(ctx, account, conditions, nil, nil)

	if err == nil {
		err = r.deleteTemplate(ctx, account)
	}

	if err == nil {
		err = r.deleteConfigMap(ctx, account)
	}

	return rq, err
}

func (r *OdhNIMAccountController) deleteTemplate(ctx context.Context, account *v1alpha1.OdhNIMAccount) error {
	if account.Status.TemplateRef != nil {
		var template tmpv1.Template
		if err := r.Get(ctx,
			client.ObjectKey{
				Namespace: account.Status.TemplateRef.Namespace,
				Name:      account.Status.TemplateRef.Name},
			&template); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		} else {
			if e := r.Delete(ctx, &template); e != nil {
				return e
			}
		}
	}
	return nil
}

func (r *OdhNIMAccountController) deleteConfigMap(ctx context.Context, account *v1alpha1.OdhNIMAccount) error {
	if account.Status.ConfigMapRef != nil {
		var cm corev1.ConfigMap
		if err := r.Get(ctx,
			client.ObjectKey{
				Namespace: account.Status.ConfigMapRef.Namespace,
				Name:      account.Status.ConfigMapRef.Name},
			&cm); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		} else {
			if e := r.Delete(ctx, &cm); e != nil {
				return e
			}
		}
	}
	return nil
}
