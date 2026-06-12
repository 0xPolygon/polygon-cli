// Package cmdutil holds the plumbing shared by every `polycli heimdall`
// command package: config-resolving client constructors, render.Options
// assembly, JSON decoding, hex normalization, the gRPC-gateway error
// envelope, and generic cobra command builders for the common
// "GET endpoint → render" and "CometBFT RPC → render" shapes.
package cmdutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// Pkg carries the per-command-package registration state. Each command
// package declares one and points Flags at the shared flag struct in
// its Register function.
type Pkg struct {
	// Name is the package label used in "not registered" errors.
	Name string
	// Flags is injected by the package's Register call.
	Flags *config.Flags
}

// resolve guards against use before Register and resolves the config.
func (p *Pkg) resolve() (*config.Config, error) {
	if p.Flags == nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("%s package not registered (flags unset)", p.Name)}
	}
	cfg, err := config.Resolve(p.Flags)
	if err != nil {
		return nil, &client.UsageError{Msg: err.Error()}
	}
	return cfg, nil
}

// RESTClient resolves the config and constructs a RESTClient. When
// --curl is set the HTTP call is replaced by a printed curl command.
func (p *Pkg) RESTClient(cmd *cobra.Command) (*client.RESTClient, *config.Config, error) {
	cfg, err := p.resolve()
	if err != nil {
		return nil, nil, err
	}
	c := client.NewRESTClient(cfg.RESTURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		c.Transport = &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
	}
	return c, cfg, nil
}

// RPCClient resolves the config and constructs an RPCClient against
// the CometBFT endpoint. When --curl is set the RPC call does not
// execute; it prints an equivalent curl command instead.
func (p *Pkg) RPCClient(cmd *cobra.Command) (*client.RPCClient, *config.Config, error) {
	cfg, err := p.resolve()
	if err != nil {
		return nil, nil, err
	}
	c := client.NewRPCClient(cfg.RPCURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		c.Transport = &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
	}
	return c, cfg, nil
}

// RenderOpts turns a resolved config into a render.Options instance,
// honouring --json, --field, --color, --raw, and TTY detection.
func RenderOpts(cmd *cobra.Command, cfg *config.Config, fields []string) render.Options {
	return render.Options{
		JSON:   cfg.JSON,
		Raw:    cfg.Raw,
		Fields: fields,
		Color:  cfg.Color,
		IsTTY:  IsTerminal(cmd.OutOrStdout()),
	}
}

// IsTerminal returns true if w is an *os.File attached to a terminal.
func IsTerminal(w io.Writer) bool {
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

// DecodeJSONMap decodes raw into a map[string]any. Used by REST
// responses whose top-level shape is always an object.
func DecodeJSONMap(raw []byte, label string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("decoding %s: %w (body=%q)", label, err, Truncate(raw, 256))
	}
	return m, nil
}

// Truncate clips b to at most n bytes for inclusion in error messages.
func Truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}

// CallEmpty issues an RPC call with an explicit empty params object.
// CometBFT's reflect-based RPC layer rejects nil params on some
// methods with "reflect: Call with too few input arguments"; passing
// map[string]any{} avoids that trap while still producing a valid
// JSON-RPC envelope.
func CallEmpty(ctx context.Context, rpc *client.RPCClient, method string) (json.RawMessage, error) {
	raw, err := rpc.Call(ctx, method, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("calling %s: %w", method, err)
	}
	return raw, nil
}

// DecodeGeneric unmarshals raw into any (map/slice/string). Used when
// we want to emit --json passthrough or pluck via --field.
func DecodeGeneric(raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return v, nil
}
