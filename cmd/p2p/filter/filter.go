package filter

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	filterParams struct {
		Threads   int
		Client    string
		NetworkID int
		NodesFile string
	}
)

var (
	inputFilterParams filterParams
)

// FilterCmd represents the filter command.
var FilterCmd = &cobra.Command{
	Use:   "filter [nodes file]",
	Short: "Filter nodes given a nodes file",
	Long:  `Given a nodes file, filter the output.`,
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputFilterParams.NodesFile = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputFilterParams.NodesFile)
		if err != nil {
			return err
		}

		for _, node := range inputSet.Nodes() {
			conn, err := p2p.Dial(node)
			if err != nil {
				log.Error().Err(err).Msg("Dial failed")
				return err
			}

			hello, status, err := conn.Peer(nil)
			if err != nil {
				log.Error().Err(err).Msg("Peer failed")
			}
			log.Debug().Interface("hello", hello).Interface("status", status).Msg("Message received")
		}

		return nil
	},
}
