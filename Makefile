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
VERSION ?= 2.0.0-alpha.1
IMAGE_TAG_BASE ?= horuszup/horusec-operator
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)
IMG ?= $(IMAGE_TAG_BASE):v$(VERSION)

fmt:
	$(GOFMT) -w $(GO_FILES)

lint:
    ifeq ($(wildcard $(GOLANG_CI_LINT)), $(GOLANG_CI_LINT))
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    else
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
		$(GOLANG_CI_LINT) run -v --timeout=5m -c .golangci.yml ./...
    endif

coverage:
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/coverage.sh | bash -s 100 .

test:
	$(GO) clean -testcache && $(GO) test -v ./... -timeout=2m -parallel=1 -failfast -short

fix-imports:
    ifeq (, $(shell which $(GO_IMPORTS)))
		$(GO) get -u golang.org/x/tools/cmd/goimports
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    else
		$(GO_IMPORTS) -local $(GO_IMPORTS_LOCAL) -w $(GO_FILES)
    endif

security:
    ifeq (, $(shell which $(HORUSEC)))
		curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s latest
		$(HORUSEC) start -p="./" -e="true"
    else
		$(HORUSEC) start -p="./" -e="true"
    endif

build:
	$(GO) build -o "./tmp/bin/api" ./cmd/app/main.go

update-swagger:
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/update-swagger.sh | bash

pipeline: fmt fix-imports lint test coverage build security

up-sample-dependencies:
	sh ./config/samples/sample_install_dependencies.sh

test-sample:
	kubectl apply -f ./config/samples/install_v2alpha1_horusec.yaml

docker-build:
	docker build -t horuszup/horusec-operator:v2.0.0-alpha.1 -f ./Dockerfile .

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

kustomize:
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

controller-gen:
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config="./config/crd/bases"

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy:
	$(KUSTOMIZE) build config/default | kubectl delete -f -
