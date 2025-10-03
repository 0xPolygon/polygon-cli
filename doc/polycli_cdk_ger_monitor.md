# `polycli cdk ger monitor`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Watch for global exit root manager events and display them on the fly.

```bash
polycli cdk ger monitor [flags]
```

## Usage

This command will keep watching for global exit root events on chain and print them on the fly.

Below are some example of how to use it

```bash
polycli cdk ger monitor
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
```

## Flags

```bash
  -h, --help   help for monitor
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
      --ger-address string              address of GER contract
      --pretty-logs                     output logs in pretty format instead of JSON (default true)
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

- [polycli cdk ger](polycli_cdk_ger.md) - Utilities for interacting with CDK global exit root manager contract.
