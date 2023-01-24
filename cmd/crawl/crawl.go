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
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"

	// "github.com/ethereum/go-ethereum/p2p/enr"
	// "github.com/maticnetwork/polygon-cli/hdwallet"
	"github.com/spf13/cobra"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type (
	crawlParams struct {
		Bootnodes *string
		FileName  *string
		// ParsedBootnodes []*enode.Node
		// PrivateKey      *ecdsa.PrivateKey
	}

	crawlSample struct{}
)

var (
	inputCrawlParams crawlParams
)

// crawlCmd represents the crawl command
var CrawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl a network",
	Long: `This is a basic function to crawl a network
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// mnemonic, err := hdwallet.NewMnemonic(*inputBootstrapNode, *inputMnemonicLang)
		// if err != nil {
		// 	return err
		// }
		// cmd.Println(mnemonic)
		// return nil

		rpc, err := ethrpc.DialContext(cmd.Context(), inputCrawlParams.URL.String())
		if err != nil {
			log.Error().Err(err).Msg("Unable to dial rpc")
			return err
		}
		rpc.SetHeader("Accept-Encoding", "identity")
		ec := ethclient.NewClient(rpc)

		log.Info().Msg("Starting crawl")

	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if strings.HasPrefix(*inputCrawlParams.Bootnodes, "enr:") {
			return fmt.Errorf("the bootnode address should start with `enr:`. Given: %s", *inputCrawlParams.Bootnodes)
		}
		return nil

	},
	Args: func(cmd *cobra.Command, args []string) error {
		// setLogLevel(inputLoadTestParams)
		if len(args) != 1 {
			return fmt.Errorf("need nodes file as argument")
		}

		var cfg discover.Config

		inputCrawlParams.FileName = &args[0]

		cfg.PrivateKey, _ = crypto.GenerateKey()

		bn, err := parseBootnodes(inputCrawlParams.Bootnodes)
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse bootnodes")
			return err
		}
		cfg.Bootnodes = bn
		return nil
	},
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

func parseBootnodes(bootnodes *string) ([]*enode.Node, error) {
	s := params.RinkebyBootnodes

	s = strings.Split(*bootnodes, ",")

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

func init() {
	cp := new(crawlParams)
	cp.Bootnodes = CrawlCmd.PersistentFlags().String("bootnodes", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")

	inputCrawlParams = *cp

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mnemonicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mnemonicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
