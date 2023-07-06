We run a lot of different blockchain technologies. Different tools often
have inconsistent tooling and this makes automation and operations
painful. The goal of this codebase is to standardize some of our
commonly needed tools and provide interfaces and formats.

# Summary

- [Install](#install)
- [Features](#features)
  - [ABI](#abi): parse an ABI and print the encoded signatures.
  - [Dumpblocks](#dumpblocks): export a range of blocks from an RPC endpoint.
  - [Forge](#forge): take blocks from the `dumpblocks` command and import them on top of a different genesis.
  - [Fork](#fork): create a fork of a chain given a blockhash.
  - [Hash](#hash): perform hashes on files, standard input, and arguments.
  - [Loadtest](#loadtest): generate various types of load against RPC endpoints.
  - [Metrics To Dash](#metrics-to-dash): given an Openmetrics / Prometheus response, create a JSON file that can be used to create a dashboard with all of the metrics in one view.
  - [Mnemonic](#mnemonic): generate [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) mnemonics.
  - [Monitor](#monitor): monitor block production on local RPC endpoints.
  - [Nodekey](#nodekey): generating a node key (e.g. for Avail or Edge).
  - [P2P](#p2p): commands related to `devp2p` (ping nodes, discover peers and crawl a network).
  - [Parse ETH Wallet](#parse-eth-wallet): extract the private key from an Ethereum wallet.
  - [RPC](#rpc): provide aliases for making RPC calls (get rid of JSON).
  - [Wallet](#wallet): generate portable wallets to be used across ETH, BTC, Polygon, Avail, etc.
- [Testing with Geth](#testing-with-geth)
- [Reference](#reference)

# Install

Requirements:

- [Go](https://go.dev/)
- make
- jq
- bc
- protoc (Only required for `make generate`)

To install, clone this repo and run:

```bash
$ make install
```

By default, the commands will be in `$HOME/go/bin/`, so for ease, we
recommend adding that path to your shell's startup file by adding the
following line:

```bash
export PATH="$HOME/go/bin:$PATH"
```

# Features

Here is a list of the different features offered by `polycli`. You can also use `polycli -h`.

## ABI

This command is useful for analyzing Solidity ABIs and decoding function selectors and input data. Most commonly, we need
this capability while analyzing raw blocks when we don't know which method is being called, but we know the smart contract.
We can the command like this to get the function signatures and selectors:

```shell
go run main.go abi --file contract.abi
```

This would output some information that would let us know the various function selectors for this contract:

```txt
Selector:19d8ac61	Signature:function lastTimestamp() view returns(uint64)
Selector:a066215c	Signature:function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
Selector:715018a6	Signature:function renounceOwnership() returns()
Selector:cfa8ed47	Signature:function trustedSequencer() view returns(address)
```

If we want to break down input data we can run something like this:

```shell
go run main.go abi --data 0x3c158267000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063ed0f8f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006eec03843b9aca0082520894d2f852ec7b4e457f6e7ff8b243c49ff5692926ea87038d7ea4c68000808204c58080642dfe2cca094f2419aad1322ec68e3b37974bd9c918e0686b9bbf02b8bd1145622a3dd64202da71549c010494fd1475d3bf232aa9028204a872fd2e531abfd31c000000000000000000000000000000000000 < contract.abi
```

In addition to the function selector data, we'll also get a breakdown of input data:

```json
{
  "batches": [
    {
      "transactions": "7AOEO5rKAIJSCJTS+FLse05Ff25/+LJDxJ/1aSkm6ocDjX6kxoAAgIIExYCAZC3+LMoJTyQZqtEyLsaOOzeXS9nJGOBoa5u/Ari9EUViKj3WQgLacVScAQSU/RR1078jKqkCggSocv0uUxq/0xw=",
      "globalExitRoot": [
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0, 0, 0, 0, 0
      ],
      "timestamp": 1676480399,
      "minForcedTimestamp": 0
    }
  ]
}
```

## Dumpblocks

For various reasons, we might want to dump a large range of blocks for
analytics or replay purposes. This is a simple util to export over RPC a
range of blocks.

This would download the first 500K blocks and zip them and then look for
blocks with transactions that create an account.

```shell
polycli dumpblocks http://172.26.26.12:8545/ 0 500000 | gzip > foo.gz

zcat < foo.gz | jq '. | select(.transactions | length > 0) | select(.transactions[].to == null)'
```

### Protobuf

Dumpblocks can also output to protobuf format. If you wish to make changes to the protobuf:

1. Install the protobuf compiler

   On GNU/Linux:

   ```bash
   sudo apt install protoc-gen-go
   ```

   On a MAC:

   ```bash
   brew install protobuf
   ```

2. Install the protobuf plugin

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```

3. Compile the proto file

   ```bash
   make generate
   ```

4. Depending on what endpoint and chain you're querying, you may be missing some fields in the proto definition.

   ```
   {"level":"error","error":"proto: (line 1:813): unknown field \"nonce\"","time":"2023-01-17T13:35:53-05:00","message":"failed to unmarshal proto message"}
   ```

   To solve this, add the unknown fields to the `.proto` files and recompile them (step 3).

## Forge

The forge tool is meant to take blocks from the `dumpblocks` command and
import them on top of a different genesis. This allows for testing with
faked state (depending on the consensus). Ideally we can use this to
support migrating clients other chains to supernets.

- Generate `genesis.json` if one doesn't exist. Full guide [here](https://wiki.polygon.technology/docs/category/deploy-a-supernet),
  but an abridged version:

  ```bash
  # Install the `polygon-edge` binary.
  $ go install github.com/0xPolygon/polygon-edge@v1.0.0

  ## New way
  # Generate new PolyBFT secrets.
  $ polygon-edge polybft-secrets --data-dir test-chain-1 --insecure

  # Generate the `genesis.json` file
  $ polygon-edge genesis \
    --block-gas-limit 10000000 \
    --epoch-size 10 \
    --validators-path ./ \
    --validators-prefix test-chain- \
    --consensus polybft \
    --premine 0x0 \
    --reward-wallet 0x61324166B0202DB1E7502924326262274Fa4358F:1000000 \
    --transactions-allow-list-admin 0x61324166B0202DB1E7502924326262274Fa4358F,0xFE5E166BA5EA50c04fCa00b07b59966E6C2E9570 \
    --transactions-allow-list-enabled 0x61324166B0202DB1E7502924326262274Fa4358F,0xFE5E166BA5EA50c04fCa00b07b59966E6C2E9570

  ## 0ld way
  # Generate new IBFT secrets.
  $ polygon-edge secrets init --data-dir test-chain-1 --insecure
  $ NODE_ID=$(polygon-edge secrets output --node-id --data-dir test-chain-1)

  # Generate the `genesis.json` file.
  # Note: you may have to add some fields to the alloc property there may be an insufficient funds error.
  $ polygon-edge genesis \
      --consensus ibft \
      --block-gas-limit 6706541 \
      --ibft-validators-prefix-path test-chain- \
      --bootnode /ip4/127.0.0.1/tcp/10001/p2p/$NODE_ID
  ```

```shell
# In this case local host is running a POA Core Archive node.
polycli dumpblocks http://127.0.0.1:8545 0 100000 --filename poa-core.0.to.100k --dump-receipts=false

# Even with disabling receipts, edge's eth_getBlockByNumber returns transactions.
# This needs to be done only if using json mode. Filter them out before forging:
cat poa-core.0.to.100k | grep '"difficulty"' > poa-core.0.to.100k.blocks

polycli forge --genesis genesis.json --mode json --blocks poa-core.0.to.100k.blocks --count 99999
```

```bash
# To do the same with using proto instead of json:
polycli dumpblocks http://127.0.0.1:8545 0 1000000 -f poa-core.0.to.100k.proto -r=false -m proto
polycli forge --genesis genesis.json --mode proto --blocks poa-core.0.to.100k.proto --count 99999
```

### Forging Filtered Blocks

Sometimes, it can be helpful to only import the blocks and transactions that are
relevant. This can be done with `dumpblocks` by providing a `--filter` flag.

```bash
polycli dumpblocks http://127.0.0.1:8545/ 0 100000 \
  --filename poa-core.0.to.100k.test \
  --dump-blocks=true \
  --dump-receipts=true \
  --filter '{"to":["0xaf93ff8c6070c4880ca5abc4051f309aa19ec385","0x2d68f0161fcd778db31c7080f6c914657f4d240"],"from":["0xcf260ea317555637c55f70e55dba8d5ad8414cb0","0xaf93ff8c6070c4880ca5abc4051f309aa19ec385","0x2d68f0161fcd778db31c7080f6c914657f4d240"]}'
```

To load the pruned blocks into Edge, a couple of flags need to be set. This will
import only the blocks that are listed in the blocks file. This can be
non-consecutive blocks. If you receive a `not enough funds to cover gas costs`
error, be sure to fund those addresses in in the `genesis.json`.

```bash
polycli forge \
  --genesis genesis.json \
  --mode json \
  --blocks poa-core.0.to.100k.test.blocks \
  --receipts poa-core.0.to.100k.test.receipts \
  --count 2 \
  --tx-fees=true \
  --base-block-reward 1000000000000000000 \
  --read-first-block=true \
  --rewrite-tx-nonces=true \
  --verify-blocks=false \
  --consecutive-blocks=false \
  --process-blocks=false
```

Start the server with:

```bash
polygon-edge server --data-dir ./forged-data --chain genesis.json --grpc-address :10000 --libp2p :10001 --jsonrpc :10002
```

and query it with:

```bash
polycli rpc http://localhost:10002 eth_getBlockByNumber 2743 false | jq
```

You will notice that block numbers that have been skipped will return `null`.

## Fork

Occasionally, we'll want to analyze the details of a side chain to understand in detail who was proposing the blocks, what was the difficultly, and just generally get better understanding of block propagation.

```shell
go run main.go fork 0x053d84d5215684c8ae810a4729f7c9b54d65a80b128a27aeddcd7dc295a0cebd https://polygon-rpc.com
```

In order to use this, you'll need to have a blockhash of a block that was part of a fork / side chain. Once you have that, you can run `fork` against a node to get the details of the fork and the canonical chain.

## Hash

The `hash` command provides a simple mechanism to perform hashes on
files, standard input, and arguments. Below shows various ways to
provide input.

```bash
$ echo -n "hello" > hello.txt
$ polycli hash sha1 --file hello.txt
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ echo -n "hello" | polycli hash sha1
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ polycli hash sha1 hello
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

We've provided many standard hashing functions

```shell
echo -n "hello" | polycli hash md4
echo -n "hello" | polycli hash md5
echo -n "hello" | polycli hash sha1
echo -n "hello" | polycli hash sha224
echo -n "hello" | polycli hash sha256
echo -n "hello" | polycli hash sha384
echo -n "hello" | polycli hash sha512
echo -n "hello" | polycli hash ripemd160
echo -n "hello" | polycli hash sha3_224
echo -n "hello" | polycli hash sha3_256
echo -n "hello" | polycli hash sha3_384
echo -n "hello" | polycli hash sha3_512
echo -n "hello" | polycli hash sha512_224
echo -n "hello" | polycli hash sha512_256
echo -n "hello" | polycli hash blake2s_256
echo -n "hello" | polycli hash blake2b_256
echo -n "hello" | polycli hash blake2b_384
echo -n "hello" | polycli hash blake2b_512
echo -n "hello" | polycli hash keccak256
echo -n "hello" | polycli hash keccak512
```

## Loadtest

The `loadtest` tool is meant to generate various types of load against
RPC end points. It leverages the
[`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
library Go Ethereum to interact with the blockchain.

```shell
polycli loadtest --help
```

Most of the options are expected in the help text. The default private
key is:
`42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We
can use `wallet inspect` to get more information about this address, in
particular its `ETHAddress` if you want to check balance or pre-mine
value for this particular account.

```shell
polycli wallet inspect  --mnemonic "code code code code code code code code code code code quality" --addresses 1
```

The `--mode` flag is important for this command:

- `t` We'll only performance transfers to the `--to-address`. This is a
  fast and common operation.
- `d` Will deploy the load testing contract over and over again
- `c` Will call random functions in our load test contract
- `f` will call a specific function on the load test contract. The
  function is specified using the `-f` flag
- `2` Will run an ERC20 transfer test. It starts out by minting
  a large amount of an ERC20 contract then transferring it in small
  amounts
- `7` Will run an ERC721 test which will mint an NFT over and over
  again
- `i` Will call the increment function repeatedly on the load test
  contract. It's a minimal example of a contract call that will
  require an update to a contract's storage.
- `r` Will call any of th eother modes randomly
- `s` Is used for Avail / Eth to store random data in large amounts
- `l` Will call a smart contract function that runs as long as it can
  (based on the block limit)

This example is very simple. It runs 1000 requests at a max rate of 1
request per second against the http rpc endpoint on localhost. t's
running in transaction mode so it will perform simple transactions send
to the default address.

```shell
polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t http://localhost:8888
```

This example runs slower and specifically calls the
[LOG4](https://www.evm.codes/#a4) function in the load test contract in
a loop for 25,078 iterations. That number was picked specifically to
require almost all of the gas for a single transaction.

```shell
polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 50 --rate-limit 0.5  --mode f --function 164 --iterations 25078 http://private.validator-001.devnet02.pos-v3.polygon.private:8545
```

### Load Test Contract

The codebase has a contract that used for load testing. It's written in
Yul and Solidity. The workflow for modifying this contract is.

1.  Make changes to <file:contracts/LoadTester.sol>
2.  Compile the contracts:
    - `solc LoadTester.sol --bin --abi -o . --overwrite`
3.  Run `abigen`
    - `abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out loadtester.go`
4.  Run the loadtester to enure it deploys and runs successfully
    - `go run main.go loadtest --verbosity 700 http://127.0.0.1:8541`

### Avail / Substrate

The loadtest tool works with Avail, but not with the same level of
functionality. There's no EVM so the functional calls will not work.
This is a basic example which would transfer value in a loop 10 times

```shell
polycli loadtest --app-id 0 --to-random=true  --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode t 'http://devnet01.dataavailability.link:8545'
```

This is a similar test but storing random nonsense hexwords

```shell
polycli loadtest --app-id 0 --data-avail --verbosity 700 --chain-id 42 --concurrency 1 --requests 10 --rate-limit 1 --mode s --byte-count 16384 'http://devnet01.dataavailability.link:8545'
```

## Metrics To Dash

Given an openmetrics / prometheus response, create a json file that can
be used to create a dashboard with all of the metrics in one view

```shell
go run main.go metrics-to-dash -i avail-metrics.txt -p avail. -t "Avail Devnet Dashboard" -T basedn -D devnet01.avail.polygon.private -T host -D validator-001 -s substrate_ -s sub_ -P true -S true
go run main.go metrics-to-dash -i avail-light-metrics.txt -p avail_light. -t "Avail Light Devnet Dashboard" -T basedn -D devnet01.avail.polygon.private -T host -D validator-001 -s substrate_ -s sub_ -P true -S true

```

## Mnemonic

The `mnemonic` command is a simple way to generate
[BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
mnemonics.

```shell
polycli mnemonic
polycli mnemonic --language spanish
polycli mnemonic --language spanish --words 12
```

## Monitor

![](assets/polycli-monitor.gif)

This is a basic tool for monitoring block production on a local RPC end
point.

If you're using the terminal UI and you'd like to be able to select text
for copying, you might need to use a modifier key.

```sh
polycli monitor --help

polycli monitor https://polygon-rpc.com
```

If you're experiencing missing blocks, try adjusting the `--batch-size` and
`--interval` flags so that you poll for more blocks or more frequently.

## Nodekey

The `nodekey` command is still in progress, but the idea is to have a
simple command for generating a node key. Most clients will generate
this on the fly, but if we want to store the key pair during an
automated provisioning process, it's helpful to have the output be
structured

```shell
# this will generate a secp256k1 key for devp2p protocol
polycli nodekey

# generate a networking keypair for avail
polycli nodekey --protocol libp2p

# generate a networking keypair for edge
polycli nodekey --protocol libp2p --key-type secp256k1 --marshal-protobuf
```

## P2P

Pinging a peer is useful to determine information about the peer and retrieving
the `Hello` and `Status` messages. By default, it will listen to the peer after
the status exchange for blocks and transactions. To disable this behavior, set
the `--listen` flag.

```bash
polycli p2p ping <enode/enr or nodes.json file>
```

Running the sensor will do peer discovery and continue to watch for blocks and
transactions from those peers. This is useful for observing the network for forks
and reorgs without the need to run the entire full node infrastructure.

```bash
go run main.go p2p sensor nodes.json --bootnodes enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303 --network-id 137 --sensor-id "sensor" --project-id "devtools-sandbox"
```

To crawl the network for nodes and write the output json to a file. This will
not engage in block or transaction propagation, but it can give a good indicator
of network size, and the output json can be used to quick start other nodes.

```bash
go run main.go p2p crawl nodes.json --bootnodes enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303 --network-id 137
```

## Parse ETH Wallet

This function can take a geth style wallet file and extract the private key as hex. It can also do the opposite.

This command would take the private key and import it into a local keystore with no password.

```bash
$ polycli parseethwallet --hexkey 42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa

$ cat UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "address": "85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "crypto": {
    "cipher": "aes-128-ctr",
    "ciphertext": "d0b4377a4ae5ebc9a5bef06ce4be99565d10cb0dedc2f7ff5aaa07ea68e7b597",
    "cipherparams": {
      "iv": "8ecd172ff7ace15ed5bc44ea89473d8e"
    },
    "kdf": "scrypt",
    "kdfparams": {
      "dklen": 32,
      "n": 262144,
      "p": 1,
      "r": 8,
      "salt": "cd6ec772dc43225297412809feaae441d578642c6a67cabf4e29bcaf594f575b"
    },
    "mac": "c992128ed466ad15a9648f4112af22929b95f511f065b12a80abcfb7e4d39a79"
  },
  "id": "82af329d-2af5-41a6-ae6b-624f3e1c224b",
  "version": 3
}
```

If we wanted to go the opposite direction, we could run a command like this.

```bash
$ polycli parseethwallet --file /tmp/keystore/UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "Address": "0x85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "PublicKey": "507cf9a75e053cda6922467721ddb10412da9bec30620347d9529cc77fca24334a4cf59685be4a2fdeabf4e7753350e42d2d3a20250fd9dc554d226463c8a3d5",
  "PrivateKey": "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
}
```

## RPC

This is a simple tool to avoid typing JSON on the command line while
making RPC calls. The implementation is generic and this is meant to be
a complete generic RPC tool.

```shell

polycli rpc https://polygon-rpc.com eth_blockNumber

polycli rpc https://polygon-rpc.com eth_getBlockByNumber 0x1e99576 true
```

## Wallet

The `wallet` command can generate portable wallets to be used across
ETH, BTC, Polygon, Avail, etc.

In the example, we're generating a wallet with a few flags that are used
to configure how many wallets are generated and how the seed phrase is
used to generate the wallets.

```shell
polycli wallet create --raw-entropy --root-only --words 15 --language english
```

In addition to generating wallets with new mnemonics, you can use a
known mnemonic to generate wallets. **Caution** entering your seed
phrase in the command line should only be done for test mnemonics. Never
do this with a real seed phrase. The example below is a test vector from
Substrate
[BIP-0039](https://github.com/paritytech/substrate-bip39/blob/eef2f86337d2dab075806c12948e8a098aa59d59/src/lib.rs#L74)
where the expected seed is
`44e9d125f037ac1d51f0a7d3649689d422c2af8b1ec8e00d71db4d7bf6d127e33f50c3d5c84fa3e5399c72d6cbbbbc4a49bf76f76d952f479d74655a2ef2d453`

```shell
polycli wallet inspect --raw-entropy --root-only --language english --password "Substrate" --mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
```

This command also leverages the BIP-0032 library for hierarchically
derived wallets.

```shell
polycli wallet create --path "m/44'/0'/0'" --addresses 5
```

# Testing with Geth

While working on some of the Polygon CLI tools, we'll run geth in dev
mode in order to make sure the various functions work properly. First,
we'll startup geth.

```shell
# Geth
./build/bin/geth --dev --dev.period 2 --http --http.addr localhost --http.port 8545 --http.api admin,debug,web3,eth,txpool,personal,miner,net --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 --dev.gaslimit 50000000
```

In the logs, we'll see a line that says IPC endpoint opened:

```example
INFO [08-14|16:09:31.451] Starting peer-to-peer node               instance=Geth/v1.10.21-stable-67109427/darwin-arm64/go1.18.1
WARN [08-14|16:09:31.451] P2P server will be useless, neither dialing nor listening
DEBUG[08-14|16:09:31.452] IPCs registered                          namespaces=admin,debug,web3,eth,txpool,personal,clique,miner,net,engine
INFO [08-14|16:09:31.452] IPC endpoint opened                      url=/var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
INFO [08-14|16:09:31.452] Generated ephemeral JWT secret           secret=0xdfa5c30e07ef1041d15a2dbf0865386305330128b792d4a461cddb9bf38e416e
```

I'll usually then use that line to attach

```shell
./build/bin/geth attach /var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
```

After attaching to geth, we can fund the default load testing account
with some currency.

```shell
eth.coinbase==eth.accounts[0]
eth.sendTransaction({from: eth.coinbase, to: "0x85da99c8a7c2c95964c8efd687e95e632fc533d6", value: web3.toWei(5000, "ether")})
```

Then we can generate some load to make sure that there are some blocks
with transactions being created. `1337` is the chain id that's used in
local geth.

```shell
polycli loadtest --verbosity 700 --chain-id 1337 --concurrency 1 --requests 1000 --rate-limit 5 --mode c http://127.0.0.1:8545
```

Then we can monitor the chain:

```bash
polycli monitor http://127.0.0.1:8545
```

# Reference

Sending some value to the default load testing account

Listening for re-orgs

```shell
socat - UNIX-CONNECT:/var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
{"id": 1, "method": "eth_subscribe", "params": ["newHeads"]}
```

Useful RPCs when testing

```shell
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "net_version", "params": []}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_blockNumber", "params": []}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getBlockByNumber", "params": ["0x1DE8531", true]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "clique_getSigner", "params": ["0x1DE8531", true]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getBalance", "params": ["0x85da99c8a7c2c95964c8efd687e95e632fc533d6", "latest"]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getCode", "params": ["0x79954f948079ee9ef1d15eff3e07ceaef7cdf3b4", "latest"]}' https://polygon-rpc.com


curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "txpool_inspect", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "txpool_status", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_gasPrice", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "admin_peers", "params": []}' http://localhost:8545
```

```shell
websocat ws://34.208.176.205:9944
{"jsonrpc":"2.0", "id": 1, "method": "chain_subscribeNewHead", "params": []}
```
