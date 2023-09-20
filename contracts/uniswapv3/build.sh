#!/bin/bash
# This script builds UniswapV3 core, periphery and swap-router contracts.

# Make sure the local chain is started.
wait_for_service() {
	ip=$1
	port=$2
	name=$3
	{
		while ! echo -n > "/dev/tcp/$ip/$port"; do
			>&2 echo "$name port is not open yet. Waiting for 5 seconds"
			sleep 5
		done
	} 2>/dev/null
	>&2 echo "$name port is now open."
}

>&2 echo "Starting status checking"
wait_for_service "127.0.0.1" 8545 "Local RPC"

# Build contracts.

solc --version
current_dir=$(pwd)

## Build WETH9 contract.
if [ "$1" -eq 1 ]; then
	echo -e "\n🏗️  Building WETH9 contract..."
	git clone https://github.com/gnosis/canonical-weth.git
	rm -rf utils
	mkdir utils
	cat canonical-weth/build/contracts/WETH9.json | jq .abi > utils/WETH9.abi
	cat canonical-weth/build/contracts/WETH9.json | jq -r .bytecode > utils/WETH9.bin
	rm -rf canonical-weth
	echo "✅ Successfully built WETH9 contract..."
fi

## Build v3-core contracts.
if [ "$1" -eq 2 ]; then
	echo -e "\n🏗️  Building v3-core contracts..."
	rm -rf v3-core
	git clone https://github.com/Uniswap/v3-core.git
	solc \
		v3-core/contracts/UniswapV3Factory.sol \
		--optimize \
		--optimize-runs 200 \
		--abi \
		--bin \
		--output-dir tmp/v3-core \
		--overwrite
	rm -rf v3-core
	mkdir v3-core
	mv tmp/v3-core/* v3-core
	rm -rf tmp
	echo "✅ Successfully built v3-core contracts..."
fi

## Build v3-periphery contracts.
if [ "$1" -eq 3 ]; then
	echo -e "\n🏗️  Building v3-periphery contracts..."
	rm -rf v3-periphery
	git clone https://github.com/Uniswap/v3-periphery.git
	pushd v3-periphery
	yarn install
	popd
	solc \
		v3-periphery/contracts/lens/UniswapInterfaceMulticall.sol \
		v3-periphery/contracts/lens/TickLens.sol \
		v3-periphery/contracts/libraries/NFTDescriptor.sol \
		v3-periphery/contracts/NonfungibleTokenPositionDescriptor.sol \
		v3-periphery/contracts/NonfungiblePositionManager.sol \
		@uniswap=$current_dir/v3-periphery/node_modules/@uniswap \
		@openzeppelin=$current_dir/v3-periphery/node_modules/@openzeppelin \
		base64-sol=$current_dir/v3-periphery/node_modules/base64-sol \
		../interfaces=$current_dir/v3-periphery/contracts/interfaces \
		--evm-version istanbul \
		--optimize \
		--optimize-runs 200 \
		--abi \
		--bin \
		--output-dir tmp/v3-periphery \
		--overwrite

	# We need to deloy the NFTDescriptor library, retrieve its address and link it inside
	# NonfungibleTokenPositionDescriptor bytecode. This is required to generate the Go binding.
	nft_descriptor_lib_address=$(cast send \
		--rpc-url http://localhost:8545 \
		--chain 1337 \
		--from 0x85da99c8a7c2c95964c8efd687e95e632fc533d6 \
		--private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
		--json \
		--create \
		"$(cat tmp/v3-periphery/NFTDescriptor.bin)" \
		| jq -r .contractAddress)
	solc \
		--link \
		--libraries v3-periphery/contracts/libraries/NFTDescriptor.sol:NFTDescriptor:$nft_descriptor_lib_address \
		tmp/v3-periphery/NonfungibleTokenPositionDescriptor.bin

	rm -rf v3-periphery
	mkdir v3-periphery
	mv tmp/v3-periphery/* v3-periphery
	rm -rf tmp
	echo "✅ Successfully built v3-periphery contracts..."
fi

## Build openzeppelin contracts.
if [ "$1" -eq 4 ]; then
	echo -e "\n🏗️  Building openzeppelin contracts..."
	rm -rf openzeppelin-contracts
	git clone https://github.com/OpenZeppelin/openzeppelin-contracts.git --branch v3.4.1-solc-0.7-2
	solc \
		openzeppelin-contracts/contracts/proxy/ProxyAdmin.sol \
		openzeppelin-contracts/contracts/proxy/TransparentUpgradeableProxy.sol \
		../access=$current_dir/openzeppelin-contracts/contracts/access \
		../utils=$current_dir/openzeppelin-contracts/contracts/utils \
		--abi \
		--bin \
		--output-dir tmp/openzeppelin-contracts \
		--overwrite
	rm -rf openzeppelin-contracts
	mkdir openzeppelin-contracts
	mv tmp/openzeppelin-contracts/* openzeppelin-contracts
	rm -rf tmp
	echo "✅ Successfully built openzeppelin contracts..."
fi
