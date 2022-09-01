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
package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethenode "github.com/ethereum/go-ethereum/p2p/enode"

	//	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"

	"github.com/spf13/cobra"
)

// libp2p (substrate/avail) - https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md
// subkey node-key

// devp2p (eth/bor) - https://github.com/ethereum/devp2p/blob/master/enr.md
// bootnode -genkey

// pex - https://github.com/tendermint/tendermint/blob/f2a8f5e054cf99ebe246818bb6d71f41f9a30faa/types/node_id.go
//

var (
	inputNodeKeyProtocol *string
	inputNodeKeyType     *string
	inputNodeKeyIP       *string
	inputNodeKeyTCP      *int
	inputNodeKeyUDP      *int
	inputNodeKeyFile     *string
)

// nodekeyCmd represents the nodekey command
var nodekeyCmd = &cobra.Command{
	Use:   "nodekey",
	Short: "Generate Node Keys",
	Long: `This is meant to be a simple utility for generating node keys for
different block chain clients and protocols. Right now we've only
implemented devp2p because that's what we needed first.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if *inputNodeKeyProtocol == "devp2p" {
			generateETHNodeKey()
			return nil
		}
		return fmt.Errorf("%s is not implemented yet", *inputNodeKeyProtocol)
	},
}

type (
	nodeKeyOut struct {
		PublicKey  string
		PrivateKey string
		ENR        string
	}
)

func generateETHNodeKey() error {
	nodeKey, err := gethcrypto.GenerateKey()

	if *inputNodeKeyFile != "" {
		nodeKey, err = gethcrypto.LoadECDSA(*inputNodeKeyFile)
	}
	if err != nil {
		return fmt.Errorf("could not generate key: %v", err)
	}

	ko := nodeKeyOut{}
	ko.PublicKey = fmt.Sprintf("%x", gethcrypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
	prvKeyBytes := gethcrypto.FromECDSA(nodeKey)
	ko.PrivateKey = hex.EncodeToString(prvKeyBytes)

	ip := net.ParseIP(*inputNodeKeyIP)
	n := gethenode.NewV4(&nodeKey.PublicKey, ip, *inputNodeKeyTCP, *inputNodeKeyUDP)

	ko.ENR = n.String()

	out, err := json.Marshal(ko)
	if err != nil {
		return fmt.Errorf("could not json marshel the key data %v", err)
	}

	fmt.Println(string(out))
	return nil
}

func init() {
	rootCmd.AddCommand(nodekeyCmd)

	inputNodeKeyProtocol = nodekeyCmd.PersistentFlags().String("protocol", "devp2p", "devp2p|libp2p|pex")
	inputNodeKeyType = nodekeyCmd.PersistentFlags().String("key-type", "ed25519", "The type of key")
	inputNodeKeyIP = nodekeyCmd.PersistentFlags().StringP("ip", "i", "0.0.0.0", "The IP to be associated with this address")
	inputNodeKeyTCP = nodekeyCmd.PersistentFlags().IntP("tcp", "t", 30303, "The tcp Port to be associated with this address")
	inputNodeKeyUDP = nodekeyCmd.PersistentFlags().IntP("udp", "u", 0, "The udp Port to be associated with this address")

	inputNodeKeyFile = nodekeyCmd.PersistentFlags().StringP("file", "f", "", "A file with the private nodekey in hex format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodekeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodekeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
