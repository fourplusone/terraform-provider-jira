PKG := $(shell go list ./... | grep -v vendor)
TEST := $(shell go list ./... |grep -v vendor)

HOSTNAME=idealo.de
NAMESPACE=pt
NAME=jira
BINARY=terraform-provider-${NAME}
VERSION=0.1
OS_ARCH=darwin_amd64


.PHONY: test

.DEFAULT_GOAL := help
help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download

build: ## Build
	go build .

release: ## Build the go binaries for various platform
	./scripts/release.sh

test: ## Run tests
	TF_ACC=1 go test -v $(TEST)
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
