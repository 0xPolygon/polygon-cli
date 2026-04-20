Topup module queries (`x/topup`) against a Heimdall v2 node.

All subcommands hit the REST gateway. Byte fields (`account_root_hash`,
`account_proof`, `tx_hash`) are rendered as `0x…`-hex by default; pass
the global `--raw` to preserve the upstream base64.

```bash
# Merkle root of all dividend accounts
polycli heimdall topup root

# Dividend account for an address (balance / fee_amount)
polycli heimdall topup account 0x02f615e95563ef16f10354dba9e584e58d2d4314

# Merkle proof for a dividend account (requires eth_rpc_url on the
# Heimdall node — an L1-less node surfaces an L1-not-configured hint)
polycli heimdall topup proof 0x02f615e95563ef16f10354dba9e584e58d2d4314

# Verify a submitted proof (proof is hex with or without 0x prefix)
polycli heimdall topup verify 0x02f615e95563ef16f10354dba9e584e58d2d4314 0x0000…0000

# Sequence / is-old replay keys for an L1 topup tx. Both require
# eth_rpc_url on the Heimdall node.
polycli heimdall topup sequence 0x48bd44a3…5c6bf8 423
polycli heimdall topup is-old   0x48bd44a3…5c6bf8 423
```

Endpoints covered (confirmed from heimdall-v2 `proto/heimdallv2/topup/query.proto`):

- `GET /topup/dividend-account-root`
- `GET /topup/dividend-account/{address}`
- `GET /topup/account-proof/{address}`
- `GET /topup/account-proof/{address}/verify?proof=…`
- `GET /topup/sequence?tx_hash=…&log_index=…`
- `GET /topup/is-old-tx?tx_hash=…&log_index=…`
