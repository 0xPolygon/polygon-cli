package nodekey

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/flag"
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
	//go:embed usage.md
	usage                       string
	inputNodeKeyProtocol        string
	inputNodeKeyType            string
	inputNodeKeyIP              string
	inputNodeKeyTCP             int
	inputNodeKeyUDP             int
	inputNodeKeyFile            string
	inputNodeKeyPrivateKey      string
	inputNodeKeySign            bool
	inputNodeKeySeed            uint64
	inputNodeKeyMarshalProtobuf bool
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
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputNodeKeyPrivateKey, err = flag.GetPrivateKey(cmd)
		if err != nil {
			return err
		}

		switch inputNodeKeyProtocol {
		case "devp2p":
			invalidFlags := []string{"seed", "marshal-protobuf"}
			return validateNodeKeyFlags(cmd, invalidFlags)
		case "libp2p":
			invalidFlags := []string{"file", "ip", "tcp", "udp", "sign", "seed"}
			return validateNodeKeyFlags(cmd, invalidFlags)
		case "seed-libp2p":
			invalidFlags := []string{"file", "ip", "tcp", "udp", "sign"}
			if err := validateNodeKeyFlags(cmd, invalidFlags); err != nil {
				return err
			}
			if inputNodeKeyType == "rsa" {
				return fmt.Errorf("the RSA key type doesn't support manual key seeding")
			}
			if inputNodeKeyType == "secp256k1" {
				return fmt.Errorf("the secp256k1 key type doesn't support manual key seeding")
			}
			return nil
		default:
			return fmt.Errorf("the protocol %s is not implemented", inputNodeKeyProtocol)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var nko nodeKeyOut
		var withSeed bool
		switch inputNodeKeyProtocol {
		case "devp2p":
			switch inputNodeKeyType {
			case "ed25519":
				var err error
				nko, err = generateDevp2pNodeKey()
				if err != nil {
					return err
				}
			case "secp256k1":
				secret := []byte(strings.TrimPrefix(inputNodeKeyPrivateKey, "0x"))
				secp256k1PrivateKey := generateSecp256k1PrivateKey(secret)
				if err := displayHeimdallV2PrivValidatorKey(secp256k1PrivateKey); err != nil {
					return err
				}
				return nil
			}

		case "seed-libp2p":
			withSeed = true
			fallthrough
		case "libp2p":
			keyType, err := keyTypeToInt(inputNodeKeyType)
			if err != nil {
				return err
			}
			nko, err = generateLibp2pNodeKey(keyType, withSeed)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s is not implemented yet", inputNodeKeyProtocol)
		}

		out, err := json.Marshal(nko)
		if err != nil {
			return fmt.Errorf("could not json marshal the key data %w", err)
		}
		fmt.Println(string(out))

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
		return fmt.Errorf("the flag %s is not valid with the %s protocol", invalidFlagName, inputNodeKeyProtocol)
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
	var nodeKey *ecdsa.PrivateKey
	var err error

	switch {
	case inputNodeKeyPrivateKey != "":
		privateKey := strings.TrimPrefix(inputNodeKeyPrivateKey, "0x")
		nodeKey, err = gethcrypto.HexToECDSA(privateKey)
		if err != nil {
			return nodeKeyOut{}, fmt.Errorf("could not create ECDSA private key from given value: %s: %w", inputNodeKeyPrivateKey, err)
		}
	case inputNodeKeyFile != "":
		nodeKey, err = gethcrypto.LoadECDSA(inputNodeKeyFile)
		if err != nil {
			return nodeKeyOut{}, fmt.Errorf("could not load ECDSA private key from file %s: %w", inputNodeKeyFile, err)
		}
	default:
		nodeKey, err = gethcrypto.GenerateKey()
		if err != nil {
			return nodeKeyOut{}, fmt.Errorf("could not generate ECDSA private key: %w", err)
		}
	}

	nko := nodeKeyOut{}
	nko.PublicKey = fmt.Sprintf("%x", gethcrypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
	prvKeyBytes := gethcrypto.FromECDSA(nodeKey)
	nko.PrivateKey = hex.EncodeToString(prvKeyBytes)

	ip := net.ParseIP(inputNodeKeyIP)
	n := gethenode.NewV4(&nodeKey.PublicKey, ip, inputNodeKeyTCP, inputNodeKeyUDP)

	if inputNodeKeySign {
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
		seedValue := inputNodeKeySeed
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
	if inputNodeKeyMarshalProtobuf {
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
	f := NodekeyCmd.Flags()
	f.StringVar(&inputNodeKeyPrivateKey, flag.PrivateKey, "", "use the provided private key (in hex format)")
	f.StringVarP(&inputNodeKeyFile, "file", "f", "", "a file with the private nodekey (in hex format)")
	NodekeyCmd.MarkFlagsMutuallyExclusive(flag.PrivateKey, "file")

	f.StringVar(&inputNodeKeyProtocol, "protocol", "devp2p", "devp2p|libp2p|pex|seed-libp2p")
	f.StringVar(&inputNodeKeyType, "key-type", "ed25519", "ed25519|secp256k1|ecdsa|rsa")
	f.StringVarP(&inputNodeKeyIP, "ip", "i", "0.0.0.0", "the IP to be associated with this address")
	f.IntVarP(&inputNodeKeyTCP, "tcp", "t", 30303, "the TCP port to be associated with this address")
	f.IntVarP(&inputNodeKeyUDP, "udp", "u", 0, "the UDP port to be associated with this address")
	f.BoolVarP(&inputNodeKeySign, "sign", "s", false, "sign the node record")
	f.Uint64VarP(&inputNodeKeySeed, "seed", "S", 271828, "a numeric seed value")
	f.BoolVarP(&inputNodeKeyMarshalProtobuf, "marshal-protobuf", "m", false, "marshal libp2p key to protobuf format instead of raw")
}
