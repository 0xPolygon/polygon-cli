package wrapcontract

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage       string
	jsonStorage *string
)

var WrapContractCmd = &cobra.Command{
	Use:     "wrap-contract bytecode|file",
	Aliases: []string{"wrapcontract", "wrapContract"},
	Short:   "Wrap deployed bytecode into create bytecode.",
	Long:    usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		deployedBytecode, err := getInputData(args)
		if err != nil {
			cmd.PrintErrf("There was an error reading input for wrapping contract: %s", err.Error())
			return err
		}
		storageBytecode, err := getStorageBytecode()
		if err != nil {
			cmd.PrintErrf("There was an error reading storage map: %s", err.Error())
		}
		createBytecode := util.WrapDeployedCode(deployedBytecode, storageBytecode)
		fmt.Println(createBytecode)
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
	flagSet := WrapContractCmd.Flags()
	jsonStorage = flagSet.String("storage", "", "Provide storage slots in json format k:v")
}

func getInputData(args []string) (string, error) {
	var deployedBytecode string
	var deployedBytecodeOrFile string

	if len(args) == 0 {
		deployedBytecodeOrFileBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		deployedBytecodeOrFile = string(deployedBytecodeOrFileBytes)
	} else {
		deployedBytecodeOrFile = args[0]
	}
	// Try to open the param as a file, otherwise treat it as bytecode
	deployedBytecodeOrFile = strings.TrimSpace(deployedBytecodeOrFile)
	deployedBytecodeBytes, err := os.ReadFile(deployedBytecodeOrFile)
	if err != nil {
		deployedBytecode = deployedBytecodeOrFile
	} else {
		deployedBytecode = string(deployedBytecodeBytes)
	}

	return strings.TrimSpace(deployedBytecode), nil
}

func getStorageBytecode() (string, error) {
	var storageBytecode string = ""

	if jsonStorage != nil && *jsonStorage != "" {
		var storage map[string]string
		err := json.Unmarshal([]byte(*jsonStorage), &storage)
		if err != nil {
			return storageBytecode, err
		}
		for k, v := range storage {
			slot := util.GetHexString(k)
			value := util.GetHexString(v)
			sLen := len(slot) / 2
			vLen := len(value) / 2
			sPushCode := 0x5f + sLen
			vPushCode := 0x5f + vLen
			storageBytecode += fmt.Sprintf("%02x%s%02x%s55", vPushCode, value, sPushCode, slot)
		}
	}

	return storageBytecode, nil
}
