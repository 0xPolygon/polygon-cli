package decode

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
	inputFileName *string
	inputData     *string
)

var ABIDecodeCmd = &cobra.Command{
	Use:   "decode Contract.abi",
	Short: "Parse an ABI and print the encoded signatures.",
	Long:  "",
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
			fmt.Printf("Selector:%s\tSignature:%s%s\n", hex.EncodeToString(meth.ID), meth.Sig, getReturnSignature(meth.Outputs))
		}
		if *inputData != "" {
			id, callData, err := parseContractInputData(*inputData)
			fmt.Printf("id: %x, %x\n", id, callData)
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

			unpackedCallData, err := meth.Inputs.UnpackValues(callData)
			if err != nil {
				return err
			}

			fmt.Println("Signature and Input")
			fmt.Printf("%s%s", meth.Sig, getReturnSignature(meth.Outputs))
			for _, unpackedCallDataArg := range unpackedCallData {
				fmt.Printf(" %v", unpackedCallDataArg)
			}
			fmt.Println()
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := ABIDecodeCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a filename to read and analyze")
	inputData = flagSet.String("data", "", "Provide input data to be unpacked based on the ABI definition")
}

func parseContractInputData(data string) ([]byte, []byte, error) {
	// "0x11223344"
	selectorLength := 8
	data = strings.TrimPrefix(data, "0x")
	if len(data) < selectorLength {
		return nil, nil, fmt.Errorf("the input %s is too short for a function call. It should start with 0x and needs at least 4 bytes for a function selector", data)
	}

	selectorId := data[:selectorLength]
	rawId, err := hex.DecodeString(selectorId)
	if err != nil {
		return nil, nil, err
	}
	rawCallData, err := hex.DecodeString(data[selectorLength:])
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

func getReturnSignature(funcReturns gethabi.Arguments) string {
	returnSig := "("
	for key, ret := range funcReturns {
		// Append comma only for non first and last elements
		if key > 0 && key < len(funcReturns) {
			returnSig += ","
		}
		returnSig += ret.Type.String()
	}
	returnSig += ")"

	return returnSig
}
