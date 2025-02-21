package ulxly

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	// TreeDepth of 32 is pulled directly from the
	// _DEPOSIT_CONTRACT_TREE_DEPTH from the smart contract. We
	// could make this a variable as well
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/54f58c8b64806429bc4d5c52248f29cf80ba401c/contracts/v2/lib/DepositContractBase.sol#L15
	TreeDepth = 32
)

var (
	ErrNotReadyForClaim        = errors.New("the claim transaction is not yet ready to be claimed, try again in a few blocks")
	ErrDepositAlreadyClaimed   = errors.New("the claim transaction has already been claimed")
	ErrUnableToRetrieveDeposit = errors.New("the bridge deposit was not found")
)

type IMT struct {
	Branches   map[uint32][]common.Hash
	Leaves     map[uint32]common.Hash
	Roots      []common.Hash
	ZeroHashes []common.Hash
	Proofs     map[uint32]Proof
}

type Proof struct {
	Siblings     [TreeDepth]common.Hash
	Root         common.Hash
	DepositCount uint32
	LeafHash     common.Hash
}

type BridgeProof struct {
	Proof struct {
		MerkleProof       []string `json:"merkle_proof"`
		RollupMerkleProof []string `json:"rollup_merkle_proof"`
		MainExitRoot      string   `json:"main_exit_root"`
		RollupExitRoot    string   `json:"rollup_exit_root"`
	} `json:"proof"`
}
type BridgeDeposit struct {
	LeafType      uint8  `json:"leaf_type"`
	OrigNet       uint32 `json:"orig_net"`
	OrigAddr      string `json:"orig_addr"`
	Amount        string `json:"amount"`
	DestNet       uint32 `json:"dest_net"`
	DestAddr      string `json:"dest_addr"`
	BlockNum      string `json:"block_num"`
	DepositCnt    uint32 `json:"deposit_cnt"`
	NetworkID     uint32 `json:"network_id"`
	TxHash        string `json:"tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash"`
	Metadata      string `json:"metadata"`
	ReadyForClaim bool   `json:"ready_for_claim"`
	GlobalIndex   string `json:"global_index"`
}

type DepositID struct {
	DepositCnt uint32 `json:"deposit_cnt"`
	NetworkID  uint32 `json:"network_id"`
}

type BridgeDepositResponse struct {
	Deposit BridgeDeposit `json:"deposit"`
	Code    *int          `json:"code"`
	Message *string       `json:"message"`
}

func readDeposit(cmd *cobra.Command) error {
	rpcUrl, err := cmd.Flags().GetString(ArgRPCURL)
	if err != nil {
		return err
	}
	bridgeAddress, err := cmd.Flags().GetString(ArgBridgeAddress)
	if err != nil {
		return err
	}

	toBlock := *inputUlxlyArgs.toBlock
	fromBlock := *inputUlxlyArgs.fromBlock
	filter := *inputUlxlyArgs.filterSize

	// Dial the Ethereum RPC server.
	rpc, err := ethrpc.DialContext(cmd.Context(), rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer rpc.Close()
	ec := ethclient.NewClient(rpc)

	bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(bridgeAddress), ec)
	if err != nil {
		return err
	}
	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := currentBlock + filter
		if endBlock > toBlock {
			endBlock = toBlock
		}

		opts := bind.FilterOpts{
			Start:   currentBlock,
			End:     &endBlock,
			Context: cmd.Context(),
		}
		evtV2Iterator, err := bridgeV2.FilterBridgeEvent(&opts)
		if err != nil {
			return err
		}

		for evtV2Iterator.Next() {
			evt := evtV2Iterator.Event
			log.Info().Uint32("deposit", evt.DepositCount).Uint64("block-number", evt.Raw.BlockNumber).Msg("Found ulxly Deposit")
			var jBytes []byte
			jBytes, err = json.Marshal(evt)
			if err != nil {
				return err
			}
			fmt.Println(string(jBytes))
		}
		err = evtV2Iterator.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing event iterator")
		}
		currentBlock = endBlock
	}

	return nil
}

func proof(args []string) error {
	depositNumber := *inputUlxlyArgs.depositNumber
	rawDepositData, err := getInputData(args)
	if err != nil {
		return err
	}
	return readDeposits(rawDepositData, uint32(depositNumber))
}

func emptyProof() error {
	p := new(Proof)

	e := generateEmptyHashes(TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(p.String())
	return nil
}

func zeroProof() error {
	p := new(Proof)

	e := generateZeroHashes(TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(p.String())
	return nil
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func logAndReturnJsonError(cmd *cobra.Command, client *ethclient.Client, tx *types.Transaction, opts *bind.TransactOpts, err error) error {

	var callErr error
	if tx != nil {
		// in case the error came down to gas estimation, we can sometimes get more information by doing a call
		_, callErr = client.CallContract(cmd.Context(), ethereum.CallMsg{
			From:          opts.From,
			To:            tx.To(),
			Gas:           tx.Gas(),
			GasPrice:      tx.GasPrice(),
			GasFeeCap:     tx.GasFeeCap(),
			GasTipCap:     tx.GasTipCap(),
			Value:         tx.Value(),
			Data:          tx.Data(),
			AccessList:    tx.AccessList(),
			BlobGasFeeCap: tx.BlobGasFeeCap(),
			BlobHashes:    tx.BlobHashes(),
		}, nil)

		if *inputUlxlyArgs.dryRun {
			castCmd := "cast call"
			castCmd += fmt.Sprintf(" --rpc-url %s", *inputUlxlyArgs.rpcURL)
			castCmd += fmt.Sprintf(" --from %s", opts.From.String())
			castCmd += fmt.Sprintf(" --gas-limit %d", tx.Gas())
			if tx.Type() == types.LegacyTxType {
				castCmd += fmt.Sprintf(" --gas-price %s", tx.GasPrice().String())
			} else {
				castCmd += fmt.Sprintf(" --gas-price %s", tx.GasFeeCap().String())
				castCmd += fmt.Sprintf(" --priority-gas-price %s", tx.GasTipCap().String())
			}
			castCmd += fmt.Sprintf(" --value %s", tx.Value().String())
			castCmd += fmt.Sprintf(" %s", tx.To().String())
			castCmd += fmt.Sprintf(" %s", common.Bytes2Hex(tx.Data()))
			log.Info().Str("cmd", castCmd).Msg("use this command to replicate the call")
		}
	}

	if err == nil {
		return nil
	}

	var jsonError JsonError
	jsonErrorBytes, jsErr := json.Marshal(err)
	if jsErr != nil {
		log.Error().Err(err).Msg("Unable to interact with the bridge contract")
		return err
	}

	jsErr = json.Unmarshal(jsonErrorBytes, &jsonError)
	if jsErr != nil {
		log.Error().Err(err).Msg("Unable to interact with the bridge contract")
		return err
	}

	errLog := log.Error().
		Err(err).
		Str("message", jsonError.Message).
		Int("code", jsonError.Code).
		Interface("data", jsonError.Data)

	if callErr != nil {
		errLog = errLog.Err(callErr)
	}

	if errCode, isValid := jsonError.Data.(string); isValid && errCode == "0x646cf558" {
		// I don't want to bother with the additional error logging for previously claimed deposits
		return err
	}

	errLog.Msg("Unable to interact with bridge contract")

	return err
}

func bridgeAsset(cmd *cobra.Command) error {
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	destinationAddress := *inputUlxlyArgs.destAddress
	chainID := *inputUlxlyArgs.chainID
	amount := *inputUlxlyArgs.value
	tokenAddr := *inputUlxlyArgs.tokenAddress
	callDataString := *inputUlxlyArgs.callData
	destinationNetwork := *inputUlxlyArgs.destNetwork
	isForced := *inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	RPCURL := *inputUlxlyArgs.rpcURL

	// Dial the Ethereum RPC server.
	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	value, _ := big.NewInt(0).SetString(amount, 0)
	tokenAddress := common.HexToAddress(tokenAddr)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	}

	bridgeTxn, err := bridgeV2.BridgeAsset(auth, destinationNetwork, toAddress, value, tokenAddress, isForced, callData)
	if err = logAndReturnJsonError(cmd, client, bridgeTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt)
}

func bridgeMessage(cmd *cobra.Command) error {
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	destinationAddress := *inputUlxlyArgs.destAddress
	chainID := *inputUlxlyArgs.chainID
	amount := *inputUlxlyArgs.value
	tokenAddr := *inputUlxlyArgs.tokenAddress
	callDataString := *inputUlxlyArgs.callData
	destinationNetwork := *inputUlxlyArgs.destNetwork
	isForced := *inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	RPCURL := *inputUlxlyArgs.rpcURL

	// Dial the Ethereum RPC server.
	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	value, _ := big.NewInt(0).SetString(amount, 0)
	tokenAddress := common.HexToAddress(tokenAddr)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	}

	bridgeTxn, err := bridgeV2.BridgeMessage(auth, destinationNetwork, toAddress, isForced, callData)
	if err = logAndReturnJsonError(cmd, client, bridgeTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt)
}

func bridgeWETHMessage(cmd *cobra.Command) error {
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	destinationAddress := *inputUlxlyArgs.destAddress
	chainID := *inputUlxlyArgs.chainID
	amount := *inputUlxlyArgs.value
	callDataString := *inputUlxlyArgs.callData
	destinationNetwork := *inputUlxlyArgs.destNetwork
	isForced := *inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	RPCURL := *inputUlxlyArgs.rpcURL

	// Dial the Ethereum RPC server.
	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}
	// Check if WETH is allowed
	wethAddress, err := bridgeV2.WETHToken(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error().Err(err).Msg("error getting WETH address from the bridge smc")
		return err
	}
	if wethAddress == (common.Address{}) {
		return fmt.Errorf("bridge WETH not allowed. Native ETH token configured in this network. This tx will fail")
	}

	value, _ := big.NewInt(0).SetString(amount, 0)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	bridgeTxn, err := bridgeV2.BridgeMessageWETH(auth, destinationNetwork, toAddress, value, isForced, callData)
	if err = logAndReturnJsonError(cmd, client, bridgeTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt)
}

func claimAsset(cmd *cobra.Command) error {
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	destinationAddress := *inputUlxlyArgs.destAddress
	chainID := *inputUlxlyArgs.chainID
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	RPCURL := *inputUlxlyArgs.rpcURL
	depositCount := *inputUlxlyArgs.depositCount
	depositNetwork := *inputUlxlyArgs.depositNetwork
	bridgeServiceUrl := *inputUlxlyArgs.bridgeServiceURL
	globalIndexOverride := *inputUlxlyArgs.globalIndex
	wait := *inputUlxlyArgs.wait

	// Dial Ethereum client
	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	globalIndex, amount, originAddress, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err := getDepositWhenReadyForClaim(bridgeServiceUrl, depositNetwork, depositCount, wait)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if leafType != 0 {
		log.Warn().Msg("Deposit leafType is not asset")
	}

	if globalIndexOverride != "" {
		globalIndex.SetString(globalIndexOverride, 10)
	}

	// Call the bridge service RPC URL to get the merkle proofs and exit roots and parses them to the correct formats.
	bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%d&net_id=%d", bridgeServiceUrl, depositCount, depositNetwork)
	merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)

	claimTxn, err := bridgeV2.ClaimAsset(auth, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), claimOriginalNetwork, originAddress, claimDestNetwork, toAddress, amount, metadata)
	if err = logAndReturnJsonError(cmd, client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}

func claimMessage(cmd *cobra.Command) error {
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	destinationAddress := *inputUlxlyArgs.destAddress
	chainID := *inputUlxlyArgs.chainID
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	RPCURL := *inputUlxlyArgs.rpcURL
	depositCount := *inputUlxlyArgs.depositCount
	depositNetwork := *inputUlxlyArgs.depositNetwork
	bridgeServiceUrl := *inputUlxlyArgs.bridgeServiceURL
	globalIndexOverride := *inputUlxlyArgs.globalIndex
	wait := *inputUlxlyArgs.wait

	// Dial Ethereum client
	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	globalIndex, amount, originAddress, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err := getDepositWhenReadyForClaim(bridgeServiceUrl, depositNetwork, depositCount, wait)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if leafType != 1 {
		log.Warn().Msg("Deposit leafType is not message")
	}
	if globalIndexOverride != "" {
		globalIndex.SetString(globalIndexOverride, 10)
	}

	// Call the bridge service RPC URL to get the merkle proofs and exit roots and parses them to the correct formats.
	bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%d&net_id=%d", bridgeServiceUrl, depositCount, depositNetwork)
	merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)

	claimTxn, err := bridgeV2.ClaimMessage(auth, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), claimOriginalNetwork, originAddress, claimDestNetwork, toAddress, amount, metadata)
	if err = logAndReturnJsonError(cmd, client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}

func getDepositWhenReadyForClaim(bridgeServiceUrl string, depositNetwork uint64, depositCount uint64, wait time.Duration) (*big.Int, *big.Int, common.Address, []byte, uint8, uint32, uint32, error) {
	var globalIndex, amount *big.Int
	var originAddress common.Address
	var metadata []byte
	var leafType uint8
	var claimDestNetwork, claimOriginalNetwork uint32
	var err error

	waiter := time.After(wait)

out:
	for {
		// Call the bridge service RPC URL to get the deposits data and parses them to the correct formats.
		bridgeServiceDepositsEndpoint := fmt.Sprintf("%s/bridge?net_id=%d&deposit_cnt=%d", bridgeServiceUrl, depositNetwork, depositCount)
		globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err = getDeposit(bridgeServiceDepositsEndpoint)
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
			if errors.Is(err, ErrNotReadyForClaim) || errors.Is(err, ErrUnableToRetrieveDeposit) {
				log.Info().Msg("retrying...")
				time.Sleep(10 * time.Second)
				continue
			}
			break out
		}
	}
	return globalIndex, amount, originAddress, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err
}

func getBridgeServiceURLs() (map[uint32]string, error) {
	bridgeServiceUrls := *inputUlxlyArgs.bridgeServiceURLs
	urlMap := make(map[uint32]string)
	for _, mapping := range bridgeServiceUrls {
		pieces := strings.Split(mapping, "=")
		if len(pieces) != 2 {
			return nil, fmt.Errorf("bridge service url mapping should contain a networkid and url separated by an equal sign. Got: %s", mapping)
		}
		networkId, err := strconv.ParseInt(pieces[0], 10, 32)
		if err != nil {
			return nil, err
		}
		urlMap[uint32(networkId)] = pieces[1]
	}
	return urlMap, nil
}

func claimEverything(cmd *cobra.Command) error {
	privateKey := *inputUlxlyArgs.privateKey
	gasLimit := *inputUlxlyArgs.gasLimit
	chainID := *inputUlxlyArgs.chainID
	timeoutTxnReceipt := *inputUlxlyArgs.timeout
	bridgeAddress := *inputUlxlyArgs.bridgeAddress
	destinationAddress := *inputUlxlyArgs.destAddress
	RPCURL := *inputUlxlyArgs.rpcURL
	limit := *inputUlxlyArgs.bridgeLimit
	offset := *inputUlxlyArgs.bridgeOffset
	urls, err := getBridgeServiceURLs()
	if err != nil {
		return err
	}

	depositMap := make(map[DepositID]*BridgeDeposit)

	for _, bridgeServiceUrl := range urls {
		deposits, bErr := getDepositsForAddress(fmt.Sprintf("%s/bridges/%s?offset=%d&limit=%d", bridgeServiceUrl, destinationAddress, offset, limit))
		if bErr != nil {
			return bErr
		}
		for idx, deposit := range deposits {
			depId := DepositID{
				DepositCnt: deposit.DepositCnt,
				NetworkID:  deposit.NetworkID,
			}
			_, hasKey := depositMap[depId]
			// if we haven't seen this deposit at all, we'll store it
			if !hasKey {
				depositMap[depId] = &deposits[idx]
				continue
			}

			// if this new deposit is ready for claim OR it has already been claimed we should override the existing value
			if deposit.ReadyForClaim || deposit.ClaimTxHash != "" {
				depositMap[depId] = &deposits[idx]
			}
		}
	}

	client, err := ethclient.DialContext(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()

	bridgeContract, _, opts, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		return err
	}
	currentNetworkID, err := bridgeContract.NetworkID(nil)
	if err != nil {
		return err
	}
	log.Info().Uint32("networkID", currentNetworkID).Msg("current network")

	concurrency := 1
	workPool := make(chan struct{}, concurrency) // Bounded chan for controlled concurrency

	nonceCounter, err := currentNonce(cmd.Context(), client, privateKey)
	if err != nil {
		return err
	}
	log.Info().Int64("nonce", nonceCounter.Int64()).Msg("starting nonce")
	nonceMutex := sync.Mutex{}
	nonceIncrement := big.NewInt(1)
	retryNonces := make(chan *big.Int, concurrency) // bounded so can async hand off

	for _, d := range depositMap {

		workPool <- struct{}{} // block until a slot is available

		go func(deposit *BridgeDeposit) {
			defer func() {
				<-workPool // release work slot
			}()

			if deposit.DestNet != currentNetworkID {
				log.Debug().Uint32("destination_network", deposit.DestNet).Msg("discarding deposit for different network")
				return
			}
			if deposit.ClaimTxHash != "" {
				log.Info().Str("txhash", deposit.ClaimTxHash).Msg("It looks like this tx was already claimed")
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

			claimTx, dErr := claimSingleDeposit(cmd, client, bridgeContract, withNonce(opts, nextNonce), *deposit, urls, currentNetworkID)
			if dErr != nil {
				log.Warn().Err(dErr).Uint32("DepositCnt", deposit.DepositCnt).
					Uint32("OrigNet", deposit.OrigNet).
					Uint32("DestNet", deposit.DestNet).
					Uint32("NetworkID", deposit.NetworkID).
					Str("OrigAddr", deposit.OrigAddr).
					Str("DestAddr", deposit.DestAddr).
					Msg("There was an error claiming")

				// Some nonces should not be reused
				if strings.Contains(dErr.Error(), "could not replace existing") {
					return
				}
				if strings.Contains(dErr.Error(), "already known") {
					return
				}

				retryNonces <- nextNonce

				return
			}
			dErr = WaitMineTransaction(cmd.Context(), client, claimTx, timeoutTxnReceipt)
			if dErr != nil {
				log.Error().Err(dErr).Msg("error while waiting for tx to mine")
				// if skip here, nonces will get screwed up? maybe bail everything here?
			}
		}(d)
	}

	return nil
}

func currentNonce(ctx context.Context, client *ethclient.Client, privateKey string) (*big.Int, error) {
	ecdsa, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to read private key")
		return nil, err
	}
	address := crypto.PubkeyToAddress(ecdsa.PublicKey)

	nonce, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		log.Error().Err(err).Str("address", address.Hex()).Msg("Failed to get nonce")
		return nil, err
	}
	n := int64(nonce)
	return big.NewInt(n), nil
}

// todo: implement for other fields in library, or find a library that does this
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

func claimSingleDeposit(cmd *cobra.Command, client *ethclient.Client, bridgeContract *ulxly.Ulxly, opts *bind.TransactOpts, deposit BridgeDeposit, bridgeURLs map[uint32]string, currentNetworkID uint32) (*types.Transaction, error) {
	networkIDForBridgeService := deposit.NetworkID
	if deposit.NetworkID == 0 {
		networkIDForBridgeService = currentNetworkID
	}
	bridgeUrl, hasKey := bridgeURLs[networkIDForBridgeService]
	if !hasKey {
		return nil, fmt.Errorf("we don't have a bridge service url for network: %d", deposit.DestNet)
	}
	bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%d&net_id=%d", bridgeUrl, deposit.DepositCnt, deposit.NetworkID)
	merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)
	if len(mainExitRoot) != 32 || len(rollupExitRoot) != 32 {
		log.Warn().
			Uint32("DepositCnt", deposit.DepositCnt).
			Uint32("OrigNet", deposit.OrigNet).
			Uint32("DestNet", deposit.DestNet).
			Uint32("NetworkID", deposit.NetworkID).
			Str("OrigAddr", deposit.OrigAddr).
			Str("DestAddr", deposit.DestAddr).
			Msg("deposit can't be claimed!")
		return nil, fmt.Errorf("the exit roots from the bridge service were empty: %s", bridgeServiceProofEndpoint)
	}

	globalIndex, isValid := new(big.Int).SetString(deposit.GlobalIndex, 10)
	if !isValid {
		return nil, fmt.Errorf("global index %s is not a valid integer", deposit.GlobalIndex)
	}
	amount, isValid := new(big.Int).SetString(deposit.Amount, 10)
	if !isValid {
		return nil, fmt.Errorf("amount %s is not a valid integer", deposit.Amount)
	}

	originAddress := common.HexToAddress(deposit.OrigAddr)
	toAddress := common.HexToAddress(deposit.DestAddr)
	metadata := common.Hex2Bytes(strings.TrimPrefix(deposit.Metadata, "0x"))

	var claimTx *types.Transaction
	var err error
	if deposit.LeafType == 0 {
		claimTx, err = bridgeContract.ClaimAsset(opts, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), deposit.OrigNet, originAddress, deposit.DestNet, toAddress, amount, metadata)
	} else {
		claimTx, err = bridgeContract.ClaimMessage(opts, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), deposit.OrigNet, originAddress, deposit.DestNet, toAddress, amount, metadata)
	}

	if err = logAndReturnJsonError(cmd, client, claimTx, opts, err); err != nil {
		log.Warn().
			Uint32("DepositCnt", deposit.DepositCnt).
			Uint32("OrigNet", deposit.OrigNet).
			Uint32("DestNet", deposit.DestNet).
			Uint32("NetworkID", deposit.NetworkID).
			Str("OrigAddr", deposit.OrigAddr).
			Str("DestAddr", deposit.DestAddr).
			Msg("attempt to claim deposit failed")
		return nil, err
	}
	log.Info().Stringer("txhash", claimTx.Hash()).Msg("sent claim")

	return claimTx, nil
}

// Wait for the transaction to be mined
func WaitMineTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, txTimeout uint64) error {
	if inputUlxlyArgs.dryRun != nil && *inputUlxlyArgs.dryRun {
		txJson, _ := tx.MarshalJSON()
		log.Info().RawJSON("tx", txJson).Msg("Skipping receipt check. Dry run is enabled")
		return nil
	}
	txnMinedTimer := time.NewTimer(time.Duration(txTimeout) * time.Second)
	defer txnMinedTimer.Stop()
	for {
		select {
		case <-txnMinedTimer.C:
			log.Info().Msg("Wait timer for transaction receipt exceeded!")
			return nil
		default:
			r, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				if err.Error() != "not found" {
					log.Error().Err(err)
					return err
				}
				time.Sleep(1 * time.Second)
				continue
			}
			if r.Status != 0 {
				log.Info().Interface("txHash", r.TxHash).Msg("Deposit transaction successful")
				return nil
			} else if r.Status == 0 {
				log.Error().Interface("txHash", r.TxHash).Msg("Deposit transaction failed")
				log.Info().Uint64("GasUsed", tx.Gas()).Uint64("cumulativeGasUsedForTx", r.CumulativeGasUsed).Msg("Perhaps try increasing the gas limit")
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func getInputData(args []string) ([]byte, error) {
	fileName := *inputUlxlyArgs.inputFileName
	if fileName != "" {
		return os.ReadFile(fileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
func readDeposits(rawDeposits []byte, depositNumber uint32) error {
	buf := bytes.NewBuffer(rawDeposits)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	imt := new(IMT)
	imt.Init()
	seenDeposit := make(map[uint32]common.Hash, 0)
	lastDeposit := uint32(0)
	for scanner.Scan() {
		evt := new(ulxly.UlxlyBridgeEvent)
		err := json.Unmarshal(scanner.Bytes(), evt)
		if err != nil {
			return err
		}
		if _, hasBeenSeen := seenDeposit[evt.DepositCount]; hasBeenSeen {
			log.Warn().Uint32("deposit", evt.DepositCount).Str("tx-hash", evt.Raw.TxHash.String()).Msg("Skipping duplicate deposit")
			continue
		}
		seenDeposit[evt.DepositCount] = evt.Raw.TxHash
		if lastDeposit+1 != evt.DepositCount && lastDeposit != 0 {
			log.Error().Uint32("missing-deposit", lastDeposit+1).Uint32("current-deposit", evt.DepositCount).Msg("Missing deposit")
			return fmt.Errorf("missing deposit: %d", lastDeposit+1)
		}
		lastDeposit = evt.DepositCount
		imt.AddLeaf(evt)
		log.Info().
			Uint64("block-number", evt.Raw.BlockNumber).
			Uint32("deposit-count", evt.DepositCount).
			Str("tx-hash", evt.Raw.TxHash.String()).
			Str("root", common.Hash(imt.Roots[len(imt.Roots)-1]).String()).
			Msg("adding event to tree")
		// There's no point adding more leaves if we can prove the deposit already?
		if evt.DepositCount >= depositNumber {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("there was an error reading the deposit file")
		return err
	}

	log.Info().Msg("finished")
	p := imt.GetProof(depositNumber)
	fmt.Println(p.String())
	return nil
}

// String will create the json representation of the proof
func (p *Proof) String() string {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling proof to json")
		return ""
	}
	return string(jsonBytes)

}

// hashDeposit create the leaf hash value for a particular deposit
func hashDeposit(deposit *ulxly.UlxlyBridgeEvent) common.Hash {
	var res common.Hash
	origNet := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(origNet, deposit.OriginNetwork)
	destNet := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(destNet, deposit.DestinationNetwork)
	var buf common.Hash
	metaHash := crypto.Keccak256Hash(deposit.Metadata)
	copy(res[:], crypto.Keccak256Hash([]byte{deposit.LeafType}, origNet, deposit.OriginAddress.Bytes(), destNet, deposit.DestinationAddress[:], deposit.Amount.FillBytes(buf[:]), metaHash.Bytes()).Bytes())
	return res
}

// Init will allocate the objects in the IMT
func (s *IMT) Init() {
	s.Branches = make(map[uint32][]common.Hash)
	s.Leaves = make(map[uint32]common.Hash)
	s.ZeroHashes = generateZeroHashes(TreeDepth)
	s.Proofs = make(map[uint32]Proof)
}

// AddLeaf will take a given deposit and add it to the collection of leaves. It will also update the
func (s *IMT) AddLeaf(deposit *ulxly.UlxlyBridgeEvent) {
	leaf := hashDeposit(deposit)
	log.Debug().Str("leaf-hash", common.Bytes2Hex(leaf[:])).Msg("Leaf hash calculated")
	// just keep a copy of the leaf indexed by deposit count for now
	s.Leaves[deposit.DepositCount] = leaf

	node := leaf
	size := uint64(deposit.DepositCount) + 1

	// copy the previous set of branches as a starting point. We're going to make copies of the branches at each deposit
	branches := make([]common.Hash, TreeDepth)
	if deposit.DepositCount == 0 {
		branches = generateEmptyHashes(TreeDepth)
	} else {
		copy(branches, s.Branches[deposit.DepositCount-1])
	}

	for height := uint64(0); height < TreeDepth; height += 1 {
		if ((size >> height) & 1) == 1 {
			copy(branches[height][:], node[:])
			break
		}
		node = crypto.Keccak256Hash(branches[height][:], node[:])
	}
	s.Branches[deposit.DepositCount] = branches
	s.Roots = append(s.Roots, s.GetRoot(deposit.DepositCount))
}

// GetRoot will return the root for a particular deposit
func (s *IMT) GetRoot(depositNum uint32) common.Hash {
	node := common.Hash{}
	size := depositNum + 1
	currentZeroHashHeight := common.Hash{}

	for height := 0; height < TreeDepth; height++ {
		if ((size >> height) & 1) == 1 {
			node = crypto.Keccak256Hash(s.Branches[depositNum][height][:], node.Bytes())

		} else {
			node = crypto.Keccak256Hash(node.Bytes(), currentZeroHashHeight.Bytes())
		}
		currentZeroHashHeight = crypto.Keccak256Hash(currentZeroHashHeight.Bytes(), currentZeroHashHeight.Bytes())
	}
	return node
}

// GetProof will return an object containing the proof data necessary for verification
func (s *IMT) GetProof(depositNum uint32) Proof {
	node := common.Hash{}
	size := depositNum + 1
	currentZeroHashHeight := common.Hash{}

	siblings := [TreeDepth]common.Hash{}
	for height := 0; height < TreeDepth; height++ {
		siblingDepositNum := getSiblingDepositNumber(depositNum, uint32(height))
		sibling := currentZeroHashHeight
		if _, hasKey := s.Branches[siblingDepositNum]; hasKey {
			sibling = s.Branches[siblingDepositNum][height]
		} else {
			sibling = currentZeroHashHeight
		}

		log.Info().Str("sibling", sibling.String()).Msg("Proof Inputs")
		siblings[height] = sibling
		if ((size >> height) & 1) == 1 {
			// node = keccak256(abi.encodePacked(_branch[height], node));
			node = crypto.Keccak256Hash(sibling.Bytes(), node.Bytes())
		} else {
			// node = keccak256(abi.encodePacked(node, currentZeroHashHeight));
			node = crypto.Keccak256Hash(node.Bytes(), sibling.Bytes())
		}
		currentZeroHashHeight = crypto.Keccak256Hash(currentZeroHashHeight.Bytes(), currentZeroHashHeight.Bytes())
	}
	p := &Proof{
		Siblings:     siblings,
		DepositCount: depositNum,
		LeafHash:     s.Leaves[depositNum],
	}

	r, err := p.Check(s.Roots)
	if err != nil {
		log.Error().Err(err).Msg("failed to validate proof")
	}
	p.Root = r
	s.Proofs[depositNum] = *p
	return *p
}

// getSiblingDepositNumber returns the sibling number of a given number at a specified level in an incremental Merkle tree.
//
// In an incremental Merkle tree, each node has a sibling node at each level of the tree.
// The sibling node can be determined by flipping the bit at the current level and setting all bits to the right of the current level to 1.
// This function calculates the sibling number based on the deposit number and the specified level.
//
// Parameters:
// - depositNumber: the original number for which the sibling is to be found.
// - level: the level in the Merkle tree at which to find the sibling.
//
// The logic works as follows:
// 1. `1 << level` creates a binary number with a single 1 bit at the position corresponding to the level.
// 2. `depositNumber ^ (1 << level)` flips the bit at the position corresponding to the level in the depositNumber.
// 3. `(1 << level) - 1` creates a binary number with all bits to the right of the current level set to 1.
// 4. `| ((1 << level) - 1)` ensures that all bits to the right of the current level are set to 1 in the result.
//
// The function effectively finds the sibling deposit number at each level of the Merkle tree by manipulating the bits accordingly.
func getSiblingDepositNumber(depositNumber, level uint32) uint32 {
	return depositNumber ^ (1 << level) | ((1 << level) - 1)
}

// Check is a sanity check of a proof in order to make sure that the
// proof that was generated creates a root that we recognize. This was
// useful while testing in order to avoid verifying that the proof
// works or doesn't work onchain
func (p *Proof) Check(roots []common.Hash) (common.Hash, error) {
	node := p.LeafHash
	index := p.DepositCount
	for height := 0; height < TreeDepth; height++ {
		if ((index >> height) & 1) == 1 {
			node = crypto.Keccak256Hash(p.Siblings[height][:], node[:])
		} else {
			node = crypto.Keccak256Hash(node[:], p.Siblings[height][:])
		}
	}

	isProofValid := false
	for i := len(roots) - 1; i >= 0; i-- {
		if roots[i].Cmp(node) == 0 {
			isProofValid = true
			break
		}
	}

	log.Info().
		Bool("is-proof-valid", isProofValid).
		Uint32("deposit-count", p.DepositCount).
		Str("leaf-hash", p.LeafHash.String()).
		Str("checked-root", node.String()).Msg("checking proof")
	if !isProofValid {
		return common.Hash{}, fmt.Errorf("invalid proof")
	}

	return node, nil
}

// https://eth2book.info/capella/part2/deposits-withdrawals/contract/
func generateZeroHashes(height uint8) []common.Hash {
	zeroHashes := make([]common.Hash, height)
	zeroHashes[0] = common.Hash{}
	for i := 1; i < int(height); i++ {
		zeroHashes[i] = crypto.Keccak256Hash(zeroHashes[i-1][:], zeroHashes[i-1][:])
	}
	return zeroHashes
}

func generateEmptyHashes(height uint8) []common.Hash {
	zeroHashes := make([]common.Hash, height)
	zeroHashes[0] = common.Hash{}
	for i := 1; i < int(height); i++ {
		zeroHashes[i] = common.Hash{}
	}
	return zeroHashes
}

func generateTransactionPayload(ctx context.Context, client *ethclient.Client, ulxlyInputArgBridge string, ulxlyInputArgPvtKey string, ulxlyInputArgGasLimit uint64, ulxlyInputArgDestAddr string, ulxlyInputArgChainID string) (bridgeV2 *ulxly.Ulxly, toAddress common.Address, opts *bind.TransactOpts, err error) {
	ulxlyInputArgPvtKey = strings.TrimPrefix(ulxlyInputArgPvtKey, "0x")
	bridgeV2, err = ulxly.NewUlxly(common.HexToAddress(ulxlyInputArgBridge), client)
	if err != nil {
		return
	}

	privateKey, err := crypto.HexToECDSA(ulxlyInputArgPvtKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve private key")
		return
	}

	// value := big.NewInt(*ulxlyInputArgs.Amount)
	gasLimit := ulxlyInputArgGasLimit

	chainID := new(big.Int)
	// For manual input of chainID, use the user's input
	if ulxlyInputArgChainID != "" {
		chainID.SetString(ulxlyInputArgChainID, 10)
	} else { // If there is no user input for chainID, infer it from context
		chainID, err = client.ChainID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Cannot get chain ID")
			return
		}
	}

	opts, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot generate transactionOpts")
		return
	}
	if inputUlxlyArgs.gasPrice != nil && *inputUlxlyArgs.gasPrice != "" {
		gasPrice := new(big.Int)
		gasPrice.SetString(*inputUlxlyArgs.gasPrice, 10)
		opts.GasPrice = gasPrice
	}
	if inputUlxlyArgs.dryRun != nil && *inputUlxlyArgs.dryRun {
		opts.NoSend = true
	}
	opts.Context = ctx
	opts.GasLimit = gasLimit
	toAddress = common.HexToAddress(ulxlyInputArgDestAddr)
	if toAddress == (common.Address{}) {
		toAddress = opts.From
	}
	return bridgeV2, toAddress, opts, err
}

func getMerkleProofsExitRoots(bridgeServiceProofEndpoint string) (merkleProofArray [32][32]byte, rollupMerkleProofArray [32][32]byte, mainExitRoot []byte, rollupExitRoot []byte) {
	reqBridgeProof, err := http.Get(bridgeServiceProofEndpoint)
	if err != nil {
		log.Error().Err(err)
		return
	}
	bodyBridgeProof, err := io.ReadAll(reqBridgeProof.Body) // Response body is []byte
	if err != nil {
		log.Error().Err(err)
		return
	}
	var bridgeProof BridgeProof
	err = json.Unmarshal(bodyBridgeProof, &bridgeProof) // Parse []byte to go struct pointer, and shadow err variable
	if err != nil {
		log.Error().Err(err).Msg("Can not unmarshal JSON")
		return
	}

	merkleProof := [][32]byte{}       // HACK: usage of common.Hash may be more consistent and considered best practice
	rollupMerkleProof := [][32]byte{} // HACK: usage of common.Hash may be more consistent and considered best practice

	for _, mp := range bridgeProof.Proof.MerkleProof {
		byteMP, _ := hexutil.Decode(mp)
		merkleProof = append(merkleProof, [32]byte(byteMP))
	}
	if len(merkleProof) == 0 {
		log.Error().Str("url", bridgeServiceProofEndpoint).Msg("The Merkle Proofs cannot be retrieved, double check the input arguments and try again.")
		return
	}
	merkleProofArray = [32][32]byte(merkleProof)
	for _, rmp := range bridgeProof.Proof.RollupMerkleProof {
		byteRMP, _ := hexutil.Decode(rmp)
		rollupMerkleProof = append(rollupMerkleProof, [32]byte(byteRMP))
	}
	if len(rollupMerkleProof) == 0 {
		log.Error().Msg("The Rollup Merkle Proofs cannot be retrieved, double check the input arguments and try again.")
		return
	}
	rollupMerkleProofArray = [32][32]byte(rollupMerkleProof)

	mainExitRoot, _ = hexutil.Decode(bridgeProof.Proof.MainExitRoot)
	rollupExitRoot, _ = hexutil.Decode(bridgeProof.Proof.RollupExitRoot)

	defer reqBridgeProof.Body.Close()

	return merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot
}

func getDeposit(bridgeServiceDepositsEndpoint string) (globalIndex *big.Int, originAddress common.Address, amount *big.Int, metadata []byte, leafType uint8, claimDestNetwork, claimOriginalNetwork uint32, err error) {
	reqBridgeDeposit, err := http.Get(bridgeServiceDepositsEndpoint)
	if err != nil {
		log.Error().Err(err)
		return
	}
	bodyBridgeDeposit, err := io.ReadAll(reqBridgeDeposit.Body) // Response body is []byte
	if err != nil {
		log.Error().Err(err)
		return
	}
	var bridgeDeposit BridgeDepositResponse
	err = json.Unmarshal(bodyBridgeDeposit, &bridgeDeposit) // Parse []byte to go struct pointer, and shadow err variable
	if err != nil {
		log.Error().Err(err).Msg("Can not unmarshal JSON")
		return
	}

	globalIndex = new(big.Int)
	amount = new(big.Int)

	defer reqBridgeDeposit.Body.Close()
	if bridgeDeposit.Code != nil {
		log.Warn().Int("code", *bridgeDeposit.Code).Str("message", *bridgeDeposit.Message).Msg("unable to retrieve bridge deposit")
		return globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, ErrUnableToRetrieveDeposit
	}

	if !bridgeDeposit.Deposit.ReadyForClaim {
		log.Error().Msg("The claim transaction is not yet ready to be claimed. Try again in a few blocks.")
		return nil, common.HexToAddress("0x0"), nil, nil, 0, 0, 0, ErrNotReadyForClaim
	} else if bridgeDeposit.Deposit.ClaimTxHash != "" {
		log.Info().Str("claimTxHash", bridgeDeposit.Deposit.ClaimTxHash).Msg("The claim transaction has already been claimed")
		return nil, common.HexToAddress("0x0"), nil, nil, 0, 0, 0, ErrDepositAlreadyClaimed
	}
	originAddress = common.HexToAddress(bridgeDeposit.Deposit.OrigAddr)
	globalIndex.SetString(bridgeDeposit.Deposit.GlobalIndex, 10)
	amount.SetString(bridgeDeposit.Deposit.Amount, 10)

	metadata = common.Hex2Bytes(strings.TrimPrefix(bridgeDeposit.Deposit.Metadata, "0x"))
	leafType = bridgeDeposit.Deposit.LeafType
	claimDestNetwork = bridgeDeposit.Deposit.DestNet
	claimOriginalNetwork = bridgeDeposit.Deposit.OrigNet
	log.Info().
		Stringer("globalIndex", globalIndex).
		Stringer("originAddress", originAddress).
		Stringer("amount", amount).
		Str("metadata", bridgeDeposit.Deposit.Metadata).
		Uint8("leafType", leafType).
		Uint32("claimDestNetwork", claimDestNetwork).
		Uint32("claimOriginalNetwork", claimOriginalNetwork).
		Msg("Got Deposit")
	return globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, nil
}

func getDepositsForAddress(bridgeRequestUrl string) ([]BridgeDeposit, error) {
	var resp struct {
		Deposits []BridgeDeposit `json:"deposits"`
		Total    int             `json:"total_cnt,string"`
	}
	httpResp, err := http.Get(bridgeRequestUrl)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Deposits) != resp.Total {
		log.Warn().Int("total_deposits", resp.Total).Int("retrieved_deposits", len(resp.Deposits)).Msg("not all deposits were retrieved")
	}

	return resp.Deposits, nil
}

//go:embed BridgeAssetUsage.md
var bridgeAssetUsage string

//go:embed BridgeMessageUsage.md
var bridgeMessageUsage string

//go:embed BridgeWETHMessageUsage.md
var bridgeWETHMessageUsage string

//go:embed ClaimAssetUsage.md
var claimAssetUsage string

//go:embed ClaimMessageUsage.md
var claimMessageUsage string

//go:embed proofUsage.md
var proofUsage string

//go:embed depositGetUsage.md
var depositGetUsage string

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the uLxLy bridge",
	Long:  "Basic utility commands for interacting with the bridge contracts, bridge services, and generating proofs",
	Args:  cobra.NoArgs,
}
var ulxlyBridgeAndClaimCmd = &cobra.Command{
	Args:   cobra.NoArgs,
	Hidden: true,
}

var ulxlyBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Commands for moving funds and sending messages from one chain to another",
	Args:  cobra.NoArgs,
}

var ulxlyClaimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Commands for claiming deposits on a particular chain",
	Args:  cobra.NoArgs,
}

type ulxlyArgs struct {
	gasLimit            *uint64
	chainID             *string
	privateKey          *string
	addressOfPrivateKey string
	value               *string
	rpcURL              *string
	bridgeAddress       *string
	destNetwork         *uint32
	destAddress         *string
	tokenAddress        *string
	forceUpdate         *bool
	callData            *string
	callDataFile        *string
	timeout             *uint64
	depositCount        *uint64
	depositNetwork      *uint64
	bridgeServiceURL    *string
	inputFileName       *string
	fromBlock           *uint64
	toBlock             *uint64
	filterSize          *uint64
	depositNumber       *uint64
	globalIndex         *string
	gasPrice            *string
	dryRun              *bool
	bridgeServiceURLs   *[]string
	bridgeLimit         *int
	bridgeOffset        *int
	wait                *time.Duration
}

var inputUlxlyArgs = ulxlyArgs{}

var (
	bridgeAssetCommand       *cobra.Command
	bridgeMessageCommand     *cobra.Command
	bridgeMessageWETHCommand *cobra.Command
	claimAssetCommand        *cobra.Command
	claimMessageCommand      *cobra.Command
	claimEverythingCommand   *cobra.Command
	emptyProofCommand        *cobra.Command
	zeroProofCommand         *cobra.Command
	proofCommand             *cobra.Command
	getDepositCommand        *cobra.Command
)

const (
	ArgGasLimit         = "gas-limit"
	ArgChainID          = "chain-id"
	ArgPrivateKey       = "private-key"
	ArgValue            = "value"
	ArgRPCURL           = "rpc-url"
	ArgBridgeAddress    = "bridge-address"
	ArgDestNetwork      = "destination-network"
	ArgDestAddress      = "destination-address"
	ArgForceUpdate      = "force-update-root"
	ArgCallData         = "call-data"
	ArgCallDataFile     = "call-data-file"
	ArgTimeout          = "transaction-receipt-timeout"
	ArgDepositCount     = "deposit-count"
	ArgDepositNetwork   = "deposit-network"
	ArgBridgeServiceURL = "bridge-service-url"
	ArgFileName         = "file-name"
	ArgFromBlock        = "from-block"
	ArgToBlock          = "to-block"
	ArgFilterSize       = "filter-size"
	ArgTokenAddress     = "token-address"
	ArgGlobalIndex      = "global-index"
	ArgDryRun           = "dry-run"
	ArgGasPrice         = "gas-price"
	ArgBridgeMappings   = "bridge-service-map"
	ArgBridgeLimit      = "bridge-limit"
	ArgBridgeOffset     = "bridge-offset"
	ArgWait             = "wait"
)

func prepInputs(cmd *cobra.Command, args []string) error {
	if *inputUlxlyArgs.dryRun && *inputUlxlyArgs.gasLimit == 0 {
		dryRunGasLimit := uint64(10_000_000)
		inputUlxlyArgs.gasLimit = &dryRunGasLimit
	}
	pvtKey := strings.TrimPrefix(*inputUlxlyArgs.privateKey, "0x")
	privateKey, err := crypto.HexToECDSA(pvtKey)
	if err != nil {
		return fmt.Errorf("invalid --%s: %w", ArgPrivateKey, err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	inputUlxlyArgs.addressOfPrivateKey = fromAddress.String()
	if *inputUlxlyArgs.destAddress == "" {
		*inputUlxlyArgs.destAddress = fromAddress.String()
		log.Info().Stringer("destAddress", fromAddress).Msg("No destination address specified. Using private key's address")
	}

	if *inputUlxlyArgs.callDataFile != "" {
		rawCallData, err := os.ReadFile(*inputUlxlyArgs.callDataFile)
		if err != nil {
			return err
		}
		if *inputUlxlyArgs.callData != "0x" {
			return fmt.Errorf("both %s and %s flags were provided", ArgCallData, ArgCallDataFile)
		}
		stringCallData := string(rawCallData)
		inputUlxlyArgs.callData = &stringCallData
	}
	return nil
}

func fatalIfError(err error) {
	if err == nil {
		return
	}
	log.Fatal().Err(err).Msg("Unexpected error occurred")
}

func init() {
	bridgeAssetCommand = &cobra.Command{
		Use:     "asset",
		Short:   "Move ETH or an ERC20 between to chains",
		Long:    bridgeAssetUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeAsset(cmd)
		},
		SilenceUsage: true,
	}
	bridgeMessageCommand = &cobra.Command{
		Use:     "message",
		Short:   "Send some ETH along with data from one chain to another chain",
		Long:    bridgeMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeMessage(cmd)
		},
		SilenceUsage: true,
	}
	bridgeMessageWETHCommand = &cobra.Command{
		Use:     "weth",
		Short:   "For L2's that use a gas token, use this to transfer WETH to another chain",
		Long:    bridgeWETHMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeWETHMessage(cmd)
		},
		SilenceUsage: true,
	}
	claimAssetCommand = &cobra.Command{
		Use:     "asset",
		Short:   "Claim a deposit",
		Long:    claimAssetUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return claimAsset(cmd)
		},
		SilenceUsage: true,
	}
	claimMessageCommand = &cobra.Command{
		Use:     "message",
		Short:   "Claim a message",
		Long:    claimMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return claimMessage(cmd)
		},
		SilenceUsage: true,
	}
	claimEverythingCommand = &cobra.Command{
		Use:     "claim-everything",
		Short:   "Attempt to claim as many deposits and messages as possible",
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return claimEverything(cmd)
		},
		SilenceUsage: true,
	}
	emptyProofCommand = &cobra.Command{
		Use:   "empty-proof",
		Short: "create an empty proof",
		Long:  "Use this command to print an empty proof response that's filled with zero-valued siblings like 0x0000000000000000000000000000000000000000000000000000000000000000. This can be useful when you need to submit a dummy proof.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return emptyProof()
		},
		SilenceUsage: true,
	}
	zeroProofCommand = &cobra.Command{
		Use:   "zero-proof",
		Short: "create a proof that's filled with zeros",
		Long: `Use this command to print a proof response that's filled with the zero
hashes. This values are very helpful for debugging because it would
tell you how populated the tree is and roughly which leaves and
siblings are empty. It's also helpful for sanity checking a proof
response to understand if the hashed value is part of the zero hashes
or if it's actually an intermediate hash.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return zeroProof()
		},
		SilenceUsage: true,
	}
	proofCommand = &cobra.Command{
		Use:   "proof",
		Short: "Generate a proof for a given range of deposits",
		Long:  proofUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return proof(args)
		},
		SilenceUsage: true,
	}
	getDepositCommand = &cobra.Command{
		Use:   "get-deposits",
		Short: "Generate ndjson for each bridge deposit over a particular range of blocks",
		Long:  depositGetUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return readDeposit(cmd)
		},
		SilenceUsage: true,
	}

	// Arguments for both bridge and claim
	inputUlxlyArgs.rpcURL = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgRPCURL, "", "the URL of the RPC to send the transaction")
	inputUlxlyArgs.bridgeAddress = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgBridgeAddress, "", "the address of the lxly bridge")
	inputUlxlyArgs.gasLimit = ulxlyBridgeAndClaimCmd.PersistentFlags().Uint64(ArgGasLimit, 0, "force a gas limit when sending a transaction")
	inputUlxlyArgs.chainID = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgChainID, "", "set the chain id to be used in the transaction")
	inputUlxlyArgs.privateKey = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgPrivateKey, "", "the hex encoded private key to be used when sending the tx")
	inputUlxlyArgs.destAddress = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgDestAddress, "", "the address where the bridge will be sent to")
	inputUlxlyArgs.timeout = ulxlyBridgeAndClaimCmd.PersistentFlags().Uint64(ArgTimeout, 60, "the amount of time to wait while trying to confirm a transaction receipt")
	inputUlxlyArgs.gasPrice = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgGasPrice, "", "the gas price to be used")
	inputUlxlyArgs.dryRun = ulxlyBridgeAndClaimCmd.PersistentFlags().Bool(ArgDryRun, false, "do all of the transaction steps but do not send the transaction")
	fatalIfError(ulxlyBridgeAndClaimCmd.MarkPersistentFlagRequired(ArgPrivateKey))
	fatalIfError(ulxlyBridgeAndClaimCmd.MarkPersistentFlagRequired(ArgRPCURL))
	fatalIfError(ulxlyBridgeAndClaimCmd.MarkPersistentFlagRequired(ArgBridgeAddress))

	// bridge specific args
	inputUlxlyArgs.forceUpdate = ulxlyBridgeCmd.PersistentFlags().Bool(ArgForceUpdate, true, "indicates if the new global exit root is updated or not")
	inputUlxlyArgs.value = ulxlyBridgeCmd.PersistentFlags().String(ArgValue, "0", "the amount in wei to be sent along with the transaction")
	inputUlxlyArgs.destNetwork = ulxlyBridgeCmd.PersistentFlags().Uint32(ArgDestNetwork, 0, "the rollup id of the destination network")
	inputUlxlyArgs.tokenAddress = ulxlyBridgeCmd.PersistentFlags().String(ArgTokenAddress, "0x0000000000000000000000000000000000000000", "the address of an ERC20 token to be used")
	inputUlxlyArgs.callData = ulxlyBridgeCmd.PersistentFlags().String(ArgCallData, "0x", "call data to be passed directly with bridge-message or as an ERC20 Permit")
	inputUlxlyArgs.callDataFile = ulxlyBridgeCmd.PersistentFlags().String(ArgCallDataFile, "", "a file containing hex encoded call data")
	fatalIfError(ulxlyBridgeCmd.MarkPersistentFlagRequired(ArgDestNetwork))

	// Claim specific args
	inputUlxlyArgs.depositCount = ulxlyClaimCmd.PersistentFlags().Uint64(ArgDepositCount, 0, "the deposit count of the bridge transaction")
	inputUlxlyArgs.depositNetwork = ulxlyClaimCmd.PersistentFlags().Uint64(ArgDepositNetwork, 0, "the rollup id of the network where the deposit was initially made")
	inputUlxlyArgs.bridgeServiceURL = ulxlyClaimCmd.PersistentFlags().String(ArgBridgeServiceURL, "", "the URL of the bridge service")
	inputUlxlyArgs.globalIndex = ulxlyClaimCmd.PersistentFlags().String(ArgGlobalIndex, "", "an override of the global index value")
	inputUlxlyArgs.wait = ulxlyClaimCmd.PersistentFlags().Duration(ArgWait, time.Duration(0), "this flag is available for claim asset and claim message. if specified, the command will retry in a loop for the deposit to be ready to claim up to duration. Once the deposit is ready to claim, the claim will actually be sent.")
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgDepositCount))
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgDepositNetwork))
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgBridgeServiceURL))

	// Claim Everything Helper Command
	inputUlxlyArgs.bridgeServiceURLs = claimEverythingCommand.Flags().StringSlice(ArgBridgeMappings, nil, "Mappings between network ids and bridge service urls. E.g. '1=http://network-1-bridgeurl,7=http://network-2-bridgeurl'")
	inputUlxlyArgs.bridgeLimit = claimEverythingCommand.Flags().Int(ArgBridgeLimit, 25, "Limit the number or responses returned by the bridge service when claiming")
	inputUlxlyArgs.bridgeOffset = claimEverythingCommand.Flags().Int(ArgBridgeOffset, 0, "The offset to specify for pagination of the underlying bridge service deposits")
	fatalIfError(claimEverythingCommand.MarkFlagRequired(ArgBridgeMappings))

	// Args that are just for the get deposit command
	inputUlxlyArgs.fromBlock = getDepositCommand.Flags().Uint64(ArgFromBlock, 0, "The start of the range of blocks to retrieve")
	inputUlxlyArgs.toBlock = getDepositCommand.Flags().Uint64(ArgToBlock, 0, "The end of the range of blocks to retrieve")
	inputUlxlyArgs.filterSize = getDepositCommand.Flags().Uint64(ArgFilterSize, 1000, "The batch size for individual filter queries")
	getDepositCommand.Flags().String(ArgRPCURL, "", "The RPC URL to read deposit data")
	getDepositCommand.Flags().String(ArgBridgeAddress, "", "The address of the ulxly bridge")
	fatalIfError(getDepositCommand.MarkFlagRequired(ArgFromBlock))
	fatalIfError(getDepositCommand.MarkFlagRequired(ArgToBlock))
	fatalIfError(getDepositCommand.MarkFlagRequired(ArgRPCURL))

	// Args for the proof command
	inputUlxlyArgs.inputFileName = proofCommand.Flags().String(ArgFileName, "", "An ndjson file with deposit data")
	inputUlxlyArgs.depositNumber = proofCommand.Flags().Uint64(ArgDepositCount, 0, "The deposit number to generate a proof for")

	// Top Level
	ULxLyCmd.AddCommand(ulxlyBridgeAndClaimCmd)
	ULxLyCmd.AddCommand(emptyProofCommand)
	ULxLyCmd.AddCommand(zeroProofCommand)
	ULxLyCmd.AddCommand(proofCommand)
	ULxLyCmd.AddCommand(getDepositCommand)

	ULxLyCmd.AddCommand(ulxlyBridgeCmd)
	ULxLyCmd.AddCommand(ulxlyClaimCmd)
	ULxLyCmd.AddCommand(claimEverythingCommand)

	// Bridge and Claim
	ulxlyBridgeAndClaimCmd.AddCommand(ulxlyBridgeCmd)
	ulxlyBridgeAndClaimCmd.AddCommand(ulxlyClaimCmd)
	ulxlyBridgeAndClaimCmd.AddCommand(claimEverythingCommand)

	// Bridge
	ulxlyBridgeCmd.AddCommand(bridgeAssetCommand)
	ulxlyBridgeCmd.AddCommand(bridgeMessageCommand)
	ulxlyBridgeCmd.AddCommand(bridgeMessageWETHCommand)

	// Claim
	ulxlyClaimCmd.AddCommand(claimAssetCommand)
	ulxlyClaimCmd.AddCommand(claimMessageCommand)
}
