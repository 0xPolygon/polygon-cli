gen:
	FOUNDRY_PROFILE=lite forge build --contracts src/loadtester
	$(call gen_go_binding,LoadTester,src/loadtester,loadtester,src/loadtester)

	forge build --skip LoadTester
	$(call gen_go_binding,ERC20,src/tokens/ERC20,tokens,src/tokens)
	$(call gen_go_binding,ERC721,src/tokens/ERC721,tokens,src/tokens)
	$(call gen_go_binding,Swapper,src/tokens/Swapper,tokens,src/tokens)
	$(call gen_go_binding,WETH9,src/tokens/WETH9,tokens,src/tokens)

# Generate go binding.
# - $1: smart contract name.
# - $2: smart contract folder path.
# - $3: go package name.
# - $4: go file output path.
define gen_go_binding
	echo "Generating ABI, bytecode and go bindings for $1..."
	cat ./out/$1.sol/$1.json | jq -r '.abi' 						> ./$2/$1.abi
	cat ./out/$1.sol/$1.json | jq -r '.bytecode.object' > ./$2/$1.bin
	abigen --abi ./$2/$1.abi --bin ./$2/$1.bin --pkg $3 --type $1 --out ./$4/$1.go
endef
