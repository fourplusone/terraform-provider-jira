PKG := $(shell go list ./... | grep -v vendor)
TEST := $(shell go list ./... |grep -v vendor)


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
