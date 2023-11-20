# Contracts

Smart contracts used to perform different types of tests:

- `tester/` to call various opcodes, precompiles, and store random data with `LoadTester` and test revert reason string with `ConformanceTester`.
- `tokens/` to perform ERC20 transfers or ERC721 mints for example.
- `uniswapv3/` to deploy the full UniswapV3 contract suite and perform some swaps.
- Other: `asm/` and `yul/`, contracts written in other languages than Solidity.

## Build

Build contracts using the `lite` profile with disables the Yul optimiser, necessary to build `LoadTester.sol`.

```bash
$ FOUNDRY_PROFILE=lite forge build
```

## Tester

Generate go bindings for the `LoadTester` contract.

```sh
$ FOUNDRY_PROFILE=lite forge build --contracts ./src/tester/LoadTester.sol \
  && cat ./out/LoadTester.sol/LoadTester.json | jq -r '.abi' > ./src/tester/LoadTester.abi \
  && cat ./out/LoadTester.sol/LoadTester.json | jq -r '.bytecode.object' > ./src/tester/LoadTester.bin \
  && abigen \
    --abi ./src/tester/LoadTester.abi \
    --bin ./src/tester/LoadTester.bin \
    --pkg tester \
    --type LoadTester \
    --out ./src/tester/loadTester.go
```

Same thing for the `ConformanceTester` contract.

```sh
$ FOUNDRY_PROFILE=lite forge build --contracts ./src/tester/ConformanceTester.sol \
  && cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.abi' > ./src/tester/ConformanceTester.abi \
  && cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.bytecode.object' > ./src/tester/ConformanceTester.bin \
  && abigen \
    --abi ./src/tester/ConformanceTester.abi \
    --bin ./src/tester/ConformanceTester.bin \
    --pkg tester \
    --type ConformanceTester \
    --out ./src/tester/conformanceTester.go
```

## Tokens

Generate go bindings for the `ERC20` contract.

```sh
$ forge build --contracts ./src/tokens/ERC20.sol \
  && cat ./out/ERC20.sol/ERC20.json | jq -r '.abi' > ./src/tokens/ERC20.abi \
  && cat ./out/ERC20.sol/ERC20.json | jq -r '.bytecode.object' > ./src/tokens/ERC20.bin \
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
