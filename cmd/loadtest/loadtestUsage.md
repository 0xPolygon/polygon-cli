The `loadtest` tool is meant to generate various types of load against RPC end points. It leverages the [`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) library Go Ethereum to interact with the blockchain.x

```bash
$ polycli wallet inspect  --mnemonic "code code code code code code code code code code code quality" --addresses 1
```

The `--mode` flag is important for this command.

- `t`/`transaction` will perform ETH transfers. This is the simplest
  and cheapest transaction that can be performed.
- `d`/`deploy` will deploy the load testing contract over and over
  again.
- `c`/`call` will call random opcodes in our load test contract. The
  random function that is called will be repeatedly called in a loop
  based on the number of iterations from the `iterations` flag
- `f`/`function` works the same way as `call` mode but instead of
  calling random op codes, the opcode can be specified with the `-f`
  flag. If you want to call `LOG4` you would pass `-f 164` which is
  the opcode for `LOG4`.
- `2`/`erc20` will run an ERC20 transfer test. The process initializes
  by minting a large amount of tokens then transferring it in small
  amounts. Each transaction is a single transfer.
- `7`/`erc721` will run an ERC721 mint test which will mint an NFT
  over and over again. The `iterations` flag here can be used to
  control if we are mint 1 NFT or many.
- `i`/`inc`/`increment` will call the increment function repeatedly on
  the load test contract. It's a minimal example of a contract call
  that will require an update to a contract's storage.
- `s`/`store` is used to store random data in the smart contract
  storage. The amount of data stored per transaction is controlled
  with the `byte-count` flag.
- `P`/`precompiles` will randomly call the commonly implemented
  precompiled functions. This functions the same way as `call` mode
  except it's hitting precompiles rather than opcodes.
- `p`/`precompile` will call a specific precompile in a loop. This
  works the same way as `function` mode except rather than specifying
  an opcode, you're specifying a precompile. E.g to call `ECRECOVER`
  you would pass `-f 1` because it's the contract at address `0x01`
- `R`/`recall` will attempt to replay all of the transactions from the
  previous blocks. You can use `--recall-blocks` to specify how many
  previous blocks should be used to seed transaction history. It's
  expected that many of the transactions in this mode would fail.
- `r`/`random` will call any of the other modes randomly. This mode
  shouldn't be used in combination with other modes. Ideally this is a
  good way to generate a lot of random activity on a test network.
- `rpc` is a unique mode that won't just simulate transactions, it
  will simulate RPC traffic (e.g. calls to get transaction receipt or
  filter logs). This is meant to stress test RPC servers rather than
  full blockchain networks. The approach is similar to `recall` mode
  where we'll fetch some recent blocks and then use that data to
  generate a variety of calls to the RPC server.

The default private key is: `42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We can use `wallet inspect` to get more information about this address, in particular its `ETHAddress` if you want to check balance or pre-mine value for this particular account.

Here is a simple example that runs 1000 requests at a max rate of 1 request per second against the http rpc endpoint on localhost. It's running in transaction mode so it will perform simple transactions send to the default address.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t --rpc-url http://localhost:8888
```

Another example, a bit slower, and that specifically calls the [LOG4](https://www.evm.codes/#a4) function in the load test contract in a loop for 25,078 iterations. That number was picked specifically to require almost all of the gas for a single transaction.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 50 --rate-limit 0.5  --mode f --function 164 --iterations 25078 --rpc-url http://private.validator-001.devnet02.pos-v3.polygon.private:8545
```

### Load Test Contract

The codebase has a contract that used for load testing. It's written in Solidity. The workflow for modifying this contract is.

1. Make changes to <file:contracts/LoadTester.sol>
2. Compile the contracts:
   - `$ solc LoadTester.sol --bin --abi -o . --overwrite`
3. Run `abigen`
   - `$ abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out loadtester.go`
4. Run the loadtester to enure it deploys and runs successfully
   - `$ polycli loadtest --verbosity 700 --rpc-url http://127.0.0.1:8541`
