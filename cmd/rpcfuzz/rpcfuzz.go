// Package rpcfuzz is meant to have some basic RPC fuzzing and
// conformance tests. Each test is meant to be self-contained i.e. the
// success or failure of a test should have no impact on other
// tests. The benefit here is that each test is an object and can be
// modified, decorated, fuzzed, etc.
//
// The conformance test should also run successful on a network that
// is or isn't isolated. In some circumstances, it might be better to
// run the conformance test in once process while there is load being
// applied. The only consideration is that you shouldn't use the same
// key to load test as is used to run the conformance tests.
package rpcfuzz

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha1"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	fuzz "github.com/google/gofuzz"
	"github.com/maticnetwork/polygon-cli/bindings/tester"
	"github.com/maticnetwork/polygon-cli/cmd/rpcfuzz/testreporter"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type (
	// RPCTest is the common interface for a test.  In the future
	// we'll need some addition methods in particular if don't
	// want to run tests that require unlocked accounts or if we
	// want to skip certain namespaces
	RPCTest interface {
		// GetName returns a more descriptive name of the test being executed
		GetName() string

		// GetMethod returns the json rpc method name
		GetMethod() string

		// GetArgs will return the list of arguments that will be used when calling the rpc
		GetArgs() []interface{}

		// Validate will return an error of the result fails validation
		Validate(result interface{}) error

		// ExpectError is used by the validation code to understand of the test typically returns an error
		ExpectError() bool
	}

	// RPCTestFlag is meant for bitmasking various flags to understand properties of the test
	RPCTestFlag uint64

	// RPCTestGeneric is the simplest implementation of the
	// RPCTest. Basically the implementation of the interface is
	// managed by just returning hard coded values for method,
	// args, validator, and error
	RPCTestGeneric struct {
		Name      string
		Method    string
		Args      []interface{}
		Validator func(result interface{}) error
		Flags     RPCTestFlag
	}

	// RPCTestDynamicArgs is a simple implementation of the
	// RPCTest that requires a function for Args which will be
	// used to generate the args for testing.
	RPCTestDynamicArgs struct {
		Name      string
		Method    string
		Args      func() []interface{}
		Validator func(result interface{}) error
		Flags     RPCTestFlag
	}

	// RPCTestTransactionArgs is used to send transactions
	RPCTestTransactionArgs struct {
		From                 string `json:"from,omitempty"`
		To                   string `json:"to,omitempty"`
		Gas                  string `json:"gas,omitempty"`
		GasPrice             string `json:"gasPrice,omitempty"`
		MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
		Value                string `json:"value,omitempty"`
		Nonce                string `json:"nonce,omitempty"`
		Data                 string `json:"data"`
	}

	// RPCTestFilterArgs is a simplified type to contain the flag
	// needed to create a filter
	RPCTestFilterArgs struct {
		FromBlock string        `json:"fromBlock,omitempty"`
		ToBlock   string        `json:"toBlock,omitempty"`
		Address   string        `json:"address,omitempty"`
		Topics    []interface{} `json:"topics,omitempty"`
	}

	// RPCTestRawHTTP is a raw RPCTest performed using HTTP requests.
	// It does not leverage advanced HTTP libraries like `github.com/ethereum/go-ethereum/rpc`.
	RPCTestRawHTTP struct {
		Name       string
		HTTPMethod string
		Args       []interface{}
		Validator  func(result interface{}) error
		Flags      RPCTestFlag
	}

	// RPCJSONError can be used to unmarshal a raw error response.
	RPCJSONError struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	// RPCJSONResponse can be used to unmarshal a raw response.
	RPCJSONResponse struct {
		Version string        `json:"jsonrpc"`
		Result  any           `json:"result,omitempty"`
		Error   *RPCJSONError `json:"error,omitempty"`
		ID      any           `json:"id"`
	}
)

const (
	FlagStrictValidation RPCTestFlag = 1 << iota // strict means the test is unsuitable for fuzzing / mutation because it most likely won't match
	FlagErrorValidation                          // error validation means the result is expected to be an error
	FlagRequiresUnlock                           // unlock means the test depends on unlocked accounts
	FlagEIP1559                                  // tests that would only exist with EIP-1559 enabled
	FlagOrderDependent                           // This flag indicates that the particular test might fail if shuffled

	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"

	defaultGas                  = "0x100000"
	defaultGasPrice             = "0x1000000000"
	defaultMaxFeePerGas         = "0x1000000000"
	defaultMaxPriorityFeePerGas = "0x1000000000"

	defaultNonceTestOffset uint64 = 0x100000000

	rpcTestRawHTTPNamespace = "raw"

	// JSON-RPC error codes.
	// https://eips.ethereum.org/EIPS/eip-1474
	invalidInputErr   = -32000
	parseErr          = -32700
	invalidRequestErr = -32600
	methodNotFoundErr = -32601
	invalidParamsErr  = -32602
	internalErr       = -32603

	// contracts/conformancetester/ConformanceTester.sol
	revertErrorMessage = "Test Revert Error Message"
)

var (
	testPrivateKey        *ecdsa.PrivateKey
	testEthAddress        ethcommon.Address
	testAccountNonce      uint64
	testAccountNonceMutex sync.Mutex
	currentChainID        *big.Int
	fuzzer                *fuzz.Fuzzer
	enabledNamespaces     []string
	// in the future allTests could be used to for
	// fuzzing.. E.g. loop over the various tests, and mutate the
	// Args before sending
	allTests = make([]RPCTest, 0)

	testResults   testreporter.TestResults
	testResultsCh = make(chan testreporter.TestResult)

	fuzzedTestsGroup sync.WaitGroup
	testResultMutex  sync.Mutex
)

func getConformanceContract(ctx context.Context, rpc *rpc.Client, chainID *big.Int) (conformanceContractAddrStr string, conformanceContract *tester.ConformanceTester, err error) {
	log.Trace().Msg("Deploying Conformance contract...")
	var conformanceContractAddr ethcommon.Address
	ec := ethclient.NewClient(rpc)
	tops, err := bind.NewKeyedTransactorWithChainID(testPrivateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return
	}

	cops := new(bind.CallOpts)
	conformanceContractAddr, conformanceContract, err = tester.DeployConformanceContract(ctx, ec, tops, cops)
	conformanceContractAddrStr = conformanceContractAddr.String()
	log.Trace().Msg("Finished Deploying Conformance contract...")

	return
}

func runRpcFuzz(ctx context.Context) error {
	if *testOutputExportPath != "" && !*testExportJson && !*testExportCSV && !*testExportMarkdown && !*testExportHTML {
		log.Warn().Msg("Setting --export-path must pair with a export type: --json, --csv, --md, or --html")
	}

	rpcClient, err := rpc.DialContext(ctx, *rpcUrl)
	if err != nil {
		return err
	}
	chainId, err := GetCurrentChainID(ctx, rpcClient)
	if err != nil {
		return err
	}
	currentChainID = chainId

	if *testContractAddress == "" {
		conformanceContractAddr, _, deploymentErr := getConformanceContract(ctx, rpcClient, currentChainID)
		if deploymentErr != nil {
			log.Error().Err(deploymentErr).Msg("Load test contract deployment error")
			return deploymentErr
		}
		testContractAddress = &conformanceContractAddr
	}
	log.Info().Str("testContractAddress", *testContractAddress).Msg("Conformance test contract address")

	nonce, err := GetTestAccountNonce(ctx, rpcClient)
	if err != nil {
		return err
	}
	testAccountNonce = nonce

	log.Trace().Uint64("nonce", nonce).Uint64("chainId", chainId.Uint64()).Msg("Doing test setup")
	setupTests(ctx, rpcClient)

	httpClient := &http.Client{}
	wrappedHTTPClient := wrappedHttpClient{httpClient, *rpcUrl}

	for _, t := range allTests {
		if !shouldRunTest(t) {
			log.Trace().Str("name", t.GetName()).Str("method", t.GetMethod()).Msg("Skipping test")
			continue
		}
		log.Trace().Str("name", t.GetName()).Str("method", t.GetMethod()).Msg("Running Test")

		currTestResult := CallRPCAndValidate(ctx, rpcClient, wrappedHTTPClient, t)
		testResults.AddTestResult(currTestResult)

		if *testFuzz {
			fuzzedTestsGroup.Add(1)

			log.Info().Str("method", t.GetMethod()).Msg("Running with fuzzed args")
			go func(t RPCTest) {
				defer fuzzedTestsGroup.Done()
				currTestResult := CallRPCWithFuzzAndValidate(ctx, rpcClient, t)
				testResultsCh <- currTestResult
			}(t)
		}
	}

	go func() {
		for currTestResult := range testResultsCh {
			testResultMutex.Lock()
			testResults.AddTestResult(currTestResult)
			testResultMutex.Unlock()
		}
	}()

	fuzzedTestsGroup.Wait()
	close(testResultsCh)

	testResults.GenerateTabularResult()
	if *testExportJson {
		testResults.ExportResultToJSON(filepath.Join(*testOutputExportPath, "output.json"))
	}
	if *testExportCSV {
		testResults.ExportResultToCSV(filepath.Join(*testOutputExportPath, "output.csv"))
	}
	if *testExportMarkdown {
		testResults.ExportResultToMarkdown(filepath.Join(*testOutputExportPath, "output.md"))
	}
	if *testExportHTML {
		testResults.ExportResultToHTML(filepath.Join(*testOutputExportPath, "output.html"))
	}
	testResults.PrintTabularResult()

	return nil
}

// setupTests will add all of the `RPCTests` to the `allTests` slice.
func setupTests(ctx context.Context, rpcClient *rpc.Client) {
	// cast rpc --rpc-url localhost:8545 net_version
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestNetVersion",
		Method:    "net_version",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^\d*$`),
	})

	// cast rpc --rpc-url localhost:8545 web3_clientVersion
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestWeb3ClientVersion",
		Method:    "web3_clientVersion",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^[[:print:]]*$`),
	})

	// cast rpc --rpc-url localhost:8545 web3_sha3 0x68656c6c6f20776f726c64
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestWeb3SHA3",
		Method:    "web3_sha3",
		Flags:     FlagStrictValidation,
		Args:      []interface{}{"0x68656c6c6f20776f726c64"},
		Validator: ValidateRegexString(`0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad`),
	})

	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestWeb3SHA3Error",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "web3_sha3",
		Args:      []interface{}{"68656c6c6f20776f726c64"},
		Validator: ValidateError(invalidParamsErr, `cannot unmarshal hex string without 0x prefix`),
	})

	// cast rpc --rpc-url localhost:8545 net_listening
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestNetListening",
		Method:    "net_listening",
		Args:      []interface{}{},
		Validator: ValidateExact(true),
	})

	// cast rpc --rpc-url localhost:8545 net_peerCount
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestNetPeerCount",
		Method:    "net_peerCount",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]*$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_protocolVersion
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthProtocolVersion",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "eth_protocolVersion",
		Args:      []interface{}{},
		Validator: ValidateError(methodNotFoundErr, `method eth_protocolVersion does not exist`),
	})

	// cast rpc --rpc-url localhost:8545 eth_syncing
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthSyncing",
		Method:    "eth_syncing",
		Args:      []interface{}{},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthSyncing),
	})

	// cast rpc --rpc-url localhost:8545 eth_coinbase
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCoinbase",
		Method:    "eth_coinbase",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{40}$`),
		Flags:     FlagRequiresUnlock,
	})

	// cast rpc --rpc-url localhost:8545 eth_chainId
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthChainID",
		Method:    "eth_chainId",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_mining
	allTests = append(allTests, &RPCTestGeneric{
		Name:   "RPCTestEthMining",
		Method: "eth_mining",
		Args:   []interface{}{},
		Validator: RequireAny(
			ValidateExact(true),
			ValidateExact(false),
		),
	})

	// cast rpc --rpc-url localhost:8545 eth_hashrate
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthHashrate",
		Method:    "eth_hashrate",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthHashrateNoETHash",
		Method:    "eth_hashrate",
		Flags:     FlagStrictValidation,
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x0$`), // A strict check for 0x0 is expected in EVM now ethhash has been yeeted
	})

	// cast rpc --rpc-url localhost:8545 eth_gasPrice
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGasPrice",
		Method:    "eth_gasPrice",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_accounts
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthAccounts",
		Method:    "eth_accounts",
		Args:      []interface{}{},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaAccountList),
		Flags:     FlagRequiresUnlock,
	})

	// cast rpc --rpc-url localhost:8545 eth_blockNumber
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthBlockNumber",
		Method:    "eth_blockNumber",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast balance --rpc-url localhost:8545 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceLatest",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "latest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceEarliest",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "earliest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBalancePending",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "pending"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceZero",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "0x0"},
		Flags:     FlagStrictValidation,
		Validator: ValidateRegexString(`^0x0$`),
	})

	// cast storage --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429 0 --block latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtLatest",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x0", "latest"},
		Flags:     FlagStrictValidation,
		Validator: ValidateRegexString(`^0x436f6e666f726d616e6365546573746572436f6e74726163744e616d6500003a$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtEarliest",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x0", "earliest"},
		Validator: ValidateRegexString(`^0x0{64}`),
	})
	// cast storage --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429 0 --block pending
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtPending",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x0", "pending"},
		Flags:     FlagStrictValidation,
		Validator: ValidateRegexString(`^0x436f6e666f726d616e6365546573746572436f6e74726163744e616d6500003a$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtZero",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x0", "0x0"},
		Flags:     FlagStrictValidation,
		Validator: ValidateRegexString(`^0x0{64}`),
	})

	// cast rpc --rpc-url localhost:8545 eth_getTransactionCount 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtLatest",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "latest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtEarliest",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "earliest"},
		Validator: ValidateRegexString(`^0x0$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtPending",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "pending"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtZero",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "0x0"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_getBlockTransactionCountByHash 0x9300b64619e167e7dbc1b41a6a6e7a8de7d6b99427dceefbd58014e328bd7f92
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthGetBlockTransactionCountByHash",
		Method:    "eth_getBlockTransactionCountByHash",
		Args:      ArgsLatestBlockHash(ctx, rpcClient),
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByHashMissing",
		Method:    "eth_getBlockTransactionCountByHash",
		Flags:     FlagStrictValidation,
		Args:      []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Validator: ValidateExact(nil),
	})

	// cast rpc --rpc-url localhost:8545 eth_getBlockTransactionCountByNumber 0x1
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberLatest",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberEarliest",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberPending",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"pending"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberZero",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_getUncleCountByBlockHash 0x9300b64619e167e7dbc1b41a6a6e7a8de7d6b99427dceefbd58014e328bd7f92
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthGetUncleCountByBlockHash",
		Method:    "eth_getUncleCountByBlockHash",
		Args:      ArgsLatestBlockHash(ctx, rpcClient),
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockHashMissing",
		Method:    "eth_getUncleCountByBlockHash",
		Flags:     FlagStrictValidation,
		Args:      []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Validator: ValidateExact(nil),
	})

	// cast rpc --rpc-url localhost:8545 eth_getUncleCountByBlockNumber 0x1
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberLatest",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberEarliest",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberPending",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"pending"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberZero",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// curl http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"method":"eth_getCode","params":["0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429","latest"],"id":48,"jsonrpc":"2.0"}' | jq -r ".result | tojson" | tr -d '\n' | sha1sum
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetCodeLatest",
		Method:    "eth_getCode",
		Args:      []interface{}{*testContractAddress, "latest"},
		Flags:     FlagStrictValidation,
		Validator: ValidateHashedResponse("b9689c32bf9284029715ff8375f8996129898db9"),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetCodePending",
		Method:    "eth_getCode",
		Args:      []interface{}{*testContractAddress, "pending"},
		Flags:     FlagStrictValidation,
		Validator: ValidateHashedResponse("b9689c32bf9284029715ff8375f8996129898db9"),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetCodeEarliest",
		Method:    "eth_getCode",
		Flags:     FlagStrictValidation,
		Args:      []interface{}{*testContractAddress, "earliest"},
		Validator: ValidateRegexString(`^0x$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_sign "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57" "0xdeadbeaf"
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSign",
		Method:    "eth_sign",
		Args:      ArgsCoinbase(ctx, rpcClient, "0xdeadbeef"),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{130}$`),
		Flags:     FlagRequiresUnlock,
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthSignFail",
		Method:    "eth_sign",
		Args:      []interface{}{testEthAddress.String(), "0xdeadbeef"},
		Validator: ValidateError(invalidInputErr, `unknown account`),
		Flags:     FlagErrorValidation | FlagStrictValidation | FlagRequiresUnlock,
	})

	// cast rpc --rpc-url localhost:8545 eth_signTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	// Weirdly, the response eth_signTransaction from geth doesn't "conform" to the spec of: https://ethereum.github.io/execution-apis/api-documentation/
	// This is the actual response with cast and curl of the above input:
	// {"raw":"0xf8...85","tx":{"type":"0x0","chainId":"0x539","nonce":"0x1","to":"0x85da99c8a7c2c95964c8efd687e95e632fc533d6","gas":"0x5208","gasPrice":"0x1","maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x0","input":"0x","v":"0xa95","r":"0x82..fe","s":"0x78..85","hash":"0xa9b..5a"}}
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSignTransaction",
		Method:    "eth_signTransaction",
		Args:      ArgsCoinbaseTransaction(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas, Nonce: "0x1"}),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaSignTxResponse),
		Flags:     FlagRequiresUnlock,
	})

	// cast rpc --rpc-url localhost:8545 eth_sendTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendTransaction",
		Method:    "eth_sendTransaction",
		Args:      ArgsCoinbaseTransaction(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{64}$`),
		Flags:     FlagRequiresUnlock,
	})

	// $ curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[],"id":1}' http://localhost:8545
	// {"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"missing value for required argument 0"}}
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthSendTransactionEmpty",
		Method:    "eth_sendTransaction",
		Args:      []interface{}{},
		Flags:     FlagErrorValidation,
		Validator: ValidateError(invalidParamsErr, "missing value for required argument 0"),
	})
	// $ curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[],"id":1}' http://localhost:8545
	// {"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"missing value for required argument 0"}}
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthSendRawTransactionEmpty",
		Method:    "eth_sendRawTransaction",
		Args:      []interface{}{},
		Flags:     FlagErrorValidation,
		Validator: ValidateError(invalidParamsErr, "missing value for required argument 0"),
	})

	// cast rpc --rpc-url localhost:8545 eth_sendRawTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransaction",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransaction(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{64}$`),
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransactionNonceTooLow",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransactionWithNonce(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}, 0),
		Validator: ValidateError(invalidInputErr, `nonce too low`),
		Flags:     FlagErrorValidation | FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransactionNonceHigh",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransactionWithNonce(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}, testAccountNonce|defaultNonceTestOffset),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{64}$`),
		Flags:     FlagOrderDependent,
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransactionNonceKnown",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransactionWithNonce(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}, testAccountNonce|defaultNonceTestOffset),
		Validator: ValidateError(invalidInputErr, `already known`),
		Flags:     FlagErrorValidation | FlagStrictValidation | FlagOrderDependent,
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransactionNonceUnderpriced",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransactionWithNonce(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x1234", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}, testAccountNonce|defaultNonceTestOffset),
		Validator: ValidateError(invalidInputErr, `replacement`),
		Flags:     FlagErrorValidation | FlagStrictValidation | FlagOrderDependent,
	})

	// curl http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"method":"eth_call","params":[{"to":"0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429","data":"0x06fdde03"},"latest"],"id":1,"jsonrpc":"2.0"}'
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCallLatest",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "latest"},
		Validator: ValidateRegexString(`0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001d436f6e666f726d616e6365546573746572436f6e74726163744e616d65000000`),
		Flags:     FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCallPending",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "pending"},
		Validator: ValidateRegexString(`0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001d436f6e666f726d616e6365546573746572436f6e74726163744e616d65000000`),
		Flags:     FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCallEarliest",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "earliest"},
		Validator: ValidateRegexString(`^0x$`),
		Flags:     FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCallZero",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "0x0"},
		Validator: ValidateRegexString(`^0x$`),
		Flags:     FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCallRevertMessage",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Data: "0xa26388bb"}, "latest"},
		Validator: ValidateErrorMsgString(revertErrorMessage),
		Flags:     FlagErrorValidation,
	})

	// curl http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"method":"eth_estimateGas","params":[{"to":"0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429","data":"0x06fdde03"},"latest"],"id":1,"jsonrpc":"2.0"}'
	// This just validates that the estimated gas response conforms to an expected hex
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthEstimateGas",
		Method:    "eth_estimateGas",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "latest"},
		Validator: ValidateRegexString(`^0x`),
		Flags:     FlagStrictValidation,
	})

	// $ curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc": "2.0", "method": "eth_estimateGas", "params": [], "id":1}' localhost:8545
	// {"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"missing value for required argument 0"}}
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthEstimateGasEmpty",
		Method:    "eth_estimateGas",
		Args:      []interface{}{},
		Flags:     FlagErrorValidation,
		Validator: ValidateError(invalidParamsErr, `missing value for required argument 0`),
	})

	// cast block --rpc-url localhost:8545 latest
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestEthGetBlockByHash",
		Method: "eth_getBlockByHash",
		Args:   ArgsLatestBlockHash(ctx, rpcClient, true),
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
			ValidateBlockHash(),
		),
	})

	allTests = append(allTests, &RPCTestGeneric{
		Name:   "RPCTestEthGetBlockByHashZero",
		Method: "eth_getBlockByHash",
		Args: []interface{}{
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			true,
		},
		Validator: ValidateExact(nil),
	})

	// cast block --rpc-url localhost:8545 0
	allTests = append(allTests, &RPCTestGeneric{
		Name:   "RPCTestEthBlockByNumber",
		Method: "eth_getBlockByNumber",
		Args:   []interface{}{"0x0", true},
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
			ValidateBlockHash(),
		),
	})

	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestEthBlockByNumberLatest",
		Method: "eth_getBlockByNumber",
		Args:   ArgsLatestBlockNumber(ctx, rpcClient, true),
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
			ValidateBlockHash(),
		),
	})

	// cast rpc --rpc-url localhost:8545 eth_getTransactionByHash 0xb27bd60d706c08a80d698b951b9ec4284b342a34b885ff5ebe567b41dab16f69
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestEthGetTransactionByHash",
		Method: "eth_getTransactionByHash",
		Args:   ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthTransaction),
			ValidateTransactionHash(),
		),
	})

	// cast rpc --rpc-url localhost:8545 eth_getTransactionByBlockHashAndIndex 0x63f86797e33513449350d0e00ef962f172a94a60b990a096a470c1ac1df5ec06 0x0
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestEthGetTransactionByBlockHashAndIndex",
		Method: "eth_getTransactionByBlockHashAndIndex",
		Args:   ArgsTransactionBlockHashAndIndex(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthTransaction),
			ValidateTransactionHash(),
		),
	})

	// cast rpc --rpc-url localhost:8545 eth_getTransactionByBlockNumberAndIndex 0xd 0x0
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestEthGetTransactionByBlockNumberAndIndex",
		Method: "eth_getTransactionByBlockNumberAndIndex",
		Args:   ArgsTransactionBlockNumberAndIndex(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: RequireAll(
			ValidateJSONSchema(rpctypes.RPCSchemaEthTransaction),
			ValidateTransactionHash(),
		),
	})

	// cast receipt --rpc-url localhost:8545 0x1bd4ec642302aa22906360af6493c230ecc41df10fffcdedc85caeb22cbb6b58
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestGetTransactionReceipt",
		Method:    "eth_getTransactionReceipt",
		Args:      ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthReceipt),
	})

	// This RPC can be validated pretty easily, but it's not clear how to create an uncle in a reproducible away in order to test this method reliably
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestGetUncleByBlockHashAndIndex",
		Method:    "eth_getUncleByBlockHashAndIndex",
		Args:      ArgsLatestBlockHash(ctx, rpcClient, "0x0"),
		Validator: RequireAny(ValidateJSONSchema(rpctypes.RPCSchemaEthBlock), ValidateExact(nil)),
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestGetUncleByBlockNumberAndIndex",
		Method:    "eth_getUncleByBlockNumberAndIndex",
		Args:      ArgsLatestBlockNumber(ctx, rpcClient, "0x0"),
		Validator: RequireAny(ValidateJSONSchema(rpctypes.RPCSchemaEthBlock), ValidateExact(nil)),
	})

	// cast rpc --rpc-url localhost:8545 eth_getCompilers
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetCompilers",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "eth_getCompilers",
		Args:      []interface{}{},
		Validator: ValidateError(methodNotFoundErr, `method eth_getCompilers does not exist`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCompileSolidity",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "eth_compileSolidity",
		Args:      []interface{}{},
		Validator: ValidateError(methodNotFoundErr, `method eth_compileSolidity does not exist`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCompileLLL",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "eth_compileLLL",
		Args:      []interface{}{},
		Validator: ValidateError(methodNotFoundErr, `method eth_compileLLL does not exist`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthCompileSerpent",
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Method:    "eth_compileSerpent",
		Args:      []interface{}{},
		Validator: ValidateError(methodNotFoundErr, `method eth_compileSerpent does not exist`),
	})

	// cast rpc --rpc-url localhost:8545 eth_newFilter "{}"
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewFilterEmpty",
		Method:    "eth_newFilter",
		Args:      []interface{}{RPCTestFilterArgs{}},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewFilterFromOnly",
		Method:    "eth_newFilter",
		Args:      []interface{}{RPCTestFilterArgs{FromBlock: "earliest"}},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewFilterToOnly",
		Method:    "eth_newFilter",
		Args:      []interface{}{RPCTestFilterArgs{ToBlock: "latest"}},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewFilterAddressOnly",
		Method:    "eth_newFilter",
		Args:      []interface{}{RPCTestFilterArgs{Address: *testContractAddress}},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewFilterTopicsOnly",
		Method:    "eth_newFilter",
		Args:      []interface{}{RPCTestFilterArgs{Topics: []interface{}{nil, nil, "0x000000000000000000000000" + testEthAddress.String()[2:]}}},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:   "RPCTestEthNewFilterAllFields",
		Method: "eth_newFilter",
		Args: []interface{}{RPCTestFilterArgs{
			FromBlock: "earliest",
			ToBlock:   "latest",
			Address:   *testContractAddress,
			Topics:    []interface{}{nil, nil, "0x000000000000000000000000" + testEthAddress.String()[2:]}},
		},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_newBlockFilter
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewBlockFilter",
		Method:    "eth_newBlockFilter",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_newPendingTransactionFilter
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthNewPendingTransactionFilter",
		Method:    "eth_newPendingTransactionFilter",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// cast rpc --rpc-url localhost:8545 eth_uninstallFilter 0x842bc0d4f68eba291ed5c00ef04541d3
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthUninstallFilterFail",
		Method:    "eth_uninstallFilter",
		Args:      []interface{}{"0xdeadbeef"},
		Validator: ValidateExact(false),
		Flags:     FlagStrictValidation,
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthUninstallFilterSucceed",
		Method:    "eth_uninstallFilter",
		Args:      ArgsBlockFilterID(ctx, rpcClient),
		Validator: ValidateExact(true),
		Flags:     FlagStrictValidation,
	})

	// cast rpc --rpc-url localhost:8545 eth_newFilter '{"fromBlock": "earliest", "toBlock": "latest", "address": "0x6fda56c57b0acadb96ed5624ac500c0429d59429", "topics": [null, null, "0x00000000000000000000000085da99c8a7c2c95964c8efd687e95e632fc533d6"]}'
	// cast rpc --rpc-url localhost:8545 eth_getFilterChanges 0xef69a30e77c9902dec23745e0bbe4586
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestGetFilterChanges",
		Method: "eth_getFilterChanges",
		Args: ArgsFilterID(ctx, rpcClient, RPCTestFilterArgs{
			FromBlock: "earliest",
			ToBlock:   "latest",
			Address:   *testContractAddress,
			Topics:    []interface{}{nil, nil, "0x000000000000000000000000" + testEthAddress.String()[2:]},
		}),
		Validator: RequireAny(
			ValidateJSONSchema(rpctypes.RPCSchemaEthFilter),
			ValidateExactJSON("[]"),
		),
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:   "RPCTestGetFilterLogs",
		Method: "eth_getFilterLogs",
		Args: ArgsFilterID(ctx, rpcClient, RPCTestFilterArgs{
			FromBlock: "earliest",
			ToBlock:   "latest",
			Address:   *testContractAddress,
			Topics:    []interface{}{nil, nil, "0x000000000000000000000000" + testEthAddress.String()[2:]},
		}),
		Validator: RequireAny(
			ValidateJSONSchema(rpctypes.RPCSchemaEthFilter),
			ValidateExactJSON("[]"),
		),
	})
	// cast rpc --rpc-url localhost:8545 eth_getLogs '{"fromBlock": "earliest", "toBlock": "latest", "address": "0x6fda56c57b0acadb96ed5624ac500c0429d59429", "topics": [null, null, "0x00000000000000000000000085da99c8a7c2c95964c8efd687e95e632fc533d6"]}'
	allTests = append(allTests, &RPCTestGeneric{
		Name:   "RPCTestGetLogs",
		Method: "eth_getLogs",
		Args: []interface{}{RPCTestFilterArgs{
			FromBlock: "earliest",
			ToBlock:   "latest",
			Address:   *testContractAddress,
			Topics:    []interface{}{nil, nil, "0x000000000000000000000000" + testEthAddress.String()[2:]},
		}},
		Validator: RequireAny(
			ValidateJSONSchema(rpctypes.RPCSchemaEthFilter),
			ValidateExactJSON("[]"),
		),
	})

	// cast rpc --rpc-url localhost:8545 eth_getWork
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestGetWork",
		Method:    "eth_getWork",
		Args:      []interface{}{},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(methodNotFoundErr, `method eth_getWork does not exist`),
	})
	// cast rpc --rpc-url localhost:8545 eth_submitWork 0x0011223344556677 0x00112233445566778899AABBCCDDEEFF 0x00112233445566778899AABBCCDDEEFF
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestSubmitWork",
		Method:    "eth_submitWork",
		Args:      []interface{}{"0x0011223344556677", "0x00112233445566778899AABBCCDDEEFF", "0x00112233445566778899AABBCCDDEEFF"},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(methodNotFoundErr, `method eth_submitWork does not exist`),
	})
	// cast rpc --rpc-url localhost:8545 eth_submitHashrate 0x00112233445566778899AABBCCDDEEFF00112233445566778899AABBCCDDEEFF 0x00112233445566778899AABBCCDDEEFF00112233445566778899AABBCCDDEEFF
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestSubmitHashrate",
		Method:    "eth_submitHashrate",
		Args:      []interface{}{"0x00112233445566778899AABBCCDDEEFF00112233445566778899AABBCCDDEEFF", "0x00112233445566778899AABBCCDDEEFF00112233445566778899AABBCCDDEEFF"},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(methodNotFoundErr, `method eth_submitHashrate does not exist`),
	})

	// cast rpc --rpc-url localhost:8545 eth_feeHistory 128 latest []
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestFeeHistory",
		Method:    "eth_feeHistory",
		Args:      []interface{}{"0xF", "latest", nil},
		Flags:     FlagEIP1559,
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthFeeHistory),
	})

	// cast rpc --rpc-url localhost:8545 eth_maxPriorityFeePerGas
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestMaxPriorityFeePerGas",
		Method:    "eth_maxPriorityFeePerGas",
		Args:      []interface{}{},
		Flags:     FlagEIP1559,
		Validator: ValidateRegexString(`^0x([1-9a-f]+[0-9a-f]*|0)$`),
	})

	// don't have a handy cast call for eth_createAccessList
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestEthCreateAccessList",
		Method:    "eth_createAccessList",
		Args:      ArgsCoinbaseTransaction(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}, "latest"),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthAccessList),
		Flags:     FlagRequiresUnlock,
	})

	// cast rpc --rpc-url localhost:8545 eth_getProof 0x6fda56c57b0acadb96ed5624ac500c0429d59429 '["0x3"]' latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestEthGetProof",
		Method:    "eth_getProof",
		Args:      []interface{}{*testContractAddress, []interface{}{"0x3"}, "latest"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthProof),
	})

	// curl http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"method":"eth_call","params":[{"to":"0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429","data":"0x06fdde03"},"latest"],"id":1,"jsonrpc":"2.0"}'
	// cast rpc --rpc-url localhost:8545 debug_traceCall '{"to": "0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429", "data":"0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001d436f6e666f726d616e6365546573746572436f6e74726163744e616d65000000"}' latest | jq '.'
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceCallSimple",
		Method:    "debug_traceCall",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "latest"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTrace),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceCallName",
		Method:    "debug_traceCall",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "latest"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTrace),
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugTraceTransactionSimple",
		Method:    "debug_traceTransaction",
		Args:      ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas, Gas: defaultGas}),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTrace),
	})
	// cast calldata "deposit(uint256)" 1
	// 0xb6b55f250000000000000000000000000000000000000000000000000000000000000001
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugTraceTransactionDeposit",
		Method:    "debug_traceTransaction",
		Args:      ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0xb6b55f250000000000000000000000000000000000000000000000000000000000000001", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas, Gas: defaultGas}),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTrace),
	})

	// cast rpc --rpc-url localhost:8545 debug_getRawBlock latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawBlockLatest",
		Method:    "debug_getRawBlock",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawBlockPending",
		Method:    "debug_getRawBlock",
		Args:      []interface{}{"pending"},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(invalidInputErr, `not found`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawBlockEarliest",
		Method:    "debug_getRawBlock",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawBlockZero",
		Method:    "debug_getRawBlock",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})

	// cast rpc --rpc-url localhost:8545 debug_getBadBlocks
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetBadBlocks",
		Method:    "debug_getBadBlocks",
		Args:      []interface{}{},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaBadBlocks),
	})

	// cast rpc --rpc-url localhost:8545 debug_getRawHeader latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawHeaderLatest",
		Method:    "debug_getRawHeader",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawHeaderPending",
		Method:    "debug_getRawHeader",
		Args:      []interface{}{"pending"},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(invalidInputErr, `not found`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawHeaderEarliest",
		Method:    "debug_getRawHeader",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawHeaderZero",
		Method:    "debug_getRawHeader",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})

	// cast rpc --rpc-url localhost:8545 debug_getRawReceipts latest
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawReceiptsLatest",
		Method:    "debug_getRawReceipts",
		Args:      []interface{}{"latest"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaHexArray),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawReceiptsPending",
		Method:    "debug_getRawReceipts",
		Args:      []interface{}{"pending"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaHexArray),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawReceiptsEarliest",
		Method:    "debug_getRawReceipts",
		Args:      []interface{}{"earliest"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaHexArray),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugGetRawReceiptsZero",
		Method:    "debug_getRawReceipts",
		Args:      []interface{}{"0x0"},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaHexArray),
	})

	// get raw tx
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugGetRawTransactionSimple",
		Method:    "debug_getRawTransaction",
		Args:      ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas}),
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})
	// cast calldata "deposit(uint256)" 1
	// 0xb6b55f250000000000000000000000000000000000000000000000000000000000000001
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugGetRawTransactionDeposit",
		Method:    "debug_getRawTransaction",
		Args:      ArgsTransactionHash(ctx, rpcClient, &RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0xb6b55f250000000000000000000000000000000000000000000000000000000000000001", MaxFeePerGas: defaultMaxFeePerGas, MaxPriorityFeePerGas: defaultMaxPriorityFeePerGas, Gas: defaultGas}),
		Validator: ValidateRegexString(`^0x[0-9a-f]*`),
	})

	// cast rpc --rpc-url localhost:8545 debug_traceBlockByNumber 0x0
	// cast rpc --rpc-url localhost:8545 debug_traceBlockByNumber 0x1
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceBlockByNumberZero",
		Method:    "debug_traceBlockByNumber",
		Args:      []interface{}{"0x0", nil},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(invalidInputErr, `genesis is not traceable`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceBlockByNumberOne",
		Method:    "debug_traceBlockByNumber",
		Args:      []interface{}{"0x1", nil},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTraceBlock),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceBlockByNumberLatest",
		Method:    "debug_traceBlockByNumber",
		Args:      []interface{}{"latest", nil},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTraceBlock),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceBlockByNumberEarliest",
		Method:    "debug_traceBlockByNumber",
		Args:      []interface{}{"earliest", nil},
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(invalidInputErr, `genesis is not traceable`),
	})
	allTests = append(allTests, &RPCTestGeneric{
		Name:      "RPCTestDebugTraceBlockByNumberPending",
		Method:    "debug_traceBlockByNumber",
		Args:      []interface{}{"pending", nil},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTraceBlock),
	})

	// debug_traceBlockByHash
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugTraceBlockByHash",
		Method:    "debug_traceBlockByHash",
		Args:      ArgsLatestBlockHash(ctx, rpcClient, nil),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTraceBlock),
	})

	// debug_traceBlock
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugTraceBlock",
		Method:    "debug_traceBlock",
		Args:      ArgsRawBlock(ctx, rpcClient, "latest", nil),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaDebugTraceBlock),
	})
	allTests = append(allTests, &RPCTestDynamicArgs{
		Name:      "RPCTestDebugTraceBlockZero",
		Method:    "debug_traceBlock",
		Args:      ArgsRawBlock(ctx, rpcClient, "0x0", nil),
		Flags:     FlagErrorValidation | FlagStrictValidation,
		Validator: ValidateError(invalidInputErr, `genesis is not traceable`),
	})

	// $ curl -X POST -H "Content-Type: application/json" --data '[]' http://localhost:8545
	// {"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"empty batch"}}
	allTests = append(allTests, &RPCTestRawHTTP{
		Name:       "EmptyBatch",
		HTTPMethod: http.MethodPost,
		Args:       []interface{}{},
		Validator:  ValidateError(invalidRequestErr, "empty batch"),
	})

	uniqueTests := make(map[RPCTest]struct{})
	uniqueTestNames := make(map[string]struct{})
	for _, v := range allTests {
		_, hasKey := uniqueTests[v]
		if hasKey {
			log.Fatal().Str("name", v.GetName()).Str("method", v.GetMethod()).Msg("Duplicate test detected")
		}
		uniqueTests[v] = struct{}{}
		_, hasKey = uniqueTestNames[v.GetName()]
		if hasKey {
			log.Fatal().Str("name", v.GetName()).Str("method", v.GetMethod()).Msg("Duplicate test name detected")
		}
		uniqueTestNames[v.GetName()] = struct{}{}
	}

}

func RequireAny(validators ...func(interface{}) error) func(result interface{}) error {
	return func(result interface{}) error {
		for _, v := range validators {
			err := v(result)
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("all Validation failed")
	}
}
func RequireAll(validators ...func(interface{}) error) func(result interface{}) error {
	return func(result interface{}) error {
		for _, v := range validators {
			err := v(result)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ValidateHashedResponse will take a hex encoded hash and return a
// function that will validate that a given result has the same
// hash. The expected hash does not start with 0x
func ValidateHashedResponse(expectedHash string) func(result interface{}) error {
	return func(result interface{}) error {
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("unable to marshal result object to json %w", err)
		}
		actualHash := fmt.Sprintf("%x", sha1.Sum(jsonBytes))
		if actualHash != expectedHash {
			return fmt.Errorf("hash mismatch expected: %s and got %s", expectedHash, actualHash)
		}
		return nil
	}
}

// ValidateJSONSchema is used to validate the response against a JSON Schema
func ValidateJSONSchema(schema string) func(result interface{}) error {
	return func(result interface{}) error {
		validatorLoader := gojsonschema.NewStringLoader(schema)

		// This is weird, but the current setup doesn't allow
		// for easy access to the initial response string...
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("unable to marshal result back to json for validation: %w", err)
		}
		responseLoader := gojsonschema.NewStringLoader(string(jsonBytes))

		validatorResult, err := gojsonschema.Validate(validatorLoader, responseLoader)
		if err != nil {
			return fmt.Errorf("unable to run json validation: %w", err)
		}
		// fmt.Println(string(jsonBytes))
		if !validatorResult.Valid() {
			errStr := ""
			for _, desc := range validatorResult.Errors() {
				errStr += desc.String() + "\n"
			}
			log.Trace().Str("resultJson", string(jsonBytes)).Msg("JSON failed to validate")
			return fmt.Errorf("the json document is not valid: %s", errStr)
		}
		return nil

	}
}

// ValidateExact will validate against the exact value expected.
func ValidateExact(expected interface{}) func(result interface{}) error {
	return func(result interface{}) error {
		if expected != result {
			return fmt.Errorf("expected %v and got %v", expected, result)
		}
		return nil
	}
}
func ValidateExactJSON(expected string) func(result interface{}) error {
	return func(result interface{}) error {
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("unable to json marshal test result: %w", err)
		}

		if expected != string(jsonResult) {
			return fmt.Errorf("expected %v and got %v", expected, string(jsonResult))
		}
		return nil
	}
}

// ValidateRegexString will match a string from the json response against a regular expression
func ValidateRegexString(regEx string) func(result interface{}) error {
	r := regexp.MustCompile(regEx)
	return func(result interface{}) error {
		resultStr, isValid := result.(string)
		if !isValid {
			return fmt.Errorf("invalid result type. Expected string but got %T", result)
		}
		if !r.MatchString(resultStr) {
			return fmt.Errorf("the regex %s failed to match result %s", regEx, resultStr)
		}
		return nil
	}
}

// ValidateErrorMsgString will check the error message text against the provide regular expression
func ValidateErrorMsgString(errorMessageRegex string) func(result interface{}) error {
	r := regexp.MustCompile(errorMessageRegex)
	return func(result interface{}) error {
		fullError, err := genericResultToError(result)
		if err != nil {
			return err
		}
		if !r.MatchString(fullError.Error()) {
			return fmt.Errorf("the regex %s failed to match result %s", errorMessageRegex, fullError.Error())
		}

		return nil
	}
}

// ValidateError will check the status code and error message text against the provide regular expression
func ValidateError(code int, errorMessageRegex string) func(result interface{}) error {
	r := regexp.MustCompile(errorMessageRegex)
	return func(result interface{}) error {
		fullError, err := genericResultToError(result)
		if err != nil {
			return err
		}
		if !r.MatchString(fullError.Error()) {
			return fmt.Errorf("the regex %s failed to match result %s", errorMessageRegex, fullError.Error())
		}
		if code != fullError.Code {
			return fmt.Errorf("expected error code %d but got %d", code, fullError.Code)
		}

		return nil
	}
}

// ValidateBlockHash will convert the result into a block and compute
// the header in order to verify that the rpc header matches the
// computed header.
func ValidateBlockHash() func(result interface{}) error {
	return func(result interface{}) error {

		blockHeader, genericHash, err := genericResultToBlockHeader(result)
		if err != nil {
			return err
		}
		if blockHeader.Hash().String() != genericHash {
			return fmt.Errorf("block hash mismatch. Computed %s and got %s in the json rpc response", blockHeader.Hash().String(), genericHash)
		}
		return nil
	}
}

// ValidateTransactionHash will compare the rpc transaction hash to the computed transaction hash
func ValidateTransactionHash() func(result interface{}) error {
	return func(result interface{}) error {

		tx, genericHash, err := genericResultToTransaction(result)
		if err != nil {
			return err
		}
		if tx.Hash().String() != genericHash {
			return fmt.Errorf("tx hash mismatch. Computed %s and got %s in the json rpc response", tx.Hash().String(), genericHash)
		}
		return nil
	}
}

func genericResultToBlockHeader(result interface{}) (*ethtypes.Header, string, error) {
	underlyingBlock, ok := result.(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("the underlying type of the result didn't match a block header. Got %T", result)
	}
	genericHash, ok := underlyingBlock["hash"].(string)
	if !ok {
		return nil, "", fmt.Errorf("could not recover the underlying hash. Expected a string and got %T", result)
	}
	log.Info().Str("blockHash", genericHash).Msg("Original block hash")
	jsonBlock, err := json.Marshal(underlyingBlock)
	if err != nil {
		return nil, "", fmt.Errorf("could not json marshal initial block result %w", err)
	}

	blockHeader := ethtypes.Header{}

	err = blockHeader.UnmarshalJSON(jsonBlock)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal json block to geth based json block: %w", err)
	}
	return &blockHeader, genericHash, nil
}
func genericResultToTransaction(result interface{}) (*ethtypes.Transaction, string, error) {
	underlyingTx, ok := result.(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("the underlying type of the result didn't match a transaction. Got %T", result)
	}
	genericHash, ok := underlyingTx["hash"].(string)
	if !ok {
		return nil, "", fmt.Errorf("could not recover the underlying hash. Expected a string and got %T", result)
	}
	log.Info().Str("txHash", genericHash).Msg("Original tx hash")
	jsonTx, err := json.Marshal(underlyingTx)
	if err != nil {
		return nil, "", fmt.Errorf("could not json marshal initial tx result %w", err)
	}

	tx := ethtypes.Transaction{}

	err = tx.UnmarshalJSON(jsonTx)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal json tx to geth based json tx: %w", err)
	}
	return &tx, genericHash, nil
}
func genericResultToError(result interface{}) (*RPCJSONError, error) {
	jsonErrorData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("unable to json marshal error result: %w", err)
	}
	fullError := new(RPCJSONError)
	err = json.Unmarshal(jsonErrorData, fullError)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json error: %w", err)
	}
	return fullError, nil

}

// ArgsLatestBlockHash is meant to generate an argument with the
// latest block hash for testing
func ArgsLatestBlockHash(ctx context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		blockData, err := getBlock(ctx, rpcClient, "latest")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve latest block hash")
			return []interface{}{"latest"}
		}
		rawHash := blockData["hash"]
		strHash, ok := rawHash.(string)
		if !ok {
			log.Error().Interface("rawHash", rawHash).Msg("The type of raw hash was expected to be string")
			return []interface{}{"latest"}
		}
		log.Trace().Str("blockHash", strHash).Msg("Got latest blockhash")

		args := []interface{}{strHash}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsLatestBlockNumber will inject arguments that correspond to the
// most recent block's number
func ArgsLatestBlockNumber(ctx context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		blockData, err := getBlock(ctx, rpcClient, "latest")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve latest block hash")
			return []interface{}{"latest"}
		}
		rawNumber := blockData["number"]
		hexNumber, ok := rawNumber.(string)
		if !ok {
			log.Error().Interface("rawNumber", rawNumber).Msg("The type of raw number was expected to be string")
			return []interface{}{"latest"}
		}

		log.Trace().Str("blockNumber", hexNumber).Msg("Got latest blockNumber")

		args := []interface{}{hexNumber}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsRawBlock will inject raw block RLP data into the arguments
func ArgsRawBlock(ctx context.Context, rpcClient *rpc.Client, blockNumOrHash string, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		blockData, err := getRawBlock(ctx, rpcClient, blockNumOrHash)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve latest raw block hash")
			return []interface{}{"latest"}
		}
		args := []interface{}{blockData}
		args = append(args, extraArgs...)
		return args
	}
}

func getBlock(ctx context.Context, rpcClient *rpc.Client, blockNumOrHash string) (map[string]interface{}, error) {
	blockData := make(map[string]interface{})
	err := rpcClient.CallContext(ctx, &blockData, "eth_getBlockByNumber", blockNumOrHash, false)
	return blockData, err
}
func getRawBlock(ctx context.Context, rpcClient *rpc.Client, blockNumOrHash string) (string, error) {
	var blockData string
	err := rpcClient.CallContext(ctx, &blockData, "debug_getRawBlock", blockNumOrHash)
	return blockData, err
}

// ArgsCoinbase would return arguments where the first argument is now
// the coinbase
func ArgsCoinbase(ctx context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		var coinbase string
		err := rpcClient.CallContext(ctx, &coinbase, "eth_coinbase")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve coinbase")
			return []interface{}{""}
		}
		log.Trace().Str("coinbase", coinbase).Msg("Got coinbase")

		args := []interface{}{coinbase}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsBlockFilterID will inject an argument that's a filter id
// corresponding to a block filte
func ArgsBlockFilterID(ctx context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		var filterId string
		err := rpcClient.CallContext(ctx, &filterId, "eth_newBlockFilter")
		if err != nil {
			log.Error().Err(err).Msg("Unable to create new block filter")
			return []interface{}{"0x0"}
		}
		log.Trace().Str("filterid", filterId).Msg("Created filter")

		args := []interface{}{filterId}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsFilterID will inject an argument that's a filter id
// corresponding to the provide filter args
func ArgsFilterID(ctx context.Context, rpcClient *rpc.Client, filterArgs RPCTestFilterArgs, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		var filterId string
		err := rpcClient.CallContext(ctx, &filterId, "eth_newFilter", filterArgs)
		if err != nil {
			log.Error().Err(err).Msg("Unable to create new block filter")
			return []interface{}{"0x0"}
		}
		log.Trace().Str("filterid", filterId).Msg("Created filter")

		args := []interface{}{filterId}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsCoinbaseTransaction will return arguments where the from is replace with the current coinbase
func ArgsCoinbaseTransaction(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		var coinbase string
		err := rpcClient.CallContext(ctx, &coinbase, "eth_coinbase")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve coinbase")
			return []interface{}{""}
		}
		tx.From = coinbase
		log.Trace().Str("coinbase", coinbase).Msg("Got coinbase")
		args := []interface{}{tx}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsSignTransaction will take the junk transaction type that we've
// created, convert it to a geth style dynamic fee transaction and
// sign it with the user provide key.
func ArgsSignTransaction(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		testAccountNonceMutex.Lock()
		defer testAccountNonceMutex.Unlock()
		curNonce := testAccountNonce

		stringTx, err := getSignedRawTx(tx, curNonce)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to sign tx")
		}

		testAccountNonce += 1

		args := []interface{}{hexutil.Encode(stringTx)}
		args = append(args, extraArgs...)
		return args
	}
}

// ArgsSignTransactionWithNonce can be used to manipulate the nonce
// directly in order to create some error cases
func ArgsSignTransactionWithNonce(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs, nonce uint64, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		stringTx, err := getSignedRawTx(tx, nonce)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to sign tx")
		}

		args := []interface{}{hexutil.Encode(stringTx)}
		args = append(args, extraArgs...)
		return args
	}
}

func getSignedRawTx(tx *RPCTestTransactionArgs, curNonce uint64) ([]byte, error) {
	chainId := currentChainID

	dft := GenericTransactionToDynamicFeeTx(tx)
	dft.ChainID = chainId
	dft.Nonce = curNonce

	londonSigner := ethtypes.NewLondonSigner(chainId)
	signedTx, err := ethtypes.SignNewTx(testPrivateKey, londonSigner, &dft)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue signing the transaction")
		return nil, err
	}
	stringTx, err := signedTx.MarshalBinary()
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshal binary for transaction")
		return nil, err
	}
	return stringTx, nil
}

// ArgsTransactionHash will execute the provided transaction and return
// the transaction hash as an argument to be used in other tests.
func ArgsTransactionHash(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) func() []interface{} {
	return func() []interface{} {
		resultHash, _, err := prepareAndSendTransaction(ctx, rpcClient, tx)
		if err != nil {
			log.Error().Err(err).Msg("Unable to execute transaction")
		}
		log.Info().Str("resultHash", resultHash).Msg("Successfully executed transaction")

		return []interface{}{resultHash}
	}
}

// ArgsTransactionBlockHashAndIndex will execute the provided transaction and return
// the block hash and index of the given transaction
func ArgsTransactionBlockHashAndIndex(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) func() []interface{} {
	return func() []interface{} {
		resultHash, receipt, err := prepareAndSendTransaction(ctx, rpcClient, tx)
		if err != nil {
			log.Error().Err(err).Msg("Unable to execute transaction")
		}
		log.Info().Str("resultHash", resultHash).Msg("Successfully executed transaction")

		return []interface{}{receipt["blockHash"], receipt["transactionIndex"]}
	}
}

// ArgsTransactionBlockNumberAndIndex will execute the provided transaction and return
// the block number and index of the given transaction
func ArgsTransactionBlockNumberAndIndex(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) func() []interface{} {
	return func() []interface{} {
		resultHash, receipt, err := prepareAndSendTransaction(ctx, rpcClient, tx)
		if err != nil {
			log.Error().Err(err).Msg("Unable to execute transaction")
		}
		log.Info().Str("resultHash", resultHash).Msg("Successfully executed transaction")

		return []interface{}{receipt["blockNumber"], receipt["transactionIndex"]}
	}
}

func prepareAndSendTransaction(ctx context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) (string, map[string]interface{}, error) {
	testAccountNonceMutex.Lock()
	defer testAccountNonceMutex.Unlock()
	curNonce := testAccountNonce

	stringTx, err := getSignedRawTx(tx, curNonce)
	if err != nil {
		return "", nil, err
	}

	resultHash, receipt, err := executeRawTxAndWait(ctx, rpcClient, stringTx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to execute transaction")
		return "", nil, err
	}

	testAccountNonce += 1

	return resultHash, receipt, nil
}

func executeRawTxAndWait(ctx context.Context, rpcClient *rpc.Client, rawTx []byte) (string, map[string]interface{}, error) {
	rawHash, err := executeRawTx(ctx, rpcClient, rawTx)
	if err != nil {
		return "", nil, err
	}

	receipt, err := waitForReceipt(ctx, rpcClient, rawHash)
	if err != nil {
		return "", nil, err
	}
	return rawHash, receipt, nil
}
func executeRawTx(ctx context.Context, rpcClient *rpc.Client, rawTx []byte) (string, error) {
	var result interface{}
	err := rpcClient.CallContext(ctx, &result, "eth_sendRawTransaction", hexutil.Encode(rawTx))
	if err != nil {
		log.Error().Err(err).Msg("Unable to send raw transaction")
		return "", err
	}
	rawHash, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid result type. Expected string but got %T", result)
	}
	log.Info().Str("txHash", rawHash).Msg("Successfully sent transaction")
	return rawHash, nil
}

func waitForReceipt(ctx context.Context, rpcClient *rpc.Client, txHash string) (map[string]interface{}, error) {
	var err error
	var result interface{}
	for i := 0; i < 30; i += 1 {
		err = rpcClient.CallContext(ctx, &result, "eth_getTransactionReceipt", txHash)
		txReceipt, ok := result.(map[string]interface{})
		if err != nil || !ok {
			time.Sleep(2 * time.Second)
			continue
		}
		log.Info().Interface("txReceipt", txReceipt).Msg("Successfully got receipt")
		return txReceipt, nil
	}
	return nil, err
}

// GenericTransactionToDynamicFeeTx convert the simple tx
// representation that we have into a standard eth type
func GenericTransactionToDynamicFeeTx(tx *RPCTestTransactionArgs) ethtypes.DynamicFeeTx {
	dft := ethtypes.DynamicFeeTx{}
	dft.GasTipCap = hexutil.MustDecodeBig(tx.MaxPriorityFeePerGas)
	dft.GasFeeCap = hexutil.MustDecodeBig(tx.MaxFeePerGas)
	dft.Gas = hexutil.MustDecodeUint64(tx.Gas)
	toAddr := ethcommon.HexToAddress(tx.To)
	dft.To = &toAddr
	dft.Value = hexutil.MustDecodeBig(tx.Value)
	dft.Data = hexutil.MustDecode(tx.Data)
	return dft
}

// GetTestAccountNonce will attempt to get the current nonce for the
// current test account
func GetTestAccountNonce(ctx context.Context, rpcClient *rpc.Client) (uint64, error) {
	ec := ethclient.NewClient(rpcClient)
	curNonce, err := ec.NonceAt(ctx, testEthAddress, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve nonce")
		curNonce = 0
	}
	log.Trace().Uint64("curNonce", curNonce).Msg("Current nonce value")
	return curNonce, err
}

// GetCurrentChainID will attempt to determine the chain for the current network
func GetCurrentChainID(ctx context.Context, rpcClient *rpc.Client) (*big.Int, error) {
	ec := ethclient.NewClient(rpcClient)
	chainId, err := ec.ChainID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get chain id")
		chainId = big.NewInt(1)

	}
	log.Trace().Uint64("chainId", chainId.Uint64()).Msg("Fetch chainid")
	return chainId, err
}

func CallRPCAndValidate(ctx context.Context, rpcClient *rpc.Client, wrappedHttpClient wrappedHttpClient, currTest RPCTest) testreporter.TestResult {
	currTestResult := testreporter.New(currTest.GetName(), currTest.GetMethod(), 1)
	args := currTest.GetArgs()

	var result interface{}
	var err error
	switch currTest.(type) {
	case *RPCTestRawHTTP:
		// Marshal the HTTP request payload.
		var payload []byte
		payload, err = json.Marshal(args)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to marshal HTTP request payload")
		}

		// Create the request.
		var request *http.Request
		request, err = http.NewRequest(currTest.GetMethod(), wrappedHttpClient.url, bytes.NewBuffer(payload)) // TODO: fix
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to create HTTP request")

		}
		request.Header.Set("Content-Type", "application/json")

		// Send the request.
		var response *http.Response
		response, err = wrappedHttpClient.client.Do(request)
		if err != nil {
			log.Error().Err(err).Msg("Unable to send HTTP request")
			break
		}
		defer response.Body.Close()

		// Read the response body.
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			log.Error().Err(err).Msg("Unable to read HTTP body")
			break
		}

		// Marshal the response and extract the error if there is any.
		var rpcResponse RPCJSONResponse
		err = json.Unmarshal(body, &rpcResponse)
		if err != nil {
			log.Error().Err(err).Msg("Unable to unmarshal HTTP body")
			break
		}
		if rpcResponse.Error != nil {
			result = &rpcResponse.Error
		}
	default:
		err = rpcClient.CallContext(ctx, &result, currTest.GetMethod(), args...)

		if err != nil && !currTest.ExpectError() {
			currTestResult.Fail(args, result, errors.New("method test failed: "+err.Error()))
			return currTestResult
		}
		if err == nil && currTest.ExpectError() {
			currTestResult.Fail(args, result, errors.New("expected an error but didn't get one"))
			return currTestResult
		}
	}

	if currTest.ExpectError() {
		err = currTest.Validate(err)
	} else {
		err = currTest.Validate(result)
	}

	if err != nil {
		currTestResult.Fail(args, result, errors.New("Failed to validate: "+err.Error()))
		return currTestResult
	}

	currTestResult.Pass(args, result, err)

	return currTestResult
}

func CallRPCWithFuzzAndValidate(ctx context.Context, rpcClient *rpc.Client, currTest RPCTest) testreporter.TestResult {
	currTestResult := testreporter.New(currTest.GetName()+"-FUZZED", currTest.GetMethod(), *testFuzzNum)

	originalArgs := currTest.GetArgs()
	for i := 0; i < *testFuzzNum; i++ {
		args := originalArgs
		fuzzer.Fuzz(&args)

		var result interface{}
		err := rpcClient.CallContext(ctx, &result, currTest.GetMethod(), args...)

		if err != nil {
			currTestResult.Fail(args, result, err)
		} else {
			currTestResult.Pass(args, result, err)
		}
	}

	return currTestResult
}

func (r *RPCTestGeneric) GetMethod() string {
	return r.Method
}
func (r *RPCTestGeneric) GetName() string {
	return r.Name
}
func (r *RPCTestGeneric) GetArgs() []interface{} {
	return r.Args
}
func (r *RPCTestGeneric) Validate(result interface{}) error {
	return r.Validator(result)
}
func (r *RPCTestGeneric) ExpectError() bool {
	return r.Flags&FlagErrorValidation != 0
}

func (r *RPCTestDynamicArgs) GetMethod() string {
	return r.Method
}
func (r *RPCTestDynamicArgs) GetName() string {
	return r.Name
}
func (r *RPCTestDynamicArgs) GetArgs() []interface{} {
	return r.Args()
}
func (r *RPCTestDynamicArgs) Validate(result interface{}) error {
	return r.Validator(result)
}
func (r *RPCTestDynamicArgs) ExpectError() bool {
	return r.Flags&FlagErrorValidation != 0
}

func (r *RPCTestRawHTTP) GetMethod() string {
	return r.HTTPMethod
}
func (r *RPCTestRawHTTP) GetName() string {
	return r.Name
}
func (r *RPCTestRawHTTP) GetArgs() []interface{} {
	return r.Args
}
func (r *RPCTestRawHTTP) Validate(result interface{}) error {
	return r.Validator(result)
}
func (r *RPCTestRawHTTP) ExpectError() bool {
	return r.Flags&FlagErrorValidation != 0
}

func (r *RPCJSONError) Error() string {
	return r.Message
}

type wrappedHttpClient struct {
	client *http.Client
	url    string
}

func shouldRunTest(t RPCTest) bool {
	var testNamespace string
	switch t.(type) {
	case *RPCTestRawHTTP:
		testNamespace = fmt.Sprintf("%s_", rpcTestRawHTTPNamespace)
	default:
		testNamespace = t.GetMethod()
	}

	for _, ns := range enabledNamespaces {
		if strings.HasPrefix(testNamespace, ns) {
			return true
		}
	}
	return false
}
