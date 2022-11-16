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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// libp2p (substrate/avail) - https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md
// subkey node-key

// devp2p (eth/bor) - https://github.com/ethereum/devp2p/blob/master/enr.md
// bootnode -genkey

// pex - https://github.com/tendermint/tendermint/blob/f2a8f5e054cf99ebe246818bb6d71f41f9a30faa/types/node_id.go
//

const RSAKeypairBits = 2048

var (
	inputNodeKeyProtocol        *string
	inputNodeKeyType            *string
	inputNodeKeyIP              *string
	inputNodeKeyTCP             *int
	inputNodeKeyUDP             *int
	inputNodeKeyFile            *string
	inputNodeKeySign            *bool
	inputNodeKeySeed            *uint64
	inputNodeKeyMarshalProtobuf *bool
)

type (
	nodeKeyOut struct {
		PublicKey      string
		PrivateKey     string
		FullPrivateKey string `json:",omitempty"`
		ENR            string `json:",omitempty"`
		Seed           uint64 `json:",omitempty"`
	}
)

// nodekeyCmd represents the nodekey command
var nodekeyCmd = &cobra.Command{
	Use:   "nodekey",
	Short: "Generate Node Keys",
	Long: `This is meant to be a simple utility for generating node keys for
different block chain clients and protocols.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var nko nodeKeyOut
		if *inputNodeKeyProtocol == "devp2p" {
			var err error
			nko, err = generateDevp2pNodeKey()
			if err != nil {
				return err
			}
		} else if *inputNodeKeyProtocol == "libp2p" {
			keyType, err := keyTypeToInt(*inputNodeKeyType)
			if err != nil {
				return err
			}
			nko, err = generateLibp2pNodeKey(keyType)
			if err != nil {
				return err
			}
		} else if *inputNodeKeyProtocol == "seed-libp2p" {
			keyType, err := keyTypeToInt(*inputNodeKeyType)
			if err != nil {
				return err
			}
			nko, err = generateSeededLibp2pNodeKey(keyType)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s is not implemented yet", *inputNodeKeyProtocol)
		}

		out, err := json.Marshal(nko)
		if err != nil {
			return fmt.Errorf("could not json marshal the key data %w", err)
		}
		fmt.Println(string(out))

		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("this command expects no arguments")
		}
		validProtocols := []string{"devp2p", "libp2p", "seed-libp2p"}
		isValidProtocol := false
		for _, p := range validProtocols {
			if p == *inputNodeKeyProtocol {
				isValidProtocol = true
			}
		}
		if !isValidProtocol {
			return fmt.Errorf("the protocol %s is not implemented", *inputNodeKeyProtocol)
		}
		if *inputNodeKeyProtocol == "devp2p" {
			invalidFlags := []string{"key-type", "seed", "marshal-protobuf"}
			err := validateNodeKeyFlags(cmd, invalidFlags)
			if err != nil {
				return err
			}
		}
		if *inputNodeKeyProtocol == "libp2p" {
			invalidFlags := []string{"file", "ip", "tcp", "udp", "sign", "seed"}
			err := validateNodeKeyFlags(cmd, invalidFlags)
			if err != nil {
				return err
			}
		}
		if *inputNodeKeyProtocol == "seed-libp2p" {
			invalidFlags := []string{"file", "ip", "tcp", "udp", "sign"}
			err := validateNodeKeyFlags(cmd, invalidFlags)
			if err != nil {
				return err
			}
			if *inputNodeKeyType == "rsa" {
				return fmt.Errorf("the RSA key type doesn't support manual key seeding")
			}
			if *inputNodeKeyType == "secp256k1" {
				return fmt.Errorf("the secp256k1 key type doesn't support manual key seeding")
			}
		}
		return nil
	},
}

func validateNodeKeyFlags(cmd *cobra.Command, invalidFlags []string) error {
	invalidFlagName := ""
	cmd.Flags().Visit(func(f *pflag.Flag) {
		for _, i := range invalidFlags {
			if f.Name == i {
				invalidFlagName = i
			}
		}
	})
	if invalidFlagName != "" {
		return fmt.Errorf("the flag %s is not valid with the %s protocol", invalidFlagName, *inputNodeKeyProtocol)
	}
	return nil
}

func keyTypeToInt(keyType string) (int, error) {
	// https://pkg.go.dev/github.com/libp2p/go-libp2p/core/crypto#pkg-constants
	switch keyType {
	case "":
		fallthrough
	case "ed25519":
		return libp2pcrypto.Ed25519, nil
	case "secp256k1":
		return libp2pcrypto.Secp256k1, nil
	case "ecdsa":
		return libp2pcrypto.ECDSA, nil
	case "rsa":
		return libp2pcrypto.RSA, nil
	default:
		return 0, fmt.Errorf("key type not implemented: %v", keyType)
	}
}

func generateDevp2pNodeKey() (nodeKeyOut, error) {
	nodeKey, err := gethcrypto.GenerateKey()

	if *inputNodeKeyFile != "" {
		nodeKey, err = gethcrypto.LoadECDSA(*inputNodeKeyFile)
	}
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("could not generate key: %w", err)
	}

	nko := nodeKeyOut{}
	nko.PublicKey = fmt.Sprintf("%x", gethcrypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
	prvKeyBytes := gethcrypto.FromECDSA(nodeKey)
	nko.PrivateKey = hex.EncodeToString(prvKeyBytes)

	ip := net.ParseIP(*inputNodeKeyIP)
	n := gethenode.NewV4(&nodeKey.PublicKey, ip, *inputNodeKeyTCP, *inputNodeKeyUDP)

	if *inputNodeKeySign {
		r := n.Record()
		err = gethenode.SignV4(r, nodeKey)
		if err != nil {
			return nodeKeyOut{}, err
		}
		n, err = gethenode.New(gethenode.ValidSchemes, r)
		if err != nil {
			return nodeKeyOut{}, err
		}
	}

	// ko.ENR = n.URLv4()
	nko.ENR = n.String()
	return nko, nil
}

func generateLibp2pNodeKey(keyType int) (nodeKeyOut, error) {
	prvKey, _, err := libp2pcrypto.GenerateKeyPairWithReader(keyType, RSAKeypairBits, rand.Reader)
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("unable to generate key pair, %w", err)
	}

	var rawPrvKey []byte
	if *inputNodeKeyMarshalProtobuf {
		rawPrvKey, err = libp2pcrypto.MarshalPrivateKey(prvKey)
	} else {
		rawPrvKey, err = prvKey.Raw()
	}
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("unable to convert the private key to a byte array, %w", err)
	}

	id, err := libp2ppeer.IDFromPrivateKey(prvKey)
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("unable to retrieve the node ID from the private key, %w", err)
	}

	return nodeKeyOut{
		PublicKey: id.String(),
		// half of the private key is the public key. Substrate doesn't handle this well and need just the 32 byte seed/private key
		// TODO: should we keep private key to 32 bytes length for all types?
		PrivateKey:     hex.EncodeToString(rawPrvKey[0:ed25519.PublicKeySize]),
		FullPrivateKey: hex.EncodeToString(rawPrvKey),
	}, nil
}

func generateSeededLibp2pNodeKey(keyType int) (nodeKeyOut, error) {
	seedValue := *inputNodeKeySeed
	seedData := make([]byte, 64)
	binary.BigEndian.PutUint64(seedData, seedValue)
	buf := bytes.NewBuffer(seedData)
	rand64 := io.LimitReader(buf, 64)

	prvKey, _, err := libp2pcrypto.GenerateKeyPairWithReader(keyType, RSAKeypairBits, rand64)
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("unable to generate key pair, %w", err)
	}

	var rawPrvKey []byte
	if *inputNodeKeyMarshalProtobuf {
		rawPrvKey, err = libp2pcrypto.MarshalPrivateKey(prvKey)
	} else {
		rawPrvKey, err = prvKey.Raw()
	}
	if err != nil {
		return nodeKeyOut{}, fmt.Errorf("unable to convert the private key to a byte array, %w", err)
	}

	id, err := libp2ppeer.IDFromPrivateKey(prvKey)
	if err != nil {
		return nodeKeyOut{}, err
	}
	return nodeKeyOut{
		PublicKey: id.String(),
		// half of the private key is the public key. Substrate doesn't handle this well and need just the 32 byte seed/private key
		// TODO: should we keep private key to 32 bytes length for all types?
		PrivateKey:     hex.EncodeToString(rawPrvKey[0:ed25519.PublicKeySize]),
		FullPrivateKey: hex.EncodeToString(rawPrvKey),
		Seed:           seedValue,
	}, nil
}

func init() {
	rootCmd.AddCommand(nodekeyCmd)

	inputNodeKeyProtocol = nodekeyCmd.PersistentFlags().String("protocol", "devp2p", "devp2p|libp2p|pex|seed-libp2p")
	inputNodeKeyType = nodekeyCmd.PersistentFlags().String("key-type", "ed25519", "ed25519|secp256k1|ecdsa|rsa")
	inputNodeKeyIP = nodekeyCmd.PersistentFlags().StringP("ip", "i", "0.0.0.0", "The IP to be associated with this address")
	inputNodeKeyTCP = nodekeyCmd.PersistentFlags().IntP("tcp", "t", 30303, "The tcp Port to be associated with this address")
	inputNodeKeyUDP = nodekeyCmd.PersistentFlags().IntP("udp", "u", 0, "The udp Port to be associated with this address")
	inputNodeKeySign = nodekeyCmd.PersistentFlags().BoolP("sign", "s", false, "Should the node record be signed?")
	inputNodeKeySeed = nodekeyCmd.PersistentFlags().Uint64P("seed", "S", 271828, "A numeric seed value")
	inputNodeKeyMarshalProtobuf = nodekeyCmd.PersistentFlags().BoolP("marshal-protobuf", "m", false, "If true the libp2p key will be marshaled to protobuf format rather than raw")

	inputNodeKeyFile = nodekeyCmd.PersistentFlags().StringP("file", "f", "", "A file with the private nodekey in hex format")
}
