package rpcfuzz

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	fuzz "github.com/google/gofuzz"
	"github.com/maticnetwork/polygon-cli/cmd/rpcfuzz/argfuzz"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	// flags
	rpcUrl               *string
	testPrivateHexKey    *string
	testContractAddress  *string
	testNamespaces       *string
	testFuzz             *bool
	testFuzzNum          *int
	seed                 *int64
	testOutputExportPath *string
	testExportJson       *bool
	testExportCSV        *bool
	testExportMarkdown   *bool
	testExportHTML       *bool
)

var RPCFuzzCmd = &cobra.Command{
	Use:   "rpcfuzz",
	Short: "Continually run a variety of RPC calls and fuzzers.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRpcFuzz(cmd.Context())
	},
}

func init() {
	flagSet := RPCFuzzCmd.PersistentFlags()

	rpcUrl = flagSet.StringP("rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	testPrivateHexKey = flagSet.String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to sending transactions")
	testContractAddress = flagSet.String("contract-address", "0x6fda56c57b0acadb96ed5624ac500c0429d59429", "The address of a contract that can be used for testing")
	testNamespaces = flagSet.String("namespaces", fmt.Sprintf("eth,web3,net,debug,%s", rpcTestRawHTTPNamespace), "Comma separated list of rpc namespaces to test")
	testFuzz = flagSet.Bool("fuzz", false, "Flag to indicate whether to fuzz input or not.")
	testFuzzNum = flagSet.Int("fuzzn", 100, "Number of times to run the fuzzer per test.")
	seed = flagSet.Int64("seed", 123456, "A seed for generating random values within the fuzzer")
	testOutputExportPath = flagSet.String("export-path", "", "The directory export path of the output of the tests. Must pair this with either --json, --csv, --md, or --html")
	testExportJson = flagSet.Bool("json", false, "Flag to indicate that output will be exported as a JSON.")
	testExportCSV = flagSet.Bool("csv", false, "Flag to indicate that output will be exported as a CSV.")
	testExportMarkdown = flagSet.Bool("md", false, "Flag to indicate that output will be exported as a Markdown.")
	testExportHTML = flagSet.Bool("html", false, "Flag to indicate that output will be exported as a HTML.")

	argfuzz.SetSeed(seed)

	fuzzer = fuzz.New()
	fuzzer.Funcs(argfuzz.FuzzRPCArgs)
}

func checkFlags() (err error) {
	// Check rpc-url flag.
	if rpcUrl == nil {
		panic("RPC URL is empty")
	}
	if err = util.ValidateUrl(*rpcUrl); err != nil {
		return
	}

	// Check private key flag.
	privateKey, err := crypto.HexToECDSA(*testPrivateHexKey)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't process the hex private key")
		return err
	}
	ethAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	log.Info().Str("ethAddress", ethAddress.String()).Msg("Loaded private key")

	// Check namespace flag.
	nsValidator := regexp.MustCompile("^[a-z0-9]*$")
	rawNameSpaces := strings.Split(*testNamespaces, ",")
	enabledNamespaces = make([]string, 0)
	for _, ns := range rawNameSpaces {
		if !nsValidator.MatchString(ns) {
			return fmt.Errorf("the namespace %s is not valid", ns)
		}
		enabledNamespaces = append(enabledNamespaces, ns+"_")
	}
	log.Info().Strs("namespaces", enabledNamespaces).Msg("Enabling namespaces")

	testPrivateKey = privateKey
	testEthAddress = ethAddress

	return nil
}

func runRpcFuzz(ctx context.Context) error {
	if *testOutputExportPath != "" && !*testExportJson && !*testExportCSV && !*testExportMarkdown && !*testExportHTML {
		log.Warn().Msg("Setting --export-path must pair with a export type: --json, --csv, --md, or --html")
	}

	rpcClient, err := rpc.DialContext(ctx, *rpcUrl)
	if err != nil {
		return err
	}
	nonce, err := GetTestAccountNonce(ctx, rpcClient)
	if err != nil {
		return err
	}
	chainId, err := GetCurrentChainID(ctx, rpcClient)
	if err != nil {
		return err
	}
	testAccountNonce = nonce
	currentChainID = chainId

	log.Trace().Uint64("nonce", nonce).Uint64("chainid", chainId.Uint64()).Msg("Doing test setup")
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
