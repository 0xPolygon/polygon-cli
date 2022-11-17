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
	"sync"
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
		Threads   *uint
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
		if *inputDumpblocks.Threads == 0 {
			*inputDumpblocks.Threads = 1
		}

		var wg sync.WaitGroup
		log.Info().Uint("thread", *inputDumpblocks.Threads).Msg("thread count")
		var pool = make(chan bool, *inputDumpblocks.Threads)
		// TODO Support parallel execution
		// TODO Support retries when there is a failure
		start := inputDumpblocks.Start
		end := inputDumpblocks.End

		for start < end {
			rangeStart := start
			rangeEnd := rangeStart + int64(inputDumpblocks.BatchSize)

			if rangeEnd > end {
				rangeEnd = end
			}

			pool <- true
			wg.Add(1)
			log.Info().Int64("start", rangeStart).Int64("end", rangeEnd).Msg("getting range")
			go func() {
				defer wg.Done()
				for {
					failCount := 0
					blocks, err := getBlockRange(ctx, rangeStart, rangeEnd, ec, inputDumpblocks.URL)
					if err != nil {
						failCount = failCount + 1
						if failCount > 5 {
							log.Error().Int64("rangeStart", rangeStart).Int64("rangeEnd", rangeEnd).Msg("unable to fetch blocks")
							break
						}
						time.Sleep(5 * time.Second)
						continue
					}

					failCount = 0
					receipts, err := getReceipts(ctx, blocks, ec, inputDumpblocks.URL)
					if err != nil {
						failCount = failCount + 1
						if failCount > 5 {
							log.Error().Int64("rangeStart", rangeStart).Int64("rangeEnd", rangeEnd).Msg("unable to fetch receipts")
							break
						}
						time.Sleep(5 * time.Second)
						continue
					}

					err = writeResponses(blocks)
					if err != nil {
						log.Error().Err(err).Msg("error writing blocks")
					}
					err = writeResponses(receipts)
					if err != nil {
						log.Error().Err(err).Msg("error writing receipts")
					}

					break
				}
				<-pool
			}()
			start = rangeEnd
		}

		log.Info().Msg("finished requesting data starting to wait")
		wg.Wait()
		log.Info().Msg("done")

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
			return fmt.Errorf("the start and end parameters need to be positive")
		}
		if end < start {
			start, end = end, start
		}
		inputDumpblocks.URL = args[0]
		inputDumpblocks.Start = start
		inputDumpblocks.End = end
		// realistically, this probably shoudln't be bigger than 999. Most Providers seem to cap at 1000
		inputDumpblocks.BatchSize = 150

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpblocksCmd)

	inputDumpblocks.Threads = dumpblocksCmd.PersistentFlags().UintP("concurrency", "c", 1, "how many go routines to leverage")

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
		log.Error().Err(err).Msg("rpc issue fetching blocks")
		return nil, err
	}
	blocks := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			return nil, b.Error
		}
		blocks = append(blocks, b.Result.(*json.RawMessage))

	}
	return blocks, nil
}

type (
	simpleRPCTransaction struct {
		Hash string `json:"hash"`
	}
	simpleRPCBlock struct {
		Transactions []simpleRPCTransaction `json:"transactions"`
	}
)

func getReceipts(ctx context.Context, rawBlocks []*json.RawMessage, c *ethrpc.Client, url string) ([]*json.RawMessage, error) {
	txHashes := make([]string, 0)
	for _, rb := range rawBlocks {
		var sb simpleRPCBlock
		err := json.Unmarshal(*rb, &sb)
		if err != nil {
			return nil, err

		}
		for _, tx := range sb.Transactions {
			txHashes = append(txHashes, tx.Hash)
		}

	}
	if len(txHashes) == 0 {
		return nil, nil
	}

	blms := make([]ethrpc.BatchElem, 0)
	for _, tx := range txHashes {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx},
			Result: r,
			Error:  err,
		})
	}

	var start uint64 = 0
	for {
		last := false
		end := start + inputDumpblocks.BatchSize
		if int(end) >= len(blms) {
			last = true
			end = uint64(len(blms) - 1)
		}

		// json: cannot unmarshal object into Go value of type []rpc.jsonrpcMessage
		// The error occurs when we call batchcallcontext with a single transaction for some reason.
		// polycli dumpblocks -c 1 http://127.0.0.1:9209/ 34457958 34458108
		// To handle this i'm making an exception when start and end are equal to make a single call
		if start == end {
			err := c.CallContext(ctx, &blms[start].Result, "eth_getTransactionReceipt", blms[start].Args[0])
			if err != nil {
				log.Error().Err(err).Uint64("start", start).Uint64("end", end).Msg("rpc issue fetching single receipt")
				return nil, err
			}
			break
		}

		err := c.BatchCallContext(ctx, blms[start:end])
		if err != nil {
			log.Error().Err(err).Str("randtx", txHashes[0]).Uint64("start", start).Uint64("end", end).Msg("rpc issue fetching receipts")
			return nil, err
		}
		start = end
		if last {
			break
		}
	}

	receipts := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			log.Error().Err(b.Error).Msg("block resp err")
			return nil, b.Error
		}
		receipts = append(receipts, b.Result.(*json.RawMessage))
	}
	log.Info().Int("hashes", len(txHashes)).Int("receipts", len(receipts)).Msg("fetched tx receipts")
	return receipts, nil
}

func writeResponses(blocks []*json.RawMessage) error {
	for _, b := range blocks {
		fmt.Println(string(*b))
	}
	return nil
}
