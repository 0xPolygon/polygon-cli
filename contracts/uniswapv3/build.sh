#!/bin/bash
# This script builds UniswapV3 contracts.

mode="${1:-0}"

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
	>&2 echo "✅ $name port is now open."
}

>&2 echo "Starting status checking"
wait_for_service "127.0.0.1" 8545 "Local RPC"
echo

# Build contracts.
build_contracts() {
	repository=$1
	url=$2
	branch=$3
	contracts=$4
	echo -e "\n🏗️  Building $repository contracts..."
	rm -rf ./$repository/*
	git clone --branch $branch $url ./tmp/$repository
	pushd tmp/$repository
	yarn install
	popd

	new_array=()
	for contract in "${contracts[@]}"; do
		new_array+=("tmp/$repository/contracts/$contract")
	done

	for element in "${new_array[@]}"; do
    echo "$element"
	done

	solc "${new_array[@]}" \
		@uniswap=$current_dir/tmp/$repository/node_modules/@uniswap \
		@openzeppelin=$current_dir/tmp/$repository/node_modules/@openzeppelin \
		base64-sol=$current_dir/tmp/$repository/node_modules/base64-sol \
		../interfaces=$current_dir/tmp/$repository/contracts/interfaces \
		../libraries=$current_dir/tmp/$repository/contracts/libraries \
		--evm-version istanbul \
		--optimize \
		--optimize-runs 200 \
		--abi \
		--bin \
		--output-dir ./$repository \
		--overwrite

	rm -rf ./tmp/$repository
	echo "✅ Successfully built $repository contracts..."
}

solc --version
current_dir=$(pwd)

## Build v3-core contracts.
if [ "$mode" -eq 1 ] || [ "$mode" -eq 0 ]; then
	contracts=("UniswapV3Factory.sol" "UniswapV3Pool.sol")
	build_contracts v3-core https://github.com/Uniswap/v3-core.git v1.0.0 $contracts
fi

## Build v3-periphery contracts.
if [ "$mode" -eq 2 ] || [ "$mode" -eq 0 ]; then
	contracts=("lens/UniswapInterfaceMulticall.sol" "lens/TickLens.sol" "libraries/NFTDescriptor.sol" "NonfungibleTokenPositionDescriptor.sol" "NonfungiblePositionManager.sol" "V3Migrator.sol")
	build_contracts v3-periphery https://github.com/Uniswap/v3-periphery.git v1.3.0 $contracts

	# We need to deloy the NFTDescriptor library, retrieve its address and link it inside
	# NonfungibleTokenPositionDescriptor bytecode. This is required to generate the Go binding.
	nft_descriptor_lib_address=$(cast send \
		--rpc-url http://localhost:8545 \
		--chain 1337 \
		--from 0x85da99c8a7c2c95964c8efd687e95e632fc533d6 \
		--private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
		--json \
		--create \
		"$(cat v3-periphery/NFTDescriptor.bin)" \
		| jq -r .contractAddress)
	solc \
		--link \
		--libraries tmp/v3-periphery/contracts/libraries/NFTDescriptor.sol:NFTDescriptor:$nft_descriptor_lib_address \
		v3-periphery/NonfungibleTokenPositionDescriptor.bin
fi

## Build v3-staker contracts.
if [ "$mode" -eq 3 ] || [ "$mode" -eq 0 ]; then
	contracts=("UniswapV3Staker.sol")
	build_contracts v3-staker https://github.com/Uniswap/v3-staker.git v1.0.2 $contracts
fi

## Build v3-swap-router contracts.
if [ "$mode" -eq 4 ] || [ "$mode" -eq 0 ]; then
	contracts=("lens/QuoterV2.sol" "SwapRouter02.sol")
	build_contracts v3-swap-router https://github.com/Uniswap/swap-router-contracts.git v1.3.0 $contracts
fi

## Build openzeppelin contracts.
if [ "$mode" -eq 5 ] || [ "$mode" -eq 0 ]; then
	contracts=("proxy/ProxyAdmin.sol" "proxy/TransparentUpgradeableProxy.sol")
	build_contracts openzeppelin https://github.com/OpenZeppelin/openzeppelin-contracts.git v3.4.1-solc-0.7-2 $contracts
fi

## Build WETH9 contract.
if [ "$mode" -eq 6 ] || [ "$mode" -eq 0 ]; then
	echo -e "\n🏗️  Building WETH9 contract..."
	git clone https://github.com/gnosis/canonical-weth.git
	rm -rf weth9
	mkdir weth9
	cat canonical-weth/build/contracts/WETH9.json | jq .abi > weth9/WETH9.abi
	cat canonical-weth/build/contracts/WETH9.json | jq -r .bytecode > weth9/WETH9.bin
	rm -rf canonical-weth
	echo "✅ Successfully built WETH9 contract..."
fi
