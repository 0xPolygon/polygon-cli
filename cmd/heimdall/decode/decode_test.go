package decode

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// buildRoot wires the decode umbrella under a throwaway root, matching
// how cmd/heimdall/heimdall.go composes it at startup.
func buildRoot(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true, SilenceErrors: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	return root, buf
}

func runCmd(t *testing.T, root *cobra.Command, args []string) (string, error) {
	t.Helper()
	buf := root.OutOrStderr().(*bytes.Buffer)
	buf.Reset()
	root.SetArgs(args)
	err := root.ExecuteContext(context.Background())
	return buf.String(), err
}

// TestDecodeMsgRoundTrip builds a MsgWithdrawFeeTx, base64-encodes its
// value, and decodes it through the CLI. The output must include the
// proposer back in the JSON body.
func TestDecodeMsgRoundTrip(t *testing.T) {
	msg := &hproto.MsgWithdrawFeeTx{Proposer: "0xabc", Amount: "1000"}
	b64 := base64.StdEncoding.EncodeToString(msg.Marshal())
	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{
		"decode", "msg", hproto.MsgWithdrawFeeTxTypeURL, b64,
	})
	if err != nil {
		t.Fatalf("decode msg: %v\n%s", err, out)
	}
	if !strings.Contains(out, "0xabc") {
		t.Errorf("expected proposer in output, got %q", out)
	}
	if !strings.Contains(out, "1000") {
		t.Errorf("expected amount in output, got %q", out)
	}
}

// TestDecodeMsgList prints the registered type URLs and must include
// the starter MsgWithdrawFeeTx.
func TestDecodeMsgList(t *testing.T) {
	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{"decode", "msg", "--list"})
	if err != nil {
		t.Fatalf("decode msg --list: %v\n%s", err, out)
	}
	if !strings.Contains(out, hproto.MsgWithdrawFeeTxTypeURL) {
		t.Errorf("expected %s in --list output, got %q", hproto.MsgWithdrawFeeTxTypeURL, out)
	}
	if !strings.Contains(out, hproto.MsgCheckpointTypeURL) {
		t.Errorf("expected %s in --list output", hproto.MsgCheckpointTypeURL)
	}
}

// TestDecodeMsgUnknownTypeURL asserts the command fails cleanly for a
// type URL that is not in the registry.
func TestDecodeMsgUnknownTypeURL(t *testing.T) {
	root, _ := buildRoot(t)
	_, err := runCmd(t, root, []string{
		"decode", "msg", "/not.a.real.type.URL", base64.StdEncoding.EncodeToString([]byte{1, 2, 3}),
	})
	if err == nil {
		t.Fatal("expected error for unknown type URL")
	}
	if !strings.Contains(err.Error(), "unknown type URL") {
		t.Errorf("expected 'unknown type URL' in error, got %v", err)
	}
}

// TestDecodeHashTx computes SHA256 over a raw input and asserts the
// output matches the known digest.
func TestDecodeHashTx(t *testing.T) {
	// A trivial 3-byte payload; the SHA256 can be hand-verified.
	payload := []byte{0x01, 0x02, 0x03}
	want := "039058C6F2C0CB492C533B0A4D14EF77CC0F78ABCCCED5287D84A1A2011CFB81"
	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{"decode", "hash-tx", "0x010203"})
	if err != nil {
		t.Fatalf("decode hash-tx: %v\n%s", err, out)
	}
	// Strip trailing newline, 0x prefix.
	got := strings.TrimSpace(out)
	got = strings.TrimPrefix(got, "0x")
	if !strings.EqualFold(got, want) {
		t.Errorf("hash mismatch\n got=%s\nwant=%s", got, want)
	}
	_ = payload
}

// TestDecodeTxRoundTrip constructs a TxRaw containing a MsgWithdrawFeeTx,
// base64-encodes it, decodes via `decode tx`, and asserts the proposer
// survives.
func TestDecodeTxRoundTrip(t *testing.T) {
	body := &hproto.TxBody{
		Messages: []*hproto.Any{
			(&hproto.MsgWithdrawFeeTx{Proposer: "0xdeadbeef", Amount: "1"}).AsAny(),
		},
	}
	auth := &hproto.AuthInfo{Fee: &hproto.Fee{GasLimit: 200000, Amount: []hproto.Coin{{Denom: "pol", Amount: "100"}}}}
	raw := &hproto.TxRaw{
		BodyBytes:     body.Marshal(),
		AuthInfoBytes: auth.Marshal(),
		Signatures:    [][]byte{{0xde, 0xad}},
	}
	b64 := base64.StdEncoding.EncodeToString(raw.Marshal())

	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{"decode", "tx", b64})
	if err != nil {
		t.Fatalf("decode tx: %v\n%s", err, out)
	}
	if !strings.Contains(out, "0xdeadbeef") {
		t.Errorf("expected proposer in output, got %q", out)
	}
	if !strings.Contains(out, hproto.MsgWithdrawFeeTxTypeURL) {
		t.Errorf("expected type URL in output, got %q", out)
	}
	if !strings.Contains(out, "fee.gas_limit=200000") {
		t.Errorf("expected fee.gas_limit in output, got %q", out)
	}
}

// TestDecodeTxJSON checks that --json emits a single JSON record.
func TestDecodeTxJSON(t *testing.T) {
	body := &hproto.TxBody{Messages: []*hproto.Any{
		(&hproto.MsgWithdrawFeeTx{Proposer: "0x1", Amount: "1"}).AsAny(),
	}}
	raw := &hproto.TxRaw{BodyBytes: body.Marshal(), AuthInfoBytes: (&hproto.AuthInfo{}).Marshal()}
	b64 := base64.StdEncoding.EncodeToString(raw.Marshal())
	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{"decode", "tx", "--json", b64})
	if err != nil {
		t.Fatalf("decode tx --json: %v\n%s", err, out)
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if _, ok := v["tx_hash_sha256"]; !ok {
		t.Errorf("missing tx_hash_sha256: %v", v)
	}
}

// TestDecodeVERoundTrip builds a VoteExtension, marshals it, and decodes
// via the CLI.
func TestDecodeVERoundTrip(t *testing.T) {
	ve := &hproto.VoteExtension{
		BlockHash: []byte{0xaa, 0xbb},
		Height:    42,
		SideTxResponses: []hproto.SideTxResponse{
			{TxHash: []byte{0xde, 0xad}, Result: hproto.VoteYes},
		},
	}
	raw := ve.Marshal()
	h := hex.EncodeToString(raw)

	root, _ := buildRoot(t)
	out, err := runCmd(t, root, []string{"decode", "ve", h})
	if err != nil {
		t.Fatalf("decode ve: %v\n%s", err, out)
	}
	if !strings.Contains(out, "0xaabb") {
		t.Errorf("expected block hash hex in output, got %q", out)
	}
	if !strings.Contains(out, "VOTE_YES") {
		t.Errorf("expected VOTE_YES in output, got %q", out)
	}
	if !strings.Contains(out, `"height": 42`) {
		t.Errorf("expected height=42 in output, got %q", out)
	}
}

// TestDecodeInputAcceptsBase64AndHex sanity-checks the shared decoder.
func TestDecodeInputAcceptsBase64AndHex(t *testing.T) {
	want := []byte{0xde, 0xad, 0xbe, 0xef}
	// 0x-hex
	got, err := decodeInput("x", "0xdeadbeef")
	if err != nil || !bytes.Equal(got, want) {
		t.Errorf("0x-hex: %v %v", got, err)
	}
	// bare hex
	got, err = decodeInput("x", "deadbeef")
	if err != nil || !bytes.Equal(got, want) {
		t.Errorf("bare hex: %v %v", got, err)
	}
	// base64
	got, err = decodeInput("x", base64.StdEncoding.EncodeToString(want))
	if err != nil || !bytes.Equal(got, want) {
		t.Errorf("base64: %v %v", got, err)
	}
}
