package msgs

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// parseHexBytes returns the decoded bytes of a 0x-prefixed or bare hex
// string. Empty input returns a nil slice and no error so msg fields
// that accept optional bytes can pass --flag="" through unchanged.
// expectedLen == 0 disables the length check.
func parseHexBytes(flagName, raw string, expectedLen int) ([]byte, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s)%2 != 0 {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s: odd hex length %d", flagName, len(s))}
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s: invalid hex: %v", flagName, err)}
	}
	if expectedLen > 0 && len(b) != expectedLen {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s must be %d bytes (got %d)", flagName, expectedLen, len(b))}
	}
	return b, nil
}

// requireNonEmptyString returns a UsageError when s is blank.
func requireNonEmptyString(flagName, s string) error {
	if strings.TrimSpace(s) == "" {
		return &client.UsageError{Msg: fmt.Sprintf("--%s is required", flagName)}
	}
	return nil
}

// requireNonZero returns a UsageError when v is zero.
func requireNonZero(flagName string, v uint64) error {
	if v == 0 {
		return &client.UsageError{Msg: fmt.Sprintf("--%s is required", flagName)}
	}
	return nil
}

// signerOrFlagAddress resolves the address used in a Msg's from /
// proposer / voter field: when flagValue is blank the signing key is
// resolved and its lowercase Eth address used; otherwise flagValue is
// validated and normalised via lowerEthAddress. flagName names the
// flag in validation errors.
func signerOrFlagAddress(cmd *cobra.Command, opts *TxOpts, flagName, flagValue string) (string, error) {
	addr := strings.TrimSpace(flagValue)
	if addr == "" {
		signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
		if err != nil {
			return "", err
		}
		return strings.ToLower(signer.Address.Hex()), nil
	}
	return lowerEthAddress(flagName, addr)
}

// executeSingleMsg wraps msg in a single-message Plan and hands it to
// the shared Execute pipeline. Every msg subcommand's RunE funnels
// through here.
func executeSingleMsg(cmd *cobra.Command, opts *TxOpts, mode Mode, msgShort, signerAddr string, msg htx.Msg) error {
	return Execute(cmd, opts, mode, &Plan{
		Msgs:          []htx.Msg{msg},
		MsgShortType:  msgShort,
		SignerAddress: signerAddr,
	})
}

// registerFromMsgFlag binds the shared --from-msg flag whose usage
// string names the proto msg type, e.g.
// "MsgSignerUpdate.from address (default: signer)".
func registerFromMsgFlag(f *pflag.FlagSet, target *string, msgShort string) {
	f.StringVar(target, "from-msg", "", msgShort+".from address (default: signer)")
}

// stakeL1Ref carries the L1 event reference flags shared by the four
// stake msgs (join / update / exit / signer-update).
type stakeL1Ref struct {
	txHashHex   string
	logIndex    uint64
	blockNumber uint64
	nonce       uint64
}

// registerFlags binds the shared L1 reference flags.
func (r *stakeL1Ref) registerFlags(f *pflag.FlagSet) {
	f.StringVar(&r.txHashHex, "tx-hash", "", "L1 tx hash (32 bytes hex)")
	f.Uint64Var(&r.logIndex, "log-index", 0, "L1 log index")
	f.Uint64Var(&r.blockNumber, "block-number", 0, "L1 block number")
	f.Uint64Var(&r.nonce, "nonce-l1", 0, "L1 stake nonce")
}

// txHash parses --tx-hash as a 32-byte hex value (empty allowed).
func (r *stakeL1Ref) txHash() ([]byte, error) {
	return parseHexBytes("tx-hash", r.txHashHex, 32)
}

// lowerEthAddress normalises s to lowercase 0x-prefixed hex. Also
// validates the 20-byte length. Returns a UsageError otherwise.
func lowerEthAddress(flagName, s string) (string, error) {
	s = strings.TrimSpace(s)
	trimmed := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(trimmed) != 40 {
		return "", &client.UsageError{Msg: fmt.Sprintf("--%s must be a 20-byte hex address", flagName)}
	}
	for _, c := range trimmed {
		ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !ok {
			return "", &client.UsageError{Msg: fmt.Sprintf("--%s must be hex", flagName)}
		}
	}
	return "0x" + strings.ToLower(trimmed), nil
}
