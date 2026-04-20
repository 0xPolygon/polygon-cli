// Package checkpoint implements the `polycli heimdall checkpoint`
// umbrella command (alias `cp`) and its subcommands: params, count,
// latest, get, buffer, last-no-ack, next, list, signatures, overview.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.1 these endpoints live under
// a single umbrella rather than at the top level of the heimdall
// tree. The umbrella also accepts a bare integer (`checkpoint 38871`)
// as a shorthand for `checkpoint get 38871`.
package checkpoint

import (
	_ "embed"
	"encoding/json"
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

// CheckpointCmd is the umbrella `checkpoint` command. Subcommands are
// attached by Register.
var CheckpointCmd = &cobra.Command{
	Use:     "checkpoint [ID]",
	Aliases: []string{"cp"},
	Short:   "Query checkpoint module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `checkpoint 38871` forwards to `checkpoint get 38871`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown checkpoint subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// Register attaches the checkpoint umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
//
// Every checkpoint subcommand is read-only, so we apply
// render.EnableWatchTree once here.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	CheckpointCmd.AddCommand(
		newParamsCmd(),
		newCountCmd(),
		newLatestCmd(),
		newGetCmd(),
		newBufferCmd(),
		newLastNoAckCmd(),
		newNextCmd(),
		newListCmd(),
		newSignaturesCmd(),
		newOverviewCmd(),
	)
	render.EnableWatchTree(CheckpointCmd)
	parent.AddCommand(CheckpointCmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "checkpoint package not registered (flags unset)"}
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

// normalizeCheckpointHash accepts a checkpoint tx hash with or without
// the `0x` prefix and returns the lower-case, unprefixed hex form
// expected by /checkpoints/signatures/{hash} on Heimdall. Returns a
// UsageError for non-hex or non-32-byte inputs.
func normalizeCheckpointHash(raw string) (string, error) {
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
	return strings.ToLower(s), nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
