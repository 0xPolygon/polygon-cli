package decode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// newHashTxCmd builds `decode hash-tx <B64_OR_HEX>`. Returns the
// upper-case SHA-256 hash CometBFT uses to key transactions. The hash
// is computed over the raw TxRaw bytes verbatim — no decode step is
// strictly required, but we still accept either encoding so the
// invocation mirrors `decode tx`.
func newHashTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-tx <tx-raw>",
		Short: "Compute the CometBFT SHA-256 hash of a TxRaw (hex or base64).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := decodeInput("tx", args[0])
			if err != nil {
				return err
			}
			sum := sha256.Sum256(raw)
			hash := strings.ToUpper(hex.EncodeToString(sum[:]))
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "0x"+hash)
			return err
		},
	}
	return cmd
}
