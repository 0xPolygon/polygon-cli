package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// signerUpdateMsgShort triggers the L1-mirroring guard.
const signerUpdateMsgShort = "MsgSignerUpdate"

func init() {
	RegisterFactory("signer-update", newSignerUpdateCmd)
}

// newSignerUpdateCmd builds `signer-update` (MsgSignerUpdate).
func newSignerUpdateCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag           string
		valID              uint64
		newSignerPubKeyHex string
		l1Ref              stakeL1Ref
	)
	cmd := &cobra.Command{
		Use:   "signer-update",
		Short: "Rotate validator signer pubkey (MsgSignerUpdate, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.stake.MsgSignerUpdate.

Produced by the bridge after a SignerChange event; manual use requires
--force.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := signerOrFlagAddress(cmd, opts, "from-msg", fromFlag)
			if err != nil {
				return err
			}
			if err = requireNonZero("val-id", valID); err != nil {
				return err
			}
			pubKey, err := parseHexBytes("new-signer-pub-key", newSignerPubKeyHex, 0)
			if err != nil {
				return err
			}
			if len(pubKey) == 0 {
				return &client.UsageError{Msg: "--new-signer-pub-key is required"}
			}
			txHash, err := l1Ref.txHash()
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, signerUpdateMsgShort, from, &htx.SignerUpdateMsg{
				From: from, ValID: valID, NewSignerPubKey: pubKey,
				TxHash: txHash, LogIndex: l1Ref.logIndex,
				BlockNumber: l1Ref.blockNumber, Nonce: l1Ref.nonce,
			})
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	registerFromMsgFlag(f, &fromFlag, signerUpdateMsgShort)
	f.Uint64Var(&valID, "val-id", 0, "validator id")
	f.StringVar(&newSignerPubKeyHex, "new-signer-pub-key", "", "new signer pubkey (hex)")
	l1Ref.registerFlags(f)
	return cmd
}
