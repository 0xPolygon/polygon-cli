The `loadtest` tool is meant to generate various types of load against RPC end points. It leverages the [`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) library Go Ethereum to interact with the blockchain.x

```bash
$ polycli wallet inspect  --mnemonic "code code code code code code code code code code code quality" --addresses 1
```

The `--mode` flag is important for this command.

- `t` will only perform transfers to the `--to-address`. This is a fast and common operation.
- `d` will deploy the load testing contract over and over again.
- `c` will call random functions in our load test contract.
- `f` will call a specific function on the load test contract. The function is specified using the `-f` flag
- `2` will run an ERC20 transfer test. It starts out by minting a large amount of an ERC20 contract then transferring it in small amounts.
- `7` will run an ERC721 test which will mint an NFT over and over again.
- `i` will call the increment function repeatedly on the load test contract. It's a minimal example of a contract call that will require an update to a contract's storage.
- `r` will call any of th eother modes randomly.
- `s` is used for Avail / Eth to store random data in large amounts.
- `l` will call a smart contract function that runs as long as it can, based on the block limit.

The default private key is: `42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We can use `wallet inspect` to get more information about this address, in particular its `ETHAddress` if you want to check balance or pre-mine value for this particular account.

Here is a simple example that runs 1000 requests at a max rate of 1 request per second against the http rpc endpoint on localhost. It's running in transaction mode so it will perform simple transactions send to the default address.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t http://localhost:8888
```

Another example, a bit slower, and that specifically calls the [LOG4](https://www.evm.codes/#a4) function in the load test contract in a loop for 25,078 iterations. That number was picked specifically to require almost all of the gas for a single transaction.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 50 --rate-limit 0.5  --mode f --function 164 --iterations 25078 http://private.validator-001.devnet02.pos-v3.polygon.private:8545
```

### Load Test Contract

The codebase has a contract that used for load testing. It's written in Yul and Solidity. The workflow for modifying this contract is.

1. Make changes to <file:contracts/LoadTester.sol>
2. Compile the contracts:
   - `$ solc LoadTester.sol --bin --abi -o . --overwrite`
3. Run `abigen`
   - `$ abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out loadtester.go`
4. Run the loadtester to enure it deploys and runs successfully
   - `$ polycli loadtest --verbosity 700 http://127.0.0.1:8541`

### Avail / Substrate

The loadtest tool works with Avail, but not with the same level of functionality. There's no EVM so the functional calls will not work. This is a basic example which would transfer value in a loop 10 times.

```bash
$ polycli loadtest --app-id 0 --to-random=true  --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode t 'http://devnet01.dataavailability.link:8545'
```

This is a similar test but storing random nonsense hexwords.

```bash
$ polycli loadtest --app-id 0 --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode s --byte-count 16384 'http://devnet01.dataavailability.link:8545'
```
