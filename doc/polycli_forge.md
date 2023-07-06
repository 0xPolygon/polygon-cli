# `polycli forge`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Forge dumped blocks on top of a genesis file.

```bash
polycli forge [flags]
```

## Usage

The forge tool is meant to take blocks from the `dumpblocks` command and import them on top of a different genesis. This allows for testing with faked state (depending on the consensus). Ideally we can use this to support migrating clients other chains to supernets.

First, dump blocks from an RPC endpoint.

```bash
# In this case local host is running a POA Core Archive node.
$ polycli dumpblocks http://127.0.0.1:8545 0 100000 --filename poa-core.0.to.100k --dump-receipts=false

# Even with disabling receipts, edge's eth_getBlockByNumber returns transactions.
# This needs to be done only if using json mode. Filter them out before forging:
$ cat poa-core.0.to.100k | grep '"difficulty"' > poa-core.0.to.100k.blocks
```

Second, generate `genesis.json` if one doesn't exist. Full guide [here](https://wiki.polygon.technology/docs/edge/get-started/set-up-ibft-locally), but an abridged version.

```bash
$ go install github.com/0xPolygon/polygon-edge@develop

$ polygon-edge secrets init --data-dir test-chain-1
$ NODE_ID=$(polygon-edge secrets output --node-id --data-dir test-chain-1)

# Generate the genesis.json file.
# Note: you may have to add some fields to the alloc property there may be an insufficient funds error.
$ polygon-edge genesis --ibft-validators-prefix-path test-chain- --bootnode /ip4/127.0.0.1/tcp/10001/p2p/$NODE_ID --block-gas-limit 6706541
```

Third, import the blocks on top of the genesis file.

```bash
$ polycli forge --genesis genesis.json --mode json --blocks poa-core.0.to.100k.blocks --count 99999
```

Here's how to do the same using `proto` instead of `json`.

```bash
polycli dumpblocks http://127.0.0.1:8545 0 1000000 -f poa-core.0.to.100k.proto -r=false -m proto
polycli forge --genesis genesis.json --mode proto --blocks poa-core.0.to.100k.proto --count 99999
```

Sometimes, it can be helpful to only import the blocks and transactions that are relevant. This can be done with `dumpblocks` by providing a `--filter` flag.

```bash
polycli dumpblocks http://127.0.0.1:8545/ 0 100000 \
  --filename poa-core.0.to.100k.test \
  --dump-blocks=true \
  --dump-receipts=true \
  --filter '{"to":["0xaf93ff8c6070c4880ca5abc4051f309aa19ec385","0x2d68f0161fcd778db31c7080f6c914657f4d240"],"from":["0xcf260ea317555637c55f70e55dba8d5ad8414cb0","0xaf93ff8c6070c4880ca5abc4051f309aa19ec385","0x2d68f0161fcd778db31c7080f6c914657f4d240"]}'
```

To load the pruned blocks into Edge, a couple of flags need to be set. This will import only the blocks that are listed in the blocks file. This can be non-consecutive blocks. If you receive a `not enough funds to cover gas costs` error, be sure to fund those addresses in in the `genesis.json`.

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

Start the server.

```bash
polygon-edge server --data-dir ./forged-data --chain genesis.json --grpc-address :10000 --libp2p :10001 --jsonrpc :10002
```

Query the server using the following command.

```bash
polycli rpc http://localhost:10002 eth_getBlockByNumber 2743 false | jq
```

You will notice that block numbers that have been skipped will return `null`.

## Flags

```bash
  -B, --base-block-reward string   The amount rewarded for mining blocks (default "2_000_000_000_000_000_000")
  -b, --blocks string              A file of encoded blocks; the format of this file should match the mode
  -c, --client string              Specify which blockchain client should be use to forge the data (default "edge")
      --consecutive-blocks         whether the blocks file has consecutive blocks (default true)
  -C, --count uint                 The number of blocks to try to forge (default 100)
  -d, --data-dir string            Specify a folder to be used to store the chain data (default "./forged-data")
  -g, --genesis string             Specify a file to be used for genesis configuration (default "genesis.json")
  -h, --help                       help for forge
  -m, --mode string                The forge mode indicates how we should get the transactions for our blocks [json, proto] (default "json")
  -p, --process-blocks             whether the transactions in blocks should be processed applied to the state (default true)
  -R, --read-first-block           whether to read the first block, leave false if first block is genesis
  -r, --receipts string            A file of encoded receipts; the format of this file should match the mode
      --rewrite-tx-nonces          whether to rewrite transaction nonces, set true if forging nonconsecutive blocks
  -t, --tx-fees                    if the transaction fees should be included when computing block rewards
  -V, --verifier string            Specify a consensus engine to use for forging (default "dummy")
      --verify-blocks              whether to verify blocks, set false if forging nonconsecutive blocks (default true)
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
