package abi

import (
	"github.com/spf13/cobra"

	_ "embed"

	"github.com/maticnetwork/polygon-cli/cmd/abi/decode"
	"github.com/maticnetwork/polygon-cli/cmd/abi/encode"
)

var (
	//go:embed usage.md
	usage string
)

var ABICmd = &cobra.Command{
	Use:   "abi",
	Short: "Provides encoding and decoding functionalities with contract signatures and ABI.",
	Long:  usage,
}

func init() {
	ABICmd.AddCommand(decode.ABIDecodeCmd)
	ABICmd.AddCommand(encode.ABIEncodeCmd)
}
