package msgs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// Clients bundles the REST + RPC clients a msg subcommand needs to
// talk to Heimdall. Resolved per-invocation (not per-registration) so
// each run can honour the network flags.
type Clients struct {
	Cfg  *config.Config
	REST *client.RESTClient
	RPC  *client.RPCClient
}

// ResolveClients builds REST + RPC clients from opts.Global. Returns
// a UsageError when the global config is missing or malformed.
func ResolveClients(cmd *cobra.Command, opts *TxOpts) (*Clients, error) {
	if opts == nil || opts.Global == nil {
		return nil, &client.UsageError{Msg: "tx flags not registered (internal wiring error)"}
	}
	cfg, err := config.Resolve(opts.Global)
	if err != nil {
		return nil, &client.UsageError{Msg: err.Error()}
	}
	rest := client.NewRESTClient(cfg.RESTURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	rpc := client.NewRPCClient(cfg.RPCURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		tr := &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
		rest.Transport = tr
		rpc.Transport = tr
	}
	return &Clients{Cfg: cfg, REST: rest, RPC: rpc}, nil
}

// Plan carries the msg-subcommand's parsed output into Execute: the
// list of Msgs to include in the TxBody and the signer's on-chain
// identifier (used as the from / proposer field and as the account
// lookup key). MsgShortType is the short msg name (e.g.
// "MsgWithdrawFeeTx") used by the L1-mirroring force guard.
type Plan struct {
	Msgs         []htx.Msg
	MsgShortType string
	// SignerAddress is the address used to fetch account
	// number/sequence and to populate Msg.Proposer-like fields when
	// the msg subcommand did not set them explicitly. Typically equal
	// to the resolved signer's Eth address.
	SignerAddress string
}

// Execute runs the full mktx/send/estimate pipeline for the given
// mode. Msg subcommands build their plan via Planner and hand it to
// Execute; everything else — account fetch, signing, broadcast,
// simulate — lives here so the mode-specific branches stay short.
//
// Output is written to cmd.OutOrStdout(). Errors propagate; callers
// should not wrap them in a generic "failed" message because the
// tx/client packages already attach context.
func Execute(cmd *cobra.Command, opts *TxOpts, mode Mode, plan *Plan) error {
	if plan == nil || len(plan.Msgs) == 0 {
		return fmt.Errorf("Execute: plan has no messages")
	}
	// L1-mirroring guard: block msgs the bridge owns unless --force.
	if err := htx.RequireForce(plan.MsgShortType, opts.Force); err != nil {
		return &client.UsageError{Msg: err.Error()}
	}

	signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
	if err != nil {
		return err
	}

	clients, err := ResolveClients(cmd, opts)
	if err != nil {
		return err
	}
	signMode, err := htx.ParseSignMode(opts.SignMode)
	if err != nil {
		return &client.UsageError{Msg: err.Error()}
	}

	ctx := cmd.Context()

	// Build the Builder with all shared inputs.
	b := htx.NewBuilder().
		WithChainID(clients.Cfg.ChainID).
		WithSignMode(signMode).
		WithMemo(opts.Memo)
	for _, m := range plan.Msgs {
		b.AddMsg(m)
	}
	// Account overrides first — if both are set we skip the network
	// fetch entirely, which is the only sane behaviour under --curl.
	if opts.AccountNumber > 0 {
		b.WithAccountNumber(opts.AccountNumber)
	}
	if opts.Sequence > 0 {
		b.WithSequence(opts.Sequence)
	}
	if opts.AccountNumber == 0 || opts.Sequence == 0 {
		// Auto-fetch the missing pieces.
		if plan.SignerAddress == "" {
			return &client.UsageError{Msg: "signer address is empty; cannot fetch account info"}
		}
		fetcher := &htx.RESTAccountFetcher{Client: clients.REST}
		if err := b.ResolveAccount(ctx, fetcher, plan.SignerAddress); err != nil {
			return err
		}
	}

	// Gas & fee. For `estimate` we don't need a real gas limit up
	// front — the simulate endpoint accepts a zero-gas tx. Everywhere
	// else we do (the builder refuses to sign without one).
	gasLimit := opts.Gas
	if gasLimit == 0 {
		// A placeholder large enough that common checkpoint/stake
		// txs pass simulation without hitting the miner-chosen cap.
		gasLimit = 200_000
	}
	b.WithGasLimit(gasLimit)

	if opts.Fee != "" {
		coin, err := parseFeeCoin(opts.Fee, clients.Cfg.Denom)
		if err != nil {
			return &client.UsageError{Msg: err.Error()}
		}
		b.WithFee(coin)
	} else if opts.GasPrice > 0 {
		coin, err := computeFeeFromGasPrice(opts.GasPrice, gasLimit, clients.Cfg.Denom)
		if err != nil {
			return err
		}
		b.WithFee(coin)
	}

	raw, err := b.Sign(signer.Key)
	if err != nil {
		return err
	}

	switch mode {
	case ModeMkTx:
		return printMkTxResult(cmd.OutOrStdout(), raw, opts.JSONOut)
	case ModeEstimate:
		sim, err := htx.Simulate(ctx, clients.REST, raw)
		if err != nil {
			return err
		}
		if sim == nil { // --curl
			return nil
		}
		return printEstimateResult(cmd.OutOrStdout(), sim, opts, clients.Cfg.Denom)
	case ModeSend:
		return runSend(ctx, cmd.OutOrStdout(), clients, raw, opts)
	default:
		return fmt.Errorf("Execute: unsupported mode %v", mode)
	}
}

// runSend handles the send-mode branch: dry-run short-circuit,
// broadcast, optional wait-for-inclusion, and optional confirmation
// polling.
func runSend(ctx context.Context, out io.Writer, clients *Clients, raw []byte, opts *TxOpts) error {
	if opts.DryRun {
		return printDryRun(out, raw, opts)
	}
	mode := htx.BroadcastModeSync
	if opts.Async {
		mode = htx.BroadcastModeAsync
	}
	res, err := htx.Broadcast(ctx, clients.REST, raw, mode)
	if err != nil {
		if res != nil {
			// Render the failure envelope so operators see the code /
			// raw_log alongside the error message.
			_ = printBroadcastResult(out, res, opts.JSONOut)
		}
		return err
	}
	if res == nil { // --curl
		return nil
	}
	if err := printBroadcastResult(out, res, opts.JSONOut); err != nil {
		return err
	}
	if opts.Async {
		return nil
	}
	height, _, err := htx.WaitForInclusion(ctx, clients.RPC, res.TxHash, 0)
	if err != nil {
		return err
	}
	if opts.Confirmations > 0 {
		if err := htx.WaitForConfirmations(ctx, clients.RPC, height, opts.Confirmations, 0); err != nil {
			return err
		}
	}
	if opts.JSONOut {
		return json.NewEncoder(out).Encode(map[string]any{
			"txhash":        res.TxHash,
			"height":        height,
			"confirmations": opts.Confirmations,
		})
	}
	_, err = fmt.Fprintf(out, "included height=%d confirmations=%d\n", height, opts.Confirmations)
	return err
}

// printMkTxResult writes the TxRaw hex (and optional JSON envelope)
// to out. Used by both `mktx` and `send --dry-run`.
func printMkTxResult(out io.Writer, raw []byte, jsonOut bool) error {
	hexStr := "0x" + hexEncode(raw)
	b64 := base64.StdEncoding.EncodeToString(raw)
	if jsonOut {
		return json.NewEncoder(out).Encode(map[string]any{
			"tx_raw_hex": hexStr,
			"tx_raw_b64": b64,
		})
	}
	_, err := fmt.Fprintln(out, hexStr)
	return err
}

// printDryRun shows what would be sent and bails out before the POST.
func printDryRun(out io.Writer, raw []byte, opts *TxOpts) error {
	hexStr := "0x" + hexEncode(raw)
	b64 := base64.StdEncoding.EncodeToString(raw)
	if opts.JSONOut {
		return json.NewEncoder(out).Encode(map[string]any{
			"dry_run":    true,
			"tx_raw_hex": hexStr,
			"tx_raw_b64": b64,
		})
	}
	_, err := fmt.Fprintf(out, "dry-run=true\ntx_raw_hex=%s\ntx_raw_b64=%s\n", hexStr, b64)
	return err
}

// printBroadcastResult renders the TxResponse envelope after
// /cosmos/tx/v1beta1/txs. Height is 0 on BROADCAST_MODE_SYNC (the
// chain hasn't included the tx yet); we print it anyway for clarity.
func printBroadcastResult(out io.Writer, res *htx.BroadcastResult, jsonOut bool) error {
	if jsonOut && res.Raw != nil {
		// Pass the raw envelope through verbatim when --json is set
		// so operators can post-process with jq without us losing
		// any fields.
		var generic any
		if err := json.Unmarshal(res.Raw, &generic); err == nil {
			return json.NewEncoder(out).Encode(generic)
		}
	}
	_, err := fmt.Fprintf(out, "txhash=0x%s\ncode=%d\nheight=%d\n", res.TxHash, res.Code, res.Height)
	if err != nil {
		return err
	}
	if res.RawLog != "" {
		_, err = fmt.Fprintf(out, "raw_log=%s\n", res.RawLog)
	}
	return err
}

// printEstimateResult renders the simulate result. When --gas-price
// is set it also reports the implied fee for the reported gas_used.
func printEstimateResult(out io.Writer, sim *htx.SimulateResult, opts *TxOpts, denom string) error {
	if opts.JSONOut && sim.Raw != nil {
		var generic any
		if err := json.Unmarshal(sim.Raw, &generic); err == nil {
			return json.NewEncoder(out).Encode(generic)
		}
	}
	if _, err := fmt.Fprintf(out, "gas_used=%d\ngas_wanted=%d\n", sim.GasUsed, sim.GasWanted); err != nil {
		return err
	}
	if opts.GasPrice > 0 && sim.GasUsed > 0 {
		coin, err := computeFeeFromGasPrice(opts.GasPrice, sim.GasUsed, denom)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "fee=%s%s\n", coin.Amount, coin.Denom)
		return err
	}
	return nil
}

// --- Fee helpers. ---

// parseFeeCoin accepts a string like "10000pol" or "10000" and
// returns a Coin. A bare number uses fallbackDenom.
func parseFeeCoin(s, fallbackDenom string) (hproto.Coin, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return hproto.Coin{}, fmt.Errorf("empty --fee")
	}
	// Split at first non-digit.
	idx := 0
	for idx < len(s) && s[idx] >= '0' && s[idx] <= '9' {
		idx++
	}
	if idx == 0 {
		return hproto.Coin{}, fmt.Errorf("--fee %q: expected leading amount", s)
	}
	amount := s[:idx]
	denom := strings.TrimSpace(s[idx:])
	if denom == "" {
		denom = fallbackDenom
	}
	if denom == "" {
		return hproto.Coin{}, fmt.Errorf("--fee %q: missing denom and no default denom configured", s)
	}
	return hproto.Coin{Denom: denom, Amount: amount}, nil
}

// computeFeeFromGasPrice returns a Coin whose amount is
// ceil(gasPrice * gasLimit). gasPrice is expressed in whole denom
// units per gas unit (e.g. 0.0000001 pol/gas). Heimdall's fee
// handling accepts decimal amounts via string-serialized big.Int, so
// we compute in big.Float for precision and convert to the smallest
// integer >= gasPrice*gasLimit.
func computeFeeFromGasPrice(gasPrice float64, gasLimit uint64, denom string) (hproto.Coin, error) {
	if gasPrice <= 0 {
		return hproto.Coin{}, fmt.Errorf("gas price must be positive, got %v", gasPrice)
	}
	if denom == "" {
		return hproto.Coin{}, fmt.Errorf("denom is empty; set --denom or configure HEIMDALL_FEE_DENOM")
	}
	bf := new(big.Float).Mul(big.NewFloat(gasPrice), new(big.Float).SetUint64(gasLimit))
	// Ceiling to integer.
	bi, _ := bf.Int(nil)
	// If bf was not already integer, increment by 1 so we round up.
	check := new(big.Float).SetInt(bi)
	if check.Cmp(bf) < 0 {
		bi = new(big.Int).Add(bi, big.NewInt(1))
	}
	return hproto.Coin{Denom: denom, Amount: bi.String()}, nil
}

// hexEncode is a lower-case hex encoder used to print TxRaw bytes.
// Avoids importing encoding/hex twice in tests.
func hexEncode(raw []byte) string {
	const alpha = "0123456789abcdef"
	out := make([]byte, 2*len(raw))
	for i, c := range raw {
		out[2*i] = alpha[c>>4]
		out[2*i+1] = alpha[c&0x0f]
	}
	return string(out)
}
