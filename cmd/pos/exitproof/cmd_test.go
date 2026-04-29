package exitproof

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// verifyProof reproduces pos-contracts/contracts/common/lib/Merkle.sol's
// checkMembership: walk the proof, hashing the running value against each
// sibling on the side dictated by the path bits of `index`.
func verifyProof(leaf common.Hash, index uint64, root common.Hash, proof []byte) bool {
	if len(proof)%32 != 0 {
		return false
	}
	h := leaf
	for i := 0; i < len(proof); i += 32 {
		var sib common.Hash
		copy(sib[:], proof[i:i+32])
		if index%2 == 0 {
			h = crypto.Keccak256Hash(h[:], sib[:])
		} else {
			h = crypto.Keccak256Hash(sib[:], h[:])
		}
		index /= 2
	}
	return h == root
}

// computeRoot rebuilds the matic-js / pos-contracts MerkleTree root: pad
// leaves to next power of 2 with zero hashes, then hash pairwise up.
func computeRoot(leaves []common.Hash) common.Hash {
	if len(leaves) == 0 {
		return common.Hash{}
	}
	size := 1
	for size < len(leaves) {
		size <<= 1
	}
	layer := make([]common.Hash, size)
	copy(layer, leaves)
	for len(layer) > 1 {
		next := make([]common.Hash, len(layer)/2)
		for i := range next {
			next[i] = crypto.Keccak256Hash(layer[i*2][:], layer[i*2+1][:])
		}
		layer = next
	}
	return layer[0]
}

// makeLeaves returns n distinct hashes derived from a counter — enough variety
// to make collisions vanishingly unlikely.
func makeLeaves(n int) []common.Hash {
	out := make([]common.Hash, n)
	for i := range out {
		out[i] = crypto.Keccak256Hash([]byte{byte(i), byte(i >> 8)})
	}
	return out
}

func TestMerkleProof_VerifiesEveryLeaf(t *testing.T) {
	// Cover both power-of-2 sizes (where the previous duplicate-last-node
	// algorithm coincidentally agreed with matic-js) and non-power-of-2 sizes
	// (where it didn't — and produced WITHDRAW_BLOCK_NOT_A_PART_OF_SUBMITTED_HEADER
	// reverts on fast-cadence devnets).
	for _, n := range []int{1, 2, 3, 4, 5, 7, 8, 16, 40, 128, 200} {
		leaves := makeLeaves(n)
		root := computeRoot(leaves)
		for idx := 0; idx < n; idx++ {
			proof := merkleProof(leaves, uint64(idx))
			if !verifyProof(leaves[idx], uint64(idx), root, proof) {
				t.Errorf("n=%d idx=%d: proof did not verify", n, idx)
			}
		}
	}
}

// TestMerkleProof_NonPowerOfTwoBugRegression locks in the specific shape that
// failed in production: a 40-leaf checkpoint range on the kurtosis-pos devnet.
// Before the fix, the proof for any leaf in this size verified against a
// duplicate-padded root but not against the matic-js / on-chain root.
func TestMerkleProof_NonPowerOfTwoBugRegression(t *testing.T) {
	leaves := makeLeaves(40)
	root := computeRoot(leaves)
	for idx := 0; idx < 40; idx++ {
		proof := merkleProof(leaves, uint64(idx))
		if !verifyProof(leaves[idx], uint64(idx), root, proof) {
			t.Fatalf("idx=%d: proof did not verify against canonical root", idx)
		}
	}
}
