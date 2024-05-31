package query

import (
	"fmt"

	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	queryParams struct {
		StartBlock uint64
		Amount     uint64
	}
)

var (
	inputQueryParams queryParams
)

var QueryCmd = &cobra.Command{
	Use:   "query [enode/enr]",
	Short: "Query block header(s) from node and prints the output.",
	Long: `Query header of single block or range of blocks given a single enode/enr.
	
This command will initially establish a handshake and exchange status message
from the peer. Then, it will query the node for block(s) given the start block
and the amount of blocks to query and print the results.`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if inputQueryParams.Amount < 1 {
			return fmt.Errorf("amount must be greater than 0")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		node, err := p2p.ParseNode(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse enode")
			return
		}

		var (
			hello  *p2p.Hello
			status *p2p.Status
			start  uint64 = inputQueryParams.StartBlock
			amount uint64 = inputQueryParams.Amount
		)

		conn, err := p2p.Dial(node)
		if err != nil {
			log.Error().Err(err).Msg("Dial failed")
			return
		}
		defer conn.Close()
		if hello, status, err = conn.Peer(); err != nil {
			log.Error().Err(err).Msg("Peer failed")
			return
		}

		log.Info().Interface("hello", hello).Interface("status", status).Msg("Peering messages received")
		log.Info().Uint64("start", start).Uint64("amount", amount).Msg("Requesting headers")

		// Handshake completed, now proceed to query headers
		if err = conn.QueryHeaders(start, amount); err != nil {
			log.Error().Err(err).Msg("Failed to request header(s)")
			return
		}

		headers, err := conn.ListenHeaders()
		if err != nil {
			log.Error().Err(err).Msg("Failed to listen for header(s)")
			return
		}

		// Verify requested headers
		if len(headers) != int(amount) {
			log.Error().Uint64("want", amount).Int("have", len(headers)).Msg("Received less headers than requested")
			return
		}

		var (
			headerStart uint64 = headers[0].Number.Uint64()
			headerEnd   uint64 = headers[len(headers)-1].Number.Uint64()
			end         uint64 = start + amount - 1
		)
		if headerStart != start || headerEnd != end {
			log.Error().Uint64("start", start).Uint64("end", end).Uint64("header start", headerStart).Uint64("header end", headerEnd).Msg("Received headers out of range")
			return
		}

		print(headers)
	},
}

// print simply prints necessary contents of headers assuming they're already verified
func print(headers eth.BlockHeadersRequest) {
	for _, header := range headers {
		log.Info().Uint64("number", header.Number.Uint64()).Str("hash", header.Hash().Hex()).Msg("Header")
	}
}

func init() {
	QueryCmd.Flags().Uint64VarP(&inputQueryParams.StartBlock, "start-block", "s", 0, "Block number to start querying from")
	QueryCmd.Flags().Uint64VarP(&inputQueryParams.Amount, "amount", "a", 1, "Amount of blocks to query")
	if err := QueryCmd.MarkFlagRequired("start-block"); err != nil {
		log.Error().Err(err).Msg("Failed to mark start-block as required persistent flag")
	}
}
