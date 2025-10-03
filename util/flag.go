package util

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// FlagRPCURL is the standard flag name for RPC endpoint URLs.
	FlagRPCURL = "rpc-url"
	// FlagPrivateKey is the standard flag name for private keys.
	FlagPrivateKey = "private-key"

	// DefaultRPCURL is the default RPC endpoint URL.
	DefaultRPCURL = "http://localhost:8545"
)

// GetFlag retrieves a flag value from Viper after binding it.
// It binds the flag to enable environment variable fallback via Viper.
func GetFlag(cmd *cobra.Command, flagName string) string {
	viper.BindPFlag(flagName, cmd.Flags().Lookup(flagName))
	return viper.GetString(flagName)
}

// GetRPCURL retrieves the rpc-url flag value from Viper after binding it and validates
// that it is a valid URL with a supported scheme (http, https, ws, wss).
func GetRPCURL(cmd *cobra.Command) (string, error) {
	rpcURL := GetFlag(cmd, FlagRPCURL)
	if err := ValidateUrl(rpcURL); err != nil {
		return "", err
	}
	return rpcURL, nil
}

// GetPrivateKey retrieves the private-key flag value from Viper after binding it.
// This is a convenience wrapper around GetFlag for the standard private key flag.
func GetPrivateKey(cmd *cobra.Command) (string, error) {
	return GetFlag(cmd, FlagPrivateKey), nil
}

// MarkFlagRequired marks a regular flag as required and logs a fatal error if marking fails.
// This helper ensures consistent error handling across all commands when marking flags as required.
func MarkFlagRequired(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkFlagRequired(flagName); err != nil {
		log.Fatal().
			Err(err).
			Str("flag", flagName).
			Str("command", cmd.Name()).
			Msg("Failed to mark flag as required")
	}
}

// MarkPersistentFlagRequired marks a persistent flag as required and logs a fatal error if marking fails.
// This helper ensures consistent error handling across all commands when marking persistent flags as required.
func MarkPersistentFlagRequired(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkPersistentFlagRequired(flagName); err != nil {
		log.Fatal().
			Err(err).
			Str("flag", flagName).
			Str("command", cmd.Name()).
			Msg("Failed to mark persistent flag as required")
	}
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
