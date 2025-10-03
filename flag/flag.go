package flag

import (
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// RPCURL is the standard flag name for RPC endpoint URLs.
	RPCURL = "rpc-url"
	// PrivateKey is the standard flag name for private keys.
	PrivateKey = "private-key"

	// DefaultRPCURL is the default RPC endpoint URL.
	DefaultRPCURL = "http://localhost:8545"
)

// GetFlag retrieves a flag value from Viper after binding it.
// It binds the flag to enable environment variable fallback via Viper.
func GetFlag(cmd *cobra.Command, flagName string) string {
	if err := viper.BindPFlag(flagName, cmd.Flags().Lookup(flagName)); err != nil {
		log.Fatal().Err(err).Str("flag", flagName).Msg("Failed to bind flag to viper")
	}
	return viper.GetString(flagName)
}

// GetRPCURL retrieves the rpc-url flag value from Viper after binding it and validates
// that it is a valid URL with a supported scheme (http, https, ws, wss).
func GetRPCURL(cmd *cobra.Command) (string, error) {
	rpcURL := GetFlag(cmd, RPCURL)
	if err := util.ValidateUrl(rpcURL); err != nil {
		return "", err
	}
	return rpcURL, nil
}

// GetPrivateKey retrieves the private-key flag value from Viper after binding it.
// This is a convenience wrapper around GetFlag for the standard private key flag.
func GetPrivateKey(cmd *cobra.Command) (string, error) {
	return GetFlag(cmd, PrivateKey), nil
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
