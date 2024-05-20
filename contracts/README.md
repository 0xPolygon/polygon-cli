# Contracts

Smart contracts used to perform different types of tests:

- `tester/` to call various opcodes, precompiles, and store random data with `LoadTester` and test revert reason string with `ConformanceTester`.
- `tokens/` to perform ERC20 transfers or ERC721 mints for example.
- Other: `asm/` contracts written in other languages than Solidity.

## Generate go bindings

```bash
$ make gen-go-bindings
FOUNDRY_PROFILE=lite forge build
[⠒] Compiling...
No files changed, compilation skipped
cat ./out/LoadTester.sol/LoadTester.json | jq -r '.abi'             > ../bindings/tester/LoadTester.abi
cat ./out/LoadTester.sol/LoadTester.json | jq -r '.bytecode.object' > ../bindings/tester/LoadTester.bin
abigen --abi ../bindings/tester/LoadTester.abi --bin ../bindings/tester/LoadTester.bin --pkg tester --type LoadTester --out ../bindings/tester/loadTester.go
✅ tester/loadTester.go generated
cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.abi'             > ../bindings/tester/ConformanceTester.abi
cat ./out/ConformanceTester.sol/ConformanceTester.json | jq -r '.bytecode.object' > ../bindings/tester/ConformanceTester.bin
abigen --abi ../bindings/tester/ConformanceTester.abi --bin ../bindings/tester/ConformanceTester.bin --pkg tester --type ConformanceTester --out ../bindings/tester/conformanceTester.go
✅ tester/conformanceTester.go generated
cat ./out/ERC20.sol/ERC20.json | jq -r '.abi'             > ../bindings/tokens/ERC20.abi
cat ./out/ERC20.sol/ERC20.json | jq -r '.bytecode.object' > ../bindings/tokens/ERC20.bin
abigen --abi ../bindings/tokens/ERC20.abi --bin ../bindings/tokens/ERC20.bin --pkg tokens --type ERC20 --out ../bindings/tokens/ERC20.go
✅ tokens/ERC20.go generated
cat ./out/ERC721.sol/ERC721.json | jq -r '.abi'             > ../bindings/tokens/ERC721.abi
cat ./out/ERC721.sol/ERC721.json | jq -r '.bytecode.object' > ../bindings/tokens/ERC721.bin
abigen --abi ../bindings/tokens/ERC721.abi --bin ../bindings/tokens/ERC721.bin --pkg tokens --type ERC721 --out ../bindings/tokens/ERC721.go
✅ tokens/ERC721.go generated
```

## Usage

```bash
$ make
Usage:
  make <target>

Help
  help             Display this help.

Build
  build            Build the smart contracts

Gen go bindings
  gen-tester-go-bindings  Generate go bindings for the tester contracts.
  gen-tokens-go-bindings  Generate go bindings for the tokens contracts.
  gen-go-bindings  Generate go bindings.
```
