package claim

import (
	_ "embed"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed claimMessageUsage.md
var claimMessageUsage string

var MessageCmd = &cobra.Command{
	Use:     "message",
	Short:   "Claim a message.",
	Long:    claimMessageUsage,
	PreRunE: ulxlycommon.PrepInputs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := claimMessage(cmd); err != nil {
			log.Fatal().Err(err).Msg("Received critical error")
		}
		return nil
	},
	SilenceUsage: true,
}

func claimMessage(cmd *cobra.Command) error {
	bridgeAddress := ulxlycommon.InputArgs.BridgeAddress
	privateKey := ulxlycommon.InputArgs.PrivateKey
	gasLimit := ulxlycommon.InputArgs.GasLimit
	destinationAddress := ulxlycommon.InputArgs.DestAddress
	chainID := ulxlycommon.InputArgs.ChainID
	timeoutTxnReceipt := ulxlycommon.InputArgs.Timeout
	RPCURL := ulxlycommon.InputArgs.RPCURL
	depositCount := ulxlycommon.InputArgs.DepositCount
	depositNetwork := ulxlycommon.InputArgs.DepositNetwork
	globalIndexOverride := ulxlycommon.InputArgs.GlobalIndex
	proofGERHash := ulxlycommon.InputArgs.ProofGER
	proofL1InfoTreeIndex := ulxlycommon.InputArgs.ProofL1InfoTreeIndex
	wait := ulxlycommon.InputArgs.Wait

	// Dial Ethereum client
	client, err := ulxlycommon.CreateEthClient(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, _, auth, err := ulxlycommon.GenerateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	deposit, err := getDepositWhenReadyForClaim(depositNetwork, depositCount, wait)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if deposit.LeafType != 1 {
		log.Warn().Msg("Deposit leafType is not message")
	}
	if globalIndexOverride != "" {
		deposit.GlobalIndex.SetString(globalIndexOverride, 10)
	}

	proof, err := getMerkleProofsExitRoots(ulxlycommon.BridgeService, *deposit, proofGERHash, proofL1InfoTreeIndex)
	if err != nil {
		log.Error().Err(err).Msg("error getting merkle proofs and exit roots from bridge service")
		return err
	}

	claimTxn, err := bridgeV2.ClaimMessage(auth, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	if err = logAndReturnJSONError(cmd.Context(), client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return ulxlycommon.WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}
