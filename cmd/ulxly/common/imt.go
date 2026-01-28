package common

import (
	"encoding/binary"
	"fmt"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
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

type RollupsProof struct {
	Siblings [TreeDepth]common.Hash
	Root     common.Hash
	RollupID uint32
	LeafHash common.Hash
}

// Init will allocate the objects in the IMT
func (s *IMT) Init() {
	s.Branches = make(map[uint32][]common.Hash)
	s.Leaves = make(map[uint32]common.Hash)
	s.ZeroHashes = GenerateZeroHashes(TreeDepth)
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
		branches = GenerateEmptyHashes(TreeDepth)
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
		siblingDepositNum := GetSiblingLeafNumber(depositNum, uint32(height))
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

// GetSiblingLeafNumber returns the sibling number of a given number at a specified level in an incremental Merkle tree.
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
func GetSiblingLeafNumber(leafNumber, level uint32) uint32 {
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

// GenerateZeroHashes generates zero hashes for the incremental Merkle tree
// https://eth2book.info/capella/part2/deposits-withdrawals/contract/
func GenerateZeroHashes(height uint8) []common.Hash {
	zeroHashes := make([]common.Hash, height)
	zeroHashes[0] = common.Hash{}
	for i := 1; i < int(height); i++ {
		zeroHashes[i] = crypto.Keccak256Hash(zeroHashes[i-1][:], zeroHashes[i-1][:])
	}
	return zeroHashes
}

// GenerateEmptyHashes generates empty hashes for the incremental Merkle tree
func GenerateEmptyHashes(height uint8) []common.Hash {
	zeroHashes := make([]common.Hash, height)
	zeroHashes[0] = common.Hash{}
	for i := 1; i < int(height); i++ {
		zeroHashes[i] = common.Hash{}
	}
	return zeroHashes
}

// ComputeRoot computes the root hash given a leaf, proof siblings, index, and height
func ComputeRoot(leafHash common.Hash, smtProof [32]common.Hash, index uint32, height uint8) common.Hash {
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

// ComputeSiblings computes the siblings for a given rollup ID and leaves
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
		if index%2 == 1 { // If it is odd
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

	computedRoot := ComputeRoot(p.LeafHash, p.Siblings, p.RollupID, TreeDepth)
	if computedRoot != p.Root {
		return nil, fmt.Errorf("error: computed root does not match the expected root")
	}

	return p, nil
}

// HashDeposit creates the leaf hash value for a particular deposit
func HashDeposit(deposit *ulxly.UlxlyBridgeEvent) common.Hash {
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
