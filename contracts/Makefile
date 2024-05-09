##@ Help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## Build the smart contracts
	forge build --force

##@ Gen go bindings

.PHONY: gen-funder-go-bindings
gen-funder-go-bindings: ## Generatego go bindings for the funder contract.
	cat ./out/Funder.sol/Funder.json | jq -r '.abi'             > ../bindings/funder/Funder.abi
	cat ./out/Funder.sol/Funder.json | jq -r '.bytecode.object' > ../bindings/funder/Funder.bin
	abigen --abi ../bindings/funder/Funder.abi --bin ../bindings/funder/Funder.bin --pkg funder --type funder --out ../bindings/funder/funder.go
	@echo "✅ funder/funder.go generated"

.PHONY: gen-tester-go-bindings
gen-tester-go-bindings: ## Generate go bindings for the tester contracts.
	cat ./out/LoadTester.sol/LoadTester.json | jq -r '.abi'             > ../bindings/tester/LoadTester.abi
	cat ./out/LoadTester.sol/LoadTester.json | jq -r '.bytecode.object' > ../bindings/tester/LoadTester.bin
	abigen --abi ../bindings/tester/LoadTester.abi --bin ../bindings/tester/LoadTester.bin --pkg tester --type LoadTester --out ../bindings/tester/loadTester.go
	@echo "✅ tester/loadTester.go generated"

	cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.abi'             > ../bindings/tester/ConformanceTester.abi
	cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.bytecode.object' > ../bindings/tester/ConformanceTester.bin
	abigen --abi ../bindings/tester/ConformanceTester.abi --bin ../bindings/tester/ConformanceTester.bin --pkg tester --type ConformanceTester --out ../bindings/tester/conformanceTester.go
	@echo "✅ tester/conformanceTester.go generated"

.PHONY: gen-tokens-go-bindings
gen-tokens-go-bindings: ## Generate go bindings for the tokens contracts.
	cat ./out/ERC20.sol/ERC20.json | jq -r '.abi'             > ../bindings/tokens/ERC20.abi
	cat ./out/ERC20.sol/ERC20.json | jq -r '.bytecode.object' > ../bindings/tokens/ERC20.bin
	abigen --abi ../bindings/tokens/ERC20.abi --bin ../bindings/tokens/ERC20.bin --pkg tokens --type ERC20 --out ../bindings/tokens/ERC20.go
	@echo "✅ tokens/ERC20.go generated"

	cat ./out/ERC721.sol/ERC721.json | jq -r '.abi'             > ../bindings/tokens/ERC721.abi
	cat ./out/ERC721.sol/ERC721.json | jq -r '.bytecode.object' > ../bindings/tokens/ERC721.bin
	abigen --abi ../bindings/tokens/ERC721.abi --bin ../bindings/tokens/ERC721.bin --pkg tokens --type ERC721 --out ../bindings/tokens/ERC721.go
	@echo "✅ tokens/ERC721.go generated"

.PHONY: gen-go-bindings
gen-go-bindings: build gen-funder-go-bindings gen-tester-go-bindings gen-tokens-go-bindings ## Generate go bindings.
