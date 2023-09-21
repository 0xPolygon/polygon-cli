# Tokens

How to generate go bindings for ERC20, ERC721 and more from Solidity contracts.

## ERC20

1. Compile contract.

```sh
$ solc ERC20/ERC20.sol --abi --bin --output-dir ERC20 --overwrite
Compiler run successful. Artifact(s) can be found in directory ERC20.
```

2. Generate go bindings.

```sh
$ abigen --abi ERC20/ERC20.abi --bin ERC20/ERC20.bin --pkg tokens --type ERC20 --out ERC20.go
# ERC20.go generated
```
