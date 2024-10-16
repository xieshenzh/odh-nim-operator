// Copyright (c) 2024 Red Hat, Inc.

package controllers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/opendatahub-io/odh-nim-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

// use the test client for testing the controllers
var testClient client.Client
var testEnv *envtest.Environment

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Tests")
}

var _ = BeforeSuite(func(ctx SpecContext) {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	By("bootstrapping testing environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("testdata", "required_crds")},
	}

	// install the scheme
	scheme := runtime.NewScheme()
	Expect(utils.InstallTypes(scheme)).To(Succeed())

	// start testing environment and get config for the client
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// create and save the test client
	testClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(testClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down testing environment")
	Expect(testEnv.Stop()).To(Succeed())
})

// #####################################
// ##### Testing utility functions #####
// #####################################
func cleanup(ctx SpecContext, objs ...client.Object) error {
	for _, obj := range objs {
		if err := testClient.Delete(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}
