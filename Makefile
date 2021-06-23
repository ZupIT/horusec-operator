GO ?= go
GOFMT ?= gofmt
GO_FILES ?= $$(find . -name '*.go' | grep -v vendor)
GOLANG_CI_LINT ?= $(shell pwd)/bin/golangci-lint
GO_IMPORTS ?= goimports
GO_IMPORTS_LOCAL ?= github.com/ZupIT/horusec-operator
HORUSEC ?= horusec
CONTROLLER_GEN ?= $(shell pwd)/bin/controller-gen
KUSTOMIZE ?= $(shell pwd)/bin/kustomize
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"
OPERATOR_VERSION ?= $(shell semver get alpha)
REGISTRY_IMAGE ?= horuszup/horusec-operator:${OPERATOR_VERSION}

fmt: # Check fmt in application
	$(GOFMT) -w $(GO_FILES)

lint: # Check lint in application
    ifeq ($(wildcard $(GOLANG_CI_LINT)), $(GOLANG_CI_LINT))
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    else
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    endif

coverage: # Check coverage in application
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/coverage.sh | bash -s 100 .

test: # Run all tests in application
	$(GO) clean -testcache && $(GO) test -v ./... -timeout=2m -parallel=1 -failfast -short

fix-imports: # Setup all imports to default mode
    ifeq (, $(shell which $(GO_IMPORTS)))
		$(GO) get -u golang.org/x/tools/cmd/goimports
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    else
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    endif

security: # Run security pipeline
    ifeq (, $(shell which $(HORUSEC)))
		curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s latest
		$(HORUSEC) start -p="./" -e="true"
    else
		$(HORUSEC) start -p="./" -e="true"
    endif

build: # Build operator image
	$(GO) build -o "./tmp/bin/operator" ./cmd/app

pipeline: fmt fix-imports lint test coverage build security  # Run all processes of the pipeline

up-sample: # Up all dev dependencies kubernetes
	sh ./config/samples/sample_install_dependencies.sh

apply-sample: # Apply yaml in kubernetes
	kubectl apply -f ./config/samples/install_v2alpha1_horusecplatform.yaml

replace-sample: # Replace to re-apply yaml in kubernetes
	kubectl replace -f ./config/samples/install_v2alpha1_horusecplatform.yaml

install-semver: # Install semver binary
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/install-semver.sh | bash

docker-up-alpha: # Update alpha in docker image
	chmod +x ./deployments/scripts/update-image.sh
	./deployments/scripts/update-image.sh alpha false

docker-up-rc: # Update alpha in docker image
	chmod +x ./deployments/scripts/update-image.sh
	./deployments/scripts/update-image.sh rc false

docker-up-release: # Update release in docker image
	chmod +x ./deployments/scripts/update-image.sh
	./deployments/scripts/update-image.sh release false

docker-up-release-latest: # Update release and latest in docker image
	chmod +x ./deployments/scripts/update-image.sh
	./deployments/scripts/update-image.sh release true

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
