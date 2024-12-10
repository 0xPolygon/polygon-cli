package wrapcontract

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"encoding/json"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string
	json_storage *string
)

var WrapContractCmd = &cobra.Command{
	Use:   "wrap-contract bytecode|file",
	Aliases: []string{"wrapcontract", "wrapContract"},
	Short: "Wrap deployed bytecode into create bytecode.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		deployed_bytecode, err := getInputData(args)
		if err != nil {
			cmd.PrintErrf("There was an error reading input for wrapping contract: %s", err.Error())
			return err
		}
		storage_bytecode, err := getStorageBytecode()
		if err != nil {
			cmd.PrintErrf("There was an error reading storage map: %s", err.Error())
		}
		create_bytecode := util.WrapDeployedCode(deployed_bytecode, storage_bytecode)
		fmt.Println(create_bytecode)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("expected at most one argument: bytecode")
		}
		return nil
	},
}

func init() {
	flagSet := WrapContractCmd.PersistentFlags()
	json_storage = flagSet.String("storage", "", "Provide storage slots in json format k:v")
}

func getInputData(args []string) (string, error) {
	var deployed_bytecode string = ""
	if len(args) == 0 {
		deployed_bytecode_bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		deployed_bytecode = string(deployed_bytecode_bytes)
	} else {
		deployed_bytecode_or_file := args[0]
		// Try to open the param as a file, otherwise treat it as bytecode
		deployed_bytecode_bytes, err := os.ReadFile(deployed_bytecode_or_file)
		if err != nil {
			deployed_bytecode = deployed_bytecode_or_file
		} else {
			deployed_bytecode = string(deployed_bytecode_bytes)
		}
	}

	return strings.TrimSpace(deployed_bytecode), nil
}

func getStorageBytecode() (string, error) {
	var storage_bytecode string = ""
	
	if json_storage != nil && *json_storage != "" {
		var storage map[string]string
		err := json.Unmarshal([]byte(*json_storage), &storage)
		if err != nil {
			return storage_bytecode, err
		}
		for k, v := range storage {
			slot := util.GetHexString(k)
			value := util.GetHexString(v)
			sLen := len(slot) / 2
			vLen := len(value) / 2
			sPushCode := 0x5f + sLen
			vPushCode := 0x5f + vLen
			storage_bytecode += fmt.Sprintf("%02x%s%02x%s55", vPushCode, value, sPushCode, slot)
		}
	}

	return storage_bytecode, nil
}
