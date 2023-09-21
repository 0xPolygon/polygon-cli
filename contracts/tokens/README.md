# Token Contracts

## Usage

To generate go bindings for ERC20, ERC721 and more from Solidity contracts, you can simply use `make gen`. Under the hood, this handy command leverages `solc` to compile the contracts and `abigen` to generate go bindings from the contract ABI and bytecode.

```sh
$ make gen
solc ERC20/ERC20.sol --abi --bin --output-dir ERC20 --overwrite
Compiler run successful. Artifact(s) can be found in directory ERC20.
abigen --abi ERC20/ERC20.abi --bin ERC20/ERC20.bin --pkg tokens --type ERC20 --out ERC20.go
solc ERC721/ERC721.sol --abi --bin --output-dir ERC721 --overwrite
Compiler run successful. Artifact(s) can be found in directory ERC721.
abigen --abi ERC721/ERC721.abi --bin ERC721/ERC721.bin --pkg tokens --type ERC721 --out ERC721.go
```

## Example

Here is a quick example on generating go bindings for the ERC20 contract.

First, compile the contract with `solc`, the Solidity compiler. Make sure that you're version matches the one defined in the contract. To manage `solc` versions, you can use [`crytic/solc-select`](https://github.com/crytic/solc-select).

```sh
$ solc ERC20/ERC20.sol --abi --bin --output-dir ERC20 --overwrite
Compiler run successful. Artifact(s) can be found in directory ERC20.
```

Then, generate go bindings using `abigen`.

```sh
$ abigen --version
abigen version 1.13.1-stable

$ abigen --abi ERC20/ERC20.abi --bin ERC20/ERC20.bin --pkg tokens --type ERC20 --out ERC20.go
# ERC20.go has been generated.
```
