// Package tx implements the cast-familiar tx/account read-only
// subcommands of `polycli heimdall`: tx, receipt, logs, nonce,
// sequence, balance, rpc, publish.
//
// The subcommands live at the top level of the heimdall tree (for
// cast parity) rather than under an intermediate group. Callers
// register them with Register(parent, flags).
package tx

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// pkg carries the package name and the flag struct injected via
// Register; cmdutil derives clients and render options from it.
var pkg = &cmdutil.Pkg{Name: "tx"}

// Register attaches the tx-group subcommands directly to parent and
// binds the shared flag struct for config resolution.
//
// Read-only subcommands (tx, receipt, logs, nonce, sequence, balance)
// get `--watch DURATION` via render.EnableWatch so operators can watch
// a value change over time. Publish is one-shot and rpc is a raw
// passthrough; neither is watchable.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	readOnly := []*cobra.Command{
		newTxCmd(),
		newReceiptCmd(),
		newLogsCmd(),
		newNonceCmd(),
		newSequenceAliasCmd(),
		newBalanceCmd(),
	}
	for _, c := range readOnly {
		render.EnableWatch(c)
		parent.AddCommand(c)
	}
	parent.AddCommand(newRPCCmd())
	parent.AddCommand(newPublishCmd())
	// The mktx/send/estimate umbrellas are attached separately so
	// their child-msg factories can live in the tx/msgs sub-package
	// without circular imports. Each umbrella owns its own copy of
	// the registered msg subcommands (cobra command trees are single-
	// parent).
	registerMktxSendEstimate(parent, f)
}

// normalizeHash accepts a tx hash with or without `0x` prefix and
// returns the upper-case `0x`-prefixed hex form expected by CometBFT's
// /tx endpoint. Returns a UsageError when the hash is not 32 bytes.
func normalizeHash(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", &client.UsageError{Msg: fmt.Sprintf("invalid tx hash %q: %v", raw, err)}
	}
	if len(b) != 32 {
		return "", &client.UsageError{Msg: fmt.Sprintf("tx hash must be 32 bytes (64 hex chars), got %d", len(b))}
	}
	return "0x" + strings.ToUpper(s), nil
}

// validateAddress accepts a 20-byte Ethereum-style address with or
// without the `0x` prefix. Returns the canonical lowercase
// `0x`-prefixed form. A bech32 decoder is out of scope here — that
// surface lives in `polycli heimdall addr`.
func validateAddress(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", &client.UsageError{Msg: fmt.Sprintf("invalid address %q: %v", raw, err)}
	}
	if len(b) != 20 {
		return "", &client.UsageError{Msg: fmt.Sprintf("address must be 20 bytes (40 hex chars), got %d", len(b))}
	}
	return "0x" + strings.ToLower(s), nil
}
