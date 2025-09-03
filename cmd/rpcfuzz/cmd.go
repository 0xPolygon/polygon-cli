package rpcfuzz

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/0xPolygon/polygon-cli/cmd/rpcfuzz/argfuzz"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/crypto"
	fuzz "github.com/google/gofuzz"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	// flags
	rpcUrl              *string
	testPrivateHexKey   *string
	testContractAddress *string
	testNamespaces      *string
	testFuzz            *bool
	testFuzzNum         *int
	seed                *int64
	streamJSON          *bool
	streamCSV           *bool
	streamCompact       *bool
	streamHTML          *bool
	streamMarkdown      *bool
	outputFilter        *string
	summaryInterval     *int
	quietMode           *bool
)

var RPCFuzzCmd = &cobra.Command{
	Use:   "rpcfuzz",
	Short: "Continually run a variety of RPC calls and fuzzers.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rpcUrl = flag_loader.GetRpcUrlFlagValue(cmd)
		testPrivateHexKey = flag_loader.GetPrivateKeyFlagValue(cmd)
	},
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
	testContractAddress = flagSet.String("contract-address", "", "The address of a contract that can be used for testing. If not specified, a contract will be deployed automatically.")
	testNamespaces = flagSet.String("namespaces", fmt.Sprintf("eth,web3,net,debug,%s", rpcTestRawHTTPNamespace), "Comma separated list of rpc namespaces to test")
	testFuzz = flagSet.Bool("fuzz", false, "Flag to indicate whether to fuzz input or not.")
	testFuzzNum = flagSet.Int("fuzzn", 100, "Number of times to run the fuzzer per test.")
	seed = flagSet.Int64("seed", 123456, "A seed for generating random values within the fuzzer")

	// Streamer type flags (mutually exclusive)
	streamJSON = flagSet.Bool("json", false, "Stream output in JSON format")
	streamCSV = flagSet.Bool("csv", false, "Stream output in CSV format")
	streamCompact = flagSet.Bool("compact", false, "Stream output in compact format (default)")
	streamHTML = flagSet.Bool("html", false, "Stream output in HTML format")
	streamMarkdown = flagSet.Bool("md", false, "Stream output in Markdown format")

	// Output control flags
	outputFilter = flagSet.String("output", "all", "What to output: all, failures, summary")
	summaryInterval = flagSet.Int("summary-interval", 0, "Print summary every N tests (0=disabled)")
	quietMode = flagSet.Bool("quiet", false, "Only show final summary")

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

	// Ensure only one streamer type is selected
	streamerCount := 0
	if *streamJSON {
		streamerCount++
	}
	if *streamCSV {
		streamerCount++
	}
	if *streamCompact {
		streamerCount++
	}
	if *streamHTML {
		streamerCount++
	}
	if *streamMarkdown {
		streamerCount++
	}

	if streamerCount > 1 {
		return fmt.Errorf("only one output format can be specified: --json, --csv, --compact, --html, or --md")
	}

	// Check private key flag.
	trimmedHexPrivateKey := strings.TrimPrefix(*testPrivateHexKey, "0x")
	privateKey, err := crypto.HexToECDSA(trimmedHexPrivateKey)
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
