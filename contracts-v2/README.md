# Generate Go Bindings

Example for `src/tokens/MyToken.sol`

```sh
$ forge build
$ cat ./out/ERC20.sol/ERC20.json | jq -r '.abi'             > ./src/tokens/ERC20/ERC20.abi
$ cat ./out/ERC20.sol/ERC20.json | jq -r '.bytecode.object' > ./src/tokens/ERC20/ERC20.bin
$ abigen \
  --abi ./src/tokens/ERC20/ERC20.abi \
  --bin ./src/tokens/ERC20/ERC20.bin \
  --pkg contracts \
  --out ./src/tokens/ERC20/ERC20.go
```
