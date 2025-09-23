package agglayer

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/ethereum/go-ethereum/common"
)

type getClaimProofResponse struct {
	L1InfoTreeLeafResponse l1InfoTreeLeafResponse `json:"l1_info_tree_leaf"`
	ProofLocalExitRoot     []string               `json:"proof_local_exit_root"`
	ProofRollupExitRoot    []string               `json:"proof_rollup_exit_root"`
}

type l1InfoTreeLeafResponse struct {
	BlockNum          uint64 `json:"block_num"`
	BlockPos          uint64 `json:"block_pos"`
	GlobalExitRoot    string `json:"global_exit_root"`
	Hash              string `json:"hash"`
	L1InfoTreeIndex   uint64 `json:"l1_info_tree_index"`
	MainnetExitRoot   string `json:"mainnet_exit_root"`
	PreviousBlockHash string `json:"previous_block_hash"`
	RollupExitRoot    string `json:"rollup_exit_root"`
	Timestamp         uint64 `json:"timestamp"`
}

func (r *getClaimProofResponse) ToProof() *bridge_service.Proof {
	p := &bridge_service.Proof{}

	var merkleProof = make([]common.Hash, len(r.ProofLocalExitRoot))
	for i, p := range r.ProofLocalExitRoot {
		merkleProof[i] = common.HexToHash(p)
	}
	p.MerkleProof = merkleProof

	var rollupMerkleProof = make([]common.Hash, len(r.ProofRollupExitRoot))
	for i, p := range r.ProofRollupExitRoot {
		rollupMerkleProof[i] = common.HexToHash(p)
	}
	p.RollupMerkleProof = rollupMerkleProof

	if len(r.L1InfoTreeLeafResponse.MainnetExitRoot) > 0 {
		mainExitRoot := common.HexToHash(r.L1InfoTreeLeafResponse.MainnetExitRoot)
		p.MainExitRoot = &mainExitRoot
	}

	if len(r.L1InfoTreeLeafResponse.RollupExitRoot) > 0 {
		rollupExitRoot := common.HexToHash(r.L1InfoTreeLeafResponse.RollupExitRoot)
		p.RollupExitRoot = &rollupExitRoot
	}

	return p
}
