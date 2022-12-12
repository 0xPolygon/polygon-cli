
INSTALL_DIR:=~/go/bin/
BIN_NAME:=polycli
BUILD_DIR:=./out

GIT_SHA:=$(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG:=$(shell git describe --tags)
CUR_DATE:=$(shell date +%s)

# strip debug and supress warnings
LD_FLAGS=-s -w
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Version=$(GIT_TAG)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Commit=$(GIT_SHA)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Date=$(CUR_DATE)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.BuiltBy=makefile\"
STATIC_LD_FLAGS=$(LD_FLAGS) -extldflags=-static

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: $(BUILD_DIR)
$(BUILD_DIR): ## Create the binary folder.
	mkdir -p $(BUILD_DIR)

.PHONY: run
run: ## Run the go program.
	go run main.go

.PHONY: build
build: $(BUILD_DIR) ## Build go binary.
	go build -o $(BUILD_DIR)/$(BIN_NAME) main.go

.PHONY: cross
cross: $(BUILD_DIR) ## Cross-compile go binaries using CGO.
	env CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-arm64-$(BIN_NAME) main.go
	env                          CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-amd64-$(BIN_NAME) main.go
  # mac builds - this will be functional but will still have secp issues
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-arm64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-amd64-$(BIN_NAME) main.go

.PHONY: simplecross
simplecross: $(BUILD_DIR) ## Cross-compile go binaries without using CGO.
	env GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/linux-arm64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/darwin-arm64-$(BIN_NAME) main.go
	env GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/darwin-amd64-$(BIN_NAME) main.go

.PHONY: install
install: build ## Install the go binary.
	$(RM) $(INSTALL_DIR)/$(BIN_NAME)
	cp $(BUILD_DIR)/$(BIN_NAME) $(INSTALL_DIR)

.PHONY: clean
clean: ## Clean the binary folder.
	$(RM) -r $(BUILD_DIR)

##@ Lint

.PHONY: tidy
tidy: ## Add missing and remove unused modules.
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

# shadow reports shadowed variables
# https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow
# go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
.PHONY: vet
vet: ## Run go vet and shadow against code.
	go vet ./...
	shadow ./...

# golangci-lint runs gofmt, govet, staticcheck and other linters
# https://golangci-lint.run/usage/install/#local-installation
.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint against code.
	golangci-lint run --fix

.PHONY: lint
lint: tidy vet golangci-lint ## Run all the linter tools against code.

##@ Test

.PHONY: test
test: lint ## Run tests.
	go test ./... -coverprofile=coverage.out
	go tool cover -func coverage.out

