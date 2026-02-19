// Package nullifier provides the compute-nullifier-tree command.
package nullifier

import (
	_ "embed"
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

const (
	ArgFileName = "file-name"
)

//go:embed usage.md
var usage string

var fileOptions = &tree.FileOptions{}

var Cmd = &cobra.Command{
	Use:   "compute-nullifier-tree",
	Short: "Compute the nullifier tree given the claims.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nullifierTree(args)
	},
	SilenceUsage: true,
}

func init() {
	Cmd.Flags().StringVar(&fileOptions.FileName, ArgFileName, "", "ndjson file with events data")
}

func nullifierTree(args []string) error {
	rawClaims, err := tree.GetInputData(args, fileOptions.FileName)
	if err != nil {
		return err
	}
	var root common.Hash
	root, err = tree.ComputeNullifierTree(rawClaims)
	if err != nil {
		return err
	}
	fmt.Printf(`
	{
		"root": "%s"
	}
	`, root.String())
	return nil
}
