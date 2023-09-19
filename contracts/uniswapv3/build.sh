#!/bin/bash
# This script builds UniswapV3 core, periphery and swap-router contracts.

solc --version
current_dir=$(pwd)

# Build v3-core contracts.
echo "\n🏗️  Building v3-core contracts..."
rm -rf v3-core
git clone https://github.com/Uniswap/v3-core.git
solc \
  --optimize \
  --optimize-runs 200 \
  --abi v3-core/contracts/UniswapV3Factory.sol \
  --bin v3-core/contracts/UniswapV3Factory.sol \
  --output-dir tmp/v3-core
rm -rf v3-core
mkdir v3-core
mv tmp/v3-core/* v3-core
rm -rf tmp
echo "✅ Successfully built v3-core contracts..."

# Build v3-periphery contracts.
echo "\n🏗️  Building v3-periphery contracts..."
rm -rf v3-periphery
git clone https://github.com/Uniswap/v3-periphery.git
pushd v3-periphery
yarn install
popd
solc \
	@uniswap=$current_dir/v3-periphery/node_modules/@uniswap \
	@openzeppelin=$current_dir/v3-periphery/node_modules/@openzeppelin \
	base64-sol=$current_dir/v3-periphery/node_modules/base64-sol \
	--evm-version istanbul \
	--optimize \
	--optimize-runs 2000 \
	--abi v3-periphery/contracts/lens/UniswapInterfaceMulticall.sol \
	--bin v3-periphery/contracts/lens/UniswapInterfaceMulticall.sol \
	--output-dir tmp/v3-periphery
rm -rf v3-periphery
mkdir v3-periphery
mv tmp/v3-periphery/* v3-periphery
rm -rf tmp
echo "✅ Successfully built v3-periphery contracts..."
