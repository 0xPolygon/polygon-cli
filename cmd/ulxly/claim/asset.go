package claim

import (
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed claimAssetUsage.md
var claimAssetUsage string

var AssetCmd = &cobra.Command{
	Use:     "asset",
	Short:   "Claim a deposit.",
	Long:    claimAssetUsage,
	PreRunE: ulxlycommon.PrepInputs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := claimAsset(cmd); err != nil {
			log.Fatal().Err(err).Msg("Received critical error")
		}
		return nil
	},
	SilenceUsage: true,
}

func claimAsset(cmd *cobra.Command) error {
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

	if deposit.LeafType != 0 {
		log.Warn().Msg("Deposit leafType is not asset")
	}
	if globalIndexOverride != "" {
		deposit.GlobalIndex.SetString(globalIndexOverride, 10)
	}

	proof, err := getMerkleProofsExitRoots(ulxlycommon.BridgeService, *deposit, proofGERHash, proofL1InfoTreeIndex)
	if err != nil {
		log.Error().Err(err).Msg("error getting merkle proofs and exit roots from bridge service")
		return err
	}

	claimTxn, err := bridgeV2.ClaimAsset(auth, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	if err = logAndReturnJSONError(cmd.Context(), client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return ulxlycommon.WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}

func getDepositWhenReadyForClaim(depositNetwork, depositCount uint32, wait time.Duration) (*bridge_service.Deposit, error) {
	var deposit *bridge_service.Deposit
	var err error

	waiter := time.After(wait)

out:
	for {
		deposit, err = getDeposit(depositNetwork, depositCount)
		if err == nil {
			log.Info().Msg("The deposit is ready to be claimed")
			break out
		}

		select {
		case <-waiter:
			if wait != 0 {
				err = fmt.Errorf("the deposit seems to be stuck after %s", wait.String())
			}
			break out
		default:
			if errors.Is(err, ErrNotReadyForClaim) || errors.Is(err, bridge_service.ErrNotFound) {
				log.Info().Msg("retrying...")
				time.Sleep(10 * time.Second)
				continue
			}
			break out
		}
	}
	return deposit, err
}

func getDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	deposit, err := ulxlycommon.BridgeService.GetDeposit(depositNetwork, depositCount)
	if err != nil {
		return nil, err
	}

	if ulxlycommon.InputArgs.Legacy {
		if !deposit.ReadyForClaim {
			log.Error().Msg("The claim transaction is not yet ready to be claimed. Try again in a few blocks.")
			return nil, ErrNotReadyForClaim
		} else if deposit.ClaimTxHash != nil {
			log.Info().Str("claimTxHash", deposit.ClaimTxHash.String()).Msg(ErrDepositAlreadyClaimed.Error())
			return nil, ErrDepositAlreadyClaimed
		}
	}

	return deposit, nil
}

func getMerkleProofsExitRoots(bridgeService bridge_service.BridgeService, deposit bridge_service.Deposit, proofGERHash string, l1InfoTreeIndex uint32) (*bridge_service.Proof, error) {
	var ger *common.Hash
	if len(proofGERHash) > 0 {
		hash := common.HexToHash(proofGERHash)
		ger = &hash
	}

	var proof *bridge_service.Proof
	var err error
	if ger != nil {
		proof, err = bridgeService.GetProofByGer(deposit.NetworkID, deposit.DepositCnt, *ger)
	} else if l1InfoTreeIndex > 0 {
		proof, err = bridgeService.GetProofByL1InfoTreeIndex(deposit.NetworkID, deposit.DepositCnt, l1InfoTreeIndex)
	} else {
		proof, err = bridgeService.GetProof(deposit.NetworkID, deposit.DepositCnt)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting proof for deposit %d on network %d: %w", deposit.DepositCnt, deposit.NetworkID, err)
	}

	if len(proof.MerkleProof) == 0 {
		errMsg := "the Merkle Proofs cannot be retrieved, double check the input arguments and try again"
		log.Error().
			Str("url", bridgeService.Url()).
			Uint32("NetworkID", deposit.NetworkID).
			Uint32("DepositCnt", deposit.DepositCnt).
			Msg(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(proof.RollupMerkleProof) == 0 {
		errMsg := "the Rollup Merkle Proofs cannot be retrieved, double check the input arguments and try again"
		log.Error().
			Str("url", bridgeService.Url()).
			Uint32("NetworkID", deposit.NetworkID).
			Uint32("DepositCnt", deposit.DepositCnt).
			Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	if proof.MainExitRoot == nil || proof.RollupExitRoot == nil {
		errMsg := "the exit roots from the bridge service were empty"
		log.Warn().
			Uint32("DepositCnt", deposit.DepositCnt).
			Uint32("OrigNet", deposit.OrigNet).
			Uint32("DestNet", deposit.DestNet).
			Uint32("NetworkID", deposit.NetworkID).
			Stringer("OrigAddr", deposit.OrigAddr).
			Stringer("DestAddr", deposit.DestAddr).
			Msg("deposit can't be claimed!")
		log.Error().
			Str("url", bridgeService.Url()).
			Uint32("NetworkID", deposit.NetworkID).
			Uint32("DepositCnt", deposit.DepositCnt).
			Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	return proof, nil
}
