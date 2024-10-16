// Copyright (c) 2024 Red Hat, Inc.

package webhooks

import ctrl "sigs.k8s.io/controller-runtime"

// WebhookOptions is encapsulating the global options for use with all webhooks
type WebhookOptions struct {
	Manager ctrl.Manager
}

// webhooksSetups is used for registering webhooks for loading
var webhooksSetups []func(WebhookOptions) error

// SetupWebhooks is used for setting up all registered webhooks
func SetupWebhooks(opts WebhookOptions) error {
	for _, webhookSetup := range webhooksSetups {
		if err := webhookSetup(opts); err != nil {
			return err
		}
	}
	return nil
}
