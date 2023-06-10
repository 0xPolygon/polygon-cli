package rpcfuzz

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

type (
	RPCTest interface {
		GetMethod() string
		GetArgs() []interface{}
		Validate(result interface{}) error
	}

	RPCTestGeneric struct {
		Method    string
		Args      []interface{}
		Validator func(result interface{}) error
	}
)

var (
	// cast rpc --rpc-url localhost:8545 net_version
	RPCTestNetVersion = RPCTestGeneric{
		Method:    "net_version",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^\d*$`),
	}

	// cast rpc --rpc-url localhost:8545 web3_clientVersion
	RPCTestWeb3ClientVersion = RPCTestGeneric{
		Method:    "web3_clientVersion",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^[[:print:]]*$`),
	}

	allTests = []RPCTest{
		&RPCTestNetVersion,
		&RPCTestWeb3ClientVersion,
	}
)

func ValidateRegexString(regEx string) func(result interface{}) error {
	r := regexp.MustCompile(regEx)
	return func(result interface{}) error {
		resultStr, isValid := result.(string)
		if !isValid {
			return fmt.Errorf("Invalid result type. Expected string but got %T", result)
		}
		if !r.MatchString(resultStr) {
			return fmt.Errorf("The regex %s failed to match result %s", regEx, resultStr)
		}
		return nil
	}
}

func (r *RPCTestGeneric) GetMethod() string {
	return r.Method
}
func (r *RPCTestGeneric) GetArgs() []interface{} {
	return r.Args
}
func (r *RPCTestGeneric) Validate(result interface{}) error {
	return r.Validator(result)
}

var RPCFuzzCmd = &cobra.Command{
	Use:   "rpcfuzz http://localhost:8545",
	Short: "Continually run a variety of RPC calls and fuzzers",
	Long: `
beep

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cxt := cmd.Context()
		rpcClient, err := rpc.DialContext(cxt, args[0])
		if err != nil {
			return err
		}
		for _, t := range allTests {
			log.Trace().Str("method", t.GetMethod()).Msg("Running Test")
			var result interface{}
			err = rpcClient.CallContext(cxt, &result, t.GetMethod(), t.GetArgs()...)
			if err != nil {
				log.Error().Err(err).Str("method", t.GetMethod()).Msg("Method test failed")
				continue
			}
			err = t.Validate(result)
			if err != nil {
				log.Error().Err(err).Str("method", t.GetMethod()).Msg("Failed to validate")
				continue
			}
			log.Info().Str("method", t.GetMethod()).Msg("Successfully validated")
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Expected 1 argument, but got %d", len(args))
		}
		return nil
	},
}

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flagSet := RPCFuzzCmd.PersistentFlags()
	_ = flagSet
}
