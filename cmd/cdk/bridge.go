package cdk

import (
	_ "embed"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Utilities for interacting with CDK bridge contract",
	Args:  cobra.NoArgs,
}

//go:embed bridgeInspectUsage.md
var bridgeInspectUsage string
var bridgeInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "List some basic information about the bridge",
	Long:  bridgeInspectUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bridgeInspect(cmd)
	},
}

//go:embed bridgeDumpUsage.md
var bridgeDumpUsage string
var bridgeDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "List detailed information about the bridge",
	Long:  bridgeDumpUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bridgeDump(cmd)
	},
}

//go:embed bridgeMonitorUsage.md
var bridgeMonitorUsage string
var bridgeMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Watch for bridge events and display them on the fly",
	Long:  bridgeMonitorUsage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bridgeMonitor(cmd)
	},
}

type BridgeData struct {
	NetworkID               uint32         `json:"networkID"`
	DepositCount            *big.Int       `json:"depositCount"`
	IsEmergencyState        bool           `json:"isEmergencyState"`
	LastUpdatedDepositCount uint32         `json:"lastUpdatedDepositCount"`
	GlobalExitRootManager   common.Address `json:"globalExitRootManager"`
	// GetDepositRoot          common.Hash    `json:"getDepositRoot"`
	// PolygonZkEVMaddress     common.Address `json:"polygonZkEVMaddress"`
}

type BridgeDumpData struct {
	Data *BridgeData `json:"data"`
}

func bridgeInspect(cmd *cobra.Command) error {
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

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
	if err != nil {
		return err
	}

	data, err := getBridgeData(bridge)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)
	return nil
}

func bridgeDump(cmd *cobra.Command) error {
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

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
	if err != nil {
		return err
	}

	data := &BridgeDumpData{}

	data.Data, err = getBridgeData(bridge)
	if err != nil {
		return err
	}

	mustPrintJSONIndent(data)
	return nil
}

func bridgeMonitor(cmd *cobra.Command) error {
	panic("not implemented")
}

func getBridgeData(bridge bridgeContractInterface) (*BridgeData, error) {
	data := &BridgeData{}
	var err error

	data.NetworkID, err = bridge.NetworkID(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.DepositCount, err = bridge.DepositCount(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.IsEmergencyState, err = bridge.IsEmergencyState(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.LastUpdatedDepositCount, err = bridge.LastUpdatedDepositCount(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.GlobalExitRootManager, err = bridge.GlobalExitRootManager(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	// data.GetDepositRoot, err = bridge.GetDepositRoot(nil)
	// if err != nil {
	// 	return nil, err
	// }
	// time.Sleep(contractRequestInterval)

	// data.PolygonZkEVMaddress, err = bridge.PolygonZkEVMaddress(nil)
	// if err != nil {
	// 	return nil, err
	// }
	// time.Sleep(contractRequestInterval)

	return data, nil
}
