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
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethenode "github.com/ethereum/go-ethereum/p2p/enode"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2ppeer "github.com/libp2p/go-libp2p/core/peer"

	//	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"

	"github.com/spf13/cobra"
)

// libp2p (substrate/avail) - https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md
// subkey node-key

// devp2p (eth/bor) - https://github.com/ethereum/devp2p/blob/master/enr.md
// bootnode -genkey

// pex - https://github.com/tendermint/tendermint/blob/f2a8f5e054cf99ebe246818bb6d71f41f9a30faa/types/node_id.go
//

const KEYPAIR_BITE_SIZE = 256

var (
	inputNodeKeyProtocol *string
	inputNodeKeyType     *string
	inputNodeKeyIP       *string
	inputNodeKeyTCP      *int
	inputNodeKeyUDP      *int
	inputNodeKeyFile     *string
	inputNodeKeySign     *bool
	inputNodeKeySeed     *uint64
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
			return generateETHNodeKey()
		}
		if *inputNodeKeyProtocol == "libp2p" {
			switch *inputNodeKeyType {
			// https://pkg.go.dev/github.com/libp2p/go-libp2p/core/crypto#pkg-constants
			case "rsa":
				return generateLibp2pNodeKey(libp2pcrypto.RSA)
			case "ed25519":
				return generateLibp2pNodeKey(libp2pcrypto.Ed25519)
			case "secp256k1":
				return generateLibp2pNodeKey(libp2pcrypto.Secp256k1)
			case "ecdsa":
				return generateLibp2pNodeKey(libp2pcrypto.ECDSA)
			default:
				return fmt.Errorf("key type not implemented %v", *inputNodeKeyType)
			}
		}
		if *inputNodeKeyProtocol == "seed-libp2p" {
			switch *inputNodeKeyType {
			// https://pkg.go.dev/github.com/libp2p/go-libp2p/core/crypto#pkg-constants
			case "rsa":
				return generateSeededLibp2pNodeKey(libp2pcrypto.RSA)
			case "ed25519":
				return generateSeededLibp2pNodeKey(libp2pcrypto.Ed25519)
			case "secp256k1":
				return generateSeededLibp2pNodeKey(libp2pcrypto.Secp256k1)
			case "ecdsa":
				return generateSeededLibp2pNodeKey(libp2pcrypto.ECDSA)
			default:
				return fmt.Errorf("key type not implemented %v", *inputNodeKeyType)
			}
		}

		return fmt.Errorf("%s is not implemented yet", *inputNodeKeyProtocol)
	},
}

type (
	nodeKeyOut struct {
		PublicKey      string
		PrivateKey     string
		FullPrivateKey string `json:",omitempty"`
		ENR            string `json:",omitempty"`
		Seed           uint64 `json:",omitempty"`
	}
)

func generateLibp2pNodeKey(keyType int) error {
	rand32 := io.LimitReader(rand.Reader, 32)
	prvKey, _, err := libp2pcrypto.GenerateKeyPairWithReader(keyType, KEYPAIR_BITE_SIZE, rand32)
	if err != nil {
		return fmt.Errorf("unable to generate key pair, %v", err)
	}

	rawPrvKey, err := libp2pcrypto.MarshalPrivateKey(prvKey)
	if err != nil {
		return fmt.Errorf("unable to convert the private key to a byte array, %v", err)
	}

	id, err := libp2ppeer.IDFromPrivateKey(prvKey)
	if err != nil {
		return fmt.Errorf("unable to retrieve the node ID from the private key, %v", err)
	}

	nko := nodeKeyOut{
		PublicKey: id.String(),
		// half of the private key is the public key. Substrate doesn't handle this well and need just the 32 byte seed/private key
		// TODO: should we keep private key to 32 bytes length for all types?
		PrivateKey:     hex.EncodeToString(rawPrvKey[0:ed25519.PublicKeySize]),
		FullPrivateKey: hex.EncodeToString(rawPrvKey),
	}

	out, err := json.Marshal(nko)
	if err != nil {
		return fmt.Errorf("could not json marshel the key data %v", err)
	}

	fmt.Println(string(out))

	return nil
}

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

	if *inputNodeKeySign {
		r := n.Record()
		err = gethenode.SignV4(r, nodeKey)
		if err != nil {
			return err
		}
		n, err = gethenode.New(gethenode.ValidSchemes, r)
		if err != nil {
			return err
		}
	}

	// ko.ENR = n.URLv4()
	ko.ENR = n.String()

	out, err := json.Marshal(ko)
	if err != nil {
		return fmt.Errorf("could not json marshel the key data %v", err)
	}

	fmt.Println(string(out))
	return nil
}

func generateSeededLibp2pNodeKey(keyType int) error {
	seedValue := *inputNodeKeySeed
	seedData := make([]byte, 32, 32)
	binary.BigEndian.PutUint64(seedData, seedValue)
	buf := bytes.NewBuffer(seedData)
	rand32 := io.LimitReader(buf, 32)

	prvKey, _, err := libp2pcrypto.GenerateKeyPairWithReader(keyType, KEYPAIR_BITE_SIZE, rand32)
	if err != nil {
		return fmt.Errorf("unable to generate key pair, %v", err)
	}

	rawPrvKey, err := libp2pcrypto.MarshalPrivateKey(prvKey)
	if err != nil {
		return fmt.Errorf("unable to convert the private key to a byte array, %v", err)
	}

	id, err := libp2ppeer.IDFromPrivateKey(prvKey)
	if err != nil {
		return err
	}
	nko := nodeKeyOut{
		PublicKey: id.String(),
		// half of the private key is the public key. Substrate doesn't handle this well and need just the 32 byte seed/private key
		// TODO: should we keep private key to 32 bytes length for all types?
		PrivateKey:     hex.EncodeToString(rawPrvKey[0:ed25519.PublicKeySize]),
		FullPrivateKey: hex.EncodeToString(rawPrvKey),
		Seed:           seedValue,
	}

	out, err := json.Marshal(nko)
	if err != nil {
		return fmt.Errorf("could not json marshel the key data %v", err)
	}

	fmt.Println(string(out))

	return nil
}
func init() {
	rootCmd.AddCommand(nodekeyCmd)

	inputNodeKeyProtocol = nodekeyCmd.PersistentFlags().String("protocol", "devp2p", "devp2p|libp2p|pex")
	inputNodeKeyType = nodekeyCmd.PersistentFlags().String("key-type", "ed25519", "rsa|ed25519|secp256k1|ecdsa")
	inputNodeKeyIP = nodekeyCmd.PersistentFlags().StringP("ip", "i", "0.0.0.0", "The IP to be associated with this address")
	inputNodeKeyTCP = nodekeyCmd.PersistentFlags().IntP("tcp", "t", 30303, "The tcp Port to be associated with this address")
	inputNodeKeyUDP = nodekeyCmd.PersistentFlags().IntP("udp", "u", 0, "The udp Port to be associated with this address")
	inputNodeKeySign = nodekeyCmd.PersistentFlags().BoolP("sign", "s", false, "Should the node record be signed?")
	inputNodeKeySeed = nodekeyCmd.PersistentFlags().Uint64P("seed", "S", 271828, "A numeric seed value")

	inputNodeKeyFile = nodekeyCmd.PersistentFlags().StringP("file", "f", "", "A file with the private nodekey in hex format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodekeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodekeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
