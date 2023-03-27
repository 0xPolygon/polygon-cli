package crawl

import (
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
)

// Hello is the RLP structure of the protocol handshake.
type Hello struct {
	Version    uint64
	Name       string
	Caps       []p2p.Cap
	ListenPort uint64
	ID         []byte // secp256k1 public key

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

func rlpxPing(n *enode.Node) (*Hello, error) {
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, err
	}

	conn := rlpx.NewConn(fd, n.Pubkey())
	ourKey, _ := crypto.GenerateKey()

	_, err = conn.Handshake(ourKey)
	if err != nil {
		return nil, err
	}

	code, data, _, err := conn.Read()
	if err != nil {
		return nil, err
	}

	switch code {
	case 0:
		var h Hello
		if err := rlp.DecodeBytes(data, &h); err != nil {
			return nil, fmt.Errorf("invalid handshake: %v", err)
		}
		fmt.Printf("%+v\n", h)
		return &h, nil

	case 1:
		var msg []p2p.DiscReason
		if rlp.DecodeBytes(data, &msg); len(msg) == 0 {
			return nil, fmt.Errorf("invalid disconnect message")
		}

		return nil, fmt.Errorf("received disconnect message: %v", msg[0])

	default:
		return nil, fmt.Errorf("invalid message code %d, expected handshake (code zero)", code)
	}

	return nil, nil
}
