package gethkeystore

import (
	"crypto/ecdsa"

	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

type RawKeystoreData struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
}

func DecryptKeystoreFile(fileData []byte, password string) (*ecdsa.PrivateKey, error) {
	var encryptedData RawKeystoreData
	err := json.Unmarshal(fileData, &encryptedData)
	if err != nil {
		return nil, err
	}
	decryptedData, err := keystore.DecryptDataV3(encryptedData.Crypto, password)
	if err != nil {
		return nil, err
	}

	privKeyStr := hex.EncodeToString(decryptedData)

	privKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
