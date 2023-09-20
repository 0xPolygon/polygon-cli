#!/bin/bash
# Generate Go bindings for Uniswap smart contracts.

gen_go_binding() {
  local repository=$1
  local contract=$2
  local current_dir=$(pwd)
  abigen \
    --abi $repository/$contract.abi \
    --bin $repository/$contract.bin \
    --pkg uniswapv3 \
    --type $contract \
    --out $contract.go \
  && echo "* $contract bindings generated."
}

abigen --version

echo -e "\n🏗️  Generating go bindings..."
gen_go_binding v3-core UniswapV3Factory

gen_go_binding v3-periphery UniswapInterfaceMulticall
gen_go_binding v3-periphery TickLens
gen_go_binding v3-periphery NFTDescriptor
gen_go_binding v3-periphery NonfungibleTokenPositionDescriptor
gen_go_binding v3-periphery NonfungiblePositionManager
gen_go_binding v3-periphery V3Migrator

gen_go_binding v3-staker UniswapV3Staker

gen_go_binding v3-swap-router QuoterV2
gen_go_binding v3-swap-router SwapRouter02

gen_go_binding openzeppelin ProxyAdmin
gen_go_binding openzeppelin TransparentUpgradeableProxy

gen_go_binding weth9 WETH9

echo "✅ Done"
