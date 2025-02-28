# `polycli cdk rollup`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with CDK rollup manager to get rollup specific information

## Flags

```bash
  -h, --help                     help for rollup
      --rollup-address string    The rollup Address
      --rollup-chain-id string   The rollup chain ID
      --rollup-id string         The rollup ID
```

The command also inherits flags from parent commands.

```bash
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

- [polycli cdk](polycli_cdk.md) - Utilities for interacting with CDK networks
- [polycli cdk rollup dump](polycli_cdk_rollup_dump.md) - List detailed information about a specific rollup

- [polycli cdk rollup inspect](polycli_cdk_rollup_inspect.md) - List some basic information about a specific rollup

- [polycli cdk rollup monitor](polycli_cdk_rollup_monitor.md) - Watch for rollup events and display them on the fly

