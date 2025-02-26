package cdk

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var rollupManagerCmd = &cobra.Command{
	Use:  "rollup-manager",
	Args: cobra.NoArgs,
}

var rollupManagerListRollupsCmd = &cobra.Command{
	Use:  "list-rollups",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerListRollups(cmd)
	},
}

var rollupManagerListRollupTypesCmd = &cobra.Command{
	Use:  "list-rollup-types",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerListRollupTypes(cmd)
	},
}

var rollupManagerInspectCmd = &cobra.Command{
	Use:  "inspect",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerInspect(cmd)
	},
}

var rollupManagerDumpCmd = &cobra.Command{
	Use:  "dump",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerDump(cmd)
	},
}

var rollupManagerMonitorCmd = &cobra.Command{
	Use:  "monitor",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerMonitor(cmd)
	},
}

type RollupData struct {
	RollupContract                 common.Address `json:"rollupContract"`
	ChainID                        uint64         `json:"chainID"`
	Verifier                       common.Address `json:"verifier"`
	ForkID                         uint64         `json:"forkID"`
	LastLocalExitRoot              common.Hash    `json:"lastLocalExitRoot"`
	LastBatchSequenced             uint64         `json:"lastBatchSequenced"`
	LastVerifiedBatch              uint64         `json:"lastVerifiedBatch"`
	LastPendingState               uint64         `json:"lastPendingState"`
	LastPendingStateConsolidated   uint64         `json:"lastPendingStateConsolidated"`
	LastVerifiedBatchBeforeUpgrade uint64         `json:"lastVerifiedBatchBeforeUpgrade"`
	RollupTypeID                   uint64         `json:"rollupTypeID"`
	RollupCompatibilityID          uint8          `json:"rollupCompatibilityID"`
}

type RollupTypeData struct {
	ConsensusImplementation common.Address `json:"consensusImplementation"`
	Verifier                common.Address `json:"verifier"`
	ForkID                  uint64         `json:"forkID"`
	RollupCompatibilityID   uint8          `json:"rollupCompatibilityID"`
	Obsolete                bool           `json:"obsolete"`
	Genesis                 common.Hash    `json:"genesis"`
}

type RollupManagerInfo struct {
	Pol                                    common.Address `json:"pol"`
	BridgeAddress                          common.Address `json:"bridgeAddress"`
	RollupCount                            uint32         `json:"rollupCount"`
	BatchFee                               *big.Int       `json:"batchFee"`
	TotalSequencedBatches                  uint64         `json:"totalSequencedBatches"`
	TotalVerifiedBatches                   uint64         `json:"totalVerifiedBatches"`
	LastAggregationTimestamp               uint64         `json:"lastAggregationTimestamp"`
	LastDeactivatedEmergencyStateTimestamp uint64         `json:"lastDeactivatedEmergencyStateTimestamp"`
	TrustedAggregatorTimeout               uint64         `json:"trustedAggregatorTimeout"`
	PendingStateTimeout                    uint64         `json:"pendingStateTimeout"`
	MultiplierBatchFee                     uint16         `json:"multiplierBatchFee"`
}

type RollupManagerDumpData struct {
	Info        *RollupManagerInfo `json:"info"`
	Rollups     []RollupData       `json:"rollups"`
	RollupTypes []RollupTypeData   `json:"rollupTypes"`
}

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

func rollupManagerListRollups(cmd *cobra.Command) error {
	cdkArgs, err := cdkInputArgs.parseCDKArgs(cmd.Context())
	if err != nil {
		return err
	}

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(cmd.Context(), *cdkArgs)
	if err != nil {
		return err
	}

	rollups, err := getRollupManagerRollups(rollupManagerArgs)
	if err != nil {
		return err
	}

	mustLogJSONIndent(rollups)
	return nil
}

func rollupManagerListRollupTypes(cmd *cobra.Command) error {
	cdkArgs, err := cdkInputArgs.parseCDKArgs(cmd.Context())
	if err != nil {
		return err
	}

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(cmd.Context(), *cdkArgs)
	if err != nil {
		return err
	}

	rollupTypes, err := getRollupManagerRollupTypes(rollupManagerArgs)
	if err != nil {
		return err
	}

	mustLogJSONIndent(rollupTypes)
	return nil
}

func rollupManagerInspect(cmd *cobra.Command) error {
	cdkArgs, err := cdkInputArgs.parseCDKArgs(cmd.Context())
	if err != nil {
		return err
	}

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(cmd.Context(), *cdkArgs)
	if err != nil {
		return err
	}

	data, err := getRollupManagerInfo(rollupManagerArgs)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)
	return nil
}

func rollupManagerDump(cmd *cobra.Command) error {
	cdkArgs, err := cdkInputArgs.parseCDKArgs(cmd.Context())
	if err != nil {
		return err
	}

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(cmd.Context(), *cdkArgs)
	if err != nil {
		return err
	}

	data := &RollupManagerDumpData{}

	data.Info, err = getRollupManagerInfo(rollupManagerArgs)
	if err != nil {
		return err
	}

	data.Rollups, err = getRollupManagerRollups(rollupManagerArgs)
	if err != nil {
		return err
	}

	data.RollupTypes, err = getRollupManagerRollupTypes(rollupManagerArgs)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)

	return nil
}

func rollupManagerMonitor(cmd *cobra.Command) error {
	panic("not implemented")
}

func getRollupManagerRollups(rollupManagerArgs *parsedRollupManagerArgs) ([]RollupData, error) {
	rollupCount, err := rollupManagerArgs.rollupManager.RollupCount(nil)
	if err != nil {
		return nil, err
	}

	rollups := make([]RollupData, 0, rollupCount)
	for i := uint32(1); i <= rollupCount; i++ {
		rollupData, err := rollupManagerArgs.rollupManager.RollupIDToRollupData(nil, i)
		if err != nil {
			return nil, err
		}
		rollups = append(rollups, RollupData{
			RollupContract:                 rollupData.RollupContract,
			ChainID:                        rollupData.ChainID,
			Verifier:                       rollupData.Verifier,
			ForkID:                         rollupData.ForkID,
			LastLocalExitRoot:              rollupData.LastLocalExitRoot,
			LastBatchSequenced:             rollupData.LastBatchSequenced,
			LastVerifiedBatch:              rollupData.LastVerifiedBatch,
			LastPendingState:               rollupData.LastPendingState,
			LastPendingStateConsolidated:   rollupData.LastPendingStateConsolidated,
			LastVerifiedBatchBeforeUpgrade: rollupData.LastVerifiedBatchBeforeUpgrade,
			RollupTypeID:                   rollupData.RollupTypeID,
			RollupCompatibilityID:          rollupData.RollupCompatibilityID,
		})
		time.Sleep(contractRequestInterval)
	}
	return rollups, nil
}

func getRollupManagerRollupTypes(rollupManagerArgs *parsedRollupManagerArgs) ([]RollupTypeData, error) {
	rollupTypeCount, err := rollupManagerArgs.rollupManager.RollupTypeCount(nil)
	if err != nil {
		return nil, err
	}

	rollupTypes := make([]RollupTypeData, 0, rollupTypeCount)
	for i := uint32(1); i <= rollupTypeCount; i++ {
		rollupType, err := rollupManagerArgs.rollupManager.RollupTypeMap(nil, i)
		if err != nil {
			return nil, err
		}
		rollupTypes = append(rollupTypes, RollupTypeData{
			ConsensusImplementation: rollupType.ConsensusImplementation,
			Verifier:                rollupType.Verifier,
			ForkID:                  rollupType.ForkID,
			RollupCompatibilityID:   rollupType.RollupCompatibilityID,
			Obsolete:                rollupType.Obsolete,
			Genesis:                 rollupType.Genesis,
		})
		time.Sleep(contractRequestInterval)
	}
	return rollupTypes, nil
}

func getRollupManagerInfo(rollupManagerArgs *parsedRollupManagerArgs) (*RollupManagerInfo, error) {
	data := &RollupManagerInfo{}
	var err error

	data.Pol, err = rollupManagerArgs.rollupManager.Pol(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.BridgeAddress, err = rollupManagerArgs.rollupManager.BridgeAddress(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.RollupCount, err = rollupManagerArgs.rollupManager.RollupCount(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.BatchFee, err = rollupManagerArgs.rollupManager.GetBatchFee(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.TotalSequencedBatches, err = rollupManagerArgs.rollupManager.TotalSequencedBatches(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.TotalVerifiedBatches, err = rollupManagerArgs.rollupManager.TotalVerifiedBatches(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastAggregationTimestamp, err = rollupManagerArgs.rollupManager.LastAggregationTimestamp(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastDeactivatedEmergencyStateTimestamp, err = rollupManagerArgs.rollupManager.LastDeactivatedEmergencyStateTimestamp(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	// data.TrustedAggregatorTimeout, err = rollupManagerArgs.rollupManager.TrustedAggregatorTimeout(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	// data.PendingStateTimeout, err = rollupManagerArgs.rollupManager.PendingStateTimeout(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	// data.MultiplierBatchFee, err = rollupManagerArgs.rollupManager.MultiplierBatchFee(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	return data, nil
}
