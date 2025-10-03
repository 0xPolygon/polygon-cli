# `polycli cdk rollup dump`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List detailed information about a specific rollup

```bash
polycli cdk rollup dump [flags]
```

## Usage

This command will reach the rollup manager contract and retrieve detailed information from a specific rollup.

Below is an example of how to use it

```bash
polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-id 1

polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-chain-id 2440

polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-address 0x89ba0ed947a88fe43c22ae305c0713ec8a7eb361
```

## Flags

```bash
  -h, --help   help for dump
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
      --pretty-logs                     output logs in pretty format instead of JSON (default true)
      --rollup-address string           rollup address
      --rollup-chain-id string          rollup chain ID
      --rollup-id string                rollup ID
      --rollup-manager-address string   address of rollup contract
      --rpc-url string                  RPC URL of network containing CDK contracts (default "http://localhost:8545")
  -v, --verbosity string                log level (string or int):
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

- [polycli cdk rollup](polycli_cdk_rollup.md) - Utilities for interacting with CDK rollup manager to get rollup specific information
