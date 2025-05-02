package ulxly

import (
	"math/big"

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

var (
	zeroHashes []common.Hash
	siblings   []common.Hash
)

type Balancer struct {
	zeroHashes []common.Hash
	siblings   []common.Hash
	depth      int
}

func NewBalanceTree() (*Balancer, error) {
	var depth uint8 = 192
	zeroHashes = generateZeroHashes(depth)
	var err error
	siblings, err = initSiblings(depth)
	if err != nil {
		return nil, err
	}
	return &Balancer{
		zeroHashes: zeroHashes,
		siblings:   siblings,
		depth:      int(depth),
	}, nil
}

func (b *Balancer) UpdateBalanceTree(token TokenInfo, leaf *big.Int) (common.Hash, error) {
	key := token.ToBits()
	return b.updateSMT(key, FromU256(leaf))
}

func (p *Balancer) updateSMT(
	key []bool,
	newValue common.Hash,
) (common.Hash, error) {
	if len(key) != p.depth || len(p.siblings) != p.depth {
		return common.Hash{}, nil
	}
	hash := newValue
	for i := 0; i < p.depth; i++ {
		index := p.depth - i - 1
		sibling := p.siblings[i]
		if key[index] {
			hash = crypto.Keccak256Hash(sibling.Bytes(), hash.Bytes())
		} else {
			hash = crypto.Keccak256Hash(hash.Bytes(), sibling.Bytes())
		}
	}

	return hash, nil
}

type Nullifier struct {
	zeroHashes []common.Hash
	siblings   []common.Hash
	depth      uint8
}

func NewNullifierTree() (*Nullifier, error) {
	var depth uint8 = 64
	zeroHashes = generateZeroHashes(depth)
	var err error
	siblings, err = initSiblings(depth)
	if err != nil {
		return nil, err
	}
	return &Nullifier{
		zeroHashes: zeroHashes,
		siblings:   siblings,
		depth:      depth,
	}, nil
}

func (n *Nullifier) UpdateNullifierTree(nullifier NullifierKey) common.Hash {
	key := nullifier.ToBits()
	return updateSMTree(key, FromBool(true), n.depth)
}

// updateSMTree returns the updated root
func updateSMTree(
	key []bool,
	newValue common.Hash,
	depth uint8,
) common.Hash {
	entry := newValue
	// First loop: From sibling length up to depth. It "fills in" the tree from the leaf up to the first real sibling.
	for i := int(depth) - 1; i >= len(siblings); i-- {
		sibling := zeroHashes[int(depth)-i-1]
		if key[i] {
			entry = crypto.Keccak256Hash(sibling.Bytes(), entry[:])
		} else {
			entry = crypto.Keccak256Hash(entry[:], sibling.Bytes())
		}
	}

	// Second loop: For provided siblings. This is the classic proof verification and root recomputation step
	for i := len(siblings) - 1; i >= 0; i-- {
		sibling := siblings[i]
		if key[i] {
			entry = crypto.Keccak256Hash(sibling.Bytes(), entry[:])
		} else {
			entry = crypto.Keccak256Hash(entry[:], sibling.Bytes())
		}
	}

	return entry
}

// initSiblings returns the siblings of the node at the given index.
// it is used to initialize the siblings array in the beginning.
func initSiblings(height uint8) ([]common.Hash, error) {
	var (
		left     common.Hash
		siblings []common.Hash
	)
	for h := 0; h < int(height); h++ {
		copy(left[:], zeroHashes[h][:])
		siblings = append(siblings, left)
	}
	return siblings, nil
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
		out[31] = 1 // Set the last byte to 1 (little-end)
	}
	return out
}
