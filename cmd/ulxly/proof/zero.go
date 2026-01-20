package proof

import (
	"fmt"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"

	"github.com/spf13/cobra"
)

var ZeroProofCmd = &cobra.Command{
	Use:   "zero-proof",
	Short: "Create a proof that's filled with zeros.",
	Long: `Use this command to print a proof response that's filled with the zero
hashes. These values are very helpful for debugging because they would
tell you how populated the tree is and roughly which leaves and
siblings are empty. It's also helpful for sanity checking a proof
response to understand if the hashed value is part of the zero hashes
or if it's actually an intermediate hash.`,
	RunE:         runZeroProof,
	SilenceUsage: true,
}

func runZeroProof(_ *cobra.Command, _ []string) error {
	return zeroProof()
}

func zeroProof() error {
	p := new(ulxlycommon.Proof)

	e := ulxlycommon.GenerateZeroHashes(ulxlycommon.TreeDepth)
	copy(p.Siblings[:], e)
	fmt.Println(proofString(p))
	return nil
}
