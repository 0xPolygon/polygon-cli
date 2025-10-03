package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

func Listen(ln *enode.LocalNode) (*net.UDPConn, error) {
	socket, err := net.ListenPacket("udp4", "0.0.0.0:0")
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

// ParseNode parses a node record and verifies its signature.
func ParseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}

	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}

	return enode.New(enode.ValidSchemes, r)
}

// ParseBootnodes parses the bootnodes string and returns a node slice.
func ParseBootnodes(bootnodes string) ([]*enode.Node, error) {
	s := strings.Split(bootnodes, ",")

	nodes := make([]*enode.Node, len(s))
	var err error
	for i, record := range s {
		nodes[i], err = ParseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}

	return nodes, nil
}

// ParsePrivateKey loads a private key from a file path or hex string.
// If file is provided, it attempts to load from the file path first.
// If the file doesn't exist, it generates a new key and saves it to the file.
// If key is provided instead, it parses the hex-encoded private key string.
// If neither file nor key is provided, it generates and returns a new private key.
func ParsePrivateKey(file string, key string) (*ecdsa.PrivateKey, error) {
	if len(file) > 0 {
		privateKey, err := crypto.LoadECDSA(file)
		if err == nil {
			return privateKey, nil
		}

		log.Warn().Err(err).Msg("Key file was not found, generating a new key file")

		privateKey, err = crypto.GenerateKey()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate new private key")
			return nil, err
		}

		if err := crypto.SaveECDSA(file, privateKey); err != nil {
			log.Error().Err(err).Msg("Failed to save private key to file")
			return nil, err
		}

		return privateKey, nil
	}

	if len(key) > 0 {
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse private key")
			return nil, err
		}
		return privateKey, nil
	}

	return crypto.GenerateKey()
}
