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

❗️ Make sure to note the address of NFTDescriptor library, this is very important.

```sh
✍️ NFTDescriptor library address: 0xf7012159bf761b312153e8c8d176932fe9aaa7ea
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

🪄 Update NonfungibleTokenPositionDescriptor deploy function to be able to update its bytecode
```

4. Update the `NonfungibleTokenPositionDescriptor` deploy function.

This is the old function.

```go
// contracts/uniswapv3/NFTPositionDescriptor.sol#L47
func DeployNFTPositionDescriptor(auth *bind.TransactOpts, backend bind.ContractBackend, _WETH9 common.Address, _nativeCurrencyLabelBytes [32]byte) (common.Address, *types.Transaction, *NFTPositionDescriptor, error) {
	parsed, err := NFTPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NFTPositionDescriptorBin), backend, _WETH9, _nativeCurrencyLabelBytes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NFTPositionDescriptor{NFTPositionDescriptorCaller: NFTPositionDescriptorCaller{contract: contract}, NFTPositionDescriptorTransactor: NFTPositionDescriptorTransactor{contract: contract}, NFTPositionDescriptorFilterer: NFTPositionDescriptorFilterer{contract: contract}}, nil
}
```

We'd like to be able to specify the new binary so here's the new function.

```go
// contracts/uniswapv3/NFTPositionDescriptor.sol#L47
func DeployNFTPositionDescriptor(auth *bind.TransactOpts, backend bind.ContractBackend, _WETH9 common.Address, _nativeCurrencyLabelBytes [32]byte, nonfungibleTokenPositionDescriptorNewBytecode string) (common.Address, *types.Transaction, *NFTPositionDescriptor, error) {
	parsed, err := NFTPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(nonfungibleTokenPositionDescriptorNewBytecode), backend, _WETH9, _nativeCurrencyLabelBytes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NFTPositionDescriptor{NFTPositionDescriptorCaller: NFTPositionDescriptorCaller{contract: contract}, NFTPositionDescriptorTransactor: NFTPositionDescriptorTransactor{contract: contract}, NFTPositionDescriptorFilterer: NFTPositionDescriptorFilterer{contract: contract}}, nil
}
```

5. Update the `NFTDescriptor` library address inside the Uniswap v3 load test module.

```go
// cmd/loadtest/uniswapv3.go#L50
var oldNFTPositionLibraryAddress = common.HexToAddress("0x3212215ccbeb5e3a808373b805f5324cebe992af")
```
