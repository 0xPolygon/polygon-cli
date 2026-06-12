Cast-familiar transaction and account queries against a Heimdall v2 node.

Read-only commands in this group talk to CometBFT's JSON-RPC (`tx`,
`receipt`, `logs`, `rpc`) or the Cosmos SDK REST gateway (`nonce`,
`balance`). Hashes may be supplied with or without a `0x` prefix;
addresses are always the 20-byte Ethereum-style hex form.

```bash
# Decode an included transaction by hash (either hex form is fine)
polycli heimdall tx 0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29

# Receipt-style view with events + logs
polycli heimdall receipt 94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29

# Wait for N confirmations past the tx's inclusion height
polycli heimdall receipt <HASH> --confirmations 3

# Full-text query over the tx index
polycli heimdall logs "message.action='/heimdallv2.topup.MsgTopupTx'" --limit 5

# Account state (nonce / sequence / balance)
polycli heimdall nonce 0x02f615e95563ef16f10354dba9e584e58d2d4314
polycli heimdall sequence 0x02f615e95563ef16f10354dba9e584e58d2d4314
polycli heimdall balance 0x02f615e95563ef16f10354dba9e584e58d2d4314 --human

# Raw JSON-RPC passthrough
polycli heimdall rpc status
polycli heimdall rpc block height=32620627

# Broadcast a pre-built TxRaw (base64 or hex). Requires --yes.
polycli heimdall publish <TXRAW> --yes
```
