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
