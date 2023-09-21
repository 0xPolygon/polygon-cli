define gen_go_binding
	solc $1/$1.sol --abi --bin --output-dir $1 --overwrite
	abigen --abi $1/$1.abi --bin $1/$1.bin --pkg tokens --type $1 --out $1.go
endef


.PHONY: gen
gen: ## Generate go bindings
	$(call gen_go_binding,ERC20)
	$(call gen_go_binding,ERC721)