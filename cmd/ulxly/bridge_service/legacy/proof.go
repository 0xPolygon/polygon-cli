package legacy

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/ethereum/go-ethereum/common"
)

type GetProofResponse struct {
	Proof ProofResponse `json:"proof"`
}

type ProofResponse struct {
	MerkleProof       []string `json:"merkle_proof"`
	RollupMerkleProof []string `json:"rollup_merkle_proof"`
	MainExitRoot      string   `json:"main_exit_root"`
	RollupExitRoot    string   `json:"rollup_exit_root"`
}

func (r *ProofResponse) ToProof() *bridge_service.Proof {
	p := &bridge_service.Proof{}

	var merkleProof = make([]common.Hash, len(r.MerkleProof))
	for i, p := range r.MerkleProof {
		merkleProof[i] = common.HexToHash(p)
	}

	var rollupMerkleProof = make([]common.Hash, len(r.RollupMerkleProof))
	for i, p := range r.RollupMerkleProof {
		rollupMerkleProof[i] = common.HexToHash(p)
	}

	if len(r.MainExitRoot) > 0 {
		mainExitRoot := common.HexToHash(r.MainExitRoot)
		p.MainExitRoot = &mainExitRoot
	}

	if len(r.RollupExitRoot) > 0 {
		rollupExitRoot := common.HexToHash(r.RollupExitRoot)
		p.RollupExitRoot = &rollupExitRoot
	}

	p.MerkleProof = merkleProof
	p.RollupMerkleProof = rollupMerkleProof

	return p
}
