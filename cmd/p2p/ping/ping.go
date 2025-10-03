package ping

import (
	"crypto/ecdsa"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/p2p"
)

type (
	pingParams struct {
		Threads    int
		OutputFile string
		NodesFile  string
		Listen     bool
		Port       int
		Addr       net.IP
		KeyFile    string
		PrivateKey string
		EnableWit  bool

		privateKey *ecdsa.PrivateKey
	}
)

var (
	inputPingParams pingParams
)

var PingCmd = &cobra.Command{
	Use:   "ping [enode/enr or nodes file]",
	Short: "Ping node(s) and return the output.",
	Long: `Ping nodes by either giving a single enode/enr or an entire nodes file.

This command will establish a handshake and status exchange to get the Hello and
Status messages and output JSON. If providing a enode/enr rather than a nodes
file, then the connection will remain open by default (--listen=true), and you
can see other messages the peer sends (e.g. blocks, transactions, etc.).`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputPingParams.privateKey, err = p2p.ParsePrivateKey(inputPingParams.KeyFile, inputPingParams.PrivateKey)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		nodes := []*enode.Node{}
		if input, err := p2p.ReadNodeSet(args[0]); err == nil {
			nodes = input
		} else if node, err := p2p.ParseNode(args[0]); err == nil {
			nodes = append(nodes, node)
		} else {
			return err
		}

		output := make(p2p.NodeSet)

		var (
			mutex sync.Mutex
			wg    sync.WaitGroup
		)

		wg.Add(len(nodes))
		sem := make(chan bool, inputPingParams.Threads)

		count := &p2p.MessageCount{}
		go func() {
			ticker := time.NewTicker(2 * time.Second)
			for {
				<-ticker.C
				c := count.Load()
				if !c.IsEmpty() {
					log.Info().Interface("counts", c).Send()
					count.Clear()
				}
			}
		}()

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
				)

				opts := p2p.DialOpts{
					EnableWit:  inputPingParams.EnableWit,
					Port:       inputPingParams.Port,
					Addr:       inputPingParams.Addr,
					PrivateKey: inputPingParams.privateKey,
				}

				conn, err := p2p.Dial(node, opts)
				if err != nil {
					log.Error().Err(err).Msg("Dial failed")
				} else {
					defer conn.Close()
					if hello, status, err = conn.Peer(); err != nil {
						log.Error().Err(err).Msg("Peer failed")
					}

					log.Info().Interface("hello", hello).Interface("status", status).Msg("Peering messages received")
				}

				if err == nil && inputPingParams.Listen {
					// If the dial and peering were successful, listen to the peer for messages.
					if err = conn.ReadAndServe(count); err != nil {
						log.Error().Err(err).Msg("Received error")
					}
				}

				// Save the results to the output map.
				mutex.Lock()
				if _, ok := output[node.ID()]; !ok {
					output[node.ID()] = []p2p.NodeJSON{}
				}

				result := p2p.NodeJSON{
					URL:    node.URLv4(),
					Hello:  hello,
					Status: status,
					Time:   time.Now().Unix(),
				}
				if err != nil {
					result.Error = err.Error()
				}

				output[node.ID()] = append(output[node.ID()], result)
				mutex.Unlock()
			}(n)
		}
		wg.Wait()

		return p2p.WriteNodeSet(inputPingParams.OutputFile, output, true)
	},
}

func init() {
	f := PingCmd.Flags()
	f.StringVarP(&inputPingParams.OutputFile, "output", "o", "", "write ping results to output file (default stdout)")
	f.IntVarP(&inputPingParams.Threads, "parallel", "p", 16, "how many parallel pings to attempt")
	f.BoolVarP(&inputPingParams.Listen, "listen", "l", true,
		`keep connection open and listen to peer. This only works if first
argument is an enode/enr, not a nodes file`)
	f.BoolVarP(&inputPingParams.EnableWit, "wit", "w", false, "enable wit/1 capability")
	f.IntVarP(&inputPingParams.Port, "port", "P", 30303, "port for discovery protocol")
	f.IPVarP(&inputPingParams.Addr, "addr", "a", net.ParseIP("127.0.0.1"), "address to bind discovery listener")
	f.StringVarP(&inputPingParams.KeyFile, "key-file", "k", "", "private key file (cannot be set with --key)")
	f.StringVar(&inputPingParams.PrivateKey, "key", "", "hex-encoded private key (cannot be set with --key-file)")
	PingCmd.MarkFlagsMutuallyExclusive("key-file", "key")
}
