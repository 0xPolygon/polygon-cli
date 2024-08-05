package hdwallet

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/oasisprotocol/curve25519-voi/primitives/sr25519"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ripemd160" //nolint:staticcheck
	"golang.org/x/crypto/sha3"
)

type (
	PolySignature int

	PolyWallet struct {
		Mnemonic       string
		Passphrase     string
		derivationPath string
		rawSeed        []byte
		kdfIterations  uint
		keyCache       map[string]*bip32.Key
		useRawEntropy  bool
	}
	PolyWalletExport struct {
		RootKey           string               `json:",omitempty"`
		Seed              string               `json:",omitempty"`
		Mnemonic          string               `json:",omitempty"`
		Passphrase        string               `json:",omitempty"`
		DerivationPath    string               `json:",omitempty"`
		AccountPublicKey  string               `json:",omitempty"`
		AccountPrivateKey string               `json:",omitempty"`
		BIP32PublicKey    string               `json:",omitempty"`
		BIP32PrivateKey   string               `json:",omitempty"`
		Addresses         []*PolyAddressExport `json:",omitempty"`
		multiAddress
	}
	multiAddress struct {
		HexPublicKey       string `json:",omitempty"`
		HexFullPublicKey   string `json:",omitempty"`
		HexPrivateKey      string `json:",omitempty"`
		ETHAddress         string `json:",omitempty"`
		BTCAddress         string `json:",omitempty"`
		WIF                string `json:",omitempty"`
		ECDSAAddress       string `json:",omitempty"`
		Sr25519Address     string `json:",omitempty"`
		Ed25519Address     string `json:",omitempty"`
		ECDSAAddressSS58   string `json:",omitempty"`
		Sr25519AddressSS58 string `json:",omitempty"`
		Ed25519AddressSS58 string `json:",omitempty"`
	}
	PolyAddressExport struct {
		Path string `json:",omitempty"`
		multiAddress
	}
)

var (
	wordsToBits    = map[int]int{12: 128, 15: 160, 18: 192, 21: 224, 24: 256}
	langToWordlist = map[string][]string{
		"chinesesimplified":  wordlists.ChineseSimplified,
		"chinesetraditional": wordlists.ChineseTraditional,
		"czech":              wordlists.Czech,
		"english":            wordlists.English,
		"french":             wordlists.French,
		"italian":            wordlists.Italian,
		"japanese":           wordlists.Japanese,
		"korean":             wordlists.Korean,
		"spanish":            wordlists.Spanish,
	}
	pathValidator   = `^m[\/0-9']*[0-9']$`
	rePathValidator *regexp.Regexp
)

const (
	SignatureSecp256k1 PolySignature = iota
	SignatureEd25519
	SignatureSr25519
)

func GenPrivKeyFromSecret(seed []byte, c PolySignature) (interface{}, error) {
	if c == SignatureEd25519 {
		return ed25519.NewKeyFromSeed(seed[0:32]), nil
	}
	if c == SignatureSr25519 {
		msk, err := sr25519.NewMiniSecretKeyFromBytes(seed[0:32])
		if err != nil {
			return nil, err
		}
		sk := msk.ExpandEd25519()
		return sk, nil
	}

	return nil, fmt.Errorf("unable to generate private key from secret")
}

func GetPublicKeyFromSeed(seed []byte, c PolySignature, compressed bool) ([]byte, error) {
	if c == SignatureEd25519 {
		prvKey := ed25519.NewKeyFromSeed(seed[0:32])
		pubKey := prvKey.Public()
		return pubKey.(ed25519.PublicKey), nil

	}
	if c == SignatureSr25519 {
		msk, err := sr25519.NewMiniSecretKeyFromBytes(seed[0:32])
		if err != nil {
			return nil, err
		}
		sk := msk.ExpandEd25519()
		pubData, err := sk.PublicKey().MarshalBinary()
		if err != nil {
			return nil, err
		}
		return pubData, nil
	}
	// default to ecdsa
	curve := secp256k1.S256()
	x1, y1 := curve.ScalarBaseMult(seed[0:32])
	if compressed {
		compressedKey := append([]byte{0x02}, x1.Bytes()...)
		return compressedKey, nil
	}
	fullKey := append(x1.Bytes(), y1.Bytes()...)

	return fullKey, nil
}

func NewPolyWallet(mnemonic, password string) (*PolyWallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("the given mnemonic is invalid: %s", mnemonic)
	}
	pw := new(PolyWallet)
	pw.Mnemonic = mnemonic
	pw.Passphrase = password
	pw.derivationPath = "m/44'/60'/0'"
	pw.kdfIterations = 2048
	pw.keyCache = make(map[string]*bip32.Key, 0)
	pw.useRawEntropy = false
	err := pw.parseMnemonic()
	if err != nil {
		return nil, err
	}

	return pw, nil
}
func NewPolyWalletFromSeed(seed []byte) (*PolyWallet, error) {
	pw := new(PolyWallet)
	pw.derivationPath = "m/44'/60'/0'"
	pw.kdfIterations = 2048
	pw.keyCache = make(map[string]*bip32.Key, 0)
	pw.useRawEntropy = false
	return pw, nil
}
func (p *PolyWallet) SetPath(path string) error {
	// TODO validate the path more carefully
	if !rePathValidator.MatchString(path) {
		return fmt.Errorf("the path %s doesn't seem to make sense", path)
	}
	p.derivationPath = path
	return nil
}
func (p *PolyWallet) SetIterations(iterations uint) error {
	p.kdfIterations = iterations
	err := p.parseMnemonic()
	if err != nil {
		return err
	}

	return nil
}
func (p *PolyWallet) SetUseRawEntropy(e bool) error {
	p.useRawEntropy = e
	err := p.parseMnemonic()
	if err != nil {
		return err
	}

	return nil
}

func (p *PolyWallet) parseMnemonic() error {
	// substrate / polkadot
	// https://github.com/paritytech/substrate-bip39/blob/master/src/lib.rs
	if p.useRawEntropy {
		r, err := bip39.EntropyFromMnemonic(p.Mnemonic)
		if err != nil {
			return err
		}
		seed := pbkdf2.Key(r, []byte("mnemonic"+p.Passphrase), int(p.kdfIterations), 64, sha512.New)
		p.rawSeed = seed
		return nil
	}

	// 2048 is the default for bip39
	if p.kdfIterations == 2048 {
		seed := bip39.NewSeed(p.Mnemonic, p.Passphrase)
		p.rawSeed = seed
		return nil
	}

	// there might be a reason why someone would want a different number of iterations
	p.rawSeed = pbkdf2.Key([]byte(p.Mnemonic), []byte("mnemonic"+p.Passphrase), int(p.kdfIterations), 64, sha512.New)
	return nil
}
func (p *PolyWallet) ExportRootAddress() (*PolyWalletExport, error) {
	pwe := new(PolyWalletExport)
	pwe.Mnemonic = p.Mnemonic
	pwe.Passphrase = p.Passphrase // ???
	pwe.Seed = hex.EncodeToString(p.rawSeed)
	// assumes bip44
	rootKey, err := p.GetKeyForPath("m")
	if err != nil {
		return nil, err
	}

	pwe.RootKey = rootKey.String()
	pwe.HexPublicKey = hex.EncodeToString(rootKey.PublicKey().Key)
	pwe.HexPrivateKey = hex.EncodeToString(rootKey.Key)
	pwe.WIF = toWIF(rootKey)
	pwe.BTCAddress = toBTCAddress(rootKey)
	rootEthAddress := toETHAddress(rootKey)
	pwe.ETHAddress = rootEthAddress.String()
	pwe.HexFullPublicKey = hex.EncodeToString(toUncompressedPubKey(rootKey))
	addr, err := GetPublicKeyFromSeed(p.rawSeed, SignatureSecp256k1, true)
	if err != nil {
		return nil, err
	}
	pwe.ECDSAAddress = hex.EncodeToString(addr)
	pwe.ECDSAAddressSS58 = substrateSS58(addr)
	addr, err = GetPublicKeyFromSeed(p.rawSeed, SignatureEd25519, true)
	if err != nil {
		return nil, err
	}
	pwe.Ed25519Address = hex.EncodeToString(addr)
	pwe.Ed25519AddressSS58 = substrateSS58(addr)
	addr, err = GetPublicKeyFromSeed(p.rawSeed, SignatureSr25519, true)
	if err != nil {
		return nil, err
	}
	pwe.Sr25519Address = hex.EncodeToString(addr)
	pwe.Sr25519AddressSS58 = substrateSS58(addr)

	return pwe, nil
}

func (p *PolyWallet) ExportHDAddresses(count int) (*PolyWalletExport, error) {
	pwe := new(PolyWalletExport)
	pwe.Mnemonic = p.Mnemonic
	pwe.Passphrase = p.Passphrase // ???
	pwe.Seed = hex.EncodeToString(p.rawSeed)
	pwe.DerivationPath = p.derivationPath
	// assumes bip44
	rootKey, err := p.GetKeyForPath("m")
	if err != nil {
		return nil, err
	}

	pwe.RootKey = rootKey.String()

	accountKey, err := p.GetKeyForPath(p.derivationPath)
	if err != nil {
		return nil, err
	}
	pwe.AccountPrivateKey = accountKey.String()
	pwe.AccountPublicKey = accountKey.PublicKey().String()

	bip32Key, err := p.GetKeyForPath(p.derivationPath + "/0")
	if err != nil {
		return nil, err
	}
	pwe.BIP32PrivateKey = bip32Key.String()
	pwe.BIP32PublicKey = bip32Key.PublicKey().String()
	pwe.Addresses = make([]*PolyAddressExport, 0)

	for i := 0; i < count; i = i + 1 {
		// TODO if we want to provide support for hardened addresses it would need to be accommodated here
		currentPath := p.derivationPath + "/0/" + fmt.Sprintf("%d", i)
		k, err := p.GetKeyForPath(currentPath)
		if err != nil {
			return nil, err
		}

		pae := new(PolyAddressExport)
		pae.Path = currentPath
		pae.HexPublicKey = hex.EncodeToString(k.PublicKey().Key)
		pae.HexPrivateKey = hex.EncodeToString(k.Key)
		pae.WIF = toWIF(k)
		pae.BTCAddress = toBTCAddress(k)
		ethAddress := toETHAddress(k)
		pae.ETHAddress = ethAddress.String()
		pae.HexFullPublicKey = hex.EncodeToString(toUncompressedPubKey(k))
		pwe.Addresses = append(pwe.Addresses, pae)
	}
	return pwe, nil
}

// https://en.bitcoin.it/wiki/Wallet_import_format
func toWIF(prvKey *bip32.Key) string {
	mainnet := []byte{0x80}
	h0 := append(mainnet, prvKey.Key...)
	h0 = append(h0, 0x01)
	h1 := sha256.Sum256(h0)
	h2 := sha256.Sum256(h1[:])
	cksum := h2[0:4]
	h3 := append(h0, cksum...)
	return base58.Encode(h3)
}

func toETHAddress(prvKey *bip32.Key) common.Address {
	concat := toUncompressedPubKey(prvKey)
	return RawPubKeyToETHAddress(concat)

}
func RawPubKeyToETHAddress(concat []byte) common.Address {
	h := sha3.NewLegacyKeccak256()
	h.Write(concat)
	b := h.Sum(nil)
	return common.BytesToAddress(b)
}
func toUncompressedPubKey(prvKey *bip32.Key) []byte {
	// the GetPublicKey method returns a compressed key so we'll manually get the public key from the curve
	curve := secp256k1.S256()
	x1, y1 := curve.ScalarBaseMult(prvKey.Key)
	concat := append(x1.Bytes(), y1.Bytes()...)
	return concat
}

// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
func toBTCAddress(prvKey *bip32.Key) string {
	publicKey := prvKey.PublicKey()
	h := sha256.Sum256(publicKey.Key)
	ripe160 := ripemd160.New()
	ripe160.Write(h[:])
	h4 := ripe160.Sum([]byte{0})
	h5 := sha256.Sum256(h4)
	h6 := sha256.Sum256(h5[:])
	h7 := h6[0:4]
	h8 := append(h4, h7...)
	return base58.Encode(h8)
}

func (p *PolyWallet) GetKey() (*bip32.Key, error) {
	if key, hasKey := p.keyCache[p.derivationPath]; hasKey {
		return key, nil
	}
	path, err := p.parseDerivationPath()
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(p.rawSeed)
	if err != nil {
		return nil, err
	}
	currentKey := masterKey
	for _, levelIndex := range path {
		currentKey, err = currentKey.NewChildKey(levelIndex)
		if err != nil {
			return nil, err
		}
	}
	p.keyCache[p.derivationPath] = currentKey
	return currentKey, nil
}

func (p *PolyWallet) GetKeyForPath(inputPath string) (*bip32.Key, error) {
	if key, hasKey := p.keyCache[inputPath]; hasKey {
		return key, nil
	}

	path, err := parseDerivationPath(inputPath)
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(p.rawSeed)
	if err != nil {
		return nil, err
	}
	currentKey := masterKey
	for _, levelIndex := range path {
		currentKey, err = currentKey.NewChildKey(levelIndex)
		if err != nil {
			return nil, err
		}

	}
	p.keyCache[inputPath] = currentKey
	return currentKey, nil
}

// bip44... It looks like polkdadot substrate can support random paths
// with different conventions that are non numeric.
//
// TODO add support for polkdadot style derivation paths
func (p *PolyWallet) parseDerivationPath() ([]uint32, error) {
	return parseDerivationPath(p.derivationPath)
}

func parseDerivationPath(inputPath string) ([]uint32, error) {
	pieces := strings.Split(inputPath, "/")
	path := make([]uint32, 0)
	for idx, piece := range pieces {
		// m
		if idx == 0 {
			if piece != "m" {
				return nil, fmt.Errorf("expected derivation path to start with \"m\" but got \"%s\" instead", piece)
			}
			continue
		}

		// purpose = 1, coin_type = 2, account = 3, change = 4, address_index = 5
		if idx >= 1 && idx <= 5 {
			val, err := parsePathElement(piece)
			if err != nil {
				return nil, err
			}
			path = append(path, val)
		}

		if idx > 5 {
			return nil, fmt.Errorf("length of derivation path exceeded 5")
		}
	}
	return path, nil

}

func parsePathElement(element string) (uint32, error) {
	var base uint32 = 0
	if strings.Contains(element, "'") {
		base = bip32.FirstHardenedChild
		element = strings.ReplaceAll(element, "'", "")
	}
	pathVal, err := strconv.ParseUint(element, 10, 32)
	if err != nil {
		return base, err
	}
	return uint32(pathVal) + base, nil

}

func NewMnemonic(wordCount int, lang string) (string, error) {
	bits, hasKey := wordsToBits[wordCount]
	if !hasKey {
		return "", fmt.Errorf("the word count needs to be 12, 15, 18, 21, or 24. Got %d", wordCount)
	}
	wordList, hasKey := langToWordlist[strings.ToLower(lang)]
	if !hasKey {
		return "", fmt.Errorf("the language %s is not recognized", lang)
	}

	bip39.SetWordList(wordList)

	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return "", fmt.Errorf("there was an error getting entropy: %s", err.Error())
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("there was an error creating the mnemonic: %s", err.Error())
	}

	return mnemonic, nil
}

func init() {
	rePathValidator = regexp.MustCompile(pathValidator)
}

// https://github.com/paritytech/substrate/blob/8310936bd25519cee81abb01d2b164805f01bc25/primitives/core/src/crypto.rs#L264
// https://wiki.polkadot.network/docs/learn-accounts#for-the-curious-how-prefixes-work
func substrateSS58(pub []byte) string {
	data := append([]byte{42}, pub...)
	prefix := []byte("SS58PRE")
	h := blake2b.Sum512(append(prefix, data...))
	return base58.Encode(append(data, h[:2]...))
}
