# `polycli`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

A Swiss Army knife of blockchain tools.

## Usage

Polycli is a collection of tools that are meant to be useful while building, testing, and running blockchain applications.
## Flags

```bash
      --config string      config file (default is $HOME/.polygon-cli.yaml)
  -h, --help               help for polycli
      --pretty-logs        output logs in pretty format instead of JSON (default true)
  -t, --toggle             help message for toggle
  -v, --verbosity string   log level (string or int):
                             0   - silent
                             100 - panic
                             200 - fatal
                             300 - error
                             400 - warn
                             500 - info (default)
                             600 - debug
                             700 - trace (default "info")
```

## See also

- [polycli abi](polycli_abi.md) - Provides encoding and decoding functionalities with contract signatures and ABI.

- [polycli cdk](polycli_cdk.md) - Utilities for interacting with CDK networks.

- [polycli contract](polycli_contract.md) - Interact with smart contracts and fetch contract information from the blockchain.

- [polycli dbbench](polycli_dbbench.md) - Perform a level/pebble db benchmark.

- [polycli dockerlogger](polycli_dockerlogger.md) - Monitor and filter Docker container logs.

- [polycli dumpblocks](polycli_dumpblocks.md) - Export a range of blocks from a JSON-RPC endpoint.

- [polycli ecrecover](polycli_ecrecover.md) - Recovers and returns the public key of the signature.

- [polycli enr](polycli_enr.md) - Convert between ENR and Enode format.

- [polycli fix-nonce-gap](polycli_fix-nonce-gap.md) - Send txs to fix the nonce gap for a specific account.

- [polycli fold-trace](polycli_fold-trace.md) - Trace an execution trace and fold it for visualization.

- [polycli fork](polycli_fork.md) - Take a forked block and walk up the chain to do analysis.

- [polycli fund](polycli_fund.md) - Bulk fund crypto wallets automatically.

- [polycli hash](polycli_hash.md) - Provide common crypto hashing functions.

- [polycli loadtest](polycli_loadtest.md) - Run a generic load test against an Eth/EVM style JSON-RPC endpoint.

- [polycli metrics-to-dash](polycli_metrics-to-dash.md) - Create a dashboard from an Openmetrics / Prometheus response.

- [polycli mnemonic](polycli_mnemonic.md) - Generate a BIP39 mnemonic seed.

- [polycli monitor](polycli_monitor.md) - Monitor blocks using a JSON-RPC endpoint.

- [polycli monitorv2](polycli_monitorv2.md) - Monitor v2 command stub.

- [polycli nodekey](polycli_nodekey.md) - Generate node keys for different blockchain clients and protocols.

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.

- [polycli parse-batch-l2-data](polycli_parse-batch-l2-data.md) - Convert batch l2 data into an ndjson stream.

- [polycli parseethwallet](polycli_parseethwallet.md) - Extract the private key from an eth wallet.

- [polycli plot](polycli_plot.md) - Plot a chart of transaction gas prices and limits.

- [polycli publish](polycli_publish.md) - Publish transactions to the network with high-throughput.

- [polycli report](polycli_report.md) - Generate a report analyzing a range of blocks from an Ethereum-compatible blockchain.

- [polycli retest](polycli_retest.md) - Convert the standard ETH test fillers into something to be replayed against an RPC.

- [polycli rpcfuzz](polycli_rpcfuzz.md) - Continually run a variety of RPC calls and fuzzers.

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions.

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge.

- [polycli version](polycli_version.md) - Get the current version of this application.

- [polycli wallet](polycli_wallet.md) - Create or inspect BIP39(ish) wallets.

- [polycli wrap-contract](polycli_wrap-contract.md) - Wrap deployed bytecode into create bytecode.

