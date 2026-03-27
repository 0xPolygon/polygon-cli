// Package exitproof generates Polygon PoS exit proofs for L2→L1 token withdrawals.
package exitproof

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/hashdb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	flagL2RPCURL         = "l2-rpc-url"
	flagL1RPCURL         = "l1-rpc-url"
	flagTxHash           = "tx-hash"
	flagRootChainAddress = "root-chain-address"
	flagCheckpointStride = "checkpoint-stride"
	flagLogIndex         = "log-index"

	defaultCheckpointStride = uint64(10000)

	// headerFetchBatchSize is the number of block headers fetched per RPC batch call.
	headerFetchBatchSize = 100

	// rootChainABI is the minimal ABI for the Polygon PoS RootChain contract.
	rootChainABI = `[
  {"name":"currentHeaderBlock","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint256"}]},
  {"name":"headerBlocks","type":"function","stateMutability":"view",
   "inputs":[{"name":"headerNumber","type":"uint256"}],
   "outputs":[
     {"name":"root","type":"bytes32"},
     {"name":"start","type":"uint256"},
     {"name":"end","type":"uint256"},
     {"name":"createdAt","type":"uint256"},
     {"name":"proposer","type":"address"}
   ]}
]`
)

//go:embed usage.md
var usage string

type inputArgs struct {
	l1RPCURL         string
	l2RPCURL         string
	txHash           string
	rootChainAddr    string
	checkpointStride uint64
	logIndex         uint
}

var args = inputArgs{}

// Cmd is the cobra command for `polycli pos exit-proof`.
var Cmd = &cobra.Command{
	Use:          "exit-proof",
	Short:        "Generate a Polygon PoS exit proof for a burn transaction.",
	Long:         usage,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		if err := util.ValidateURL(args.l1RPCURL); err != nil {
			return fmt.Errorf("--l1-rpc-url: %w", err)
		}
		if err := util.ValidateURL(args.l2RPCURL); err != nil {
			return fmt.Errorf("--l2-rpc-url: %w", err)
		}
		if len(args.txHash) != 66 || !strings.HasPrefix(args.txHash, "0x") {
			return fmt.Errorf("--tx-hash must be a 0x-prefixed 32-byte hex string")
		}
		if !common.IsHexAddress(args.rootChainAddr) {
			return fmt.Errorf("--root-chain-address is not a valid hex address: %s", args.rootChainAddr)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return run(cmd.Context())
	},
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&args.l1RPCURL, flagL1RPCURL, "", "Ethereum RPC URL")
	f.StringVar(&args.l2RPCURL, flagL2RPCURL, "", "Polygon PoS RPC URL")
	f.StringVar(&args.rootChainAddr, flagRootChainAddress, "", "RootChain contract address on L1")
	f.StringVar(&args.txHash, flagTxHash, "", "burn transaction hash on L2")
	f.Uint64Var(&args.checkpointStride, flagCheckpointStride, defaultCheckpointStride, "number of L2 blocks per checkpoint (10000 on mainnet; override for local testnets)")
	f.UintVar(&args.logIndex, flagLogIndex, 0, "index of the burn log within the receipt (0 works for most ERC20 withdrawals; increase if the token emits extra logs before the burn event)")
	_ = Cmd.MarkFlagRequired(flagL1RPCURL)
	_ = Cmd.MarkFlagRequired(flagL2RPCURL)
	_ = Cmd.MarkFlagRequired(flagTxHash)
	_ = Cmd.MarkFlagRequired(flagRootChainAddress)
}

// checkpointInfo holds a single RootChain checkpoint.
type checkpointInfo struct {
	HeaderNumber *big.Int
	Root         [32]byte
	Start        *big.Int
	End          *big.Int
	CreatedAt    *big.Int
	Proposer     common.Address
}

func run(ctx context.Context) error {
	// Connect to L2. We need the raw RPC client for eth_getBlockReceipts and
	// batched eth_getBlockByNumber, and the typed ethclient for everything else.
	rawRPC, err := ethrpc.DialContext(ctx, args.l2RPCURL)
	if err != nil {
		return fmt.Errorf("dial L2: %w", err)
	}
	defer rawRPC.Close()
	l2Client := ethclient.NewClient(rawRPC)

	// Connect to L1.
	l1Client, err := ethclient.DialContext(ctx, args.l1RPCURL)
	if err != nil {
		return fmt.Errorf("dial L1: %w", err)
	}
	defer l1Client.Close()

	txHash := common.HexToHash(args.txHash)

	// Step 1: fetch burn receipt.
	burnReceipt, err := l2Client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return fmt.Errorf("fetch burn receipt: %w", err)
	}
	if burnReceipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("burn transaction %s failed (status=%d)", txHash.Hex(), burnReceipt.Status)
	}
	if uint(args.logIndex) >= uint(len(burnReceipt.Logs)) {
		return fmt.Errorf("log-index %d out of range (receipt has %d logs)", args.logIndex, len(burnReceipt.Logs))
	}
	log.Info().
		Str("txHash", txHash.Hex()).
		Uint64("blockNumber", burnReceipt.BlockNumber.Uint64()).
		Uint("logIndex", args.logIndex).
		Msg("Burn receipt fetched")

	// Step 2: fetch the burn block header.
	burnBlock, err := l2Client.BlockByNumber(ctx, burnReceipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("fetch burn block %d: %w", burnReceipt.BlockNumber.Uint64(), err)
	}

	// Step 3: find the tx's position in the block.
	txIndex, err := findTxIndex(burnBlock, txHash)
	if err != nil {
		return err
	}
	log.Info().
		Uint("txIndex", txIndex).
		Str("receiptsRoot", burnBlock.ReceiptHash().Hex()).
		Msg("Transaction index found")

	// Step 4: fetch all block receipts.
	blockReceipts, err := getBlockReceipts(ctx, rawRPC, burnReceipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("fetch block receipts: %w", err)
	}
	log.Info().Int("count", len(blockReceipts)).Msg("Block receipts fetched")

	// Step 5: build the receipts MPT and generate a proof.
	proofNodes, branchMask, err := buildReceiptProof(blockReceipts, txIndex, burnBlock.ReceiptHash())
	if err != nil {
		return fmt.Errorf("build receipt proof: %w", err)
	}
	log.Info().Int("proofDepth", len(proofNodes)).Msg("Receipt MPT proof generated")

	// Step 6: find the checkpoint that covers the burn block.
	cp, err := findCheckpoint(ctx, l1Client, burnReceipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("find checkpoint: %w", err)
	}
	log.Info().
		Str("checkpointId", cp.HeaderNumber.String()).
		Str("start", cp.Start.String()).
		Str("end", cp.End.String()).
		Msg("Checkpoint found")

	// Step 7: build the binary Merkle block proof.
	blockProof, err := buildBlockProof(ctx, rawRPC, cp, burnReceipt.BlockNumber.Uint64())
	if err != nil {
		return fmt.Errorf("build block proof: %w", err)
	}
	log.Info().Int("proofSiblings", len(blockProof)/32).Msg("Block proof generated")

	// Step 8: RLP-encode the burn receipt.
	receiptBytes, err := burnReceipt.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal burn receipt: %w", err)
	}

	// Step 9: RLP-encode the receipt proof nodes.
	// Each proof node is already RLP-encoded (it's a trie node list like [path, value]).
	// Using rlp.RawValue preserves each node as its own RLP list item in the outer list,
	// so the contract's MerklePatriciaProof.verify can keccak256(toRlpBytes(node)) and
	// match the trie hash. Encoding [][]byte directly would wrap each node in a byte-string
	// prefix, causing a hash mismatch.
	rawProofNodes := make([]rlp.RawValue, len(proofNodes))
	for i, node := range proofNodes {
		rawProofNodes[i] = rlp.RawValue(node)
	}
	rlpProofNodes, err := rlp.EncodeToBytes(rawProofNodes)
	if err != nil {
		return fmt.Errorf("rlp-encode receipt proof: %w", err)
	}

	// Step 10: RLP-encode the full exit payload.
	txRoot := burnBlock.TxHash()
	receiptRoot := burnBlock.ReceiptHash()
	payload, err := encodeExitPayload(
		cp.HeaderNumber,
		blockProof,
		burnReceipt.BlockNumber,
		new(big.Int).SetUint64(burnBlock.Time()),
		txRoot,
		receiptRoot,
		receiptBytes,
		rlpProofNodes,
		branchMask,
		new(big.Int).SetUint64(uint64(args.logIndex)),
	)
	if err != nil {
		return fmt.Errorf("abi-encode exit payload: %w", err)
	}

	// Output: 0x-prefixed hex to stdout for shell capture.
	fmt.Println(hexutil.Encode(payload))
	return nil
}


// getBlockReceipts calls eth_getBlockReceipts for the given block number.
func getBlockReceipts(ctx context.Context, rpc *ethrpc.Client, blockNum *big.Int) ([]*types.Receipt, error) {
	var receipts []*types.Receipt
	if err := rpc.CallContext(ctx, &receipts, "eth_getBlockReceipts", hexutil.EncodeBig(blockNum)); err != nil {
		return nil, fmt.Errorf("eth_getBlockReceipts %s: %w", hexutil.EncodeBig(blockNum), err)
	}
	if len(receipts) == 0 {
		return nil, fmt.Errorf("eth_getBlockReceipts returned empty list for block %s", hexutil.EncodeBig(blockNum))
	}
	return receipts, nil
}

// findTxIndex returns the position of txHash within the block's transaction list.
func findTxIndex(block *types.Block, txHash common.Hash) (uint, error) {
	for i, tx := range block.Transactions() {
		if tx.Hash() == txHash {
			return uint(i), nil
		}
	}
	return 0, fmt.Errorf("transaction %s not found in block %d", txHash.Hex(), block.NumberU64())
}

// buildReceiptProof constructs the receipts MPT and returns the root-to-leaf proof
// nodes and the compact-encoded branch mask for the given tx index.
// Returns an error if the reconstructed trie root does not match expectedRoot.
func buildReceiptProof(receipts []*types.Receipt, txIndex uint, expectedRoot common.Hash) ([][]byte, []byte, error) {
	trieDB := triedb.NewDatabase(rawdb.NewMemoryDatabase(), &triedb.Config{HashDB: hashdb.Defaults})
	t := trie.NewEmpty(trieDB)

	rs := types.Receipts(receipts)
	var buf bytes.Buffer
	for i := 0; i < rs.Len(); i++ {
		key := rlp.AppendUint64(nil, uint64(i)) // matches core/types/hashing.go DeriveSha
		buf.Reset()
		rs.EncodeIndex(i, &buf) // consensus-only encoding identical to DeriveSha
		if err := t.Update(key, buf.Bytes()); err != nil {
			return nil, nil, fmt.Errorf("trie insert receipt %d: %w", i, err)
		}
	}

	gotRoot := t.Hash()
	if gotRoot != expectedRoot {
		return nil, nil, fmt.Errorf("receipts trie root mismatch: computed %s, block header has %s",
			gotRoot.Hex(), expectedRoot.Hex())
	}

	targetKey := rlp.AppendUint64(nil, uint64(txIndex))
	proofDB := memorydb.New()
	if err := t.Prove(targetKey, proofDB); err != nil {
		return nil, nil, fmt.Errorf("trie prove for tx index %d: %w", txIndex, err)
	}

	nodes, err := extractProofNodes(gotRoot, targetKey, proofDB)
	if err != nil {
		return nil, nil, fmt.Errorf("extract ordered proof nodes: %w", err)
	}

	// Encode the trie key as HP (hex-prefix) without the leaf terminator.
	// WithdrawManager.verifyInclusion requires branchMaskBytes[0] == 0, which means
	// the HP prefix byte must be 0x00 (even-length extension path). Including the
	// terminator nibble (16) causes hexToCompact to produce a leaf prefix (0x20),
	// which fails the on-chain check.
	nibbles := keybytesToHex(targetKey)
	mask := hexToCompact(nibbles[:len(nibbles)-1])
	return nodes, mask, nil
}

// extractProofNodes walks the trie proof from root to leaf, returning nodes in
// that order. proofDB maps keccak256(nodeRLP) → nodeRLP.
func extractProofNodes(root common.Hash, key []byte, proofDB *memorydb.Database) ([][]byte, error) {
	hexKey := keybytesToHex(key)
	var nodes [][]byte
	wantHash := root[:]

	for {
		nodeRLP, err := proofDB.Get(wantHash)
		if err != nil || len(nodeRLP) == 0 {
			break
		}
		nodes = append(nodes, nodeRLP)
		nextHash, remaining, err := followTrieNode(nodeRLP, hexKey)
		if err != nil {
			return nil, fmt.Errorf("follow trie node: %w", err)
		}
		if nextHash == nil {
			break // leaf reached
		}
		wantHash = nextHash
		hexKey = remaining
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no proof nodes found for root %s", root.Hex())
	}
	return nodes, nil
}

// followTrieNode decodes a RLP-encoded trie node and returns the hash of the
// next child along hexKey, and the remaining key after consuming the node's prefix.
// Returns (nil, nil, nil) when the node is a leaf (no further traversal needed).
func followTrieNode(nodeRLP []byte, hexKey []byte) (nextHash []byte, remaining []byte, err error) {
	var decoded [][]byte
	if err := rlp.DecodeBytes(nodeRLP, &decoded); err != nil {
		// The node may be an inline node (< 32 bytes); ignore traversal errors at leaves.
		return nil, nil, nil
	}

	switch len(decoded) {
	case 17: // branch node: [child0..child15, value]
		if len(hexKey) == 0 {
			return nil, nil, nil // value slot
		}
		child := decoded[hexKey[0]]
		if len(child) == 0 {
			return nil, nil, nil
		}
		if len(child) == 32 {
			return child, hexKey[1:], nil
		}
		// Inline child — hash it ourselves to look it up (shouldn't happen in practice
		// for non-trivial receipt tries, but handle it gracefully).
		h := crypto.Keccak256(child)
		return h, hexKey[1:], nil

	case 2: // extension or leaf node: [encodedPath, value/hash]
		// Decode the compact-encoded path prefix.
		prefix, isLeaf := compactToHex(decoded[0])
		if len(hexKey) < len(prefix) {
			return nil, nil, nil
		}
		if !bytes.HasPrefix(hexKey, prefix) {
			return nil, nil, nil
		}
		remaining = hexKey[len(prefix):]
		if isLeaf {
			return nil, nil, nil
		}
		// Extension node — decoded[1] is the hash of the next node.
		child := decoded[1]
		if len(child) == 32 {
			return child, remaining, nil
		}
		h := crypto.Keccak256(child)
		return h, remaining, nil

	default:
		return nil, nil, nil
	}
}

// keybytesToHex converts key bytes to hex nibbles with a terminating byte (16).
// This mirrors the unexported trie.keybytesToHex function.
func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	nibbles := make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16 // terminator
	return nibbles
}

// hexToCompact encodes hex nibbles into compact (HP) encoding.
// This mirrors the unexported trie.hexToCompact function.
func hexToCompact(hex []byte) []byte {
	terminator := byte(0)
	if len(hex) > 0 && hex[len(hex)-1] == 16 {
		terminator = 1
		hex = hex[:len(hex)-1]
	}
	buf := make([]byte, len(hex)/2+1)
	buf[0] = terminator << 5
	if len(hex)&1 == 1 {
		buf[0] |= 1 << 4
		buf[0] |= hex[0]
		hex = hex[1:]
	}
	for i := 0; i < len(hex); i += 2 {
		buf[i/2+1] = hex[i]<<4 | hex[i+1]
	}
	return buf
}

// compactToHex decodes a HP-encoded path, returning the nibbles and whether it
// is a leaf (terminator present).
func compactToHex(compact []byte) (nibbles []byte, isLeaf bool) {
	if len(compact) == 0 {
		return nil, false
	}
	first := compact[0]
	isLeaf = first>>5&1 == 1
	oddLen := first>>4&1 == 1

	rest := compact[1:]
	var decoded []byte
	if oddLen {
		decoded = append(decoded, first&0x0f)
	}
	for _, b := range rest {
		decoded = append(decoded, b>>4, b&0x0f)
	}
	return decoded, isLeaf
}

// findCheckpoint binary-searches the RootChain contract for the checkpoint
// whose [start, end] range contains burnBlockNumber.
func findCheckpoint(ctx context.Context, l1Client *ethclient.Client, burnBlockNumber *big.Int) (*checkpointInfo, error) {
	parsedABI, err := abi.JSON(strings.NewReader(rootChainABI))
	if err != nil {
		return nil, fmt.Errorf("parse rootchain ABI: %w", err)
	}
	contractAddr := common.HexToAddress(args.rootChainAddr)
	stride := new(big.Int).SetUint64(args.checkpointStride)

	// currentHeaderBlock() returns the latest submitted checkpoint key (= N * stride).
	currentHeaderBlock, err := callUint256(ctx, l1Client, parsedABI, contractAddr, "currentHeaderBlock")
	if err != nil {
		return nil, fmt.Errorf("get currentHeaderBlock: %w", err)
	}
	if currentHeaderBlock.Sign() == 0 {
		return nil, fmt.Errorf("no checkpoints submitted yet (currentHeaderBlock=0)")
	}

	numCheckpoints := new(big.Int).Div(currentHeaderBlock, stride).Uint64()
	log.Info().
		Uint64("numCheckpoints", numCheckpoints).
		Str("currentHeaderBlock", currentHeaderBlock.String()).
		Msg("Searching for checkpoint")

	lo, hi := uint64(1), numCheckpoints
	for lo <= hi {
		mid := (lo + hi) / 2
		key := new(big.Int).Mul(new(big.Int).SetUint64(mid), stride)
		cp, err := fetchCheckpointByKey(ctx, l1Client, parsedABI, contractAddr, key)
		if err != nil {
			return nil, fmt.Errorf("fetch checkpoint at key %s: %w", key, err)
		}
		switch {
		case burnBlockNumber.Cmp(cp.Start) < 0:
			hi = mid - 1
		case burnBlockNumber.Cmp(cp.End) > 0:
			lo = mid + 1
		default:
			return cp, nil
		}
	}
	return nil, fmt.Errorf("no checkpoint found covering block %s (searched %d checkpoints)", burnBlockNumber, numCheckpoints)
}

// callUint256 calls a no-argument view function that returns a single uint256.
func callUint256(ctx context.Context, l1Client *ethclient.Client, parsedABI abi.ABI, contractAddr common.Address, method string) (*big.Int, error) {
	callData, err := parsedABI.Pack(method)
	if err != nil {
		return nil, fmt.Errorf("pack %s: %w", method, err)
	}
	result, err := l1Client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: callData}, nil)
	if err != nil {
		return nil, fmt.Errorf("call %s: %w", method, err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("%s returned empty data — is --root-chain-address correct? (%s)", method, contractAddr.Hex())
	}
	res, err := parsedABI.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("unpack %s: %w", method, err)
	}
	v, ok := res[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("%s returned unexpected type %T", method, res[0])
	}
	return v, nil
}

// fetchCheckpointByKey fetches a single checkpoint by its raw headerBlocks key.
func fetchCheckpointByKey(ctx context.Context, l1Client *ethclient.Client, parsedABI abi.ABI, contractAddr common.Address, key *big.Int) (*checkpointInfo, error) {
	callData, err := parsedABI.Pack("headerBlocks", key)
	if err != nil {
		return nil, fmt.Errorf("pack headerBlocks: %w", err)
	}
	result, err := l1Client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: callData}, nil)
	if err != nil {
		return nil, fmt.Errorf("call headerBlocks: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("headerBlocks returned empty data — is --root-chain-address correct? (%s)", contractAddr.Hex())
	}
	res, err := parsedABI.Unpack("headerBlocks", result)
	if err != nil {
		return nil, fmt.Errorf("unpack headerBlocks: %w", err)
	}
	return unpackCheckpoint(key, res)
}

// unpackCheckpoint converts the abi.Unpack result of headerBlocks() to checkpointInfo.
func unpackCheckpoint(id *big.Int, res []any) (*checkpointInfo, error) {
	if len(res) != 5 {
		return nil, fmt.Errorf("headerBlocks returned %d values, expected 5", len(res))
	}
	root, ok := res[0].([32]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected root type %T", res[0])
	}
	start, ok := res[1].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected start type %T", res[1])
	}
	end, ok := res[2].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected end type %T", res[2])
	}
	createdAt, ok := res[3].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected createdAt type %T", res[3])
	}
	proposer, ok := res[4].(common.Address)
	if !ok {
		return nil, fmt.Errorf("unexpected proposer type %T", res[4])
	}
	return &checkpointInfo{
		HeaderNumber: new(big.Int).Set(id),
		Root:         root,
		Start:        start,
		End:          end,
		CreatedAt:    createdAt,
		Proposer:     proposer,
	}, nil
}

// buildBlockProof fetches all block hashes in the checkpoint range and constructs
// the binary Merkle proof for burnBlockNumber (sibling hashes concatenated).
func buildBlockProof(ctx context.Context, rpc *ethrpc.Client, cp *checkpointInfo, burnBlockNumber uint64) ([]byte, error) {
	start := cp.Start.Uint64()
	end := cp.End.Uint64()

	log.Info().
		Uint64("start", start).
		Uint64("end", end).
		Uint64("count", end-start+1).
		Msg("Fetching block hashes for checkpoint range")

	hashes, err := fetchBlockHashesBatched(ctx, rpc, start, end)
	if err != nil {
		return nil, fmt.Errorf("fetch block hashes [%d..%d]: %w", start, end, err)
	}

	leafIdx := burnBlockNumber - start
	proof := merkleProof(hashes, leafIdx)
	return proof, nil
}

// fetchBlockHashesBatched fetches block headers for [start, end] and computes the
// Polygon checkpoint Merkle leaf for each block:
//
//	keccak256(abi.encodePacked(blockNumber_32, blockTime_32, txRoot_32, receiptsRoot_32))
//
// This matches the leaf formula used by getBlockHeader() in the matic-js SDK and
// verified by WithdrawManager.checkBlockMembershipInCheckpoint on-chain.
func fetchBlockHashesBatched(ctx context.Context, rpc *ethrpc.Client, start, end uint64) ([]common.Hash, error) {
	count := end - start + 1
	hashes := make([]common.Hash, count)

	type blockHeader struct {
		Number           hexutil.Big  `json:"number"`
		Timestamp        hexutil.Big  `json:"timestamp"`
		TransactionsRoot common.Hash  `json:"transactionsRoot"`
		ReceiptsRoot     common.Hash  `json:"receiptsRoot"`
	}

	for batchStart := uint64(0); batchStart < count; batchStart += headerFetchBatchSize {
		batchEnd := batchStart + headerFetchBatchSize
		if batchEnd > count {
			batchEnd = count
		}
		batchLen := batchEnd - batchStart

		elems := make([]ethrpc.BatchElem, batchLen)
		results := make([]blockHeader, batchLen)
		for i := uint64(0); i < batchLen; i++ {
			blockNum := start + batchStart + i
			elems[i] = ethrpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []any{hexutil.EncodeUint64(blockNum), false},
				Result: &results[i],
			}
		}

		if err := rpc.BatchCallContext(ctx, elems); err != nil {
			return nil, fmt.Errorf("batch RPC call: %w", err)
		}
		for i, elem := range elems {
			if elem.Error != nil {
				blockNum := start + batchStart + uint64(i)
				return nil, fmt.Errorf("fetch block %d: %w", blockNum, elem.Error)
			}
			h := results[i]
			// Compute the Polygon checkpoint leaf: keccak256(n_32 || ts_32 || txRoot || receiptsRoot)
			var buf [128]byte
			h.Number.ToInt().FillBytes(buf[0:32])
			h.Timestamp.ToInt().FillBytes(buf[32:64])
			copy(buf[64:96], h.TransactionsRoot[:])
			copy(buf[96:128], h.ReceiptsRoot[:])
			hashes[batchStart+uint64(i)] = crypto.Keccak256Hash(buf[:])
		}

		log.Debug().
			Uint64("fetched", batchStart+batchLen).
			Uint64("total", count).
			Msg("Block headers fetched")
	}

	return hashes, nil
}

// merkleProof builds a binary Merkle tree from the given leaf hashes and returns
// the concatenated sibling hashes (proof) for leafIdx.
// Construction matches the matic.js MerkleTree: odd-length layers duplicate the last leaf.
// Internal nodes: keccak256(left || right).
func merkleProof(leaves []common.Hash, leafIdx uint64) []byte {
	layer := make([]common.Hash, len(leaves))
	copy(layer, leaves)

	var siblings []common.Hash
	pos := leafIdx

	for len(layer) > 1 {
		if len(layer)%2 == 1 {
			layer = append(layer, layer[len(layer)-1])
		}
		sibling := pos ^ 1
		siblings = append(siblings, layer[sibling])

		next := make([]common.Hash, len(layer)/2)
		for i := range next {
			next[i] = crypto.Keccak256Hash(layer[i*2][:], layer[i*2+1][:])
		}
		layer = next
		pos = pos / 2
	}

	result := make([]byte, len(siblings)*32)
	for i, h := range siblings {
		copy(result[i*32:], h[:])
	}
	return result
}

// encodeExitPayload RLP-encodes the 10-field exit payload for startExitWithBurntTokens(bytes).
// The Polygon contracts (ExitPayloadReader.toExitPayload) RLP-decode the payload, so it must
// be an RLP list — not ABI-encoded. The format matches the matic-js buildReferenceTxPayload:
// [headerNumber, blockProof, blockNumber, blockTimestamp, txRoot, receiptRoot, receipt,
//  receiptParentNodes, branchMask, logIndex]
func encodeExitPayload(
	headerNumber *big.Int,
	blockProof []byte,
	blockNumber *big.Int,
	blockTimestamp *big.Int,
	txRoot common.Hash,
	receiptRoot common.Hash,
	receipt []byte,
	receiptParentNodes []byte,
	branchMask []byte,
	logIndex *big.Int,
) ([]byte, error) {
	type rlpPayload struct {
		HeaderNumber       *big.Int
		BlockProof         []byte
		BlockNumber        *big.Int
		BlockTimestamp     *big.Int
		TxRoot             []byte
		ReceiptRoot        []byte
		Receipt            []byte
		ReceiptParentNodes []byte
		BranchMask         []byte
		LogIndex           *big.Int
	}

	p := rlpPayload{
		HeaderNumber:       headerNumber,
		BlockProof:         blockProof,
		BlockNumber:        blockNumber,
		BlockTimestamp:     blockTimestamp,
		TxRoot:             txRoot[:],
		ReceiptRoot:        receiptRoot[:],
		Receipt:            receipt,
		ReceiptParentNodes: receiptParentNodes,
		BranchMask:         branchMask,
		LogIndex:           logIndex,
	}

	return rlp.EncodeToBytes(p)
}
