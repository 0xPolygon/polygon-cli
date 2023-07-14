package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	_ "embed"

	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage         string
	inputFileName *string
	inputData     *string
)

var ABICmd = &cobra.Command{
	Use:   "abi Contract.abi",
	Short: "Parse an ABI and print the encoded signatures.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		// it would be nice to have a generic reader

		rawData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		buf := bytes.NewReader(rawData)
		abi, err := gethabi.JSON(buf)
		if err != nil {
			return err
		}
		for _, meth := range abi.Methods {
			fmt.Printf("Selector:%s\tSignature:%s\n", hex.EncodeToString(meth.ID), meth)
		}
		if *inputData != "" {
			id, callData, err := parseContractInputData(*inputData)
			if err != nil {
				return err
			}
			meth, err := abi.MethodById(id)
			if err != nil {
				return err
			}
			if meth == nil {
				return fmt.Errorf("the function selector %s wasn't matched in the given abi", hex.EncodeToString(id))
			}
			inputVals := make(map[string]any, 0)
			err = meth.Inputs.UnpackIntoMap(inputVals, callData)
			if err != nil {
				return err
			}
			fmt.Println("Input data:")
			prettyInput, _ := json.MarshalIndent(inputVals, "", "  ")
			fmt.Println(string(prettyInput))
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := ABICmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a filename to read and analyze")
	inputData = flagSet.String("data", "", "Provide input data to be unpacked based on the ABI definition")
}

func parseContractInputData(data string) ([]byte, []byte, error) {
	// "0x11223344"
	if len(data) < 10 {
		return nil, nil, fmt.Errorf("the input %s is too short for a function call. It should start with 0x and needs at least 4 bytes for a function selector", data)
	}
	if data[0:2] != "0x" {
		return nil, nil, fmt.Errorf("the input data must start with 0x")
	}
	// drop the 0x and select the next bytes to represent the selector
	stringId := data[2:10]
	rawId, err := hex.DecodeString(stringId)
	if err != nil {
		return nil, nil, err
	}
	rawCallData, err := hex.DecodeString(data[10:])
	if err != nil {
		return nil, nil, err
	}
	return rawId, rawCallData, err
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
