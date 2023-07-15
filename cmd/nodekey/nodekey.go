package nodekey

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

	_ "embed"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethenode "github.com/ethereum/go-ethereum/p2p/enode"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

// libp2p (substrate/avail) - https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md
// subkey node-key

// devp2p (eth/bor) - https://github.com/ethereum/devp2p/blob/master/enr.md
// bootnode -genkey

// pex - https://github.com/tendermint/tendermint/blob/f2a8f5e054cf99ebe246818bb6d71f41f9a30faa/types/node_id.go
//

const RSAKeypairBits = 2048

var (
	//go:embed usage.md
	usage                       string
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

// NodekeyCmd represents the nodekey command
var NodekeyCmd = &cobra.Command{
	Use:   "nodekey",
	Short: "Generate node keys for different blockchain clients and protocols.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		var nko nodeKeyOut
		var withSeed bool
		switch *inputNodeKeyProtocol {
		case "devp2p":
			var err error
			nko, err = generateDevp2pNodeKey()
			if err != nil {
				return err
			}
		case "seed-libp2p":
			withSeed = true
			fallthrough
		case "libp2p":
			keyType, err := keyTypeToInt(*inputNodeKeyType)
			if err != nil {
				return err
			}
			nko, err = generateLibp2pNodeKey(keyType, withSeed)
			if err != nil {
				return err
			}
		default:
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
		ok := slices.Contains(validProtocols, *inputNodeKeyProtocol)
		if !ok {
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
	case "", "ed25519":
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

// That function can generate seeded keys but it shouldn't be used for production environments.
// It was created to allow us to create keys that work with the avail light client's method of generating keys
// User shouldn't encounter these problems but devs, be aware of:
// - generating a seeded ECDSA key does not return the same key even though you use the same seed
// - generating a seeded secp256k1 key does not return the same key even though you use the same seed
// - it's not possible to generate a seeded rsa key, it returns an "unexpected EOF" error
func generateLibp2pNodeKey(keyType int, seed bool) (nodeKeyOut, error) {
	var nko nodeKeyOut
	reader := rand.Reader
	if seed {
		seedValue := *inputNodeKeySeed
		seedData := make([]byte, 64)
		binary.BigEndian.PutUint64(seedData, seedValue)
		buf := bytes.NewBuffer(seedData)
		reader = io.LimitReader(buf, 64)
		nko.Seed = seedValue
	}

	prvKey, _, err := libp2pcrypto.GenerateKeyPairWithReader(keyType, RSAKeypairBits, reader)
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

	nko.PublicKey = id.String()
	// half of the private key is the public key. Substrate doesn't handle this well and need just the 32 byte seed/private key
	// TODO: should we keep private key to 32 bytes length for all types?
	nko.PrivateKey = hex.EncodeToString(rawPrvKey[0:ed25519.PublicKeySize])
	nko.FullPrivateKey = hex.EncodeToString(rawPrvKey)
	return nko, nil
}

func init() {
	inputNodeKeyProtocol = NodekeyCmd.PersistentFlags().String("protocol", "devp2p", "devp2p|libp2p|pex|seed-libp2p")
	inputNodeKeyType = NodekeyCmd.PersistentFlags().String("key-type", "ed25519", "ed25519|secp256k1|ecdsa|rsa")
	inputNodeKeyIP = NodekeyCmd.PersistentFlags().StringP("ip", "i", "0.0.0.0", "The IP to be associated with this address")
	inputNodeKeyTCP = NodekeyCmd.PersistentFlags().IntP("tcp", "t", 30303, "The tcp Port to be associated with this address")
	inputNodeKeyUDP = NodekeyCmd.PersistentFlags().IntP("udp", "u", 0, "The udp Port to be associated with this address")
	inputNodeKeySign = NodekeyCmd.PersistentFlags().BoolP("sign", "s", false, "Should the node record be signed?")
	inputNodeKeySeed = NodekeyCmd.PersistentFlags().Uint64P("seed", "S", 271828, "A numeric seed value")
	inputNodeKeyMarshalProtobuf = NodekeyCmd.PersistentFlags().BoolP("marshal-protobuf", "m", false, "If true the libp2p key will be marshaled to protobuf format rather than raw")

	inputNodeKeyFile = NodekeyCmd.PersistentFlags().StringP("file", "f", "", "A file with the private nodekey in hex format")
}
