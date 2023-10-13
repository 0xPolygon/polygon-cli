# `polycli loadtest`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Run a generic load test against an Eth/EVM style JSON-RPC endpoint.

```bash
polycli loadtest url [flags]
```

## Usage

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

```

## Flags

```bash
      --adaptive-backoff-factor float              When using adaptive rate limiting, this flag controls our multiplicative decrease value. (default 2)
      --adaptive-cycle-duration-seconds uint       When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates (default 10)
      --adaptive-rate-limit                        Enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint         When using adaptive rate limiting, this flag controls the size of the additive increases. (default 50)
      --bar-chart-num-bucket int                   The number of buckets to visualize latency time distribution (default 10)
      --batch-size uint                            Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time. (default 999)
  -b, --byte-count uint                            If we're in store mode, this controls how many bytes we'll try to store in our contract (default 1024)
      --call-only                                  When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.
      --call-only-latest                           When using call only mode with recall, should we execute on the latest block or on the original block
      --chain-id uint                              The chain id for the transactions.
  -c, --concurrency int                            Number of requests to perform concurrently. Default is one request at a time. (default 1)
      --contract-call-block-interval uint          During deployment, this flag controls if we should check every block, every other block, or every nth block to determine that the contract has been deployed (default 1)
      --contract-call-nb-blocks-to-wait-for uint   The number of blocks to wait for before giving up on a contract deployment (default 30)
      --erc20-address string                       The address of a pre-deployed erc 20 contract
      --erc721-address string                      The address of a pre-deployed erc 721 contract
      --force-contract-deploy                      Some load test modes don't require a contract deployment. Set this flag to true to force contract deployments. This will still respect the --lt-address flags.
  -f, --function --mode f                          A specific function to be called if running with --mode f or a specific precompiled contract when running with `--mode a` (default 1)
      --gas-limit uint                             In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas
      --gas-price uint                             In environments where the gas price can't be determined automatically, we can specify it manually
  -h, --help                                       help for loadtest
  -i, --iterations uint                            If we're making contract calls, this controls how many times the contract will execute the instruction in a loop. If we are making ERC721 Mints, this indicates the minting batch size (default 1)
      --legacy                                     Send a legacy transaction instead of an EIP1559 transaction.
      --lt-address string                          The address of a pre-deployed load test contract
  -m, --mode strings                               The testing mode to use. It can be multiple like: "t,c,d,f"
                                                   t - sending transactions
                                                   d - deploy contract
                                                   c - call random contract functions
                                                   f - call specific contract function
                                                   p - call random precompiled contracts
                                                   a - call a specific precompiled contract address
                                                   s - store mode
                                                   r - random modes
                                                   2 - ERC20 Transfers
                                                   7 - ERC721 Mints
                                                   R - total recall
                                                   rpc - call random rpc methods (default [t])
      --output-mode string                         Format mode for summary output (json | text) (default "text")
      --priority-gas-price uint                    Specify Gas Tip Price in the case of EIP-1559
      --private-key string                         The hex encoded private key that we'll use to send transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --rate-limit float                           An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together (default 4)
      --recall-blocks uint                         The number of blocks that we'll attempt to fetch for recall (default 50)
  -n, --requests int                               Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results. (default 1)
      --seed int                                   A seed for generating random values and addresses (default 123456)
      --send-amount string                         The amount of wei that we'll send every transaction (default "0x38D7EA4C68000")
      --steady-state-tx-pool-size uint             When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off. (default 1000)
      --summarize                                  Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time
  -t, --time-limit int                             Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit. (default -1)
      --to-address string                          The address that we're going to send to (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --to-random                                  When doing a transfer test, should we send to random addresses rather than DEADBEEFx5
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
