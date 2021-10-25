
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/stenic/sql-operator:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)



clean: ## Clean build directories across the project
	rm -rf ./bin
	rm -rf ./testbin
	rm -rf ./release-artifacts
	rm -rf ./charts/*/charts ./charts/*/Chart.lock
	rm -rf ./cover.out
	rm -rf ./generated-check

mod-tidy: ## Make sure the go mod files are up-to-date
	export GO111MODULE=on; go mod tidy

##@ Development

manifests: controller-gen helm-docs ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=sql-operator-role webhook paths="./api/..." paths="./controllers/." output:crd:artifacts:config=config/crd/bases
	CONFIG_DIRECTORY=$(or $(TMP_CONFIG_OUTPUT_DIRECTORY),config) HELM_DIRECTORY=$(or $(TMP_HELM_OUTPUT_DIRECTORY),charts) ./hack/config/copy_crds_roles_helm.sh
	$(HELM_DOCS)
	

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="./hack/headers/header.go.txt" paths=$(or $(TMP_API_DIRECTORY),"./...")

##@ Tests and Checks

check: lint test ## Do all checks, lints and tests for the Solr Operator

lint: check-mod vet check-manifests check-generated check-helm ## Lint the codebase to check for formatting and correctness

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

check-manifests: ## Ensure the manifest files (CRDs, RBAC, etc) are up-to-date across the project, including the helm charts
	rm -rf generated-check
	mkdir -p generated-check
	cp -r charts generated-check/charts
	cp -r config generated-check/config
	TMP_CONFIG_OUTPUT_DIRECTORY=generated-check/config TMP_HELM_OUTPUT_DIRECTORY=generated-check/charts make manifests
	@echo "Check to make sure the manifests are up to date"
	diff --recursive config generated-check/config
	diff --recursive charts generated-check/charts

check-generated: ## Ensure the generated code is up-to-date
	rm -rf generated-check
	mkdir -p generated-check
	cp -r api generated-check/api
	cp -r config generated-check/config
	TMP_API_DIRECTORY="./generated-check/api/..." make generate
	@echo "Check to make sure the generated code is up to date"
	diff --recursive api generated-check/api

check-mod: ## Ensure the go mod files are up-to-date
	rm -rf generated-check
	mkdir -p generated-check/existing-go-mod generated-check/go-mod
	cp go.* generated-check/existing-go-mod/.
	make mod-tidy
	cp go.* generated-check/go-mod/.
	mv generated-check/existing-go-mod/go.* .
	@echo "Check to make sure the go mod info is up to date"
	diff go.mod generated-check/go-mod/go.mod
	diff go.sum generated-check/go-mod/go.sum

check-helm: ## Ensure the helm charts lint successfully
	helm lint charts/*

check-git: ## Check to make sure the repo does not have uncommitted code
	git diff --exit-code

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

##@ Build

build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

run: manifests generate ## Run a controller from your host.
	go run ./main.go 2>&1

docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -


CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.5.0)

HELM_DOCS = $(shell pwd)/bin/helm-docs
helm-docs: ## Download helm-docs locally if necessary.
	$(call go-get-tool,$(HELM_DOCS),github.com/norwoodj/helm-docs/cmd/helm-docs@v1.6.0)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
