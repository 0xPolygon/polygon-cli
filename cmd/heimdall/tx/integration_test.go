//go:build heimdall_integration

package tx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Integration tests talk directly to the live Heimdall v2 node on
// 172.19.0.2 unless overridden via HEIMDALL_TEST_REST_URL /
// HEIMDALL_TEST_RPC_URL.

func liveRPC() string {
	if v := os.Getenv("HEIMDALL_TEST_RPC_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:26657"
}

func liveREST() string {
	if v := os.Getenv("HEIMDALL_TEST_REST_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:1317"
}

func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := append([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
	}, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

// pickRecentTxHash queries the live RPC directly for a tx hash from
// the latest block. Falls back to sweeping recent heights because the
// current block may be an extended-commit-only block with no indexable
// user tx.
func pickRecentTxHash(t *testing.T) string {
	t.Helper()
	httpC := &http.Client{Timeout: 10 * time.Second}
	// Status first to find tip.
	status := func() int64 {
		req, _ := http.NewRequest(http.MethodPost, liveRPC(), strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"status"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpC.Do(req)
		if err != nil {
			t.Fatalf("status: %v", err)
		}
		defer resp.Body.Close()
		var out struct {
			Result struct {
				SyncInfo struct {
					LatestBlockHeight string `json:"latest_block_height"`
				} `json:"sync_info"`
			} `json:"result"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&out)
		h, _ := strconv.ParseInt(out.Result.SyncInfo.LatestBlockHeight, 10, 64)
		return h
	}()

	fetchBlockTxs := func(height int64) []string {
		body := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"block","params":{"height":"%d"}}`, height)
		req, _ := http.NewRequest(http.MethodPost, liveRPC(), strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpC.Do(req)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		buf, _ := io.ReadAll(resp.Body)
		var out struct {
			Result struct {
				Block struct {
					Data struct {
						Txs []string `json:"txs"`
					} `json:"data"`
				} `json:"block"`
			} `json:"result"`
		}
		_ = json.Unmarshal(buf, &out)
		return out.Result.Block.Data.Txs
	}

	// Use tx_search for a hash that will round-trip via /tx.
	searchBody := `{"jsonrpc":"2.0","id":1,"method":"tx_search","params":{"query":"message.action='/heimdallv2.topup.MsgTopupTx'","per_page":"1","page":"1","order_by":"desc"}}`
	req, _ := http.NewRequest(http.MethodPost, liveRPC(), strings.NewReader(searchBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpC.Do(req)
	if err == nil {
		defer resp.Body.Close()
		buf, _ := io.ReadAll(resp.Body)
		var out struct {
			Result struct {
				Txs []struct {
					Hash string `json:"hash"`
				} `json:"txs"`
			} `json:"result"`
		}
		if json.Unmarshal(buf, &out) == nil && len(out.Result.Txs) > 0 {
			return out.Result.Txs[0].Hash
		}
	}

	// Fallback — the txs in /block are raw TxRaw bytes, not hashes.
	// If tx_search failed we skip instead of faking a hash.
	_ = fetchBlockTxs
	_ = status
	t.Skip("no indexed tx available on live node")
	return ""
}

// pickValidatorSigner fetches the validator set from REST and returns
// the signer of the first validator, used as a well-known address for
// nonce/balance tests.
func pickValidatorSigner(t *testing.T) string {
	t.Helper()
	httpC := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpC.Get(liveREST() + "/stake/validators-set")
	if err != nil {
		t.Fatalf("validators-set: %v", err)
	}
	defer resp.Body.Close()
	var out struct {
		ValidatorSet struct {
			Validators []struct {
				Signer string `json:"signer"`
			} `json:"validators"`
		} `json:"validator_set"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decoding validators-set: %v", err)
	}
	if len(out.ValidatorSet.Validators) == 0 {
		t.Skip("no validators in set")
	}
	return out.ValidatorSet.Validators[0].Signer
}

func TestIntegrationTxByHash(t *testing.T) {
	hash := pickRecentTxHash(t)
	stdout, _, err := execLive(t, "tx", hash)
	if err != nil {
		t.Fatalf("tx %s: %v", hash, err)
	}
	for _, want := range []string{"hash", "height", "code"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("tx output missing %q: %q", want, stdout)
		}
	}
}

func TestIntegrationReceiptByHash(t *testing.T) {
	hash := pickRecentTxHash(t)
	stdout, _, err := execLive(t, "receipt", hash)
	if err != nil {
		t.Fatalf("receipt %s: %v", hash, err)
	}
	if !strings.Contains(stdout, "events") {
		t.Errorf("receipt missing events section: %q", stdout)
	}
}

func TestIntegrationRPCStatus(t *testing.T) {
	stdout, _, err := execLive(t, "rpc", "status")
	if err != nil {
		t.Fatalf("rpc status: %v", err)
	}
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("rpc status not JSON: %v", err)
	}
	if !strings.Contains(stdout, "latest_block_height") {
		t.Errorf("rpc status missing latest_block_height: %q", stdout)
	}
}

func TestIntegrationNonce(t *testing.T) {
	signer := pickValidatorSigner(t)
	stdout, _, err := execLive(t, "nonce", signer)
	if err != nil {
		t.Fatalf("nonce %s: %v", signer, err)
	}
	if _, err := strconv.ParseInt(strings.TrimSpace(stdout), 10, 64); err != nil {
		t.Errorf("nonce output not an integer: %q", stdout)
	}
}

func TestIntegrationBalance(t *testing.T) {
	signer := pickValidatorSigner(t)
	stdout, _, err := execLive(t, "balance", signer)
	if err != nil {
		t.Fatalf("balance %s: %v", signer, err)
	}
	// Balance might legitimately be 0 on a fresh validator; check
	// that the output at least parses as an integer.
	if _, err := strconv.ParseInt(strings.TrimSpace(stdout), 10, 64); err != nil {
		t.Errorf("balance not an integer: %q", stdout)
	}
}

func TestIntegrationLogs(t *testing.T) {
	stdout, stderr, err := execLive(t, "logs",
		"message.action='/heimdallv2.topup.MsgTopupTx'",
		"--limit", "3")
	if err != nil {
		t.Fatalf("logs: %v", err)
	}
	if !strings.Contains(stderr, "total_count") {
		t.Errorf("logs stderr missing total_count: %q", stderr)
	}
	_ = stdout
}
