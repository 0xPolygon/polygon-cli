package cdk

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/0xPolygon/polygon-cli/custom_marshaller"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	banana_rollup "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonrollupbaseetrog"
	banana_rollup_manager "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonrollupmanager"
	banana_bridge "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmbridgev2"
	banana_ger "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmglobalexitrootv2"

	elderberry_rollup "github.com/0xPolygon/cdk-contracts-tooling/contracts/elderberry/polygonrollupbaseetrog"
	elderberry_rollup_manager "github.com/0xPolygon/cdk-contracts-tooling/contracts/elderberry/polygonrollupmanager"
	elderberry_bridge "github.com/0xPolygon/cdk-contracts-tooling/contracts/elderberry/polygonzkevmbridgev2"
	elderberry_ger "github.com/0xPolygon/cdk-contracts-tooling/contracts/elderberry/polygonzkevmglobalexitrootv2"

	etrog_rollup "github.com/0xPolygon/cdk-contracts-tooling/contracts/etrog/polygonrollupbaseetrog"
	etrog_rollup_manager "github.com/0xPolygon/cdk-contracts-tooling/contracts/etrog/polygonrollupmanager"
	etrog_bridge "github.com/0xPolygon/cdk-contracts-tooling/contracts/etrog/polygonzkevmbridgev2"
	etrog_ger "github.com/0xPolygon/cdk-contracts-tooling/contracts/etrog/polygonzkevmglobalexitrootv2"
)

const (
	ArgRpcURL = "rpc-url"
	ArgForkID = "fork-id"

	ArgRollupManagerAddress = "rollup-manager-address"

	ArgRollupChainID = "rollup-chain-id"
	ArgRollupID      = "rollup-id"
	ArgRollupAddress = "rollup-address"
	ArgBridgeAddress = "bridge-address"
	ArgGERAddress    = "ger-address"

	defaultRPCURL = "http://localhost:8545"
	defaultForkId = "12"

	// forks
	blueberry   = uint64(4)
	dragonfruit = uint64(5)
	incaberry   = uint64(6)
	etrog       = uint64(7)
	elderberry  = uint64(9)
	feijoa      = uint64(10)
	banana      = uint64(12)
	durian      = uint64(13)

	contractRequestInterval = 200 * time.Millisecond
)

var (
	knownRollupManagerAddresses = map[string]string{
		"bali":    "0xe2ef6215adc132df6913c8dd16487abf118d1764",
		"cardona": "0x32d33D5137a7cFFb54c5Bf8371172bcEc5f310ff",
		"mainnet": "0x5132a183e9f3cb7c848b0aac5ae0c4f0491b7ab2",
	}
	knownForks = map[string]uint64{
		// "4":           blueberry,
		// "blueberry":   blueberry,
		// "5":           dragonfruit,
		// "dragonfruit": dragonfruit,
		// "6":           incaberry,
		// "incaberry":   incaberry,
		"7":          etrog,
		"etrog":      etrog,
		"9":          elderberry,
		"elderberry": elderberry,
		// "10":          feijoa,
		// "feijoa":      feijoa,
		"12":     banana,
		"banana": banana,
		// "13":          durian,
		// "durian":      durian,
	}
)

var CDKCmd = &cobra.Command{
	Use:   "cdk",
	Short: "Utilities for interacting with CDK networks",
	Long:  "Basic utility commands for interacting with the cdk contracts",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cdkInputArgs.rpcURL = flag_loader.GetRpcUrlFlagValue(cmd)
	},
	Args: cobra.NoArgs,
}

type inputArgs struct {
	rpcURL *string

	forkID *string

	rollupManagerAddress *string

	rollupID      *string
	rollupChainID *string
	rollupAddress *string
	bridgeAddress *string
	gerAddress    *string
}

type parsedCDKArgs struct {
	rpcURL string
	forkID uint64
}

type parsedRollupManagerArgs struct {
	rollupManagerAddress common.Address
}

type parsedRollupArgs struct {
	rollupID      uint32
	rollupChainID uint64
	rollupAddress common.Address
}

var cdkInputArgs = inputArgs{}

func checkAddressArg(argFlagName, address string) error {
	prefix := fmt.Sprintf("invalid flag %s: ", argFlagName)

	if !common.IsHexAddress(address) {
		return errors.New(prefix + "invalid address")
	}

	return nil
}

func (inputArgs *inputArgs) parseCDKArgs(ctx context.Context) (parsedCDKArgs, error) {
	args := parsedCDKArgs{}

	args.rpcURL = *inputArgs.rpcURL

	if inputArgs.forkID != nil && len(*inputArgs.forkID) > 0 {
		_, found := knownForks[*inputArgs.forkID]
		if !found {
			return parsedCDKArgs{}, invalidForkIDErr()
		}
		args.forkID = knownForks[*inputArgs.forkID]
	}

	return args, nil
}

func (inputArgs *inputArgs) parseRollupManagerArgs(ctx context.Context, cdkArgs parsedCDKArgs) (*parsedRollupManagerArgs, error) {
	args := &parsedRollupManagerArgs{}

	if knownRollupManagerAddress, found := knownRollupManagerAddresses[*cdkInputArgs.rollupManagerAddress]; found {
		args.rollupManagerAddress = common.HexToAddress(knownRollupManagerAddress)
	} else {
		err := checkAddressArg(ArgRollupManagerAddress, *inputArgs.rollupManagerAddress)
		if err != nil {
			return nil, err
		}
		args.rollupManagerAddress = common.HexToAddress(*cdkInputArgs.rollupManagerAddress)
	}

	return args, nil
}

func (inputArgs *inputArgs) parseRollupArgs(ctx context.Context, rollupManager rollupManagerContractInterface) (*parsedRollupArgs, error) {
	args := &parsedRollupArgs{}

	var rollupChainID uint64
	if cdkInputArgs.rollupChainID != nil && len(*cdkInputArgs.rollupChainID) > 0 {
		rollupChainIDN, err := strconv.ParseInt(*cdkInputArgs.rollupChainID, 10, 64)
		if err != nil || rollupChainIDN < 0 {
			return nil, fmt.Errorf("invalid rollupChainID: %s, it must be a valid uint64", *cdkInputArgs.rollupID)
		}
		rollupChainID = uint64(rollupChainIDN)
	}
	args.rollupChainID = rollupChainID

	args.rollupAddress = common.Address{}
	if inputArgs.rollupAddress != nil && len(*inputArgs.rollupAddress) > 0 {
		err := checkAddressArg(ArgRollupAddress, *inputArgs.rollupAddress)
		if err != nil {
			return nil, err
		}
		args.rollupAddress = common.HexToAddress(*inputArgs.rollupAddress)
	}

	args.rollupID = 0
	if cdkInputArgs.rollupID != nil && len(*cdkInputArgs.rollupID) > 0 {
		rollupIDN, err := strconv.Atoi(*cdkInputArgs.rollupID)
		if err != nil || rollupIDN < 0 {
			return nil, fmt.Errorf("invalid rollupID: %s, it must be a valid uint32", *cdkInputArgs.rollupID)
		}
		args.rollupID = uint32(rollupIDN)
	} else {
		if rollupChainID == 0 && args.rollupAddress == (common.Address{}) {
			return nil, fmt.Errorf("%s, %s, or %s must be provided", ArgRollupID, ArgRollupChainID, ArgRollupAddress)
		}
		if rollupChainID != 0 {
			rollupID, err := rollupManager.ChainIDToRollupID(nil, rollupChainID)
			if err != nil {
				return nil, err
			}
			args.rollupID = rollupID
		} else if args.rollupAddress != (common.Address{}) {
			rollupID, err := rollupManager.RollupAddressToID(nil, args.rollupAddress)
			if err != nil {
				return nil, err
			}
			args.rollupID = rollupID
		}
	}

	return args, nil
}

func mustGetRPCClient(ctx context.Context, rpcURL string) *ethclient.Client {
	rpcClient, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to Dial RPC")
	}
	return rpcClient
}

func getRollupManager(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, addr common.Address) (*rollupManager, *abi.ABI, error) {
	var contract *rollupManager
	var contractABI *abi.ABI
	switch cdkArgs.forkID {
	case etrog:
		contractInstance, err := etrog_rollup_manager.NewPolygonrollupmanager(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &rollupManager{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = etrog_rollup_manager.PolygonrollupmanagerMetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case elderberry:
		contractInstance, err := elderberry_rollup_manager.NewPolygonrollupmanager(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &rollupManager{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = elderberry_rollup_manager.PolygonrollupmanagerMetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case banana:
		contractInstance, err := banana_rollup_manager.NewPolygonrollupmanager(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &rollupManager{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = banana_rollup_manager.PolygonrollupmanagerMetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, invalidForkIDErr()
	}
	return contract, contractABI, nil
}

func getRollup(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, addr common.Address) (rollupContractInterface, error) {
	var rollup rollupContractInterface
	var err error
	switch cdkArgs.forkID {
	case etrog:
		rollup, err = etrog_rollup.NewPolygonrollupbaseetrog(addr, rpcClient)
		if err != nil {
			return nil, err
		}
	case elderberry:
		rollup, err = elderberry_rollup.NewPolygonrollupbaseetrog(addr, rpcClient)
		if err != nil {
			return nil, err
		}
	case banana:
		rollup, err = banana_rollup.NewPolygonrollupbaseetrog(addr, rpcClient)
		if err != nil {
			return nil, err
		}
	default:
		return nil, invalidForkIDErr()
	}
	return rollup, nil
}

func getBridge(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, addr common.Address) (*bridge, *abi.ABI, error) {
	var contract *bridge
	var contractABI *abi.ABI
	switch cdkArgs.forkID {
	case etrog:
		contractInstance, err := etrog_bridge.NewPolygonzkevmbridgev2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &bridge{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = etrog_bridge.Polygonzkevmbridgev2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case elderberry:
		contractInstance, err := elderberry_bridge.NewPolygonzkevmbridgev2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &bridge{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = elderberry_bridge.Polygonzkevmbridgev2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case banana:
		contractInstance, err := banana_bridge.NewPolygonzkevmbridgev2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &bridge{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = banana_bridge.Polygonzkevmbridgev2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, invalidForkIDErr()
	}
	return contract, contractABI, nil
}

func getGER(cdkArgs parsedCDKArgs, rpcClient *ethclient.Client, addr common.Address) (*ger, *abi.ABI, error) {
	var contract *ger
	var contractABI *abi.ABI
	switch cdkArgs.forkID {
	case etrog:
		contractInstance, err := etrog_ger.NewPolygonzkevmglobalexitrootv2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &ger{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = etrog_ger.Polygonzkevmglobalexitrootv2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case elderberry:
		contractInstance, err := elderberry_ger.NewPolygonzkevmglobalexitrootv2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &ger{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = elderberry_ger.Polygonzkevmglobalexitrootv2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	case banana:
		contractInstance, err := banana_ger.NewPolygonzkevmglobalexitrootv2(addr, rpcClient)
		if err != nil {
			return nil, nil, err
		}
		contract = &ger{contractInstance, reflect.ValueOf(contractInstance)}
		contractABI, err = banana_ger.Polygonzkevmglobalexitrootv2MetaData.GetAbi()
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, invalidForkIDErr()
	}
	return contract, contractABI, nil
}

func mustPrintJSONIndent(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))
}

func invalidForkIDErr() error {
	forkIDs := make([]string, 0, len(knownForks))
	for forkID := range knownForks {
		forkIDs = append(forkIDs, forkID)
	}
	slices.Sort(forkIDs)
	v := strings.Join(forkIDs, ", ")
	return fmt.Errorf("invalid forkID. supported forkIDs are %s", v)
}

// watchNewLogs watches for new logs for the given filter and print them to the console
// - rpcClient is used to fetch the logs
// - filter is used to set which logs must be fetched
// - contractInstance and contractABI are used to parse the logs
// - logs are printed in JSON format
func watchNewLogs(ctx context.Context, rpcClient *ethclient.Client, filter ethereum.FilterQuery, contractInstance reflect.Value, contractABI *abi.ABI) error {
	log.Info().Msgf("Waiting for events")

	latestBlockNumber, err := rpcClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	time.Sleep(contractRequestInterval)

	// rewind 1 block to force reading the current block
	if latestBlockNumber > 0 {
		latestBlockNumber--
	}

	for {
		currentBlockNumber, err := rpcClient.BlockNumber(ctx)
		if err != nil {
			return err
		}
		time.Sleep(contractRequestInterval)

		// wait for the new block
		if currentBlockNumber <= latestBlockNumber {
			time.Sleep(time.Second)
		}

		for blockNumber := latestBlockNumber + 1; blockNumber <= currentBlockNumber; blockNumber++ {
			log.Info().Msgf("New block detected %d", blockNumber)

			filter.FromBlock = big.NewInt(0).SetUint64(blockNumber)
			filter.ToBlock = big.NewInt(0).SetUint64(blockNumber)

			logs, err := rpcClient.FilterLogs(ctx, filter)
			if err != nil {
				return err
			}
			time.Sleep(contractRequestInterval)

			mustPrintLogs(logs, contractInstance, contractABI)
		}
		latestBlockNumber = currentBlockNumber
	}
}

// mustPrintLogs prints the logs in JSON format
// - logs are parsed using the contractInstance and contractABI
// - logs are printed in JSON format
// - if the log cannot be parsed, the log is printed as is
func mustPrintLogs(logs []types.Log, contractInstance reflect.Value, contractABI *abi.ABI) {
	eventsFound := false
	for _, l := range logs {
		e, _ := contractABI.EventByID(l.Topics[0])
		if e == nil {
			mustPrintJSONIndent(l)
			continue
		}
		eventsFound = true

		var parsedEvent any
		parseLogMethod, methodFound := contractInstance.Type().MethodByName(fmt.Sprintf("Parse%s", e.Name))
		if !methodFound {
			log.Warn().Msgf("Method Parse%s not found", e.Name)
		} else {
			parsedLogValues := parseLogMethod.Func.Call([]reflect.Value{contractInstance, reflect.ValueOf(l)})
			parsedEventValue := parsedLogValues[0].Interface()
			errValue := parsedLogValues[1].Interface()
			if errValue != nil {
				log.Warn().Err(errValue.(error)).Msgf("Error parsing log %v", l)
			} else {
				parsedEvent = parsedEventValue
			}
		}

		customMarshaller := custom_marshaller.New(parsedEvent)

		mustPrintJSONIndent(struct {
			Name      string `json:"name"`
			Signature string `json:"signature"`
			Event     any    `json:"event"`
		}{
			Name:      e.Name,
			Signature: e.Sig,
			Event:     customMarshaller,
		})
	}
	if !eventsFound {
		log.Info().Msg("No events found")
	}
}

func init() {
	// cdk
	cdkInputArgs.rpcURL = CDKCmd.PersistentFlags().String(ArgRpcURL, defaultRPCURL, "The RPC URL of the network containing the CDK contracts")
	cdkInputArgs.forkID = CDKCmd.PersistentFlags().String(ArgForkID, defaultForkId, "The ForkID of the cdk networks")
	cdkInputArgs.rollupManagerAddress = CDKCmd.PersistentFlags().String(ArgRollupManagerAddress, "", "The address of the rollup contract")

	// rollup manager

	// rollup
	cdkInputArgs.rollupID = rollupCmd.PersistentFlags().String(ArgRollupID, "", "The rollup ID")
	cdkInputArgs.rollupChainID = rollupCmd.PersistentFlags().String(ArgRollupChainID, "", "The rollup chain ID")
	cdkInputArgs.rollupAddress = rollupCmd.PersistentFlags().String(ArgRollupAddress, "", "The rollup Address")

	// bridge
	cdkInputArgs.bridgeAddress = bridgeCmd.PersistentFlags().String(ArgBridgeAddress, "", "The address of the bridge contract")

	// ger
	cdkInputArgs.gerAddress = gerCmd.PersistentFlags().String(ArgGERAddress, "", "The address of the GER contract")

	CDKCmd.AddCommand(rollupManagerCmd)
	CDKCmd.AddCommand(rollupCmd)
	CDKCmd.AddCommand(bridgeCmd)
	CDKCmd.AddCommand(gerCmd)

	rollupManagerCmd.AddCommand(rollupManagerListRollupsCmd)
	rollupManagerCmd.AddCommand(rollupManagerListRollupTypesCmd)
	rollupManagerCmd.AddCommand(rollupManagerInspectCmd)
	rollupManagerCmd.AddCommand(rollupManagerDumpCmd)
	rollupManagerCmd.AddCommand(rollupManagerMonitorCmd)

	rollupCmd.AddCommand(rollupInspectCmd)
	rollupCmd.AddCommand(rollupDumpCmd)
	rollupCmd.AddCommand(rollupMonitorCmd)

	bridgeCmd.AddCommand(bridgeInspectCmd)
	bridgeCmd.AddCommand(bridgeDumpCmd)
	bridgeCmd.AddCommand(bridgeMonitorCmd)

	gerCmd.AddCommand(gerInspectCmd)
	gerCmd.AddCommand(gerDumpCmd)
	gerCmd.AddCommand(gerMonitorCmd)
}
