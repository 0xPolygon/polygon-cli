// Package clerk implements the `polycli heimdall state-sync` umbrella
// command (aliases `clerk` and `ss`) and its subcommands targeting
// Heimdall v2's `x/clerk` module: count, latest-id, get, list, range,
// sequence, is-old.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.5 these endpoints live under a
// single umbrella rather than at the top level of the heimdall tree.
// The umbrella also accepts a bare integer (`state-sync 36610`) as a
// shorthand for `state-sync get 36610`.
//
// Pagination note: `/clerk/event-records/list` is page-based (page +
// limit query params), NOT Cosmos pagination, and rejects `page=0`
// with HTTP 400. `/clerk/time` is Cosmos-paginated (pagination.limit).
package clerk

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

// ClerkCmd is the umbrella `state-sync` command (aliases `clerk`,
// `ss`). Subcommands are attached by Register.
var ClerkCmd = &cobra.Command{
	Use:     "state-sync [ID]",
	Aliases: []string{"clerk", "ss"},
	Short:   "Query state-sync (clerk) module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `state-sync 36610` forwards to `state-sync get 36610`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown state-sync subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0], false)
	},
}

// Register attaches the state-sync umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	ClerkCmd.AddCommand(
		newCountCmd(),
		newLatestIDCmd(),
		newGetCmd(),
		newListCmd(),
		newRangeCmd(),
		newSequenceCmd(),
		newIsOldCmd(),
	)
	render.EnableWatchTree(ClerkCmd)
	parent.AddCommand(ClerkCmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "clerk package not registered (flags unset)"}
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
// honouring --json, --field, --color, --raw, and TTY detection. base64
// is a clerk-specific sugar for Raw (since the spec talks about `data`
// specifically rather than all byte-fields).
func renderOpts(cmd *cobra.Command, cfg *config.Config, fields []string, base64 bool) render.Options {
	return render.Options{
		JSON:   cfg.JSON,
		Raw:    cfg.Raw || base64,
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

// normalizeTxHash accepts a tx hash with or without the `0x` prefix and
// returns the lower-case, `0x`-prefixed form. The clerk REST endpoints
// expect a 0x-prefixed value.
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

// gRPCCodeUnavailable is the L1-unreachable code surfaced by clerk
// endpoints (`/clerk/event-records/latest-id`, `/clerk/sequence`,
// `/clerk/is-old-tx`) when the Heimdall node lacks `eth_rpc_url`.
const gRPCCodeUnavailable = 13

// isL1Unreachable inspects a REST body / error pair and returns true if
// the response looks like "gRPC code 13 because L1 RPC isn't configured
// on this Heimdall". Shape mirrors validator.isL1Unreachable — clerk
// repeats the logic locally rather than cross-import, keeping command
// packages independent.
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
