# `polycli cdk bridge dump`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List detailed information about the bridge

```bash
polycli cdk bridge dump [flags]
```

## Usage

This command will reach the bridge contract and retrieve detailed information.

Below is an example of how to use it

```bash
polycli cdk bridge dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
```

## Flags

```bash
  -h, --help   help for dump
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string           The address of the bridge contract
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  The ForkID of the cdk networks (default "12")
      --pretty-logs                     Should logs be in pretty format or JSON (default true)
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

- [polycli cdk bridge](polycli_cdk_bridge.md) - Utilities for interacting with CDK bridge contract
