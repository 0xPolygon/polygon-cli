package decode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// newMsgCmd builds `decode msg <TYPE_URL> <B64>`. Renders the decoded
// message as JSON. With --list the command prints every type URL
// registered locally and exits; --list is handy for discovery without
// needing a live Heimdall node.
func newMsgCmd() *cobra.Command {
	var listOnly bool
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "msg <type-url> <value-b64>",
		Short: "Decode a single Any.value for type-url (base64 value).",
		Long: strings.TrimSpace(`
Decode a single Any.value for a registered type URL.

Example:
  polycli heimdall decode msg /heimdallv2.topup.MsgWithdrawFeeTx \
    CioweDAxNzE3MDAyN2YwYzVjZDE5MDRmOGI0MDU1OGRhZjUwN2FiNGViNjJhEgEw
`),
		Args: func(cmd *cobra.Command, args []string) error {
			if listOnly {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listOnly {
				out := cmd.OutOrStdout()
				for _, u := range hproto.KnownTypeURLs() {
					if _, err := fmt.Fprintln(out, u); err != nil {
						return err
					}
				}
				return nil
			}
			typeURL := strings.TrimSpace(args[0])
			if typeURL == "" {
				return &client.UsageError{Msg: "type-url is required"}
			}
			val, err := decodeInput("value", args[1])
			if err != nil {
				return err
			}
			decoded, err := hproto.Decode(typeURL, val)
			if err != nil {
				return err
			}
			env := map[string]interface{}{
				"type_url": typeURL,
				"value":    decoded,
			}
			if jsonOut {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(env)
			}
			buf, err := json.MarshalIndent(env, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(buf))
			return err
		},
	}
	cmd.Flags().BoolVar(&listOnly, "list", false, "print every registered type URL and exit")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit single-line JSON")
	return cmd
}
