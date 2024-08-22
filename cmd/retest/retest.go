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
	CurrentBaseFee       any `json:"currentBaseFee"`
	CurrentCoinbase      any `json:"currentCoinbase"`
	CurrentDifficulty    any `json:"currentDifficulty"`
	CurrentExcessBlobGas any `json:"currentExcessBlobGas"`
	CurrentGasLimit      any `json:"currentGasLimit"`
	CurrentNumber        any `json:"currentNumber"`
	CurrentRandom        any `json:"currentRandom"`
	CurrentTimestamp     any `json:"currentTimestamp"`
}
type EthTestTx struct {
	AccessList           any `json:"accessList"`
	BlobVersionedHashes  any `json:"blobVersionedHashes"`
	ChainID              any `json:"chainId"`
	Data                 any `json:"data"`
	ExpectException      any `json:"expectException"`
	GasLimit             any `json:"gasLimit"`
	GasPrice             any `json:"gasPrice"`
	MaxFeePerBlobGas     any `json:"maxFeePerBlobGas"`
	MaxFeePerGas         any `json:"maxFeePerGas"`
	MaxPriorityFeePerGas any `json:"maxPriorityFeePerGas"`
	Nonce                any `json:"nonce"`
	R                    any `json:"r"`
	S                    any `json:"s"`
	SecretKey            any `json:"secretKey"`
	To                   any `json:"to"`
	V                    any `json:"v"`
	Value                any `json:"value"`
}
type EthTestGenesis struct {
	BaseFeePerGas         any `json:"baseFeePerGas"`
	BeaconRoot            any `json:"beaconRoot"`
	BlobGasUsed           any `json:"blobGasUsed"`
	Bloom                 any `json:"bloom"`
	Coinbase              any `json:"coinbase"`
	Difficulty            any `json:"difficulty"`
	ExcessBlobGas         any `json:"excessBlobGas"`
	ExtraData             any `json:"extraData"`
	GasLimit              any `json:"gasLimit"`
	GasUsed               any `json:"gasUsed"`
	Hash                  any `json:"hash"`
	MixHash               any `json:"mixHash"`
	Nonce                 any `json:"nonce"`
	Number                any `json:"number"`
	ParentBeaconBlockRoot any `json:"parentBeaconBlockRoot"`
	ParentHash            any `json:"parentHash"`
	ReceiptTrie           any `json:"receiptTrie"`
	StateRoot             any `json:"stateRoot"`
	Timestamp             any `json:"timestamp"`
	TransactionsTrie      any `json:"transactionsTrie"`
	UncleHash             any `json:"uncleHash"`
	WithdrawalsRoot       any `json:"withdrawalsRoot"`
}
type EthTestBlocks struct {
	Transactions []EthTestTx `json:"transactions"`
}

type EthTestPre struct {
	Nonce   any `json:"nonce"`
	Balance any `json:"balance"`
	Storage any `json:"storage"`
	Code    any `json:"code"`
}

// This is based on examination of the data. We've left out the fields at are used in less then 10 tests
type EthTest struct {
	Pre                map[string]EthTestPre `json:"pre"`                // Pre determines preallocated accounts in the test case
	Expect             any                   `json:"expect"`             // Expect are the success conditions. In our cases we won't be able to validate them
	Transaction        EthTestTx             `json:"transaction"`        // Transaction is usually a test transaction
	Env                EthTestEnv            `json:"env"`                // Env sets the environment for execution. This is not relvant usually for creating load
	Info               any                   `json:"_info"`              // Info mostly comments
	GenesisBlockHeader EthTestGenesis        `json:"genesisBlockHeader"` // GenesisBlockHeader sets the typical genesis params of the network.
	Blocks             []EthTestBlocks       `json:"blocks"`             // Blocks are the blocks that lead to the state being tested
	SealEngine         string                `json:"sealEngine"`         // SealEngine is either Null or NoProof
	ExpectException    any                   `json:"expectException"`    // ExpectException also determins if errors are expected for certain levels of HF
	Vectors            any                   `json:"vectors"`            // Vectors are specific cases which we probably won't use
	TxBytes            any                   `json:"txbytes"`            // TxBytes are RLP tests to send directly without manipulation
	Result             any                   `json:"result"`             // Result are specific expected results
	Network            any                   `json:"network"`            // Network in rare cases seems to specify ArrayGlacier, Istanbul, or GrayGlacier
	Exceptions         any                   `json:"exceptions"`         // Exceptions specifies resulting errors
}

type EthTestSuite map[string]EthTest

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
