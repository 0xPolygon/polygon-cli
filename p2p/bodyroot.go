package p2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
)

// bodyRoot is a single hash over every part of a block body that its header
// commits to: the transaction root, the uncle hash and, when the block has
// them, the withdrawals root. A body matches its header exactly when the two
// bodyRoots are equal, so checking it covers the uncles and withdrawals as
// well as the transactions.
func bodyRoot(txRoot, uncleHash common.Hash, withdrawalsRoot *common.Hash) common.Hash {
	var w common.Hash
	if withdrawalsRoot != nil {
		w = *withdrawalsRoot
	}
	return crypto.Keccak256Hash(txRoot[:], uncleHash[:], w[:])
}

// headerBodyRoot is the bodyRoot a header commits its block body to.
func headerBodyRoot(h *types.Header) common.Hash {
	return bodyRoot(h.TxHash, h.UncleHash, h.WithdrawalsHash)
}

// decodedBodyRoot is the bodyRoot of a decoded body. Withdrawals count only when
// the body carries a withdrawals list, as a header's withdrawals root does.
func decodedBodyRoot(txs *bortx.Body, uncles []*types.Header, withdrawals []*types.Withdrawal) common.Hash {
	var w *common.Hash
	if withdrawals != nil {
		root := types.DeriveSha(types.Withdrawals(withdrawals), trie.NewStackTrie(nil))
		w = &root
	}
	return bodyRoot(txs.Root(), types.CalcUncleHash(uncles), w)
}

// rawBodyRoot is the bodyRoot of a body whose uncles and withdrawals are still
// encoded, as a BlockBodies response carries them.
func rawBodyRoot(body *eth.BlockBody, txs *bortx.Body) (common.Hash, error) {
	uncles, err := body.Uncles.Items()
	if err != nil {
		return common.Hash{}, fmt.Errorf("decode uncles: %w", err)
	}
	var withdrawals []*types.Withdrawal
	if body.Withdrawals != nil {
		if withdrawals, err = body.Withdrawals.Items(); err != nil {
			return common.Hash{}, fmt.Errorf("decode withdrawals: %w", err)
		}
		if withdrawals == nil {
			withdrawals = []*types.Withdrawal{}
		}
	}
	return decodedBodyRoot(txs, uncles, withdrawals), nil
}
