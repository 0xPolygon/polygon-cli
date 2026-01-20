package common

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	bridge_service_factory "github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/factory"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgGasLimit             = "gas-limit"
	ArgChainID              = "chain-id"
	ArgPrivateKey           = flag.PrivateKey
	ArgValue                = "value"
	ArgRPCURL               = flag.RPCURL
	ArgBridgeAddress        = "bridge-address"
	ArgRollupManagerAddress = "rollup-manager-address"
	ArgDestNetwork          = "destination-network"
	ArgDestAddress          = "destination-address"
	ArgForceUpdate          = "force-update-root"
	ArgCallData             = "call-data"
	ArgCallDataFile         = "call-data-file"
	ArgTimeout              = "transaction-receipt-timeout"
	ArgDepositCount         = "deposit-count"
	ArgDepositNetwork       = "deposit-network"
	ArgRollupID             = "rollup-id"
	ArgCompleteMT           = "complete-merkle-tree"
	ArgBridgeServiceURL     = "bridge-service-url"
	ArgFileName             = "file-name"
	ArgL2ClaimsFileName     = "l2-claims-file"
	ArgL2DepositsFileName   = "l2-deposits-file"
	ArgL2NetworkID          = "l2-network-id"
	ArgFromBlock            = "from-block"
	ArgToBlock              = "to-block"
	ArgFilterSize           = "filter-size"
	ArgTokenAddress         = "token-address"
	ArgGlobalIndex          = "global-index"
	ArgDryRun               = "dry-run"
	ArgGasPrice             = "gas-price"
	ArgBridgeMappings       = "bridge-service-map"
	ArgBridgeLimit          = "bridge-limit"
	ArgBridgeOffset         = "bridge-offset"
	ArgWait                 = "wait"
	ArgConcurrency          = "concurrency"
	ArgInsecure             = "insecure"
	ArgLegacy               = "legacy"
	ArgProofGER             = "proof-ger"
	ArgProofL1InfoTreeIndex = "proof-l1-info-tree-index"
)

// UlxlyArgs holds the arguments for ulxly commands
type UlxlyArgs struct {
	GasLimit             uint64
	ChainID              string
	PrivateKey           string
	AddressOfPrivateKey  string
	Value                string
	RPCURL               string
	BridgeAddress        string
	DestNetwork          uint32
	DestAddress          string
	TokenAddress         string
	ForceUpdate          bool
	CallData             string
	CallDataFile         string
	Timeout              uint64
	DepositCount         uint32
	DepositNetwork       uint32
	BridgeServiceURL     string
	GlobalIndex          string
	GasPrice             string
	DryRun               bool
	BridgeServiceURLs    []string
	BridgeLimit          int
	BridgeOffset         int
	Wait                 time.Duration
	Concurrency          uint
	Insecure             bool
	Legacy               bool
	ProofGER             string
	ProofL1InfoTreeIndex uint32
}

// JSONError represents a JSON-RPC error
type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// GetEvent holds options for retrieving events
type GetEvent struct {
	URL                            string
	FromBlock, ToBlock, FilterSize uint64
	Insecure                       bool
}

// AddFlags adds event retrieval flags to a command
func (o *GetEvent) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&o.URL, ArgRPCURL, "u", "", "RPC URL to read the events data")
	f.Uint64VarP(&o.FromBlock, ArgFromBlock, "f", 0, "start of the range of blocks to retrieve")
	f.Uint64VarP(&o.ToBlock, ArgToBlock, "t", 0, "end of the range of blocks to retrieve")
	f.Uint64VarP(&o.FilterSize, ArgFilterSize, "i", 1000, "batch size for individual filter queries")
	f.BoolVarP(&o.Insecure, ArgInsecure, "", false, "skip TLS certificate verification")
	flag.MarkFlagsRequired(cmd, ArgFromBlock, ArgToBlock, ArgRPCURL)
}

// GetSmcOptions holds options for smart contract retrieval
type GetSmcOptions struct {
	BridgeAddress string
}

// AddFlags adds smart contract option flags to a command
func (o *GetSmcOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.BridgeAddress, ArgBridgeAddress, "a", "", "address of the ulxly bridge")
}

// GetVerifyBatchesOptions holds options for verify batches retrieval
type GetVerifyBatchesOptions struct {
	RollupManagerAddress string
}

// AddFlags adds verify batches option flags to a command
func (o *GetVerifyBatchesOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.RollupManagerAddress, ArgRollupManagerAddress, "a", "", "address of the rollup manager contract")
}

// DecodeGlobalIndex decodes a global index into its components
func DecodeGlobalIndex(globalIndex *big.Int) (bool, uint32, uint32, error) {
	const lengthGlobalIndexInBytes = 32
	var buf [32]byte
	gIBytes := globalIndex.FillBytes(buf[:])
	if len(gIBytes) != lengthGlobalIndexInBytes {
		return false, 0, 0, fmt.Errorf("invalid globalIndex length. Should be 32. Current length: %d", len(gIBytes))
	}
	mainnetFlag := big.NewInt(0).SetBytes([]byte{gIBytes[23]}).Uint64() == 1
	rollupIndex := big.NewInt(0).SetBytes(gIBytes[24:28])
	localRootIndex := big.NewInt(0).SetBytes(gIBytes[28:32])
	if rollupIndex.Uint64() > 0xFFFFFFFF {
		return false, 0, 0, fmt.Errorf("invalid rollupIndex length. Should be fit into uint32 type")
	}
	if localRootIndex.Uint64() > 0xFFFFFFFF {
		return false, 0, 0, fmt.Errorf("invalid localRootIndex length. Should be fit into uint32 type")
	}
	return mainnetFlag, uint32(rollupIndex.Uint64()), uint32(localRootIndex.Uint64()), nil //nolint:gosec
}

// InputArgs is the global instance of UlxlyArgs
var InputArgs = UlxlyArgs{}

var (
	BridgeService  bridge_service.BridgeService
	BridgeServices = make(map[uint32]bridge_service.BridgeService)
)

// PrepInputs prepares common inputs for ulxly commands.
// It loads RPC URL and private key from environment/config if not set via flags,
// then validates and processes all inputs.
func PrepInputs(cmd *cobra.Command, _ []string) (err error) {
	// Load RPC URL and private key from env/config if not provided via flags
	InputArgs.RPCURL, err = flag.GetRequiredRPCURL(cmd)
	if err != nil {
		return err
	}
	InputArgs.PrivateKey, err = flag.GetRequiredPrivateKey(cmd)
	if err != nil {
		return err
	}

	if InputArgs.DryRun && InputArgs.GasLimit == 0 {
		InputArgs.GasLimit = uint64(10_000_000)
	}
	pvtKey := strings.TrimPrefix(InputArgs.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(pvtKey)
	if err != nil {
		return fmt.Errorf("invalid --%s: %w", ArgPrivateKey, err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	InputArgs.AddressOfPrivateKey = fromAddress.String()
	if InputArgs.DestAddress == "" {
		InputArgs.DestAddress = fromAddress.String()
		log.Info().Stringer("destAddress", fromAddress).Msg("No destination address specified. Using private key's address")
	}

	if InputArgs.CallDataFile != "" {
		rawCallData, iErr := os.ReadFile(InputArgs.CallDataFile)
		if iErr != nil {
			return iErr
		}
		if InputArgs.CallData != "0x" {
			return fmt.Errorf("both %s and %s flags were provided", ArgCallData, ArgCallDataFile)
		}
		InputArgs.CallData = string(rawCallData)
	}

	BridgeService, err = bridge_service_factory.NewBridgeService(InputArgs.BridgeServiceURL, InputArgs.Insecure, InputArgs.Legacy)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create bridge service")
		return err
	}

	bridgeServicesURLs, err := GetBridgeServiceURLs()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get bridge service URLs")
		return err
	}

	for networkID, url := range bridgeServicesURLs {
		bs, bsErr := bridge_service_factory.NewBridgeService(url, InputArgs.Insecure, InputArgs.Legacy)
		if bsErr != nil {
			log.Error().Err(bsErr).Str("url", url).Msg("Unable to create bridge service")
			return bsErr
		}
		if _, exists := BridgeServices[networkID]; exists {
			log.Warn().Uint32("networkID", networkID).Str("url", url).Msg("Duplicate network ID found for bridge service URL. Overwriting previous entry.")
		}
		BridgeServices[networkID] = bs
		log.Info().Uint32("networkID", networkID).Str("url", url).Msg("Added bridge service")
	}

	return nil
}

// GetBridgeServiceURLs returns a map of network IDs to bridge service URLs
func GetBridgeServiceURLs() (map[uint32]string, error) {
	bridgeServiceUrls := InputArgs.BridgeServiceURLs
	urlMap := make(map[uint32]string)
	for _, mapping := range bridgeServiceUrls {
		pieces := strings.Split(mapping, "=")
		if len(pieces) != 2 {
			return nil, fmt.Errorf("bridge service url mapping should contain a networkid and url separated by an equal sign. Got: %s", mapping)
		}
		networkID, err := strconv.ParseInt(pieces[0], 10, 32)
		if err != nil {
			return nil, err
		}
		urlMap[uint32(networkID)] = pieces[1]
	}
	return urlMap, nil
}

// CreateInsecureEthClient creates an Ethereum client with TLS certificate verification disabled
func CreateInsecureEthClient(rpcURL string) (*ethclient.Client, error) {
	// WARNING: This disables TLS certificate verification
	log.Warn().Msg("WARNING: TLS certificate verification is disabled. This is unsafe for production use.")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}

	rpcClient, err := ethrpc.DialOptions(context.Background(), rpcURL, ethrpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}

// CreateEthClient creates either a secure or insecure client based on the Insecure flag
func CreateEthClient(ctx context.Context, rpcURL string) (*ethclient.Client, error) {
	if InputArgs.Insecure {
		return CreateInsecureEthClient(rpcURL)
	}
	return ethclient.DialContext(ctx, rpcURL)
}

// GenerateTransactionPayload generates the transaction payload for bridge operations
func GenerateTransactionPayload(ctx context.Context, client *ethclient.Client, ulxlyInputArgBridge string, ulxlyInputArgPvtKey string, ulxlyInputArgGasLimit uint64, ulxlyInputArgDestAddr string, ulxlyInputArgChainID string) (bridgeV2 *ulxly.Ulxly, toAddress ethcommon.Address, opts *bind.TransactOpts, err error) {
	// checks if bridge address has code
	err = ensureCodePresent(ctx, client, ulxlyInputArgBridge)
	if err != nil {
		err = fmt.Errorf("bridge code check err: %w", err)
		return
	}

	ulxlyInputArgPvtKey = strings.TrimPrefix(ulxlyInputArgPvtKey, "0x")
	bridgeV2, err = ulxly.NewUlxly(ethcommon.HexToAddress(ulxlyInputArgBridge), client)
	if err != nil {
		return
	}

	privateKey, err := crypto.HexToECDSA(ulxlyInputArgPvtKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve private key")
		return
	}

	gasLimit := ulxlyInputArgGasLimit

	chainID := new(big.Int)
	// For manual input of chainID, use the user's input
	if ulxlyInputArgChainID != "" {
		chainID.SetString(ulxlyInputArgChainID, 10)
	} else { // If there is no user input for chainID, infer it from context
		chainID, err = client.ChainID(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Cannot get chain ID")
			return
		}
	}

	opts, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot generate transactionOpts")
		return
	}
	if InputArgs.GasPrice != "" {
		gasPrice := new(big.Int)
		gasPrice.SetString(InputArgs.GasPrice, 10)
		opts.GasPrice = gasPrice
	}
	if InputArgs.DryRun {
		opts.NoSend = true
	}
	opts.Context = ctx
	opts.GasLimit = gasLimit
	toAddress = ethcommon.HexToAddress(ulxlyInputArgDestAddr)
	if toAddress == (ethcommon.Address{}) {
		toAddress = opts.From
	}
	return bridgeV2, toAddress, opts, err
}

// ensureCodePresent checks if there is code at the given address
func ensureCodePresent(ctx context.Context, client *ethclient.Client, address string) error {
	code, err := client.CodeAt(ctx, ethcommon.HexToAddress(address), nil)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("error getting code at address")
		return err
	}
	if len(code) == 0 {
		return fmt.Errorf("address %s has no code", address)
	}
	return nil
}

// WaitMineTransaction waits for a transaction to be mined
func WaitMineTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, txTimeout uint64) error {
	if InputArgs.DryRun {
		txJSON, _ := tx.MarshalJSON()
		log.Info().RawJSON("tx", txJSON).Msg("Skipping receipt check. Dry run is enabled")
		return nil
	}
	txnMinedTimer := time.NewTimer(time.Duration(txTimeout) * time.Second)
	defer txnMinedTimer.Stop()
	for {
		select {
		case <-txnMinedTimer.C:
			log.Info().Msg("Wait timer for transaction receipt exceeded!")
			return nil
		default:
			r, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				if err.Error() != "not found" {
					log.Error().Err(err)
					return err
				}
				time.Sleep(1 * time.Second)
				continue
			}
			if r.Status != 0 {
				log.Info().Interface("txHash", r.TxHash).Msg("transaction successful")
				return nil
			} else if r.Status == 0 {
				log.Error().Interface("txHash", r.TxHash).Msg("Deposit transaction failed")
				log.Info().Uint64("GasUsed", tx.Gas()).Uint64("cumulativeGasUsedForTx", r.CumulativeGasUsed).Msg("Perhaps try increasing the gas limit")
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// AddTransactionFlags adds the common transaction flags shared by bridge and claim commands.
// These flags are needed for any command that sends transactions.
func AddTransactionFlags(cmd *cobra.Command) {
	f := cmd.PersistentFlags()
	f.StringVar(&InputArgs.RPCURL, ArgRPCURL, "", "RPC URL to send the transaction")
	f.StringVar(&InputArgs.BridgeAddress, ArgBridgeAddress, "", "address of the lxly bridge")
	f.Uint64Var(&InputArgs.GasLimit, ArgGasLimit, 0, "force specific gas limit for transaction")
	f.StringVar(&InputArgs.ChainID, ArgChainID, "", "chain ID to use in the transaction")
	f.StringVar(&InputArgs.PrivateKey, ArgPrivateKey, "", "hex encoded private key for sending transaction")
	f.StringVar(&InputArgs.DestAddress, ArgDestAddress, "", "destination address for the bridge")
	f.Uint64Var(&InputArgs.Timeout, ArgTimeout, 60, "timeout in seconds to wait for transaction receipt confirmation")
	f.StringVar(&InputArgs.GasPrice, ArgGasPrice, "", "gas price to use")
	f.BoolVar(&InputArgs.DryRun, ArgDryRun, false, "do all of the transaction steps but do not send the transaction")
	f.BoolVar(&InputArgs.Insecure, ArgInsecure, false, "skip TLS certificate verification")
	f.BoolVar(&InputArgs.Legacy, ArgLegacy, true, "force usage of legacy bridge service")
	flag.MarkPersistentFlagsRequired(cmd, ArgBridgeAddress)
}
