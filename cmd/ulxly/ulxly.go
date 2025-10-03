package ulxly

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/0xPolygon/polygon-cli/bindings/ulxly/polygonrollupmanager"
	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	bridge_service_factory "github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/factory"
	smcerror "github.com/0xPolygon/polygon-cli/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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
	ErrNotReadyForClaim      = errors.New("the claim transaction is not yet ready to be claimed, try again in a few blocks")
	ErrDepositAlreadyClaimed = errors.New("the claim transaction has already been claimed")
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
type RollupsProof struct {
	Siblings [TreeDepth]common.Hash
	Root     common.Hash
	RollupID uint32
	LeafHash common.Hash
}

type DepositID struct {
	DepositCnt uint32 `json:"deposit_cnt"`
	NetworkID  uint32 `json:"network_id"`
}

func readDeposit(cmd *cobra.Command) error {
	bridgeAddress := getSmcOptions.BridgeAddress
	rpcUrl := getEvent.URL
	toBlock := getEvent.ToBlock
	fromBlock := getEvent.FromBlock
	filter := getEvent.FilterSize

	// Use the new helper function
	var rpc *ethrpc.Client
	var err error

	if getEvent.Insecure {
		client, clientErr := createInsecureEthClient(rpcUrl)
		if clientErr != nil {
			log.Error().Err(clientErr).Msg("Unable to create insecure client")
			return clientErr
		}
		defer client.Close()
		rpc = client.Client()
	} else {
		rpc, err = ethrpc.DialContext(cmd.Context(), rpcUrl)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer rpc.Close()
	}

	ec := ethclient.NewClient(rpc)

	bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(bridgeAddress), ec)
	if err != nil {
		return err
	}
	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := min(currentBlock+filter, toBlock)

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
		currentBlock = endBlock + 1
	}

	return nil
}

func DecodeGlobalIndex(globalIndex *big.Int) (bool, uint32, uint32, error) {
	const lengthGlobalIndexInBytes = 32
	var buf [32]byte
	gIBytes := globalIndex.FillBytes(buf[:])
	if len(gIBytes) != lengthGlobalIndexInBytes {
		return false, 0, 0, fmt.Errorf("invalid globalIndex length. Should be 32. Current length: %d", len(gIBytes))
	}
	mainnetFlag := big.NewInt(0).SetBytes([]byte{gIBytes[23]}).Uint64() == 1
	rollupIndex := big.NewInt(0).SetBytes(gIBytes[24:28])
	localRootIndex := big.NewInt(0).SetBytes(gIBytes[28:32])
	if rollupIndex.Uint64() > math.MaxUint32 {
		return false, 0, 0, fmt.Errorf("invalid rollupIndex length. Should be fit into uint32 type")
	}
	if localRootIndex.Uint64() > math.MaxUint32 {
		return false, 0, 0, fmt.Errorf("invalid localRootIndex length. Should be fit into uint32 type")
	}
	return mainnetFlag, uint32(rollupIndex.Uint64()), uint32(localRootIndex.Uint64()), nil // nolint:gosec
}

func readClaim(cmd *cobra.Command) error {
	bridgeAddress := getSmcOptions.BridgeAddress
	rpcUrl := getEvent.URL
	toBlock := getEvent.ToBlock
	fromBlock := getEvent.FromBlock
	filter := getEvent.FilterSize

	// Use the new helper function
	var rpc *ethrpc.Client
	var err error

	if getEvent.Insecure {
		client, clientErr := createInsecureEthClient(rpcUrl)
		if clientErr != nil {
			log.Error().Err(clientErr).Msg("Unable to create insecure client")
			return clientErr
		}
		defer client.Close()
		rpc = client.Client()
	} else {
		rpc, err = ethrpc.DialContext(cmd.Context(), rpcUrl)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer rpc.Close()
	}

	ec := ethclient.NewClient(rpc)

	bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(bridgeAddress), ec)
	if err != nil {
		return err
	}
	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := min(currentBlock+filter, toBlock)

		opts := bind.FilterOpts{
			Start:   currentBlock,
			End:     &endBlock,
			Context: cmd.Context(),
		}
		evtV2Iterator, err := bridgeV2.FilterClaimEvent(&opts)
		if err != nil {
			return err
		}

		for evtV2Iterator.Next() {
			evt := evtV2Iterator.Event
			var (
				mainnetFlag                     bool
				rollupIndex, localExitRootIndex uint32
			)
			mainnetFlag, rollupIndex, localExitRootIndex, err = DecodeGlobalIndex(evt.GlobalIndex)
			if err != nil {
				log.Error().Err(err).Msg("error decoding globalIndex")
				return err
			}
			log.Info().Bool("claim-mainnetFlag", mainnetFlag).Uint32("claim-RollupIndex", rollupIndex).Uint32("claim-LocalExitRootIndex", localExitRootIndex).Uint64("block-number", evt.Raw.BlockNumber).Msg("Found Claim")
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
		currentBlock = endBlock + 1
	}

	return nil
}

func readVerifyBatches(cmd *cobra.Command) error {
	rollupManagerAddress := getVerifyBatchesOptions.RollupManagerAddress
	rpcUrl := getEvent.URL
	toBlock := getEvent.ToBlock
	fromBlock := getEvent.FromBlock
	filter := getEvent.FilterSize

	// Use the new helper function
	var rpc *ethrpc.Client
	var err error

	if getEvent.Insecure {
		client, clientErr := createInsecureEthClient(rpcUrl)
		if clientErr != nil {
			log.Error().Err(clientErr).Msg("Unable to create insecure client")
			return clientErr
		}
		defer client.Close()
		rpc = client.Client()
	} else {
		rpc, err = ethrpc.DialContext(cmd.Context(), rpcUrl)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer rpc.Close()
	}

	client := ethclient.NewClient(rpc)
	rm := common.HexToAddress(rollupManagerAddress)
	rollupManager, err := polygonrollupmanager.NewPolygonrollupmanager(rm, client)
	if err != nil {
		return err
	}
	verifyBatchesTrustedAggregatorSignatureHash := crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)"))

	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := min(currentBlock+filter, toBlock)
		// Filter 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3
		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(currentBlock),
			ToBlock:   new(big.Int).SetUint64(endBlock),
			Addresses: []common.Address{rm},
			Topics:    [][]common.Hash{{verifyBatchesTrustedAggregatorSignatureHash}},
		}
		logs, err := client.FilterLogs(cmd.Context(), query)
		if err != nil {
			return err
		}

		for _, vLog := range logs {
			vb, err := rollupManager.ParseVerifyBatchesTrustedAggregator(vLog)
			if err != nil {
				return err
			}
			log.Info().Uint32("RollupID", vb.RollupID).Uint64("block-number", vb.Raw.BlockNumber).Msg("Found rollupmanager VerifyBatchesTrustedAggregator event")
			var jBytes []byte
			jBytes, err = json.Marshal(vb)
			if err != nil {
				return err
			}
			fmt.Println(string(jBytes))
		}
		currentBlock = endBlock + 1
	}

	return nil
}

func proof(args []string) error {
	depositNumber := proofOptions.DepositCount
	rawDepositData, err := getInputData(args)
	if err != nil {
		return err
	}
	return readDeposits(rawDepositData, uint32(depositNumber))
}

func balanceTree() error {
	l2NetworkID := balanceTreeOptions.L2NetworkID
	bridgeAddress := common.HexToAddress(balanceTreeOptions.BridgeAddress)

	var client *ethclient.Client
	var err error

	if balanceTreeOptions.Insecure {
		client, err = createInsecureEthClient(balanceTreeOptions.RpcURL)
	} else {
		client, err = ethclient.DialContext(context.Background(), balanceTreeOptions.RpcURL)
	}

	if err != nil {
		return err
	}
	defer client.Close()
	l2RawClaimsData, l2RawDepositsData, err := getBalanceTreeData()
	if err != nil {
		return err
	}
	root, err := computeBalanceTree(client, bridgeAddress, l2RawClaimsData, l2NetworkID, l2RawDepositsData)
	if err != nil {
		return err
	}
	fmt.Printf(`
	{
		"root": "%s"
	}
	`, root.String())
	return nil
}

func nullifierTree(args []string) error {
	rawClaims, err := getInputData(args)
	if err != nil {
		return err
	}
	root, err := computeNullifierTree(rawClaims)
	if err != nil {
		return err
	}
	fmt.Printf(`
	{
		"root": "%s"
	}
	`, root.String())
	return nil
}

func nullifierAndBalanceTree(args []string) error {
	l2NetworkID := balanceTreeOptions.L2NetworkID
	bridgeAddress := common.HexToAddress(balanceTreeOptions.BridgeAddress)

	var client *ethclient.Client
	var err error

	if balanceTreeOptions.Insecure {
		client, err = createInsecureEthClient(balanceTreeOptions.RpcURL)
	} else {
		client, err = ethclient.DialContext(context.Background(), balanceTreeOptions.RpcURL)
	}

	if err != nil {
		return err
	}
	defer client.Close()
	l2RawClaimsData, l2RawDepositsData, err := getBalanceTreeData()
	if err != nil {
		return err
	}
	bridgeV2, err := ulxly.NewUlxly(bridgeAddress, client)
	if err != nil {
		return err
	}
	ler_count, err := bridgeV2.LastUpdatedDepositCount(&bind.CallOpts{Pending: false})
	if err != nil {
		return err
	}
	log.Info().Msgf("Last LER count: %d", ler_count)
	balanceTreeRoot, err := computeBalanceTree(client, bridgeAddress, l2RawClaimsData, l2NetworkID, l2RawDepositsData)
	if err != nil {
		return err
	}
	nullifierTreeRoot, err := computeNullifierTree(l2RawClaimsData)
	if err != nil {
		return err
	}
	initPessimisticRoot := crypto.Keccak256Hash(balanceTreeRoot.Bytes(), nullifierTreeRoot.Bytes(), Uint32ToBytesLittleEndian(ler_count))
	fmt.Printf(`
	{
		"balanceTreeRoot": "%s",
		"nullifierTreeRoot": "%s",
		"initPessimisticRoot": "%s"
	}
	`, balanceTreeRoot.String(), nullifierTreeRoot.String(), initPessimisticRoot.String())
	return nil
}

func computeNullifierTree(rawClaims []byte) (common.Hash, error) {
	buf := bytes.NewBuffer(rawClaims)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	nTree, err := NewNullifierTree()
	if err != nil {
		return common.Hash{}, err
	}
	var root common.Hash
	for scanner.Scan() {
		claim := new(ulxly.UlxlyClaimEvent)
		err := json.Unmarshal(scanner.Bytes(), claim)
		if err != nil {
			return common.Hash{}, err
		}
		mainnetFlag, rollupIndex, localExitRootIndex, err := DecodeGlobalIndex(claim.GlobalIndex)
		if err != nil {
			log.Error().Err(err).Msg("error decoding globalIndex")
			return common.Hash{}, err
		}
		log.Info().Bool("MainnetFlag", mainnetFlag).Uint32("RollupIndex", rollupIndex).Uint32("LocalExitRootIndex", localExitRootIndex).Uint64("block-number", claim.Raw.BlockNumber).Msg("Adding Claim")
		nullifierKey := NullifierKey{
			NetworkID: claim.OriginNetwork,
			Index:     localExitRootIndex,
		}
		root, err = nTree.UpdateNullifierTree(nullifierKey)
		if err != nil {
			log.Error().Err(err).Uint32("OriginNetwork: ", claim.OriginNetwork).Msg("error computing nullifierTree. Claim information: GlobalIndex: " + claim.GlobalIndex.String() + ", OriginAddress: " + claim.OriginAddress.String() + ", Amount: " + claim.Amount.String())
			return common.Hash{}, err
		}
	}
	log.Info().Msgf("Final nullifierTree root: %s", root.String())
	return root, nil
}

func computeBalanceTree(client *ethclient.Client, bridgeAddress common.Address, l2RawClaims []byte, l2NetworkID uint32, l2RawDeposits []byte) (common.Hash, error) {
	buf := bytes.NewBuffer(l2RawClaims)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	bTree, err := NewBalanceTree()
	if err != nil {
		return common.Hash{}, err
	}
	balances := make(map[string]*big.Int)
	for scanner.Scan() {
		l2Claim := new(ulxly.UlxlyClaimEvent)
		err := json.Unmarshal(scanner.Bytes(), l2Claim)
		if err != nil {
			return common.Hash{}, err
		}
		token := TokenInfo{
			OriginNetwork:      big.NewInt(0).SetUint64(uint64(l2Claim.OriginNetwork)),
			OriginTokenAddress: l2Claim.OriginAddress,
		}
		isMessage, err := checkClaimCalldata(client, bridgeAddress, l2Claim.Raw.TxHash)
		if err != nil {
			return common.Hash{}, err
		}
		if isMessage {
			token.OriginNetwork = big.NewInt(0)
			token.OriginTokenAddress = common.Address{}
		}
		log.Info().Msgf("L2 Claim. isMessage: %v OriginNetwork: %d. TokenAddress: %s. Amount: %s", isMessage, token.OriginNetwork, token.OriginTokenAddress.String(), l2Claim.Amount.String())
		if _, ok := balances[token.String()]; !ok {
			balances[token.String()] = big.NewInt(0)
		}
		balances[token.String()] = big.NewInt(0).Add(balances[token.String()], l2Claim.Amount)

	}
	l2Buf := bytes.NewBuffer(l2RawDeposits)
	l2Scanner := bufio.NewScanner(l2Buf)
	l2ScannerBuf := make([]byte, 0)
	l2Scanner.Buffer(l2ScannerBuf, 1024*1024)
	for l2Scanner.Scan() {
		l2Deposit := new(ulxly.UlxlyBridgeEvent)
		err := json.Unmarshal(l2Scanner.Bytes(), l2Deposit)
		if err != nil {
			return common.Hash{}, err
		}
		token := TokenInfo{
			OriginNetwork:      big.NewInt(0).SetUint64(uint64(l2Deposit.OriginNetwork)),
			OriginTokenAddress: l2Deposit.OriginAddress,
		}
		if _, ok := balances[token.String()]; !ok {
			balances[token.String()] = big.NewInt(0)
		}
		balances[token.String()] = big.NewInt(0).Sub(balances[token.String()], l2Deposit.Amount)
	}
	// Now, the balance map is complete. Let's build the tree.
	var root common.Hash
	for t, balance := range balances {
		if balance.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		token, err := TokenInfoStringToStruct(t)
		if err != nil {
			return common.Hash{}, err
		}
		if token.OriginNetwork.Uint64() == uint64(l2NetworkID) {
			continue
		}
		root, err = bTree.UpdateBalanceTree(token, balance)
		if err != nil {
			return common.Hash{}, err
		}
		log.Info().Msgf("New balanceTree leaf. OriginNetwork: %s, TokenAddress: %s, Balance: %s, Root: %s", token.OriginNetwork.String(), token.OriginTokenAddress.String(), balance.String(), root.String())
	}
	log.Info().Msgf("Final balanceTree root: %s", root.String())

	return root, nil
}

func rollupsExitRootProof(args []string) error {
	rollupID := rollupsProofOptions.RollupID
	completeMT := rollupsProofOptions.CompleteMerkleTree
	rawLeavesData, err := getInputData(args)
	if err != nil {
		return err
	}
	return readRollupsExitRootLeaves(rawLeavesData, rollupID, completeMT)
}

func emptyProof() error {
	p := new(Proof)

	e := generateEmptyHashes(TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(String(p))
	return nil
}

func zeroProof() error {
	p := new(Proof)

	e := generateZeroHashes(TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(String(p))
	return nil
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func logAndReturnJsonError(ctx context.Context, client *ethclient.Client, tx *types.Transaction, opts *bind.TransactOpts, err error) error {

	var callErr error
	if tx != nil {
		// in case the error came down to gas estimation, we can sometimes get more information by doing a call
		_, callErr = client.CallContract(ctx, ethereum.CallMsg{
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

		if inputUlxlyArgs.dryRun {
			castCmd := "cast call"
			castCmd += fmt.Sprintf(" --rpc-url %s", inputUlxlyArgs.rpcURL)
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

	reason, decodeErr := smcerror.DecodeSmcErrorCode(jsonError.Data)
	if decodeErr != nil {
		log.Error().Err(err).Msg("unable to decode smart contract error")
		return err
	}
	errLog := log.Error().
		Err(err).
		Str("message", jsonError.Message).
		Int("code", jsonError.Code).
		Interface("data", jsonError.Data).
		Str("reason", reason)

	if callErr != nil {
		errLog = errLog.Err(callErr)
	}

	customErr := errors.New(err.Error() + ": " + reason)
	if errCode, isValid := jsonError.Data.(string); isValid && errCode == "0x646cf558" {
		// I don't want to bother with the additional error logging for previously claimed deposits
		return customErr
	}

	errLog.Msg("Unable to interact with bridge contract")
	return customErr
}

// Function to parse deposit count from bridge transaction logs
func ParseBridgeDepositCount(logs []types.Log, bridgeContract *ulxly.Ulxly) (uint32, error) {
	for _, log := range logs {
		// Try to parse the log as a BridgeEvent using the contract's filterer
		bridgeEvent, err := bridgeContract.ParseBridgeEvent(log)
		if err != nil {
			// This log is not a bridge event, continue to next log
			continue
		}

		// Successfully parsed a bridge event, return the deposit count
		return bridgeEvent.DepositCount, nil
	}

	return 0, fmt.Errorf("bridge event not found in logs")
}

// parseDepositCountFromTransaction extracts the deposit count from a bridge transaction receipt
func parseDepositCountFromTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash, bridgeContract *ulxly.Ulxly) (uint32, error) {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return 0, err
	}

	// Check if the transaction was successful before trying to parse logs
	if receipt.Status == 0 {
		log.Error().Str("txHash", receipt.TxHash.String()).Msg("Bridge transaction failed")
		return 0, fmt.Errorf("bridge transaction failed with hash: %s", receipt.TxHash.String())
	}

	// Convert []*types.Log to []types.Log
	logs := make([]types.Log, len(receipt.Logs))
	for i, log := range receipt.Logs {
		logs[i] = *log
	}

	depositCount, err := ParseBridgeDepositCount(logs, bridgeContract)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse deposit count from logs")
		return 0, err
	}

	return depositCount, nil
}

func bridgeAsset(cmd *cobra.Command) error {
	bridgeAddr := inputUlxlyArgs.bridgeAddress
	privateKey := inputUlxlyArgs.privateKey
	gasLimit := inputUlxlyArgs.gasLimit
	destinationAddress := inputUlxlyArgs.destAddress
	chainID := inputUlxlyArgs.chainID
	amount := inputUlxlyArgs.value
	tokenAddr := inputUlxlyArgs.tokenAddress
	callDataString := inputUlxlyArgs.callData
	destinationNetwork := inputUlxlyArgs.destNetwork
	isForced := inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	RPCURL := inputUlxlyArgs.rpcURL

	client, err := createEthClient(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()

	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddr, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	bridgeAddress := common.HexToAddress(bridgeAddr)
	value, _ := big.NewInt(0).SetString(amount, 0)
	tokenAddress := common.HexToAddress(tokenAddr)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	} else {
		// in case it's a token transfer, we need to ensure that the bridge contract
		// has enough allowance to transfer the tokens on behalf of the user
		tokenContract, iErr := tokens.NewERC20(tokenAddress, client)
		if iErr != nil {
			log.Error().Err(iErr).Msg("error getting token contract")
			return iErr
		}

		allowance, iErr := tokenContract.Allowance(&bind.CallOpts{Pending: false}, auth.From, bridgeAddress)
		if iErr != nil {
			log.Error().Err(iErr).Msg("error getting token allowance")
			return iErr
		}

		if allowance.Cmp(value) < 0 {
			log.Info().
				Str("amount", value.String()).
				Str("tokenAddress", tokenAddress.String()).
				Str("bridgeAddress", bridgeAddress.String()).
				Str("userAddress", auth.From.String()).
				Msg("approving bridge contract to spend tokens on behalf of user")

			// Approve the bridge contract to spend the tokens on behalf of the user
			approveTxn, iErr := tokenContract.Approve(auth, bridgeAddress, value)
			if iErr = logAndReturnJsonError(cmd.Context(), client, approveTxn, auth, iErr); iErr != nil {
				return iErr
			}
			log.Info().Msg("approveTxn: " + approveTxn.Hash().String())
			if iErr = WaitMineTransaction(cmd.Context(), client, approveTxn, timeoutTxnReceipt); iErr != nil {
				return iErr
			}
		}
	}

	bridgeTxn, err := bridgeV2.BridgeAsset(auth, destinationNetwork, toAddress, value, tokenAddress, isForced, callData)
	if err = logAndReturnJsonError(cmd.Context(), client, bridgeTxn, auth, err); err != nil {
		log.Info().Err(err).Str("calldata", callDataString).Msg("Bridge transaction failed")
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	if err = WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt); err != nil {
		return err
	}
	depositCount, err := parseDepositCountFromTransaction(cmd.Context(), client, bridgeTxn.Hash(), bridgeV2)
	if err != nil {
		return err
	}

	log.Info().Uint32("depositCount", depositCount).Msg("Bridge deposit count parsed from logs")
	return nil
}

func bridgeMessage(cmd *cobra.Command) error {
	bridgeAddress := inputUlxlyArgs.bridgeAddress
	privateKey := inputUlxlyArgs.privateKey
	gasLimit := inputUlxlyArgs.gasLimit
	destinationAddress := inputUlxlyArgs.destAddress
	chainID := inputUlxlyArgs.chainID
	amount := inputUlxlyArgs.value
	tokenAddr := inputUlxlyArgs.tokenAddress
	callDataString := inputUlxlyArgs.callData
	destinationNetwork := inputUlxlyArgs.destNetwork
	isForced := inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	RPCURL := inputUlxlyArgs.rpcURL

	// Dial the Ethereum RPC server.
	client, err := createEthClient(cmd.Context(), RPCURL)
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
	if err = logAndReturnJsonError(cmd.Context(), client, bridgeTxn, auth, err); err != nil {
		log.Info().Err(err).Str("calldata", callDataString).Msg("Bridge transaction failed")
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	if err = WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt); err != nil {
		return err
	}
	depositCount, err := parseDepositCountFromTransaction(cmd.Context(), client, bridgeTxn.Hash(), bridgeV2)
	if err != nil {
		return err
	}

	log.Info().Uint32("depositCount", depositCount).Msg("Bridge deposit count parsed from logs")
	return nil
}

func bridgeWETHMessage(cmd *cobra.Command) error {
	bridgeAddress := inputUlxlyArgs.bridgeAddress
	privateKey := inputUlxlyArgs.privateKey
	gasLimit := inputUlxlyArgs.gasLimit
	destinationAddress := inputUlxlyArgs.destAddress
	chainID := inputUlxlyArgs.chainID
	amount := inputUlxlyArgs.value
	callDataString := inputUlxlyArgs.callData
	destinationNetwork := inputUlxlyArgs.destNetwork
	isForced := inputUlxlyArgs.forceUpdate
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	RPCURL := inputUlxlyArgs.rpcURL

	// Dial the Ethereum RPC server.
	client, err := createEthClient(cmd.Context(), RPCURL)
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
	if err = logAndReturnJsonError(cmd.Context(), client, bridgeTxn, auth, err); err != nil {
		log.Info().Err(err).Str("calldata", callDataString).Msg("Bridge transaction failed")
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	if err = WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt); err != nil {
		return err
	}
	depositCount, err := parseDepositCountFromTransaction(cmd.Context(), client, bridgeTxn.Hash(), bridgeV2)
	if err != nil {
		return err
	}

	log.Info().Uint32("depositCount", depositCount).Msg("Bridge deposit count parsed from logs")
	return nil
}

func claimAsset(cmd *cobra.Command) error {
	bridgeAddress := inputUlxlyArgs.bridgeAddress
	privateKey := inputUlxlyArgs.privateKey
	gasLimit := inputUlxlyArgs.gasLimit
	destinationAddress := inputUlxlyArgs.destAddress
	chainID := inputUlxlyArgs.chainID
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	RPCURL := inputUlxlyArgs.rpcURL
	depositCount := inputUlxlyArgs.depositCount
	depositNetwork := inputUlxlyArgs.depositNetwork
	globalIndexOverride := inputUlxlyArgs.globalIndex
	proofGERHash := inputUlxlyArgs.proofGER
	wait := inputUlxlyArgs.wait

	// Dial Ethereum client
	client, err := createEthClient(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, _, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
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

	proof, err := getMerkleProofsExitRoots(bridgeService, *deposit, proofGERHash)
	if err != nil {
		log.Error().Err(err).Msg("error getting merkle proofs and exit roots from bridge service")
		return err
	}

	claimTxn, err := bridgeV2.ClaimAsset(auth, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	if err = logAndReturnJsonError(cmd.Context(), client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}

func claimMessage(cmd *cobra.Command) error {
	bridgeAddress := inputUlxlyArgs.bridgeAddress
	privateKey := inputUlxlyArgs.privateKey
	gasLimit := inputUlxlyArgs.gasLimit
	destinationAddress := inputUlxlyArgs.destAddress
	chainID := inputUlxlyArgs.chainID
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	RPCURL := inputUlxlyArgs.rpcURL
	depositCount := inputUlxlyArgs.depositCount
	depositNetwork := inputUlxlyArgs.depositNetwork
	globalIndexOverride := inputUlxlyArgs.globalIndex
	proofGERHash := inputUlxlyArgs.proofGER
	wait := inputUlxlyArgs.wait

	// Dial Ethereum client
	client, err := createEthClient(cmd.Context(), RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()
	// Initialize and assign variables required to send transaction payload
	bridgeV2, _, auth, err := generateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
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

	proof, err := getMerkleProofsExitRoots(bridgeService, *deposit, proofGERHash)
	if err != nil {
		log.Error().Err(err).Msg("error getting merkle proofs and exit roots from bridge service")
		return err
	}

	claimTxn, err := bridgeV2.ClaimMessage(auth, bridge_service.HashSliceToBytesArray(proof.MerkleProof), bridge_service.HashSliceToBytesArray(proof.RollupMerkleProof), deposit.GlobalIndex, *proof.MainExitRoot, *proof.RollupExitRoot, deposit.OrigNet, deposit.OrigAddr, deposit.DestNet, deposit.DestAddr, deposit.Amount, deposit.Metadata)
	if err = logAndReturnJsonError(cmd.Context(), client, claimTxn, auth, err); err != nil {
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
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

func getBridgeServiceURLs() (map[uint32]string, error) {
	bridgeServiceUrls := inputUlxlyArgs.bridgeServiceURLs
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
	privateKey := inputUlxlyArgs.privateKey
	claimerAddress := inputUlxlyArgs.addressOfPrivateKey

	gasLimit := inputUlxlyArgs.gasLimit
	chainID := inputUlxlyArgs.chainID
	timeoutTxnReceipt := inputUlxlyArgs.timeout
	bridgeAddress := inputUlxlyArgs.bridgeAddress
	destinationAddress := inputUlxlyArgs.destAddress
	RPCURL := inputUlxlyArgs.rpcURL
	limit := inputUlxlyArgs.bridgeLimit
	offset := inputUlxlyArgs.bridgeOffset
	concurrency := inputUlxlyArgs.concurrency

	depositMap := make(map[DepositID]*bridge_service.Deposit)

	for networkID, bridgeService := range bridgeServices {
		deposits, _, bErr := getDepositsForAddress(bridgeService, destinationAddress, offset, limit)
		if bErr != nil {
			log.Err(bErr).Uint32("id", networkID).Str("url", bridgeService.Url()).Msgf("Error getting deposits for bridge: %s", bErr.Error())
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
			if inputUlxlyArgs.legacy {
				if deposit.ReadyForClaim || deposit.ClaimTxHash != nil {
					depositMap[depId] = &deposits[idx]
				}
			}
		}
	}

	client, err := createEthClient(cmd.Context(), RPCURL)
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

			claimTx, dErr := claimSingleDeposit(cmd, client, bridgeContract, withNonce(opts, nextNonce), *deposit, bridgeServices, currentNetworkID)
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
			dErr = WaitMineTransaction(cmd.Context(), client, claimTx, timeoutTxnReceipt)
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

func claimSingleDeposit(cmd *cobra.Command, client *ethclient.Client, bridgeContract *ulxly.Ulxly, opts *bind.TransactOpts, deposit bridge_service.Deposit, bridgeServices map[uint32]bridge_service.BridgeService, currentNetworkID uint32) (*types.Transaction, error) {
	networkIDForBridgeService := deposit.NetworkID
	if deposit.NetworkID == 0 {
		networkIDForBridgeService = currentNetworkID
	}

	bridgeServiceFromMap, hasKey := bridgeServices[networkIDForBridgeService]
	if !hasKey {
		return nil, fmt.Errorf("we don't have a bridge service url for network: %d", deposit.DestNet)
	}

	proof, err := getMerkleProofsExitRoots(bridgeServiceFromMap, deposit, "")
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

	if err = logAndReturnJsonError(cmd.Context(), client, claimTx, opts, err); err != nil {
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

// Wait for the transaction to be mined
func WaitMineTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, txTimeout uint64) error {
	if inputUlxlyArgs.dryRun {
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
				log.Info().Interface("txHash", r.TxHash).Msg("transaction successful")
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
	fileName := fileOptions.FileName
	if fileName != "" {
		return os.ReadFile(fileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

func getBalanceTreeData() ([]byte, []byte, error) {
	claimsFileName := balanceTreeOptions.L2ClaimsFile
	file, err := os.Open(claimsFileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // Ensure the file is closed after reading

	// Read the entire file content
	l2Claims, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	l2FileName := balanceTreeOptions.L2DepositsFile
	file2, err := os.Open(l2FileName)
	if err != nil {
		return nil, nil, err
	}
	defer file2.Close() // Ensure the file is closed after reading

	// Read the entire file content
	l2Deposits, err := io.ReadAll(file2)
	if err != nil {
		return nil, nil, err
	}
	return l2Claims, l2Deposits, nil
}

func readRollupsExitRootLeaves(rawLeaves []byte, rollupID uint32, completeMT bool) error {
	buf := bytes.NewBuffer(rawLeaves)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	leaves := make(map[uint32]*polygonrollupmanager.PolygonrollupmanagerVerifyBatchesTrustedAggregator, 0)
	highestRollupID := uint32(0)
	for scanner.Scan() {
		evt := new(polygonrollupmanager.PolygonrollupmanagerVerifyBatchesTrustedAggregator)
		err := json.Unmarshal(scanner.Bytes(), evt)
		if err != nil {
			return err
		}
		if highestRollupID < evt.RollupID {
			highestRollupID = evt.RollupID
		}
		leaves[evt.RollupID] = evt
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("there was an error reading the deposit file")
		return err
	}
	if rollupID > highestRollupID && !completeMT {
		return fmt.Errorf("rollupID %d required is higher than the highest rollupID %d provided in the file. Please use --complete-merkle-tree option if you know what you are doing", rollupID, highestRollupID)
	} else if completeMT {
		highestRollupID = rollupID
	}
	var ls []common.Hash
	var i uint32 = 0
	for ; i <= highestRollupID; i++ {
		var exitRoot common.Hash
		if leaf, exists := leaves[i]; exists {
			exitRoot = leaf.ExitRoot
			log.Info().
				Uint64("block-number", leaf.Raw.BlockNumber).
				Uint32("rollupID", leaf.RollupID).
				Str("exitRoot", exitRoot.String()).
				Str("tx-hash", leaf.Raw.TxHash.String()).
				Msg("latest event received for the tree")
		} else {
			log.Warn().Uint32("rollupID", i).Msg("No event found for this rollup")
		}
		ls = append(ls, exitRoot)
	}
	p, err := ComputeSiblings(rollupID, ls, TreeDepth)
	if err != nil {
		return err
	}
	log.Info().Str("root", p.Root.String()).Msg("finished")
	fmt.Println(String(p))
	return nil
}

func ComputeSiblings(rollupID uint32, leaves []common.Hash, height uint8) (*RollupsProof, error) {
	initLeaves := leaves
	var ns [][][]byte
	if len(leaves) == 0 {
		leaves = append(leaves, common.Hash{})
	}
	currentZeroHashHeight := common.Hash{}
	var siblings []common.Hash
	index := rollupID
	for h := uint8(0); h < height; h++ {
		if len(leaves)%2 == 1 {
			leaves = append(leaves, currentZeroHashHeight)
		}
		if index%2 == 1 { //If it is odd
			siblings = append(siblings, leaves[index-1])
		} else { // It is even
			if len(leaves) > 1 {
				siblings = append(siblings, leaves[index+1])
			}
		}
		var (
			nsi    [][][]byte
			hashes []common.Hash
		)
		for i := 0; i < len(leaves); i += 2 {
			var left, right = i, i + 1
			hash := crypto.Keccak256Hash(leaves[left][:], leaves[right][:])
			nsi = append(nsi, [][]byte{hash[:], leaves[left][:], leaves[right][:]})
			hashes = append(hashes, hash)
		}
		// Find the index of the leave in the next level of the tree.
		// Divide the index by 2 to find the position in the upper level
		index = uint32(float64(index) / 2) //nolint:gomnd
		ns = nsi
		leaves = hashes
		currentZeroHashHeight = crypto.Keccak256Hash(currentZeroHashHeight.Bytes(), currentZeroHashHeight.Bytes())
	}
	if len(ns) != 1 {
		return nil, fmt.Errorf("error: more than one root detected: %+v", ns)
	}
	if len(siblings) != TreeDepth {
		return nil, fmt.Errorf("error: invalid number of siblings: %+v", siblings)
	}
	if leaves[0] != common.BytesToHash(ns[0][0]) {
		return nil, fmt.Errorf("latest leave (root of the tree) does not match with the root (ns[0][0])")
	}
	sb := [32]common.Hash{}
	for i := range TreeDepth {
		sb[i] = siblings[i]
	}
	p := &RollupsProof{
		Siblings: sb,
		RollupID: rollupID,
		LeafHash: initLeaves[rollupID],
		Root:     common.BytesToHash(ns[0][0]),
	}

	computedRoot := computeRoot(p.LeafHash, p.Siblings, p.RollupID, TreeDepth)
	if computedRoot != p.Root {
		return nil, fmt.Errorf("error: computed root does not match the expected root")
	}

	return p, nil
}

func computeRoot(leafHash common.Hash, smtProof [32]common.Hash, index uint32, height uint8) common.Hash {
	var node common.Hash
	copy(node[:], leafHash[:])

	// Check merkle proof
	var h uint8
	for h = 0; h < height; h++ {
		if ((index >> h) & 1) == 1 {
			node = crypto.Keccak256Hash(smtProof[h].Bytes(), node.Bytes())
		} else {
			node = crypto.Keccak256Hash(node.Bytes(), smtProof[h].Bytes())
		}
	}
	return common.BytesToHash(node[:])
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
		leaf := hashDeposit(evt)
		log.Debug().Str("leaf-hash", common.Bytes2Hex(leaf[:])).Msg("Leaf hash calculated")
		imt.AddLeaf(leaf, evt.DepositCount)
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
	fmt.Println(String(p))
	return nil
}

func ensureCodePresent(ctx context.Context, client *ethclient.Client, address string) error {
	code, err := client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("error getting code at address")
		return err
	}
	if len(code) == 0 {
		return fmt.Errorf("address %s has no code", address)
	}
	return nil
}

// String will create the json representation of the proof
func String[T any](p T) string {
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
func (s *IMT) AddLeaf(leaf common.Hash, position uint32) {
	// just keep a copy of the leaf indexed by deposit count for now
	s.Leaves[position] = leaf

	node := leaf
	size := uint64(position) + 1

	// copy the previous set of branches as a starting point. We're going to make copies of the branches at each deposit
	branches := make([]common.Hash, TreeDepth)
	if position == 0 {
		branches = generateEmptyHashes(TreeDepth)
	} else {
		copy(branches, s.Branches[position-1])
	}

	for height := uint64(0); height < TreeDepth; height += 1 {
		if ((size >> height) & 1) == 1 {
			copy(branches[height][:], node[:])
			break
		}
		node = crypto.Keccak256Hash(branches[height][:], node[:])
	}
	s.Branches[position] = branches
	s.Roots = append(s.Roots, s.GetRoot(position))
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
		siblingDepositNum := getSiblingLeafNumber(depositNum, uint32(height))
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

	r, err := Check(s.Roots, p.LeafHash, p.DepositCount, p.Siblings)
	if err != nil {
		log.Error().Err(err).Msg("failed to validate proof")
	}
	p.Root = r
	s.Proofs[depositNum] = *p
	return *p
}

// getSiblingLeafNumber returns the sibling number of a given number at a specified level in an incremental Merkle tree.
//
// In an incremental Merkle tree, each node has a sibling node at each level of the tree.
// The sibling node can be determined by flipping the bit at the current level and setting all bits to the right of the current level to 1.
// This function calculates the sibling number based on the deposit number and the specified level.
//
// Parameters:
// - LeafNumber: the original number for which the sibling is to be found.
// - level: the level in the Merkle tree at which to find the sibling.
//
// The logic works as follows:
// 1. `1 << level` creates a binary number with a single 1 bit at the position corresponding to the level.
// 2. `LeafNumber ^ (1 << level)` flips the bit at the position corresponding to the level in the LeafNumber.
// 3. `(1 << level) - 1` creates a binary number with all bits to the right of the current level set to 1.
// 4. `| ((1 << level) - 1)` ensures that all bits to the right of the current level are set to 1 in the result.
//
// The function effectively finds the sibling deposit number at each level of the Merkle tree by manipulating the bits accordingly.
func getSiblingLeafNumber(leafNumber, level uint32) uint32 {
	return leafNumber ^ (1 << level) | ((1 << level) - 1)
}

// Check is a sanity check of a proof in order to make sure that the
// proof that was generated creates a root that we recognize. This was
// useful while testing in order to avoid verifying that the proof
// works or doesn't work onchain
func Check(roots []common.Hash, leaf common.Hash, position uint32, siblings [32]common.Hash) (common.Hash, error) {
	node := leaf
	index := position
	for height := 0; height < TreeDepth; height++ {
		if ((index >> height) & 1) == 1 {
			node = crypto.Keccak256Hash(siblings[height][:], node[:])
		} else {
			node = crypto.Keccak256Hash(node[:], siblings[height][:])
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
		Uint32("leaf-position", position).
		Str("leaf-hash", leaf.String()).
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
	// checks if bridge address has code
	err = ensureCodePresent(ctx, client, ulxlyInputArgBridge)
	if err != nil {
		err = fmt.Errorf("bridge code check err: %w", err)
		return
	}

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
	if inputUlxlyArgs.gasPrice != "" {
		gasPrice := new(big.Int)
		gasPrice.SetString(inputUlxlyArgs.gasPrice, 10)
		opts.GasPrice = gasPrice
	}
	if inputUlxlyArgs.dryRun {
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

func getMerkleProofsExitRoots(bridgeService bridge_service.BridgeService, deposit bridge_service.Deposit, proofGERHash string) (*bridge_service.Proof, error) {
	var ger *common.Hash
	if len(proofGERHash) > 0 {
		hash := common.HexToHash(proofGERHash)
		ger = &hash
	}

	proof, err := bridgeService.GetProof(deposit.NetworkID, deposit.DepositCnt, ger)
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

func getDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	deposit, err := bridgeService.GetDeposit(depositNetwork, depositCount)
	if err != nil {
		return nil, err
	}

	if inputUlxlyArgs.legacy {
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

// Add the helper function to create an insecure client
func createInsecureEthClient(rpcURL string) (*ethclient.Client, error) {
	// WARNING: This disables TLS certificate verification
	log.Warn().Msg("WARNING: TLS certificate verification is disabled. This is unsafe for production use.")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rpcClient, err := ethrpc.DialOptions(context.Background(), rpcURL, ethrpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}

// Add helper function to create either secure or insecure client based on flag
func createEthClient(ctx context.Context, rpcURL string) (*ethclient.Client, error) {
	if inputUlxlyArgs.insecure {
		return createInsecureEthClient(rpcURL)
	}
	return ethclient.DialContext(ctx, rpcURL)
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

//go:embed rollupsProofUsage.md
var rollupsProofUsage string

//go:embed balanceTreeUsage.md
var balanceTreeUsage string

//go:embed nullifierAndBalanceTreeUsage.md
var nullifierAndBalanceTreeUsage string

//go:embed nullifierTreeUsage.md
var nullifierTreeUsage string

//go:embed depositGetUsage.md
var depositGetUsage string

//go:embed claimGetUsage.md
var claimGetUsage string

//go:embed verifyBatchesGetUsage.md
var verifyBatchesGetUsage string

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the uLxLy bridge",
	Long:  "Basic utility commands for interacting with the bridge contracts, bridge services, and generating proofs",
	Args:  cobra.NoArgs,
}
var ulxlyBridgeAndClaimCmd = &cobra.Command{
	Args:   cobra.NoArgs,
	Hidden: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		rpcURL, err := flag_loader.GetRequiredRpcUrlFlagValue(cmd)
		if err != nil {
			return err
		}
		if rpcURL != nil {
			inputUlxlyArgs.rpcURL = *rpcURL
		}

		privateKey, err := flag_loader.GetRequiredPrivateKeyFlagValue(cmd)
		if err != nil {
			return err
		}
		if privateKey != nil {
			inputUlxlyArgs.privateKey = *privateKey
		}
		return nil
	},
}

var ulxlyGetEventsCmd = &cobra.Command{
	Args:   cobra.NoArgs,
	Hidden: true,
}

var ulxlyProofsCmd = &cobra.Command{
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
	gasLimit            uint64
	chainID             string
	privateKey          string
	addressOfPrivateKey string
	value               string
	rpcURL              string
	bridgeAddress       string
	destNetwork         uint32
	destAddress         string
	tokenAddress        string
	forceUpdate         bool
	callData            string
	callDataFile        string
	timeout             uint64
	depositCount        uint32
	depositNetwork      uint32
	bridgeServiceURL    string
	globalIndex         string
	gasPrice            string
	dryRun              bool
	bridgeServiceURLs   []string
	bridgeLimit         int
	bridgeOffset        int
	wait                time.Duration
	concurrency         uint
	insecure            bool
	legacy              bool
	proofGER            string
}

var inputUlxlyArgs = ulxlyArgs{}

var (
	bridgeAssetCommand             *cobra.Command
	bridgeMessageCommand           *cobra.Command
	bridgeMessageWETHCommand       *cobra.Command
	claimAssetCommand              *cobra.Command
	claimMessageCommand            *cobra.Command
	claimEverythingCommand         *cobra.Command
	emptyProofCommand              *cobra.Command
	zeroProofCommand               *cobra.Command
	proofCommand                   *cobra.Command
	rollupsProofCommand            *cobra.Command
	balanceTreeCommand             *cobra.Command
	nullifierAndBalanceTreeCommand *cobra.Command
	nullifierTreeCommand           *cobra.Command
	getDepositCommand              *cobra.Command
	getClaimCommand                *cobra.Command
	getVerifyBatchesCommand        *cobra.Command

	getEvent                = &GetEvent{}
	getSmcOptions           = &GetSmcOptions{}
	getVerifyBatchesOptions = &GetVerifyBatchesOptions{}
	fileOptions             = &FileOptions{}
	balanceTreeOptions      = &BalanceTreeOptions{}
	proofOptions            = &ProofOptions{}
	rollupsProofOptions     = &RollupsProofOptions{}
)

const (
	ArgGasLimit             = "gas-limit"
	ArgChainID              = "chain-id"
	ArgPrivateKey           = "private-key"
	ArgValue                = "value"
	ArgRPCURL               = "rpc-url"
	ArgBridgeAddress        = "bridge-address"
	ArgRollupManagerAddress = "rollup-manager-address"
	ArgDestNetwork          = "destination-network"
	ArgDestAddress          = "destination-address"
	ArgForceUpdate          = "force-update-root"
	ArgCallData             = "call-data"
	ArgCallDataFile         = "call-data-file"
	ArgTimeout              = "transaction-receipt-timeout"
	ArgDepositCount         = "deposit-count"
	ArgDepositNetwork       = "deposit-network"
	ArgRollupID             = "rollup-id"
	ArgCompleteMT           = "complete-merkle-tree"
	ArgBridgeServiceURL     = "bridge-service-url"
	ArgFileName             = "file-name"
	ArgL2ClaimsFileName     = "l2-claims-file"
	ArgL2DepositsFileName   = "l2-deposits-file"
	ArgL2NetworkID          = "l2-network-id"
	ArgFromBlock            = "from-block"
	ArgToBlock              = "to-block"
	ArgFilterSize           = "filter-size"
	ArgTokenAddress         = "token-address"
	ArgGlobalIndex          = "global-index"
	ArgDryRun               = "dry-run"
	ArgGasPrice             = "gas-price"
	ArgBridgeMappings       = "bridge-service-map"
	ArgBridgeLimit          = "bridge-limit"
	ArgBridgeOffset         = "bridge-offset"
	ArgWait                 = "wait"
	ArgConcurrency          = "concurrency"
	ArgInsecure             = "insecure"
	ArgLegacy               = "legacy"
	ArgProofGER             = "proof-ger"
)

var (
	bridgeService  bridge_service.BridgeService
	bridgeServices map[uint32]bridge_service.BridgeService = make(map[uint32]bridge_service.BridgeService)
)

func prepInputs(cmd *cobra.Command, args []string) error {
	if inputUlxlyArgs.dryRun && inputUlxlyArgs.gasLimit == 0 {
		inputUlxlyArgs.gasLimit = uint64(10_000_000)
	}
	pvtKey := strings.TrimPrefix(inputUlxlyArgs.privateKey, "0x")
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
	if inputUlxlyArgs.destAddress == "" {
		inputUlxlyArgs.destAddress = fromAddress.String()
		log.Info().Stringer("destAddress", fromAddress).Msg("No destination address specified. Using private key's address")
	}

	if inputUlxlyArgs.callDataFile != "" {
		rawCallData, iErr := os.ReadFile(inputUlxlyArgs.callDataFile)
		if iErr != nil {
			return iErr
		}
		if inputUlxlyArgs.callData != "0x" {
			return fmt.Errorf("both %s and %s flags were provided", ArgCallData, ArgCallDataFile)
		}
		inputUlxlyArgs.callData = string(rawCallData)
	}

	bridgeService, err = bridge_service_factory.NewBridgeService(inputUlxlyArgs.bridgeServiceURL, inputUlxlyArgs.insecure, inputUlxlyArgs.legacy)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create bridge service")
		return err
	}

	bridgeServicesURLs, err := getBridgeServiceURLs()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get bridge service URLs")
		return err
	}

	for networkID, url := range bridgeServicesURLs {
		bs, err := bridge_service_factory.NewBridgeService(url, inputUlxlyArgs.insecure, inputUlxlyArgs.legacy)
		if err != nil {
			log.Error().Err(err).Str("url", url).Msg("Unable to create bridge service")
			return err
		}
		if _, exists := bridgeServices[networkID]; exists {
			log.Warn().Uint32("networkID", networkID).Str("url", url).Msg("Duplicate network ID found for bridge service URL. Overwriting previous entry.")
		}
		bridgeServices[networkID] = bs
		log.Info().Uint32("networkID", networkID).Str("url", url).Msg("Added bridge service")
	}

	return nil
}

func fatalIfError(err error) {
	if err == nil {
		return
	}
	log.Fatal().Err(err).Msg("Unexpected error occurred")
}

type FileOptions struct {
	FileName string
}

func (o *FileOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.FileName, ArgFileName, "", "", "ndjson file with events data")
}

type BalanceTreeOptions struct {
	L2ClaimsFile, L2DepositsFile, BridgeAddress, RpcURL string
	L2NetworkID                                         uint32
	Insecure                                            bool
}

func (o *BalanceTreeOptions) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&o.L2ClaimsFile, ArgL2ClaimsFileName, "", "", "ndjson file with l2 claim events data")
	f.StringVarP(&o.L2DepositsFile, ArgL2DepositsFileName, "", "", "ndjson file with l2 deposit events data")
	f.StringVarP(&o.BridgeAddress, ArgBridgeAddress, "", "", "bridge address")
	f.StringVarP(&o.RpcURL, ArgRPCURL, "r", "", "RPC URL")
	f.Uint32VarP(&o.L2NetworkID, ArgL2NetworkID, "", 0, "L2 network ID")
	f.BoolVarP(&o.Insecure, ArgInsecure, "", false, "skip TLS certificate verification")
}

type ProofOptions struct {
	DepositCount uint32
}

func (o *ProofOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32VarP(&o.DepositCount, ArgDepositCount, "", 0, "deposit number to generate a proof for")
}

type RollupsProofOptions struct {
	RollupID           uint32
	CompleteMerkleTree bool
}

func (o *RollupsProofOptions) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.Uint32VarP(&o.RollupID, ArgRollupID, "", 0, "rollup ID number to generate a proof for")
	f.BoolVarP(&o.CompleteMerkleTree, ArgCompleteMT, "", false, "get proof for a leave higher than the highest rollup ID")
}

type GetEvent struct {
	URL                            string
	FromBlock, ToBlock, FilterSize uint64
	Insecure                       bool
}

func (o *GetEvent) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&o.URL, ArgRPCURL, "u", "", "RPC URL to read the events data")
	f.Uint64VarP(&o.FromBlock, ArgFromBlock, "f", 0, "start of the range of blocks to retrieve")
	f.Uint64VarP(&o.ToBlock, ArgToBlock, "t", 0, "end of the range of blocks to retrieve")
	f.Uint64VarP(&o.FilterSize, ArgFilterSize, "i", 1000, "batch size for individual filter queries")
	f.BoolVarP(&o.Insecure, ArgInsecure, "", false, "skip TLS certificate verification")
	fatalIfError(cmd.MarkFlagRequired(ArgFromBlock))
	fatalIfError(cmd.MarkFlagRequired(ArgToBlock))
	fatalIfError(cmd.MarkFlagRequired(ArgRPCURL))
}

type GetSmcOptions struct {
	BridgeAddress string
}

func (o *GetSmcOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.BridgeAddress, ArgBridgeAddress, "a", "", "address of the ulxly bridge")
}

type GetVerifyBatchesOptions struct {
	RollupManagerAddress string
}

func (o *GetVerifyBatchesOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.RollupManagerAddress, ArgRollupManagerAddress, "a", "", "address of the rollup manager contract")
}

func init() {
	bridgeAssetCommand = &cobra.Command{
		Use:     "asset",
		Short:   "Move ETH or an ERC20 between to chains",
		Long:    bridgeAssetUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bridgeAsset(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
		},
		SilenceUsage: true,
	}
	bridgeMessageCommand = &cobra.Command{
		Use:     "message",
		Short:   "Send some ETH along with data from one chain to another chain",
		Long:    bridgeMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bridgeMessage(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
		},
		SilenceUsage: true,
	}
	bridgeMessageWETHCommand = &cobra.Command{
		Use:     "weth",
		Short:   "For L2's that use a gas token, use this to transfer WETH to another chain",
		Long:    bridgeWETHMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bridgeWETHMessage(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
		},
		SilenceUsage: true,
	}
	claimAssetCommand = &cobra.Command{
		Use:     "asset",
		Short:   "Claim a deposit",
		Long:    claimAssetUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := claimAsset(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
		},
		SilenceUsage: true,
	}
	claimMessageCommand = &cobra.Command{
		Use:     "message",
		Short:   "Claim a message",
		Long:    claimMessageUsage,
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := claimMessage(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
		},
		SilenceUsage: true,
	}
	claimEverythingCommand = &cobra.Command{
		Use:     "claim-everything",
		Short:   "Attempt to claim as many deposits and messages as possible",
		PreRunE: prepInputs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := claimEverything(cmd); err != nil {
				log.Fatal().Err(err).Msg("Received critical error")
			}
			return nil
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
	fileOptions.AddFlags(proofCommand)
	proofOptions.AddFlags(proofCommand)
	ulxlyProofsCmd.AddCommand(proofCommand)
	ULxLyCmd.AddCommand(proofCommand)

	rollupsProofCommand = &cobra.Command{
		Use:   "rollups-proof",
		Short: "Generate a proof for a given range of rollups",
		Long:  rollupsProofUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rollupsExitRootProof(args)
		},
		SilenceUsage: true,
	}
	fileOptions.AddFlags(rollupsProofCommand)
	rollupsProofOptions.AddFlags(rollupsProofCommand)
	ulxlyProofsCmd.AddCommand(rollupsProofCommand)
	ULxLyCmd.AddCommand(rollupsProofCommand)

	balanceTreeCommand = &cobra.Command{
		Use:   "compute-balance-tree",
		Short: "Compute the balance tree given the deposits",
		Long:  balanceTreeUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return balanceTree()
		},
		SilenceUsage: true,
	}
	balanceTreeOptions.AddFlags(balanceTreeCommand)
	ULxLyCmd.AddCommand(balanceTreeCommand)

	nullifierAndBalanceTreeCommand = &cobra.Command{
		Use:   "compute-balance-nullifier-tree",
		Short: "Compute the balance tree and the nullifier tree given the deposits and claims",
		Long:  nullifierAndBalanceTreeUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nullifierAndBalanceTree(args)
		},
		SilenceUsage: true,
	}
	balanceTreeOptions.AddFlags(nullifierAndBalanceTreeCommand)
	ULxLyCmd.AddCommand(nullifierAndBalanceTreeCommand)

	nullifierTreeCommand = &cobra.Command{
		Use:   "compute-nullifier-tree",
		Short: "Compute the nullifier tree given the claims",
		Long:  nullifierTreeUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nullifierTree(args)
		},
		SilenceUsage: true,
	}
	fileOptions.AddFlags(nullifierTreeCommand)
	ULxLyCmd.AddCommand(nullifierTreeCommand)

	getDepositCommand = &cobra.Command{
		Use:   "get-deposits",
		Short: "Generate ndjson for each bridge deposit over a particular range of blocks",
		Long:  depositGetUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return readDeposit(cmd)
		},
		SilenceUsage: true,
	}
	getEvent.AddFlags(getDepositCommand)
	getSmcOptions.AddFlags(getDepositCommand)
	ulxlyGetEventsCmd.AddCommand(getDepositCommand)
	ULxLyCmd.AddCommand(getDepositCommand)

	getClaimCommand = &cobra.Command{
		Use:   "get-claims",
		Short: "Generate ndjson for each bridge claim over a particular range of blocks",
		Long:  claimGetUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return readClaim(cmd)
		},
		SilenceUsage: true,
	}
	getEvent.AddFlags(getClaimCommand)
	getSmcOptions.AddFlags(getClaimCommand)
	ulxlyGetEventsCmd.AddCommand(getClaimCommand)
	ULxLyCmd.AddCommand(getClaimCommand)

	getVerifyBatchesCommand = &cobra.Command{
		Use:   "get-verify-batches",
		Short: "Generate ndjson for each verify batch over a particular range of blocks",
		Long:  verifyBatchesGetUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return readVerifyBatches(cmd)
		},
		SilenceUsage: true,
	}
	getEvent.AddFlags(getVerifyBatchesCommand)
	getVerifyBatchesOptions.AddFlags(getVerifyBatchesCommand)
	ulxlyGetEventsCmd.AddCommand(getVerifyBatchesCommand)
	ULxLyCmd.AddCommand(getVerifyBatchesCommand)

	// Arguments for both bridge and claim
	fBridgeAndClaim := ulxlyBridgeAndClaimCmd.PersistentFlags()
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.rpcURL, ArgRPCURL, "", "RPC URL to send the transaction")
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.bridgeAddress, ArgBridgeAddress, "", "address of the lxly bridge")
	fBridgeAndClaim.Uint64Var(&inputUlxlyArgs.gasLimit, ArgGasLimit, 0, "force specific gas limit for transaction")
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.chainID, ArgChainID, "", "chain ID to use in the transaction")
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.privateKey, ArgPrivateKey, "", "hex encoded private key for sending transaction")
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.destAddress, ArgDestAddress, "", "destination address for the bridge")
	fBridgeAndClaim.Uint64Var(&inputUlxlyArgs.timeout, ArgTimeout, 60, "timeout in seconds to wait for transaction receipt confirmation")
	fBridgeAndClaim.StringVar(&inputUlxlyArgs.gasPrice, ArgGasPrice, "", "gas price to use")
	fBridgeAndClaim.BoolVar(&inputUlxlyArgs.dryRun, ArgDryRun, false, "do all of the transaction steps but do not send the transaction")
	fBridgeAndClaim.BoolVar(&inputUlxlyArgs.insecure, ArgInsecure, false, "skip TLS certificate verification")
	fatalIfError(ulxlyBridgeAndClaimCmd.MarkPersistentFlagRequired(ArgBridgeAddress))

	// bridge specific args
	fBridge := ulxlyBridgeCmd.PersistentFlags()
	fBridge.BoolVar(&inputUlxlyArgs.forceUpdate, ArgForceUpdate, true, "update the new global exit root")
	fBridge.StringVar(&inputUlxlyArgs.value, ArgValue, "0", "amount in wei to send with the transaction")
	fBridge.Uint32Var(&inputUlxlyArgs.destNetwork, ArgDestNetwork, 0, "rollup ID of the destination network")
	fBridge.StringVar(&inputUlxlyArgs.tokenAddress, ArgTokenAddress, "0x0000000000000000000000000000000000000000", "address of ERC20 token to use")
	fBridge.StringVar(&inputUlxlyArgs.callData, ArgCallData, "0x", "call data to be passed directly with bridge-message or as an ERC20 Permit")
	fBridge.StringVar(&inputUlxlyArgs.callDataFile, ArgCallDataFile, "", "a file containing hex encoded call data")
	fatalIfError(ulxlyBridgeCmd.MarkPersistentFlagRequired(ArgDestNetwork))

	// Claim specific args
	fClaim := ulxlyClaimCmd.PersistentFlags()
	fClaim.Uint32Var(&inputUlxlyArgs.depositCount, ArgDepositCount, 0, "deposit count of the bridge transaction")
	fClaim.Uint32Var(&inputUlxlyArgs.depositNetwork, ArgDepositNetwork, 0, "rollup ID of the network where the deposit was made")
	fClaim.StringVar(&inputUlxlyArgs.bridgeServiceURL, ArgBridgeServiceURL, "", "URL of the bridge service")
	fClaim.StringVar(&inputUlxlyArgs.globalIndex, ArgGlobalIndex, "", "an override of the global index value")
	fClaim.DurationVar(&inputUlxlyArgs.wait, ArgWait, time.Duration(0), "retry claiming until deposit is ready, up to specified duration (available for claim asset and claim message)")
	fClaim.StringVar(&inputUlxlyArgs.proofGER, ArgProofGER, "", "if specified and using legacy mode, the proof will be generated against this GER")
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgDepositCount))
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgDepositNetwork))
	fatalIfError(ulxlyClaimCmd.MarkPersistentFlagRequired(ArgBridgeServiceURL))

	// Claim Everything Helper Command
	fClaimEverything := claimEverythingCommand.Flags()
	fClaimEverything.StringSliceVar(&inputUlxlyArgs.bridgeServiceURLs, ArgBridgeMappings, nil, "network ID to bridge service URL mappings (e.g. '1=http://network-1-bridgeurl,7=http://network-2-bridgeurl')")
	fClaimEverything.IntVar(&inputUlxlyArgs.bridgeLimit, ArgBridgeLimit, 25, "limit the number or responses returned by the bridge service when claiming")
	fClaimEverything.IntVar(&inputUlxlyArgs.bridgeOffset, ArgBridgeOffset, 0, "offset to specify for pagination of underlying bridge service deposits")
	fClaimEverything.UintVar(&inputUlxlyArgs.concurrency, ArgConcurrency, 1, "worker pool size for claims")

	fatalIfError(claimEverythingCommand.MarkFlagRequired(ArgBridgeMappings))

	// Top Level
	ULxLyCmd.AddCommand(ulxlyBridgeAndClaimCmd)
	ULxLyCmd.AddCommand(ulxlyGetEventsCmd)
	ULxLyCmd.AddCommand(ulxlyProofsCmd)
	ULxLyCmd.AddCommand(emptyProofCommand)
	ULxLyCmd.AddCommand(zeroProofCommand)
	ULxLyCmd.AddCommand(proofCommand)

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
