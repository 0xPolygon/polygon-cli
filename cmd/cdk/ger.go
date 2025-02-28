package cdk

import (
	_ "embed"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var gerCmd = &cobra.Command{
	Use:   "ger",
	Short: "Utilities for interacting with CDK global exit root manager contract",
	Args:  cobra.NoArgs,
}

//go:embed gerInspectUsage.md
var gerInspectUsage string
var gerInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "List some basic information about the global exit root manager",
	Long:  gerInspectUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return gerInspect(cmd)
	},
}

//go:embed gerDumpUsage.md
var gerDumpUsage string
var gerDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "List detailed information about the global exit root manager",
	Long:  gerDumpUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return gerDump(cmd)
	},
}

//go:embed gerMonitorUsage.md
var gerMonitorUsage string
var gerMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for global exit root manager events and display them on the fly",
	Long:  gerMonitorUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return gerMonitor(cmd)
	},
}

type gerData struct {
	BridgeAddress         common.Address `json:"bridgeAddress"`
	GetLastGlobalExitRoot common.Hash    `json:"getLastGlobalExitRoot"`
	LastMainnetExitRoot   common.Hash    `json:"lastMainnetExitRoot"`
	LastRollupExitRoot    common.Hash    `json:"lastRollupExitRoot"`
	// RollupAddress         common.Address `json:"rollupAddress"`
}

type gerDumpData struct {
	Data *gerData `json:"data"`
}

func gerInspect(cmd *cobra.Command) error {
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

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
	if err != nil {
		return err
	}

	bridgeData, err := getBridgeData(bridge)
	if err != nil {
		return err
	}

	ger, err := getGER(cdkArgs, rpcClient, bridgeData.GlobalExitRootManager)
	if err != nil {
		return err
	}

	data, err := getGERData(ger)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)
	return nil
}

func gerDump(cmd *cobra.Command) error {
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

	rollupManager, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
	if err != nil {
		return err
	}

	bridgeData, err := getBridgeData(bridge)
	if err != nil {
		return err
	}

	ger, err := getGER(cdkArgs, rpcClient, bridgeData.GlobalExitRootManager)
	if err != nil {
		return err
	}

	data := &gerDumpData{}

	data.Data, err = getGERData(ger)
	if err != nil {
		return err
	}

	mustLogJSONIndent(data)
	return nil
}

func gerMonitor(cmd *cobra.Command) error {
	panic("not implemented")
}

func getGERData(ger gerContractInterface) (*gerData, error) {
	data := &gerData{}
	var err error

	data.BridgeAddress, err = ger.BridgeAddress(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.GetLastGlobalExitRoot, err = ger.GetLastGlobalExitRoot(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastMainnetExitRoot, err = ger.LastMainnetExitRoot(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastRollupExitRoot, err = ger.LastRollupExitRoot(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	// data.RollupAddress, err = ger.RollupAddress(nil)
	// if err != nil {
	// 	return nil, err
	// }
	// time.Sleep(contractRequestInterval)

	return data, nil
}
