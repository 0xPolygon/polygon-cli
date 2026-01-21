package tree

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/0xPolygon/cdk-rpc/types"
	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang-collections/collections/stack"
	"github.com/rs/zerolog/log"
)

// TokenInfo struct
type TokenInfo struct {
	OriginNetwork      *big.Int
	OriginTokenAddress common.Address // 20 bytes, Ethereum address
}

// ToBits convert TokenInfo to an array of 192 bits (bool)
func (t *TokenInfo) ToBits() []bool {
	bits := make([]bool, 192)

	// First 32 bits: OriginNetwork
	for i := 0; i < 32; i++ {
		if t.OriginNetwork.Bit(i) == 1 {
			bits[i] = true
		}
	}

	// The next 160 bits: OriginTokenAddress (20 bytes * 8 bits = 160)
	for i := 32; i < 192; i++ {
		byteIndex := (i - 32) / 8
		bitIndex := (i % 8)
		if (t.OriginTokenAddress.Bytes()[byteIndex]>>bitIndex)&1 == 1 {
			bits[i] = true
		}
	}

	return bits
}

func (t *TokenInfo) String() string {
	return fmt.Sprintf("%s-%s", t.OriginNetwork.String(), t.OriginTokenAddress.Hex())
}

func TokenInfoStringToStruct(key string) (TokenInfo, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return TokenInfo{}, fmt.Errorf("invalid key format: %s", key)
	}

	originNetwork, ok := big.NewInt(0).SetString(parts[0], 10) // Parse the first part as a big.Int
	if !ok {
		return TokenInfo{}, fmt.Errorf("invalid origin network value: %s", parts[0])
	}

	originTokenAddress := common.HexToAddress(parts[1]) // Parse the second part as an address

	return TokenInfo{
		OriginNetwork:      originNetwork,
		OriginTokenAddress: originTokenAddress,
	}, nil
}

// NullifierKey struct
type NullifierKey struct {
	NetworkID uint32
	Index     uint32
}

func (n *NullifierKey) ToBits() []bool {
	bits := make([]bool, 64)

	// First 32 bits: NetworkID
	for i := 0; i < 32; i++ {
		if (n.NetworkID>>i)&1 == 1 {
			bits[i] = true
		}
	}

	// Next 32 bits: Index
	for i := 0; i < 32; i++ {
		if (n.Index>>i)&1 == 1 {
			bits[i+32] = true
		}
	}

	return bits
}

type Tree struct {
	zeroHashes []common.Hash
	depth      uint8
	Tree       map[common.Hash]Node
}
type Balancer struct {
	tree     Tree
	lastRoot common.Hash
}

func NewBalanceTree() (*Balancer, error) {
	var depth uint8 = 192
	zeroHashes := generateZeroHashes(depth)
	initRoot := crypto.Keccak256Hash(zeroHashes[depth-1].Bytes(), zeroHashes[depth-1].Bytes())
	log.Info().Msg("Initial Root: " + initRoot.String())
	return &Balancer{
		tree: Tree{
			zeroHashes: zeroHashes,
			depth:      depth,
			Tree:       make(map[common.Hash]Node),
		},
		lastRoot: initRoot,
	}, nil
}

func (b *Balancer) UpdateBalanceTree(token TokenInfo, leaf *big.Int) (common.Hash, error) {
	key := token.ToBits()
	newRoot, err := b.tree.insertHelper(b.lastRoot, 0, key, FromU256(leaf), true)
	if err != nil {
		return common.Hash{}, err
	}
	b.lastRoot = newRoot
	return newRoot, nil
}

type Nullifier struct {
	tree     Tree
	lastRoot common.Hash
}

func NewNullifierTree() (*Nullifier, error) {
	var depth uint8 = 64
	zeroHashes := generateZeroHashes(depth)
	initRoot := crypto.Keccak256Hash(zeroHashes[depth-1].Bytes(), zeroHashes[depth-1].Bytes())
	log.Info().Msg("Initial Root: " + initRoot.String())
	return &Nullifier{
		tree: Tree{
			zeroHashes: zeroHashes,
			depth:      depth,
			Tree:       make(map[common.Hash]Node),
		},
		lastRoot: initRoot,
	}, nil
}

func (n *Nullifier) UpdateNullifierTree(nullifier NullifierKey) (common.Hash, error) {
	key := nullifier.ToBits()
	newRoot, err := n.tree.insertHelper(n.lastRoot, 0, key, FromBool(true), false)
	if err != nil {
		return common.Hash{}, err
	}
	n.lastRoot = newRoot
	return newRoot, nil
}

func FromU256(u *big.Int) common.Hash {
	var aux [32]byte
	// Get the byte slice in big-endian format
	bytes := u.Bytes()

	// Fill the last bytes (right-aligned) of out
	copy(aux[32-len(bytes):], bytes)
	return aux
}

func FromBool(b bool) common.Hash {
	var out [32]byte
	if b {
		out[0] = 1
	}
	return out
}

type Node struct {
	Left  common.Hash
	Right common.Hash
}

func (t *Tree) insertHelper(
	hash common.Hash,
	depth uint8,
	bits []bool,
	value common.Hash,
	update bool,
) (common.Hash, error) {
	if depth > t.depth {
		return common.Hash{}, fmt.Errorf("depth exceeds maximum")
	}
	if depth == t.depth {
		if !update && hash != t.zeroHashes[0] {
			return common.Hash{}, fmt.Errorf("key already exists")
		}
		return value, nil
	}

	// Get node at this hash or initialize a default one
	node, ok := t.Tree[hash]
	if !ok {
		defaultChild := t.zeroHashes[t.depth-depth-1]
		node = Node{
			Left:  defaultChild,
			Right: defaultChild,
		}
	}

	// Recurse to update or insert value
	var childHash common.Hash
	var err error
	if bits[depth] {
		childHash, err = t.insertHelper(node.Right, depth+1, bits, value, update)
		if err != nil {
			return common.Hash{}, err
		}
		node.Right = childHash
	} else {
		childHash, err = t.insertHelper(node.Left, depth+1, bits, value, update)
		if err != nil {
			return common.Hash{}, err
		}
		node.Left = childHash
	}

	// Compute hash of updated node and store
	newHash := crypto.Keccak256Hash(node.Left.Bytes(), node.Right.Bytes())
	t.Tree[newHash] = node

	return newHash, nil
}

var methodIDClaimMessage = common.Hex2Bytes("f5efcd79")

func IsMessageClaim(input []byte) (bool, error) {
	methodID := input[:4]
	// Ignore ClaimAsset method
	if bytes.Equal(methodID, methodIDClaimMessage) {
		return true, nil
	} else {
		return false, nil
	}
}

type call struct {
	To    common.Address `json:"to"`
	Value *types.ArgBig  `json:"value"`
	Err   *string        `json:"error"`
	Input types.ArgBytes `json:"input"`
	Calls []call         `json:"calls"`
}

type tracerCfg struct {
	Tracer string `json:"tracer"`
}

func checkClaimCalldata(client *ethclient.Client, bridge common.Address, claimHash common.Hash) (bool, error) {
	c := &call{}
	err := client.Client().Call(c, "debug_traceTransaction", claimHash, tracerCfg{Tracer: "callTracer"})
	if err != nil {
		return false, err
	}

	// find the claim linked to the event using DFS
	callStack := stack.New()
	callStack.Push(*c)
	for {
		if callStack.Len() == 0 {
			break
		}

		currentCallInterface := callStack.Pop()
		currentCall, ok := currentCallInterface.(call)
		if !ok {
			return false, fmt.Errorf("unexpected type for 'currentCall'. Expected 'call', got '%T'", currentCallInterface)
		}

		if currentCall.To == bridge {
			isMessage, err := IsMessageClaim(currentCall.Input)
			if err != nil {
				return false, err
			}
			return isMessage, err
		}
		for _, c := range currentCall.Calls {
			callStack.Push(c)
		}
	}
	return false, fmt.Errorf("claim not found")
}

func Uint32ToBytesLittleEndian(num uint32) []byte {
	bytes := make([]byte, 4) // uint32 is 4 bytes
	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
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

// computeNullifierTree computes the nullifier tree root from raw claims data
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
		err = json.Unmarshal(scanner.Bytes(), claim)
		if err != nil {
			return common.Hash{}, err
		}
		mainnetFlag, rollupIndex, localExitRootIndex, err := ulxlycommon.DecodeGlobalIndex(claim.GlobalIndex)
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

// computeBalanceTree computes the balance tree root from claims and deposits data
func computeBalanceTree(client *ethclient.Client, bridgeAddress common.Address, l2RawClaims []byte, l2NetworkID uint32, l2RawDeposits []byte) (common.Hash, map[string]*big.Int, error) {
	buf := bytes.NewBuffer(l2RawClaims)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	bTree, err := NewBalanceTree()
	if err != nil {
		return common.Hash{}, nil, err
	}
	balances := make(map[string]*big.Int)
	for scanner.Scan() {
		l2Claim := new(ulxly.UlxlyClaimEvent)
		err = json.Unmarshal(scanner.Bytes(), l2Claim)
		if err != nil {
			return common.Hash{}, nil, err
		}
		token := TokenInfo{
			OriginNetwork:      big.NewInt(0).SetUint64(uint64(l2Claim.OriginNetwork)),
			OriginTokenAddress: l2Claim.OriginAddress,
		}
		isMessage, err := checkClaimCalldata(client, bridgeAddress, l2Claim.Raw.TxHash)
		if err != nil {
			return common.Hash{}, nil, err
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
			return common.Hash{}, nil, err
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
			return common.Hash{}, nil, err
		}
		if token.OriginNetwork.Uint64() == uint64(l2NetworkID) {
			continue
		}
		root, err = bTree.UpdateBalanceTree(token, balance)
		if err != nil {
			return common.Hash{}, nil, err
		}
		log.Info().Msgf("New balanceTree leaf. OriginNetwork: %s, TokenAddress: %s, Balance: %s, Root: %s", token.OriginNetwork.String(), token.OriginTokenAddress.String(), balance.String(), root.String())
	}
	log.Info().Msgf("Final balanceTree root: %s", root.String())

	return root, balances, nil
}

// FileOptions holds options for file input
type FileOptions struct {
	FileName string
}

// BalanceTreeOptions holds options for the balance tree command
type BalanceTreeOptions struct {
	L2ClaimsFile, L2DepositsFile, BridgeAddress, RpcURL string
	L2NetworkID                                         uint32
	Insecure                                            bool
}

// getInputData reads input data from file, args, or stdin
func getInputData(args []string, fileName string) ([]byte, error) {
	if fileName != "" {
		return os.ReadFile(fileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

// getBalanceTreeData reads the claims and deposits files
func getBalanceTreeData(opts *BalanceTreeOptions) ([]byte, []byte, error) {
	claimsFileName := opts.L2ClaimsFile
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

	l2FileName := opts.L2DepositsFile
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
