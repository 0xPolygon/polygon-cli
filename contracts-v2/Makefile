gen:
	forge build
	$(call gen_go_binding,ERC20,tokens)
	$(call gen_go_binding,ERC721,tokens)


# Generate go binding.
# - $1: smart contract name
# - $2: smart contract folder under `./src`
define gen_go_binding
	cat ./out/$1.sol/$1.json | jq -r '.abi' 						> ./src/$2/$1/$1.abi
	cat ./out/$1.sol/$1.json | jq -r '.bytecode.object' > ./src/$2/$1/$1.bin
	abigen --abi ./src/$2/$1/$1.abi --bin ./src/$2/$1/$1.bin --pkg contracts --type $1 --out ./src/$2/$1/$1.go
endef
