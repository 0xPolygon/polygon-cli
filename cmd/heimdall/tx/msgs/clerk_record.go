package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
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
		fromFlag        string
		txHash          string
		logIndex        uint64
		blockNumber     uint64
		contractAddr    string
		dataHex         string
		recordID        uint64
		chainID         string
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
			from := strings.TrimSpace(fromFlag)
			if from == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				from = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("from-msg", from)
				if err != nil {
					return err
				}
				from = p
			}
			if strings.TrimSpace(txHash) == "" {
				return &client.UsageError{Msg: "--tx-hash is required"}
			}
			if strings.TrimSpace(contractAddr) == "" {
				return &client.UsageError{Msg: "--contract-address is required"}
			}
			caddr, err := lowerEthAddress("contract-address", contractAddr)
			if err != nil {
				return err
			}
			if recordID == 0 {
				return &client.UsageError{Msg: "--id is required"}
			}
			data, err := parseHexBytes("data", dataHex, 0)
			if err != nil {
				return err
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.ClerkEventRecordMsg{
					From: from, TxHash: strings.TrimSpace(txHash),
					LogIndex: logIndex, BlockNumber: blockNumber,
					ContractAddress: caddr, Data: data,
					ID: recordID, ChainID: strings.TrimSpace(chainID),
				}},
				MsgShortType:  clerkRecordMsgShort,
				SignerAddress: from,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&fromFlag, "from-msg", "", "MsgEventRecord.from address (default: signer)")
	f.StringVar(&txHash, "tx-hash", "", "L1 tx hash (hex string; proto field is string)")
	f.Uint64Var(&logIndex, "log-index", 0, "L1 log index")
	f.Uint64Var(&blockNumber, "block-number", 0, "L1 block number")
	f.StringVar(&contractAddr, "contract-address", "", "L1 contract emitting the event")
	f.StringVar(&dataHex, "data", "", "event payload (hex-encoded bytes)")
	f.Uint64Var(&recordID, "id", 0, "record id")
	f.StringVar(&chainID, "source-chain-id", "", "source L1 chain id")
	return cmd
}
