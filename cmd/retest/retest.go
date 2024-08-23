package retest

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/maticnetwork/polygon-cli/abi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"math"
	"math/big"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

var (
	//go:embed usage.md
	usage         string
	inputFileName *string

	validBase10    *regexp.Regexp
	dataLabel      *regexp.Regexp
	typeIndidcator *regexp.Regexp
	abiSpec        *regexp.Regexp
	normalizeWs    *regexp.Regexp
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
	AccessList           []EthTestAccessList `json:"accessList"`
	BlobVersionedHashes  []EthTestHash       `json:"blobVersionedHashes"`
	ChainID              EthTestNumeric      `json:"chainId"`
	Data                 EthTestData         `json:"data"`
	GasLimit             EthTestNumeric      `json:"gasLimit"`
	GasPrice             EthTestNumeric      `json:"gasPrice"`
	MaxFeePerBlobGas     EthTestNumeric      `json:"maxFeePerBlobGas"`
	MaxFeePerGas         EthTestNumeric      `json:"maxFeePerGas"`
	MaxPriorityFeePerGas EthTestNumeric      `json:"maxPriorityFeePerGas"`
	Nonce                EthTestNumeric      `json:"nonce"`
	R                    EthTestNumeric      `json:"r"`
	S                    EthTestNumeric      `json:"s"`
	To                   EthTestAddress      `json:"to"`
	V                    EthTestNumeric      `json:"v"`
	Value                EthTestNumeric      `json:"value"`
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

// EthTest is based on examination of the data. We've left out the fields at are used in less then 10 tests
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

type WrappedNumeric struct {
	raw   EthTestNumeric
	inner *big.Int
}

type WrappedData struct {
	raw   EthTestData
	inner string
}

func (wr *WrappedNumeric) ToBigInt() *big.Int {
	if wr.inner != nil {
		return wr.inner
	}
	wr.inner = EthTestNumericToBigInt(wr.raw)
	return wr.inner
}

func EthTestNumericToBigInt(num EthTestNumeric) *big.Int {
	if num == nil {
		return nil
	}
	v := reflect.ValueOf(num)
	if v.IsZero() {
		return nil
	}

	switch v.Kind() {
	case reflect.Float64:
		return float64ToBigInt(v.Float())
	case reflect.String:
		return processNumericString(v.String())
	case reflect.Slice:
		if v.Len() == 0 {
			log.Warn().Msg("The slice is empty; returning nil")
			return nil
		}
		first := v.Index(0)
		// TODO this indicates this is a matrixed parameter.. WE're losing out right now
		log.Debug().Any("first", first).Msg("A numeric field is multi-valued. This isn't currently supported and we'll use the first element")
		return EthTestNumericToBigInt(first.Interface().(EthTestNumeric))
	default:
		log.Fatal().Any("input", num).Str("kind", v.Kind().String()).Msg("Attempted to convert unknown type to number")
	}
	return new(big.Int)
}
func float64ToBigInt(f float64) *big.Int {
	// Step 1: Round the float64 (e.g., using math.Round)
	roundedFloat := math.Round(f)

	// Step 2: Convert the rounded float to a string
	str := fmt.Sprintf("%.0f", roundedFloat)

	// Step 3: Create a big.Int and set it from the string
	bigInt := new(big.Int)
	bigInt.SetString(str, 10)

	return bigInt
}

func processNumericString(s string) *big.Int {
	log.Trace().Str("numString", s).Msg("Converting numeric string to big int")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x:bigint ")
	if strings.Contains(s, ":") {
		log.Fatal().Str("string", s).Msg("Unknown number format")
	}
	if strings.HasPrefix(s, "0x") {
		s = strings.TrimPrefix(s, "0x")
		num, fullRead := new(big.Int).SetString(s, 16)
		if !fullRead {
			log.Fatal().Str("input", s).Msg("Unable to read the full hex data?!")
		}
		return num
	}

	if validBase10.MatchString(s) {
		num, fullRead := new(big.Int).SetString(s, 10)
		if !fullRead {
			log.Fatal().Str("input", s).Msg("Unable to read the full numeric data")
		}
		return num
	}
	num, fullRead := new(big.Int).SetString(s, 16)
	if !fullRead {
		log.Fatal().Str("input", s).Msg("Unable to read the full numeric data")
	}
	return num
}

func (wd *WrappedData) ToString() string {
	if wd.inner != "" {
		return wd.inner
	}
	wd.inner = EthTestDataToString(wd.raw)
	return wd.inner
}
func EthTestDataToString(data EthTestData) string {
	if data == nil {
		return ""
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.String:
		return processTestDataString(v.String())
	case reflect.Map:
		dataMap := data.(map[string]any)
		_, hasData := dataMap["data"]
		_, hasAccessList := dataMap["accessList"]
		if !hasData || !hasAccessList {
			log.Fatal().Msg("Got a data field with a map type that wasn't data + access list")
		}
		// TODO - we're losing the ability to send lots of different access list tests
		return EthTestDataToString(dataMap["data"])
	case reflect.Slice:
		if v.Len() == 0 {
			log.Warn().Msg("The slice is empty; returning nil")
			return ""
		}
		// TODO these tests should be collected somehow
		for i := 0; i < v.Len(); i = i + 1 {
			EthTestDataToString(v.Index(i).Interface().(EthTestData))
		}
		return EthTestDataToString(v.Index(0).Interface().(EthTestData))
	default:
		log.Fatal().Any("input", data).Str("kind", v.Kind().String()).Msg("Attempted to convert unknown type to raw data")
	}
	return ""
}

func processTestDataString(data string) string {
	data = strings.TrimSpace(data)
	if dataLabel.MatchString(data) {
		label := dataLabel.FindStringSubmatch(data)
		data = dataLabel.ReplaceAllString(data, "")
		log.Trace().Str("label", label[1]).Msg("stripping label")
	}
	data = strings.TrimSpace(data)
	if data == "" {
		return ""
	}
	if typeIndidcator.MatchString(data) {
		rawType := typeIndidcator.FindStringSubmatch(data)[1]
		switch rawType {
		case "raw":
			return processRawStringToString(data)
		case "yul":
			return processYulToString(data)
		case "abi":
			return processAbiStringToString(data)
		default:
			log.Fatal().Str("type", rawType).Msg("unknown type designation")
		}
	} else if strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}") {
		return processLLLToString(data)
	} else if strings.HasPrefix(data, "0x") {
		return processRawStringToString(data)
	} else {
		log.Fatal().Str("data", data).Msg("unknown data format")
	}

	return ""
}

func processYulToString(data string) string {
	data = preProcessTypedString(data, true)
	if !strings.HasPrefix(data, "berlin ") {
		// at this point it seems like every yul contract is prefixed with berlin
		// https://github.com/ethereum/tests/commit/fd26aad70e24f042fcd135b2f0338b1c6bf1a324
		log.Fatal().Str("contract", data).Msg("The contract didn't have a berlin prefix")
	}
	data = strings.TrimPrefix(data, "berlin")
	yulInput, err := os.CreateTemp("", "yul-")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create yul file")
	}
	_, err = yulInput.WriteString(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to write yul file")
	}
	err = yulInput.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to close yul file")
	}
	cmd := exec.Command("solc", "--strict-assembly", "--evm-version", "berlin", "--bin", "--input-file", yulInput.Name())
	solcOut, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Str("contract", data).Msg("there was an error running solc/solidity for yul contracts")
	}
	lines := strings.Split(string(solcOut), "\n")
	if len(lines) != 6 {
		log.Fatal().Int("lines", len(lines)).Str("contract", data).Msg("YUL contract does not contain 6 lines")
	}
	return lines[len(lines)-2]
}

func processLLLToString(data string) string {
	data = preProcessTypedString(data, true)
	lllcInput, err := os.CreateTemp("", "lllc-")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create LLL file")
	}
	_, err = lllcInput.WriteString(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to write lll file")
	}
	err = lllcInput.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to close lll file")
	}
	cmd := exec.Command("lllc", lllcInput.Name())
	lllcOut, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Str("contract", data).Msg("there was an error running solc/solidity for yul contracts")
	}
	lines := strings.Split(string(lllcOut), "\n")
	if len(lines) != 2 {
		log.Fatal().Int("lines", len(lines)).Str("contract", data).Msg("LLLC output does not contain 2 lines")
	}
	fmt.Println(lines[0])
	return lines[0]
}

func processAbiStringToString(data string) string {
	data = preProcessTypedString(data, true)
	matches := abiSpec.FindAllStringSubmatch(data, -1)
	if len(matches) != 1 {
		log.Fatal().Int("matches", len(matches)).Str("abi", data).Msg("unrecognized abi spec")
	}
	if len(matches[0]) != 4 {
		log.Fatal().Int("matches", len(matches[0])).Str("abi", data).Msg("unrecognized abi spec")
	}
	funcName := matches[0][1]
	funcParams := matches[0][2]
	funcInputs := matches[0][3]
	params := strings.Split(funcParams, ",")
	processedArgs := rawArgsToStrings(funcInputs, params)
	encodedArgs, err := abi.AbiEncode(fmt.Sprintf("%s(%s)", funcName, funcParams), processedArgs)
	if err != nil {
		log.Fatal().Err(err).Str("funcName", funcName).Str("funcParams", funcParams).Str("funcInputs", funcInputs).Msg("failed to encode args in abi")
	}

	return encodedArgs
}
func rawArgsToStrings(rawArgs string, params []string) []string {
	rawArgs = strings.TrimSpace(rawArgs)
	if rawArgs == "" {
		return []string{}
	}
	count := len(params)
	rawArgs = strings.ReplaceAll(rawArgs, "0x", " 0x")
	rawArgs = normalizeWs.ReplaceAllString(rawArgs, " ")
	argList := strings.Split(rawArgs, " ")
	if argList[0] == "" {
		argList = argList[1:]
	}
	if len(argList) == 1 && count > 1 {
		for k := 1; k < count; k += 1 {
			argList = append(argList, argList[0])
		}
	}
	if len(argList) != count {
		log.Fatal().Str("rawArgs", rawArgs).Int("argListLength", len(argList)).Int("paramCount", count).Msg("arg length mismatch")
	}

	processedArgs := make([]string, count)
	for k, arg := range argList {
		if strings.HasPrefix(params[k], "uint") {
			if strings.HasPrefix(arg, "0x") {
				arg = strings.TrimPrefix(arg, "0x")
				if len(arg) > 64 {
					// i think this is a bug but there is a test case that's somehow longer than 32 bytes
					// https://github.com/ethereum/tests/blob/fd26aad70e24f042fcd135b2f0338b1c6bf1a324/src/GeneralStateTestsFiller/Cancun/stEIP1153-transientStorage/transStorageOKFiller.yml#L801
					arg = arg[len(arg)-64:]
				}
				n, _ := new(big.Int).SetString(arg, 16)
				processedArgs[k] = n.String()
			} else {
				processedArgs[k] = arg
			}
		} else if params[k] == "bool" {
			if arg == "0x01" {
				processedArgs[k] = "true"
			} else if arg == "0x00" {
				processedArgs[k] = "false"
			} else {
				log.Fatal().Str("arg", arg).Msg("unrecognized bool type input")
			}
		} else {
			log.Fatal().Str("type", params[k]).Msg("unknown type designation")
		}
	}

	return processedArgs
}

func preProcessTypedString(data string, preserveSpace bool) string {
	data = strings.TrimSpace(data)
	data = typeIndidcator.ReplaceAllString(data, "")
	data = strings.TrimPrefix(data, "0x")
	if !preserveSpace {
		data = strings.Replace(data, " ", "", -1)
	}
	return data
}

func processRawStringToString(data string) string {
	data = preProcessTypedString(data, false)

	byteData, err := hex.DecodeString(data)
	if err != nil {
		log.Fatal().Str("data", data).Err(err).Msg("Unable to decode the raw data")
	}
	return "0x" + hex.EncodeToString(byteData)
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
			n := tests[k].Transaction.GasLimit
			wn := WrappedNumeric{raw: n}
			if wn.ToBigInt() != nil {
				log.Trace().Uint64("nonce", wn.ToBigInt().Uint64()).Msg("Parsing nonce")
			}
			d := tests[k].Transaction.Data
			wd := WrappedData{raw: d}
			wd.ToString()
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

	validBase10 = regexp.MustCompile(`^[0-9]*$`)
	dataLabel = regexp.MustCompile(`^:label ([^ ]*) `)
	typeIndidcator = regexp.MustCompile(`^:([^ ]*) `)
	abiSpec = regexp.MustCompile(`^([a-zA-Z0-9]*)\((.*)\)(.*)$`)
	normalizeWs = regexp.MustCompile(` +`)
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	return io.ReadAll(os.Stdin)
}
