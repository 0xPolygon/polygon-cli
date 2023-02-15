package abi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	inputFileName *string
)

var ABICmd = &cobra.Command{
	Use:   "abi Contract.abi",
	Short: "A simple tool to parse an ABI and print the encoded signatures",
	Long: `
When looking at raw contract calls, sometimes we have an ABI and we just want
to quickly figure out which method is being called. This is a quick way to
get all of the function selectors for an ABI

go run main.go abi --file ../zkevm-node/etherman/smartcontracts/abi/polygonzkevm.abi

go run main.go abi < ../zkevm-node/etherman/smartcontracts/abi/polygonzkevm.abi
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// it would be nice to have a generic reader

		rawData, err := getInputData(cmd, args)
		if err != nil {
			return nil
		}
		buf := bytes.NewReader(rawData)
		abi, err := gethabi.JSON(buf)
		if err != nil {
			return err
		}
		for _, meth := range abi.Methods {
			fmt.Printf("Selector:%s\tSignature:%s\n", hex.EncodeToString(meth.ID), meth)
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
