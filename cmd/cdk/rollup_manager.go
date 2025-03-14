package cdk

import (
	_ "embed"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var rollupManagerCmd = &cobra.Command{
	Use:   "rollup-manager",
	Short: "Utilities for interacting with CDK rollup manager contract",
	Args:  cobra.NoArgs,
}

//go:embed rollupManagerListRollupsUsage.md
var rollupManagerListRollupsUsage string

var rollupManagerListRollupsCmd = &cobra.Command{
	Use:   "list-rollups",
	Short: "List some basic information about each rollup",
	Long:  rollupManagerListRollupsUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerListRollups(cmd)
	},
}

//go:embed rollupManagerListRollupTypesUsage.md
var rollupManagerListRollupTypesUsage string

var rollupManagerListRollupTypesCmd = &cobra.Command{
	Use:   "list-rollup-types",
	Short: "List some basic information about each rollup type",
	Long:  rollupManagerListRollupTypesUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerListRollupTypes(cmd)
	},
}

//go:embed rollupManagerInspectUsage.md
var rollupManagerInspectUsage string

var rollupManagerInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "List some basic information about the rollup manager",
	Long:  rollupManagerInspectUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerInspect(cmd)
	},
}

//go:embed rollupManagerDumpUsage.md
var rollupManagerDumpUsage string

var rollupManagerDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "List detailed information about the rollup manager",
	Long:  rollupManagerDumpUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerDump(cmd)
	},
}

//go:embed rollupManagerMonitorUsage.md
var rollupManagerMonitorUsage string

var rollupManagerMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for rollup manager events and display them on the fly",
	Long:  rollupManagerMonitorUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rollupManagerMonitor(cmd)
	},
}

type rollupManager struct {
	rollupManagerContractInterface
	instance reflect.Value
}

type RollupManagerData struct {
	Pol                                    common.Address `json:"pol"`
	BridgeAddress                          common.Address `json:"bridgeAddress"`
	RollupCount                            uint32         `json:"rollupCount"`
	BatchFee                               *big.Int       `json:"batchFee"`
	TotalSequencedBatches                  uint64         `json:"totalSequencedBatches"`
	TotalVerifiedBatches                   uint64         `json:"totalVerifiedBatches"`
	LastAggregationTimestamp               uint64         `json:"lastAggregationTimestamp"`
	LastDeactivatedEmergencyStateTimestamp uint64         `json:"lastDeactivatedEmergencyStateTimestamp"`
	// TrustedAggregatorTimeout               uint64         `json:"trustedAggregatorTimeout"`
	// PendingStateTimeout                    uint64         `json:"pendingStateTimeout"`
	// MultiplierBatchFee                     uint16         `json:"multiplierBatchFee"`
}

type RollupManagerDumpData struct {
	Data        *RollupManagerData `json:"data"`
	Rollups     []RollupData       `json:"rollups"`
	RollupTypes []RollupTypeData   `json:"rollupTypes"`
}

func rollupManagerListRollups(cmd *cobra.Command) error {
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollups, err := getRollupManagerRollups(cdkArgs, rpcClient, rollupManager)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(rollups)
	return nil
}

func rollupManagerListRollupTypes(cmd *cobra.Command) error {
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupTypes, err := getRollupManagerRollupTypes(rollupManager)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(rollupTypes)
	return nil
}

func rollupManagerInspect(cmd *cobra.Command) error {
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	data, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)
	return nil
}

func rollupManagerDump(cmd *cobra.Command) error {
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	data := &RollupManagerDumpData{}

	data.Data, err = getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	data.Rollups, err = getRollupManagerRollups(cdkArgs, rpcClient, rollupManager)
	if err != nil {
		return err
	}

	data.RollupTypes, err = getRollupManagerRollupTypes(rollupManager)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)

	return nil
}

func rollupManagerMonitor(cmd *cobra.Command) error {
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

	rollupManager, rollupManagerABI, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	filter := ethereum.FilterQuery{
		Addresses: []common.Address{rollupManagerArgs.rollupManagerAddress},
	}

	err = watchNewLogs(ctx, rpcClient, filter, rollupManager.instance, rollupManagerABI)
	if err != nil {
		return err
	}

	return nil
}

func getRollupManagerRollups(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, rollupManager rollupManagerContractInterface) ([]RollupData, error) {
	rollupCount, err := rollupManager.RollupCount(nil)
	if err != nil {
		return nil, err
	}

	rollups := make([]RollupData, 0, rollupCount)
	for i := uint32(1); i <= rollupCount; i++ {
		rollupData, err := getRollupData(cdkArgs, rpcClient, rollupManager, i)
		if err != nil {
			return nil, err
		}
		rollups = append(rollups, *rollupData)
		time.Sleep(contractRequestInterval)
	}
	return rollups, nil
}

func getRollupManagerRollupTypes(rollupManager rollupManagerContractInterface) ([]RollupTypeData, error) {
	rollupTypeCount, err := rollupManager.RollupTypeCount(nil)
	if err != nil {
		return nil, err
	}

	rollupTypes := make([]RollupTypeData, 0, rollupTypeCount)
	for i := uint32(1); i <= rollupTypeCount; i++ {
		rollupTypeData, err := getRollupTypeData(rollupManager, uint64(i))
		if err != nil {
			return nil, err
		}
		rollupTypes = append(rollupTypes, *rollupTypeData)
		time.Sleep(contractRequestInterval)
	}
	return rollupTypes, nil
}

func getRollupManagerData(rollupManager rollupManagerContractInterface) (*RollupManagerData, error) {
	data := &RollupManagerData{}
	var err error

	data.Pol, err = rollupManager.Pol(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.BridgeAddress, err = rollupManager.BridgeAddress(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.RollupCount, err = rollupManager.RollupCount(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.BatchFee, err = rollupManager.GetBatchFee(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.TotalSequencedBatches, err = rollupManager.TotalSequencedBatches(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.TotalVerifiedBatches, err = rollupManager.TotalVerifiedBatches(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastAggregationTimestamp, err = rollupManager.LastAggregationTimestamp(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastDeactivatedEmergencyStateTimestamp, err = rollupManager.LastDeactivatedEmergencyStateTimestamp(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	// data.TrustedAggregatorTimeout, err = rollupManager.TrustedAggregatorTimeout(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	// data.PendingStateTimeout, err = rollupManager.PendingStateTimeout(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	// data.MultiplierBatchFee, err = rollupManager.MultiplierBatchFee(nil)
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(contractRequestInterval)

	return data, nil
}
