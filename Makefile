# Copyright (c) 2024 Red Hat, Inc.

######################################
###### NIM OpenDataHub Operator ######
######################################
default: help

OPERATOR_NAMESPACE ?= opendatahub-operator-system##@ Set the target namespace for deploying the operator, defaults to 'opendatahub-operator-system'

############################################################################
###### Create working directories (note .gitignore) and fetch OS info ######
############################################################################
LOCALBIN = $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

LOCALBUILD = $(shell pwd)/build
$(LOCALBUILD):
	mkdir -p $(LOCALBUILD)

OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

#####################################
###### Image related variables ######
#####################################
IMAGE_REGISTRY ?= quay.io##@ Set the image registry, defaults to 'quay.io'
IMAGE_NAMESPACE ?= ecosystem-appeng##@ Set the image namespace, defaults to 'ecosystem-appeng'
IMAGE_NAME ?= odh-nim-operator##@ Set the operator image name, defaults to 'odh-nim-operator'
IMAGE_TAG ?= $(strip $(shell cat VERSION))##@ Set the operator image tag, defaults to content of the VERSION file
IMAGE_BUILDER = podman

######################################
###### Bundle related variables ######
######################################
BUNDLE_PACKAGE_NAME ?= $(IMAGE_NAME)##@ Set the bundle package name, defaults to IMAGE_NAME
BUNDLE_CHANNELS ?= alpha##@ Set a comma-seperated list of channels the bundle belongs too, defaults to 'alpha'
BUNDLE_DEFAULT_CHANNEL ?= alpha##@ Set the default channel for the bundle, defaults to 'alpha'
BUNDLE_IMAGE_NAME ?= $(IMAGE_NAME)-bundle##@ Set the image name for the bundle, defaults to IMAGE_NAME-bundle
BUNDLE_NAMESPACE ?= $(OPERATOR_NAMESPACE)##@ Set the target namespace for running the bundle, defaults to OPERATOR_NAMESPACE
BUNDLE_SCORECARD_NAMESPACE ?= $(IMAGE_NAME)-scorecard##@ Set the target namespace for running scorecard tests, defaults to IMAGE_NAME-scorecard
BUNDLE_TEST_VERBOSE ?= false##@ If true, will display full log for scorecard tests and exit with and error if any test fails

####################################################
###### Required tools customization variables ######
####################################################
REQ_BIN_AWK ?=""##@ Set a custom 'awk'/'gwak' binary path if not in PATH
REQ_BIN_OC ?= oc##@ Set a custom 'oc' binary path if not in PATH
REQ_BIN_GO ?= go##@ Set a custom 'go' binary path if not in PATH (useful for multi versions environment)
REQ_BIN_CURL ?= curl##@ Set a custom 'curl' binary path if not in PATH
REQ_BIN_YQ ?= yq##@ Set a custom 'yq' binary path if not in PATH
REQ_BIN_JQ ?= jq##@ Set a custom 'jq' binary path if not in PATH

# set default awk if one not provided
ifeq ($(REQ_BIN_AWK),"")
ifeq ($(OS),darwin)
REQ_BIN_AWK = gawk
else
REQ_BIN_AWK = awk
endif
endif

######################################################
###### Downloaded tools customization variables ######
######################################################
BIN_CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen##@ Set custom 'controller-gen', if not supplied will install in ./bin
BIN_OPERATOR_SDK ?= $(LOCALBIN)/operator-sdk##@ Set custom 'operator-sdk', if not supplied will install in ./bin
BIN_KUSTOMIZE ?= $(LOCALBIN)/kustomize##@ Set custom 'kustomize', if not supplied will install in ./bin
BIN_ENVTEST ?= $(LOCALBIN)/setup-envtest##@ Set custom 'setup-envtest', if not supplied will install in ./bin
BIN_GOLINTCI ?= $(LOCALBIN)/golangci-lint##@ Set custom 'golangci-lint', if not supplied will install in ./bin
BIN_KUTTL ?= $(LOCALBIN)/kubectl-kuttl##@ Set custom 'kuttl', if not supplied will install in ./bin

################################################
###### Downloaded tools version variables ######
################################################
VERSION_CONTROLLER_GEN = v0.16.4
VERSION_OPERATOR_SDK = v1.37.0
VERSION_KUSTOMIZE = v5.5.0
VERSION_GOLANG_CI_LINT = v1.61.0
VERSION_KUTTL = 0.19.0

#####################################
###### Build related variables ######
#####################################
ifeq ($(OS),darwin)
DATE_BIN = gdate
else
DATE_BIN = date
endif
BUILD_DATE = $(strip $(shell $(DATE_BIN) +%FT%T))
BUILD_TIMESTAMP = $(strip $(shell $(DATE_BIN) -d "$(BUILD_DATE)" +%s))
COMMIT_HASH = $(strip $(shell git rev-parse --short HEAD))
LDFLAGS=-ldflags="\
-X 'github.com/opendatahub-io/odh-nim-operator/pkg/version.tag=${IMAGE_TAG}' \
-X 'github.com/opendatahub-io/odh-nim-operator/pkg/version.commit=${COMMIT_HASH}' \
-X 'github.com/opendatahub-io/odh-nim-operator/pkg/version.date=${BUILD_DATE}' \
"

####################################
###### Test related variables ######
####################################
ENVTEST_K8S_VERSION = 1.29.x
OPERATOR_RUN_ARGS ?=##@ Use for setting custom run arguments for development local run

#########################
###### Image names ######
#########################
FULL_OPERATOR_IMAGE_NAME = $(strip $(IMAGE_REGISTRY)/$(IMAGE_NAMESPACE)/$(IMAGE_NAME):$(IMAGE_TAG))
FULL_OPERATOR_IMAGE_NAME_UNIQUE = $(FULL_OPERATOR_IMAGE_NAME)_$(COMMIT_HASH)_$(BUILD_TIMESTAMP)
FULL_BUNDLE_IMAGE_NAME = $(strip $(IMAGE_REGISTRY)/$(IMAGE_NAMESPACE)/$(BUNDLE_IMAGE_NAME):$(IMAGE_TAG))
FULL_BUNDLE_IMAGE_NAME_UNIQUE = $(FULL_BUNDLE_IMAGE_NAME)_$(COMMIT_HASH)_$(BUILD_TIMESTAMP)

####################################
###### Build and push project ######
####################################
build/all/image: build/operator/image build/bundle/image ## Build both the operator and bundle images

build/all/image/push: build/operator/image/push build/bundle/image/push ## Build and push both the operator and bundle images

.PHONY: build build/operator
build build/operator: $(LOCALBUILD) ## Build the project as a binary in ./build
	GOOS="linux" GOARCH="amd64" $(REQ_BIN_GO) build $(LDFLAGS) -o $(LOCALBUILD)/odhnimoperator ./main.go

build/operator/image: build/operator ## Build the operator image - Builds locally, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, IMAGE_NAME, and IMAGE_TAG
	$(IMAGE_BUILDER) build --platform linux/amd64 --tag $(FULL_OPERATOR_IMAGE_NAME) -f ./Containerfile

build/operator/image/push: build/operator/image ## Build and push the operator image, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, IMAGE_NAME, and IMAGE_TAG
	$(IMAGE_BUILDER) tag $(FULL_OPERATOR_IMAGE_NAME) $(FULL_OPERATOR_IMAGE_NAME_UNIQUE)
	$(IMAGE_BUILDER) push $(FULL_OPERATOR_IMAGE_NAME_UNIQUE)
	$(IMAGE_BUILDER) push $(FULL_OPERATOR_IMAGE_NAME)

.PHONY: build/bundle/image
build/bundle/image: ## Build the bundle image, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, BUNDLE_IMAGE_NAME, and IMAGE_TAG
	$(IMAGE_BUILDER) build --ignorefile ./.gitignore --tag $(FULL_BUNDLE_IMAGE_NAME) -f ./bundle.Containerfile

build/bundle/image/push: build/bundle/image ## Build and push the bundle image, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, BUNDLE_IMAGE_NAME, and IMAGE_TAG
	$(IMAGE_BUILDER) tag $(FULL_BUNDLE_IMAGE_NAME) $(FULL_BUNDLE_IMAGE_NAME_UNIQUE)
	$(IMAGE_BUILDER) push $(FULL_BUNDLE_IMAGE_NAME_UNIQUE)
	$(IMAGE_BUILDER) push $(FULL_BUNDLE_IMAGE_NAME)

###########################################
###### Code and Manifests generation ######
###########################################
generate generate/all: generate/manifests generate/webhooks generate/code generate/bundle generate/e2e ## Generate rbac, crd, webhooks, and e2e manifests, as well as code and olm bundle files

.PHONY: generate/manifests
generate/manifests: $(BIN_CONTROLLER_GEN) ## Generate rbac and crd manifest files
	$(BIN_CONTROLLER_GEN) rbac:roleName=role paths="./pkg/controllers/..."
	$(BIN_CONTROLLER_GEN) crd paths="./api/..."

.PHONY: generate/webhooks
generate/webhooks: $(BIN_CONTROLLER_GEN) ## Generate admission webhooks manifest files
	$(BIN_CONTROLLER_GEN) webhook paths="./pkg/webhooks/..."

.PHONY: generate/code
generate/code: $(BIN_CONTROLLER_GEN) ## Generate API boiler-plate code
	rm -rf ./api/**/zz_generated.deepcopy.go
	$(BIN_CONTROLLER_GEN) object:headerFile="hack/header.txt" paths="./api/..."

.PHONY: generate/bundle
generate/bundle: $(BIN_OPERATOR_SDK) $(BIN_KUSTOMIZE) ## Generate olm bundle
	@$(call kustomize-setup)
	$(BIN_OPERATOR_SDK) generate kustomize manifests
	$(BIN_KUSTOMIZE) build config/manifests | $(BIN_OPERATOR_SDK) generate bundle --quiet --version $(IMAGE_TAG) \
	--package $(BUNDLE_PACKAGE_NAME) --channels $(BUNDLE_CHANNELS) --default-channel $(BUNDLE_DEFAULT_CHANNEL)
	@mv -f ./bundle.Dockerfile ./bundle.Containerfile

.PHONY: generate/e2e
generate/e2e: $(BIN_KUSTOMIZE) ## Generate deployment files for E2E testing
	cp config/e2e/kustomization.yaml config/e2e/kustomization.yaml.tmp
	(cd config/e2e && $(BIN_KUSTOMIZE) edit set image odh-nim-operator-image==$(FULL_OPERATOR_IMAGE_NAME))
	$(BIN_KUSTOMIZE) build config/e2e > e2e/kuttl/manifests/manifests.yaml
	-mv config/e2e/kustomization.yaml.tmp config/e2e/kustomization.yaml

################################################
###### Install and Uninstall the operator ######
################################################
.PHONY: operator/run
operator/run: ## Run the Operator in your local environment for development purposes, use OPERATOR_RUN_ARGS for run args
	go run main.go --debug $(OPERATOR_RUN_ARGS)

.PHONY: operator/deploy
operator/deploy: $(BIN_KUSTOMIZE) ## Deploy the Operator
	@$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	@$(call kustomize-setup)
	$(BIN_KUSTOMIZE) build config/default | $(REQ_BIN_OC) apply -f -

.PHONY: operator/deploy/stdout
operator/deploy/stdout: $(BIN_KUSTOMIZE) ## Build the Operator manifests to STDOUT
	@$(call kustomize-setup)
	$(BIN_KUSTOMIZE) build config/default

.PHONY: operator/undeploy
operator/undeploy: $(BIN_KUSTOMIZE) ## Undeploy the Operator
	@$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	@$(call kustomize-setup)
	$(BIN_KUSTOMIZE) build config/default | $(REQ_BIN_OC) delete --ignore-not-found --warnings-as-errors -f -

.PHONY: bundle/run
bundle/run: $(BIN_OPERATOR_SDK) ## Run the Operator OLM Bundle from image
	@$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	-$(REQ_BIN_OC) create ns $(BUNDLE_NAMESPACE)
	$(BIN_OPERATOR_SDK) run bundle $(FULL_BUNDLE_IMAGE_NAME) -n $(BUNDLE_NAMESPACE)

.PHONY: bundle/cleanup
bundle/cleanup: $(BIN_OPERATOR_SDK) ## Cleanup the Operator OLM Bundle package installed
	$(BIN_OPERATOR_SDK) cleanup $(BUNDLE_PACKAGE_NAME) -n $(BUNDLE_NAMESPACE)

.PHONY: bundle/cleanup/namespace
bundle/cleanup/namespace: ## DELETE the Operator OLM Bundle namespace (BE CAREFUL)
	@$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	$(REQ_BIN_OC) delete ns $(BUNDLE_NAMESPACE)

.PHONY: operator/api/install
operator/api/install: ## Install all owned CustomResourceDefinitions (not required for the deploy target)
	$(BIN_KUSTOMIZE) build config/crd | $(REQ_BIN_OC) apply -f -

.PHONY: operator/api/uninstall
operator/api/uninstall: ## Uninstall all owned CustomResourceDefinitions (not required for the undeploy target)
	$(BIN_KUSTOMIZE) build config/crd | $(REQ_BIN_OC) delete -f -

###########################
###### Test codebase ######
###########################
kubeAssets = "KUBEBUILDER_ASSETS=$(shell $(BIN_ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)"

testCmd = "$(kubeAssets) $(REQ_BIN_GO) test -v ./pkg/... -ginkgo.v"
ifdef TEST_NAME
testCmd += " -ginkgo.focus \"$(TEST_NAME)\""
endif

.PHONY: test
test: $(BIN_ENVTEST) ## Run all unit tests, Use TEST_NAME to run a specific test
	@eval $(testCmd)

covTestCmd = "$(kubeAssets) $(REQ_BIN_GO) test -failfast -coverprofile=cov.out -v ./pkg/controllers/... -ginkgo.v"

.PHONY: test/cov
test/cov: $(BIN_ENVTEST) ## Run all unit tests and print coverage report
	@eval $(covTestCmd)
	$(REQ_BIN_GO) tool cover -func=cov.out
	$(REQ_BIN_GO) tool cover -html=cov.out -o cov.html

testBundleCmd = "$(BIN_OPERATOR_SDK) scorecard ./bundle -n $(BUNDLE_SCORECARD_NAMESPACE) --pod-security=restricted"
ifneq ($(BUNDLE_TEST_VERBOSE),true)
testBundleCmd += " --output json | $(REQ_BIN_JQ) '[ .items[].status.results | del(.[].creationTimestamp, .[].log) | .[] ]'"
endif

.PHONY: test/bundle
test/bundle: $(BIN_OPERATOR_SDK) ## Run Scorecard Bundle Tests (requires connected cluster)
	$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	$(call verify-essential-tool,$(REQ_BIN_JQ),REQ_BIN_JQ)
	@ { \
	if $(REQ_BIN_OC) create ns $(BUNDLE_SCORECARD_NAMESPACE); then \
		$(call run-scorecard-tests); \
		$(REQ_BIN_OC) delete ns $(BUNDLE_SCORECARD_NAMESPACE); \
	else \
		$(call run-scorecard-tests); \
	fi \
	}

# if BUNDLE_TEST_VERBOSE=true, the results are processed with jq and this will not error out for failing tests
define run-scorecard-tests
if !(eval $(testBundleCmd)); then \
	echo "bundle test failed"; \
	exit 1; \
fi
endef

.PHONY: test/bundle/delete/ns
test/bundle/delete/ns: ## DELETE the Scorecard namespace (BE CAREFUL)
	@$(call verify-essential-tool,$(REQ_BIN_OC),REQ_BIN_OC)
	-$(REQ_BIN_OC) delete ns $(BUNDLE_SCORECARD_NAMESPACE)

.PHONY: test/e2e/kuttl
test/e2e/kuttl: $(BIN_KUTTL) generate/e2e build/operator/image ## Run End-to-End tests, will build the operator and image (requires kind)
	$(call verify-essential-tool,$(REQ_BIN_YQ),REQ_BIN_YQ)
	cp e2e/kuttl/kuttl-test.yaml e2e/kuttl/kuttl-test-patched.yaml
	$(REQ_BIN_YQ) -i '.kindContainers = ["$(FULL_OPERATOR_IMAGE_NAME)"]' e2e/kuttl/kuttl-test-patched.yaml
	$(BIN_KUTTL) test --config e2e/kuttl/kuttl-test-patched.yaml
	@ rm e2e/kuttl/kuttl-test-patched.yaml

###########################
###### Lint codebase ######
###########################
lint/all: lint/code lint/containerfile lint/bundle ## Lint the entire project (code, containerfile, bundle)

.PHONY: lint lint/code
lint lint/code: $(BIN_GOLINTCI) ## Lint the code
	$(REQ_BIN_GO) fmt ./...
	$(BIN_GOLINTCI) run

.PHONY: lint/containerfile
lint/containerfile: ## Lint the Containerfile (using Hadolint image, do not use inside a container)
	$(IMAGE_BUILDER) run --rm -i docker.io/hadolint/hadolint:latest < ./Containerfile

.PHONY: lint/bundle
lint/bundle: $(BIN_OPERATOR_SDK) ## Validate OLM bundle
	$(BIN_OPERATOR_SDK) bundle validate ./bundle --select-optional suite=operatorframework

####################################
###### Install required tools ######
####################################
$(BIN_KUSTOMIZE): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(REQ_BIN_GO) install sigs.k8s.io/kustomize/kustomize/v5@$(VERSION_KUSTOMIZE)

$(BIN_CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(REQ_BIN_GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@$(VERSION_CONTROLLER_GEN)

$(BIN_ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(REQ_BIN_GO) install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

$(BIN_GOLINTCI): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(REQ_BIN_GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(VERSION_GOLANG_CI_LINT)

$(BIN_OPERATOR_SDK): $(LOCALBIN)
	@$(call verify-essential-tool,$(REQ_BIN_CURL),REQ_BIN_CURL)
	$(REQ_BIN_CURL) -sSLo $(BIN_OPERATOR_SDK) https://github.com/operator-framework/operator-sdk/releases/download/$(VERSION_OPERATOR_SDK)/operator-sdk_$(OS)_$(ARCH)
	chmod +x $(BIN_OPERATOR_SDK)

$(BIN_KUTTL): $(LOCALBIN)
	@$(call verify-essential-tool,$(REQ_BIN_CURL),REQ_BIN_CURL)
	$(REQ_BIN_CURL) -sSLo $(BIN_KUTTL) https://github.com/kudobuilder/kuttl/releases/download/v$(VERSION_KUTTL)/kubectl-kuttl_$(VERSION_KUTTL)_$(OS)_$(ARCH)
	chmod +x $(BIN_KUTTL)

###############################
###### Utility functions ######
###############################
# the namespace is being set with yq because components are loaded before transformers in kustomize (see prometheus component)
define kustomize-setup
$(call verify-essential-tool,$(REQ_BIN_YQ),REQ_BIN_YQ)
cd config/default && \
$(BIN_KUSTOMIZE) edit set image odh-nim-operator-image=$(FULL_OPERATOR_IMAGE_NAME) && \
$(BIN_KUSTOMIZE) edit set namespace $(OPERATOR_NAMESPACE)
$(REQ_BIN_YQ) -i '.labels[1].pairs."app.kubernetes.io/instance" = "odh-nim-operator-$(IMAGE_TAG)"' config/default/kustomization.yaml
$(REQ_BIN_YQ) -i '.labels[1].pairs."app.kubernetes.io/version" = "$(IMAGE_TAG)"' config/default/kustomization.yaml
$(REQ_BIN_YQ) -i '.metadata.name = "$(OPERATOR_NAMESPACE)"' config/manager/namespace.yaml
endef

# arg1 = name of the tool to look for | arg2 = name of the variable for a custom replacement
TOOL_MISSING_ERR_MSG = Please install '$(1)' or specify a custom path using the '$(2)' variable
define verify-essential-tool
@if !(which $(1) &> /dev/null); then \
	echo $(call TOOL_MISSING_ERR_MSG,$(1),$(2)); \
	exit 1; \
fi
endef

################################
###### Display build help ######
################################
.PHONY: help
help: ## Show this help message
	$(call verify-essential-tool,$(REQ_BIN_AWK),REQ_BIN_AWK)
	@$(REQ_BIN_AWK) 'BEGIN {\
			FS = ".*##@";\
			print "\033[1;31mNIM Open Data Hub Operator\033[0m";\
			print "\033[1;32mUsage\033[0m";\
			printf "\t\033[1;37mmake <target> |";\
			printf "\tmake <target> [Variables Set] |";\
            printf "\tmake [Variables Set] <target> |";\
            print "\t[Variables Set] make <target>\033[0m";\
			print "\033[1;32mAvailable Variables\033[0m" }\
		/^(\s|[a-zA-Z_0-9-]|\/)+ \?=.*?##@/ {\
			split($$0,t,"?=");\
			printf "\t\033[1;36m%-35s \033[0;37m%s\033[0m\n",t[1], $$2 | "sort" }'\
		$(MAKEFILE_LIST)
	@$(REQ_BIN_AWK) 'BEGIN {\
			FS = ":.*##";\
			SORTED = "sort";\
            print "\033[1;32mAvailable Targets\033[0m"}\
		/^(\s|[a-zA-Z_0-9-]|\/)+:.*?##/ {\
			if($$0 ~ /deploy/)\
				printf "\t\033[1;36m%-35s \033[0;33m%s\033[0m\n", $$1, $$2 | SORTED;\
			else if($$0 ~ /push/)\
				printf "\t\033[1;36m%-35s \033[0;35m%s\033[0m\n", $$1, $$2 | SORTED;\
			else if($$0 ~ /DELETE/)\
				printf "\t\033[1;36m%-35s \033[0;31m%s\033[0m\n", $$1, $$2 | SORTED;\
			else\
				printf "\t\033[1;36m%-35s \033[0;37m%s\033[0m\n", $$1, $$2 | SORTED; }\
		END { \
			close(SORTED);\
			print "\033[1;32mFurther Information\033[0m";\
			print "\t\033[0;37m* Source code: \033[38;5;26mhttps://github.com/RHEcosystemappeng/odh-nim-operator\33[0m"}'\
		$(MAKEFILE_LIST)
