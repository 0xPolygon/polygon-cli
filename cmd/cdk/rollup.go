package cdk

import (
	_ "embed"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var rollupCmd = &cobra.Command{
	Use:   "rollup",
	Short: "Utilities for interacting with CDK rollup manager to get rollup specific information",
	Args:  cobra.NoArgs,
}

//go:embed rollupInspectUsage.md
var rollupInspectUsage string

var rollupInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "List some basic information about a specific rollup",
	Long:  rollupInspectUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupInspect(cmd)
	},
}

//go:embed rollupDumpUsage.md
var rollupDumpUsage string

var rollupDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "List detailed information about a specific rollup",
	Long:  rollupDumpUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupDump(cmd)
	},
}

//go:embed rollupMonitorUsage.md
var rollupMonitorUsage string

var rollupMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for rollup events and display them on the fly",
	Long:  rollupMonitorUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupMonitor(cmd)
	},
}

type RollupData struct {
	// from rollup manager sc
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

	// from rollup sc
	Admin               common.Address `json:"admin"`
	GasTokenAddress     common.Address `json:"gasTokenAddress"`
	GasTokenNetwork     uint32         `json:"gasTokenNetwork"`
	LastAccInputHash    common.Hash    `json:"lastAccInputHash"`
	NetworkName         string         `json:"networkName"`
	TrustedSequencer    common.Address `json:"trustedSequencer"`
	TrustedSequencerURL string         `json:"trustedSequencerURL"`
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

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(ctx, cdkArgs)
	if err != nil {
		return err
	}

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupArgs, err := cdkInputArgs.parseRollupArgs(ctx, rollupManager)
	if err != nil {
		return err
	}

	data, err := getRollupData(cdkArgs, rpcClient, rollupManager, rollupArgs.rollupID)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)
	return nil
}

func rollupDump(cmd *cobra.Command) error {
	ctx := cmd.Context()

	cdkArgs, err := cdkInputArgs.parseCDKArgs(ctx)
	if err != nil {
		return err
	}

	rpcClient := mustGetRPCClient(ctx, cdkArgs.rpcURL)

	rollupManagerArgs, err := cdkInputArgs.parseRollupManagerArgs(ctx, cdkArgs)
	if err != nil {
		return err
	}

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupArgs, err := cdkInputArgs.parseRollupArgs(ctx, rollupManager)
	if err != nil {
		return err
	}

	data := &RollupDumpData{}

	data.Data, err = getRollupData(cdkArgs, rpcClient, rollupManager, rollupArgs.rollupID)
	if err != nil {
		return err
	}

	data.Type, err = getRollupTypeData(rollupManager, data.Data.RollupTypeID)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)

	return nil
}

func rollupMonitor(cmd *cobra.Command) error {
	panic("not implemented")
}

func getRollupData(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, rollupManager rollupManagerContractInterface, rollupID uint32) (*RollupData, error) {
	rollupData, err := rollupManager.RollupIDToRollupData(nil, rollupID)
	if err != nil {
		return nil, err
	}

	rollup, err := getRollup(cdkArgs, rpcClient, rollupData.RollupContract)
	if err != nil {
		return nil, err
	}

	admin, err := rollup.Admin(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	gasTokenAddress, err := rollup.GasTokenAddress(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	gasTokenNetwork, err := rollup.GasTokenNetwork(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	lastAccInputHash, err := rollup.LastAccInputHash(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	networkName, err := rollup.NetworkName(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	trustedSequencer, err := rollup.TrustedSequencer(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	trustedSequencerURL, err := rollup.TrustedSequencerURL(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

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

		Admin:               admin,
		GasTokenAddress:     gasTokenAddress,
		GasTokenNetwork:     gasTokenNetwork,
		LastAccInputHash:    lastAccInputHash,
		NetworkName:         networkName,
		TrustedSequencer:    trustedSequencer,
		TrustedSequencerURL: trustedSequencerURL,
	}, nil
}

func getRollupTypeData(rollupManager rollupManagerContractInterface, rollupTypeID uint64) (*RollupTypeData, error) {
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
