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
package dumpblocks

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	dumpblocksArgs struct {
		URL          string
		Start        uint64
		End          uint64
		BatchSize    uint64
		Threads      *uint
		DumpBlocks   bool
		DumpReceipts bool
	}
)

var inputDumpblocks dumpblocksArgs = dumpblocksArgs{}

// dumpblocksCmd represents the dumpblocks command
var DumpblocksCmd = &cobra.Command{
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
		log.Info().Uint("thread", *inputDumpblocks.Threads).Msg("Thread count")
		var pool = make(chan bool, *inputDumpblocks.Threads)
		start := inputDumpblocks.Start
		end := inputDumpblocks.End

		for start < end {
			rangeStart := start
			rangeEnd := rangeStart + inputDumpblocks.BatchSize

			if rangeEnd > end {
				rangeEnd = end
			}

			pool <- true
			wg.Add(1)
			log.Info().Uint64("start", rangeStart).Uint64("end", rangeEnd).Msg("Getting range")
			go func() {
				defer wg.Done()
				for {
					failCount := 0
					blocks, err := util.GetBlockRange(ctx, rangeStart, rangeEnd, ec)
					if err != nil {
						failCount = failCount + 1
						if failCount > 5 {
							log.Error().Uint64("rangeStart", rangeStart).Uint64("rangeEnd", rangeEnd).Msg("Unable to fetch blocks")
							break
						}
						time.Sleep(5 * time.Second)
						continue
					}

					failCount = 0
					receipts, err := util.GetReceipts(ctx, blocks, ec, inputDumpblocks.BatchSize)
					if err != nil {
						failCount = failCount + 1
						if failCount > 5 {
							log.Error().Uint64("rangeStart", rangeStart).Uint64("rangeEnd", rangeEnd).Msg("Unable to fetch receipts")
							break
						}
						time.Sleep(5 * time.Second)
						continue
					}

					if inputDumpblocks.DumpBlocks {
						err = writeResponses(blocks)
						if err != nil {
							log.Error().Err(err).Msg("error writing blocks")
						}
					}

					if inputDumpblocks.DumpReceipts {
						err = writeResponses(receipts)
						if err != nil {
							log.Error().Err(err).Msg("error writing receipts")
						}
					}

					break
				}
				<-pool
			}()
			start = rangeEnd
		}

		log.Info().Msg("Finished requesting data starting to wait")
		wg.Wait()
		log.Info().Msg("Done")

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
		inputDumpblocks.Start = uint64(start)
		inputDumpblocks.End = uint64(end)
		// realistically, this probably shouldn't be bigger than 999. Most Providers seem to cap at 1000
		inputDumpblocks.BatchSize = 150

		return nil
	},
}

func init() {
	inputDumpblocks.Threads = DumpblocksCmd.PersistentFlags().UintP("concurrency", "c", 1, "how many go routines to leverage")
	inputDumpblocks.DumpBlocks = *DumpblocksCmd.PersistentFlags().BoolP("dump-blocks", "b", true, "if the blocks will be dumped")
	inputDumpblocks.DumpReceipts = *DumpblocksCmd.PersistentFlags().BoolP("dump-receipts", "r", true, "if the receipts will be dumped")
}

func writeResponses(blocks []*json.RawMessage) error {
	for _, b := range blocks {
		fmt.Println(string(*b))
	}
	return nil
}
