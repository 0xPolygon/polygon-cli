package encode

import (
	"fmt"

	_ "embed"

	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/abi"
)

var (
	functionSignature string
	functionInputs    []string
)

var ABIEncodeCmd = &cobra.Command{
	Use:   "encode [function signature] [args...]",
	Short: "ABI encodes a function signature and the inputs",
	Long:  "[function-signature] is required and is a fragment in the form <function name>(<types...>). If the function signature has parameters, then those values would have to be passed as arguments after the function signature.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		functionSignature = args[0]
		if len(args) > 1 {
			functionInputs = args[1:]
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fullEncoding, err := abi.AbiEncode(functionSignature, functionInputs)
		if err != nil {
			return err
		}
		fmt.Println(fullEncoding)

		return nil
	},
}

func init() {}
