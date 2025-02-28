package flag_loader

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func GetRpcUrlFlagValue(cmd *cobra.Command, required bool) (*string, error) {
	return getFlagValue(cmd, "rpc-url", "ETH_RPC_URL", required)
}

func GetPrivateKeyFlagValue(cmd *cobra.Command, required bool) (*string, error) {
	return getFlagValue(cmd, "private-key", "PRIVATE_KEY", required)
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
