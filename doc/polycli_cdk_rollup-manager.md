# `polycli cdk rollup-manager`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with CDK rollup manager contract.

## Flags

```bash
  -h, --help   help for rollup-manager
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
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

- [polycli cdk](polycli_cdk.md) - Utilities for interacting with CDK networks
- [polycli cdk rollup-manager dump](polycli_cdk_rollup-manager_dump.md) - List detailed information about the rollup manager.

- [polycli cdk rollup-manager inspect](polycli_cdk_rollup-manager_inspect.md) - List some basic information about the rollup manager.

- [polycli cdk rollup-manager list-rollup-types](polycli_cdk_rollup-manager_list-rollup-types.md) - List some basic information about each rollup type.

- [polycli cdk rollup-manager list-rollups](polycli_cdk_rollup-manager_list-rollups.md) - List some basic information about each rollup.

- [polycli cdk rollup-manager monitor](polycli_cdk_rollup-manager_monitor.md) - Watch for rollup manager events and display them on the fly.

