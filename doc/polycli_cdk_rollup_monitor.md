# `polycli cdk rollup monitor`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Watch for rollup events and display them on the fly

```bash
polycli cdk rollup monitor [flags]
```

## Usage

This command will keep watching for rollup manager events from a specific rollup on chain and print them on the fly.

Below are some example of how to use it

```bash
polycli cdk rollup monitor
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-id 1

polycli cdk rollup monitor
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-chain-id 2440

polycli cdk rollup monitor
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-address 0x89ba0ed947a88fe43c22ae305c0713ec8a7eb361
```

## Flags

```bash
  -h, --help   help for monitor
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  The ForkID of the cdk networks (default "12")
      --pretty-logs                     Should logs be in pretty format or JSON (default true)
      --rollup-address string           The rollup Address
      --rollup-chain-id string          The rollup chain ID
      --rollup-id string                The rollup ID
      --rollup-manager-address string   The address of the rollup contract
      --rpc-url string                  The RPC URL of the network containing the CDK contracts (default "http://localhost:8545")
  -v, --verbosity int                   0 - Silent
                                        100 Panic
                                        200 Fatal
                                        300 Error
                                        400 Warning
                                        500 Info
                                        600 Debug
                                        700 Trace (default 500)
```

## See also

- [polycli cdk rollup](polycli_cdk_rollup.md) - Utilities for interacting with CDK rollup manager to get rollup specific information
