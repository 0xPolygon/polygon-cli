// Package topup implements the `polycli heimdall topup` umbrella
// command and its subcommands targeting Heimdall v2's `x/topup`
// module: root, account, proof, verify, sequence, is-old.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.6 these endpoints live under a
// single umbrella rather than at the top level of the heimdall tree.
//
// Endpoints confirmed from heimdall-v2
// proto/heimdallv2/topup/query.proto:
//
//   - GET /topup/dividend-account-root
//   - GET /topup/dividend-account/{address}
//   - GET /topup/account-proof/{address}
//   - GET /topup/account-proof/{address}/verify?proof=…
//   - GET /topup/sequence?tx_hash=…&log_index=…
//   - GET /topup/is-old-tx?tx_hash=…&log_index=…
//
// The `proof`, `sequence`, and `is-old` endpoints fan out to L1 on the
// server side; a Heimdall node without `eth_rpc_url` returns gRPC code
// 13, which we surface as an L1-not-configured hint on stderr before
// propagating the error.
package topup

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// flags is injected by the caller via Register.
var flags *config.Flags

// TopupCmd is the umbrella `topup` command. Subcommands are attached
// by Register.
var TopupCmd = &cobra.Command{
	Use:   "topup",
	Short: "Query topup (dividend account) module endpoints.",
	Long:  usage,
	Args:  cobra.NoArgs,
}

// Register attaches the topup umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	TopupCmd.AddCommand(
		newRootCmd(),
		newAccountCmd(),
		newProofCmd(),
		newVerifyCmd(),
		newSequenceCmd(),
		newIsOldCmd(),
	)
	render.EnableWatchTree(TopupCmd)
	parent.AddCommand(TopupCmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "topup package not registered (flags unset)"}
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

// decodeJSONMap decodes raw into a map[string]any. Used by REST
// responses whose top-level shape is always an object.
func decodeJSONMap(raw []byte, label string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("decoding %s: %w (body=%q)", label, err, truncate(raw, 256))
	}
	return m, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}

// normalizeAddress accepts an Ethereum address with or without the
// `0x` prefix and returns the lower-case, `0x`-prefixed form (20 bytes
// / 40 hex chars). The topup REST endpoints expect a hex address in
// the URL path.
func normalizeAddress(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 40 {
		return "", &client.UsageError{Msg: fmt.Sprintf("address must be 20 bytes (40 hex chars), got %d", len(s))}
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return "", &client.UsageError{Msg: fmt.Sprintf("invalid address %q (non-hex character %q)", raw, r)}
		}
	}
	return "0x" + strings.ToLower(s), nil
}

// normalizeTxHash accepts a tx hash with or without the `0x` prefix and
// returns the lower-case, `0x`-prefixed form.
func normalizeTxHash(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 64 {
		return "", &client.UsageError{Msg: fmt.Sprintf("tx hash must be 32 bytes (64 hex chars), got %d", len(s))}
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return "", &client.UsageError{Msg: fmt.Sprintf("invalid tx hash %q (non-hex character %q)", raw, r)}
		}
	}
	return "0x" + strings.ToLower(s), nil
}

// normalizeHexBytes accepts a hex string with or without the `0x`
// prefix and returns the lower-case form WITHOUT the prefix (for use
// as a bare query param). Empty input is an error.
func normalizeHexBytes(raw, label string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if s == "" {
		return "", &client.UsageError{Msg: fmt.Sprintf("%s must not be empty", label)}
	}
	if len(s)%2 != 0 {
		return "", &client.UsageError{Msg: fmt.Sprintf("%s must have an even number of hex chars, got %d", label, len(s))}
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return "", &client.UsageError{Msg: fmt.Sprintf("invalid %s %q (non-hex character %q)", label, raw, r)}
		}
	}
	return strings.ToLower(s), nil
}

// gRPCErrorBody is the standard gRPC-gateway error envelope returned on
// 4xx/5xx from Heimdall REST. Only `code` and `message` are used here.
type gRPCErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// gRPCCodeUnavailable is the L1-unreachable code surfaced by topup
// endpoints (`/topup/account-proof/{address}`, `/topup/sequence`,
// `/topup/is-old-tx`) when the Heimdall node lacks `eth_rpc_url`.
const gRPCCodeUnavailable = 13

// isL1Unreachable inspects a REST body / error pair and returns true
// if the response looks like "gRPC code 13 because L1 RPC isn't
// configured on this Heimdall". Shape mirrors clerk.isL1Unreachable —
// topup repeats the logic locally rather than cross-import, keeping
// command packages independent.
func isL1Unreachable(body []byte, err error) bool {
	var hErr *client.HTTPError
	if errors.As(err, &hErr) && len(hErr.Body) > 0 {
		body = hErr.Body
	}
	if len(body) == 0 {
		if err != nil {
			msg := err.Error()
			return strings.Contains(msg, "connection refused") ||
				strings.Contains(msg, "dial tcp")
		}
		return false
	}
	var g gRPCErrorBody
	if jerr := json.Unmarshal(body, &g); jerr == nil && g.Code == gRPCCodeUnavailable {
		return true
	}
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "connection refused") ||
			strings.Contains(msg, "dial tcp") {
			return true
		}
	}
	return false
}
