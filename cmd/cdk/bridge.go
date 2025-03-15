package cdk

import (
	_ "embed"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum"
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
	WETHToken               common.Address `json:"wethToken"`
	DepositCount            *big.Int       `json:"depositCount"`
	GasTokenAddress         common.Address `json:"gasTokenAddress"`
	GasTokenMetadata        common.Hash    `json:"gasTokenMetadata"`
	GasTokenNetwork         uint32         `json:"gasTokenNetwork"`
	GetRoot                 common.Hash    `json:"getRoot"`
	GlobalExitRootManager   common.Address `json:"globalExitRootManager"`
	IsEmergencyState        bool           `json:"isEmergencyState"`
	LastUpdatedDepositCount uint32         `json:"lastUpdatedDepositCount"`
	NetworkID               uint32         `json:"networkID"`
	PolygonRollupManager    common.Address `json:"polygonRollupManager"`
}

type bridge struct {
	bridgeContractInterface
	instance reflect.Value
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, _, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
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

	rollupManager, _, err := getRollupManager(cdkArgs, rpcClient, rollupManagerArgs.rollupManagerAddress)
	if err != nil {
		return err
	}

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, _, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
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

	rollupManagerData, err := getRollupManagerData(rollupManager)
	if err != nil {
		return err
	}

	bridge, bridgeABI, err := getBridge(cdkArgs, rpcClient, rollupManagerData.BridgeAddress)
	if err != nil {
		return err
	}

	filter := customFilter{
		contractInstance: bridge.instance,
		contractABI:      bridgeABI,
		blockchainFilter: ethereum.FilterQuery{
			Addresses: []common.Address{rollupManagerData.BridgeAddress},
		},
	}

	err = watchNewLogs(ctx, rpcClient, filter)
	if err != nil {
		return err
	}

	return nil
}

func getBridgeData(bridge bridgeContractInterface) (*BridgeData, error) {
	data := &BridgeData{}
	var err error

	data.WETHToken, err = bridge.WETHToken(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.DepositCount, err = bridge.DepositCount(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.GasTokenAddress, err = bridge.GasTokenAddress(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	gasTokenMetadata, err := bridge.GasTokenMetadata(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)
	data.GasTokenMetadata = common.BytesToHash(gasTokenMetadata)

	data.GasTokenNetwork, err = bridge.GasTokenNetwork(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.GetRoot, err = bridge.GetRoot(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.GlobalExitRootManager, err = bridge.GlobalExitRootManager(nil)
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

	data.NetworkID, err = bridge.NetworkID(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	data.PolygonRollupManager, err = bridge.PolygonRollupManager(nil)
	if err != nil {
		return nil, err
	}
	time.Sleep(contractRequestInterval)

	return data, nil
}
