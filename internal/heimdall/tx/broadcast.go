package tx

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// BroadcastMode mirrors cosmos.tx.v1beta1.BroadcastMode. We expose only
// SYNC (wait for CheckTx) and ASYNC (return immediately) — BLOCK is
// deprecated in cosmos-sdk and Heimdall has instant finality so SYNC +
// polling covers the real inclusion flow.
type BroadcastMode string

const (
	BroadcastModeSync  BroadcastMode = "BROADCAST_MODE_SYNC"
	BroadcastModeAsync BroadcastMode = "BROADCAST_MODE_ASYNC"
)

// BroadcastResult carries the fields polycli's send command surfaces
// after a broadcast round-trip. Raw holds the undecoded JSON response
// so --json callers can pass it through verbatim.
type BroadcastResult struct {
	TxHash    string
	Code      int
	Codespace string
	RawLog    string
	Height    uint64
	// Raw is the undecoded JSON body of the /cosmos/tx/v1beta1/txs
	// response. May be nil if the transport short-circuited (--curl).
	Raw []byte
}

// Broadcast submits txBytes via /cosmos/tx/v1beta1/txs and returns the
// decoded tx response. Callers typically pass the bytes returned by
// Builder.Sign; this function base64-encodes them before POSTing.
//
// On non-zero CheckTx code the function returns BroadcastResult with
// code/raw_log populated AND a non-nil error so callers can surface
// both without double-handling.
func Broadcast(ctx context.Context, rest *client.RESTClient, txBytes []byte, mode BroadcastMode) (*BroadcastResult, error) {
	if rest == nil {
		return nil, fmt.Errorf("Broadcast: rest client is nil")
	}
	if len(txBytes) == 0 {
		return nil, fmt.Errorf("Broadcast: tx bytes are empty")
	}
	if mode == "" {
		mode = BroadcastModeSync
	}
	req := map[string]any{
		"tx_bytes": base64.StdEncoding.EncodeToString(txBytes),
		"mode":     string(mode),
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding broadcast request: %w", err)
	}
	resp, _, err := rest.Post(ctx, "/cosmos/tx/v1beta1/txs", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("broadcasting tx: %w", err)
	}
	if resp == nil {
		return nil, nil // --curl short-circuit
	}
	var parsed struct {
		TxResponse struct {
			Height    string `json:"height"`
			TxHash    string `json:"txhash"`
			Code      int    `json:"code"`
			Codespace string `json:"codespace"`
			RawLog    string `json:"raw_log"`
		} `json:"tx_response"`
	}
	if err := json.Unmarshal(resp, &parsed); err != nil {
		return nil, fmt.Errorf("decoding broadcast response: %w (body=%q)", err, truncate(resp, 256))
	}
	result := &BroadcastResult{
		TxHash:    strings.ToLower(parsed.TxResponse.TxHash),
		Code:      parsed.TxResponse.Code,
		Codespace: parsed.TxResponse.Codespace,
		RawLog:    parsed.TxResponse.RawLog,
		Raw:       resp,
	}
	if parsed.TxResponse.Height != "" {
		// Height is 0 for SYNC (not yet included).
		if h, perr := parseUint(parsed.TxResponse.Height); perr == nil {
			result.Height = h
		}
	}
	if result.Code != 0 {
		return result, fmt.Errorf("broadcast returned code %d (%s): %s", result.Code, result.Codespace, result.RawLog)
	}
	return result, nil
}

// WaitForInclusion polls CometBFT /tx for txHash until the tx is found
// or ctx is cancelled. pollInterval defaults to 500ms when zero.
// Returns the height and raw /tx JSON body.
//
// Callers who need `--confirmations N` should call this to get
// inclusion height, then poll /status until
// latest_block_height >= height + N.
func WaitForInclusion(ctx context.Context, rpc *client.RPCClient, txHash string, pollInterval time.Duration) (uint64, []byte, error) {
	if rpc == nil {
		return 0, nil, fmt.Errorf("WaitForInclusion: rpc client is nil")
	}
	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}
	hashHex := strings.TrimPrefix(strings.TrimPrefix(txHash, "0x"), "0X")
	raw, err := hex.DecodeString(hashHex)
	if err != nil {
		return 0, nil, fmt.Errorf("WaitForInclusion: decoding hash %q: %w", txHash, err)
	}
	params := map[string]any{
		"hash":  base64.StdEncoding.EncodeToString(raw),
		"prove": false,
	}
	timer := time.NewTimer(0) // fire immediately on first iteration
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return 0, nil, ctx.Err()
		case <-timer.C:
		}
		resRaw, rerr := rpc.Call(ctx, "tx", params)
		if rerr == nil && resRaw != nil {
			var parsed struct {
				Height string `json:"height"`
			}
			if err := json.Unmarshal(resRaw, &parsed); err == nil && parsed.Height != "" {
				h, perr := parseUint(parsed.Height)
				if perr == nil {
					return h, resRaw, nil
				}
			}
		}
		timer.Reset(pollInterval)
	}
}

// WaitForConfirmations waits until the chain's latest_block_height is
// at least txHeight + confirmations. When confirmations is zero the
// function returns immediately.
func WaitForConfirmations(ctx context.Context, rpc *client.RPCClient, txHeight uint64, confirmations uint64, pollInterval time.Duration) error {
	if confirmations == 0 {
		return nil
	}
	if rpc == nil {
		return fmt.Errorf("WaitForConfirmations: rpc client is nil")
	}
	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}
	target := txHeight + confirmations
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		resRaw, err := rpc.Call(ctx, "status", nil)
		if err == nil && resRaw != nil {
			var parsed struct {
				SyncInfo struct {
					LatestBlockHeight string `json:"latest_block_height"`
				} `json:"sync_info"`
			}
			if err := json.Unmarshal(resRaw, &parsed); err == nil {
				if h, perr := parseUint(parsed.SyncInfo.LatestBlockHeight); perr == nil && h >= target {
					return nil
				}
			}
		}
		timer.Reset(pollInterval)
	}
}

// SimulateResult carries the gas estimate plus the raw response JSON.
type SimulateResult struct {
	GasUsed   uint64
	GasWanted uint64
	Raw       []byte
}

// Simulate calls /cosmos/tx/v1beta1/simulate with txBytes and returns
// the simulation result. Used by `estimate` and by `--gas auto`.
// Callers sign the tx before simulating — simulate still validates
// signatures, so a throwaway sig produced over the final doc works.
func Simulate(ctx context.Context, rest *client.RESTClient, txBytes []byte) (*SimulateResult, error) {
	if rest == nil {
		return nil, fmt.Errorf("Simulate: rest client is nil")
	}
	if len(txBytes) == 0 {
		return nil, fmt.Errorf("Simulate: tx bytes are empty")
	}
	req := map[string]any{
		"tx_bytes": base64.StdEncoding.EncodeToString(txBytes),
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding simulate request: %w", err)
	}
	resp, _, err := rest.Post(ctx, "/cosmos/tx/v1beta1/simulate", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("simulating tx: %w", err)
	}
	if resp == nil {
		return nil, nil // --curl short-circuit
	}
	var parsed struct {
		GasInfo struct {
			GasWanted string `json:"gas_wanted"`
			GasUsed   string `json:"gas_used"`
		} `json:"gas_info"`
	}
	if err := json.Unmarshal(resp, &parsed); err != nil {
		return nil, fmt.Errorf("decoding simulate response: %w (body=%q)", err, truncate(resp, 256))
	}
	out := &SimulateResult{Raw: resp}
	if parsed.GasInfo.GasUsed != "" {
		if n, perr := parseUint(parsed.GasInfo.GasUsed); perr == nil {
			out.GasUsed = n
		}
	}
	if parsed.GasInfo.GasWanted != "" {
		if n, perr := parseUint(parsed.GasInfo.GasWanted); perr == nil {
			out.GasWanted = n
		}
	}
	return out, nil
}

// parseUint accepts decimal-string uint64 values as returned by the
// Cosmos SDK REST gateway.
func parseUint(s string) (uint64, error) {
	var n uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("parseUint: non-digit %q in %q", c, s)
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}

// truncate clips b at n bytes with a trailing "..." for error messages.
func truncate(b []byte, n int) []byte {
	if len(b) <= n {
		return b
	}
	return append(append([]byte{}, b[:n]...), []byte("...")...)
}
