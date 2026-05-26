# internal/heimdall/proto — hand-rolled Cosmos SDK tx wire encoding

This package implements the subset of the Cosmos SDK / Heimdall tx
protobuf wire format that polycli's heimdall tx builder
(`internal/heimdall/tx/`) needs to build, sign, and broadcast
transactions against a Heimdall v2 node.

## Why hand-rolled?

The obvious choice is `buf generate` against heimdall-v2's `.proto` files.
We did not take it because heimdall-v2's dependency closure (cosmos-sdk,
cometbft, go-ethereum) requires three Polygon-fork `replace` directives
in go.mod:

- `github.com/ethereum/go-ethereum` → `github.com/0xPolygon/bor`
- `github.com/cosmos/cosmos-sdk` → `github.com/0xPolygon/cosmos-sdk`
- `github.com/cometbft/cometbft` → `github.com/0xPolygon/cometbft`

Adopting any of those in polycli's root `go.mod` would risk breaking
every existing command that already pulls upstream go-ethereum
(loadtest, monitor, fund, wallet, rpcfuzz, abi, ...). The W2 brief is
explicit: do not add replace directives.

Proto wire format is a small, well-specified byte-level protocol
(<https://protobuf.dev/programming-guides/encoding>). The tx builder
only needs a handful of messages:

- `cosmos.tx.v1beta1.{TxBody, AuthInfo, SignerInfo, ModeInfo,
  ModeInfo.Single, Fee, TxRaw, SignDoc}`
- `cosmos.base.v1beta1.Coin`
- `cosmos.crypto.secp256k1.PubKey`
- `google.protobuf.Any`
- `heimdallv2.topup.MsgWithdrawFeeTx` (starter Msg; others added in W3/W4)

Each is encoded by appending tag + length-prefixed bytes using
`google.golang.org/protobuf/encoding/protowire`, which is part of the
standard protobuf distribution and already a direct dep. Decoding is
the mirror operation. Total hand-rolled code is ~300 lines and can be
audited in one sitting.

If the tx-builder surface grows large enough that hand-rolling stops
being tractable, the escape hatch is `buf generate` with a local
buf.gen.yaml pointing at heimdall-v2's proto tree, outputting to this
package. That path stays open.

## Signing specifics

Heimdall v2 uses `PubKeySecp256k1eth`: standard secp256k1 curve, but
signing digests are `keccak256(signBytes)` rather than `sha256`, and
the 65-byte signature (r||s||v) is truncated to 64 bytes (r||s) before
storing in `TxRaw.signatures`. See `internal/heimdall/tx/sign.go` for
the implementation.
