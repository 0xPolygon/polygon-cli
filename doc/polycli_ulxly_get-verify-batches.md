# `polycli ulxly get-verify-batches`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate ndjson for each verify batch over a particular range of blocks

```bash
polycli ulxly get-verify-batches [flags]
```

## Usage

This command will attempt to scan a range of blocks and look for verify batch Events. This is the specific signature that we're interested
in:

```solidity
    /**
     * @dev Emitted when the trusted aggregator verifies batches
     */
    event VerifyBatchesTrustedAggregator(
        uint32 indexed rollupID,
        uint64 numBatch,
        bytes32 stateRoot,
        bytes32 exitRoot,
        address indexed aggregator
    );

```

If you're looking at the raw topics from on chain or in an explorer, this is the associated value:

`0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3`

Each event that we counter will be parsed and written as JSON to
stdout. Example usage:

```bash
polycli ulxly get-verify-batches \
        --rollup-manager-address 0x32d33d5137a7cffb54c5bf8371172bcec5f310ff \
        --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo \
        --from-block 4880876 \
        --to-block 6028159 \
        --filter-size 9999 > verify-batches-cardona-4880876-to-6028159.ndjson
```

This command will look for verify batch events from block `4880876` to
block `6028159` in increments of `9999` blocks at a time for the
contract address `0x32d33d5137a7cffb54c5bf8371172bcec5f310ff`. The
output will be written as newline delimited JSON.

This command is very specific for the ulxly bridge, and it's meant to
serve as the input to the rollup-proof command.



## Flags

```bash
  -i, --filter-size uint                The batch size for individual filter queries (default 1000)
  -f, --from-block uint                 The start of the range of blocks to retrieve
  -h, --help                            help for get-verify-batches
      --insecure                        skip TLS certificate verification
  -a, --rollup-manager-address string   The address of the rollup manager contract
  -u, --rpc-url string                  The RPC URL to read the events data
  -t, --to-block uint                   The end of the range of blocks to retrieve
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge
