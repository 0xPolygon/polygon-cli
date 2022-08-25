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
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	dumpblocksArgs struct {
		URL       string
		Start     int64
		End       int64
		BatchSize uint64
	}
)

var inputDumpblocks dumpblocksArgs = dumpblocksArgs{}

// dumpblocksCmd represents the dumpblocks command
var dumpblocksCmd = &cobra.Command{
	Use:   "dumpblocks URL start end",
	Short: "Export a range of blocks from an RPC endpoint",
	Long:  `This is a simple function to export a range of blocks from a JSON RPC endpoint.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ec, err := ethrpc.DialContext(ctx, args[0])
		if err != nil {
			return err
		}

		// TODO Support parallel execution
		// TODO Support retries when there is a failure
		start := inputDumpblocks.Start
		end := inputDumpblocks.End
		failCount := 0
		for start < end {
			rangeStart := start
			rangeEnd := rangeStart + int64(inputDumpblocks.BatchSize)
			log.Info().Int64("start", rangeStart).Int64("end", rangeEnd).Msg("getting range")

			if rangeEnd > end {
				rangeEnd = end
			}
			blocks, err := getBlockRange(ctx, rangeStart, rangeEnd, ec, inputDumpblocks.URL)
			if err != nil {
				failCount = failCount + 1
				if failCount > 5 {
					return fmt.Errorf("Failed to get blockrange(%d - %d) after %d attempts", rangeStart, rangeEnd, failCount)
				}
				time.Sleep(5 * time.Second)
				continue
			}
			err = writeBlocks(blocks)
			if err != nil {
				return err
			}
			failCount = 0
			start = rangeEnd

		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("command needs at least three arguments. A URL a start block and an end block")
		}

		_, err := url.Parse(args[0])
		if err != nil {
			return err
		}
		start, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		end, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}
		if start < 0 || end < 0 {
			return fmt.Errorf("The start and end parameters need to be positive")
		}
		if end < start {
			start, end = end, start
		}
		inputDumpblocks.URL = args[0]
		inputDumpblocks.Start = start
		inputDumpblocks.End = end
		inputDumpblocks.BatchSize = 500

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpblocksCmd)
}

func getBlockRange(ctx context.Context, from, to int64, c *ethrpc.Client, url string) ([]*json.RawMessage, error) {
	blms := make([]ethrpc.BatchElem, 0)
	for i := from; i <= to; i = i + 1 {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + strconv.FormatInt(i, 16), true},
			Result: r,
			Error:  err,
		})
	}
	err := c.BatchCallContext(ctx, blms)
	if err != nil {
		return nil, err
	}
	blocks := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			return nil, err
		}
		blocks = append(blocks, b.Result.(*json.RawMessage))

	}
	return blocks, nil
}

func writeBlocks(blocks []*json.RawMessage) error {
	for _, b := range blocks {
		fmt.Println(string(*b))
	}
	return nil
}
