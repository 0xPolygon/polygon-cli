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
	"math/big"
	"net/url"

	"github.com/maticnetwork/polygon-cli/jsonrpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	monitorStatus struct {
		ChainID   *big.Int
		HeadBlock *big.Int
		PeerCount uint64
		GasPrice  *big.Int

		Blocks            map[string]*jsonrpc.RawBlockResponse
		MaxBlockRetrieved *big.Int
	}
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor [rpc-url]",
	Short: "A simple terminal monitor for a blockchain",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := jsonrpc.NewClient()
		resps, err := c.MakeRequestBatch(args[0], []string{"eth_blockNumber", "net_version", "net_peerCount", "eth_gasPrice"}, [][]any{nil, nil, nil, nil})
		if err != nil {
			return err
		}
		ms := new(monitorStatus)

		ms.MaxBlockRetrieved = big.NewInt(0)
		ms.Blocks = make(map[string]*jsonrpc.RawBlockResponse, 0)
		ms.HeadBlock = jsonrpc.MustConvHexToBigInt(resps[0].Result)
		ms.ChainID = jsonrpc.MustConvHexToBigInt(resps[1].Result)
		ms.PeerCount = jsonrpc.MustConvHexToUint64(resps[2].Result)
		ms.GasPrice = jsonrpc.MustConvHexToBigInt(resps[3].Result)

		from := big.NewInt(0)
		from.Sub(ms.HeadBlock, big.NewInt(10))
		getBlockRange(from, ms.HeadBlock, c, args[0], ms)
		jsonData, _ := json.Marshal(ms)
		fmt.Println(string(jsonData))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
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
		return nil
	},
}

func getBlockRange(from, to *big.Int, c *jsonrpc.Client, url string, ms *monitorStatus) (any, error) {
	one := big.NewInt(1)
	methods := make([]string, 0)
	params := make([][]any, 0)
	for i := from; i.Cmp(to) != 1; i.Add(i, one) {
		methods = append(methods, "eth_getBlockByNumber")
		params = append(params, []any{"0x" + i.Text(16), true})
	}
	var resps []jsonrpc.RPCBlockResp
	err := c.MakeRequestBatchGenric(url, methods, params, &resps)
	if err != nil {
		return nil, err
	}
	for _, r := range resps {
		block := r.Result
		ms.Blocks[string(block.Number)] = &block
		bi, err := jsonrpc.ConvHexToBigInt(block.Number)
		if err != nil {
			// unclear why this would happen
			log.Error().Err(err).Msg("Could not convert block number")
		}
		if ms.MaxBlockRetrieved.Cmp(bi) == -1 {
			ms.MaxBlockRetrieved = bi

		}
	}
	return nil, nil
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
