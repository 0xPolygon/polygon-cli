package retest

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
	"reflect"
)

var (
	//go:embed usage.md
	usage         string
	inputFileName *string
)

type EthTestEnv struct {
	CurrentBaseFee       EthTestNumeric `json:"currentBaseFee"`
	CurrentCoinbase      EthTestAddress `json:"currentCoinbase"`
	CurrentDifficulty    EthTestNumeric `json:"currentDifficulty"`
	CurrentExcessBlobGas EthTestNumeric `json:"currentExcessBlobGas"`
	CurrentGasLimit      EthTestNumeric `json:"currentGasLimit"`
	CurrentNumber        EthTestNumeric `json:"currentNumber"`
	CurrentRandom        EthTestNumeric `json:"currentRandom"`
	CurrentTimestamp     EthTestNumeric `json:"currentTimestamp"`
}
type EthTestAccessList struct {
	Address     EthTestAddress   `json:"address"`
	StorageKeys []EthTestNumeric `json:"storageKeys"`
}
type EthTestTx struct {
	AccessList          []EthTestAccessList `json:"accessList"`
	BlobVersionedHashes []EthTestHash       `json:"blobVersionedHashes"`
	ChainID             EthTestNumeric      `json:"chainId"`
	Data                EthTestData         `json:"data"`

	GasLimit             EthTestNumeric `json:"gasLimit"`
	GasPrice             EthTestNumeric `json:"gasPrice"`
	MaxFeePerBlobGas     EthTestNumeric `json:"maxFeePerBlobGas"`
	MaxFeePerGas         EthTestNumeric `json:"maxFeePerGas"`
	MaxPriorityFeePerGas EthTestNumeric `json:"maxPriorityFeePerGas"`
	Nonce                EthTestNumeric `json:"nonce"`
	R                    EthTestNumeric `json:"r"`
	S                    EthTestNumeric `json:"s"`
	To                   EthTestAddress `json:"to"`
	V                    EthTestNumeric `json:"v"`
	Value                EthTestNumeric `json:"value"`
	// Unused
	ExpectException any `json:"expectException"`
	SecretKey       any `json:"secretKey"`
}
type EthTestGenesis struct {
	BaseFeePerGas         EthTestNumeric `json:"baseFeePerGas"`
	BeaconRoot            EthTestHash    `json:"beaconRoot"`
	BlobGasUsed           EthTestNumeric `json:"blobGasUsed"`
	Bloom                 EthTestData    `json:"bloom"`
	Coinbase              EthTestAddress `json:"coinbase"`
	Difficulty            EthTestNumeric `json:"difficulty"`
	ExcessBlobGas         EthTestNumeric `json:"excessBlobGas"`
	ExtraData             EthTestData    `json:"extraData"`
	GasLimit              EthTestNumeric `json:"gasLimit"`
	GasUsed               EthTestNumeric `json:"gasUsed"`
	Hash                  EthTestHash    `json:"hash"`
	MixHash               EthTestHash    `json:"mixHash"`
	Nonce                 EthTestNumeric `json:"nonce"`
	Number                EthTestNumeric `json:"number"`
	ParentBeaconBlockRoot EthTestHash    `json:"parentBeaconBlockRoot"`
	ParentHash            EthTestHash    `json:"parentHash"`
	ReceiptTrie           EthTestHash    `json:"receiptTrie"`
	StateRoot             EthTestHash    `json:"stateRoot"`
	Timestamp             EthTestNumeric `json:"timestamp"`
	TransactionsTrie      EthTestHash    `json:"transactionsTrie"`
	UncleHash             EthTestHash    `json:"uncleHash"`
	WithdrawalsRoot       EthTestHash    `json:"withdrawalsRoot"`
}
type EthTestBlocks struct {
	Transactions []EthTestTx `json:"transactions"`
}

type EthTestPre struct {
	Nonce   EthTestNumeric            `json:"nonce"`
	Balance EthTestNumeric            `json:"balance"`
	Storage map[string]EthTestNumeric `json:"storage"`
	Code    EthTestData               `json:"code"`
}

// This is based on examination of the data. We've left out the fields at are used in less then 10 tests
type EthTest struct {
	Pre                map[string]EthTestPre `json:"pre"`                // Pre determines preallocated accounts in the test case
	Transaction        EthTestTx             `json:"transaction"`        // Transaction is usually a test transaction
	Env                EthTestEnv            `json:"env"`                // Env sets the environment for execution. This is not relvant usually for creating load
	GenesisBlockHeader EthTestGenesis        `json:"genesisBlockHeader"` // GenesisBlockHeader sets the typical genesis params of the network.
	Blocks             []EthTestBlocks       `json:"blocks"`             // Blocks are the blocks that lead to the state being tested
	SealEngine         string                `json:"sealEngine"`         // SealEngine is either Null or NoProof
	Expect             any                   `json:"expect"`             // Expect are the success conditions. In our cases we won't be able to validate them
	Info               any                   `json:"_info"`              // Info mostly comments
	ExpectException    any                   `json:"expectException"`    // ExpectException also determins if errors are expected for certain levels of HF
	Vectors            any                   `json:"vectors"`            // Vectors are specific cases which we probably won't use
	TxBytes            any                   `json:"txbytes"`            // TxBytes are RLP tests to send directly without manipulation
	Result             any                   `json:"result"`             // Result are specific expected results
	Network            any                   `json:"network"`            // Network in rare cases seems to specify ArrayGlacier, Istanbul, or GrayGlacier
	Exceptions         any                   `json:"exceptions"`         // Exceptions specifies resulting errors
}

type EthTestSuite map[string]EthTest

type EthTestNumeric any
type EthTestData any
type EthTestHash any
type EthTestAddress any

func decodeToUint64(v any) (uint64, error) {
	switch v := reflect.ValueOf(v); v.Kind() {
	case reflect.Float64:
		// fmt.Println(v.Int())
	case reflect.String:
	// 0x0abc123
	case reflect.Slice:
		// idk
	default:
		fmt.Printf("unhandled kind %s\n", v.Kind())
	}
	return 0, nil
}

var RetestCmd = &cobra.Command{
	Use:   "retest [flags]",
	Short: "Convert the standard ETH test fillers into something to be replayed against an RPC",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("starting")
		rawData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		tests := make(EthTestSuite, 0)
		err = json.Unmarshal(rawData, &tests)
		if err != nil {
			return err
		}
		for k := range tests {
			log.Debug().Str("testname", k).Msg("Parsing test")
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := RetestCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a file that's filed with test transaction fillers")
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	return io.ReadAll(os.Stdin)
}
