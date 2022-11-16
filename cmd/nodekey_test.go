package cmd

import (
	"testing"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func TestGenerateSeededLibp2pNodeKey(t *testing.T) {
	type test struct {
		name     string
		param    int
		expected nodeKeyOut
	}

	tests := []test{
		{
			name:  "ed25519 key with default seed",
			param: libp2pcrypto.Ed25519,
			expected: nodeKeyOut{
				PublicKey:      "12D3KooWMQMaVofHvQgjffhQGy3RmBRERFJxSJy59BgwTFnYASqX",
				PrivateKey:     "00000000000425d4000000000000000000000000000000000000000000000000",
				FullPrivateKey: "00000000000425d4000000000000000000000000000000000000000000000000ac25a2d49f6266b0a513cf0caf9ea45a9d74d74a1131d5530ac3291d70e81d7a",
			},
		},
		{
			name:  "secp256k1 key with default seed",
			param: libp2pcrypto.Secp256k1,
			expected: nodeKeyOut{
				PublicKey:      "16Uiu2HAm7WUXemKojg36e3P5pBwUAZoJCxBUmTXqeMMUaDxqkU6H",
				PrivateKey:     "08021220f59d925597656394ac83db5c280338c8c253982ffeea3f71f4dbe329",
				FullPrivateKey: "08021220f59d925597656394ac83db5c280338c8c253982ffeea3f71f4dbe329eae5638d",
			},
		},
		{
			name:  "ecdsa key with default seed",
			param: libp2pcrypto.ECDSA,
			expected: nodeKeyOut{
				PublicKey:      "",
				PrivateKey:     "",
				FullPrivateKey: "",
			},
		},
		{
			name:  "rsa key with default seed",
			param: libp2pcrypto.RSA,
			expected: nodeKeyOut{
				PublicKey:      "16Uiu2HAm7WUXemKojg36e3P5pBwUAZoJCxBUmTXqeMMUaDxqkU6H",
				PrivateKey:     "08021220f59d925597656394ac83db5c280338c8c253982ffeea3f71f4dbe329",
				FullPrivateKey: "08021220f59d925597656394ac83db5c280338c8c253982ffeea3f71f4dbe329eae5638d",
			},
		},
	}

	for _, test := range tests {
		res, err := generateSeededLibp2pNodeKey(test.param)
		if err != nil {
			t.Errorf("test %v: %v", test.name, err)
		}
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
