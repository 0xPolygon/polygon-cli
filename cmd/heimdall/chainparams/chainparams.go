// Package chainparams implements the `polycli heimdall chainmanager`
// umbrella command (alias `cm`) and its subcommands targeting Heimdall
// v2's `x/chainmanager` module.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.7 the chainmanager module holds
// the L1/L2 chain ids, tx confirmation depths, and L1 contract
// addresses. Upstream exposes a single HTTP route
// (`/chainmanager/params`, confirmed in
// heimdall-v2/proto/heimdallv2/chainmanager/query.proto); the
// `addresses` subcommand is a derived view over the same response.
//
// Package directory is named `chainparams` (not `chain`) because the
// top-level `chain` command is already claimed by the CometBFT-facing
// cast-like commands in cmd/heimdall/chain.
package chainparams

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// flags is injected by the caller via Register.
var flags *config.Flags

// Cmd is the umbrella `chainmanager` command (alias `cm`).
// Subcommands are attached by Register.
var Cmd = &cobra.Command{
	Use:     "chainmanager",
	Aliases: []string{"cm"},
	Short:   "Query chainmanager module endpoints.",
	Long:    usage,
	Args:    cobra.NoArgs,
}

// Register attaches the chainmanager umbrella command and its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	Cmd.AddCommand(
		newParamsCmd(),
		newAddressesCmd(),
	)
	parent.AddCommand(Cmd)
}

// newRESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func newRESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	if flags == nil {
		return nil, nil, &client.UsageError{Msg: "chainmanager package not registered (flags unset)"}
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

// decodeJSONMap decodes raw into a map[string]any.
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
