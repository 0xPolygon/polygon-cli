package flag_loader

import (
	"fmt"
	"os"
	"math/big"

	"github.com/spf13/cobra"
)

const (
	rpcUrlFlagName, rpcUrlEnvVar         = "rpc-url", "ETH_RPC_URL"
	privateKeyFlagName, privateKeyEnvVar = "private-key", "PRIVATE_KEY"
)

type BigIntValue struct {
    Val *big.Int
}

func (b *BigIntValue) String() string {
    // Return the decimal representation
    return b.Val.String()
}

func (b *BigIntValue) Set(s string) error {
    // Parse the string in base 10
    if _, ok := b.Val.SetString(s, 10); !ok {
        return fmt.Errorf("invalid big integer: %q", s)
    }
    return nil
}

func (b *BigIntValue) Type() string {
    return "big.Int"
}

func GetRpcUrlFlagValue(cmd *cobra.Command) *string {
	v, _ := getFlagValue(cmd, rpcUrlFlagName, rpcUrlEnvVar, false)
	return v
}

func GetRequiredRpcUrlFlagValue(cmd *cobra.Command) (*string, error) {
	return getFlagValue(cmd, rpcUrlFlagName, rpcUrlEnvVar, true)
}

func GetPrivateKeyFlagValue(cmd *cobra.Command) *string {
	v, _ := getFlagValue(cmd, privateKeyFlagName, privateKeyEnvVar, false)
	return v
}

func GetRequiredPrivateKeyFlagValue(cmd *cobra.Command) (*string, error) {
	return getFlagValue(cmd, privateKeyFlagName, privateKeyEnvVar, true)
}

func getFlagValue(cmd *cobra.Command, flagName, envVarName string, required bool) (*string, error) {
	flag := cmd.Flag(flagName)
	var flagValue string
	if flag.Changed {
		flagValue = flag.Value.String()
	}
	flagDefaultValue := flag.DefValue

	envVarValue := os.Getenv(envVarName)

	value := flagDefaultValue
	if envVarValue != "" {
		value = envVarValue
	}
	if flag.Changed {
		value = flagValue
	}

	if required && (!flag.Changed && envVarValue == "") {
		return nil, fmt.Errorf("required flag(s) \"%s\" not set", flagName)
	}

	return &value, nil
}
