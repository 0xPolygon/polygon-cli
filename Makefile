
INSTALL_DIR:=~/go/bin
BIN_DIR=./bin
BIN_NAME:=polycli
BUILD_DIR:=./out

GIT_SHA:=$(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG:=$(shell git describe --tags)
CUR_DATE:=$(shell date +%s)

# strip debug and supress warnings
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
	go build -ldflags "-s -w -X \"github.com/maticnetwork/polygon-cli/cmd/version.Version=dev ($(GIT_SHA))\"" -o $(BUILD_DIR)/$(BIN_NAME) main.go

.PHONY: install
install: build ## Install the go binary.
	$(RM) $(INSTALL_DIR)/$(BIN_NAME)
	cp $(BUILD_DIR)/$(BIN_NAME) $(INSTALL_DIR)

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
	golangci-lint run --fix --timeout 5m

.PHONY: lint
lint: tidy vet golangci-lint ## Run all the linter tools against code.

##@ Clients
PORT?=8545

.PHONY: $(BIN_DIR)
$(BIN_DIR): ## Create the binary folder.
	mkdir -p $(BIN_DIR)

.PHONY: geth
geth: $(BIN_DIR) ## Start a local geth node.
	geth --dev --dev.period 2 --http --http.addr localhost --http.port $(PORT) --http.api admin,debug,web3,eth,txpool,personal,miner,net --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 --dev.gaslimit 50000000

.PHONY: avail
avail: $(BIN_DIR) ## Start a local avail node.
	avail --dev --rpc-port $(PORT)

##@ Test

.PHONY: test
test: ## Run tests.
	go test ./... -coverprofile=coverage.out
	go tool cover -func coverage.out

LOADTEST_ACCOUNT=0x85da99c8a7c2c95964c8efd687e95e632fc533d6
LOADTEST_FUNDING_AMOUNT_ETH=5000
eth_coinbase := $(shell curl -s -H 'Content-Type: application/json' -d '{"jsonrpc": "2.0", "id": 2, "method": "eth_coinbase", "params": []}' http://127.0.0.1:${PORT} | jq -r ".result")
hex_funding_amount := $(shell echo "obase=16; ${LOADTEST_FUNDING_AMOUNT_ETH}*10^18" | bc)
.PHONY: geth-loadtest
geth-loadtest: build ## Fund test account with 5k ETH and run loadtest against an EVM/Geth chain.
	curl -H "Content-Type: application/json" -d '{"jsonrpc":"2.0", "method":"eth_sendTransaction", "params":[{"from": "${eth_coinbase}","to": "${LOADTEST_ACCOUNT}","value": "0x${hex_funding_amount}"}], "id":1}' http://127.0.0.1:${PORT}
	sleep 5
	$(BUILD_DIR)/$(BIN_NAME) loadtest --verbosity 700 --chain-id 1337 --concurrency 1 --requests 1000 --rate-limit 5 --mode c http://127.0.0.1:$(PORT)

.PHONY: avail-loadtest
avail-loadtest: build ## Run loadtest against an Avail chain.
	$(BUILD_DIR)/$(BIN_NAME) loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 5 --mode t --data-avail http://127.0.0.1:$(PORT)