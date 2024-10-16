// Copyright (c) 2024 Red Hat, Inc.

package operator

import (
	"github.com/opendatahub-io/odh-nim-operator/api/v1alpha1"
	"github.com/opendatahub-io/odh-nim-operator/pkg/controllers"
	"github.com/opendatahub-io/odh-nim-operator/pkg/utils"
	"github.com/opendatahub-io/odh-nim-operator/pkg/webhooks"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// OdhNimOperator is the receiver for running the operator and binding the options
type OdhNimOperator struct {
	Options OdhNimOperatorOptions
}

// OdhNimOperatorOptions is used for encapsulating the operator options
type OdhNimOperatorOptions struct {
	MetricAddr     string
	LeaderElection bool
	ProbeAddr      string
	Debug          bool
	EnableWebhooks bool
	controllers.ControllerOptions
}

// NewOdhNimOperator is a factory function for creating a Nim operator instance
func NewOdhNimOperator() OdhNimOperator {
	return OdhNimOperator{Options: OdhNimOperatorOptions{}}
}

// Run is the function to run for the command, used for running the operator. It creates and configures all aspects of
// the operator, i.e. manager, scheme, health checks, webhooks, and controllers.
func (o *OdhNimOperator) Run(cmd *cobra.Command, args []string) error {
	// set logging and create initial logger
	ctrl.SetLogger(zap.New(zap.UseDevMode(o.Options.Debug)))
	logger := log.Log.WithName("odh-nim-operator")

	// create the scheme and install the required types
	scheme := runtime.NewScheme()
	if err := utils.InstallTypes(scheme); err != nil {
		logger.Error(nil, "failed installing scheme")
		return err
	}

	// create the manager
	kubeConfig := config.GetConfigOrDie()
	mgr, err := ctrl.NewManager(kubeConfig, ctrl.Options{
		Scheme:                 scheme,
		Logger:                 logger,
		LeaderElection:         o.Options.LeaderElection,
		LeaderElectionID:       "odh-nim-leader-election-id",
		MetricsBindAddress:     o.Options.MetricAddr,
		HealthProbeBindAddress: o.Options.ProbeAddr,
		NewCache: cache.BuilderWithOptions(cache.Options{
			SelectorsByObject: cache.SelectorsByObject{
				&v1alpha1.OdhNimApp{}: {Label: labels.Everything()},
			},
		}),
	})
	if err != nil {
		logger.Error(err, "failed creating k8s manager")
		return err
	}

	// setup controllers
	o.Options.ControllerOptions.Manager = mgr
	if err = controllers.SetupControllers(o.Options.ControllerOptions); err != nil {
		logger.Error(err, "failed setting up the controllers")
		return err
	}

	// setup webhooks
	if o.Options.EnableWebhooks {
		wopts := webhooks.WebhookOptions{Manager: mgr}
		if err = webhooks.SetupWebhooks(wopts); err != nil {
			logger.Error(err, "failed setting up the webhooks")
			return err
		}
	}

	// setup health checks
	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logger.Error(err, "failed setting up health check")
		return err
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		logger.Error(err, "failed setting up ready check")
		return err
	}

	// start the manager to run the operator
	return mgr.Start(cmd.Context())
}
