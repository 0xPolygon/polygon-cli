# Contracts

Smart contracts used to perform different types of tests:

- `LoadTester` to call various opcodes, precompiles, and store random data.
- `Tokens` to perform ERC20 transfers or ERC721 mints for example.
- `UniswapV3` to deploy the full UniswapV3 contract suite and perform some swaps.
- Other: `asm` and `yul`, contracts written in other languages than Solidity.

## LoadTester

Generate go bindings for the `LoadTester` contract.

```sh
$ FOUNDRY_PROFILE=lite forge build --contracts ./src/loadtester/LoadTester.sol \
  && cat ./out/LoadTester.sol/LoadTester.json| jq -r '.abi' > ./src/loadtester/LoadTester.abi \
  && cat ./out/LoadTester.sol/LoadTester.json| jq -r '.bytecode.object' > ./src/loadtester/LoadTester.bin \
  && abigen \
    --abi ./src/loadtester/LoadTester.abi \
    --bin ./src/loadtester/LoadTester.bin \
    --pkg loadtester \
    --type loadTester \
    --out ./src/loadtester/loadTester.go
```

## Tokens

Generate go bindings for the `ERC20` contract.

```sh
$ forge build --contracts ./src/tokens/ERC20.sol \
  && cat ./out/ERC20.sol/ERC20.json| jq -r '.abi' > ./src/tokens/ERC20.abi \
  && cat ./out/ERC20.sol/ERC20.json| jq -r '.bytecode.object' > ./src/tokens/ERC20.bin \
  && abigen \
    --abi ./src/tokens/ERC20.abi \
    --bin ./src/tokens/ERC20.bin \
    --pkg tokens \
    --type ERC20 \
    --out ./src/tokens/ERC20.go
```

Generate go bindings for the `ERC721` contract.

```sh
$ forge build --contracts ./src/tokens/ERC721.sol \
  && cat ./out/ERC721.sol/ERC721.json| jq -r '.abi' > ./src/tokens/ERC721.abi \
  && cat ./out/ERC721.sol/ERC721.json| jq -r '.bytecode.object' > ./src/tokens/ERC721.bin \
  && abigen \
    --abi ./src/tokens/ERC721.abi \
    --bin ./src/tokens/ERC721.bin \
    --pkg tokens \
    --type ERC721 \
    --out ./src/tokens/ERC721.go
```

## UniswapV3

The UniswapV3 go bindings have been generated in a certain way. Check `uniswapv3/README.org` for more details.
