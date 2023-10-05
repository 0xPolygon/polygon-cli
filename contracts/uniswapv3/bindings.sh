#!/bin/bash
# Generate Go bindings for Uniswap smart contracts.

v3core=v3-core-v1.0.0
v3periphery_v1_1=v3-periphery-v1.1.1
v3periphery_v1_3=v3-periphery-v1.3.0
v3staker=v3-staker-v1.0.2
v3router=v3-swap-router-v1.3.0
openzeppelin=openzeppelin-v3.4.1-solc-0.7-2

gen_go_binding() {
  local repository=$1
  local contract=$2
  local current_dir=$(pwd)
  abigen \
    --abi $repository/$contract.abi \
    --bin $repository/$contract.bin \
    --pkg uniswapv3 \
    --type $contract \
    --out $contract.go
  echo "* $contract bindings generated."
}

abigen --version

echo -e "\n🏗️  Generating go bindings..."

gen_go_binding $v3core UniswapV3Factory
gen_go_binding $v3core UniswapV3Pool

gen_go_binding ${v3periphery_v1_1} UniswapInterfaceMulticall
gen_go_binding ${v3periphery_v1_1} TickLens
gen_go_binding ${v3periphery_v1_1} NFPositionManager
gen_go_binding ${v3periphery_v1_1} V3Migrator

gen_go_binding ${v3periphery_v1_3} NFTDescriptor
gen_go_binding ${v3periphery_v1_3} NFTPositionDescriptor

gen_go_binding $v3staker UniswapV3Staker

gen_go_binding $v3router QuoterV2
gen_go_binding $v3router SwapRouter02

gen_go_binding $openzeppelin ProxyAdmin
gen_go_binding $openzeppelin TransparentUpgradeableProxy

gen_go_binding weth9 WETH9
gen_go_binding erc20 Swapper

echo "✅ Done"

echo -e "\n❗️ Make sure to read contracts/uniswapv3/README.md to update the deploy function of NFTPositionDescriptor.go and the address of the NFTDescriptor library."
