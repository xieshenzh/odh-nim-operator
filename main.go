// Copyright (c) 2024 Red Hat, Inc.

package main

import (
	"github.com/opendatahub-io/odh-nim-operator/pkg/operator"
	"github.com/opendatahub-io/odh-nim-operator/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/component-base/cli"
)

// command for running the operator
var cmd = &cobra.Command{
	Use:   "odhnimoperator",
	Short: "Open Data Hub NIM Operator",
}

// init is used for binding the operator to the command
func init() {
	oper := operator.NewOdhNimOperator()

	cmd.Flags().StringVar(
		&oper.Options.MetricAddr,
		"metric-address",
		":8080",
		"The address the metric endpoint binds to.")
	cmd.Flags().StringVar(
		&oper.Options.ProbeAddr,
		"probe-address",
		":8081",
		"The address the probe endpoint binds to.")
	cmd.Flags().BoolVar(
		&oper.Options.LeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controllers manager.")
	cmd.Flags().BoolVar(
		&oper.Options.Debug,
		"debug",
		false,
		"Enable debug logging")
	cmd.Flags().BoolVar(
		&oper.Options.EnableWebhooks,
		"enable-webhooks",
		false,
		"Enable admission webhooks")

	cmd.RunE = oper.Run
	cmd.Version = version.Get().GitVersion
}

// main is used for running the odh nim operator command
func main() {
	if err := cli.RunNoErrOutput(cmd); err != nil {
		panic(err)
	}
}
