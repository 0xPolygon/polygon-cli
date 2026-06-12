package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// validatorJoinMsgShort triggers the L1-mirroring guard.
const validatorJoinMsgShort = "MsgValidatorJoin"

func init() {
	RegisterFactory("stake-join", newStakeJoinCmd)
}

// newStakeJoinCmd builds `stake-join` (MsgValidatorJoin). L1-mirroring.
func newStakeJoinCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag        string
		valID           uint64
		activationEpoch uint64
		amount          string
		signerPubKeyHex string
		l1Ref           stakeL1Ref
	)
	cmd := &cobra.Command{
		Use:   "stake-join",
		Short: "Register a validator (MsgValidatorJoin, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.stake.MsgValidatorJoin.

Produced by the bridge after a StakingInfo event; manual use requires
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
			if err = requireNonEmptyString("amount", amount); err != nil {
				return err
			}
			pubKey, err := parseHexBytes("signer-pub-key", signerPubKeyHex, 0)
			if err != nil {
				return err
			}
			if len(pubKey) == 0 {
				return &client.UsageError{Msg: "--signer-pub-key is required"}
			}
			txHash, err := l1Ref.txHash()
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, validatorJoinMsgShort, from, &htx.ValidatorJoinMsg{
				From: from, ValID: valID, ActivationEpoch: activationEpoch,
				Amount: strings.TrimSpace(amount), SignerPubKey: pubKey,
				TxHash: txHash, LogIndex: l1Ref.logIndex,
				BlockNumber: l1Ref.blockNumber, Nonce: l1Ref.nonce,
			})
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	registerFromMsgFlag(f, &fromFlag, validatorJoinMsgShort)
	f.Uint64Var(&valID, "val-id", 0, "validator id")
	f.Uint64Var(&activationEpoch, "activation-epoch", 0, "activation epoch")
	f.StringVar(&amount, "amount", "", "stake amount (decimal string)")
	f.StringVar(&signerPubKeyHex, "signer-pub-key", "", "validator signer pubkey (hex)")
	l1Ref.registerFlags(f)
	return cmd
}
