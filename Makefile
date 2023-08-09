include scripts/lint.mk
include scripts/clients.mk
.DEFAULT_GOAL := help

INSTALL_DIR := ~/go/bin
BIN_NAME := polycli
BUILD_DIR := ./out

GIT_SHA := $(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG := $(shell git describe --tags)
CUR_DATE := $(shell date +%s)

# Strip debug and suppress warnings.
LD_FLAGS=-s -w
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/cmd/version.Version=$(GIT_TAG)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/cmd/version.Commit=$(GIT_SHA)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/cmd/version.Date=$(CUR_DATE)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/cmd/version.BuiltBy=makefile\"
STATIC_LD_FLAGS=$(LD_FLAGS) -extldflags=-static

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: $(BUILD_DIR)
$(BUILD_DIR): ## Create the build folder.
	mkdir -p $(BUILD_DIR)

.PHONY: generate
generate: ## Generate protobuf stubs.
	protoc --proto_path=proto --go_out=proto/gen/pb --go_opt=paths=source_relative $(wildcard proto/*.proto)

.PHONY: build
build: $(BUILD_DIR) ## Build go binary.
	go build -ldflags "-X \"github.com/maticnetwork/polygon-cli/cmd/version.Version=dev ($(GIT_SHA))\"" -o $(BUILD_DIR)/$(BIN_NAME) main.go

.PHONY: install
install: build ## Install the go binary.
	$(RM) $(INSTALL_DIR)/$(BIN_NAME)
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BIN_NAME) $(INSTALL_DIR)/

.PHONY: cross
cross: $(BUILD_DIR) ## Cross-compile go binaries using CGO.
	env CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-arm64-$(BIN_NAME) main.go
	env                          CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-amd64-$(BIN_NAME) main.go
# MAC builds - this will be functional but will still have secp issues.
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-arm64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-amd64-$(BIN_NAME) main.go

.PHONY: simplecross
simplecross: $(BUILD_DIR) ## Cross-compile go binaries without using CGO.
	env GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/linux-arm64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/darwin-arm64-$(BIN_NAME) main.go
	env GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/darwin-amd64-$(BIN_NAME) main.go

.PHONY: clean
clean: ## Clean the binary folder.
	$(RM) -r $(BUILD_DIR)

##@ Test

.PHONY: test
test: ## Run tests.
	go test ./... -coverprofile=coverage.out
	go tool cover -func coverage.out

##@ Documentation

.PHONY: gen-doc
gen-doc: ## Generate documentation for `polycli`.
	go run docutil/*.go
