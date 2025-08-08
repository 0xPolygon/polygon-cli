package util

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

func ReadPrivateKeysFromFile(sendingAddressesFile string) ([]*ecdsa.PrivateKey, error) {
	file, err := os.Open(sendingAddressesFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open sending addresses file: %w", err)
	}
	defer file.Close()

	var privateKeys []*ecdsa.PrivateKey
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		privateKey, err := ethcrypto.HexToECDSA(strings.TrimPrefix(line, "0x"))
		if err != nil {
			log.Error().Err(err).Str("key", line).Msg("Unable to parse private key")
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
		privateKeys = append(privateKeys, privateKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading sending address file: %w", err)
	}

	return privateKeys, nil
}

// Returns the address and private key of the given private key
func GetAddressAndPrivateKeyHex(ctx context.Context, privateKey *ecdsa.PrivateKey) (string, string) {
	privateKeyHex := GetPrivateKeyHex(privateKey)
	address := GetAddress(ctx, privateKey)
	return address.String(), privateKeyHex
}

func GetPrivateKeyHex(privateKey *ecdsa.PrivateKey) string {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	return fmt.Sprintf("0x%x", privateKeyBytes)
}

func GetAddress(ctx context.Context, privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKey)
}
