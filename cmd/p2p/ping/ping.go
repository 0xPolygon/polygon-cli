package ping

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

var PingCmd = &cobra.Command{
	Use:   "ping [enode/enr]",
	Short: "Ping a node given the enode or enr",
	Long:  `Pinging a node will return the Hello and Status messages.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		node, err := p2p.ParseNode(args[0])
		if err != nil {
			log.Error().Err(err).Msg("failed to parse node")
			return err
		}

		conn, err := p2p.Dial(node)
		if err != nil {
			log.Error().Err(err).Msg("dial failed")
			return err
		}

		hello, status, err := conn.Peer(nil)
		log.Debug().Interface("hello", hello).Interface("status", status).Msg("Message received")
		return err
	},
	Args: cobra.MinimumNArgs(1),
}
