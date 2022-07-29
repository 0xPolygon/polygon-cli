/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/maticnetwork/polygon-cli/jsonrpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	inputLoadTestParams loadTestParams
	OneETH              = big.NewInt(1000000000000000000)
)

// loadtestCmd represents the loadtest command
var loadtestCmd = &cobra.Command{
	Use:   "loadtest [options] rpc-endpoint",
	Short: "A simple script for quickly running a load test",
	Long:  `Loadtest gives us a simple way to run a generic load test against an eth/EVM style json RPC endpoint`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug().Msg("Starting Loadtest")
		log.Trace().Interface("Input Params", inputLoadTestParams).Msg("Params")
		c := jsonrpc.NewClient()
		c.SetTimeout(time.Duration(*inputLoadTestParams.Timeout) * time.Second)
		if *inputLoadTestParams.Auth != "" {
			c.SetAuth(*inputLoadTestParams.Auth)
		}
		if *inputLoadTestParams.Proxy != "" {
			c.SetProxy(*inputLoadTestParams.Proxy, *inputLoadTestParams.ProxyAuth)
		}
		c.SetKeepAlive(*inputLoadTestParams.KeepAlive)
		_, err := getInitialAccountValues(c)
		if err != nil {
			return err
		}

		c.SendTx(
			*inputLoadTestParams.PrivateKey,
			OneETH,
			inputLoadTestParams.CurrentGas,
			"0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d", // TODO
			*inputLoadTestParams.CurrentNonce,
			big.NewInt(int64(*inputLoadTestParams.ChainID)),
			inputLoadTestParams.URL.String(),
		)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		setLogLevel(inputLoadTestParams)
		if len(args) != 1 {
			return fmt.Errorf("Expected exactly one argument")
		}
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}
		if url.Scheme != "http" && url.Scheme != "https" {
			return fmt.Errorf("The scheme %s is not supported", url.Scheme)
		}
		inputLoadTestParams.URL = url
		return nil
	},
}

func setLogLevel(ltp loadTestParams) {
	verbosity := *ltp.Verbosity
	if verbosity < 100 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if verbosity < 200 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if verbosity < 300 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if verbosity < 400 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if verbosity < 500 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if verbosity < 600 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	if *ltp.PrettyLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Starting logger in console mode")
	} else {
		log.Debug().Msg("Starting logger in JSON mode")
	}
}

type (
	loadTestParams struct {
		Requests     *int64
		Concurrency  *int64
		TimeLimit    *int64
		Timeout      *int64
		PostFile     *string
		Verbosity    *int64
		Auth         *string
		Proxy        *string
		ProxyAuth    *string
		KeepAlive    *bool
		PrettyLogs   *bool
		URL          *url.URL
		ChainID      *uint64
		PrivateKey   *string
		CurrentGas   *big.Int
		CurrentNonce *uint64
	}
)

func init() {
	rootCmd.AddCommand(loadtestCmd)

	ltp := new(loadTestParams)
	// Apache Bench Parameters
	ltp.Requests = loadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = loadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of multiple requests to perform at a time. Default is one request at a time.")
	ltp.TimeLimit = loadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no timelimit.")
	ltp.Timeout = loadtestCmd.PersistentFlags().Int64P("timeout", "s", 30, "Maximum number of seconds to wait before the socket times out. Default is 30 seconds.")
	ltp.PostFile = loadtestCmd.PersistentFlags().StringP("post-file", "p", "", "File containing data to POST.")
	// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
	ltp.Verbosity = loadtestCmd.PersistentFlags().Int64P("verbosity", "v", 200, "0 - Silent\n100 Fatals\n200 Errors\n300 Warnings\n400 INFO\n500 Debug\n600 Trace")
	ltp.Auth = loadtestCmd.PersistentFlags().StringP("auth", "A", "", "username:password used for www basic auth")
	ltp.Proxy = loadtestCmd.PersistentFlags().StringP("proxy", "X", "", "proxy:port combination to use a proxy server for the requests.")
	ltp.ProxyAuth = loadtestCmd.PersistentFlags().StringP("proxy-auth", "P", "", "Supply BASIC Authentication credentials to a proxy en-route. The username and password are separated by a single : and sent on the wire base64 encoded. The string is sent regardless of whether the proxy needs it (i.e., has sent an 407 proxy authentication needed).")
	ltp.KeepAlive = loadtestCmd.PersistentFlags().BoolP("keep-alive", "k", true, "Enable the HTTP KeepAlive feature, i.e., perform multiple requests within one HTTP session.")

	// extended parameters
	ltp.PrettyLogs = loadtestCmd.PersistentFlags().Bool("pretty-logs", true, "Should we log in pretty format or JSON")
	ltp.PrivateKey = loadtestCmd.PersistentFlags().String("private-key", "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa", "The hex encoded private key that we'll use to sending transactions")
	ltp.ChainID = loadtestCmd.PersistentFlags().Uint64("chain-id", 1256, "The chain id for the transactions that we're going to send")

	inputLoadTestParams = *ltp

	// TODO batch size
	// TODO Compression
}

func getInitialAccountValues(c *jsonrpc.Client) (interface{}, error) {
	resp, err := c.MakeRequest(inputLoadTestParams.URL.String(), "eth_gasPrice", nil)
	if err != nil {
		return nil, err
	}
	log.Trace().Interface("current gas price", resp.Result).Msg("Retreived current gas price")

	gasHexString, ok := resp.Result.(string)
	if !ok {
		return nil, fmt.Errorf("Could not assert %v as a string", resp.Result)
	}
	rawGas, err := hex.DecodeString(gasHexString[2:])
	if err != nil {
		return nil, err
	}
	gas := big.NewInt(0)
	gas.SetBytes(rawGas)
	log.Trace().Interface("current gas price big int", gas).Msg("Converted gas to big int")

	privateKey, err := ethcrypto.HexToECDSA(*inputLoadTestParams.PrivateKey)
	if err != nil {
		return nil, err
	}
	ethAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	address := "0xa0ebe20d02245b6540ae2c16c695dc815ea38f7e"
	resp, err = c.MakeRequest(inputLoadTestParams.URL.String(), "eth_getTransactionCount", []any{ethAddress.Hex(), "latest"})
	if err != nil {
		return nil, err
	}
	log.Trace().Interface("count", resp.Result).Str("address", address).Msg("Retrieved the current transaction count")
	var nonce uint64 = 2
	if resp.Result == nil {
		nonce = 1
	} else {
		nonceHexString, ok := resp.Result.(string)
		if !ok {
			return nil, fmt.Errorf("Could not assert %v as a string", resp.Result)
		}
		nonce = hex2int(nonceHexString)
	}

	resp, err = c.MakeRequest(inputLoadTestParams.URL.String(), "eth_getBalance", []any{address, "latest"})
	if err != nil {
		return nil, err
	}
	log.Trace().Interface("balance", resp.Result).Msg("Current account balance")

	inputLoadTestParams.CurrentGas = gas
	inputLoadTestParams.CurrentNonce = &nonce

	return nil, nil
}

func hex2int(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}
