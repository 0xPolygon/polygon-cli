package ulxly

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func TokenInfoStringToStruct(key string) (TokenInfo, error){
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
	fmt.Println("Initial Root: ", initRoot.String())
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
	fmt.Println("Initial Root: ", initRoot.String())
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
