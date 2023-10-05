#!/bin/bash
# This script builds UniswapV3 contracts.
set -e
mode="${1:-all}"

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

	# Clone repository.
	git clone --branch $branch $url ./tmp/$repository-$branch

	# Install dependencies.
	pushd tmp/$repository-$branch
	yarn install
	popd

	# Update contract's array path.
	new_array=()
	for contract in "${contracts[@]}"; do
		new_array+=("tmp/$repository-$branch/contracts/$contract.sol")
	done

	# Remove old artefacts.
	rm -rf ./$repository-$branch/*

	# Compile contracts.
	for element in "${new_array[@]}"; do
    echo "$element"
	done
	solc "${new_array[@]}" \
		@uniswap=$current_dir/tmp/$repository-$branch/node_modules/@uniswap \
		@openzeppelin=$current_dir/tmp/$repository-$branch/node_modules/@openzeppelin \
		base64-sol=$current_dir/tmp/$repository-$branch/node_modules/base64-sol \
		../interfaces=$current_dir/tmp/$repository-$branch/contracts/interfaces \
		../libraries=$current_dir/tmp/$repository-$branch/contracts/libraries \
		--evm-version istanbul \
		--optimize \
		--optimize-runs 200 \
		--abi \
		--bin \
		--output-dir ./$repository-$branch \
		--overwrite

	# Clean up.
	rm -rf ./tmp/$repository-$branch

	echo "✅ Successfully built $repository-$branch contracts..."
}

# Use solc 0.7.6.
solc-select use 0.7.6
solc --version
current_dir=$(pwd)

# Clean up
rm -rf ./tmp

## Build v3-core contracts.
if [ "$mode" == "v3-core" ] || [ "$mode" == "all" ]; then
	contracts=("UniswapV3Factory" "UniswapV3Pool")
	build_contracts v3-core https://github.com/Uniswap/v3-core.git v1.0.0 $contracts
fi

## Build v3-periphery v1.1.1 contracts
if [ "$mode" == "v3-periphery" ] || [ "$mode" == "all" ]; then
	repository=v3-periphery
	branch=v1.1.1
	contracts=("lens/UniswapInterfaceMulticall" "lens/TickLens" "NonfungiblePositionManager" "V3Migrator")
	build_contracts $repository https://github.com/Uniswap/v3-periphery.git $branch $contracts

	# Rename the NonfungiblePositionManager contract.
	mv $repository-$branch/NonfungiblePositionManager.abi $repository-$branch/NFPositionManager.abi
	mv $repository-$branch/NonfungiblePositionManager.bin $repository-$branch/NFPositionManager.bin
fi

## Build v3-periphery v1.3.0 contracts
if [ "$mode" == "v3-periphery" ] || [ "$mode" == "all" ]; then
	repository=v3-periphery
	branch=v1.3.0
	contracts=("libraries/NFTDescriptor" "NonfungibleTokenPositionDescriptor")
	build_contracts $repository https://github.com/Uniswap/v3-periphery.git $branch $contracts

	# Rename the NonfungibleTokenPositionDescriptor contract.
	mv $repository-$branch/NonfungibleTokenPositionDescriptor.abi $repository-$branch/NFTPositionDescriptor.abi
	mv $repository-$branch/NonfungibleTokenPositionDescriptor.bin $repository-$branch/NFTPositionDescriptor.bin

	# We need to deloy the NFTDescriptor library, retrieve its address and link it inside
	# NonfungibleTokenPositionDescriptor bytecode. This is required to generate the Go binding.
	nft_descriptor_lib_address=$(cast send \
		--rpc-url http://localhost:8545 \
		--chain 1337 \
		--from 0x85da99c8a7c2c95964c8efd687e95e632fc533d6 \
		--private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
		--json \
		--create \
		"$(cat $repository-$branch/NFTDescriptor.bin)" \
		| jq -r .contractAddress)
	solc \
		--link \
		--libraries tmp/$repository-$branch/contracts/libraries/NFTDescriptor.sol:NFTDescriptor:$nft_descriptor_lib_address \
		$repository-$branch/NFTPositionDescriptor.bin
	echo "✍️  NFTDescriptor library address: $nft_descriptor_lib_address"
fi

## Build v3-staker contracts.
if [ "$mode" == "v3-staker" ] || [ "$mode" == "all" ]; then
	contracts=("UniswapV3Staker")
	build_contracts v3-staker https://github.com/Uniswap/v3-staker.git v1.0.2 $contracts
fi

## Build v3-swap-router contracts.
if [ "$mode" == "v3-swap-router" ] || [ "$mode" == "all" ]; then
	contracts=("lens/QuoterV2" "SwapRouter02")
	build_contracts v3-swap-router https://github.com/Uniswap/swap-router-contracts.git v1.3.0 $contracts
fi

## Build openzeppelin contracts.
if [ "$mode" == "openzeppelin" ] || [ "$mode" == "all" ]; then
	contracts=("proxy/ProxyAdmin" "proxy/TransparentUpgradeableProxy")
	build_contracts openzeppelin https://github.com/OpenZeppelin/openzeppelin-contracts.git v3.4.1-solc-0.7-2 $contracts
fi

## Build WETH9 contract.
if [ "$mode" == "weth9" ] || [ "$mode" == "all" ]; then
	echo -e "\n🏗️  Building WETH9 contract..."
	git clone https://github.com/gnosis/canonical-weth.git
	rm -rf weth9
	mkdir weth9
	cat canonical-weth/build/contracts/WETH9.json | jq .abi > weth9/WETH9.abi
	cat canonical-weth/build/contracts/WETH9.json | jq -r .bytecode > weth9/WETH9.bin
	rm -rf canonical-weth
	echo "✅ Successfully built WETH9 contract..."
fi

## Build ERC20 contract.
if [ "$mode" == "erc20" ] || [ "$mode" == "all" ]; then
	echo -e "\n🏗️  Building ERC20 contract..."

	# Clone repository
	git clone --branch v3.4.1-solc-0.7-2 https://github.com/OpenZeppelin/openzeppelin-contracts.git ./tmp/openzeppelin-contracts

	# Install dependencies.
	pushd tmp/openzeppelin-contracts
	yarn install
	popd

	# Compile contract.
	solc erc20/Swapper.sol \
		@openzeppelin=$current_dir/tmp/openzeppelin-contracts/contracts \
		--evm-version istanbul \
		--optimize \
		--optimize-runs 200 \
		--abi \
		--bin \
		--output-dir erc20 \
		--overwrite

	# Clean up.
	#rm -rf ./tmp/openzeppelin-contracts

	echo "✅ Successfully built ERC20 contract..."
fi
