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

package crawl

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	crawlParams struct {
		Bootnodes string
		Timeout   string
		Threads   int
		Client    string
		NetworkID int
		NodesFile string
	}
)

var (
	inputCrawlParams crawlParams
)

// crawlCmd represents the crawl command
var CrawlCmd = &cobra.Command{
	Use:   "crawl [nodes file]",
	Short: "Crawl a network",
	Long:  `This is a basic function to crawl a network.`,
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputCrawlParams.NodesFile = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputCrawlParams.NodesFile)
		if err != nil {
			return err
		}

		var cfg discover.Config
		cfg.PrivateKey, _ = crypto.GenerateKey()
		bn, err := parseBootnodes(inputCrawlParams.Bootnodes)
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse bootnodes")
			return err
		}
		cfg.Bootnodes = bn

		db, err := enode.OpenDB("")
		if err != nil {
			return err
		}

		ln := enode.NewLocalNode(db, cfg.PrivateKey)
		socket, err := listen(ln)
		if err != nil {
			return err
		}

		disc, err := discover.ListenV4(socket, ln, cfg)
		if err != nil {
			return err
		}
		defer disc.Close()

		c := newCrawler(inputSet, disc, disc.RandomNodes())
		c.revalidateInterval = 10 * time.Minute

		timeout, err := time.ParseDuration(inputCrawlParams.Timeout)
		if err != nil {
			return err
		}

		log.Info().Msg("Starting crawl")

		output := c.run(timeout, inputCrawlParams.Threads)
		return p2p.WriteNodesJSON(inputCrawlParams.NodesFile, output)
	},
}

func init() {
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")
	CrawlCmd.MarkPersistentFlagRequired("bootnodes")
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Timeout, "timeout", "t", "30m0s", "Time limit for the crawl.")
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Client, "client", "c", "", "Name of client to filter the node information for.")
	CrawlCmd.PersistentFlags().IntVarP(&inputCrawlParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	CrawlCmd.PersistentFlags().IntVarP(&inputCrawlParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network id.")
}

func listen(ln *enode.LocalNode) (*net.UDPConn, error) {
	addr := "0.0.0.0:0"

	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		return nil, err
	}

	// Configure UDP endpoint in ENR from listener address.
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)

	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}

	ln.SetFallbackUDP(uaddr.Port)

	return usocket, nil
}

func decodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}

	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)

	return dec, err == nil
}

func decodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}

	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)

	return dec[:n], err == nil
}

// parseRecord parses a node record from hex, base64, or raw binary input.
func parseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)

	if d, ok := decodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := decodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}

	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)

	return &r, err
}

// parseNode parses a node record and verifies its signature.
func parseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}

	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}

	return enode.New(enode.ValidSchemes, r)
}

// parseBootnodes parses the bootnodes string and returns a node slice.
func parseBootnodes(bootnodes string) ([]*enode.Node, error) {
	s := strings.Split(bootnodes, ",")

	nodes := make([]*enode.Node, len(s))
	var err error
	for i, record := range s {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}

	return nodes, nil
}
