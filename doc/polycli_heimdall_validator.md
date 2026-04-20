# `polycli heimdall validator`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query stake module endpoints.

```bash
polycli heimdall validator [ID] [flags]
```

## Usage

Validator / staking queries (`x/stake`) against a Heimdall v2 node.

Alias: `val`. `validator <ID>` is a shorthand for `validator get <ID>`.
The top-level `validators` command is an alias for `validator set`.

All subcommands hit the REST gateway.

```bash
# Full current validator set (power-desc by default)
polycli heimdall validator set
polycli heimdall validators --limit 5 --sort signer

# Aggregate voting power across the set
polycli heimdall validator total-power

# By numeric id
polycli heimdall validator 4
polycli heimdall validator get 4

# By hex signer (0x prefix optional)
polycli heimdall validator signer 0x4ad84f7014b7b44f723f284a85b1662337971439

# Membership check. Note: the upstream field `is_old` is surfaced as
# `is_current` because the upstream name is misleading — a response of
# `true` means the address is still in the current validator set.
polycli heimdall validator status 0x4ad84f7014b7b44f723f284a85b1662337971439

# Current proposer / upcoming proposers
polycli heimdall validator proposer
polycli heimdall validator proposers 5

# L1 replay check on a stake event (requires eth_rpc_url on the node)
polycli heimdall validator is-old-stake-tx 0x94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29 0
```

## Flags

```bash
  -h, --help             help for validator
      --watch duration   repeat every DURATION (e.g. 5s) until Ctrl-C; 0 disables
```

The command also inherits flags from parent commands.

```bash
      --amoy                     shortcut for --network amoy (default)
      --chain-id string          chain id used for signing
      --color string             color mode (auto|always|never) (default "auto")
      --config string            config file (default is $HOME/.polygon-cli.yaml)
      --curl                     print the equivalent curl command instead of executing
      --denom string             fee denom
      --heimdall-config string   path to heimdall config TOML (default ~/.polycli/heimdall.toml)
  -k, --insecure                 accept invalid TLS certs
      --json                     emit JSON instead of key/value
      --mainnet                  shortcut for --network mainnet
  -N, --network string           named network preset (amoy|mainnet)
      --no-color                 disable color output
      --pretty-logs              output logs in pretty format instead of JSON (default true)
      --raw                      preserve raw bytes (no 0x-hex normalization)
  -r, --rest-url string          heimdall REST gateway URL
      --rpc-headers string       extra request headers, comma-separated key=value pairs
  -R, --rpc-url string           cometBFT RPC URL
      --timeout int              HTTP timeout in seconds
  -v, --verbosity string         log level (string or int):
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

- [polycli heimdall](polycli_heimdall.md) - Query and interact with a Heimdall v2 node.
- [polycli heimdall validator get](polycli_heimdall_validator_get.md) - Fetch one validator by numeric id.

- [polycli heimdall validator is-old-stake-tx](polycli_heimdall_validator_is-old-stake-tx.md) - Check whether an L1 stake event was already replayed.

- [polycli heimdall validator proposer](polycli_heimdall_validator_proposer.md) - Show the current proposer.

- [polycli heimdall validator proposers](polycli_heimdall_validator_proposers.md) - Show the next N proposers (default 1).

- [polycli heimdall validator set](polycli_heimdall_validator_set.md) - Print the current validator set.

- [polycli heimdall validator signer](polycli_heimdall_validator_signer.md) - Fetch a validator by hex signer address.

- [polycli heimdall validator status](polycli_heimdall_validator_status.md) - Check whether an address is in the current validator set.

- [polycli heimdall validator total-power](polycli_heimdall_validator_total-power.md) - Print aggregate validator voting power.

