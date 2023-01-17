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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/proto/gen/pb"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type (
	dumpblocksArgs struct {
		URL          string
		Start        uint64
		End          uint64
		BatchSize    uint64
		Threads      uint
		DumpBlocks   bool
		DumpReceipts bool
		Filename     string
		Format       string
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
		if inputDumpblocks.Threads == 0 {
			inputDumpblocks.Threads = 1
		}

		var wg sync.WaitGroup
		log.Info().Uint("thread", inputDumpblocks.Threads).Msg("thread count")
		var pool = make(chan bool, inputDumpblocks.Threads)
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
						err = writeBlocks(blocks)
						if err != nil {
							log.Error().Err(err).Msg("Error writing blocks")
						}
					}

					if inputDumpblocks.DumpReceipts {
						err = writeTxs(receipts)
						if err != nil {
							log.Error().Err(err).Msg("Error writing receipts")
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

		if !slices.Contains([]string{"json", "proto"}, inputDumpblocks.Format) {
			return fmt.Errorf("output format must one of [json, proto]")
		}

		return nil
	},
}

func init() {
	DumpblocksCmd.PersistentFlags().UintVarP(&inputDumpblocks.Threads, "concurrency", "c", 1, "how many go routines to leverage")
	DumpblocksCmd.PersistentFlags().BoolVarP(&inputDumpblocks.DumpBlocks, "dump-blocks", "B", true, "if the blocks will be dumped")
	DumpblocksCmd.PersistentFlags().BoolVarP(&inputDumpblocks.DumpReceipts, "dump-receipts", "r", true, "if the receipts will be dumped")
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.Filename, "filename", "f", "", "where to write the output to (default stdout)")
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.Format, "format", "F", "json", "the output format [json, proto]")
	DumpblocksCmd.PersistentFlags().Uint64VarP(&inputDumpblocks.BatchSize, "batch-size", "b", 150, "the batch size. Realistically, this probably shouldn't be bigger than 999. Most providers seem to cap at 1000.")
}

// writeBlock writes the blocks.
func writeBlocks(msg []*json.RawMessage) error {
	switch inputDumpblocks.Format {
	case "json":
		if err := writeJSON(msg); err != nil {
			log.Error().Err(err).Msg("Failed to write block json")
		}
	case "proto":
		for _, b := range msg {
			protoMsg := &pb.Block{}
			err := protojson.Unmarshal(*b, protoMsg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal json to block proto")
				continue
			}

			out, err := proto.Marshal(protoMsg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal block proto")
				continue
			}

			if err = writeProto(out); err != nil {
				log.Error().Err(err).Msg("Failed to write block proto")
				continue
			}
		}
	}

	return nil
}

// writeTxs writes the transactions receipts.
func writeTxs(msg []*json.RawMessage) error {
	switch inputDumpblocks.Format {
	case "json":
		if err := writeJSON(msg); err != nil {
			log.Error().Err(err).Msg("Failed to write tx json")
		}
	case "proto":
		for _, b := range msg {
			protoMsg := &pb.Transaction{}
			err := protojson.Unmarshal(*b, protoMsg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal json to tx proto")
				continue
			}

			out, err := proto.Marshal(protoMsg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal tx proto")
				continue
			}

			if err = writeProto(out); err != nil {
				log.Error().Err(err).Msg("Failed to write tx proto")
				continue
			}
		}
	}

	return nil
}

// writeJSON writes the json raw messages to stdout by default and to a file if
// provided.
func writeJSON(msg []*json.RawMessage) error {
	f := os.Stdout
	if inputDumpblocks.Filename != "" {
		var err error
		f, err = os.OpenFile(inputDumpblocks.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
	}

	for _, b := range msg {
		fmt.Fprintln(f, string(*b))
	}

	return nil
}

// writeProto writes the buffer data to stdout by default and to a file if
// provided.
//
// It will write first the length of the buffer and then the buffer.
func writeProto(out []byte) error {
	f := os.Stdout
	// Open the file for writing if the filename is provided.
	if inputDumpblocks.Filename != "" {
		var err error
		f, err = os.OpenFile(inputDumpblocks.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
	}

	// Because protobuf isn't a self delimiting format, we write the length of the
	// bytes to the file as a header. This allows us to correctly read back in the
	// file.
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(out)))

	if _, err := f.Write(buf); err != nil {
		return err
	}

	if _, err := f.Write(out); err != nil {
		return err
	}

	return nil
}
