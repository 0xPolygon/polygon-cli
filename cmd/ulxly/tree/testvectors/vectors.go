package vectors

import (
	"github.com/ethereum/go-ethereum/common"
)

type NullifierLeaf struct {
	NetworkID uint32 `json:"network_id"`
	Index     uint32 `json:"let_index"`
}

type BalanceLeaf struct {
	OriginNetwork      uint32         `json:"origin_network"`
	OriginTokenAddress common.Address `json:"origin_token_address"`
}

type UpdatedLeaf[T BalanceLeaf | NullifierLeaf] struct {
	Key   T           `json:"key"`
	Value common.Hash `json:"value"`
	Path  string      `json:"path"`
}

type Transition[T BalanceLeaf | NullifierLeaf] struct {
	PrevRoot   common.Hash    `json:"prev_root"`
	NewRoot    common.Hash    `json:"new_root"`
	UpdateLeaf UpdatedLeaf[T] `json:"updated_leaf"`
}

type TestVector[T BalanceLeaf | NullifierLeaf] struct {
	Transitions []Transition[T] `json:"transitions"`
}
