package decode

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// newTxCmd builds `decode tx <B64_OR_HEX>`. It unmarshals a TxRaw,
// decodes the body + auth info, resolves every Any.type_url against
// the internal registry, and prints either a human-readable summary or
// JSON when --json is set.
func newTxCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "tx <tx-raw>",
		Short: "Decode a TxRaw (base64 or 0x-hex) and pretty-print its contents.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := decodeInput("tx", args[0])
			if err != nil {
				return err
			}
			out, err := decodeTxRaw(raw)
			if err != nil {
				return err
			}
			if jsonOut {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
			}
			return writeTxSummary(cmd, out)
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit JSON instead of key/value output")
	return cmd
}

// decodedTx is the JSON-friendly shape produced by `decode tx`.
type decodedTx struct {
	TxHashSHA256 string         `json:"tx_hash_sha256"`
	Body         decodedTxBody  `json:"body"`
	AuthInfo     decodedAuth    `json:"auth_info"`
	Signatures   []string       `json:"signatures"`
}

type decodedTxBody struct {
	Memo          string                   `json:"memo,omitempty"`
	TimeoutHeight uint64                   `json:"timeout_height,omitempty"`
	Messages      []map[string]interface{} `json:"messages"`
}

type decodedAuth struct {
	Fee         *decodedFee     `json:"fee,omitempty"`
	SignerInfos []decodedSigner `json:"signer_infos,omitempty"`
}

type decodedFee struct {
	Amount   []hproto.Coin `json:"amount,omitempty"`
	GasLimit uint64        `json:"gas_limit"`
	Payer    string        `json:"payer,omitempty"`
	Granter  string        `json:"granter,omitempty"`
}

type decodedSigner struct {
	PublicKey *decodedAny `json:"public_key,omitempty"`
	ModeInfo  string      `json:"mode_info,omitempty"`
	Sequence  uint64      `json:"sequence"`
}

type decodedAny struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value_b64,omitempty"`
}

// decodeTxRaw is shared between `decode tx` (summary) and
// `decode hash-tx` (hash-only).
func decodeTxRaw(raw []byte) (*decodedTx, error) {
	txRaw, err := hproto.UnmarshalTxRaw(raw)
	if err != nil {
		return nil, err
	}
	// Decode body.
	body, err := unmarshalTxBody(txRaw.BodyBytes)
	if err != nil {
		return nil, fmt.Errorf("tx body: %w", err)
	}
	msgs := make([]map[string]interface{}, 0, len(body.Messages))
	for _, any := range body.Messages {
		m := map[string]interface{}{"type_url": any.TypeURL}
		if decoded, err := hproto.Decode(any.TypeURL, any.Value); err == nil {
			m["value"] = decoded
		} else {
			m["value_b64"] = base64.StdEncoding.EncodeToString(any.Value)
			m["decode_error"] = err.Error()
		}
		msgs = append(msgs, m)
	}

	// Decode auth info.
	auth, err := unmarshalAuthInfo(txRaw.AuthInfoBytes)
	if err != nil {
		return nil, fmt.Errorf("auth info: %w", err)
	}

	sigs := make([]string, 0, len(txRaw.Signatures))
	for _, s := range txRaw.Signatures {
		sigs = append(sigs, "0x"+hex.EncodeToString(s))
	}

	h := sha256.Sum256(raw)
	return &decodedTx{
		TxHashSHA256: strings.ToUpper(hex.EncodeToString(h[:])),
		Body: decodedTxBody{
			Memo:          body.Memo,
			TimeoutHeight: body.TimeoutHeight,
			Messages:      msgs,
		},
		AuthInfo:   auth,
		Signatures: sigs,
	}, nil
}

// writeTxSummary emits a human-readable summary of decodedTx.
func writeTxSummary(cmd *cobra.Command, d *decodedTx) error {
	w := cmd.OutOrStdout()
	if _, err := fmt.Fprintf(w, "tx_hash=%s\n", d.TxHashSHA256); err != nil {
		return err
	}
	if d.Body.Memo != "" {
		if _, err := fmt.Fprintf(w, "memo=%s\n", d.Body.Memo); err != nil {
			return err
		}
	}
	if d.Body.TimeoutHeight != 0 {
		if _, err := fmt.Fprintf(w, "timeout_height=%d\n", d.Body.TimeoutHeight); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "messages=%d\n", len(d.Body.Messages)); err != nil {
		return err
	}
	for i, m := range d.Body.Messages {
		if _, err := fmt.Fprintf(w, "  [%d] %v\n", i, m["type_url"]); err != nil {
			return err
		}
		if m["decode_error"] != nil {
			if _, err := fmt.Fprintf(w, "       error: %v\n", m["decode_error"]); err != nil {
				return err
			}
			continue
		}
		// For registered messages we render as JSON at two-space indent so
		// nested proto values remain readable.
		buf, err := json.MarshalIndent(m["value"], "       ", "  ")
		if err != nil {
			continue
		}
		if _, err := fmt.Fprintf(w, "       %s\n", string(buf)); err != nil {
			return err
		}
	}
	if d.AuthInfo.Fee != nil {
		fee := d.AuthInfo.Fee
		if _, err := fmt.Fprintf(w, "fee.gas_limit=%d\n", fee.GasLimit); err != nil {
			return err
		}
		for _, c := range fee.Amount {
			if _, err := fmt.Fprintf(w, "fee.amount=%s%s\n", c.Amount, c.Denom); err != nil {
				return err
			}
		}
		if fee.Payer != "" {
			if _, err := fmt.Fprintf(w, "fee.payer=%s\n", fee.Payer); err != nil {
				return err
			}
		}
	}
	for i, si := range d.AuthInfo.SignerInfos {
		if _, err := fmt.Fprintf(w, "signer[%d].sequence=%d\n", i, si.Sequence); err != nil {
			return err
		}
		if si.PublicKey != nil {
			if _, err := fmt.Fprintf(w, "signer[%d].pubkey.type_url=%s\n", i, si.PublicKey.TypeURL); err != nil {
				return err
			}
		}
		if si.ModeInfo != "" {
			if _, err := fmt.Fprintf(w, "signer[%d].mode=%s\n", i, si.ModeInfo); err != nil {
				return err
			}
		}
	}
	for i, s := range d.Signatures {
		if _, err := fmt.Fprintf(w, "signature[%d]=%s\n", i, s); err != nil {
			return err
		}
	}
	return nil
}

// unmarshalTxBody parses a TxBody. Only the fields we render (messages,
// memo, timeout_height) are kept.
type txBodyParsed struct {
	Messages      []*hproto.Any
	Memo          string
	TimeoutHeight uint64
}

func unmarshalTxBody(b []byte) (*txBodyParsed, error) {
	out := &txBodyParsed{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, err
		}
		switch num {
		case 1:
			any, err := hproto.UnmarshalAny(val)
			if err != nil {
				return nil, err
			}
			out.Messages = append(out.Messages, any)
		case 2:
			out.Memo = string(val)
		case 3:
			v, err := rawVarint(val)
			if err != nil {
				return nil, err
			}
			out.TimeoutHeight = v
		}
		b = b[n:]
	}
	return out, nil
}

// unmarshalAuthInfo parses a cosmos AuthInfo.
func unmarshalAuthInfo(b []byte) (decodedAuth, error) {
	out := decodedAuth{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, err
		}
		switch num {
		case 1:
			si, err := unmarshalSignerInfo(val)
			if err != nil {
				return out, err
			}
			out.SignerInfos = append(out.SignerInfos, si)
		case 2:
			fee, err := unmarshalFee(val)
			if err != nil {
				return out, err
			}
			out.Fee = fee
		}
		b = b[n:]
	}
	return out, nil
}

func unmarshalSignerInfo(b []byte) (decodedSigner, error) {
	out := decodedSigner{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, err
		}
		switch num {
		case 1:
			any, err := hproto.UnmarshalAny(val)
			if err != nil {
				return out, err
			}
			out.PublicKey = &decodedAny{
				TypeURL: any.TypeURL,
				Value:   base64.StdEncoding.EncodeToString(any.Value),
			}
		case 2:
			out.ModeInfo = parseModeInfo(val)
		case 3:
			v, err := rawVarint(val)
			if err != nil {
				return out, err
			}
			out.Sequence = v
		}
		b = b[n:]
	}
	return out, nil
}

// parseModeInfo extracts ModeInfo.Single.mode as a SignMode string.
func parseModeInfo(b []byte) string {
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return ""
		}
		if num == 1 {
			for len(val) > 0 {
				inNum, _, inVal, inN, err := consumeField(val)
				if err != nil {
					return ""
				}
				if inNum == 1 {
					mode, err := rawVarint(inVal)
					if err != nil {
						return ""
					}
					return signModeString(int32(mode))
				}
				val = val[inN:]
			}
		}
		b = b[n:]
	}
	return ""
}

func signModeString(m int32) string {
	switch m {
	case hproto.SignModeUnspecif:
		return "UNSPECIFIED"
	case hproto.SignModeDirect:
		return "DIRECT"
	case hproto.SignModeAminoJSON:
		return "LEGACY_AMINO_JSON"
	default:
		return fmt.Sprintf("MODE(%d)", m)
	}
}

func unmarshalFee(b []byte) (*decodedFee, error) {
	out := &decodedFee{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, err
		}
		switch num {
		case 1:
			c, err := unmarshalCoin(val)
			if err != nil {
				return nil, err
			}
			out.Amount = append(out.Amount, c)
		case 2:
			v, err := rawVarint(val)
			if err != nil {
				return nil, err
			}
			out.GasLimit = v
		case 3:
			out.Payer = string(val)
		case 4:
			out.Granter = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

func unmarshalCoin(b []byte) (hproto.Coin, error) {
	out := hproto.Coin{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, err
		}
		switch num {
		case 1:
			out.Denom = string(val)
		case 2:
			out.Amount = string(val)
		}
		b = b[n:]
	}
	return out, nil
}
