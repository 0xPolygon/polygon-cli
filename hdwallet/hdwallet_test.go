package hdwallet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestBIP32Vec2(t *testing.T) {
	seed := "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	table := map[string][2]string{
		"m":                               {"xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB", "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"},
		"m/0":                             {"xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH", "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt"},
		"m/0/2147483647'":                 {"xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a", "xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9"},
		"m/0/2147483647'/1":               {"xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon", "xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef"},
		"m/0/2147483647'/1/2147483646'":   {"xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL", "xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc"},
		"m/0/2147483647'/1/2147483646'/2": {"xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt", "xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j"},
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

func TestBIP32Vec3(t *testing.T) {
	seed := "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	table := map[string][2]string{
		"m":    {"xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13", "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6"},
		"m/0'": {"xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y", "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L"},
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
func TestBIP32Vec4(t *testing.T) {
	seed := "3ddd5602285899a946114506157c7997e5444528f3003f6134712147db19b678"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}

	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	table := map[string][2]string{
		"m":       {"xpub661MyMwAqRbcGczjuMoRm6dXaLDEhW1u34gKenbeYqAix21mdUKJyuyu5F1rzYGVxyL6tmgBUAEPrEz92mBXjByMRiJdba9wpnN37RLLAXa", "xprv9s21ZrQH143K48vGoLGRPxgo2JNkJ3J3fqkirQC2zVdk5Dgd5w14S7fRDyHH4dWNHUgkvsvNDCkvAwcSHNAQwhwgNMgZhLtQC63zxwhQmRv"},
		"m/0'":    {"xpub69AUMk3qDBi3uW1sXgjCmVjJ2G6WQoYSnNHyzkmdCHEhSZ4tBok37xfFEqHd2AddP56Tqp4o56AePAgCjYdvpW2PU2jbUPFKsav5ut6Ch1m", "xprv9vB7xEWwNp9kh1wQRfCCQMnZUEG21LpbR9NPCNN1dwhiZkjjeGRnaALmPXCX7SgjFTiCTT6bXes17boXtjq3xLpcDjzEuGLQBM5ohqkao9G"},
		"m/0'/1'": {"xpub6BJA1jSqiukeaesWfxe6sNK9CCGaujFFSJLomWHprUL9DePQ4JDkM5d88n49sMGJxrhpjazuXYWdMf17C9T5XnxkopaeS7jGk1GyyVziaMt", "xprv9xJocDuwtYCMNAo3Zw76WENQeAS6WGXQ55RCy7tDJ8oALr4FWkuVoHJeHVAcAqiZLE7Je3vZJHxspZdFHfnBEjHqU5hG1Jaj32dVoS6XLT1"},
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

func TestPolyWalletSetters(t *testing.T) {
	seed := "9C9B913EB1B6254F4737CE947EFD16F16E916F9D6EE5C1102A2002E48D4C88BD"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	err = pw.SetPath("invalid path")
	if err == nil {
		t.Fatal("This path should fail")
	}
	newPath := "m/44'/60'/0'"
	err = pw.SetPath(newPath)
	assert.Equal(t, pw.derivationPath, newPath, "Paths should be equal")
	if err != nil {
		t.Fatalf("Failed to set path: %v", err)
	}
}

func TestNewPolyWalletSetIterations(t *testing.T) {
	seed := "9C9B913EB1B6254F4737CE947EFD16F16E916F9D6EE5C1102A2002E48D4C88BD"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}
	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	_ = pw.SetUseRawEntropy(true)
	err = pw.SetIterations(1)
	if err == nil {
		t.Fatalf("Set iteration should fail")
	}

	_ = pw.SetUseRawEntropy(false)
	err = pw.SetIterations(1)
	if err != nil {
		t.Fatalf("Failed to set iteration: %v", err)
	}
}

func TestNewPolyWalletFail(t *testing.T) {
	mnemonic := "invalid mnemonic"
	_, err := NewPolyWallet(mnemonic, "password")

	if err == nil {
		t.Fatalf("Should fail with invalid mnemonic.")
	}
}

func TestPolyWalletParseMnemonic(t *testing.T) {
	mnemonic := "bottom drive obey lake curtain smoke basket hold race lonely fit walk"
	pw, err := NewPolyWallet(mnemonic, "password")
	if err != nil {
		t.Fatalf("Failed to create new poly wallet: %v", err)
	}

	_ = pw.SetUseRawEntropy(false)
	err = pw.parseMnemonic()
	if err != nil {
		t.Fatalf("Failed to parse mnemonic %v", err)
	}

	_ = pw.SetUseRawEntropy(true)
	err = pw.parseMnemonic()
	if err != nil {
		t.Fatalf("Failed to parse mnemonic %v", err)
	}
}

func TestPolyWalletExportRootAddress(t *testing.T) {
	mnemonic := "bottom drive obey lake curtain smoke basket hold race lonely fit walk"
	pw, _ := NewPolyWallet(mnemonic, "password")

	_, err := pw.ExportRootAddress()
	if err != nil {
		t.Fatalf("Failed to export root address %v", err)
	}
}

func TestPolyWalletExportHDAddresses(t *testing.T) {
	mnemonic := "bottom drive obey lake curtain smoke basket hold race lonely fit walk"
	pw, _ := NewPolyWallet(mnemonic, "password")

	_, err := pw.ExportHDAddresses(2)
	if err != nil {
		t.Fatalf("Failed to export HD address %v", err)
	}
}

func TestPolyWalletGetKey(t *testing.T) {
	mnemonic := "bottom drive obey lake curtain smoke basket hold race lonely fit walk"
	pw, _ := NewPolyWallet(mnemonic, "password")

	_, err := pw.GetKey()
	if err != nil {
		t.Fatalf("Failed to getKey %v", err)
	}

	pw.derivationPath = "invalid derivation path"
	_, err2 := pw.GetKey()
	if err2 == nil {
		t.Fatalf("should fail with invalid path")
	}

}

func TestPolyWalletGetKeyForPath(t *testing.T) {
	mnemonic := "bottom drive obey lake curtain smoke basket hold race lonely fit walk"
	pw, _ := NewPolyWallet(mnemonic, "password")

	_, err := pw.GetKeyForPath("invalid path")
	if err == nil {
		t.Fatalf("should fail with invalid path")
	}
}

func TestNewMnemonic(t *testing.T) {
	_, err := NewMnemonic(12, "fake language")
	if err == nil {
		t.Fatalf("should not create mnemonic - unrecognized language")
	}

	_, err1 := NewMnemonic(0, "")
	if err1 == nil {
		t.Fatalf("should not create mnemonic - invalid word count")
	}

	_, err2 := NewMnemonic(12, "english")
	if err2 != nil {
		t.Fatalf("Failed to create new mnemonic%v", err)
	}
}

func TestGetPublicKeyFromSeed(t *testing.T) {
	seed := "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"
	rawSeed, err := hex.DecodeString(seed)
	if err != nil {
		t.Fatalf("Failed to parse seed hex: %v", err)
	}

	pw, _ := NewPolyWalletFromSeed(rawSeed)
	pw.rawSeed = rawSeed

	_, err2 := GetPublicKeyFromSeed(rawSeed, SignatureSecp256k1, false)
	if err2 != nil {
		t.Fatalf("Failed to get public key from seed SignatureSecp256k1: %v", err)
	}
	_, err3 := GetPublicKeyFromSeed(rawSeed, SignatureEd25519, false)
	if err3 != nil {
		t.Fatalf("Failed to get public key from seed SignatureSecp256k1: %v", err)
	}
	_, err4 := GetPublicKeyFromSeed(rawSeed, SignatureSr25519, false)
	if err4 != nil {
		t.Fatalf("Failed to get public key from seed SignatureSecp256k1: %v", err)
	}

}
