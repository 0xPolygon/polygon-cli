# `polycli ulxly get-deposits`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Get a range of deposits

```bash
polycli ulxly get-deposits [flags]
```

## Usage

This command will attempt to scan a range of blocks and look for uLxLy
Bridge Events. This is is the specific signature that we're interested
in:

```solidity
    /**
     * @dev Emitted when bridge assets or messages to another network
     */
    event BridgeEvent(
        uint8 leafType,
        uint32 originNetwork,
        address originAddress,
        uint32 destinationNetwork,
        address destinationAddress,
        uint256 amount,
        bytes metadata,
        uint32 depositCount
    );

```

If you're looking at the raw topics from on chain or in an explorer, this is the associated value:

`0x501781209a1f8899323b96b4ef08b168df93e0a90c673d1e4cce39366cb62f9b`

Each event that we counter will be parsed and written as JSON to
stdout. Example usage:

```bash
polycli ulxly deposits \
        --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
        --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo \
        --from-block 4880876 \
        --to-block 6028159 \
        --filter-size 9999 > cardona-4880876-to-6028159.ndjson
```

This command would look for bridge events from block `4880876` to
block `6028159` in increments of `9999` blocks at a time for the
contract address `0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582`. The
output will be written as newline delimited JSON.

This command is very specific for the ulxly bridge and it's meant to
serve as the input to the proof command.



## Flags

```bash
      --bridge-address string   The address of the lxly bridge
      --filter-size uint        The batch size for individual filter queries (default 1000)
      --from-block uint         The block height to start query at.
  -h, --help                    help for get-deposits
      --rpc-url string          The RPC to query for events (default "http://127.0.0.1:8545")
      --to-block uint           The block height to start query at.
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the lxly bridge
