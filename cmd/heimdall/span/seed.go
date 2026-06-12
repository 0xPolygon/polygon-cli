package span

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newSeedCmd builds `span seed <ID>` → GET /bor/spans/seed/{id}.
// Prints seed (already 0x-hex upstream) and seed_author.
func newSeedCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "seed <ID>",
		Short: "Show seed and seed_author for a span.",
		Args:  cobra.ExactArgs(1),
		Label: "span seed",
		Build: func(cmd *cobra.Command, args []string) (string, url.Values, error) {
			id, err := parseSpanID("span id", args[0])
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("/bor/spans/seed/%d", id), nil, nil
		},
	})
}
