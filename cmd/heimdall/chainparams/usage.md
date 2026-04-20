Chainmanager module queries (`x/chainmanager`) against a Heimdall v2 node.

The chainmanager module holds the L1 / L2 chain ids, the tx confirmation
depths, and the Ethereum contract addresses Heimdall uses to interact
with the root chain. Only one REST route exists upstream
(`/chainmanager/params`); `addresses` is a derived view for quick
copy-paste into a block explorer.

```bash
# Full params (chain ids + confirmations + contract addresses)
polycli heimdall chainmanager params

# Pluck a single field
polycli heimdall chainmanager params --field params.chain_params.root_chain_address

# Just the address map, one address per line, for etherscan copy-paste
polycli heimdall chainmanager addresses

# Alias
polycli heimdall cm params
```

Endpoints covered (confirmed from heimdall-v2 `proto/heimdallv2/chainmanager/query.proto`):

- `GET /chainmanager/params`
