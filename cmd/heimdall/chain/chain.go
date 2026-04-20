// Package chain implements the cast-familiar CometBFT-facing
// subcommands of `polycli heimdall`: block, block-number, age,
// find-block, chain-id, chain, client. All calls target the CometBFT
// JSON-RPC endpoint — the REST gateway is unused here.
//
// The subcommands live at the top level of the heimdall tree (for
// cast parity) rather than under an intermediate `chain` group.
// Callers register them with Register(parent).
package chain

import (
	_ "embed"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// Flags is injected by the caller. The persistent flag set on the
// root heimdall command is the source of truth for global config
// resolution. Store it as a package variable so each RunE can call
// config.Resolve.
var flags *config.Flags

// Register attaches the chain-group subcommands directly to parent
// and binds the shared flag struct for config resolution. parent is
// typically the root heimdall cobra command.
//
// Every chain subcommand is read-only, so we wire in render.EnableWatch
// to give them a `--watch DURATION` flag that repeats the call.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	subs := []*cobra.Command{
		newBlockCmd(),
		newBlockNumberCmd(),
		newAgeCmd(),
		newFindBlockCmd(),
		newChainIDCmd(),
		newChainCmd(),
		newClientCmd(),
	}
	for _, s := range subs {
		render.EnableWatch(s)
		parent.AddCommand(s)
	}
}

// chainNames maps the well-known Heimdall v2 chain ids to their
// human-readable marketing names. Additional ids fall through to a
// generic "unknown chain" response in the `chain` subcommand.
var chainNames = map[string]string{
	"heimdallv2-137":   "Polygon Mainnet",
	"heimdallv2-80002": "Polygon Amoy Testnet",
}

// --- shared helpers ---

// newRPCClient resolves the config and constructs an RPCClient. When
// --curl is set the RPC call does not execute; it prints an
// equivalent curl command instead.
func newRPCClient(cmd *cobra.Command) (*client.RPCClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "chain package not registered (flags unset)"}
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

// resolveHeight converts a CLI height argument (empty, "latest",
// "earliest", or a bare decimal) into the `height` param accepted by
// CometBFT's /block endpoint. An empty string is returned for the
// latest-block shorthand; CometBFT interprets a missing `height` as
// "latest".
func resolveHeight(ctx context.Context, rpc *client.RPCClient, arg string) (string, error) {
	tag := strings.ToLower(strings.TrimSpace(arg))
	switch tag {
	case "", "latest":
		return "", nil
	case "earliest":
		st, err := fetchStatus(ctx, rpc)
		if err != nil {
			return "", err
		}
		if st == nil {
			return "", nil // --curl
		}
		if st.SyncInfo.EarliestBlockHeight == "" {
			return "", fmt.Errorf("status did not contain earliest_block_height")
		}
		return st.SyncInfo.EarliestBlockHeight, nil
	case "finalized", "safe", "pending":
		return "", &client.UsageError{Msg: fmt.Sprintf(
			"block tag %q is not valid on Heimdall (instant finality); use `latest` or `earliest`", tag)}
	}
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return "", &client.UsageError{Msg: fmt.Sprintf("invalid height %q (want integer, `latest`, or `earliest`)", arg)}
	}
	if n <= 0 {
		return "", &client.UsageError{Msg: fmt.Sprintf("height must be positive, got %d", n)}
	}
	return strconv.FormatInt(n, 10), nil
}

// --- CometBFT response types ---
//
// We decode only the fields actually used; everything else stays in
// the raw JSON (and is available via --json / --field).

type cometStatus struct {
	NodeInfo struct {
		Version string `json:"version"`
		Network string `json:"network"`
		Moniker string `json:"moniker"`
	} `json:"node_info"`
	SyncInfo struct {
		LatestBlockHeight   string `json:"latest_block_height"`
		LatestBlockTime     string `json:"latest_block_time"`
		EarliestBlockHeight string `json:"earliest_block_height"`
		CatchingUp          bool   `json:"catching_up"`
	} `json:"sync_info"`
}

type cometBlock struct {
	BlockID struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			ChainID         string `json:"chain_id"`
			Height          string `json:"height"`
			Time            string `json:"time"`
			ProposerAddress string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
	} `json:"block"`
}

type cometABCIInfo struct {
	Response struct {
		Data            string `json:"data"`
		Version         string `json:"version"`
		LastBlockHeight string `json:"last_block_height"`
	} `json:"response"`
}

// fetchStatus calls the CometBFT /status RPC and returns the decoded
// struct. Returns (nil, nil) when running under --curl.
func fetchStatus(ctx context.Context, rpc *client.RPCClient) (*cometStatus, error) {
	raw, err := rpc.Call(ctx, "status", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching status: %w", err)
	}
	if raw == nil {
		return nil, nil
	}
	var st cometStatus
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, fmt.Errorf("decoding status: %w", err)
	}
	return &st, nil
}

// fetchBlock calls CometBFT /block at the given height (empty ==
// latest). Returns the typed struct, the raw result bytes (for --json
// passthrough), and any error. Both return values are nil when --curl
// short-circuits.
func fetchBlock(ctx context.Context, rpc *client.RPCClient, height string) (*cometBlock, json.RawMessage, error) {
	// CometBFT's reflect-based RPC requires an explicit `height`
	// key in params; a missing or empty params object returns
	// "reflect: Call with too few input arguments". Pass nil height
	// to request the latest block.
	params := map[string]any{"height": nil}
	if height != "" {
		params["height"] = height
	}
	raw, err := rpc.Call(ctx, "block", params)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching block: %w", err)
	}
	if raw == nil {
		return nil, nil, nil
	}
	var blk cometBlock
	if err := json.Unmarshal(raw, &blk); err != nil {
		return nil, nil, fmt.Errorf("decoding block: %w", err)
	}
	return &blk, raw, nil
}

// fetchABCIInfo calls CometBFT /abci_info.
func fetchABCIInfo(ctx context.Context, rpc *client.RPCClient) (*cometABCIInfo, error) {
	raw, err := rpc.Call(ctx, "abci_info", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching abci_info: %w", err)
	}
	if raw == nil {
		return nil, nil
	}
	var out cometABCIInfo
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decoding abci_info: %w", err)
	}
	return &out, nil
}

// unixFromRFC3339Nano parses a CometBFT block timestamp into a unix-
// second integer string. CometBFT emits RFC3339Nano UTC by default
// but older nodes may drop the nanoseconds.
func unixFromRFC3339Nano(ts string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		t, err = time.Parse(time.RFC3339, ts)
		if err != nil {
			return "", fmt.Errorf("parsing block time %q: %w", ts, err)
		}
	}
	return strconv.FormatInt(t.Unix(), 10), nil
}
