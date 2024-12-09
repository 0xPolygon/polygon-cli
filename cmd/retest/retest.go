package retest

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/0xPolygon/polygon-cli/abi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/0xPolygon/polygon-cli/util"
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
	solcPath      string
	inputFileName *string

	validBase10         *regexp.Regexp
	dataLabel           *regexp.Regexp
	typeIndicator       *regexp.Regexp
	abiSpec             *regexp.Regexp
	normalizeWs         *regexp.Regexp
	solidityCompileInfo *regexp.Regexp
	removablePreamble   *regexp.Regexp
	solcCompileMultiOut *regexp.Regexp
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
	ExpectException    any                   `json:"expectException"`    // ExpectException also determines if errors are expected for certain levels of HF
	Vectors            any                   `json:"vectors"`            // Vectors are specific cases which we probably won't use
	TxBytes            any                   `json:"txbytes"`            // TxBytes are RLP tests to send directly without manipulation
	Result             any                   `json:"result"`             // Result are specific expected results
	Network            any                   `json:"network"`            // Network in rare cases seems to specify ArrayGlacier, Istanbul, or GrayGlacier
	Exceptions         any                   `json:"exceptions"`         // Exceptions specifies resulting errors
	Solidity           string                `json:"solidity"`           // Solidity contains a standard solidity file with one or more contracts
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
type WrappedAddress struct {
	raw   EthTestAddress
	inner *ethcommon.Address
}

func (wr *WrappedAddress) ToString() *ethcommon.Address {
	if wr.inner != nil {
		return wr.inner
	}
	v := reflect.ValueOf(wr.raw)
	if v.Kind() == reflect.Invalid {
		wr.inner = new(ethcommon.Address)
		return wr.inner
	}
	if v.Kind() == reflect.Float64 {
		// FIXME this case shouldn't be necessary
		// this is a weird case where the address is specified as a number... There are a dozen or so cases that seem
		// like in the yml to json conversion, numbers like 095e7baea6a6c7c4c2dfeb977efac326af552d87 are prefixed with
		// 0x and interpreted as a number rather than a string. This seems to work fine in retest ETH but cases an issue
		// in this workflow
		// GeneralStateTestsFiller/Cancun/stEIP4844-blobtransactions/blobhashListBounds5Filler.yml
		f := new(big.Float)
		f = f.SetFloat64(v.Interface().(float64))
		fInt, _ := f.Int(nil)
		addr := ethcommon.BytesToAddress(fInt.Bytes())
		wr.inner = &addr
		return wr.inner
	}

	if v.Kind() != reflect.String {
		log.Fatal().Any("addr", wr.raw).Str("kind", v.Kind().String()).Msg("unknown source address type")
	}
	addr := ethcommon.HexToAddress(v.Interface().(string))
	wr.inner = &addr
	return wr.inner
}
func (wr *WrappedNumeric) ToBigInt() *big.Int {
	if wr.inner != nil {
		return wr.inner
	}
	wr.inner = EthTestNumericToBigInt(wr.raw)
	if wr.inner == nil {
		wr.inner = new(big.Int)
	}
	return wr.inner
}

func (wr *WrappedNumeric) ToString() string {
	bi := wr.ToBigInt()
	return hexutil.EncodeBig(bi)
}

func (wr *WrappedData) IsSlice() bool {
	v := reflect.ValueOf(wr.raw)
	k := v.Kind()
	return k == reflect.Slice
}
func (wr *WrappedData) ToSlice() []*WrappedData {
	if !wr.IsSlice() {
		return []*WrappedData{wr}
	}
	v := reflect.ValueOf(wr.raw)
	if v.Len() == 0 {
		return []*WrappedData{}
	}
	wrappedDatas := make([]*WrappedData, v.Len())
	for i := 0; i < v.Len(); i = i + 1 {
		nwd := new(WrappedData)
		nwd.raw = v.Index(i).Interface().(EthTestData)
		wrappedDatas[i] = nwd

	}
	return wrappedDatas
}
func (wn *WrappedNumeric) IsSlice() bool {
	v := reflect.ValueOf(wn.raw)
	k := v.Kind()
	return k == reflect.Slice
}
func (wn *WrappedNumeric) ToSlice() []*WrappedNumeric {
	if !wn.IsSlice() {
		return []*WrappedNumeric{wn}
	}
	v := reflect.ValueOf(wn.raw)
	if v.Len() == 0 {
		return []*WrappedNumeric{}
	}
	wrappedDatas := make([]*WrappedNumeric, v.Len())
	for i := 0; i < v.Len(); i = i + 1 {
		nwd := new(WrappedNumeric)
		nwd.raw = v.Index(i).Interface().(EthTestData)
		wrappedDatas[i] = nwd

	}
	return wrappedDatas
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
		// remove the separators that might be in the string like 100_000
		s = strings.ReplaceAll(s, "_", "")
		num, fullRead := new(big.Int).SetString(s, 10)
		if !fullRead {
			log.Fatal().Str("input", s).Msg("Unable to read the full numeric data in base 10")
		}
		return num
	}
	num, fullRead := new(big.Int).SetString(s, 16)
	if !fullRead {
		log.Fatal().Str("input", s).Msg("Unable to read the full numeric data in base 16")
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
	case reflect.Float64:
		// We have few tests with numeric code, ex:
		// "code": 16449,
		return processStorageData(data.(float64))
	
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
	if typeIndicator.MatchString(data) {
		rawType := typeIndicator.FindStringSubmatch(data)[1]
		switch rawType {
		case "raw":
			return processRawStringToString(data)
		case "yul":
			return processSolidityToString(data, true)
		case "abi":
			return processAbiStringToString(data)
		case "solidity":
			return processSolidityToString(data, false)
		default:
			log.Fatal().Str("type", rawType).Msg("unknown type designation")
		}
	} else if strings.HasPrefix(data, "{") {
		return processLLLToString(data)
	} else if strings.HasPrefix(data, "0x") {
		return processRawStringToString(data)
	} else if strings.HasPrefix(data, "(asm ") {
		return processLLLToString(data)
	} else if isStandardSolidityString(data) {
		return processSolidityToString(data, false)
	} else {
		log.Fatal().Str("data", data).Msg("unknown data format")
	}

	return ""
}

func processStorageData(data any) string {
	var result string
	if reflect.TypeOf(data).Kind() == reflect.Float64 {
		result = fmt.Sprintf("%x", int64(data.(float64)))
	} else if reflect.TypeOf(data).Kind() == reflect.String {
		if strings.HasPrefix(data.(string), "0x") {
			result = strings.TrimPrefix(data.(string), "0x")
		} else {
			result = data.(string)
		}
	} else {
		log.Fatal().Any("data", data).Msg("unknown storage data type")
	}
	if len(result) % 2 != 0 {
		result = "0" + result
	}
	return result
}

// isStandardSolidityString will do a rough check to see if the string looks like a typical solidity file rather than
// the contracts that are usually in the retest code base
func isStandardSolidityString(contract string) bool {
	if strings.Contains(contract, "pragma solidity") && strings.Contains(contract, "SPDX-License-Identifier") {
		return true
	}
	return false
}

func processSolidityFlags(contract string) (string, bool) {
	if isStandardSolidityString(contract) {
		return "", true
	}
	shouldOptimize := false
	solidityInfo := solidityCompileInfo.FindStringSubmatch(contract)
	if len(solidityInfo) == 0 {
		return "", shouldOptimize
	}
	compilerOptions := strings.Split(strings.TrimSpace(solidityInfo[1]), " ")
	evmVersion := strings.TrimSpace(compilerOptions[0])
	if evmVersion == "" {
		return "", shouldOptimize
	}
	if len(compilerOptions) == 2 {
		if strings.TrimSpace(compilerOptions[1]) != "optimise" {
			log.Fatal().Str("setting", compilerOptions[1]).Msg("only aware of the optimise setting... what is this?")
		}
		shouldOptimize = true
	}
	if len(compilerOptions) > 2 && compilerOptions[1] != "object" {
		fmt.Println(contract)

		log.Fatal().Strs("opts", compilerOptions).Msg("There are more settings that we realized")
	}
	return evmVersion, shouldOptimize
}

// There are a few contracts that are structured like `london object c {`
// the goal is to remove the london part
func stripVersions(contract string) string {
	stripper := regexp.MustCompile("^(london|berlin) object")
	contract = stripper.ReplaceAllString(contract, "object")
	return contract
}

func processSolidityToString(data string, isYul bool) string {
	data = preProcessTypedString(data, true)
	solidityVersion, optimize := processSolidityFlags(data)
	matches := solidityCompileInfo.FindStringSubmatch(data)
	if len(matches) != 3 {
		fmt.Println(data)
		log.Fatal().Str("contractData", data).Msg("The format of this contract is unique and it's not clear what it is")
	}

	data = stripVersions(data)
	data = removablePreamble.ReplaceAllString(data, "")
	// data = solidityCompileInfo.ReplaceAllString(data, matches[2])
	solInput, err := os.CreateTemp("", "sol-")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create solidity file")
	}
	_, err = solInput.WriteString(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to write solidity file")
	}
	err = solInput.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to close solidity file")
	}
	args := []string{"--bin", "--input-file", solInput.Name()}
	if solidityVersion != "" {
		args = append(args, "--evm-version", solidityVersion)
	}
	if optimize {
		args = append(args, "--optimize")
	}
	if isYul {
		args = append(args, "--strict-assembly")
	}

	cmd := exec.Command(solcPath, args...)
	errOut := ""
	bufErr := bytes.NewBufferString(errOut)
	cmd.Stderr = bufErr
	solcOut, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Str("filename", solInput.Name()).Str("stdErr", bufErr.String()).Str("contract", data).Strs("args", args).Msg("there was an error running solc/solidity for contracts")
	}
	lines := strings.Split(string(solcOut), "\n")
	if len(lines) < 4 {
		log.Warn().Strs("args", args).Str("filename", solInput.Name()).Str("stdErr", bufErr.String()).Int("lines", len(lines)).Str("contract", data).Msg("soldity output does not contain 4 lines")
	}

	os.Remove(solInput.Name())
	return lines[len(lines)-2]
}

// solidityStringToBin is specifically meant for the solidity property at the test level rather than embedded contracts
func solidityStringToBin(data string) map[string]string {
	// data = solidityCompileInfo.ReplaceAllString(data, matches[2])
	solInput, err := os.CreateTemp("", "sol-")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create solidity file")
	}
	_, err = solInput.WriteString(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to write solidity file")
	}
	err = solInput.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to close solidity file")
	}
	args := []string{"--bin", "--input-file", solInput.Name()}

	cmd := exec.Command(solcPath, args...)
	errOut := ""
	bufErr := bytes.NewBufferString(errOut)
	cmd.Stderr = bufErr
	solcOut, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Str("filename", solInput.Name()).Str("stdErr", bufErr.String()).Str("contract", data).Strs("args", args).Msg("there was an error running solc/solidity for contracts")
	}
	matches := solcCompileMultiOut.FindAllStringSubmatch(string(solcOut), -1)
	contractMap := make(map[string]string)
	for _, contract := range matches {
		if len(contract) != 3 {
			log.Fatal().Int("contractLen", len(contract)).Msg("the number of matches in this compiled contract doesn't look right")
		}
		contractMap[contract[1]] = contract[2]
	}

	return contractMap
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
		log.Fatal().Err(err).Stack().Str("contract", data).Msg("there was an error compiling the lll contract")
	}
	lines := strings.Split(string(lllcOut), "\n")
	if len(lines) != 2 {
		log.Fatal().Int("lines", len(lines)).Str("contract", data).Msg("LLLC output does not contain 2 lines")
	}
	os.Remove(lllcInput.Name())
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
	data = typeIndicator.ReplaceAllString(data, "")
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

// WrapPredeployeCode will wrap a predeployed contract so that it can be deployed for testing. For now we're just wrapping
// the code so that it should match what the precondition is. In the future, this should also have a constructor that
// would initialize the storage slots to match the predeployed state. This will never be 100% right, but useful for
// smoke testing
func WrapPredeployedCode(pre EthTestPre) string {
	rawCode := WrappedData{raw: pre.Code}
	deployedCode := rawCode.ToString()
	storageInitCode := storageToByteCode(pre.Storage)

	return util.WrapDeployedCode(deployedCode, storageInitCode)
}

// storageToByteCode
func storageToByteCode(storage map[string]EthTestNumeric) string {
	if len(storage) == 0 {
		return ""
	}
	var bytecode string
	for slot, value := range storage {
		if slot == "0x" || slot == "0" {
			// special case that we encountered...
			// https://github.com/ethereum/tests/blob/fd26aad70e24f042fcd135b2f0338b1c6bf1a324/src/EIPTestsFiller/InvalidRLP/bcForgedTest/bcForkBlockTestFiller.json#L222
			log.Warn().Str("slot", slot).Msg("found a storage entry for invalid slot")
		}

		s := processStorageData(slot)
		v := processStorageData(value)
		sLen := len(s) / 2
		vLen := len(v) / 2
		sPushCode := 0x5F + sLen
		vPushCode := 0x5F + vLen
		bytecode += fmt.Sprintf("%02x%s%02x%s55", vPushCode, v, sPushCode, s)
	}

	log.Info().Str("storageInit", bytecode).Msg("produced code to initialize storage")
	return bytecode
}
func WrapCode(inputData EthTestData) string {
	rawCode := WrappedData{raw: inputData}
	return rawCode.ToString()
}

func checkContractMap(input any, contractMap map[string]string) (string, bool) {
	if len(contractMap) == 0 {
		return "", false
	}
	v := reflect.ValueOf(input)
	k := v.Kind()

	if k != reflect.String {
		return "", false
	}

	inputString := input.(string)
	if !strings.HasPrefix(inputString, ":solidity") {
		return "", false
	}
	inputString = strings.TrimSpace(strings.TrimPrefix(inputString, ":solidity "))
	for k, v := range contractMap {
		if k == inputString {
			return v, true
		}
	}
	return "", false
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

		// TODO in the future we might want to support various output modes. E.g. Assertoor, Foundry, web3.js, ethers, whatever
		simpleTests := make([]any, 0)
		for testName, t := range tests {
			// Some of the tests have a specific `solidity` field at the top level which must be compiled and then made
			// accessible to the rest of the properties
			contractMap := make(map[string]string)
			if t.Solidity != "" {
				contractMap = solidityStringToBin(t.Solidity)
			}

			st := make(map[string]any)
			st["name"] = testName
			preDeploys := make([]map[string]string, 0)
			for addr, p := range t.Pre {
				if input, matched := checkContractMap(p.Code, contractMap); matched {
					p.Code = ":raw 0x" + input
				}

				dep := make(map[string]string)
				dep["label"] = fmt.Sprintf("pre:%s", addr)
				dep["addr"] = addr
				dep["code"] = WrapPredeployedCode(p)
				preDeploys = append(preDeploys, dep)
			}
			st["dependencies"] = preDeploys

			gasLimit := WrappedNumeric{raw: t.Transaction.GasLimit}
			gasLimits := gasLimit.ToSlice()
			txValue := WrappedNumeric{raw: t.Transaction.Value}
			txValues := txValue.ToSlice()
			wd := WrappedData{raw: t.Transaction.Data}
			wds := wd.ToSlice()

			testCases := make([]map[string]any, 0)

			// This is a little bit awkward, but each test can have multiple inputs, gas limits, and values. If we have
			// 2 inputs, 2 values, and 2 gas limits, that means in total, we'll have 8 test cases
			for _, singleGas := range gasLimits {
				for _, singleValue := range txValues {
					for _, singleTx := range wds {
						tc := make(map[string]any)
						tc["name"] = testName

						if input, matched := checkContractMap(singleTx.raw, contractMap); matched {
							// as far as I can tell this isn't used yet, but i suppose it should work
							tc["input"] = input
						} else {
							tc["input"] = singleTx.ToString()
						}

						wTo := WrappedAddress{raw: t.Transaction.To}
						tc["to"] = wTo.ToString()
						tc["originalTo"] = t.Transaction.To

						tc["gas"] = singleGas.ToString()
						tc["originalGas"] = t.Transaction.GasLimit

						tc["value"] = singleValue.ToString()
						tc["originalValue"] = t.Transaction.Value

						testCases = append(testCases, tc)
					}
				}
			}

			// There is also a block filed of a test which might have a bunch of transactions.
			for _, singleBlock := range t.Blocks {
				for _, singleTx := range singleBlock.Transactions {
					// my assumption is that the tx within a block won't be multi-valued??
					tc := make(map[string]any)
					tc["name"] = testName
					if input, matched := checkContractMap(singleTx.ChainID, contractMap); matched {
						// as far as I can tell this isn't used yet, but i suppose it should work
						tc["input"] = input
					} else {
						wrappedTxData := WrappedData{raw: singleTx.Data}
						tc["input"] = wrappedTxData.ToString()
					}

					wTo := WrappedAddress{raw: singleTx.To}
					tc["to"] = wTo.ToString()
					tc["originalTo"] = singleTx.To

					wrappedGas := WrappedNumeric{raw: singleTx.GasLimit}
					tc["gas"] = wrappedGas.ToString()
					tc["originalGas"] = singleTx.GasLimit

					wrappedVal := WrappedNumeric{raw: singleTx.Value}
					tc["value"] = wrappedVal.ToString()
					tc["originalValue"] = singleTx.Value

					testCases = append(testCases, tc)
				}

			}

			st["testCases"] = testCases

			simpleTests = append(simpleTests, st)
		}
		testOut, err := json.Marshal(simpleTests)
		if err != nil {
			return err
		}
		fmt.Println(string(testOut))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := RetestCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a file that's filed with test transaction fillers")
	solcPath = os.Getenv("SOLC_PATH")
	if solcPath == "" {
		solcPath = "solc"
	} else {
		log.Info().Str("path", solcPath).Msg("Setting solc path from environment")
	}

	validBase10 = regexp.MustCompile(`^[0-9_]*$`) // the numbers can be formatted like 100_000
	dataLabel = regexp.MustCompile(`^:label ([^ ]*) `)
	typeIndicator = regexp.MustCompile(`^:([^ ]*) `)
	abiSpec = regexp.MustCompile(`^([a-zA-Z0-9]*)\((.*)\)(.*)$`)
	normalizeWs = regexp.MustCompile(` +`)
	solidityCompileInfo = regexp.MustCompile(`^([^\n\r{]*)([\n\r{])`)
	removablePreamble = regexp.MustCompile(`^\b(london|berlin|byzantium|shanghai|optimise)\b(\s+\b(london|berlin|byzantium|shanghai|optimise)\b)*`)
	solcCompileMultiOut = regexp.MustCompile(`======= [^=]*:([^=]*) =======\nBinary:\n([a-fA-F0-9]*)`)
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	return io.ReadAll(os.Stdin)
}
