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
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/spf13/cobra"
)

// rpcCmd represents the rpc command
var rpcCmd = &cobra.Command{
	Use:   "rpc URL method param_0 param_1 ... param_n",
	Short: "A simple wrapper for making RPC requests",
	Long: `
Use this function to make JSON-RPC calls.

## ETH Examples
rpc http://127.0.0.1:8541 eth_getBlockByNumber 0x10E true

rpc http://127.0.0.1:8541 eth_getBlockByHash 0x15206ab0a5b408214127f5c445a86b7cfe6ae48fdcd9172b14e013dae7a7f470 true

rpc http://127.0.0.1:8541 eth_getTransactionByHash 0x97c070cb07bfac783ca73f08fb5999ae1ab509bf644197ef4a2c4e4f4a3c1516
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ec, err := ethrpc.DialContext(ctx, args[0])
		if err != nil {
			return err
		}

		params := toGenericParams(args[2:])
		var res = new(json.RawMessage)
		err = ec.Call(res, args[1], params...)
		if err != nil {
			return err
		}
		body, err := res.MarshalJSON()
		if err != nil {
			fmt.Println("gyahhhhhhhh1")
			return err
		}
		fmt.Println(string(body))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("rpc needs at least two arguments. A URL and a method")
		}

		_, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		return nil
	},
}

func toGenericParams(args []string) []any {
	retArgs := make([]any, len(args))
	for k := range args {
		retArgs[k] = convertStringTypes(args[k])
	}
	return retArgs
}

func convertStringTypes(param string) any {
	lowerParam := strings.ToLower(param)

	if lowerParam == "true" {
		return true
	}
	if lowerParam == "false" {
		return false
	}
	if lowerParam == "null" {
		return nil
	}

	return param
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rpcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
