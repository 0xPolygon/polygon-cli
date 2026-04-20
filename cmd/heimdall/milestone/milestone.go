// Package milestone implements the `polycli heimdall milestone`
// umbrella command (alias `ms`) and its subcommands targeting Heimdall
// v2's `x/milestone` module: params, count, latest, get.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.3 these endpoints live under a
// single umbrella rather than at the top level. The umbrella also
// accepts a bare integer (`milestone 11602043`) as a shorthand for
// `milestone get 11602043`.
package milestone

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// flags is injected by the caller via Register.
var flags *config.Flags

// MilestoneCmd is the umbrella `milestone` command. Subcommands are
// attached by Register.
var MilestoneCmd = &cobra.Command{
	Use:     "milestone [NUMBER]",
	Aliases: []string{"ms"},
	Short:   "Query milestone module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-number shorthand: `milestone 11602043` forwards to
	// `milestone get 11602043`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown milestone subcommand or number %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// Register attaches the milestone umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	MilestoneCmd.AddCommand(
		newParamsCmd(),
		newCountCmd(),
		newLatestCmd(),
		newGetCmd(),
	)
	parent.AddCommand(MilestoneCmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "milestone package not registered (flags unset)"}
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
