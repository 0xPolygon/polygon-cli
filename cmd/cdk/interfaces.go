package cdk

import (
	_ "embed"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type rollupManagerContractInterface interface {
	// rollup manager methods
	BridgeAddress(opts *bind.CallOpts) (common.Address, error)
	CalculateRewardPerBatch(opts *bind.CallOpts) (*big.Int, error)
	GetBatchFee(opts *bind.CallOpts) (*big.Int, error)
	GetForcedBatchFee(opts *bind.CallOpts) (*big.Int, error)
	GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error)
	GetRollupExitRoot(opts *bind.CallOpts) ([32]byte, error)
	GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error)
	HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error)
	IsEmergencyState(opts *bind.CallOpts) (bool, error)
	LastAggregationTimestamp(opts *bind.CallOpts) (uint64, error)
	LastDeactivatedEmergencyStateTimestamp(opts *bind.CallOpts) (uint64, error)
	MultiplierBatchFee(opts *bind.CallOpts) (uint16, error)
	PendingStateTimeout(opts *bind.CallOpts) (uint64, error)
	Pol(opts *bind.CallOpts) (common.Address, error)
	RollupCount(opts *bind.CallOpts) (uint32, error)
	RollupTypeCount(opts *bind.CallOpts) (uint32, error)
	TotalSequencedBatches(opts *bind.CallOpts) (uint64, error)
	TotalVerifiedBatches(opts *bind.CallOpts) (uint64, error)
	TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error)
	VerifyBatchTimeTarget(opts *bind.CallOpts) (uint64, error)

	// rollup methods
	ChainIDToRollupID(opts *bind.CallOpts, chainID uint64) (uint32, error)
	RollupAddressToID(opts *bind.CallOpts, rollupAddress common.Address) (uint32, error)
	GetLastVerifiedBatch(opts *bind.CallOpts, rollupID uint32) (uint64, error)
	GetRollupBatchNumToStateRoot(opts *bind.CallOpts, rollupID uint32, batchNum uint64) ([32]byte, error)
	GetInputSnarkBytes(opts *bind.CallOpts, rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error)
	IsPendingStateConsolidable(opts *bind.CallOpts, rollupID uint32, pendingStateNum uint64) (bool, error)
	RollupIDToRollupData(opts *bind.CallOpts, rollupID uint32) (struct {
		RollupContract                 common.Address
		ChainID                        uint64
		Verifier                       common.Address
		ForkID                         uint64
		LastLocalExitRoot              [32]byte
		LastBatchSequenced             uint64
		LastVerifiedBatch              uint64
		LastPendingState               uint64
		LastPendingStateConsolidated   uint64
		LastVerifiedBatchBeforeUpgrade uint64
		RollupTypeID                   uint64
		RollupCompatibilityID          uint8
	}, error)
	RollupTypeMap(opts *bind.CallOpts, rollupTypeID uint32) (struct {
		ConsensusImplementation common.Address
		Verifier                common.Address
		ForkID                  uint64
		RollupCompatibilityID   uint8
		Obsolete                bool
		Genesis                 [32]byte
	}, error)
}

type rollupContractInterface interface {
	Admin(opts *bind.CallOpts) (common.Address, error)
	BridgeAddress(opts *bind.CallOpts) (common.Address, error)
	CalculatePolPerForceBatch(opts *bind.CallOpts) (*big.Int, error)
	ForceBatchAddress(opts *bind.CallOpts) (common.Address, error)
	ForceBatchTimeout(opts *bind.CallOpts) (uint64, error)
	ForcedBatches(opts *bind.CallOpts, arg0 uint64) ([32]byte, error)
	GasTokenAddress(opts *bind.CallOpts) (common.Address, error)
	GasTokenNetwork(opts *bind.CallOpts) (uint32, error)
	GenerateInitializeTransaction(opts *bind.CallOpts, networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error)
	GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error)
	LastAccInputHash(opts *bind.CallOpts) ([32]byte, error)
	LastForceBatch(opts *bind.CallOpts) (uint64, error)
	LastForceBatchSequenced(opts *bind.CallOpts) (uint64, error)
	NetworkName(opts *bind.CallOpts) (string, error)
	PendingAdmin(opts *bind.CallOpts) (common.Address, error)
	Pol(opts *bind.CallOpts) (common.Address, error)
	RollupManager(opts *bind.CallOpts) (common.Address, error)
	TrustedSequencer(opts *bind.CallOpts) (common.Address, error)
	TrustedSequencerURL(opts *bind.CallOpts) (string, error)
}

type validiumContractInterface interface {
	DataAvailabilityProtocol(opts *bind.CallOpts) (common.Address, error)
	IsSequenceWithDataAvailabilityAllowed(opts *bind.CallOpts) (bool, error)
}

type committeeContractInterface interface {
	CommitteeHash(opts *bind.CallOpts) ([32]byte, error)
	GetAmountOfMembers(opts *bind.CallOpts) (*big.Int, error)
	GetProtocolName(opts *bind.CallOpts) (string, error)
	Members(opts *bind.CallOpts, arg0 *big.Int) (struct {
		Url  string
		Addr common.Address
	}, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
	RequiredAmountOfSignatures(opts *bind.CallOpts) (*big.Int, error)
}

type bridgeContractInterface interface {
	WETHToken(opts *bind.CallOpts) (common.Address, error)
	DepositCount(opts *bind.CallOpts) (*big.Int, error)
	GasTokenAddress(opts *bind.CallOpts) (common.Address, error)
	GasTokenMetadata(opts *bind.CallOpts) ([]byte, error)
	GasTokenNetwork(opts *bind.CallOpts) (uint32, error)
	GetRoot(opts *bind.CallOpts) ([32]byte, error)
	GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error)
	IsEmergencyState(opts *bind.CallOpts) (bool, error)
	LastUpdatedDepositCount(opts *bind.CallOpts) (uint32, error)
	NetworkID(opts *bind.CallOpts) (uint32, error)
	PolygonRollupManager(opts *bind.CallOpts) (common.Address, error)
}

type gerContractInterface interface {
	BridgeAddress(opts *bind.CallOpts) (common.Address, error)
	DepositCount(opts *bind.CallOpts) (*big.Int, error)
	GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error)
	GetRoot(opts *bind.CallOpts) ([32]byte, error)
	GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)
	LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error)
	LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error)
	RollupManager(opts *bind.CallOpts) (common.Address, error)
}
