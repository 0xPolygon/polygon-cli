package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// clerkRecordMsgShort triggers the L1-mirroring guard.
const clerkRecordMsgShort = "MsgEventRecord"

func init() {
	RegisterFactory("clerk-record", newClerkRecordCmd)
}

// newClerkRecordCmd builds `clerk-record` (MsgEventRecord). This is the
// clerk state-sync event the bridge submits after observing an L1
// StateSync event; manual use requires --force.
func newClerkRecordCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag     string
		txHash       string
		logIndex     uint64
		blockNumber  uint64
		contractAddr string
		dataHex      string
		recordID     uint64
		chainID      string
	)
	cmd := &cobra.Command{
		Use:   "clerk-record",
		Short: "Submit an L1 state-sync record (MsgEventRecord, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.clerk.MsgEventRecord.

Produced by the bridge after an L1 StateSync event; manual use requires
--force.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := signerOrFlagAddress(cmd, opts, "from-msg", fromFlag)
			if err != nil {
				return err
			}
			if err = requireNonEmptyString("tx-hash", txHash); err != nil {
				return err
			}
			if err = requireNonEmptyString("contract-address", contractAddr); err != nil {
				return err
			}
			caddr, err := lowerEthAddress("contract-address", contractAddr)
			if err != nil {
				return err
			}
			if err = requireNonZero("id", recordID); err != nil {
				return err
			}
			data, err := parseHexBytes("data", dataHex, 0)
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, clerkRecordMsgShort, from, &htx.ClerkEventRecordMsg{
				From: from, TxHash: strings.TrimSpace(txHash),
				LogIndex: logIndex, BlockNumber: blockNumber,
				ContractAddress: caddr, Data: data,
				ID: recordID, ChainID: strings.TrimSpace(chainID),
			})
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	registerFromMsgFlag(f, &fromFlag, clerkRecordMsgShort)
	f.StringVar(&txHash, "tx-hash", "", "L1 tx hash (hex string; proto field is string)")
	f.Uint64Var(&logIndex, "log-index", 0, "L1 log index")
	f.Uint64Var(&blockNumber, "block-number", 0, "L1 block number")
	f.StringVar(&contractAddr, "contract-address", "", "L1 contract emitting the event")
	f.StringVar(&dataHex, "data", "", "event payload (hex-encoded bytes)")
	f.Uint64Var(&recordID, "id", 0, "record id")
	f.StringVar(&chainID, "source-chain-id", "", "source L1 chain id")
	return cmd
}
