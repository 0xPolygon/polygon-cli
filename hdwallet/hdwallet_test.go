package hdwallet

import (
	"encoding/hex"
	"testing"
)

// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
func TestBIP32(t *testing.T) {
	seed := "000102030405060708090a0b0c0d0e0f"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed
	key, err := pw.GetKeyForPath("m")
	if err != nil {
		t.Fatalf("Failed to get key for path 'm': %v", err)
	}

	if key.String() != "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi" {
		t.Fatalf("Bip 32 private key failed for chain: m")
	}

	if key.PublicKey().String() != "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8" {
		t.Fatalf("Bip 32 public key failed for chain: m")
	}
}
func TestBIP32Vec1(t *testing.T) {
	seed := "000102030405060708090a0b0c0d0e0f"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	table := map[string][2]string{
		"m":                      {"xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8", "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"},
		"m/0'":                   {"xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw", "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7"},
		"m/0'/1":                 {"xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ", "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs"},
		"m/0'/1/2'":              {"xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5", "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM"},
		"m/0'/1/2'/2":            {"xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV", "xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334"},
		"m/0'/1/2'/2/1000000000": {"xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy", "xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76"},
	}
	for k, v := range table {
		key, err := pw.GetKeyForPath(k)
		if err != nil {
			t.Fatalf("Unable to generate key for path %s %v", k, err)
		}
		prvData := key.String()
		pubData := key.PublicKey().String()
		if prvData != v[1] {
			t.Fatalf("Private key for path %s was mismatched. Expected %s got %s", k, v[1], prvData)
		}
		if pubData != v[0] {
			t.Fatalf("Public key for path %s was mismatched. Expected %s got %s", k, v[0], pubData)
		}

	}
}
