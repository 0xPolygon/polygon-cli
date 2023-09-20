# 🦄 UniswapV3 contracts

Simple steps to build UniswapV3 contracts and generate go bindings.

1. Make sure you have `solc@0.7.6` installed. This is required to build UniswapV3 contracts. A handy way to manage `solc` versions is to use [crytic/solc-select](https://github.com/crytic/solc-select).

2. Make sure you have a local RPC running. Some contract require their libraries to be deployed in order to link them directly in the bytecode.

3. Build UniswapV3 contracts using `build.sh`.

```sh
$ ./build.sh
Starting status checking
✅ Local RPC port is now open.

solc, the solidity compiler commandline interface
Version: 0.7.6+commit.7338295f.Darwin.appleclang

🏗️  Building v3-core contracts...
...
```

3. Generate Go bindings for those contracts using `bindings.sh`.

```sh
$ ./bindings.sh
abigen version 1.13.1-stable

🏗️  Generating go bindings...
* UniswapV3Factory bindings generated.
* UniswapInterfaceMulticall bindings generated.
* TickLens bindings generated.
* NFTDescriptor bindings generated.
* NonfungibleTokenPositionDescriptor bindings generated.
* NonfungiblePositionManager bindings generated.
* V3Migrator bindings generated.
* UniswapV3Staker bindings generated.
* QuoterV2 bindings generated.
* SwapRouter02 bindings generated.
* ProxyAdmin bindings generated.
* TransparentUpgradeableProxy bindings generated.
* WETH9 bindings generated.
✅ Done
```
