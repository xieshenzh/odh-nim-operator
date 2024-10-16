# Contributing to ODH NIM Operator

```shell
$ make help

NIM Open Data Hub Operator
Usage
        make <target> | make <target> [Variables Set] | make [Variables Set] <target> | [Variables Set] make <target>
Available Variables
        BIN_CONTROLLER_GEN                   Set custom 'controller-gen', if not supplied will install in ./bin
        BIN_ENVTEST                          Set custom 'setup-envtest', if not supplied will install in ./bin
        BIN_GOLINTCI                         Set custom 'golangci-lint', if not supplied will install in ./bin
        BIN_KUSTOMIZE                        Set custom 'kustomize', if not supplied will install in ./bin
        BIN_KUTTL                            Set custom 'kuttl', if not supplied will install in ./bin
        BIN_OPERATOR_SDK                     Set custom 'operator-sdk', if not supplied will install in ./bin
        BUNDLE_CHANNELS                      Set a comma-seperated list of channels the bundle belongs too, defaults to 'alpha'
        BUNDLE_DEFAULT_CHANNEL               Set the default channel for the bundle, defaults to 'alpha'
        BUNDLE_IMAGE_NAME                    Set the image name for the bundle, defaults to IMAGE_NAME-bundle
        BUNDLE_NAMESPACE                     Set the target namespace for running the bundle, defaults to OPERATOR_NAMESPACE
        BUNDLE_PACKAGE_NAME                  Set the bundle package name, defaults to IMAGE_NAME
        BUNDLE_SCORECARD_NAMESPACE           Set the target namespace for running scorecard tests, defaults to IMAGE_NAME-scorecard
        BUNDLE_TEST_VERBOSE                  If true, will display full log for scorecard tests and exit with and error if any test fails
        IMAGE_NAME                           Set the operator image name, defaults to 'odh-nim-operator'
        IMAGE_NAMESPACE                      Set the image namespace, defaults to 'ecosystem-appeng'
        IMAGE_REGISTRY                       Set the image registry, defaults to 'quay.io'
        IMAGE_TAG                            Set the operator image tag, defaults to content of the VERSION file
        OPERATOR_NAMESPACE                   Set the target namespace for deploying the operator, defaults to 'opendatahub-operator-system'
        OPERATOR_RUN_ARGS                    Use for setting custom run arguments for development local run
        REQ_BIN_AWK                          Set a custom 'awk'/'gwak' binary path if not in PATH
        REQ_BIN_CURL                         Set a custom 'curl' binary path if not in PATH
        REQ_BIN_GO                           Set a custom 'go' binary path if not in PATH (useful for multi versions environment)
        REQ_BIN_JQ                           Set a custom 'jq' binary path if not in PATH
        REQ_BIN_OC                           Set a custom 'oc' binary path if not in PATH
        REQ_BIN_YQ                           Set a custom 'yq' binary path if not in PATH
Available Targets
        build build/operator                 Build the project as a binary in ./build - Requires VPN Access
        build/all/image                      Build both the operator and bundle images
        build/all/image/push                 Build and push both the operator and bundle images
        build/bundle/image                   Build the bundle image, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, BUNDLE_IMAGE_NAME, and IMAGE_TAG
        build/bundle/image/push              Build and push the bundle image, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, BUNDLE_IMAGE_NAME, and IMAGE_TAG
        build/operator/image                 Build the operator image - Builds locally, requires VPN Access, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, IMAGE_NAME, and IMAGE_TAG
        build/operator/image/push            Build and push the operator image - Requires VPN Access, customized with IMAGE_REGISTRY, IMAGE_NAMESPACE, IMAGE_NAME, and IMAGE_TAG
        bundle/cleanup                       Cleanup the Operator OLM Bundle package installed
        bundle/cleanup/namespace             DELETE the Operator OLM Bundle namespace (BE CAREFUL)
        bundle/run                           Run the Operator OLM Bundle from image
        generate generate/all                Generate rbac, crd, webhooks, and e2e manifests, as well as code and olm bundle files
        generate/bundle                      Generate olm bundle
        generate/code                        Generate API boiler-plate code
        generate/e2e                         Generate deployment files for E2E testing
        generate/manifests                   Generate rbac and crd manifest files
        generate/webhooks                    Generate admission webhooks manifest files
        help                                 Show this help message
        lint lint/code                       Lint the code
        lint/all                             Lint the entire project (code, containerfile, bundle)
        lint/bundle                          Validate OLM bundle
        lint/containerfile                   Lint the Containerfile (using Hadolint image, do not use inside a container)
        operator/api/install                 Install all owned CustomResourceDefinitions (not required for the deploy target)
        operator/api/uninstall               Uninstall all owned CustomResourceDefinitions (not required for the undeploy target)
        operator/deploy                      Deploy the Operator
        operator/deploy/stdout               Build the Operator manifests to STDOUT
        operator/run                         Run the Operator in your local environment for development purposes, use OPERATOR_RUN_ARGS for run args
        operator/undeploy                    Undeploy the Operator
        test                                 Run all unit tests, Use TEST_NAME to run a specific test
        test/bundle                          Run Scorecard Bundle Tests (requires connected cluster)
        test/bundle/delete/ns                DELETE the Scorecard namespace (BE CAREFUL)
        test/cov                             Run all unit tests and print coverage report
        test/e2e/kuttl                       Run End-to-End tests, will build the operator and image (requires kind)
Further Information
        * Source code: https://github.com/RHEcosystemAppEng/odh-nim-operator

```
