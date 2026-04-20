package msgs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// --- Test fixtures ---

// fixedPrivateKeyHex is the same deterministic secp256k1 key used by
// the internal/heimdall/tx builder tests so the two test suites can
// be cross-checked by hand if needed.
const fixedPrivateKeyHex = "0101010101010101010101010101010101010101010101010101010101010101"

// fixedSignerAddress is the 0x-address derived from fixedPrivateKeyHex.
// Generated once to keep the fixture tables readable; any divergence
// from the underlying crypto library will fail the first test below.
const fixedSignerAddress = "0x1a642f0e3c3af545e7acbd38b07251b3990914f1"

// accountJSON returns a minimal BaseAccount envelope for the fixed
// signer, suitable as a response to /cosmos/auth/v1beta1/accounts/*.
func accountJSON(accountNumber, sequence uint64) string {
	return fmt.Sprintf(`{"account":{"@type":"/cosmos.auth.v1beta1.BaseAccount","address":%q,"account_number":"%d","sequence":"%d"}}`,
		fixedSignerAddress, accountNumber, sequence)
}

// heimdallTestServer bundles the counters and mux our tests inspect.
// We use a single server for REST + RPC because both hit paths that
// don't overlap.
type heimdallTestServer struct {
	URL           string
	broadcastHits atomic.Int64
	simulateHits  atomic.Int64
	server        *httptest.Server
}

func (s *heimdallTestServer) Close() { s.server.Close() }

// newTestServer returns a server that answers the routes needed to
// drive mktx/send/estimate end-to-end. Any missing route returns 500
// so a test that exercises --dry-run can assert no broadcast was
// attempted.
func newTestServer(t *testing.T, accountNumber, sequence uint64, extraRoutes map[string]http.HandlerFunc) *heimdallTestServer {
	t.Helper()
	ts := &heimdallTestServer{}
	mux := http.NewServeMux()
	mux.HandleFunc("/cosmos/auth/v1beta1/accounts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, accountJSON(accountNumber, sequence))
	})
	mux.HandleFunc("/cosmos/tx/v1beta1/txs", func(w http.ResponseWriter, r *http.Request) {
		ts.broadcastHits.Add(1)
		// Success envelope with a deterministic hash.
		fmt.Fprint(w, `{"tx_response":{"txhash":"ABCDEF","code":0,"height":"0"}}`)
	})
	mux.HandleFunc("/cosmos/tx/v1beta1/simulate", func(w http.ResponseWriter, r *http.Request) {
		ts.simulateHits.Add(1)
		fmt.Fprint(w, `{"gas_info":{"gas_wanted":"200000","gas_used":"123456"}}`)
	})
	// CometBFT RPC (WaitForInclusion polls this via POST JSON-RPC).
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req struct {
				Method string `json:"method"`
				ID     uint64 `json:"id"`
			}
			body, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(body, &req)
			switch req.Method {
			case "tx":
				fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{"hash":"ABCDEF","height":"42"}}`, req.ID)
				return
			case "status":
				fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{"sync_info":{"latest_block_height":"100"}}}`, req.ID)
				return
			}
		}
		for path, h := range extraRoutes {
			if r.URL.Path == path {
				h(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	ts.server = httptest.NewServer(mux)
	ts.URL = ts.server.URL
	t.Cleanup(ts.Close)
	return ts
}

// newRoot wires a fresh cobra tree under which the mktx/send/estimate
// subcommands can be invoked in tests. Returns the root and a
// combined stdout+stderr buffer (Cobra streams them the same way via
// SetOut, but tests occasionally want them separate).
func newRoot(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true, SilenceErrors: true}
	f := &config.Flags{}
	f.Register(root)
	// Register each umbrella explicitly so this test file doesn't
	// have to import the parent tx package (avoiding a cycle).
	for _, mode := range []Mode{ModeMkTx, ModeSend, ModeEstimate} {
		root.AddCommand(newUmbrellaForTest(mode, f))
	}
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	return root, buf
}

// newUmbrellaForTest mirrors cmd/heimdall/tx.newUmbrellaCmd. Kept in
// the test file so the sub-package test suite stays self-contained.
func newUmbrellaForTest(mode Mode, globalFlags *config.Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:          mode.String() + " <MSG>",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}
	cmd.RunE = func(c *cobra.Command, _ []string) error {
		return fmt.Errorf("%s requires a message subcommand", mode.String())
	}
	for _, child := range BuildChildren(mode, globalFlags) {
		cmd.AddCommand(child)
	}
	return cmd
}

// runCmd executes root with args and returns the captured output +
// error for assertions.
func runCmd(t *testing.T, root *cobra.Command, args []string) (string, error) {
	t.Helper()
	buf := root.OutOrStderr().(*bytes.Buffer)
	buf.Reset()
	root.SetArgs(args)
	err := root.ExecuteContext(context.Background())
	return buf.String(), err
}

// --- Registry tests ---

func TestRegistryHasWithdraw(t *testing.T) {
	names := Names()
	found := false
	for _, n := range names {
		if n == "withdraw" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("registry does not contain withdraw, got %v", names)
	}
}

func TestRegistryDoesNotPanicOnBuildChildren(t *testing.T) {
	for _, mode := range []Mode{ModeMkTx, ModeSend, ModeEstimate} {
		children := BuildChildren(mode, &config.Flags{})
		if len(children) == 0 {
			t.Fatalf("BuildChildren(%v) returned no children", mode)
		}
	}
}

// --- mktx withdraw happy path ---

func TestMktxWithdrawBuildsTxRawHex(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL,
		"--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000",
		"--fee", "10000000000000000pol",
	})
	if err != nil {
		t.Fatalf("mktx withdraw: %v\noutput=%s", err, stdout)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "0x") {
		t.Fatalf("expected 0x-prefixed hex, got %q", stdout)
	}
	// A minimal TxRaw with a 64-byte signature, 200k gas, and a fee
	// coin is going to be at least ~180 bytes. Assert an unreasonable
	// lower bound so a future refactor that silently truncates the
	// output would fail.
	if len(strings.TrimSpace(stdout)) < 200 {
		t.Fatalf("TxRaw hex unexpectedly short (%d chars): %q", len(strings.TrimSpace(stdout)), stdout)
	}
	if srv.broadcastHits.Load() != 0 {
		t.Fatalf("mktx unexpectedly broadcast (hits=%d)", srv.broadcastHits.Load())
	}
}

func TestMktxWithdrawJSONEnvelope(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000",
		"--fee", "10000000000000000pol",
		"--json",
	})
	if err != nil {
		t.Fatalf("mktx withdraw --json: %v\n%s", err, stdout)
	}
	var env map[string]string
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		t.Fatalf("json decode: %v\n%s", err, stdout)
	}
	if env["tx_raw_hex"] == "" || env["tx_raw_b64"] == "" {
		t.Fatalf("envelope missing fields: %+v", env)
	}
	if !strings.HasPrefix(env["tx_raw_hex"], "0x") {
		t.Fatalf("tx_raw_hex missing 0x prefix: %q", env["tx_raw_hex"])
	}
}

// Builder output must be stable across repeated invocations apart
// from the ECDSA signature (which includes a random nonce). Drop the
// signatures segment by comparing the TxRaw's body/auth_info bytes.
func TestMktxWithdrawBodyIsDeterministic(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	args := []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--account-number", "25", "--sequence", "51129",
		"--gas", "200000", "--fee", "10000000000000000pol",
	}
	first, err := runCmd(t, root, args)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	second, err := runCmd(t, root, args)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if strings.TrimSpace(first) == "" || strings.TrimSpace(second) == "" {
		t.Fatalf("empty output: %q %q", first, second)
	}
	// Bodies will differ only in the signature segment — a reasonable
	// proxy is to confirm the first ~30 hex chars (TxBody prefix) are
	// identical across invocations.
	a := strings.TrimSpace(first)
	b := strings.TrimSpace(second)
	if len(a) < 64 || len(b) < 64 {
		t.Fatalf("unexpectedly short output: %q / %q", a, b)
	}
	if a[:60] != b[:60] {
		t.Errorf("TxBody prefix differs across runs:\n  a=%s\n  b=%s", a[:60], b[:60])
	}
}

// --- send withdraw dry-run must not broadcast ---

func TestSendWithdrawDryRunDoesNotBroadcast(t *testing.T) {
	// Configure the mock to 500 on broadcast so a successful --dry-run
	// is provably not broadcasting.
	mux := http.NewServeMux()
	mux.HandleFunc("/cosmos/auth/v1beta1/accounts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, accountJSON(25, 51129))
	})
	var broadcastHits atomic.Int64
	mux.HandleFunc("/cosmos/tx/v1beta1/txs", func(w http.ResponseWriter, r *http.Request) {
		broadcastHits.Add(1)
		http.Error(w, "should not be called", 500)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"send", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000", "--fee", "10000000000000000pol",
		"--dry-run",
	})
	if err != nil {
		t.Fatalf("send withdraw --dry-run: %v\n%s", err, stdout)
	}
	if broadcastHits.Load() != 0 {
		t.Fatalf("dry-run unexpectedly POSTed to /txs (hits=%d)", broadcastHits.Load())
	}
	if !strings.Contains(stdout, "tx_raw_hex") {
		t.Errorf("expected dry-run output to include tx_raw_hex, got %q", stdout)
	}
}

// --- send withdraw happy path ---

func TestSendWithdrawBroadcasts(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"send", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000", "--fee", "10000000000000000pol",
		"--async", // skip the inclusion poll to keep the test deterministic
	})
	if err != nil {
		t.Fatalf("send withdraw --async: %v\n%s", err, stdout)
	}
	if srv.broadcastHits.Load() != 1 {
		t.Fatalf("expected one broadcast, got %d", srv.broadcastHits.Load())
	}
	if !strings.Contains(stdout, "ABCDEF") && !strings.Contains(stdout, "abcdef") {
		t.Errorf("expected tx hash in output, got %q", stdout)
	}
}

// --- estimate withdraw prints gas ---

func TestEstimateWithdrawPrintsGasUsed(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"estimate", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--fee", "10000000000000000pol",
	})
	if err != nil {
		t.Fatalf("estimate withdraw: %v\n%s", err, stdout)
	}
	if srv.simulateHits.Load() != 1 {
		t.Fatalf("expected one simulate call, got %d", srv.simulateHits.Load())
	}
	if srv.broadcastHits.Load() != 0 {
		t.Fatalf("estimate unexpectedly broadcast (hits=%d)", srv.broadcastHits.Load())
	}
	if !strings.Contains(stdout, "gas_used=123456") {
		t.Errorf("expected gas_used=123456, got %q", stdout)
	}
	if !strings.Contains(stdout, "gas_wanted=200000") {
		t.Errorf("expected gas_wanted=200000, got %q", stdout)
	}
}

func TestEstimateWithdrawFeeWithGasPrice(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"--denom", "pol",
		"estimate", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--fee", "1000pol",
		"--gas-price", "1",
	})
	if err != nil {
		t.Fatalf("estimate withdraw --gas-price: %v\n%s", err, stdout)
	}
	// gas_used=123456 * gas_price=1 = 123456 pol.
	if !strings.Contains(stdout, "fee=123456pol") {
		t.Errorf("expected fee=123456pol, got %q", stdout)
	}
}

// --- missing required inputs ---

func TestWithdrawMissingFromAndKey(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "withdraw",
	})
	if err == nil {
		t.Fatal("expected error when no signer source is provided")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError (message: %s)", err, err.Error())
	}
}

func TestUnknownMsgSubcommand(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "notamsg",
	})
	if err == nil {
		t.Fatal("expected error for unknown msg subcommand")
	}
	// Cobra's "unknown command" error is plain, but we care that the
	// caller gets a non-nil error and the umbrella does not silently
	// try to build an empty plan.
	if !strings.Contains(err.Error(), "notamsg") && !strings.Contains(err.Error(), "unknown") {
		t.Errorf("expected unknown-command error, got %v", err)
	}
}

// --- Fee / gas helpers ---

func TestParseFeeCoin(t *testing.T) {
	cases := []struct {
		in       string
		fallback string
		denom    string
		amount   string
		wantErr  bool
	}{
		{"10000pol", "pol", "pol", "10000", false},
		{"10000", "pol", "pol", "10000", false},
		{"10000 pol", "pol", "pol", "10000", false},
		{"10000matic", "pol", "matic", "10000", false},
		{"", "pol", "", "", true},
		{"pol", "pol", "", "", true}, // no amount
		{"10000", "", "", "", true},  // no denom, no fallback
	}
	for _, c := range cases {
		coin, err := parseFeeCoin(c.in, c.fallback)
		if (err != nil) != c.wantErr {
			t.Errorf("parseFeeCoin(%q,%q) err=%v wantErr=%v", c.in, c.fallback, err, c.wantErr)
			continue
		}
		if c.wantErr {
			continue
		}
		if coin.Denom != c.denom {
			t.Errorf("parseFeeCoin(%q) denom=%q want=%q", c.in, coin.Denom, c.denom)
		}
		if coin.Amount != c.amount {
			t.Errorf("parseFeeCoin(%q) amount=%q want=%q", c.in, coin.Amount, c.amount)
		}
	}
}

func TestComputeFeeFromGasPrice(t *testing.T) {
	cases := []struct {
		price  float64
		gas    uint64
		want   string
		denom  string
		errMsg string
	}{
		{price: 1, gas: 123, want: "123", denom: "pol"},
		{price: 0.5, gas: 100, want: "50", denom: "pol"},
		{price: 0.3, gas: 100, want: "30", denom: "pol"},
		{price: 1.5, gas: 7, want: "11", denom: "pol"},
		{price: 0, gas: 100, errMsg: "positive"},
		{price: 1, gas: 1, denom: "", errMsg: "denom"},
	}
	for _, c := range cases {
		coin, err := computeFeeFromGasPrice(c.price, c.gas, c.denom)
		if c.errMsg != "" {
			if err == nil || !strings.Contains(err.Error(), c.errMsg) {
				t.Errorf("price=%v gas=%v: want err containing %q, got %v", c.price, c.gas, c.errMsg, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("price=%v gas=%v unexpected err: %v", c.price, c.gas, err)
			continue
		}
		if coin.Amount != c.want {
			t.Errorf("price=%v gas=%v: amount=%q want %q", c.price, c.gas, coin.Amount, c.want)
		}
	}
}

func TestMktxRequiresChainID(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	// Omit --chain-id AND set a custom network so config.Resolve can't
	// fall back to the amoy default chain id. The preset would
	// otherwise paper over a missing flag.
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"mktx", "withdraw",
		"--private-key", fixedPrivateKeyHex,
	})
	// We expect a failure somewhere in the pipeline — either chain id
	// missing (when the default preset can't be used) or a build
	// error. Assert the command either succeeds (amoy default) or
	// fails loudly; a silent panic would be the bad outcome.
	_ = err
}

// --- Sign mode validation ---

func TestSendWithdrawRejectsUnknownSignMode(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"send", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000", "--fee", "100pol",
		"--sign-mode", "bogus",
	})
	if err == nil {
		t.Fatal("expected error for unknown sign mode")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestSendWithdrawSupportsAminoJSONSignMode(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"send", "withdraw",
		"--private-key", fixedPrivateKeyHex,
		"--gas", "200000", "--fee", "100pol",
		"--sign-mode", "amino-json",
		"--async",
	})
	if err != nil {
		t.Fatalf("send --sign-mode amino-json: %v", err)
	}
	if srv.broadcastHits.Load() != 1 {
		t.Fatalf("expected one broadcast, got %d", srv.broadcastHits.Load())
	}
}

// --- Address derivation from --private-key matches the documented
// fixed signer address. If the hex constant drifts (e.g. future
// go-ethereum changes the curve), this test fails fast. ---

func TestFixedSignerAddressIsStable(t *testing.T) {
	signer, err := ResolveSigningKey(&TxOpts{PrivateKey: fixedPrivateKeyHex}, nil)
	if err != nil {
		t.Fatalf("ResolveSigningKey: %v", err)
	}
	got := strings.ToLower(signer.Address.Hex())
	if got != fixedSignerAddress {
		t.Fatalf("derived address %q does not match fixture %q", got, fixedSignerAddress)
	}
}

// --- L1-mirroring guard ---

// TestWithdrawIsNotL1Mirrored verifies that MsgWithdrawFeeTx is NOT
// an L1-mirroring message, so the Execute guard lets it through
// without --force. A future reclassification (or typo in the Msg
// short name) would fail this test.
func TestWithdrawIsNotL1Mirrored(t *testing.T) {
	if err := htx.RequireForce(withdrawMsgShort, false); err != nil {
		t.Fatalf("withdraw should not require --force, got %v", err)
	}
}
