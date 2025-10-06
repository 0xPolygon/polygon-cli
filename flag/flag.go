// Package flag provides utilities for managing command flags with environment variable fallback support.
// It implements a priority system: flag value > environment variable > default value.
package flag

import (
	"fmt"
	"os"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	// RPCURL is the flag name for RPC URL
	RPCURL = "rpc-url"
	// RPCURLEnvVar is the environment variable name for RPC URL
	RPCURLEnvVar = "ETH_RPC_URL"
	// DefaultRPCURL is the default RPC URL when no flag or env var is set
	DefaultRPCURL = "http://localhost:8545"
	// PrivateKey is the flag name for private key
	PrivateKey = "private-key"
	// PrivateKeyEnvVar is the environment variable name for private key
	PrivateKeyEnvVar = "PRIVATE_KEY"
)

// GetRPCURL retrieves the RPC URL from the command flag or environment variable.
// Returns the flag value if set, otherwise the environment variable value, otherwise empty string.
// Validates the URL format if a non-empty value is provided and returns an error if validation fails.
// Returns empty string and nil error if no value is set.
func GetRPCURL(cmd *cobra.Command) (string, error) {
	value, err := getValue(cmd, RPCURL, RPCURLEnvVar, false)
	if err != nil || value == "" {
		return value, err
	}

	if err := util.ValidateUrl(value); err != nil {
		return "", err
	}

	return value, nil
}

// GetRequiredRPCURL retrieves the RPC URL from the command flag or environment variable.
// Returns the flag value if set, otherwise the environment variable value.
// Validates the URL format and returns an error if the value is not set, empty, or invalid.
func GetRequiredRPCURL(cmd *cobra.Command) (string, error) {
	value, err := getValue(cmd, RPCURL, RPCURLEnvVar, true)
	if err != nil {
		return "", err
	}

	if err := util.ValidateUrl(value); err != nil {
		return "", err
	}

	return value, nil
}

// GetPrivateKey retrieves the private key from the command flag or environment variable.
// Returns the flag value if set, otherwise the environment variable value, otherwise the default.
// Returns empty string and nil error if none are set.
func GetPrivateKey(cmd *cobra.Command) (string, error) {
	return getValue(cmd, PrivateKey, PrivateKeyEnvVar, false)
}

// GetRequiredPrivateKey retrieves the private key from the command flag or environment variable.
// Returns an error if the value is not set or empty.
func GetRequiredPrivateKey(cmd *cobra.Command) (string, error) {
	return getValue(cmd, PrivateKey, PrivateKeyEnvVar, true)
}

// getValue retrieves a flag value with environment variable fallback support.
// It implements a priority system where flag values take precedence over environment variables,
// which take precedence over default values.
//
// Parameters:
//   - cmd: The cobra command to retrieve the flag from
//   - flagName: The name of the flag to retrieve
//   - envVarName: The environment variable name to check as fallback
//   - required: Whether the value is required (returns error if empty)
//
// Returns the resolved value and an error if required validation fails.
func getValue(cmd *cobra.Command, flagName, envVarName string, required bool) (string, error) {
	flag := cmd.Flag(flagName)
	if flag == nil {
		return "", fmt.Errorf("flag %q not found", flagName)
	}

	// Priority: flag > env var > default
	value := flag.DefValue

	envVarValue := os.Getenv(envVarName)
	if envVarValue != "" {
		value = envVarValue
	}

	if flag.Changed {
		value = flag.Value.String()
	}

	if required && value == "" {
		return "", fmt.Errorf("required flag(s) %q not set", flagName)
	}

	return value, nil
}

// MarkFlagRequired marks one or more regular flags as required and logs a fatal error if marking fails.
// This helper ensures consistent error handling across all commands when marking flags as required.
func MarkFlagRequired(cmd *cobra.Command, flagNames ...string) {
	for _, flagName := range flagNames {
		if err := cmd.MarkFlagRequired(flagName); err != nil {
			log.Fatal().
				Err(err).
				Str("flag", flagName).
				Str("command", cmd.Name()).
				Msg("Failed to mark flag as required")
		}
	}
}

// MarkPersistentFlagRequired marks one or more persistent flags as required and logs a fatal error if marking fails.
// This helper ensures consistent error handling across all commands when marking persistent flags as required.
func MarkPersistentFlagRequired(cmd *cobra.Command, flagNames ...string) {
	for _, flagName := range flagNames {
		if err := cmd.MarkPersistentFlagRequired(flagName); err != nil {
			log.Fatal().
				Err(err).
				Str("flag", flagName).
				Str("command", cmd.Name()).
				Msg("Failed to mark persistent flag as required")
		}
	}
}
