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
		ExpectError() bool
	}

	RPCTestGeneric struct {
		Method    string
		Args      []interface{}
		Validator func(result interface{}) error
		IsError   bool
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

	// cast rpc --rpc-url localhost:8545 web3_sha3 0x68656c6c6f20776f726c64
	RPCTestWeb3SHA3 = RPCTestGeneric{
		Method:    "web3_sha3",
		Args:      []interface{}{"0x68656c6c6f20776f726c64"},
		Validator: ValidateRegexString(`0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad`),
	}

	RPCTestWeb3SHA3Error = RPCTestGeneric{
		IsError:   true,
		Method:    "web3_sha3",
		Args:      []interface{}{"68656c6c6f20776f726c64"},
		Validator: ValidatorError(`cannot unmarshal hex string without 0x prefix`),
	}

	// cast rpc --rpc-url localhost:8545 net_listening
	RPCTestNetListening = RPCTestGeneric{
		Method:    "net_listening",
		Args:      []interface{}{},
		Validator: ValidateExact(true),
	}

	// cast rpc --rpc-url localhost:8545 net_peerCount
	RPCTestNetPeerCount = RPCTestGeneric{
		Method:    "net_peerCount",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]*$`),
	}

	// cast rpc --rpc-url localhost:8545 eth_protocolVersion
	RPCTestEthProtocolVersion = RPCTestGeneric{
		IsError:   true,
		Method:    "eth_protocolVersion",
		Args:      []interface{}{},
		Validator: ValidatorError(`method eth_protocolVersion does not exist`),
	}

	allTests = []RPCTest{
		&RPCTestNetVersion,
		&RPCTestWeb3ClientVersion,
		&RPCTestWeb3SHA3,
		&RPCTestWeb3SHA3Error,
		&RPCTestNetListening,
		&RPCTestNetPeerCount,
		&RPCTestEthProtocolVersion,
	}
)

func ValidateExact(expected interface{}) func(result interface{}) error {
	return func(result interface{}) error {
		if expected != result {
			return fmt.Errorf("Expected %v and got %v", expected, result)
		}
		return nil
	}
}

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
func ValidatorError(errorMessageRegex string) func(result interface{}) error {
	r := regexp.MustCompile(errorMessageRegex)
	return func(result interface{}) error {
		resultError, isValid := result.(error)
		if !isValid {
			return fmt.Errorf("Invalid result type. Expected error but got %T", result)
		}
		if !r.MatchString(resultError.Error()) {
			return fmt.Errorf("The regex %s failed to match result %s", errorMessageRegex, resultError.Error())
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
func (r *RPCTestGeneric) ExpectError() bool {
	return r.IsError
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
			if err != nil && !t.ExpectError() {
				log.Error().Err(err).Str("method", t.GetMethod()).Msg("Method test failed")
				continue
			}

			if t.ExpectError() {
				err = t.Validate(err)
			} else {
				err = t.Validate(result)
			}

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
