package rpcfuzz

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
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

	// cast rpc --rpc-url localhost:8545 eth_syncing
	RPCTestEthSyncing = RPCTestGeneric{
		Method: "eth_syncing",
		Args:   []interface{}{},
		Validator: ChainValidator(
			ValidateExact(false),
			ValidateJSONSchema(rpctypes.RPCSchemaEthSyncing),
		),
	}

	// I probably need to put these giant strings somewhere else
	// cast block --rpc-url localhost:8545 0
	RPCTestEthBlockByNumber = RPCTestGeneric{
		Method:    "eth_getBlockByNumber",
		Args:      []interface{}{"0x0", true},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
	}

	allTests = []RPCTest{
		&RPCTestNetVersion,
		&RPCTestWeb3ClientVersion,
		&RPCTestWeb3SHA3,
		&RPCTestWeb3SHA3Error,
		&RPCTestNetListening,
		&RPCTestNetPeerCount,
		&RPCTestEthProtocolVersion,
		&RPCTestEthSyncing,
		&RPCTestEthBlockByNumber,
	}
)

// ChainValidator would take a list of validation functions to be
// applied in order. The idea is that if first validator is true, then
// the rest won't be applied.
func ChainValidator(validators ...func(interface{}) error) func(result interface{}) error {
	return func(result interface{}) error {
		for _, v := range validators {
			err := v(result)
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("All Validation failed")
	}

}
func ValidateJSONSchema(schema string) func(result interface{}) error {
	return func(result interface{}) error {
		validatorLoader := gojsonschema.NewStringLoader(schema)

		// This is weird, but the current setup doesn't allow
		// for easy access to the initial response string...
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("Unable to marshal result back to json for validation: %w", err)
		}
		responseLoader := gojsonschema.NewStringLoader(string(jsonBytes))

		validatorResult, err := gojsonschema.Validate(validatorLoader, responseLoader)
		if err != nil {
			return fmt.Errorf("Unable to run json validation: %w", err)
		}
		// fmt.Println(string(jsonBytes))
		if !validatorResult.Valid() {
			errStr := ""
			for _, desc := range validatorResult.Errors() {
				errStr += desc.String() + "\n"
			}
			return fmt.Errorf("The json document is not valid: %s", errStr)
		}
		return nil

	}
}
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

- https://ethereum.github.io/execution-apis/api-documentation/
- https://ethereum.org/en/developers/docs/apis/json-rpc/
- https://json-schema.org/
- https://www.liquid-technologies.com/online-json-to-schema-converter

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
