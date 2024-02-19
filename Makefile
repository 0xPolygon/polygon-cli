include scripts/lint.mk
include scripts/clients.mk
.DEFAULT_GOAL := help

INSTALL_DIR := ~/go/bin
BIN_NAME := polycli
BUILD_DIR := ./out

GIT_SHA := $(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG := $(shell git describe --tags)
DATE := $(shell date +%s)
VERSION_FLAGS=\
  -X github.com/maticnetwork/polygon-cli/cmd/version.Version=$(GIT_TAG) \
  -X github.com/maticnetwork/polygon-cli/cmd/version.Commit=$(GIT_SHA) \
  -X github.com/maticnetwork/polygon-cli/cmd/version.Date=$(DATE) \
  -X github.com/maticnetwork/polygon-cli/cmd/version.BuiltBy=makefile

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: $(BUILD_DIR)
$(BUILD_DIR): ## Create the build folder.
	mkdir -p $(BUILD_DIR)

.PHONY: build
build: $(BUILD_DIR) ## Build go binary.
	go build -ldflags "$(VERSION_FLAGS)" -o $(BUILD_DIR)/$(BIN_NAME) main.go

.PHONY: install
install: build ## Install the go binary.
	$(RM) $(INSTALL_DIR)/$(BIN_NAME)
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BIN_NAME) $(INSTALL_DIR)/

.PHONY: cross
cross: $(BUILD_DIR) ## Cross-compile go binaries using CGO.
# Notes:
# - `-s -w` enables to strip debug and suppress warnings.
# - `-linkmode external -extldflags "-static-libgo"` allows dynamic linking.
	echo "Building $(BIN_NAME)_$(GIT_TAG)_linux_arm64..."
	CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build \
			-ldflags '$(VERSION_FLAGS) -s -w -linkmode external -extldflags "-static-libgo"' \
			-tags netgo \
			-o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_linux_arm64 \
			main.go

	echo "Building $(BIN_NAME)_$(GIT_TAG)_linux_amd64..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
			-ldflags '$(VERSION_FLAGS) -s -w -linkmode external -extldflags "-static-libgo"' \
			-tags netgo \
			-o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_linux_amd64 \
			main.go

.PHONY: simplecross
simplecross: $(BUILD_DIR) ## Cross-compile go binaries without using CGO.
	GOOS=linux  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_linux_arm64  main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_darwin_arm64 main.go
	GOOS=linux  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_linux_amd64  main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_darwin_amd64 main.go

.PHONY: clean
clean: ## Clean the binary folder.
	$(RM) -r $(BUILD_DIR)

##@ Test

.PHONY: test
test: ## Run tests.
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

##@ Generation

.PHONY: gen
gen: gen-doc gen-proto gen-go-bindings gen-loadtest-modes gen-json-rpctypes ## Generate everything.

.PHONY: gen-doc
gen-doc: ## Generate documentation for `polycli`.
	go run docutil/*.go

.PHONY: gen-proto
gen-proto: ## Generate protobuf stubs.
	protoc --proto_path=proto --go_out=proto/gen/pb --go_opt=paths=source_relative $(wildcard proto/*.proto)

.PHONY: gen-go-bindings
gen-go-bindings: ## Generate go bindings for smart contracts.
	cd contracts && make gen-go-bindings

.PHONY: gen-loadtest-modes
gen-loadtest-modes: ## Generate loadtest modes strings.
	cd cmd/loadtest && stringer -type=loadTestMode

.PHONY: gen-json-rpctypes
gen-json-rpctypes: ## Generate JSON rpc types.
	./scripts/rpctypes.sh rpctypes/jsonschemas/
