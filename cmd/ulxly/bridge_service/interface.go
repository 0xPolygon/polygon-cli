package bridge_service

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrUnableToRetrieveDeposit = errors.New("the bridge deposit was not found")
)

type BridgeService interface {
	GetDeposit(depositNetwork, depositCount uint32) (*Deposit, error)
	GetDeposits(destinationAddress string, offset, limit int) (deposits []Deposit, total int, err error)
	GetProof(depositNetwork, depositCount uint32) (*Proof, error)
	Url() string
}

type BridgeServiceBase struct {
	url string
}

func (b *BridgeServiceBase) Url() string {
	return b.url
}

func NewBridgeServiceBase(url string) BridgeServiceBase {
	return BridgeServiceBase{
		url: url,
	}
}

type Deposit struct {
	LeafType      uint8
	OrigNet       uint32
	OrigAddr      common.Address
	Amount        *big.Int
	DestNet       uint32
	DestAddr      common.Address
	BlockNum      uint64
	DepositCnt    uint32
	NetworkID     uint32
	TxHash        common.Hash
	ClaimTxHash   *common.Hash
	Metadata      []byte
	ReadyForClaim bool
	GlobalIndex   *big.Int
}

type Proof struct {
	MerkleProof       []common.Hash
	RollupMerkleProof []common.Hash
	MainExitRoot      *common.Hash
	RollupExitRoot    *common.Hash
}

func HashArrayToBytesArray(hashes []common.Hash) [32][32]byte {
	var array [32][32]byte
	for i, h := range hashes {
		if i >= 32 {
			break
		}
		array[i] = h
	}
	return array
}
