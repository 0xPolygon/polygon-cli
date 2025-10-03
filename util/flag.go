package util

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// FlagRPCURL is the standard flag name for RPC endpoint URLs.
	FlagRPCURL = "rpc-url"
	// FlagPrivateKey is the standard flag name for private keys.
	FlagPrivateKey = "private-key"
)

// GetRequiredFlag retrieves a flag value from Viper after binding it and marking it as required.
// It binds the flag to enable environment variable fallback via Viper, marks the flag as required,
// and returns the flag's value.
func GetRequiredFlag(cmd *cobra.Command, flagName string) (string, error) {
	viper.BindPFlag(flagName, cmd.Flags().Lookup(flagName))
	if err := cmd.MarkFlagRequired(flagName); err != nil {
		return "", err
	}
	return viper.GetString(flagName), nil
}

// GetRPCURL retrieves the rpc-url flag value from Viper after binding it, marking it as required,
// and validating that it is a valid URL with a supported scheme (http, https, ws, wss).
func GetRPCURL(cmd *cobra.Command) (string, error) {
	rpcURL, err := GetRequiredFlag(cmd, FlagRPCURL)
	if err != nil {
		return "", err
	}
	if err := ValidateUrl(rpcURL); err != nil {
		return "", err
	}
	return rpcURL, nil
}

// GetPrivateKey retrieves the private-key flag value from Viper after binding it and marking it as required.
// This is a convenience wrapper around GetRequiredFlag for the standard private key flag.
func GetPrivateKey(cmd *cobra.Command) (string, error) {
	return GetRequiredFlag(cmd, FlagPrivateKey)
}

// BigIntValue is a custom flag type for big.Int values.
// It implements the pflag.Value interface to enable using *big.Int with Cobra flags.
type BigIntValue struct {
	Val *big.Int
}

// String returns the decimal string representation of the big.Int value.
func (b *BigIntValue) String() string {
	return b.Val.String()
}

// Set parses a decimal string and sets the big.Int value.
func (b *BigIntValue) Set(s string) error {
	if _, ok := b.Val.SetString(s, 10); !ok {
		return fmt.Errorf("invalid big integer: %q", s)
	}
	return nil
}

// Type returns the type string for this flag value.
func (b *BigIntValue) Type() string {
	return "big.Int"
}
