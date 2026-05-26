package tx

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// broadcastRequest is the shape accepted by /cosmos/tx/v1beta1/txs.
// BROADCAST_MODE_SYNC is the sane default: wait for CheckTx,
// don't wait for block inclusion.
type broadcastRequest struct {
	TxBytes string `json:"tx_bytes"`
	Mode    string `json:"mode"`
}

// newPublishCmd builds `publish <TX>`. Accepts a TxRaw as either
// base64 (the REST gateway's native format) or hex (as `cast publish`
// emits). Requires --yes because it is the only state-changing
// subcommand in this group.
func newPublishCmd() *cobra.Command {
	var yes bool
	var mode string
	var fields []string
	cmd := &cobra.Command{
		Use:   "publish <TX>",
		Short: "Broadcast a signed TxRaw (base64 or hex).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBytesB64, err := normalizeTxBytes(args[0])
			if err != nil {
				return err
			}
			if mode == "" {
				mode = "BROADCAST_MODE_SYNC"
			}
			payload := broadcastRequest{TxBytes: txBytesB64, Mode: mode}
			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshalling broadcast request: %w", err)
			}

			if !yes {
				// Not a usage error per se — the user supplied
				// enough information, they just haven't opted in.
				// Cast-style exit 3 communicates "aborted" via
				// UsageError so the caller gets a non-zero rc.
				if _, werr := fmt.Fprintf(cmd.OutOrStdout(),
					"would broadcast tx_bytes=%s mode=%s\nre-run with --yes to send\n",
					txBytesB64, mode); werr != nil {
					return werr
				}
				return &client.UsageError{Msg: "publish requires --yes"}
			}

			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			respBody, status, err := rest.Post(cmd.Context(), "/cosmos/tx/v1beta1/txs", "application/json", body)
			if err != nil {
				return err
			}
			if status == 0 && respBody == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			var generic any
			if err := json.Unmarshal(respBody, &generic); err != nil {
				return fmt.Errorf("decoding broadcast response: %w", err)
			}
			opts.JSON = true
			return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&yes, "yes", false, "confirm broadcast (required)")
	f.StringVar(&mode, "mode", "BROADCAST_MODE_SYNC", "broadcast mode")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// normalizeTxBytes accepts TxRaw as base64 or hex (with optional 0x
// prefix) and returns the base64 form the REST gateway expects.
func normalizeTxBytes(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", &client.UsageError{Msg: "empty tx"}
	}
	// Try hex first — unambiguous when prefixed.
	hs := s
	if strings.HasPrefix(hs, "0x") || strings.HasPrefix(hs, "0X") {
		hs = hs[2:]
		b, err := hex.DecodeString(hs)
		if err != nil {
			return "", &client.UsageError{Msg: fmt.Sprintf("invalid hex tx: %v", err)}
		}
		return base64.StdEncoding.EncodeToString(b), nil
	}
	// Plain hex (even length, hex chars only)?
	if len(hs)%2 == 0 && looksHex(hs) {
		if b, err := hex.DecodeString(hs); err == nil {
			return base64.StdEncoding.EncodeToString(b), nil
		}
	}
	// Fall back to base64 decode-then-reencode to normalise form.
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return base64.StdEncoding.EncodeToString(b), nil
	}
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return base64.StdEncoding.EncodeToString(b), nil
	}
	return "", &client.UsageError{Msg: "tx is neither hex nor base64"}
}

func looksHex(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}
