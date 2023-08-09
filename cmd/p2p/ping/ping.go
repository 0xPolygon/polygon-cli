package ping

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	pingParams struct {
		Threads    int
		OutputFile string
		NodesFile  string
		Listen     bool
	}
	pingNodeJSON struct {
		Record *enode.Node `json:"record"`
		Hello  *p2p.Hello  `json:"hello,omitempty"`
		Status *p2p.Status `json:"status,omitempty"`
		Error  string      `json:"error,omitempty"`
	}
	pingNodeSet map[enode.ID]pingNodeJSON
)

var (
	inputPingParams pingParams
)

var PingCmd = &cobra.Command{
	Use:   "ping [enode/enr or nodes file]",
	Short: "Ping node(s) and return the output.",
	Long: `Ping nodes by either giving a single enode/enr or an entire nodes file.

This command will establish a handshake and status exchange to get the Hello and
Status messages and output JSON. If providing a enode/enr rather than a node file,
then the connection will remain open by default (--listen=true), and you can see
other messages the peer sends (e.g. blocks, transactions, etc.).`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodes := []*enode.Node{}
		if inputSet, err := p2p.ReadNodeSet(args[0]); err == nil {
			nodes = inputSet.Nodes()
			inputPingParams.Listen = false
		} else if node, err := p2p.ParseNode(args[0]); err == nil {
			nodes = append(nodes, node)
		} else {
			return err
		}

		output := make(pingNodeSet)

		var (
			mutex sync.Mutex
			wg    sync.WaitGroup
		)

		wg.Add(len(nodes))
		sem := make(chan bool, inputPingParams.Threads)

		count := &p2p.MessageCount{}
		go p2p.LogMessageCount(count, time.NewTicker(time.Second))

		// Ping each node in the slice.
		for _, n := range nodes {
			sem <- true
			go func(node *enode.Node) {
				defer func() {
					<-sem
					wg.Done()
				}()

				var (
					hello  *p2p.Hello
					status *p2p.Status
					errStr string
				)

				conn, err := p2p.Dial(node)
				if err != nil {
					log.Error().Err(err).Msg("Dial failed")
				} else {
					defer conn.Close()
					if hello, status, err = conn.Peer(); err != nil {
						log.Error().Err(err).Msg("Peer failed")
					}

					log.Info().Interface("hello", hello).Interface("status", status).Msg("Peering messages received")
				}

				if err != nil {
					errStr = err.Error()
				} else if inputPingParams.Listen {
					// If the dial and peering were successful, listen to the peer for messages.
					if err := conn.ReadAndServe(count); err != nil {
						log.Error().Err(err).Msg("Received error")
					}
				}

				// Save the results to the output map.
				mutex.Lock()
				output[node.ID()] = pingNodeJSON{node, hello, status, errStr}
				mutex.Unlock()
			}(n)
		}
		wg.Wait()

		// Write the output.
		nodesJSON, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		if inputPingParams.OutputFile == "" {
			os.Stdout.Write(nodesJSON)
		} else if err := os.WriteFile(inputPingParams.OutputFile, nodesJSON, 0644); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	PingCmd.PersistentFlags().StringVarP(&inputPingParams.OutputFile, "output", "o", "", "Write ping results to output file. (default stdout)")
	PingCmd.PersistentFlags().IntVarP(&inputPingParams.Threads, "parallel", "p", 16, "How many parallel pings to attempt.")
	PingCmd.PersistentFlags().BoolVarP(&inputPingParams.Listen, "listen", "l", true,
		`Keep the connection open and listen to the peer. This only works if the first
argument is an enode/enr, not a nodes file.`)
}
