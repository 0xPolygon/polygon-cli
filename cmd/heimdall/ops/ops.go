// Package ops implements the `polycli heimdall ops` umbrella command
// and its CometBFT JSON-RPC-facing subcommands: status, health, peers,
// consensus, tx-pool, abci-info, commit, and validators-cometbft.
//
// All calls target the CometBFT RPC endpoint (`:26657`). The
// Heimdall REST gateway is unused here. The umbrella keeps node-ops
// commands grouped so operators can find them without scrolling
// through the flat cast-like tree.
package ops

import (
	_ "embed"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// flags is injected by Register. Each subcommand reads it via
// config.Resolve when building its RPC client.
var flags *config.Flags

// newOpsCmd builds a fresh `ops` umbrella. Constructed per Register
// call so tests that re-wire a parent do not accumulate duplicate
// subcommands on a shared command tree.
func newOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Node-operator commands backed by CometBFT JSON-RPC.",
		Long:  usage,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		newStatusCmd(),
		newHealthCmd(),
		newPeersCmd(),
		newConsensusCmd(),
		newTxPoolCmd(),
		newABCIInfoCmd(),
		newCommitCmd(),
		newValidatorsCometBFTCmd(),
	)
	return cmd
}

// Register attaches the ops umbrella and its subcommands to parent,
// wiring in the shared flag struct used for config resolution.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	parent.AddCommand(newOpsCmd())
}

// newRPCClient resolves the heimdall config and constructs an RPC
// client. When --curl is set the client's Transport is swapped for one
// that prints the equivalent POST command instead of executing it.
func newRPCClient(cmd *cobra.Command) (*client.RPCClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "ops package not registered (flags unset)"}
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

// renderOpts returns a render.Options honouring --json/--field/--color
// plus TTY detection.
func renderOpts(cmd *cobra.Command, cfg *config.Config, fields []string) render.Options {
	return render.Options{
		JSON:   cfg.JSON,
		Raw:    cfg.Raw,
		Fields: fields,
		Color:  cfg.Color,
		IsTTY:  isTerminal(cmd.OutOrStdout()),
	}
}

// isTerminal returns true if w is an *os.File attached to a terminal.
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

// callEmpty issues an RPC call with an explicit empty params object.
// CometBFT's reflect-based RPC layer rejects nil params on some
// methods with "reflect: Call with too few input arguments"; passing
// map[string]any{} avoids that trap while still producing a valid
// JSON-RPC envelope.
func callEmpty(ctx context.Context, rpc *client.RPCClient, method string) (json.RawMessage, error) {
	raw, err := rpc.Call(ctx, method, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("calling %s: %w", method, err)
	}
	return raw, nil
}

// decodeGeneric unmarshals raw into any (map/slice/string). Used when
// we want to emit --json passthrough or pluck via --field.
func decodeGeneric(raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return v, nil
}
