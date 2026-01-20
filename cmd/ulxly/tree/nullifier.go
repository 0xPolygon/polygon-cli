package tree

import (
	_ "embed"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

const (
	ArgFileName = "file-name"
)

//go:embed computeNullifierTreeUsage.md
var computeNullifierTreeUsage string

var fileOptions = &FileOptions{}

var NullifierTreeCmd = &cobra.Command{
	Use:   "compute-nullifier-tree",
	Short: "Compute the nullifier tree given the claims.",
	Long:  computeNullifierTreeUsage,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nullifierTree(args)
	},
	SilenceUsage: true,
}

func init() {
	NullifierTreeCmd.Flags().StringVar(&fileOptions.FileName, ArgFileName, "", "ndjson file with events data")
}

func nullifierTree(args []string) error {
	rawClaims, err := getInputData(args, fileOptions.FileName)
	if err != nil {
		return err
	}
	var root common.Hash
	root, err = computeNullifierTree(rawClaims)
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
