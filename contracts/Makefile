# Generate go binding.
# - $1: dir
# - $2: name
# - $3: pkg
define gen_go_binding
	solc $1/$2.sol --abi --bin --output-dir $1 --overwrite
	abigen --abi $1/$2.abi --bin $1/$2.bin --pkg $3 --type $2 --out $3/$2.go
endef


.PHONY: gen
gen: ## Generate go bindings
	$(call gen_go_binding,tokens/ERC20,ERC20,tokens)
	$(call gen_go_binding,tokens/ERC721,ERC721,tokens)
	$(call gen_go_binding,loadtester,LoadTester,loadtester)