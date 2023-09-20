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
  && echo ✅ $contract bindings generated.
}

abigen --version
gen_go_binding utils WETH9
gen_go_binding v3-core UniswapV3Factory
gen_go_binding v3-periphery UniswapInterfaceMulticall
gen_go_binding openzeppelin-contracts ProxyAdmin
gen_go_binding openzeppelin-contracts TransparentUpgradeableProxy
gen_go_binding v3-periphery TickLens
gen_go_binding v3-periphery NFTDescriptor
gen_go_binding v3-periphery NonfungibleTokenPositionDescriptor
