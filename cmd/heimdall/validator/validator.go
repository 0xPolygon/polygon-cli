// Package validator implements the `polycli heimdall validator`
// umbrella command (alias `val`) and its subcommands targeting Heimdall
// v2's `x/stake` module: set/validators, total-power, get, signer,
// status, proposer, proposers, is-old-stake-tx.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.4 these endpoints live under a
// single umbrella, and the umbrella also accepts a bare integer
// (`validator 4`) as a shorthand for `validator get 4`. The top-level
// `validators` command is registered separately as an alias for
// `validator set`.
package validator

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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

// ValidatorCmd is the umbrella `validator` command. Subcommands are
// attached by Register.
var ValidatorCmd = &cobra.Command{
	Use:     "validator [ID]",
	Aliases: []string{"val"},
	Short:   "Query stake module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `validator 4` forwards to `validator get 4`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown validator subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// ValidatorsCmd is the top-level `validators` alias for `validator set`.
// It is attached to the root heimdall command alongside ValidatorCmd so
// operators can type either form.
var ValidatorsCmd = &cobra.Command{
	Use:   "validators",
	Short: "Alias for `validator set`.",
	Args:  cobra.NoArgs,
	RunE:  runSet,
}

// setFlags keeps the shared flag state for `validator set` /
// `validators`. Both commands share RunE (runSet) and must therefore
// read from the same variables.
var setFlags = struct {
	sort   string
	limit  int
	fields []string
}{}

// Register attaches the validator umbrella command (and the top-level
// `validators` alias) to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	ValidatorCmd.AddCommand(
		newSetCmd(),
		newTotalPowerCmd(),
		newGetCmd(),
		newSignerCmd(),
		newStatusCmd(),
		newProposerCmd(),
		newProposersCmd(),
		newIsOldStakeTxCmd(),
	)
	// Attach shared flags to the top-level `validators` alias as well.
	attachSetFlags(ValidatorsCmd.Flags())
	// Read-only umbrella: wire `--watch` into every descendant plus
	// the top-level alias.
	render.EnableWatchTree(ValidatorCmd)
	render.EnableWatch(ValidatorsCmd)
	parent.AddCommand(ValidatorCmd)
	parent.AddCommand(ValidatorsCmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "validator package not registered (flags unset)"}
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

// normalizeSignerAddress accepts a signer hex with or without the `0x`
// prefix and returns the lower-case, `0x`-prefixed form consumed by
// /stake/signer/{addr}. Returns a UsageError for non-hex or wrong-length
// inputs.
func normalizeSignerAddress(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 40 {
		return "", &client.UsageError{Msg: fmt.Sprintf("signer must be 20 bytes (40 hex chars), got %d", len(s))}
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return "", &client.UsageError{Msg: fmt.Sprintf("invalid signer %q (non-hex character %q)", raw, r)}
		}
	}
	return "0x" + strings.ToLower(s), nil
}

// normalizeTxHash accepts a tx hash with or without the `0x` prefix and
// returns the lower-case, `0x`-prefixed form. The REST stake endpoints
// expect the prefix and will 500 without it, so we re-add it
// unconditionally.
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

// gRPCErrorBody is the standard gRPC-gateway error envelope returned on
// 4xx/5xx from Heimdall REST. Only `code` and `message` are used here.
type gRPCErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// gRPCCodeUnavailable is the L1-unreachable code surfaced by
// /stake/is-old-tx when the node lacks `eth_rpc_url`.
const gRPCCodeUnavailable = 13

// isL1Unreachable inspects a REST body / error pair and returns true if
// the response looks like "gRPC code 13 because L1 RPC isn't configured
// on this Heimdall". The body may come either from a successful 2xx
// response that still carries a gRPC-error envelope, or from an
// HTTPError (4xx/5xx).
func isL1Unreachable(body []byte, err error) bool {
	var hErr *client.HTTPError
	if errors.As(err, &hErr) && len(hErr.Body) > 0 {
		body = hErr.Body
	}
	if len(body) == 0 {
		// Fall through to error-string inspection: the transport layer
		// surfaces "dial tcp" / "connection refused" when the REST
		// gateway itself can't reach L1.
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
