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
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// flags is injected by the caller via Register; every RunE uses
// config.Resolve on it to obtain a resolved *config.Config.
var flags *config.Flags

// Register attaches the tx-group subcommands directly to parent and
// binds the shared flag struct for config resolution.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	parent.AddCommand(
		newTxCmd(),
		newReceiptCmd(),
		newLogsCmd(),
		newNonceCmd(),
		newSequenceAliasCmd(),
		newBalanceCmd(),
		newRPCCmd(),
		newPublishCmd(),
	)
}

// newRPCClient resolves the config and constructs an RPCClient. When
// --curl is set the RPC call does not execute; it prints an
// equivalent curl command instead.
func newRPCClient(cmd *cobra.Command) (*client.RPCClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "tx package not registered (flags unset)"}
	}
	cfg, err := config.Resolve(flags)
	if err != nil {
		return nil, nil, &client.UsageError{Msg: err.Error()}
	}
	c := client.NewRPCClient(cfg.RPCURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		c.Transport = &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
	}
	return c, cfg, nil
}

// newRESTClient resolves the config and constructs a RESTClient.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "tx package not registered (flags unset)"}
	}
	cfg, err := config.Resolve(flags)
	if err != nil {
		return nil, nil, &client.UsageError{Msg: err.Error()}
	}
	c := client.NewRESTClient(cfg.RESTURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		c.Transport = &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
	}
	return c, cfg, nil
}

// renderOpts turns a resolved config into a render.Options instance,
// honouring --json, --field, --color, --raw, and TTY detection.
func renderOpts(cmd *cobra.Command, cfg *config.Config, fields []string) render.Options {
	return render.Options{
		JSON:   cfg.JSON,
		Raw:    cfg.Raw,
		Fields: fields,
		Color:  cfg.Color,
		IsTTY:  isTerminal(cmd.OutOrStdout()),
	}
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
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

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
