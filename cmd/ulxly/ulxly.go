package ulxly

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

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
	Deposit struct {
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
	} `json:"deposit"`
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

func readDeposit(cmd *cobra.Command) error {
	rpcUrl, err := cmd.Flags().GetString(ArgRPCURL)
	if err != nil {
		return err
	}

	bridgeAddress := *inputUlxlyArgs.bridgeAddress
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
	callData := common.Hex2Bytes(callDataString)

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	}

	bridgeTxn, err := bridgeV2.BridgeAsset(auth, destinationNetwork, toAddress, value, tokenAddress, isForced, callData)
	if err != nil {
		log.Error().Err(err).Msg("Unable to interact with bridge contract")
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
	callData := common.Hex2Bytes(callDataString)

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	}

	bridgeTxn, err := bridgeV2.BridgeMessage(auth, destinationNetwork, toAddress, isForced, callData)
	if err != nil {
		log.Error().Err(err).Msg("Unable to interact with bridge contract")
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
	callData := common.Hex2Bytes(callDataString)

	bridgeTxn, err := bridgeV2.BridgeMessageWETH(auth, destinationNetwork, toAddress, value, isForced, callData)
	if err != nil {
		log.Error().Err(err).Msg("Unable to interact with bridge contract")
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

	///////////////////////////////////////// TODO BORRAR
	// gasPrice, err := client.SuggestGasPrice(ctx) // This call is done automatically if it is not set
	// if err != nil {
	// 	log.Error().Err(err).Msg("Cannot get suggested gas price")
	// }
	// auth.GasPrice = big.NewInt(0).Mul(gasPrice,big.NewInt(10))
	/////////////////////////////////////////////////////////

	// Call the bridge service RPC URL to get the merkle proofs and exit roots and parses them to the correct formats.
	bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%d&net_id=%d", bridgeServiceUrl, depositCount, depositNetwork)
	merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)

	// Call the bridge service RPC URL to get the deposits data and parses them to the correct formats.
	bridgeServiceDepositsEndpoint := fmt.Sprintf("%s/bridge?net_id=%d&deposit_cnt=%d", bridgeServiceUrl, depositNetwork, depositCount)
	globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err := getDeposit(bridgeServiceDepositsEndpoint)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	if leafType != 0 {
		log.Warn().Msg("Deposit leafType is not asset")
	}

	claimTxn, err := bridgeV2.ClaimAsset(auth, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), claimOriginalNetwork, originAddress, claimDestNetwork, toAddress, amount, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Unable to interact with bridge contract")
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

	// Call the bridge service RPC URL to get the merkle proofs and exit roots and parses them to the correct formats.
	bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%d&net_id=%d", bridgeServiceUrl, depositCount, depositNetwork)
	merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)

	// Call the bridge service RPC URL to get the deposits data and parses them to the correct formats.
	bridgeServiceDepositsEndpoint := fmt.Sprintf("%s/bridge?net_id=%d&deposit_cnt=%d", bridgeServiceUrl, depositNetwork, depositCount)
	globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, err := getDeposit(bridgeServiceDepositsEndpoint)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	if leafType != 1 {
		log.Warn().Msg("Deposit leafType is not message")
	}

	claimTxn, err := bridgeV2.ClaimMessage(auth, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), claimOriginalNetwork, originAddress, claimDestNetwork, toAddress, amount, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Unable to interact with bridge contract")
		return err
	}
	log.Info().Msg("claimTxn: " + claimTxn.Hash().String())
	return WaitMineTransaction(cmd.Context(), client, claimTxn, timeoutTxnReceipt)
}

// Wait for the transaction to be mined
func WaitMineTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, txTimeout uint64) error {
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
	}

	// value := big.NewInt(*ulxlyInputArgs.Amount)
	gasLimit := ulxlyInputArgGasLimit
	// gasPrice, err := client.SuggestGasPrice(ctx) // This call is done automatically if it is not set
	// if err != nil {
	// 	log.Error().Err(err).Msg("Cannot get suggested gas price")
	// }
	// gasTipCap, err := client.SuggestGasTipCap(ctx)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Cannot get suggested gas tip cap")
	// }

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
		log.Error().Msg("The Merkle Proofs cannot be retrieved, double check the input arguments and try again.")
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
	var bridgeDeposit BridgeDeposit
	err = json.Unmarshal(bodyBridgeDeposit, &bridgeDeposit) // Parse []byte to go struct pointer, and shadow err variable
	if err != nil {
		log.Error().Err(err).Msg("Can not unmarshal JSON")
		return
	}

	globalIndex = new(big.Int)
	amount = new(big.Int)

	defer reqBridgeDeposit.Body.Close()
	if bridgeDeposit.Code != nil {
		return globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, fmt.Errorf("error code received getting the deposit. Code: %d, Message: %s", *bridgeDeposit.Code, *bridgeDeposit.Message)
	}

	if !bridgeDeposit.Deposit.ReadyForClaim {
		log.Error().Msg("The claim transaction is not yet ready to be claimed. Try again in a few blocks.")
		return nil, common.HexToAddress("0x0"), nil, nil, 0, 0, 0, errors.New("the claim transaction is not yet ready to be claimed, try again in a few blocks")
	} else if bridgeDeposit.Deposit.ClaimTxHash != "" {
		log.Info().Str("claimTxHash", bridgeDeposit.Deposit.ClaimTxHash).Msg("The claim transaction has already been claimed")
		return nil, common.HexToAddress("0x0"), nil, nil, 0, 0, 0, errors.New("the claim transaction has already been claimed")
	}
	originAddress = common.HexToAddress(bridgeDeposit.Deposit.OrigAddr)
	globalIndex.SetString(bridgeDeposit.Deposit.GlobalIndex, 10)
	amount.SetString(bridgeDeposit.Deposit.Amount, 10)
	metadata = common.Hex2Bytes(bridgeDeposit.Deposit.Metadata)
	leafType = bridgeDeposit.Deposit.LeafType
	claimDestNetwork = bridgeDeposit.Deposit.DestNet
	claimOriginalNetwork = bridgeDeposit.Deposit.OrigNet
	return globalIndex, originAddress, amount, metadata, leafType, claimDestNetwork, claimOriginalNetwork, nil
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
	Short: "Utilities for interacting with the lxly bridge",
	Long:  "These are low level tools for directly scanning bridge events and constructing proofs.",
	Args:  cobra.NoArgs,
}
var ulxlyBridgeAndClaimCmd = &cobra.Command{
	Args:   cobra.NoArgs,
	Hidden: true,
}

var ulxlxBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "commands for making deposits to the uLxLy bridge",
	Args:  cobra.NoArgs,
}

var ulxlyClaimCmd = &cobra.Command{
	Use:   "claim",
	Short: "commands for making claims of deposits from the uLxLy bridge",
	Args:  cobra.NoArgs,
}

type ulxlyArgs struct {
	gasLimit         *uint64
	chainID          *string
	privateKey       *string
	value            *string
	rpcURL           *string
	bridgeAddress    *string
	destNetwork      *uint32
	destAddress      *string
	tokenAddress     *string
	forceUpdate      *bool
	callData         *string
	timeout          *uint64
	depositCount     *uint64
	depositNetwork   *uint64
	bridgeServiceURL *string
	inputFileName    *string
	fromBlock        *uint64
	toBlock          *uint64
	filterSize       *uint64
	depositNumber    *uint64
}

var inputUlxlyArgs = ulxlyArgs{}

var (
	bridgeAssetCommand       *cobra.Command
	bridgeMessageCommand     *cobra.Command
	bridgeMessageWETHCommand *cobra.Command
	claimAssetCommand        *cobra.Command
	claimMessageCommand      *cobra.Command
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
	ArgTimeout          = "transaction-receipt-timeout"
	ArgDepositCount     = "deposit-count"
	ArgDepositNetwork   = "deposit-network"
	ArgBridgeServiceURL = "bridge-service-url"
	ArgFileName         = "file-name"
	ArgFromBlock        = "from-block"
	ArgToBlock          = "to-block"
	ArgFilterSize       = "filter-size"
	ArgTokenAddress     = "token-address"
)

func init() {
	bridgeAssetCommand = &cobra.Command{
		Use:   "asset",
		Short: "send a single deposit of value or an ERC20 into the bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeAsset(cmd)
		},
	}
	bridgeMessageCommand = &cobra.Command{
		Use:   "message",
		Short: "send some value along with call data into the bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeMessage(cmd)
		},
	}
	bridgeMessageWETHCommand = &cobra.Command{
		Use:   "weth",
		Short: "send some WETH into the bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeWETHMessage(cmd)
		},
	}
	claimAssetCommand = &cobra.Command{
		Use:   "asset",
		Short: "perform a claim of a given deposit in the bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return claimAsset(cmd)
		},
	}
	claimMessageCommand = &cobra.Command{
		Use:   "message",
		Short: "perform a claim of a given message in the birdge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	emptyProofCommand = &cobra.Command{
		Use:   "empty-proof",
		Short: "create an empty proof",
		Long:  "Use this command to print an empty proof response that's filled with zero-valued siblings like 0x0000000000000000000000000000000000000000000000000000000000000000. This can be useful when you need to submit a dummy proof.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return emptyProof()
		},
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
	}
	proofCommand = &cobra.Command{
		Use:   "proof",
		Short: "generate a proof for a given range of deposits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return proof(args)
		},
	}
	getDepositCommand = &cobra.Command{
		Use:   "get-deposits",
		Short: "generate ndjson for each bridge deposit over a particular range of blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return readDeposit(cmd)
		},
	}

	// Arguments for both bridge and claim
	inputUlxlyArgs.gasLimit = ulxlyBridgeAndClaimCmd.PersistentFlags().Uint64(ArgGasLimit, 0, "force a gas limit when sending a transaction")
	inputUlxlyArgs.chainID = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgChainID, "", "set the chain id to be used in the transaction")
	inputUlxlyArgs.privateKey = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgPrivateKey, "", "the hex encoded private key to be used when sending the tx")
	inputUlxlyArgs.rpcURL = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgRPCURL, "", "the URL of the RPC to send the transaction")
	inputUlxlyArgs.bridgeAddress = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgBridgeAddress, "", "the address of the lxly bridge")
	inputUlxlyArgs.destAddress = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgDestAddress, "", "the address where the bridge will be sent to")
	inputUlxlyArgs.callData = ulxlyBridgeAndClaimCmd.PersistentFlags().String(ArgCallData, "0x", "call data to be passed directly with bridge-message or as an ERC20 Permit")
	inputUlxlyArgs.timeout = ulxlyBridgeAndClaimCmd.PersistentFlags().Uint64(ArgTimeout, 60, "the amount of time to wait while trying to confirm a transaction receipt")
	ulxlyBridgeAndClaimCmd.MarkFlagRequired(ArgPrivateKey)
	ulxlyBridgeAndClaimCmd.MarkFlagRequired(ArgRPCURL)
	ulxlyBridgeAndClaimCmd.MarkFlagRequired(ArgBridgeAddress)
	ulxlyBridgeAndClaimCmd.MarkFlagRequired(ArgDestAddress)

	// bridge specific args
	inputUlxlyArgs.forceUpdate = ulxlxBridgeCmd.PersistentFlags().Bool(ArgForceUpdate, true, "indicates if the new global exit root is updated or not")
	inputUlxlyArgs.value = ulxlxBridgeCmd.PersistentFlags().String(ArgValue, "", "the amount in wei to be sent along with the transaction")
	inputUlxlyArgs.destNetwork = ulxlxBridgeCmd.PersistentFlags().Uint32(ArgDestNetwork, 0, "the rollup id of the destination network")
	inputUlxlyArgs.tokenAddress = ulxlxBridgeCmd.PersistentFlags().String(ArgTokenAddress, "0x0000000000000000000000000000000000000000", "the address of an ERC20 token to be used")
	ulxlxBridgeCmd.MarkFlagRequired(ArgDestNetwork)

	// Claim specific args
	inputUlxlyArgs.depositCount = ulxlyClaimCmd.PersistentFlags().Uint64(ArgDepositCount, 0, "the deposit count of the bridge transaction")
	inputUlxlyArgs.depositNetwork = ulxlyClaimCmd.PersistentFlags().Uint64(ArgDepositNetwork, 0, "the rollup id of the network where the bridge is being claimed")
	inputUlxlyArgs.bridgeServiceURL = ulxlyClaimCmd.PersistentFlags().String(ArgBridgeServiceURL, "", "the URL of the bridge service")
	ulxlyClaimCmd.MarkFlagRequired(ArgDepositCount)
	ulxlyClaimCmd.MarkFlagRequired(ArgDepositNetwork)
	ulxlyClaimCmd.MarkFlagRequired(ArgBridgeServiceURL)

	// Args that are just for the get deposit command
	inputUlxlyArgs.inputFileName = getDepositCommand.Flags().String(ArgFileName, "", "An ndjson file with deposit data")
	inputUlxlyArgs.fromBlock = getDepositCommand.Flags().Uint64(ArgFromBlock, 0, "The start of the range of blocks to retrieve")
	inputUlxlyArgs.toBlock = getDepositCommand.Flags().Uint64(ArgToBlock, 0, "The end of the range of blocks to retrieve")
	inputUlxlyArgs.filterSize = getDepositCommand.Flags().Uint64(ArgFilterSize, 1000, "The batch size for individual filter queries")
	getDepositCommand.Flags().String(ArgRPCURL, "", "The RPC URL to read deposit data")
	getDepositCommand.MarkFlagRequired(ArgFromBlock)
	getDepositCommand.MarkFlagRequired(ArgToBlock)
	getDepositCommand.MarkFlagRequired(ArgRPCURL)

	// Top Level
	ULxLyCmd.AddCommand(ulxlyBridgeAndClaimCmd)
	ULxLyCmd.AddCommand(emptyProofCommand)
	ULxLyCmd.AddCommand(zeroProofCommand)
	ULxLyCmd.AddCommand(proofCommand)
	ULxLyCmd.AddCommand(getDepositCommand)

	ULxLyCmd.AddCommand(ulxlxBridgeCmd)
	ULxLyCmd.AddCommand(ulxlyClaimCmd)

	// Bridge and Claim
	ulxlyBridgeAndClaimCmd.AddCommand(ulxlxBridgeCmd)
	ulxlyBridgeAndClaimCmd.AddCommand(ulxlyClaimCmd)

	// Bridge
	ulxlxBridgeCmd.AddCommand(bridgeAssetCommand)
	ulxlxBridgeCmd.AddCommand(bridgeMessageCommand)
	ulxlxBridgeCmd.AddCommand(bridgeMessageWETHCommand)

	// Claim
	ulxlyClaimCmd.AddCommand(claimAssetCommand)
	ulxlyClaimCmd.AddCommand(claimMessageCommand)

}
