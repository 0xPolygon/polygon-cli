package wrapcontract

import (
	"fmt"
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/0xPolygon/polygon-cli/util"
)

var (
	//go:embed usage.md
	usage                  string
)

var WrapContractCmd = &cobra.Command{
	Use:   "wrap-contract bytecode",
	Aliases: []string{"wrapcontract", "wrapContract"},
	Short: "Wrap deployed bytecode into create bytecode.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		deployed_bytecode := args[0]
		storage_bytecode := ""
		create_bytecode := util.WrapDeployedCode(deployed_bytecode, storage_bytecode)
		fmt.Println(create_bytecode)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected exactly one argument: bytecode")
		}
		return nil
	},
}
