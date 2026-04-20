package tx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// --- normalizeHash ---

func TestNormalizeHash(t *testing.T) {
	const raw = "94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29"
	cases := []struct {
		in      string
		wantErr bool
	}{
		{raw, false},
		{"0x" + raw, false},
		{"0X" + raw, false},
		{strings.ToLower(raw), false},
		{"", true},
		{"0x12", true},
		{"zz", true},
	}
	for _, c := range cases {
		got, err := normalizeHash(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("in=%q err=%v wantErr=%v", c.in, err, c.wantErr)
			continue
		}
		if err != nil {
			continue
		}
		if !strings.HasPrefix(got, "0x") || len(got) != 66 {
			t.Errorf("in=%q got=%q", c.in, got)
		}
	}
}

// --- validateAddress ---

func TestValidateAddress(t *testing.T) {
	const raw = "02f615e95563ef16f10354dba9e584e58d2d4314"
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{raw, "0x" + raw, false},
		{"0x" + raw, "0x" + raw, false},
		{"0X" + strings.ToUpper(raw), "0x" + raw, false},
		{raw[:10], "", true},
		{"", "", true},
		{"garbage", "", true},
	}
	for _, c := range cases {
		got, err := validateAddress(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("in=%q err=%v wantErr=%v", c.in, err, c.wantErr)
		}
		if !c.wantErr && got != c.want {
			t.Errorf("in=%q got=%q want=%q", c.in, got, c.want)
		}
	}
}

// --- tx command ---

func TestTxCmdHumanOutput(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	hash := "94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29"
	stdout, _, err := runCmd(t, "", srv.URL, "tx", hash)
	if err != nil {
		t.Fatalf("tx: %v", err)
	}
	for _, want := range []string{"hash", "height", "code", "gas_used", "num_events"} {
		mustContain(t, stdout, want)
	}
	mustContain(t, stdout, "31579117")
}

func TestTxCmdAcceptsHashWithoutPrefix(t *testing.T) {
	// Same as TestTxCmdHumanOutput but without the 0x — covers the
	// prefix-tolerance requirement.
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "tx",
		"94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29")
	if err != nil {
		t.Fatalf("tx: %v", err)
	}
	mustContain(t, stdout, "31579117")
}

func TestTxCmdJSON(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "tx",
		"0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29",
		"--json")
	if err != nil {
		t.Fatalf("tx --json: %v", err)
	}
	var out any
	if uerr := json.Unmarshal([]byte(stdout), &out); uerr != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", uerr, stdout)
	}
}

func TestTxCmdRawPreservesTxBody(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "tx",
		"0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29",
		"--raw")
	if err != nil {
		t.Fatalf("tx --raw: %v", err)
	}
	mustContain(t, stdout, "tx")
}

func TestTxCmdInvalidHash(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	_, _, err := runCmd(t, "", srv.URL, "tx", "0x1234")
	if err == nil {
		t.Fatal("expected error for short hash")
	}
	var uErr *client.UsageError
	if !errorsAs(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- receipt command ---

func TestReceiptRendersEvents(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx": loadFixture(t, "rpc", "tx.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "receipt",
		"0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29")
	if err != nil {
		t.Fatalf("receipt: %v", err)
	}
	mustContain(t, stdout, "events")
	mustContain(t, stdout, "coin_spent")
	mustContain(t, stdout, "spender = 0x")
}

func TestReceiptConfirmationsNegativeRejected(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{})
	_, _, err := runCmd(t, "", srv.URL, "receipt",
		"0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29",
		"--confirmations", "-1")
	if err == nil {
		t.Fatal("expected error for negative confirmations")
	}
	var uErr *client.UsageError
	if !errorsAs(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// TestReceiptConfirmationsWaits stubs /tx + /status and verifies the
// receipt command re-polls status until the tip reaches the target
// height.
func TestReceiptConfirmationsWaits(t *testing.T) {
	txBody := loadFixture(t, "rpc", "tx.json")
	// tx fixture height is 31579117; require 2 confirmations so target=31579119.
	var statusCalls atomic.Int64
	fakeStatus := func(tip int64) []byte {
		env := map[string]any{
			"jsonrpc": "2.0", "id": 1,
			"result": map[string]any{
				"node_info": map[string]any{},
				"sync_info": map[string]any{
					"latest_block_height": strconv.FormatInt(tip, 10),
				},
			},
		}
		b, _ := json.Marshal(env)
		return b
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string `json:"method"`
			ID     uint64 `json:"id"`
		}
		_ = json.Unmarshal(body, &req)
		switch req.Method {
		case "tx":
			var env map[string]any
			_ = json.Unmarshal(txBody, &env)
			env["id"] = req.ID
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		case "status":
			n := statusCalls.Add(1)
			tip := int64(31579117) + n - 1 // first call: tip=31579117; second: 31579118; third: 31579119 (ok)
			var env map[string]any
			_ = json.Unmarshal(fakeStatus(tip), &env)
			env["id"] = req.ID
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		default:
			http.Error(w, "no route "+req.Method, 404)
		}
	}))
	defer srv.Close()

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"receipt", "0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29",
		"--confirmations", "2",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tighten the poll interval via a custom wrapper — we can't easily
	// poke pollInt from outside, so we just let the server advance
	// the tip deterministically on each call. The default 500ms means
	// this test takes ~1s; acceptable.
	if err := root.ExecuteContext(ctx); err != nil {
		t.Fatalf("receipt --confirmations 2: %v", err)
	}
	if statusCalls.Load() < 2 {
		t.Errorf("expected at least 2 status calls, got %d", statusCalls.Load())
	}
}

// TestReceiptConfirmationsCancels verifies the poll loop exits
// promptly when its context is cancelled.
func TestReceiptConfirmationsCancels(t *testing.T) {
	txBody := loadFixture(t, "rpc", "tx.json")
	// A status response with tip well below target, so the loop never
	// completes on its own.
	stuckStatus := func() []byte {
		env := map[string]any{
			"jsonrpc": "2.0", "id": 1,
			"result": map[string]any{
				"sync_info": map[string]any{"latest_block_height": "1"},
			},
		}
		b, _ := json.Marshal(env)
		return b
	}()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string `json:"method"`
			ID     uint64 `json:"id"`
		}
		_ = json.Unmarshal(body, &req)
		switch req.Method {
		case "tx":
			var env map[string]any
			_ = json.Unmarshal(txBody, &env)
			env["id"] = req.ID
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		case "status":
			var env map[string]any
			_ = json.Unmarshal(stuckStatus, &env)
			env["id"] = req.ID
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		default:
			http.Error(w, "no route", 404)
		}
	}))
	defer srv.Close()

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"receipt", "0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29",
		"--confirmations", "100",
	})
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	var errOut error
	start := time.Now()
	go func() {
		defer wg.Done()
		errOut = root.ExecuteContext(ctx)
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	wg.Wait()
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("receipt took %s to cancel, expected < 2s", elapsed)
	}
	if errOut == nil {
		t.Errorf("expected cancellation error")
	}
	if !errors.Is(errOut, context.Canceled) {
		t.Logf("note: got %v (context.Canceled preferred but any non-nil is acceptable)", errOut)
	}
}

// --- logs command ---

func TestLogsRendersHeightHashPairs(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx_search": loadFixture(t, "rpc", "tx_search.json"),
	})
	stdout, stderr, err := runCmd(t, "", srv.URL, "logs",
		"message.action='/heimdallv2.topup.MsgTopupTx'",
		"--limit", "2", "--page", "1")
	if err != nil {
		t.Fatalf("logs: %v", err)
	}
	// Expect at least one "<height>  0x<hash>" line.
	found := false
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "0x") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("no matching line in logs output: %q", stdout)
	}
	mustContain(t, stderr, "total_count")
}

func TestLogsJSONPassthrough(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"tx_search": loadFixture(t, "rpc", "tx_search.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "logs",
		"message.action='/heimdallv2.topup.MsgTopupTx'",
		"--json")
	if err != nil {
		t.Fatalf("logs --json: %v", err)
	}
	var v any
	if uerr := json.Unmarshal([]byte(stdout), &v); uerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", uerr, stdout)
	}
}

// --- nonce command ---

func TestNonceBareOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]byte{
		"/cosmos/auth/v1beta1/accounts/0x02f615e95563ef16f10354dba9e584e58d2d4314": loadFixture(t, "rest", "cosmos_auth_account.json"),
	})
	stdout, _, err := runCmd(t, srv.URL, "", "nonce", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("nonce: %v", err)
	}
	if strings.TrimSpace(stdout) != "51129" {
		t.Errorf("nonce = %q, want 51129", stdout)
	}
}

func TestSequenceIsAliasOfNonce(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]byte{
		"/cosmos/auth/v1beta1/accounts/0x02f615e95563ef16f10354dba9e584e58d2d4314": loadFixture(t, "rest", "cosmos_auth_account.json"),
	})
	stdout, _, err := runCmd(t, srv.URL, "", "sequence", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("sequence: %v", err)
	}
	if strings.TrimSpace(stdout) != "51129" {
		t.Errorf("sequence = %q, want 51129", stdout)
	}
}

// --- balance command ---

func TestBalanceRawOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]byte{
		"/cosmos/bank/v1beta1/balances/0x02f615e95563ef16f10354dba9e584e58d2d4314/by_denom": loadFixture(t, "rest", "cosmos_bank_balance_pol.json"),
	})
	stdout, _, err := runCmd(t, srv.URL, "", "balance", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if strings.TrimSpace(stdout) != "7779000000000000000" {
		t.Errorf("balance = %q", stdout)
	}
}

func TestBalanceHumanOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]byte{
		"/cosmos/bank/v1beta1/balances/0x02f615e95563ef16f10354dba9e584e58d2d4314/by_denom": loadFixture(t, "rest", "cosmos_bank_balance_pol.json"),
	})
	stdout, _, err := runCmd(t, srv.URL, "", "balance", "0x02f615e95563ef16f10354dba9e584e58d2d4314", "--human")
	if err != nil {
		t.Fatalf("balance --human: %v", err)
	}
	// 7779000000000000000 with 18 decimals = 7.779
	mustContain(t, stdout, "7.779 pol")
}

func TestFormatDecimal(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"1000000000000000000", "1"},
		{"500000000000000000", "0.5"},
		{"0", "0"},
		{"1", "0.000000000000000001"},
		{"7779000000000000000", "7.779"},
	}
	for _, c := range cases {
		got, err := formatDecimal(c.in, 18)
		if err != nil {
			t.Errorf("formatDecimal(%s) err=%v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("formatDecimal(%s) = %s, want %s", c.in, got, c.want)
		}
	}
}

// --- rpc command ---

func TestRPCPassthrough(t *testing.T) {
	srv := newRPCFixtureServer(t, map[string][]byte{
		"status": loadFixture(t, "rpc", "status.json"),
	})
	stdout, _, err := runCmd(t, "", srv.URL, "rpc", "status")
	if err != nil {
		t.Fatalf("rpc status: %v", err)
	}
	var v any
	if uerr := json.Unmarshal([]byte(stdout), &v); uerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", uerr, stdout)
	}
	mustContain(t, stdout, "heimdallv2-80002")
}

func TestParseRPCArgs(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{"empty", nil, nil, false},
		{"string value", []string{"hash=abc"}, map[string]any{"hash": "abc"}, false},
		{"numeric JSON", []string{"height=42"}, map[string]any{"height": float64(42)}, false},
		{"bool JSON", []string{"prove=true"}, map[string]any{"prove": true}, false},
		{"null JSON", []string{"height=null"}, map[string]any{"height": nil}, false},
		{"invalid shape", []string{"no-equals-sign"}, nil, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseRPCArgs(c.args)
			if (err != nil) != c.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, c.wantErr)
			}
			if err != nil {
				return
			}
			if len(got) != len(c.want) {
				t.Fatalf("got=%v want=%v", got, c.want)
			}
			for k, v := range c.want {
				if got[k] != v {
					t.Errorf("key %q: got=%v want=%v", k, got[k], v)
				}
			}
		})
	}
}

// --- publish command ---

func TestPublishRequiresYes(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]byte{})
	stdout, _, err := runCmd(t, srv.URL, "", "publish", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error without --yes")
	}
	var uErr *client.UsageError
	if !errorsAs(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
	mustContain(t, stdout, "would broadcast")
}

func TestPublishHappyPath(t *testing.T) {
	// Respond with a minimal TxResponse envelope.
	resp := []byte(`{"tx_response":{"txhash":"ABCD","code":0,"raw_log":""}}`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/cosmos/tx/v1beta1/txs" {
			http.Error(w, "bad route", 404)
			return
		}
		var body broadcastRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Mode == "" || body.TxBytes == "" {
			http.Error(w, "missing fields", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	}))
	defer srv.Close()
	stdout, _, err := runCmd(t, srv.URL, "", "publish", "0xdeadbeef", "--yes")
	if err != nil {
		t.Fatalf("publish --yes: %v", err)
	}
	mustContain(t, stdout, "ABCD")
}

func TestNormalizeTxBytes(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"0x-hex", "0xdeadbeef", "3q2+7w==", false},
		{"plain hex", "deadbeef", "3q2+7w==", false},
		{"std base64", "3q2+7w==", "3q2+7w==", false},
		{"raw base64", "3q2+7w", "3q2+7w==", false},
		{"empty", "", "", true},
		{"bogus", "!!!not-valid!!!", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := normalizeTxBytes(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, c.wantErr)
			}
			if c.wantErr {
				return
			}
			if got != c.want {
				t.Errorf("got=%q want=%q", got, c.want)
			}
		})
	}
}
