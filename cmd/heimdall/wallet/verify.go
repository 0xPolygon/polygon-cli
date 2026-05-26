package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newVerifyCmd builds `wallet verify <address> <message> <sig>`:
// recover the signer from a signature and compare with the supplied
// address. Prints `ok` / `mismatch` and returns a usage-grade error
// on a mismatch so the exit code is 3 (same convention as cast).
func newVerifyCmd() *cobra.Command {
	var raw bool
	cmd := &cobra.Command{
		Use:   "verify <address> <message-or-hash> <signature>",
		Short: "Verify a signature against an address.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !common.IsHexAddress(args[0]) {
				return &client.UsageError{Msg: fmt.Sprintf("invalid address %q", args[0])}
			}
			addr := common.HexToAddress(args[0])
			sig, err := parseSignatureHex(args[2])
			if err != nil {
				return err
			}
			var ok bool
			if raw {
				hash, err := parseHex(args[1], "hash")
				if err != nil {
					return err
				}
				if len(hash) != 32 {
					return &client.UsageError{Msg: fmt.Sprintf("--raw input must decode to 32 bytes, got %d", len(hash))}
				}
				ok, err = verifyRaw(addr, hash, sig)
				if err != nil {
					return err
				}
			} else {
				payload := []byte(args[1])
				if decoded, err := parseHex(args[1], "message"); err == nil {
					payload = decoded
				}
				ok, err = verifyPersonal(addr, payload, sig)
				if err != nil {
					return err
				}
			}
			w := cmd.OutOrStdout()
			if !ok {
				fmt.Fprintln(w, "mismatch")
				return &client.UsageError{Msg: "signature does not match address"}
			}
			fmt.Fprintln(w, "ok")
			return nil
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "verify against a 32-byte hash (no EIP-191 framing)")
	return cmd
}
