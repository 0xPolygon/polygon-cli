package util

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/rs/zerolog/log"
)

// HexToBigInt converts a hexadecimal string to a big integer.
func HexToBigInt(hexString string) (*big.Int, error) {
	// Clean up the string.
	// - Remove the `0x` prefix.
	// - If the length of the string is odd, pad it with a leading zero.
	hexString = strings.TrimPrefix(hexString, "0x")
	if len(hexString)%2 != 0 {
		log.Trace().Str("original", hexString).Msg("Hexadecimal string of odd length padded with a leading zero")
		hexString = "0" + hexString
	}

	// Decode the hexadecimal string into a byte slice and return the `big.Int` value.
	rawGas, err := hex.DecodeString(hexString)
	if err != nil {
		log.Error().Err(err).Str("hex", hexString).Msg("Unable to decode hex string")
		return nil, err
	}
	bigInt := big.NewInt(0)
	bigInt.SetBytes(rawGas)
	return bigInt, nil
}

// EthToWei converts a given amount of Ether to Wei.
func EthToWei(ethAmount float64) *big.Int {
	weiBigFloat := new(big.Float).SetFloat64(ethAmount)
	weiBigFloat.Mul(weiBigFloat, new(big.Float).SetInt64(1e18))
	weiBigInt, _ := weiBigFloat.Int(nil)
	return weiBigInt
}
