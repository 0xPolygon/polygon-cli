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
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/comet"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// pkg carries the package name and the flag struct injected via
// Register; cmdutil derives clients and render options from it.
var pkg = &cmdutil.Pkg{Name: "chain"}

// Register attaches the chain-group subcommands directly to parent
// and binds the shared flag struct for config resolution. parent is
// typically the root heimdall cobra command.
//
// Every chain subcommand is read-only, so we wire in render.EnableWatch
// to give them a `--watch DURATION` flag that repeats the call.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
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
// The shared response types and fetch helpers live in
// internal/heimdall/comet so the milestone family can reuse them; the
// aliases below keep this package's call sites unchanged.

type cometStatus = comet.Status

type cometBlock = comet.Block

type cometABCIInfo struct {
	Response struct {
		Data            string `json:"data"`
		Version         string `json:"version"`
		LastBlockHeight string `json:"last_block_height"`
	} `json:"response"`
}

// fetchStatus delegates to comet.FetchStatus.
func fetchStatus(ctx context.Context, rpc *client.RPCClient) (*cometStatus, error) {
	return comet.FetchStatus(ctx, rpc)
}

// fetchBlock delegates to comet.FetchBlock.
func fetchBlock(ctx context.Context, rpc *client.RPCClient, height string) (*cometBlock, json.RawMessage, error) {
	return comet.FetchBlock(ctx, rpc, height)
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
