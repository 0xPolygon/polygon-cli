package dumpblocks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "embed"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/proto/gen/pb"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type (
	dumpblocksParams struct {
		RpcUrl             string
		Start              uint64
		End                uint64
		BatchSize          uint64
		Threads            uint
		ShouldDumpBlocks   bool
		ShouldDumpReceipts bool
		Filename           string
		Mode               string
		FilterStr          string
		filter             Filter
	}
	Filter struct {
		To   []string `json:"to"`
		From []string `json:"from"`
	}
)

var (
	//go:embed usage.md
	usage           string
	inputDumpblocks dumpblocksParams = dumpblocksParams{}
)

// dumpblocksCmd represents the dumpblocks command
var DumpblocksCmd = &cobra.Command{
	Use:   "dumpblocks start end",
	Short: "Export a range of blocks from a JSON-RPC endpoint.",
	Long:  usage,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ec, err := ethrpc.DialContext(ctx, inputDumpblocks.RpcUrl)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		log.Info().Uint("thread", inputDumpblocks.Threads).Msg("Thread count")
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

					blocks = filterBlocks(blocks)

					if inputDumpblocks.ShouldDumpBlocks {
						err = writeResponses(blocks, "block")
						if err != nil {
							log.Error().Err(err).Msg("Error writing blocks")
						}
					}

					if inputDumpblocks.ShouldDumpReceipts {
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

						err = writeResponses(receipts, "transaction")
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
		if len(args) < 2 {
			return fmt.Errorf("command needs at least two arguments. A start block and an end block")
		}

		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		end, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		if start < 0 || end < 0 {
			return fmt.Errorf("the start and end parameters need to be positive")
		}
		if end < start {
			start, end = end, start
		}

		inputDumpblocks.Start = uint64(start)
		inputDumpblocks.End = uint64(end)

		if inputDumpblocks.Threads == 0 {
			inputDumpblocks.Threads = 1
		}
		if !slices.Contains([]string{"json", "proto"}, inputDumpblocks.Mode) {
			return fmt.Errorf("output format must one of [json, proto]")
		}

		if err := json.Unmarshal([]byte(inputDumpblocks.FilterStr), &inputDumpblocks.filter); err != nil {
			return fmt.Errorf("could not unmarshal filter string")
		}

		// Make sure the filters are all lowercase.
		for i := 0; i < len(inputDumpblocks.filter.To); i++ {
			inputDumpblocks.filter.To[i] = strings.ToLower(inputDumpblocks.filter.To[i])
		}
		for i := 0; i < len(inputDumpblocks.filter.From); i++ {
			inputDumpblocks.filter.From[i] = strings.ToLower(inputDumpblocks.filter.From[i])
		}

		return nil
	},
}

func init() {
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.RpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	DumpblocksCmd.PersistentFlags().UintVarP(&inputDumpblocks.Threads, "concurrency", "c", 1, "how many go routines to leverage")
	DumpblocksCmd.PersistentFlags().BoolVarP(&inputDumpblocks.ShouldDumpBlocks, "dump-blocks", "B", true, "if the blocks will be dumped")
	DumpblocksCmd.PersistentFlags().BoolVar(&inputDumpblocks.ShouldDumpReceipts, "dump-receipts", true, "if the receipts will be dumped")
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.Filename, "filename", "f", "", "where to write the output to (default stdout)")
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.Mode, "mode", "m", "json", "the output format [json, proto]")
	DumpblocksCmd.PersistentFlags().Uint64VarP(&inputDumpblocks.BatchSize, "batch-size", "b", 150, "the batch size. Realistically, this probably shouldn't be bigger than 999. Most providers seem to cap at 1000.")
	DumpblocksCmd.PersistentFlags().StringVarP(&inputDumpblocks.FilterStr, "filter", "F", "{}", "filter output based on tx to and from, not setting a filter means all are allowed")
}

func checkFlags() error {
	// Check rpc url flag.
	if err := util.ValidateUrl(inputDumpblocks.RpcUrl); err != nil {
		return err
	}

	return nil
}

// writeResponses writes the data to either stdout or a file if one is provided.
// The message type can be either "block" or "transaction". The format of the
// output is either "json" or "proto" depending on the mode.
func writeResponses(msg []*json.RawMessage, msgType string) error {
	switch inputDumpblocks.Mode {
	case "json":
		if err := writeJSON(msg); err != nil {
			log.Error().Err(err).Msgf("Failed to write %s json", msgType)
		}
	case "proto":
		for _, b := range msg {
			var protoMsg proto.Message
			switch msgType {
			case "block":
				protoMsg = &pb.Block{}
			case "transaction":
				protoMsg = &pb.Transaction{}
			}

			err := protojson.Unmarshal(*b, protoMsg)
			if err != nil {
				log.Error().Err(err).RawJSON("msg", *b).Msgf("Failed to unmarshal json to %s proto", msgType)
				continue
			}

			out, err := proto.Marshal(protoMsg)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to marshal %s proto", msgType)
				continue
			}

			if err = writeProto(out); err != nil {
				log.Error().Err(err).Msgf("Failed to write %s proto", msgType)
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

// filterBlocks will filter blocks that having transactions with a matching to or
// from field. If the to or from is an empty slice, then it will match all.
func filterBlocks(blocks []*json.RawMessage) []*json.RawMessage {
	// No filtering is done if there filters are not set.
	if len(inputDumpblocks.filter.To) == 0 && len(inputDumpblocks.filter.From) == 0 {
		return blocks
	}

	filtered := []*json.RawMessage{}
	for _, msg := range blocks {
		var block rpctypes.RawBlockResponse
		if err := json.Unmarshal(*msg, &block); err != nil {
			log.Error().Bytes("block", *msg).Msg("Unable to unmarshal block")
			continue
		}

		for _, tx := range block.Transactions {
			if (len(inputDumpblocks.filter.To) > 0 && slices.Contains(inputDumpblocks.filter.To, strings.ToLower(string(tx.To)))) ||
				(len(inputDumpblocks.filter.From) > 0 && slices.Contains(inputDumpblocks.filter.From, strings.ToLower(string(tx.From)))) {
				filtered = append(filtered, msg)
			}
		}
	}

	return filtered
}
