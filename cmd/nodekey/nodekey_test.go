package nodekey

import (
	"testing"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func TestGenerateDevp2pNodeKey(t *testing.T) {
	res, err := generateDevp2pNodeKey()
	if err != nil {
		t.Errorf("could not generate eth key: %v", err)
	} else {
		if len(res.PublicKey) == 0 {
			t.Errorf("eth public key is empty")
		}
		if len(res.PrivateKey) == 0 {
			t.Errorf("eth private key is empty")
		}
		if len(res.ENR) == 0 {
			t.Errorf("enr is empty")
		}
	}
}

func TestGenerateLibp2pNodeKey(t *testing.T) {
	type test struct {
		name    string
		keyType int
	}

	tests := []test{
		{
			name:    "ed25519",
			keyType: libp2pcrypto.Ed25519,
		},
		{
			name:    "secp256k1",
			keyType: libp2pcrypto.Secp256k1,
		},
		{
			name:    "ecdsa",
			keyType: libp2pcrypto.ECDSA,
		},
		{
			name:    "rsa",
			keyType: libp2pcrypto.RSA,
		},
	}

	for _, test := range tests {
		res, err := generateLibp2pNodeKey(test.keyType, false)
		if err != nil {
			t.Errorf("could not generate %v key: %v", test.name, err)
		} else {
			if len(res.PublicKey) == 0 {
				t.Errorf("%v public key is empty", test.name)
			}
			if len(res.PrivateKey) == 0 {
				t.Errorf("%v private key is empty", test.name)
			}
			if len(res.FullPrivateKey) == 0 {
				t.Errorf("%v full private key is empty", test.name)
			}
		}
	}
}

func TestGenerateSeededLibp2pNodeKey(t *testing.T) {
	type test struct {
		name     string
		keyType  int
		expected nodeKeyOut
	}

	tests := []test{
		{
			name:    "ed25519 key with default seed",
			keyType: libp2pcrypto.Ed25519,
			expected: nodeKeyOut{
				PublicKey:      "12D3KooWMQMaVofHvQgjffhQGy3RmBRERFJxSJy59BgwTFnYASqX",
				PrivateKey:     "00000000000425d4000000000000000000000000000000000000000000000000",
				FullPrivateKey: "00000000000425d4000000000000000000000000000000000000000000000000ac25a2d49f6266b0a513cf0caf9ea45a9d74d74a1131d5530ac3291d70e81d7a",
			},
		},
	}

	for _, test := range tests {
		res, err := generateLibp2pNodeKey(test.keyType, true)
		if err != nil {
			t.Errorf("test %v: %v", test.name, err)
		} else {
			if res.PublicKey != test.expected.PublicKey {
				t.Errorf("test %v: public keys do not match", test.name)
			}
			if res.PrivateKey != test.expected.PrivateKey {
				t.Errorf("test %v: private keys do not match", test.name)
			}
			if res.FullPrivateKey != test.expected.FullPrivateKey {
				t.Errorf("test %v: full private keys do not match", test.name)
			}
		}
	}
}
