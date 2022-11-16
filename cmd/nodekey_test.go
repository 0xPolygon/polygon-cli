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
				PublicKey:      "Qme51zsx66xx4zAo7KZoRKf45hQNEX3jC9Wb3ghUi5GduZ",
				PrivateKey:     "30770201010420000425d3fffbda2c000000000001164e43c6b2acbc73b3c524",
				FullPrivateKey: "30770201010420000425d3fffbda2c000000000001164e43c6b2acbc73b3c524c9f680144089c1a00a06082a8648ce3d030107a144034200047f2634af6d55f95d23b41e6f65e46e5f604977d0610df0b1bb9078a4a512f8bcaefadd9bdded533bca2c5c300d838ca894cced1d1006045f51025109fda29688",
			},
		},
		{
			name:     "rsa key with default seed",
			param:    libp2pcrypto.RSA,
			expected: nodeKeyOut{},
		},
	}

	for _, test := range tests {
		res, err := generateSeededLibp2pNodeKey(test.param)
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
