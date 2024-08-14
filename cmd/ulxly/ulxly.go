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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/maticnetwork/polygon-cli/bindings/ulxly"
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

type uLxLyArgs struct {
	FromBlock     *uint64
	ToBlock       *uint64
	RPCURL        *string
	BridgeAddress *string
	FilterSize    *uint64

	ClaimIndex              *string
	ClaimAddress            *string
	ClaimOriginNetwork      *string
	ClaimDestinationNetwork *string
	BridgeServiceRPCURL     *string
	ClaimRPCURL             *string
	ClaimPrivateKey         *string
	ClaimBridgeAddress      *string
	ClaimGasLimit           *uint64
	ClaimChainID            *string
	ClaimTimeoutTxnReceipt  *uint32
	ClaimMessage            *bool
	ClaimWETH               *bool

	InputFileName            *string
	DepositNum               *uint32
	DepositPrivateKey        *string
	DepositGasLimit          *uint64
	Amount                   *int64 // HACK: This should be big.Int but depositNewCmd.PersistentFlags() doesn't support that type.
	DestinationNetwork       *uint32
	DestinationAddress       *string
	DepositRPCURL            *string
	DepositBridgeAddress     *string
	DepositChainID           *string
	TokenAddress             *string
	IsForced                 *bool
	CallData                 *string
	DepositTimeoutTxnReceipt *uint32
	DepositMessage           *bool
	DepositWETH              *bool
}

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

type BridgeDeposits struct {
	Deposit []struct {
		LeafType      int    `json:"leaf_type"`
		OrigNet       int    `json:"orig_net"`
		OrigAddr      string `json:"orig_addr"`
		Amount        string `json:"amount"`
		DestNet       int    `json:"dest_net"`
		DestAddr      string `json:"dest_addr"`
		BlockNum      string `json:"block_num"`
		DepositCnt    string `json:"deposit_cnt"`
		NetworkID     int    `json:"network_id"`
		TxHash        string `json:"tx_hash"`
		ClaimTxHash   string `json:"claim_tx_hash"`
		Metadata      string `json:"metadata"`
		ReadyForClaim bool   `json:"ready_for_claim"`
		GlobalIndex   string `json:"global_index"`
	} `json:"deposits"`
	TotalCnt string `json:"total_cnt"`
}

var ulxlyInputArgs uLxLyArgs

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the lxly bridge",
	Long:  "These are low level tools for directly scanning bridge events and constructing proofs.",
	Args:  cobra.NoArgs,
}

//go:embed depositGetUsage.md
var depositGetUsage string
var depositGetCmd = &cobra.Command{
	Use:     "deposit-get",
	Short:   "Get a range of deposits",
	Long:    depositGetUsage,
	Args:    cobra.NoArgs,
	PreRunE: checkGetDepositArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Dial the Ethereum RPC server.
		rpc, err := ethrpc.DialContext(ctx, *ulxlyInputArgs.RPCURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer rpc.Close()
		ec := ethclient.NewClient(rpc)

		bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(*ulxlyInputArgs.BridgeAddress), ec)
		if err != nil {
			return err
		}
		fromBlock := *ulxlyInputArgs.FromBlock
		toBlock := *ulxlyInputArgs.ToBlock
		currentBlock := fromBlock
		for currentBlock < toBlock {
			endBlock := currentBlock + *ulxlyInputArgs.FilterSize
			if endBlock > toBlock {
				endBlock = toBlock
			}

			opts := bind.FilterOpts{
				Start:   currentBlock,
				End:     &endBlock,
				Context: ctx,
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
	},
}

//go:embed depositNewUsage.md
var depositNewUsage string
var depositNewCmd = &cobra.Command{
	Use:     "deposit-new",
	Short:   "Make a uLxLy deposit transaction",
	Long:    depositNewUsage,
	Args:    cobra.NoArgs,
	PreRunE: checkDepositArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Dial the Ethereum RPC server.
		client, err := ethclient.DialContext(ctx, *ulxlyInputArgs.DepositRPCURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer client.Close()
		// Initialize and assign variables required to send transaction payload
		bridgeV2, privateKey, fromAddress, gasLimit, gasPrice, toAddress, signer := generateTransactionPayload(ctx, client, *ulxlyInputArgs.DepositBridgeAddress, *ulxlyInputArgs.DepositPrivateKey, *ulxlyInputArgs.DepositGasLimit, *ulxlyInputArgs.DestinationAddress, *ulxlyInputArgs.DepositChainID)

		value := big.NewInt(*ulxlyInputArgs.Amount)
		tokenAddress := common.HexToAddress(*ulxlyInputArgs.TokenAddress)
		callData := common.Hex2Bytes(*ulxlyInputArgs.CallData)

		tops := &bind.TransactOpts{
			Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
				return types.SignTx(transaction, signer, privateKey)
			},
			From:      fromAddress,
			Context:   ctx,
			GasLimit:  gasLimit,
			GasPrice:  gasPrice,
			GasFeeCap: nil,
			GasTipCap: nil,
		}
		if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
			tops = &bind.TransactOpts{
				Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
					return types.SignTx(transaction, signer, privateKey)
				},
				Value:     value,
				From:      fromAddress,
				Context:   ctx,
				GasLimit:  gasLimit,
				GasPrice:  gasPrice,
				GasFeeCap: nil,
				GasTipCap: nil,
			}
		}

		var bridgeTxn *types.Transaction
		switch {
		case *ulxlyInputArgs.DepositMessage:
			bridgeTxn, err = bridgeV2.BridgeMessage(tops, *ulxlyInputArgs.DestinationNetwork, toAddress, *ulxlyInputArgs.IsForced, callData)
			if err != nil {
				log.Error().Err(err).Msg("Unable to interact with bridge contract")
				return err
			}
		case *ulxlyInputArgs.DepositWETH:
			bridgeTxn, err = bridgeV2.BridgeMessageWETH(tops, *ulxlyInputArgs.DestinationNetwork, toAddress, value, *ulxlyInputArgs.IsForced, callData)
			if err != nil {
				log.Error().Err(err).Msg("Unable to interact with bridge contract")
				return err
			}
		default:
			bridgeTxn, err = bridgeV2.BridgeAsset(tops, *ulxlyInputArgs.DestinationNetwork, toAddress, value, tokenAddress, *ulxlyInputArgs.IsForced, callData)
			if err != nil {
				log.Error().Err(err).Msg("Unable to interact with bridge contract")
				return err
			}
		}

		// Wait for the transaction to be mined
		// TODO: Consider creating a function for this section
		txnMinedTimer := time.NewTimer(time.Duration(*ulxlyInputArgs.DepositTimeoutTxnReceipt) * time.Second)
		defer txnMinedTimer.Stop()
		for {
			select {
			case <-txnMinedTimer.C:
				fmt.Printf("Wait timer for transaction receipt exceeded!")
				return nil
			default:
				r, err := client.TransactionReceipt(ctx, bridgeTxn.Hash())
				if err != nil {
					if err.Error() != "not found" {
						log.Error().Err(err)
						return err
					}
					time.Sleep(1 * time.Second)
					continue
				}
				if r.Status != 0 {
					fmt.Printf("Deposit Transaction Successful: %s\n", r.TxHash)
					return nil
				} else if r.Status == 0 {
					fmt.Printf("Deposit Transaction Failed: %s\n", r.TxHash)
					fmt.Printf("Perhaps try increasing the gas limit:\n")
					fmt.Printf("Current gas limit: %d\n", gasLimit)
					fmt.Printf("Cumulative gas used for transaction: %d\n", r.CumulativeGasUsed)
					return nil
				}
				time.Sleep(1 * time.Second)
			}
		}
	},
}

//go:embed depositClaimUsage.md
var depositClaimUsage string
var depositClaimCmd = &cobra.Command{
	Use:     "deposit-claim",
	Short:   "Make a uLxLy claim transaction",
	Long:    depositClaimUsage,
	Args:    cobra.NoArgs,
	PreRunE: checkClaimArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Dial Ethereum client
		client, err := ethclient.DialContext(ctx, *ulxlyInputArgs.ClaimRPCURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer client.Close()
		// Initialize and assign variables required to send transaction payload
		bridgeV2, privateKey, fromAddress, gasLimit, gasPrice, toAddress, signer := generateTransactionPayload(ctx, client, *ulxlyInputArgs.ClaimBridgeAddress, *ulxlyInputArgs.ClaimPrivateKey, *ulxlyInputArgs.ClaimGasLimit, *ulxlyInputArgs.ClaimAddress, *ulxlyInputArgs.ClaimChainID)

		// Call the bridge service RPC URL to get the merkle proofs and exit roots and parses them to the correct formats.
		bridgeServiceProofEndpoint := fmt.Sprintf("%s/merkle-proof?deposit_cnt=%s&net_id=%s", *ulxlyInputArgs.BridgeServiceRPCURL, *ulxlyInputArgs.ClaimIndex, *ulxlyInputArgs.ClaimOriginNetwork)
		merkleProofArray, rollupMerkleProofArray, mainExitRoot, rollupExitRoot := getMerkleProofsExitRoots(bridgeServiceProofEndpoint)

		tops := &bind.TransactOpts{
			Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
				return types.SignTx(transaction, signer, privateKey)
			},
			// Value:     value,
			From:      fromAddress,
			Context:   ctx,
			GasLimit:  gasLimit,
			GasPrice:  gasPrice,
			GasFeeCap: nil,
			GasTipCap: nil,
		}

		// Call the bridge service RPC URL to get the deposits data and parses them to the correct formats.
		bridgeServiceDepositsEndpoint := fmt.Sprintf("%s/bridges/%s", *ulxlyInputArgs.BridgeServiceRPCURL, *ulxlyInputArgs.ClaimAddress)
		globalIndex, originAddress, amount, metadata, err := getDeposits(bridgeServiceDepositsEndpoint)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		claimOriginNetwork, _ := strconv.Atoi(*ulxlyInputArgs.ClaimOriginNetwork)           // Convert ClaimOriginNetwork to int
		claimDestinationNetwork, _ := strconv.Atoi(*ulxlyInputArgs.ClaimDestinationNetwork) // Convert ClaimDestinationNetwork to int
		var claimTxn *types.Transaction
		switch {
		case *ulxlyInputArgs.ClaimMessage:
			claimTxn, err = bridgeV2.ClaimMessage(tops, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), uint32(claimOriginNetwork), originAddress, uint32(claimDestinationNetwork), toAddress, amount, metadata)
			if err != nil {
				log.Error().Err(err).Msg("Unable to interact with bridge contract")
				return err
			}
		default:
			claimTxn, err = bridgeV2.ClaimAsset(tops, merkleProofArray, rollupMerkleProofArray, globalIndex, [32]byte(mainExitRoot), [32]byte(rollupExitRoot), uint32(claimOriginNetwork), originAddress, uint32(claimDestinationNetwork), toAddress, amount, metadata)
			if err != nil {
				log.Error().Err(err).Msg("Unable to interact with bridge contract")
				return err
			}
		}

		// Wait for the transaction to be mined
		// TODO: Consider creating a function for this section
		txnMinedTimer := time.NewTimer(time.Duration(*ulxlyInputArgs.ClaimTimeoutTxnReceipt) * time.Second)
		defer txnMinedTimer.Stop()
		for {
			select {
			case <-txnMinedTimer.C:
				fmt.Printf("Wait timer for transaction receipt exceeded!")
				return nil
			default:
				r, err := client.TransactionReceipt(ctx, claimTxn.Hash())
				if err != nil {
					if err.Error() != "not found" {
						log.Error().Err(err)
						return err
					}
					time.Sleep(1 * time.Second)
					continue
				}
				if r.Status != 0 {
					fmt.Printf("Claim Transaction Successful: %s\n", r.TxHash)
					return nil
				} else if r.Status == 0 {
					fmt.Printf("Claim Transaction Failed: %s\n", r.TxHash)
					return nil
				}
				time.Sleep(1 * time.Second)
			}
		}
	},
}

//go:embed proofUsage.md
var proofUsage string
var ProofCmd = &cobra.Command{
	Use:     "proof",
	Short:   "generate a merkle proof",
	Long:    proofUsage,
	Args:    cobra.NoArgs,
	PreRunE: checkProofArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rawDepositData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		return readDeposits(rawDepositData)
	},
}

var EmptyProofCmd = &cobra.Command{
	Use:   "empty-proof",
	Short: "print an empty proof structure",
	Long: `Use this command to print an empty proof response that's filled with
zero-valued siblings like
0x0000000000000000000000000000000000000000000000000000000000000000. This
can be useful when you need to submit a dummy proof.`,
	Args:    cobra.NoArgs,
	PreRunE: checkProofArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := new(Proof)

		e := generateEmptyHashes(TreeDepth)
		copy(p.Siblings[:], e)
		fmt.Println(p.String())
		return nil
	},
}

var ZeroProofCmd = &cobra.Command{
	Use:   "zero-proof",
	Short: "print a proof structure with the zero hashes",
	Long: `Use this command to print a proof response that's filled with the zero
hashes. This values are very helpful for debugging because it would
tell you how populated the tree is and roughly which leaves and
siblings are empty. It's also helpful for sanity checking a proof
response to understand if the hashed value is part of the zero hashes
or if it's actually an intermediate hash.`,
	Args:    cobra.NoArgs,
	PreRunE: checkProofArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := new(Proof)

		e := generateZeroHashes(TreeDepth)
		copy(p.Siblings[:], e)
		fmt.Println(p.String())
		return nil
	},
}

func checkProofArgs(cmd *cobra.Command, args []string) error {
	return nil
}
func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if ulxlyInputArgs.InputFileName != nil && *ulxlyInputArgs.InputFileName != "" {
		return os.ReadFile(*ulxlyInputArgs.InputFileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
func readDeposits(rawDeposits []byte) error {
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
		if evt.DepositCount >= *ulxlyInputArgs.DepositNum {
			break
		}
	}

	p := imt.GetProof(*ulxlyInputArgs.DepositNum)

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
	sibling := common.Hash{}
	size := depositNum + 1
	currentZeroHashHeight := common.Hash{}

	siblings := [TreeDepth]common.Hash{}
	for height := 0; height < TreeDepth; height++ {
		siblingDepositNum := getSiblingDepositNumber(depositNum, uint32(height))

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

func generateTransactionPayload(ctx context.Context, client *ethclient.Client, ulxlyInputArgBridge string, ulxlyInputArgPvtKey string, ulxlyInputArgGasLimit uint64, ulxlyInputArgDestAddr string, ulxlyInputArgChainID string) (bridgeV2 *ulxly.Ulxly, privateKey *ecdsa.PrivateKey, fromAddress common.Address, gasLimit uint64, gasPrice *big.Int, toAddress common.Address, signer types.Signer) {
	var err error
	bridgeV2, err = ulxly.NewUlxly(common.HexToAddress(ulxlyInputArgBridge), client)
	if err != nil {
		return
	}

	privateKey, err = crypto.HexToECDSA(ulxlyInputArgPvtKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve private key")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error().Msg("Error casting public key to ECDSA")
	}

	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

	// value := big.NewInt(*ulxlyInputArgs.Amount)
	gasLimit = ulxlyInputArgGasLimit
	gasPrice, err = client.SuggestGasPrice(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get suggested gas price")
	}
	// gasTipCap, err := client.SuggestGasTipCap(ctx)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Cannot get suggested gas tip cap")
	// }

	toAddress = common.HexToAddress(ulxlyInputArgDestAddr)

	chainID := new(big.Int)
	// For manual input of chainID, use the user's input
	if ulxlyInputArgChainID != "" {
		chainID.SetString(ulxlyInputArgChainID, 10)
	} else { // If there is no user input for chainID, infer it from context
		chainID, err = client.NetworkID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Cannot get chain ID")
			return
		}
	}

	signer = types.LatestSignerForChainID(chainID)

	return bridgeV2, privateKey, fromAddress, gasLimit, gasPrice, toAddress, signer
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

func getDeposits(bridgeServiceDepositsEndpoint string) (globalIndex *big.Int, originAddress common.Address, amount *big.Int, metadata []byte, err error) {
	reqBridgeDeposits, err := http.Get(bridgeServiceDepositsEndpoint)
	if err != nil {
		log.Error().Err(err)
		return
	}
	bodyBridgeDeposit, err := io.ReadAll(reqBridgeDeposits.Body) // Response body is []byte
	if err != nil {
		log.Error().Err(err)
		return
	}
	var bridgeDeposit BridgeDeposits
	err = json.Unmarshal(bodyBridgeDeposit, &bridgeDeposit) // Parse []byte to go struct pointer, and shadow err variable
	if err != nil {
		log.Error().Err(err).Msg("Can not unmarshal JSON")
		return
	}

	globalIndex = new(big.Int)
	amount = new(big.Int)

	intClaimIndex, _ := strconv.Atoi(*ulxlyInputArgs.ClaimIndex) // Convert deposit_cnt to int
	for index, deposit := range bridgeDeposit.Deposit {
		intDepositCnt, _ := strconv.Atoi(deposit.DepositCnt) // Convert deposit_cnt to int
		if intDepositCnt == intClaimIndex {                  // deposit_cnt must match the user's input value
			if !bridgeDeposit.Deposit[index].ReadyForClaim {
				log.Error().Msg("The claim transaction is not yet ready to be claimed. Try again in a few blocks.")
				return nil, common.HexToAddress("0x0"), nil, nil, errors.New("The claim transaction is not yet ready to be claimed. Try again in a few blocks.")
			} else if bridgeDeposit.Deposit[index].ClaimTxHash != "" {
				fmt.Printf("The claim transaction has already been claimed at %s.", bridgeDeposit.Deposit[index].ClaimTxHash)
				return nil, common.HexToAddress("0x0"), nil, nil, errors.New("The claim transaction has already been claimed.")
			}
			originAddress = common.HexToAddress(bridgeDeposit.Deposit[index].OrigAddr)
			globalIndex.SetString(bridgeDeposit.Deposit[index].GlobalIndex, 10)
			amount.SetString(bridgeDeposit.Deposit[index].Amount, 10)
			metadata = common.Hex2Bytes(bridgeDeposit.Deposit[index].Metadata)
			return globalIndex, originAddress, amount, metadata, nil
		}
	}
	defer reqBridgeDeposits.Body.Close()

	return nil, common.HexToAddress("0x0"), nil, nil, errors.New("Failed to correctly get deposits...")
}

func checkGetDepositArgs(cmd *cobra.Command, args []string) error {
	if *ulxlyInputArgs.BridgeAddress == "" {
		return fmt.Errorf("please provide the bridge address")
	}
	if *ulxlyInputArgs.FromBlock > *ulxlyInputArgs.ToBlock {
		return fmt.Errorf("the from block should be less than the to block")
	}
	return nil
}

func checkDepositArgs(cmd *cobra.Command, args []string) error {
	if *ulxlyInputArgs.DepositBridgeAddress == "" {
		return fmt.Errorf("please provide the bridge address")
	}
	if *ulxlyInputArgs.DepositGasLimit < 130000 && *ulxlyInputArgs.DepositGasLimit != 0 {
		return fmt.Errorf("the gas limit may be too low for the transaction to pass")
	}
	if *ulxlyInputArgs.DepositMessage && *ulxlyInputArgs.DepositWETH {
		return fmt.Errorf("choose a single deposit mode (asset, message, or WETH)")
	}
	return nil
}

func checkClaimArgs(cmd *cobra.Command, args []string) error {
	if *ulxlyInputArgs.ClaimGasLimit < 150000 && *ulxlyInputArgs.ClaimGasLimit != 0 {
		return fmt.Errorf("the gas limit may be too low for the transaction to pass")
	}
	if *ulxlyInputArgs.ClaimMessage && *ulxlyInputArgs.ClaimWETH {
		return fmt.Errorf("choose a single claim mode (asset, message, or WETH)")
	}
	return nil
}

func init() {
	ULxLyCmd.AddCommand(depositClaimCmd)
	ULxLyCmd.AddCommand(depositNewCmd)
	ULxLyCmd.AddCommand(depositGetCmd)
	ULxLyCmd.AddCommand(ProofCmd)
	ULxLyCmd.AddCommand(EmptyProofCmd)
	ULxLyCmd.AddCommand(ZeroProofCmd)

	ulxlyInputArgs.ClaimIndex = depositClaimCmd.PersistentFlags().String("claim-index", "0", "The deposit count, or index to initiate a claim transaction for.")
	ulxlyInputArgs.ClaimAddress = depositClaimCmd.PersistentFlags().String("claim-address", "", "The address that is receiving the bridged asset.")
	ulxlyInputArgs.ClaimOriginNetwork = depositClaimCmd.PersistentFlags().String("origin-network", "0", "The network ID of the origin network.")
	ulxlyInputArgs.ClaimDestinationNetwork = depositClaimCmd.PersistentFlags().String("destination-network", "1", "The network ID of the destination network.")
	ulxlyInputArgs.ClaimRPCURL = depositClaimCmd.PersistentFlags().String("rpc-url", "http://127.0.0.1:8545", "The RPC endpoint of the destination network")
	ulxlyInputArgs.BridgeServiceRPCURL = depositClaimCmd.PersistentFlags().String("bridge-service-url", "", "The RPC endpoint of the bridge service component.")
	ulxlyInputArgs.ClaimPrivateKey = depositClaimCmd.PersistentFlags().String("private-key", "", "The private key of the sender account.")
	ulxlyInputArgs.ClaimBridgeAddress = depositClaimCmd.PersistentFlags().String("bridge-address", "", "The address of the bridge contract.")
	ulxlyInputArgs.ClaimGasLimit = depositClaimCmd.PersistentFlags().Uint64("gas-limit", 0, "The gas limit for the transaction. Setting this value to 0 will estimate the gas limit.")
	ulxlyInputArgs.ClaimChainID = depositClaimCmd.PersistentFlags().String("chain-id", "", "The chainID.")
	ulxlyInputArgs.ClaimTimeoutTxnReceipt = depositClaimCmd.PersistentFlags().Uint32("transaction-receipt-timeout", 60, "The timeout limit to check for the transaction receipt of the claim.")
	ulxlyInputArgs.ClaimMessage = depositClaimCmd.PersistentFlags().Bool("claim-message", false, "Claim a message instead of an asset.")
	ulxlyInputArgs.ClaimWETH = depositClaimCmd.PersistentFlags().Bool("claim-weth", false, "Claim a weth instead of an asset.")

	ulxlyInputArgs.DepositGasLimit = depositNewCmd.PersistentFlags().Uint64("gas-limit", 0, "The gas limit for the transaction. Setting this value to 0 will estimate the gas limit.")
	ulxlyInputArgs.DepositChainID = depositNewCmd.PersistentFlags().String("chain-id", "", "The chainID.")
	ulxlyInputArgs.DepositPrivateKey = depositNewCmd.PersistentFlags().String("private-key", "", "The private key of the sender account.")
	ulxlyInputArgs.Amount = depositNewCmd.PersistentFlags().Int64("amount", 0, "The amount to send.")
	ulxlyInputArgs.DepositRPCURL = depositNewCmd.PersistentFlags().String("rpc-url", "http://127.0.0.1:8545", "The RPC endpoint of the network")
	ulxlyInputArgs.DepositBridgeAddress = depositNewCmd.PersistentFlags().String("bridge-address", "", "The address of the bridge contract.")
	ulxlyInputArgs.DestinationNetwork = depositNewCmd.PersistentFlags().Uint32("destination-network", 1, "The destination network number.")
	ulxlyInputArgs.DestinationAddress = depositNewCmd.PersistentFlags().String("destination-address", "", "The address of receiver in destination network.")
	ulxlyInputArgs.TokenAddress = depositNewCmd.PersistentFlags().String("token-address", "0x0000000000000000000000000000000000000000", "The address of the token to send.")
	ulxlyInputArgs.IsForced = depositNewCmd.PersistentFlags().Bool("force-update-root", true, "Force the update of the Global Exit Root.")
	ulxlyInputArgs.CallData = depositNewCmd.PersistentFlags().String("call-data", "0x", "For bridging assets - raw data of the call `permit` of the token. For bridging messages - the metadata.")
	ulxlyInputArgs.DepositTimeoutTxnReceipt = depositNewCmd.PersistentFlags().Uint32("transaction-receipt-timeout", 60, "The timeout limit to check for the transaction receipt of the deposit.")
	ulxlyInputArgs.DepositMessage = depositNewCmd.PersistentFlags().Bool("bridge-message", false, "Bridge a message instead of an asset.")
	ulxlyInputArgs.DepositWETH = depositNewCmd.PersistentFlags().Bool("bridge-weth", false, "Bridge a weth instead of an asset.")

	ulxlyInputArgs.FromBlock = depositGetCmd.PersistentFlags().Uint64("from-block", 0, "The block height to start query at.")
	ulxlyInputArgs.ToBlock = depositGetCmd.PersistentFlags().Uint64("to-block", 0, "The block height to start query at.")
	ulxlyInputArgs.RPCURL = depositGetCmd.PersistentFlags().String("rpc-url", "http://127.0.0.1:8545", "The RPC to query for events")
	ulxlyInputArgs.FilterSize = depositGetCmd.PersistentFlags().Uint64("filter-size", 1000, "The batch size for individual filter queries")

	ulxlyInputArgs.BridgeAddress = depositGetCmd.Flags().String("bridge-address", "", "The address of the lxly bridge")
	ulxlyInputArgs.InputFileName = ProofCmd.PersistentFlags().String("file-name", "", "The filename with ndjson data of deposits")
	ulxlyInputArgs.DepositNum = ProofCmd.PersistentFlags().Uint32("deposit-number", 0, "The deposit that we would like to prove")
}
