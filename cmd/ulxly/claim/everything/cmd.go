// Package everything provides the claim-everything command.
package everything

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type DepositID struct {
	DepositCnt uint32 `json:"deposit_cnt"`
	NetworkID  uint32 `json:"network_id"`
}

var Cmd = &cobra.Command{
	Use:     "claim-everything",
	Short:   "Attempt to claim as many deposits and messages as possible.",
	PreRunE: ulxlycommon.PrepInputs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return claimEverything(cmd)
	},
	SilenceUsage: true,
}

func init() {
	// Add shared transaction flags (rpc-url, bridge-address, private-key, etc.)
	ulxlycommon.AddTransactionFlags(Cmd)

	// Claim-everything specific flags
	f := Cmd.Flags()
	f.StringSliceVar(&ulxlycommon.InputArgs.BridgeServiceURLs, ulxlycommon.ArgBridgeMappings, nil, "network ID to bridge service URL mappings (e.g. '1=http://network-1-bridgeurl,7=http://network-2-bridgeurl')")
	f.IntVar(&ulxlycommon.InputArgs.BridgeLimit, ulxlycommon.ArgBridgeLimit, 25, "limit the number or responses returned by the bridge service when claiming")
	f.IntVar(&ulxlycommon.InputArgs.BridgeOffset, ulxlycommon.ArgBridgeOffset, 0, "offset to specify for pagination of underlying bridge service deposits")
	f.UintVar(&ulxlycommon.InputArgs.Concurrency, ulxlycommon.ArgConcurrency, 1, "worker pool size for claims")
	flag.MarkFlagsRequired(Cmd, ulxlycommon.ArgBridgeMappings)
}

func claimEverything(cmd *cobra.Command) error {
	privateKey := ulxlycommon.InputArgs.PrivateKey
	claimerAddress := ulxlycommon.InputArgs.AddressOfPrivateKey

	gasLimit := ulxlycommon.InputArgs.GasLimit
	chainID := ulxlycommon.InputArgs.ChainID
	timeoutTxnReceipt := ulxlycommon.InputArgs.Timeout
	bridgeAddress := ulxlycommon.InputArgs.BridgeAddress
	destinationAddress := ulxlycommon.InputArgs.DestAddress
	RPCURL := ulxlycommon.InputArgs.RPCURL
	limit := ulxlycommon.InputArgs.BridgeLimit
	offset := ulxlycommon.InputArgs.BridgeOffset
	concurrency := ulxlycommon.InputArgs.Concurrency

	depositMap := make(map[DepositID]*bridge_service.Deposit)

	for networkID, bs := range ulxlycommon.BridgeServices {
		deposits, _, bErr := getDepositsForAddress(bs, destinationAddress, offset, limit)
		if bErr != nil {
			log.Err(bErr).Uint32("id", networkID).Str("url", bs.Url()).Msgf("Error getting deposits for bridge: %s", bErr.Error())
			return bErr
		}
		for idx, deposit := range deposits {
			depID := DepositID{
				DepositCnt: deposit.DepositCnt,
				NetworkID:  deposit.NetworkID,
			}
			_, hasKey := depositMap[depID]
			// if we haven't seen this deposit at all, we'll store it
			if !hasKey {
				depositMap[depID] = &deposits[idx]
				continue
			}

			// if this new deposit is ready for claim OR it has already been claimed we should override the existing value
			if ulxlycommon.InputArgs.Legacy {
				if deposit.ReadyForClaim || deposit.ClaimTxHash != nil {
					depositMap[depID] = &deposits[idx]
				}
			}
		}
	}

	client, err := ulxlycommon.CreateEthClient(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()

	bridgeContract, _, opts, err := ulxlycommon.GenerateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		return err
	}
	currentNetworkID, err := bridgeContract.NetworkID(nil)
	if err != nil {
		return err
	}
	log.Info().Uint32("networkID", currentNetworkID).Msg("current network")

	workPool := make(chan *bridge_service.Deposit, concurrency) // bounded chan for controlled concurrency

	nonceCounter, err := currentNonce(cmd.Context(), client, claimerAddress)
	if err != nil {
		return err
	}
	log.Info().Int64("nonce", nonceCounter.Int64()).Msg("starting nonce")
	nonceMutex := sync.Mutex{}
	nonceIncrement := big.NewInt(1)
	retryNonces := make(chan *big.Int, concurrency) // bounded same as workPool

	wg := sync.WaitGroup{} // wg so the last ones can get processed

	for _, d := range depositMap {
		wg.Add(1)
		workPool <- d // block until a slot is available
		go func(deposit *bridge_service.Deposit) {
			defer func() {
				<-workPool // release work slot
			}()
			defer wg.Done()

			if deposit.DestNet != currentNetworkID {
				log.Debug().Uint32("destination_network", deposit.DestNet).Msg("discarding deposit for different network")
				return
			}
			if deposit.ClaimTxHash != nil {
				log.Info().Str("txhash", deposit.ClaimTxHash.String()).Msg("It looks like this tx was already claimed")
				return
			}
			// Either use the next retry nonce, or set and increment the next one
			var nextNonce *big.Int
			select {
			case n := <-retryNonces:
				nextNonce = n
			default:
				nonceMutex.Lock()
				nextNonce = big.NewInt(nonceCounter.Int64())
				nonceCounter = nonceCounter.Add(nonceCounter, nonceIncrement)
				nonceMutex.Unlock()
			}
			log.Info().Int64("nonce", nextNonce.Int64()).Msg("Next nonce")

			claimTx, dErr := claimSingleDeposit(cmd, client, bridgeContract, withNonce(opts, nextNonce), *deposit, ulxlycommon.BridgeServices, currentNetworkID)
			if dErr != nil {
				log.Warn().Err(dErr).Uint32("DepositCnt", deposit.DepositCnt).
					Uint32("OrigNet", deposit.OrigNet).
					Uint32("DestNet", deposit.DestNet).
					Uint32("NetworkID", deposit.NetworkID).
					Stringer("OrigAddr", deposit.OrigAddr).
					Stringer("DestAddr", deposit.DestAddr).
					Int64("nonce", nextNonce.Int64()).
					Msg("There was an error claiming")

				// Some nonces should not be reused
				if strings.Contains(dErr.Error(), "could not replace existing") {
					return
				}
				if strings.Contains(dErr.Error(), "already known") {
					return
				}
				if strings.Contains(dErr.Error(), "nonce is too low") {
					return
				}
				// are there other cases?
				retryNonces <- nextNonce
				return
			}
			dErr = ulxlycommon.WaitMineTransaction(cmd.Context(), client, claimTx, timeoutTxnReceipt)
			if dErr != nil {
				log.Error().Err(dErr).Msg("error while waiting for tx to mine")
			}
		}(d)
	}

	wg.Wait()
	return nil
}

func currentNonce(ctx context.Context, client *ethclient.Client, address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	nonce, err := client.NonceAt(ctx, addr, nil)
	if err != nil {
		log.Error().Err(err).Str("address", addr.Hex()).Msg("Failed to get nonce")
		return nil, err
	}
	n := int64(nonce)
	return big.NewInt(n), nil
}

// withNonce creates a copy of the TransactOpts with a new nonce
func withNonce(opts *bind.TransactOpts, newNonce *big.Int) *bind.TransactOpts {
	if opts == nil {
		return nil
	}
	clone := &bind.TransactOpts{
		From:     opts.From,
		Signer:   opts.Signer,
		GasLimit: opts.GasLimit,
		Context:  opts.Context, // Usually OK to share, unless you need a separate context
		NoSend:   opts.NoSend,
	}
	// Deep-copy big.Int fields
	if opts.Value != nil {
		clone.Value = new(big.Int).Set(opts.Value)
	}
	if opts.GasFeeCap != nil {
		clone.GasFeeCap = new(big.Int).Set(opts.GasFeeCap)
	}
	if opts.GasTipCap != nil {
		clone.GasTipCap = new(big.Int).Set(opts.GasTipCap)
	}
	// Set the new nonce
	if newNonce != nil {
		clone.Nonce = new(big.Int).Set(newNonce)
	}

	return clone
}

func claimSingleDeposit(cmd *cobra.Command, client *ethclient.Client, bridgeContract *ulxly.Ulxly, opts *bind.TransactOpts, deposit bridge_service.Deposit, bridgeServices map[uint32]bridge_service.BridgeService, currentNetworkID uint32) (*types.Transaction, error) {
	networkIDForBridgeService := deposit.NetworkID
	if deposit.NetworkID == 0 {
		networkIDForBridgeService = currentNetworkID
	}

	bridgeServiceFromMap, hasKey := bridgeServices[networkIDForBridgeService]
	if !hasKey {
		return nil, fmt.Errorf("we don't have a bridge service url for network: %d", deposit.DestNet)
	}

	proof, err := getMerkleProofsExitRoots(bridgeServiceFromMap, deposit, "", 0)
	if err != nil {
		log.Error().Err(err).Msg("error getting merkle proofs and exit roots from bridge service")
		return nil, err
	}

	var claimTx *types.Transaction
	if deposit.LeafType == 0 {
		claimTx, err = bridgeContract.ClaimAsset(opts, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	} else {
		claimTx, err = bridgeContract.ClaimMessage(opts, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	}

	if err = ulxlycommon.LogAndReturnJSONError(cmd.Context(), client, claimTx, opts, err); err != nil {
		log.Warn().
			Uint32("DepositCnt", deposit.DepositCnt).
			Uint32("OrigNet", deposit.OrigNet).
			Uint32("DestNet", deposit.DestNet).
			Uint32("NetworkID", deposit.NetworkID).
			Stringer("OrigAddr", deposit.OrigAddr).
			Stringer("DestAddr", deposit.DestAddr).
			Msg("attempt to claim deposit failed")
		return nil, err
	}
	log.Info().Stringer("txhash", claimTx.Hash()).Msg("sent claim")

	return claimTx, nil
}

func getDepositsForAddress(bridgeService bridge_service.BridgeService, destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	deposits, total, err := bridgeService.GetDeposits(destinationAddress, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(deposits) != total {
		log.Warn().Int("total_deposits", total).Int("retrieved_deposits", len(deposits)).Msg("not all deposits were retrieved")
	}

	return deposits, total, nil
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
		return nil, fmt.Errorf("%s", errMsg)
	}
	if len(proof.RollupMerkleProof) == 0 {
		errMsg := "the Rollup Merkle Proofs cannot be retrieved, double check the input arguments and try again"
		log.Error().
			Str("url", bridgeService.Url()).
			Uint32("NetworkID", deposit.NetworkID).
			Uint32("DepositCnt", deposit.DepositCnt).
			Msg(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
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
		return nil, fmt.Errorf("%s", errMsg)
	}

	return proof, nil
}
