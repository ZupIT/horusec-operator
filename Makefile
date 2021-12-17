GO ?= go
GOFMT ?= gofmt
GO_FILES ?= $$(find . -name '*.go' | grep -v vendor | grep -v _gen.go | grep -v zz_generated.deepcopy.go | grep -v wire.go)
GOLANG_CI_LINT ?= golangci-lint
GO_IMPORTS ?= goimports
GO_IMPORTS_LOCAL ?= github.com/ZupIT/horusec-operator
HORUSEC ?= horusec
CONTROLLER_GEN ?= $(shell pwd)/bin/controller-gen
KUSTOMIZE ?= $(shell pwd)/bin/kustomize
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"
OPERATOR_VERSION ?= $(shell curl -sL https://api.github.com/repos/ZupIT/horusec-operator/releases/latest | jq -r ".tag_name") # Get the latest version of the operator
REGISTRY_IMAGE ?= horuszup/horusec-operator:${OPERATOR_VERSION}
ADDLICENSE ?= addlicense
GO_GCI ?= gci
GO_FUMPT ?= gofumpt

lint: # Run install and run golangci lint tool
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...

format: install-format-dependencies # Install format dependencies and run format
	$(GOFMT) -s -l -w $(GO_FILES)
	$(GO_IMPORTS) -w -local $(GO_IMPORTS_LOCAL) $(GO_FILES)
	$(GO_FUMPT) -l -w $(GO_FILES)
	$(GO_GCI) -w -local $(GO_IMPORTS_LOCAL) $(GO_FILES)

install-format-dependencies: # Install format dependencies
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/daixiang0/gci@latest

coverage: # Check coverage in application
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/coverage.sh | bash -s 0 .

tests: # Run all tests in application
	$(GO) clean -testcache && $(GO) test -v ./... -timeout=2m -parallel=1 -failfast -short

security: # Run security pipeline
    ifeq (, $(shell which $(HORUSEC)))
		curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s latest
		$(HORUSEC) start -p="./" -e="true"
    else
		$(HORUSEC) start -p="./" -e="true"
    endif

build: # Build operator image
	$(GO) build -o "./tmp/bin/operator" ./cmd/app

license: # Check for missing license headers
	$(GO) install github.com/google/addlicense@latest
	@$(ADDLICENSE) -check -f ./copyright.txt $(shell find -regex '.*\.\(go\|js\|ts\|yml\|yaml\|sh\|dockerfile\)')

license-fix: # Add missing license headers
	$(GO) install github.com/google/addlicense@latest
	@$(ADDLICENSE) -f ./copyright.txt $(shell find -regex '.*\.\(go\|js\|ts\|yml\|yaml\|sh\|dockerfile\)')

pipeline: format lint test coverage build security license  # Run all processes of the pipeline

######### Operator commands #########
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

kustomize: # Install kustomize binary
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

controller-gen: # Install controller-gen binary
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

manifests: controller-gen  # Update all manifests in config
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config="./config/crd/bases"

generate: controller-gen # Generate new controller
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

generate-service-yaml: kustomize
	mkdir -p $(shell pwd)/tmp
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(REGISTRY_IMAGE)
	$(KUSTOMIZE) build config/default > $(shell pwd)/tmp/horusec-operator.yaml

install: manifests kustomize # install horusec crd in kubernetes
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize # uninstall horusec crd in kubernetes
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize # deploy horusec-operator in environment
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(REGISTRY_IMAGE)
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: # undeploy horusec-operator in environment
	$(KUSTOMIZE) build config/default | kubectl delete -f -

mock: # generate source code for a mock
	mockgen -package=test -destination test/kubernetes_client.go -source=internal/horusec/usecase/kubeclient.go KubernetesClient

up-sample: # Starts a kind cluster with all platform dependencies and secrets
	chmod +x ./config/samples/sample_install_dependencies.sh
	./config/samples/sample_install_dependencies.sh

apply-sample: # Apply a sample operator from platform
	kubectl apply -f ./config/samples/install_v2alpha1_horusecplatform.yaml

replace-sample: # Replace sample operator from platform
	kubectl replace -f ./config/samples/install_v2alpha1_horusecplatform.yaml

up-local-operator: up-sample install deploy apply-sample # Start a kind cluster and install operator with a example
