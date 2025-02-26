package cdk

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var rollupCmd = &cobra.Command{
	Use:  "rollup",
	Args: cobra.NoArgs,
}

var rollupInspectCmd = &cobra.Command{
	Use:  "inspect",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupInspect(cmd)
	},
}

var rollupDumpCmd = &cobra.Command{
	Use:  "dump",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupDump(cmd)
	},
}

var rollupMonitorCmd = &cobra.Command{
	Use:  "monitor",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupMonitor(cmd)
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

type RollupDumpData struct {
	Data *RollupData     `json:"data"`
	Type *RollupTypeData `json:"type"`
}

func rollupInspect(cmd *cobra.Command) error {
	ctx := cmd.Context()

	cdkArgs, err := cdkInputArgs.parseCDKArgs(ctx)
	if err != nil {
		return err
	}

	rpcClient := mustGetRPCClient(ctx, cdkArgs.rpcURL)

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(ctx, *cdkArgs)
	if err != nil {
		return err
	}

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs)
	if err != nil {
		return err
	}

	rollupArgs, err := cdkInputArgs.parseRollupArgs(ctx, rollupManager)
	if err != nil {
		return err
	}

	data, err := getRollupData(rollupManager, rollupArgs.rollupID)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)
	return nil
}

func rollupDump(cmd *cobra.Command) error {
	ctx := cmd.Context()

	cdkArgs, err := cdkInputArgs.parseCDKArgs(ctx)
	if err != nil {
		return err
	}

	rpcClient := mustGetRPCClient(ctx, cdkArgs.rpcURL)

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(ctx, *cdkArgs)
	if err != nil {
		return err
	}

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs)
	if err != nil {
		return err
	}

	rollupArgs, err := cdkInputArgs.parseRollupArgs(ctx, rollupManager)
	if err != nil {
		return err
	}

	data := &RollupDumpData{}

	data.Data, err = getRollupData(rollupManager, rollupArgs.rollupID)
	if err != nil {
		return err
	}

	data.Type, err = getRollupType(rollupManager, data.Data.RollupTypeID)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)

	return nil
}

func rollupMonitor(cmd *cobra.Command) error {
	panic("not implemented")
}

func getRollupData(rollupManager rollupManagerContractInterface, rollupID uint32) (*RollupData, error) {
	rollupData, err := rollupManager.RollupIDToRollupData(nil, rollupID)
	if err != nil {
		return nil, err
	}

	return &RollupData{
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
	}, nil
}

func getRollupType(rollupManager rollupManagerContractInterface, rollupTypeID uint64) (*RollupTypeData, error) {
	rollupType, err := rollupManager.RollupTypeMap(nil, uint32(rollupTypeID))
	if err != nil {
		return nil, err
	}
	return &RollupTypeData{
		ConsensusImplementation: rollupType.ConsensusImplementation,
		Verifier:                rollupType.Verifier,
		ForkID:                  rollupType.ForkID,
		RollupCompatibilityID:   rollupType.RollupCompatibilityID,
		Obsolete:                rollupType.Obsolete,
		Genesis:                 rollupType.Genesis,
	}, nil
}
