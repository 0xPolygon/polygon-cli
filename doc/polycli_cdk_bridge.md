# `polycli cdk bridge`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with CDK bridge contract

## Flags

```bash
      --bridge-address string   The address of the bridge contract
  -h, --help                    help for bridge
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
- [polycli cdk bridge dump](polycli_cdk_bridge_dump.md) - List detailed information about the bridge

- [polycli cdk bridge inspect](polycli_cdk_bridge_inspect.md) - List some basic information about the bridge

- [polycli cdk bridge monitor](polycli_cdk_bridge_monitor.md) - Watch for bridge events and display them on the fly

