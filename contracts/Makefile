# Generate go binding.
# - $1: input_dir
# - $2: name
# - $3: pkg
# - $4: output_dir
define gen_go_binding
	solc $1/$2.sol --abi --bin --output-dir $1 --overwrite
	abigen --abi $1/$2.abi --bin $1/$2.bin --pkg $3 --type $2 --out $4/$2.go
endef


.PHONY: gen-go-bindings
gen-go-bindings: ## Generate go bindings
	$(call gen_go_binding,tokens/ERC20,ERC20,tokens,tokens)
	$(call gen_go_binding,tokens/ERC721,ERC721,tokens,tokens)
	$(call gen_go_binding,loadtester,LoadTester,contracts,.)

# solc tokens/ERC20/ERC20.sol --abi --bin --output-dir tokens/ERC20 --overwrite
# abigen --abi tokens/ERC20/ERC20.abi --bin tokens/ERC20/ERC20.bin --pkg tokens --type ERC20 --out tokens/ERC20.go

# solc loadtester/LoadTester.sol --abi --bin --output-dir loadtester --overwrite
# abigen --abi loadtester/LoadTester.abi --bin loadtester/LoadTester.bin --pkg contracts --type LoadTester --out LoadTester.go