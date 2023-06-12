// Package rpcfuzz is meant to have some basic RPC fuzzing and
// conformance tests
package rpcfuzz

// TODO add configuration for name space
// TODO add the open rpc schemas
// TODO refactor to remove names

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"math/big"
	"os"
	"regexp"
)

type (
	// RPCTest is the common interface for a test.  In the future
	// we'll need some addition methods in particular if don't
	// want to run tests that require unlocked accounts or if we
	// want to skip certain namepaces
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

	// RPCTestGeneric is the simplist implementation of the
	// RPCTest. Basically the implementation of the interface is
	// managed by just returning hard coded values for method,
	// args, validator, and error
	RPCTestGeneric struct {
		Name           string
		Method         string
		Args           []interface{}
		Validator      func(result interface{}) error
		IsError        bool
		RequiresUnlock bool
	}

	// RPCTestDynamicArgs is a simple implementation of the
	// RPCTest that requires a function for Args which will be
	// used to generate the args for testing.
	RPCTestDynamicArgs struct {
		Name           string
		Method         string
		Args           func() []interface{}
		Validator      func(result interface{}) error
		IsError        bool
		RequiresUnlock bool
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
)

const (
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

var (
	testPrivateHexKey   *string
	testContractAddress *string
	testPrivateKey      *ecdsa.PrivateKey
	testEthAddress      ethcommon.Address
)

var (
	RPCTestNetVersion                                  RPCTestGeneric
	RPCTestWeb3ClientVersion                           RPCTestGeneric
	RPCTestWeb3SHA3                                    RPCTestGeneric
	RPCTestWeb3SHA3Error                               RPCTestGeneric
	RPCTestNetListening                                RPCTestGeneric
	RPCTestNetPeerCount                                RPCTestGeneric
	RPCTestEthProtocolVersion                          RPCTestGeneric
	RPCTestEthSyncing                                  RPCTestGeneric
	RPCTestEthCoinbase                                 RPCTestGeneric
	RPCTestEthChainID                                  RPCTestGeneric
	RPCTestEthMining                                   RPCTestGeneric
	RPCTestEthHashrate                                 RPCTestGeneric
	RPCTestEthGasPrice                                 RPCTestGeneric
	RPCTestEthAccounts                                 RPCTestGeneric
	RPCTestEthBlockNumber                              RPCTestGeneric
	RPCTestEthGetBalanceLatest                         RPCTestGeneric
	RPCTestEthGetBalanceEarliest                       RPCTestGeneric
	RPCTestEthGetBalancePending                        RPCTestGeneric
	RPCTestEthGetBalanceZero                           RPCTestGeneric
	RPCTestEthGetStorageAtLatest                       RPCTestGeneric
	RPCTestEthGetStorageAtEarliest                     RPCTestGeneric
	RPCTestEthGetStorageAtPending                      RPCTestGeneric
	RPCTestEthGetStorageAtZero                         RPCTestGeneric
	RPCTestEthGetTransactionCountAtLatest              RPCTestGeneric
	RPCTestEthGetTransactionCountAtEarliest            RPCTestGeneric
	RPCTestEthGetTransactionCountAtPending             RPCTestGeneric
	RPCTestEthGetTransactionCountAtZero                RPCTestGeneric
	RPCTestEthGetBlockTransactionCountByHash           RPCTestDynamicArgs
	RPCTestEthGetBlockTransactionCountByHashMissing    RPCTestGeneric
	RPCTestEthGetBlockTransactionCountByNumberLatest   RPCTestGeneric
	RPCTestEthGetBlockTransactionCountByNumberEarliest RPCTestGeneric
	RPCTestEthGetBlockTransactionCountByNumberPending  RPCTestGeneric
	RPCTestEthGetBlockTransactionCountByNumberZero     RPCTestGeneric
	RPCTestEthGetUncleCountByBlockHash                 RPCTestDynamicArgs
	RPCTestEthGetUncleCountByBlockHashMissing          RPCTestGeneric
	RPCTestEthGetUncleCountByBlockNumberLatest         RPCTestGeneric
	RPCTestEthGetUncleCountByBlockNumberEarliest       RPCTestGeneric
	RPCTestEthGetUncleCountByBlockNumberPending        RPCTestGeneric
	RPCTestEthGetUncleCountByBlockNumberZero           RPCTestGeneric
	RPCTestEthGetCodeLatest                            RPCTestGeneric
	RPCTestEthGetCodePending                           RPCTestGeneric
	RPCTestEthGetCodeEarliest                          RPCTestGeneric
	RPCTestEthSign                                     RPCTestDynamicArgs
	RPCTestEthSignFail                                 RPCTestGeneric
	RPCTestEthSignTransaction                          RPCTestDynamicArgs
	RPCTestEthSendTransaction                          RPCTestDynamicArgs
	RPCTestEthSendRawTransaction                       RPCTestDynamicArgs
	RPCTestEthCallLatest                               RPCTestGeneric
	RPCTestEthCallEarliest                             RPCTestGeneric
	RPCTestEthCallPending                              RPCTestGeneric
	RPCTestEthCallZero                                 RPCTestGeneric
	RPCTestEthEstimateGas                              RPCTestGeneric
	RPCTestEthGetBlockByHash                           RPCTestDynamicArgs
	RPCTestEthGetBlockByHashNoTx                       RPCTestDynamicArgs
	RPCTestEthGetBlockByHashZero                       RPCTestGeneric

	allTests                = make([]RPCTest, 0)
	RPCTestEthBlockByNumber RPCTestGeneric
)

func setupTests(cxt context.Context, rpcClient *rpc.Client) {
	// cast rpc --rpc-url localhost:8545 net_version
	RPCTestNetVersion = RPCTestGeneric{
		Name:      "RPCTestNetVersion",
		Method:    "net_version",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^\d*$`),
	}
	allTests = append(allTests, &RPCTestNetVersion)

	// cast rpc --rpc-url localhost:8545 web3_clientVersion
	RPCTestWeb3ClientVersion = RPCTestGeneric{
		Name:      "RPCTestWeb3ClientVersion",
		Method:    "web3_clientVersion",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^[[:print:]]*$`),
	}
	allTests = append(allTests, &RPCTestWeb3ClientVersion)

	// cast rpc --rpc-url localhost:8545 web3_sha3 0x68656c6c6f20776f726c64
	RPCTestWeb3SHA3 = RPCTestGeneric{
		Name:      "RPCTestWeb3SHA3",
		Method:    "web3_sha3",
		Args:      []interface{}{"0x68656c6c6f20776f726c64"},
		Validator: ValidateRegexString(`0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad`),
	}
	allTests = append(allTests, &RPCTestWeb3SHA3)

	RPCTestWeb3SHA3Error = RPCTestGeneric{
		Name:      "RPCTestWeb3SHA3Error",
		IsError:   true,
		Method:    "web3_sha3",
		Args:      []interface{}{"68656c6c6f20776f726c64"},
		Validator: ValidateError(`cannot unmarshal hex string without 0x prefix`),
	}
	allTests = append(allTests, &RPCTestWeb3SHA3Error)

	// cast rpc --rpc-url localhost:8545 net_listening
	RPCTestNetListening = RPCTestGeneric{
		Name:      "RPCTestNetListening",
		Method:    "net_listening",
		Args:      []interface{}{},
		Validator: ValidateExact(true),
	}
	allTests = append(allTests, &RPCTestNetListening)

	// cast rpc --rpc-url localhost:8545 net_peerCount
	RPCTestNetPeerCount = RPCTestGeneric{
		Name:      "RPCTestNetPeerCount",
		Method:    "net_peerCount",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]*$`),
	}
	allTests = append(allTests, &RPCTestNetPeerCount)

	// cast rpc --rpc-url localhost:8545 eth_protocolVersion
	RPCTestEthProtocolVersion = RPCTestGeneric{
		Name:      "RPCTestEthProtocolVersion",
		IsError:   true,
		Method:    "eth_protocolVersion",
		Args:      []interface{}{},
		Validator: ValidateError(`method eth_protocolVersion does not exist`),
	}
	allTests = append(allTests, &RPCTestEthProtocolVersion)

	// cast rpc --rpc-url localhost:8545 eth_syncing
	RPCTestEthSyncing = RPCTestGeneric{
		Name:   "RPCTestEthSyncing",
		Method: "eth_syncing",
		Args:   []interface{}{},
		Validator: ChainValidator(
			ValidateExact(false),
			ValidateJSONSchema(rpctypes.RPCSchemaEthSyncing),
		),
	}
	allTests = append(allTests, &RPCTestEthSyncing)

	// cast rpc --rpc-url localhost:8545 eth_coinbase
	RPCTestEthCoinbase = RPCTestGeneric{
		Name:           "RPCTestEthCoinbase",
		Method:         "eth_coinbase",
		Args:           []interface{}{},
		Validator:      ValidateRegexString(`^0x[[:xdigit:]]{40}$`),
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthCoinbase)

	// cast rpc --rpc-url localhost:8545 eth_chainId
	RPCTestEthChainID = RPCTestGeneric{
		Name:      "RPCTestEthChainID",
		Method:    "eth_chainId",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthChainID)

	// cast rpc --rpc-url localhost:8545 eth_mining
	RPCTestEthMining = RPCTestGeneric{
		Name:   "RPCTestEthMining",
		Method: "eth_mining",
		Args:   []interface{}{},
		Validator: ChainValidator(
			ValidateExact(true),
			ValidateExact(false),
		),
	}
	allTests = append(allTests, &RPCTestEthMining)

	// cast rpc --rpc-url localhost:8545 eth_hashrate
	RPCTestEthHashrate = RPCTestGeneric{
		Name:      "RPCTestEthHashrate",
		Method:    "eth_hashrate",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthHashrate)

	// cast rpc --rpc-url localhost:8545 eth_gasPrice
	RPCTestEthGasPrice = RPCTestGeneric{
		Name:      "RPCTestEthGasPrice",
		Method:    "eth_gasPrice",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGasPrice)

	// cast rpc --rpc-url localhost:8545 eth_accounts
	RPCTestEthAccounts = RPCTestGeneric{
		Name:           "RPCTestEthAccounts",
		Method:         "eth_accounts",
		Args:           []interface{}{},
		Validator:      ValidateJSONSchema(rpctypes.RPCSchemaAccountList),
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthAccounts)

	// cast rpc --rpc-url localhost:8545 eth_blockNumber
	RPCTestEthBlockNumber = RPCTestGeneric{
		Name:      "RPCTestEthBlockNumber",
		Method:    "eth_blockNumber",
		Args:      []interface{}{},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthBlockNumber)

	// cast balance --rpc-url localhost:8545 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
	RPCTestEthGetBalanceLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceLatest",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "latest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBalanceLatest)
	RPCTestEthGetBalanceEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceEarliest",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "earliest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBalanceEarliest)
	RPCTestEthGetBalancePending = RPCTestGeneric{
		Name:      "RPCTestEthGetBalancePending",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "pending"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBalancePending)
	RPCTestEthGetBalanceZero = RPCTestGeneric{
		Name:      "RPCTestEthGetBalanceZero",
		Method:    "eth_getBalance",
		Args:      []interface{}{testEthAddress.String(), "0x0"},
		Validator: ValidateRegexString(`^0x0$`),
	}
	allTests = append(allTests, &RPCTestEthGetBalanceZero)

	// cast storage --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429 3
	RPCTestEthGetStorageAtLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtLatest",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x3", "latest"},
		Validator: ValidateRegexString(`^0x536f6c6964697479206279204578616d706c6500000000000000000000000026$`),
	}
	allTests = append(allTests, &RPCTestEthGetStorageAtLatest)
	RPCTestEthGetStorageAtEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtEarliest",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x3", "earliest"},
		Validator: ValidateRegexString(`^0x0{64}`),
	}
	allTests = append(allTests, &RPCTestEthGetStorageAtEarliest)
	RPCTestEthGetStorageAtPending = RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtPending",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x3", "pending"},
		Validator: ValidateRegexString(`^0x536f6c6964697479206279204578616d706c6500000000000000000000000026$`),
	}
	allTests = append(allTests, &RPCTestEthGetStorageAtZero)
	RPCTestEthGetStorageAtZero = RPCTestGeneric{
		Name:      "RPCTestEthGetStorageAtZero",
		Method:    "eth_getStorageAt",
		Args:      []interface{}{*testContractAddress, "0x3", "0x0"},
		Validator: ValidateRegexString(`^0x0{64}`),
	}
	allTests = append(allTests, &RPCTestEthGetStorageAtPending)

	// cast rpc --rpc-url localhost:8545 eth_getTransactionCount 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 latest
	RPCTestEthGetTransactionCountAtLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtLatest",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "latest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetTransactionCountAtLatest)
	RPCTestEthGetTransactionCountAtEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtEarliest",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "earliest"},
		Validator: ValidateRegexString(`^0x0$`),
	}
	allTests = append(allTests, &RPCTestEthGetTransactionCountAtEarliest)
	RPCTestEthGetTransactionCountAtPending = RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtPending",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "pending"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetTransactionCountAtPending)
	RPCTestEthGetTransactionCountAtZero = RPCTestGeneric{
		Name:      "RPCTestEthGetTransactionCountAtZero",
		Method:    "eth_getTransactionCount",
		Args:      []interface{}{testEthAddress.String(), "0x0"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetTransactionCountAtZero)

	// cast rpc --rpc-url localhost:8545 eth_getBlockTransactionCountByHash 0x9300b64619e167e7dbc1b41a6a6e7a8de7d6b99427dceefbd58014e328bd7f92
	RPCTestEthGetBlockTransactionCountByHash = RPCTestDynamicArgs{
		Name:      "RPCTestEthGetBlockTransactionCountByHash",
		Method:    "eth_getBlockTransactionCountByHash",
		Args:      ArgsLatestBlockHash(cxt, rpcClient),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByHash)
	RPCTestEthGetBlockTransactionCountByHashMissing = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByHashMissing",
		Method:    "eth_getBlockTransactionCountByHash",
		Args:      []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Validator: ValidateExact(nil),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByHashMissing)

	// cast rpc --rpc-url localhost:8545 eth_getBlockTransactionCountByNumber 0x1
	RPCTestEthGetBlockTransactionCountByNumberLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberLatest",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByNumberLatest)
	RPCTestEthGetBlockTransactionCountByNumberEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberEarliest",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByNumberEarliest)
	RPCTestEthGetBlockTransactionCountByNumberPending = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberPending",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"pending"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByNumberPending)
	RPCTestEthGetBlockTransactionCountByNumberZero = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockTransactionCountByNumberZero",
		Method:    "eth_getBlockTransactionCountByNumber",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetBlockTransactionCountByNumberZero)

	// cast rpc --rpc-url localhost:8545 eth_getUncleCountByBlockHash 0x9300b64619e167e7dbc1b41a6a6e7a8de7d6b99427dceefbd58014e328bd7f92
	RPCTestEthGetUncleCountByBlockHash = RPCTestDynamicArgs{
		Name:      "RPCTestEthGetUncleCountByBlockHash",
		Method:    "eth_getUncleCountByBlockHash",
		Args:      ArgsLatestBlockHash(cxt, rpcClient),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockHash)
	RPCTestEthGetUncleCountByBlockHashMissing = RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockHashMissing",
		Method:    "eth_getUncleCountByBlockHash",
		Args:      []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Validator: ValidateExact(nil),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockHashMissing)

	// cast rpc --rpc-url localhost:8545 eth_getUncleCountByBlockNumber 0x1
	RPCTestEthGetUncleCountByBlockNumberLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberLatest",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"latest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockNumberLatest)
	RPCTestEthGetUncleCountByBlockNumberEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberEarliest",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"earliest"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockNumberEarliest)
	RPCTestEthGetUncleCountByBlockNumberPending = RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberPending",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"pending"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockNumberPending)
	RPCTestEthGetUncleCountByBlockNumberZero = RPCTestGeneric{
		Name:      "RPCTestEthGetUncleCountByBlockNumberZero",
		Method:    "eth_getUncleCountByBlockNumber",
		Args:      []interface{}{"0x0"},
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{1,}$`),
	}
	allTests = append(allTests, &RPCTestEthGetUncleCountByBlockNumberZero)

	// cast code --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429
	RPCTestEthGetCodeLatest = RPCTestGeneric{
		Name:      "RPCTestEthGetCodeLatest",
		Method:    "eth_getCode",
		Args:      []interface{}{*testContractAddress, "latest"},
		Validator: ValidateHashedResponse("e39381f1654cf6a3b7eac2a789b9adf7319312cb"),
	}
	allTests = append(allTests, &RPCTestEthGetCodeLatest)
	RPCTestEthGetCodePending = RPCTestGeneric{
		Name:      "RPCTestEthGetCodePending",
		Method:    "eth_getCode",
		Args:      []interface{}{*testContractAddress, "pending"},
		Validator: ValidateHashedResponse("e39381f1654cf6a3b7eac2a789b9adf7319312cb"),
	}
	allTests = append(allTests, &RPCTestEthGetCodePending)
	RPCTestEthGetCodeEarliest = RPCTestGeneric{
		Name:      "RPCTestEthGetCodeEarliest",
		Method:    "eth_getCode",
		Args:      []interface{}{*testContractAddress, "earliest"},
		Validator: ValidateRegexString(`^0x$`),
	}
	allTests = append(allTests, &RPCTestEthGetCodeEarliest)

	// cast rpc --rpc-url localhost:8545 eth_sign "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57" "0xdeadbeaf"
	RPCTestEthSign = RPCTestDynamicArgs{
		Name:           "RPCTestEthSign",
		Method:         "eth_sign",
		Args:           ArgsCoinbase(cxt, rpcClient, "0xdeadbeef"),
		Validator:      ValidateRegexString(`^0x[[:xdigit:]]{72,}$`),
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthSign)
	RPCTestEthSignFail = RPCTestGeneric{
		Name:           "RPCTestEthSignFail",
		Method:         "eth_sign",
		Args:           []interface{}{testEthAddress.String(), "0xdeadbeef"},
		Validator:      ValidateError(`unknown account`),
		IsError:        true,
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthSignFail)

	// cast rpc --rpc-url localhost:8545 eth_signTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	RPCTestEthSignTransaction = RPCTestDynamicArgs{
		Name:           "RPCTestEthSignTransaction",
		Method:         "eth_signTransaction",
		Args:           ArgsCoinbaseTransaction(cxt, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: "0x6FC23AC00", MaxPriorityFeePerGas: "0x1", Nonce: "0x1"}),
		Validator:      ValidateJSONSchema(rpctypes.RPCSchemaSignTxResponse),
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthSignTransaction)

	// cast rpc --rpc-url localhost:8545 eth_sendTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	RPCTestEthSendTransaction = RPCTestDynamicArgs{
		Name:           "RPCTestEthSendTransaction",
		Method:         "eth_sendTransaction",
		Args:           ArgsCoinbaseTransaction(cxt, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: "0x6FC23AC00", MaxPriorityFeePerGas: "0x1"}),
		Validator:      ValidateRegexString(`^0x[[:xdigit:]]{64}$`),
		RequiresUnlock: true,
	}
	allTests = append(allTests, &RPCTestEthSendTransaction)

	// cast rpc --rpc-url localhost:8545 eth_sendRawTransaction '{"from": "0xb9b1cf51a65b50f74ed8bcb258413c02cba2ec57", "to": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6", "data": "0x", "gas": "0x5208", "gasPrice": "0x1", "nonce": "0x1"}'
	RPCTestEthSendRawTransaction = RPCTestDynamicArgs{
		Name:      "RPCTestEthSendRawTransaction",
		Method:    "eth_sendRawTransaction",
		Args:      ArgsSignTransaction(cxt, rpcClient, &RPCTestTransactionArgs{To: testEthAddress.String(), Value: "0x123", Gas: "0x5208", Data: "0x", MaxFeePerGas: "0x6FC23AC00", MaxPriorityFeePerGas: "0x1"}),
		Validator: ValidateRegexString(`^0x[[:xdigit:]]{64}$`),
	}
	allTests = append(allTests, &RPCTestEthSendRawTransaction)

	// cat contracts/ERC20.abi| go run main.go abi
	// cast call --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429  'function name() view returns(string)'
	RPCTestEthCallLatest = RPCTestGeneric{
		Name:      "RPCTestEthCallLatest",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "latest"},
		Validator: ValidateRegexString(`536f6c6964697479206279204578616d706c65`),
	}
	allTests = append(allTests, &RPCTestEthCallLatest)
	RPCTestEthCallPending = RPCTestGeneric{
		Name:      "RPCTestEthCallPending",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "pending"},
		Validator: ValidateRegexString(`536f6c6964697479206279204578616d706c65`),
	}
	allTests = append(allTests, &RPCTestEthCallPending)
	RPCTestEthCallEarliest = RPCTestGeneric{
		Name:      "RPCTestEthCallEarliest",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "earliest"},
		Validator: ValidateRegexString(`^0x$`),
	}
	allTests = append(allTests, &RPCTestEthCallEarliest)
	RPCTestEthCallZero = RPCTestGeneric{
		Name:      "RPCTestEthCallZero",
		Method:    "eth_call",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0x06fdde03"}, "0x0"},
		Validator: ValidateRegexString(`^0x$`),
	}
	allTests = append(allTests, &RPCTestEthCallZero)

	// cat contracts/ERC20.abi| go run main.go abi
	// cast estimate --rpc-url localhost:8545 0x6fda56c57b0acadb96ed5624ac500c0429d59429  'function mint(uint256 amount) returns()' 10000
	// cast abi-encode 'function mint(uint256 amount) returns()' 10000
	RPCTestEthEstimateGas = RPCTestGeneric{
		Name:      "RPCTestEthEstimateGas",
		Method:    "eth_estimateGas",
		Args:      []interface{}{&RPCTestTransactionArgs{To: *testContractAddress, Value: "0x0", Data: "0xa0712d680000000000000000000000000000000000000000000000000000000000002710"}, "latest"},
		Validator: ValidateRegexString(`0x10b0d`),
	}
	allTests = append(allTests, &RPCTestEthEstimateGas)

	// cast block --rpc-url localhost:8545 latest
	RPCTestEthGetBlockByHash = RPCTestDynamicArgs{
		Name:      "RPCTestEthGetBlockByHash",
		Method:    "eth_getBlockByHash",
		Args:      ArgsLatestBlockHash(cxt, rpcClient, true),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
	}
	allTests = append(allTests, &RPCTestEthGetBlockByHash)
	RPCTestEthGetBlockByHashNoTx = RPCTestDynamicArgs{
		Name:      "RPCTestEthGetBlockByHashNoTx",
		Method:    "eth_getBlockByHash",
		Args:      ArgsLatestBlockHash(cxt, rpcClient, false),
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
	}
	allTests = append(allTests, &RPCTestEthGetBlockByHashNoTx)
	RPCTestEthGetBlockByHashZero = RPCTestGeneric{
		Name:      "RPCTestEthGetBlockByHashZero",
		Method:    "eth_getBlockByHash",
		Args:      []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000000", true},
		Validator: ValidateExact(nil),
	}
	allTests = append(allTests, &RPCTestEthGetBlockByHashZero)

	// cast block --rpc-url localhost:8545 0
	RPCTestEthBlockByNumber = RPCTestGeneric{
		Name:      "RPCTestEthBlockByNumber",
		Method:    "eth_getBlockByNumber",
		Args:      []interface{}{"0x0", true},
		Validator: ValidateJSONSchema(rpctypes.RPCSchemaEthBlock),
	}
	allTests = append(allTests, &RPCTestEthBlockByNumber)

	uniqueTests := make(map[RPCTest]struct{})
	uniqueTestNames := make(map[string]struct{})
	for _, v := range allTests {
		_, hasKey := uniqueTests[v]
		if hasKey {
			log.Fatal().Str("name", v.GetName()).Str("method", v.GetMethod()).Msg("duplicate test detected")
		}
		uniqueTests[v] = struct{}{}
		_, hasKey = uniqueTestNames[v.GetName()]
		if hasKey {
			log.Fatal().Str("name", v.GetName()).Str("method", v.GetMethod()).Msg("duplicate test name detected")
		}
		uniqueTestNames[v.GetName()] = struct{}{}
	}

}

// ChainValidator would take a list of validation functions to be
// applied in order. The idea is that if first validator is true, then
// the rest won't be applied. This is needed to support responses that
// might be different types.
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

// ValidateHashedResponse will take a hex encoded hash and return a
// function that will validate that a given result has the same
// hash. The expected has does not start with 0x
func ValidateHashedResponse(expectedHash string) func(result interface{}) error {
	return func(result interface{}) error {
		resultStr, isValid := result.(string)
		if !isValid {
			return fmt.Errorf("Invalid result type. Expected string but got %T", result)
		}
		rawData, err := hex.DecodeString(resultStr[2:])
		if err != nil {
			return fmt.Errorf("The result string could be hex decoded: %w", err)
		}
		actualHash := fmt.Sprintf("%x", sha1.Sum(rawData))
		if actualHash != expectedHash {
			return fmt.Errorf("Hash mismatch expected: %s and got %s", expectedHash, actualHash)
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

// ValidateExact will validate against the exact value expected.
func ValidateExact(expected interface{}) func(result interface{}) error {
	return func(result interface{}) error {
		if expected != result {
			return fmt.Errorf("Expected %v and got %v", expected, result)
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
			return fmt.Errorf("Invalid result type. Expected string but got %T", result)
		}
		if !r.MatchString(resultStr) {
			return fmt.Errorf("The regex %s failed to match result %s", regEx, resultStr)
		}
		return nil
	}
}

// ValidateError will check the error message text against the provide regular expression
func ValidateError(errorMessageRegex string) func(result interface{}) error {
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

// ArgsLatestBlockHash is meant to generate an argument with the
// latest block hash for testing
func ArgsLatestBlockHash(cxt context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		blockData := make(map[string]interface{})
		err := rpcClient.CallContext(cxt, &blockData, "eth_getBlockByNumber", "latest", false)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retreive latest block hash")
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
		for _, v := range extraArgs {
			args = append(args, v)
		}
		return args
	}
}

// ArgsCoinbase would return arguments where the first argument is now
// the coinbase
func ArgsCoinbase(cxt context.Context, rpcClient *rpc.Client, extraArgs ...interface{}) func() []interface{} {
	return func() []interface{} {
		var coinbase string
		err := rpcClient.CallContext(cxt, &coinbase, "eth_coinbase")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retreive coinbase")
			return []interface{}{""}
		}
		log.Trace().Str("coinbase", coinbase).Msg("Got coinbase")

		args := []interface{}{coinbase}
		for _, v := range extraArgs {
			args = append(args, v)
		}
		return args
	}
}

// ArgsCoinbaseTransaction will return arguments where the from is replace with the current coinbase
func ArgsCoinbaseTransaction(cxt context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) func() []interface{} {
	return func() []interface{} {
		var coinbase string
		err := rpcClient.CallContext(cxt, &coinbase, "eth_coinbase")
		if err != nil {
			log.Error().Err(err).Msg("Unable to retreive coinbase")
			return []interface{}{""}
		}
		tx.From = coinbase
		log.Trace().Str("coinbase", coinbase).Msg("Got coinbase")
		return []interface{}{tx}
	}
}

// ArgsSignTransaction will take the junk transaction type that we've
// created, convert it to a geth style dynamic fee transaction and
// sign it with the user provide key.
func ArgsSignTransaction(cxt context.Context, rpcClient *rpc.Client, tx *RPCTestTransactionArgs) func() []interface{} {
	return func() []interface{} {
		ec := ethclient.NewClient(rpcClient)
		curNonce, err := ec.NonceAt(cxt, testEthAddress, nil)
		if err != nil {
			log.Error().Err(err).Msg("Unable to retreive nonce")
			curNonce = 0
		}
		log.Trace().Uint64("curNonce", curNonce).Msg("current nonce value")

		chainId, err := ec.ChainID(cxt)
		if err != nil {
			log.Error().Err(err).Msg("Unable to get chain id")
			chainId = big.NewInt(1)

		}
		log.Trace().Uint64("chainId", chainId.Uint64()).Msg("fetch chainid")

		dft := ethtypes.DynamicFeeTx{}
		dft.ChainID = chainId
		dft.Nonce = curNonce
		dft.GasTipCap = hexutil.MustDecodeBig(tx.MaxPriorityFeePerGas)
		dft.GasFeeCap = hexutil.MustDecodeBig(tx.MaxFeePerGas)
		dft.Gas = hexutil.MustDecodeUint64(tx.Gas)
		toAddr := ethcommon.HexToAddress(tx.To)
		dft.To = &toAddr
		dft.Value = hexutil.MustDecodeBig(tx.Value)
		dft.Data = hexutil.MustDecode(tx.Data)

		londonSigner := ethtypes.NewLondonSigner(chainId)
		signedTx, err := ethtypes.SignNewTx(testPrivateKey, londonSigner, &dft)
		if err != nil {
			log.Fatal().Err(err).Msg("There was an issue signing the transaction")
		}
		stringTx, err := signedTx.MarshalBinary()
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to marshal binary for transaction")
		}

		return []interface{}{hexutil.Encode(stringTx)}
	}
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
	return r.IsError
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
	return r.IsError
}

var RPCFuzzCmd = &cobra.Command{
	Use:   "rpcfuzz http://localhost:8545",
	Short: "Continually run a variety of RPC calls and fuzzers",
	Long: `

This command will run a series of RPC calls against a given json rpc
endpoint. The idea is to be able to check for various features and
function to see if the RPC generally conforms to typical geth
standards for the RPC

Some setup might be neede depending on how you're testing. We'll
demonstrate with geth. In order to quickly test this, you can run geth
in dev mode:

# ./build/bin/geth --dev --dev.period 5 --http --http.addr localhost \
    --http.port 8545 \
    --http.api admin,debug,web3,eth,txpool,personal,miner,net \
    --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 \
    --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 \
    --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 \
    --dev.gaslimit 50000000

Once your Eth client is running and the RPC is functional, you'll need
to transfer some amount of ether to a known account that ca be used
for testing

# cast send --from "$(cast rpc --rpc-url localhost:8545 eth_coinbase | jq -r '.')" \
    --rpc-url localhost:8545 --unlocked --value 100ether \
    0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6

Then we might want to deploy some test smart contracts. For the
purposes of testing we'll our ERC20 contract:

# cast send --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
    --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
    --rpc-url localhost:8545 --create \
    "$(cat ./contracts/ERC20.bin)"

Once this has been completed this will be the address of the contract:
0x6fda56c57b0acadb96ed5624ac500c0429d59429

# docker run -v $PWD/contracts:/contracts ethereum/solc:stable --storage-layout /contracts/ERC20.sol

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
		log.Trace().Msg("Doing test setup")
		setupTests(cxt, rpcClient)

		for _, t := range allTests {
			log.Trace().Str("name", t.GetName()).Str("method", t.GetMethod()).Msg("Running Test")
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

		privateKey, err := ethcrypto.HexToECDSA(*testPrivateHexKey)
		if err != nil {
			log.Error().Err(err).Msg("Couldn't process the hex private key")
			return err
		}

		ethAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)
		log.Info().Str("ethAddress", ethAddress.String()).Msg("Loaded private key")

		testPrivateKey = privateKey
		testEthAddress = ethAddress

		return nil
	},
}

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flagSet := RPCFuzzCmd.PersistentFlags()
	testPrivateHexKey = flagSet.String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to sending transactions")
	testContractAddress = flagSet.String("contract-address", "0x6fda56c57b0acadb96ed5624ac500c0429d59429", "The address of a contract that can be used for testing")

}
