#!/bin/bash
# This script builds UniswapV3 core, periphery and swap-router contracts.

solc --version
current_dir=$(pwd)

# Build v3-core contracts.
echo -e "\n🏗️  Building v3-core contracts..."
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
build_v3_periphery_contract() {
	contract_path=$1
	solc \
		@uniswap=$current_dir/v3-periphery/node_modules/@uniswap \
		@openzeppelin=$current_dir/v3-periphery/node_modules/@openzeppelin \
		base64-sol=$current_dir/v3-periphery/node_modules/base64-sol \
		../interfaces=$current_dir/v3-periphery/contracts/interfaces \
		--evm-version istanbul \
		--optimize \
		--optimize-runs 2000 \
		--abi v3-periphery/contracts/$contract_path \
		--bin v3-periphery/contracts/$contract_path \
		--output-dir tmp/v3-periphery
}

echo -e "\n🏗️  Building v3-periphery contracts..."
rm -rf v3-periphery
git clone https://github.com/Uniswap/v3-periphery.git
pushd v3-periphery
yarn install
popd
build_v3_periphery_contract lens/UniswapInterfaceMulticall.sol
build_v3_periphery_contract lens/TickLens.sol
rm -rf v3-periphery
mkdir v3-periphery
mv tmp/v3-periphery/* v3-periphery
rm -rf tmp
echo "✅ Successfully built v3-periphery contracts..."

# Build openzeppelin contracts.
echo -e "\n🏗️  Building openzeppelin contracts..."
rm -rf openzeppelin-contracts
git clone https://github.com/OpenZeppelin/openzeppelin-contracts.git --branch v3.4.1-solc-0.7-2
solc \
	../access=$current_dir/openzeppelin-contracts/contracts/access \
	../utils=$current_dir/openzeppelin-contracts/contracts/utils \
  --abi openzeppelin-contracts/contracts/proxy/ProxyAdmin.sol \
  --bin openzeppelin-contracts/contracts/proxy/ProxyAdmin.sol \
  --output-dir tmp/openzeppelin-contracts
rm -rf openzeppelin-contracts
mkdir openzeppelin-contracts
mv tmp/openzeppelin-contracts/* openzeppelin-contracts
rm -rf tmp
echo "✅ Successfully built openzeppelin contracts..."
