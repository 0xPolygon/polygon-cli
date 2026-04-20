# `polycli heimdall send`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Build, sign, and broadcast a transaction.

```bash
polycli heimdall send <MSG> [flags]
```

## Usage

Build a Heimdall v2 transaction for the chosen message type, sign
it, and POST it to the REST gateway. The default mode is
BROADCAST_MODE_SYNC: polycli waits for CheckTx to return, prints
the tx hash, and then polls CometBFT for inclusion. --async skips
both waits. --confirmations N waits for N blocks past inclusion.
--dry-run stops after building (useful for CI).
## Flags

```bash
  -h, --help   help for send
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
- [polycli heimdall send checkpoint](polycli_heimdall_send_checkpoint.md) - Propose a checkpoint (MsgCheckpoint).

- [polycli heimdall send checkpoint-ack](polycli_heimdall_send_checkpoint-ack.md) - Acknowledge a checkpoint on L2 (MsgCpAck, L1-mirroring).

- [polycli heimdall send checkpoint-noack](polycli_heimdall_send_checkpoint-noack.md) - Mark missed checkpoint ack (MsgCpNoAck, L1-mirroring).

- [polycli heimdall send clerk-record](polycli_heimdall_send_clerk-record.md) - Submit an L1 state-sync record (MsgEventRecord, L1-mirroring).

- [polycli heimdall send signer-update](polycli_heimdall_send_signer-update.md) - Rotate validator signer pubkey (MsgSignerUpdate, L1-mirroring).

- [polycli heimdall send span-backfill](polycli_heimdall_send_span-backfill.md) - Trigger span backfill (MsgBackfillSpans).

- [polycli heimdall send span-propose](polycli_heimdall_send_span-propose.md) - Propose a new bor span (MsgProposeSpan).

- [polycli heimdall send span-set-downtime](polycli_heimdall_send_span-set-downtime.md) - Record producer downtime window (MsgSetProducerDowntime).

- [polycli heimdall send span-vote-producers](polycli_heimdall_send_span-vote-producers.md) - Vote for producers in the next span (MsgVoteProducers).

- [polycli heimdall send stake-exit](polycli_heimdall_send_stake-exit.md) - Mark validator exit (MsgValidatorExit, L1-mirroring).

- [polycli heimdall send stake-join](polycli_heimdall_send_stake-join.md) - Register a validator (MsgValidatorJoin, L1-mirroring).

- [polycli heimdall send stake-update](polycli_heimdall_send_stake-update.md) - Update validator stake (MsgStakeUpdate, L1-mirroring).

- [polycli heimdall send topup](polycli_heimdall_send_topup.md) - Credit validator fee balance (MsgTopupTx, L1-mirroring).

- [polycli heimdall send withdraw](polycli_heimdall_send_withdraw.md) - Withdraw accumulated validator fees.

