# Contracts

## Usage

To generate go bindings for ERC20, ERC721 and more from Solidity contracts, you can simply use `make gen-go-bindings`. Under the hood, this handy command leverages `solc` to compile the contracts and `abigen` to generate go bindings from the contract ABI and bytecode.

```sh
~/polycli/contracts $ make gen-go-bindings
solc tokens/ERC20/ERC20.sol --abi --bin --output-dir tokens/ERC20 --overwrite
Compiler run successful. Artifact(s) can be found in directory "tokens/ERC20".
abigen --abi tokens/ERC20/ERC20.abi --bin tokens/ERC20/ERC20.bin --pkg tokens --type ERC20 --out tokens/ERC20.go
...
```

## Example

Here is a quick example on generating go bindings for the ERC20 contract.

First, compile the contract with `solc`, the Solidity compiler. Make sure that you're version matches the one defined in the contract. To manage `solc` versions, you can use [`crytic/solc-select`](https://github.com/crytic/solc-select).

```sh
$ solc tokens/ERC20/ERC20.sol --abi --bin --output-dir tokens/ERC20 --overwrite
Compiler run successful. Artifact(s) can be found in directory ERC20.
```

Then, generate go bindings using `abigen`.

```sh
$ abigen --version
abigen version 1.13.1-stable

$ abigen --abi tokens/ERC20/ERC20.abi --bin tokens/ERC20/ERC20.bin --pkg tokens --type ERC20 --out tokens/ERC20.go
# ERC20.go has been generated.
```
