package proof

import (
	"fmt"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"

	"github.com/spf13/cobra"
)

var EmptyProofCmd = &cobra.Command{
	Use:          "empty-proof",
	Short:        "Create an empty proof.",
	Long:         "Use this command to print an empty proof response that's filled with zero-valued siblings like 0x0000000000000000000000000000000000000000000000000000000000000000. This can be useful when you need to submit a dummy proof.",
	RunE:         runEmptyProof,
	SilenceUsage: true,
}

func runEmptyProof(_ *cobra.Command, _ []string) error {
	return emptyProof()
}

func emptyProof() error {
	p := new(ulxlycommon.Proof)

	e := ulxlycommon.GenerateEmptyHashes(ulxlycommon.TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(proofString(p))
	return nil
}
