package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// defaultFieldsUsage is the usage string for the --field flag shared by
// every query command.
const defaultFieldsUsage = "pluck one or more fields (repeatable)"

// Get describes a REST GET command: fetch one endpoint, decode the
// JSON-object body, and render it as JSON or KV. The zero hooks cover
// the plain case; the optional ones absorb the per-command quirks.
type Get struct {
	Use     string
	Short   string
	Aliases []string
	// Args defaults to cobra.NoArgs.
	Args cobra.PositionalArgs
	// Label names the response in decode/gRPC errors, e.g. "milestone params".
	Label string
	// Path is the fixed REST path. Mutually exclusive with Build.
	Path string
	// Build computes the path and query from the CLI args.
	Build func(cmd *cobra.Command, args []string) (string, url.Values, error)
	// L1Hint surfaces render.HintL1NotConfigured when the call fails
	// with gRPC code 13 / a transport-level dial error, and checks the
	// 2xx body for a gRPC error envelope before rendering.
	L1Hint bool
	// FieldsUsage overrides the --field usage string.
	FieldsUsage string
	// Flags registers extra command flags.
	Flags func(fs *pflag.FlagSet)
	// Opts post-processes the render options (e.g. a --base64 sugar flag).
	Opts func(cmd *cobra.Command, opts *render.Options)
	// Render renders the decoded map when --json is NOT set. nil falls
	// back to RenderKV, after unwrapping UnwrapKey if present.
	Render func(cmd *cobra.Command, m map[string]any, opts render.Options) error
	// RenderBody takes over non-JSON rendering before map decoding
	// (e.g. printing one bare field from a typed struct).
	RenderBody func(cmd *cobra.Command, body []byte, opts render.Options) error
	// UnwrapKey names a single-key envelope to unwrap for KV output.
	UnwrapKey string
}

// NewGetCmd builds the cobra command for a Get spec.
func (p *Pkg) NewGetCmd(spec Get) *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:     spec.Use,
		Short:   spec.Short,
		Aliases: spec.Aliases,
		Args:    orNoArgs(spec.Args),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, query := spec.Path, url.Values(nil)
			if spec.Build != nil {
				var err error
				if path, query, err = spec.Build(cmd, args); err != nil {
					return err
				}
			}
			rest, cfg, err := p.RESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), path, query)
			opts := RenderOpts(cmd, cfg, fields)
			if spec.Opts != nil {
				spec.Opts(cmd, &opts)
			}
			if err != nil {
				if spec.L1Hint && IsL1Unreachable(body, err) {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return err
			}
			if status == 0 && body == nil {
				return nil // --curl mode: the command was printed, not executed
			}
			if spec.L1Hint {
				var gerr GRPCErrorBody
				if jerr := json.Unmarshal(body, &gerr); jerr == nil && gerr.Code != 0 {
					if gerr.Code == GRPCCodeUnavailable {
						_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
					}
					return fmt.Errorf("%s failed: code=%d %s", spec.Label, gerr.Code, gerr.Message)
				}
			}
			if !opts.JSON && spec.RenderBody != nil {
				return spec.RenderBody(cmd, body, opts)
			}
			m, err := DecodeJSONMap(body, spec.Label)
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			if spec.Render != nil {
				return spec.Render(cmd, m, opts)
			}
			if spec.UnwrapKey != "" {
				if inner, ok := m[spec.UnwrapKey].(map[string]any); ok {
					return render.RenderKV(cmd.OutOrStdout(), inner, opts)
				}
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	registerCommon(cmd, &fields, spec.FieldsUsage, spec.Flags)
	return cmd
}

// RPC describes a CometBFT JSON-RPC command: call one method with
// empty params, pass the raw result through for --json, and hand it to
// Render otherwise.
type RPC struct {
	Use     string
	Short   string
	Aliases []string
	// Args defaults to cobra.NoArgs.
	Args cobra.PositionalArgs
	// Method is the CometBFT RPC method, e.g. "abci_info".
	Method string
	// FieldsUsage overrides the --field usage string.
	FieldsUsage string
	// Flags registers extra command flags.
	Flags func(fs *pflag.FlagSet)
	// Render renders the raw RPC result when --json is NOT set.
	Render func(cmd *cobra.Command, raw json.RawMessage, opts render.Options) error
}

// NewRPCCmd builds the cobra command for an RPC spec.
func (p *Pkg) NewRPCCmd(spec RPC) *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:     spec.Use,
		Short:   spec.Short,
		Aliases: spec.Aliases,
		Args:    orNoArgs(spec.Args),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := p.RPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := CallEmpty(cmd.Context(), rpc, spec.Method)
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}
			opts := RenderOpts(cmd, cfg, fields)
			if opts.JSON {
				generic, derr := DecodeGeneric(raw)
				if derr != nil {
					return derr
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			return spec.Render(cmd, raw, opts)
		},
	}
	registerCommon(cmd, &fields, spec.FieldsUsage, spec.Flags)
	return cmd
}

func orNoArgs(args cobra.PositionalArgs) cobra.PositionalArgs {
	if args == nil {
		return cobra.NoArgs
	}
	return args
}

func registerCommon(cmd *cobra.Command, fields *[]string, fieldsUsage string, extra func(fs *pflag.FlagSet)) {
	if fieldsUsage == "" {
		fieldsUsage = defaultFieldsUsage
	}
	f := cmd.Flags()
	f.StringArrayVarP(fields, "field", "f", nil, fieldsUsage)
	if extra != nil {
		extra(f)
	}
}
