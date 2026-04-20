package tx

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/heimdall/tx/msgs"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// registerMktxSendEstimate attaches the three umbrella commands
// (`mktx`, `send`, `estimate`) to parent. Each umbrella carries a
// fresh copy of every registered msg subcommand — cobra commands can
// only have one parent, so we rebuild the subtree per umbrella
// rather than sharing pointers.
//
// Called from Register alongside the other cast-style top-level tx
// subcommands (tx/receipt/logs/nonce/...). The msgs package's
// registry holds the msg-name -> factory map; W4 will add more
// factories via its own init() calls, which this function picks up
// automatically the next time it runs.
func registerMktxSendEstimate(parent *cobra.Command, flags *config.Flags) {
	parent.AddCommand(
		newUmbrellaCmd(msgs.ModeMkTx, flags),
		newUmbrellaCmd(msgs.ModeSend, flags),
		newUmbrellaCmd(msgs.ModeEstimate, flags),
	)
}

// newUmbrellaCmd builds a single umbrella (mktx|send|estimate). The
// short/long strings are generated per-mode; the subcommand list is
// populated from msgs.BuildChildren.
func newUmbrellaCmd(mode msgs.Mode, globalFlags *config.Flags) *cobra.Command {
	var short, long string
	switch mode {
	case msgs.ModeMkTx:
		short = "Build a signed TxRaw without broadcasting."
		long = strings.TrimSpace(`
Construct a Heimdall v2 transaction for the chosen message type and
print the signed TxRaw bytes as 0x-prefixed hex. Nothing is sent.
Use --json for an envelope that also carries the base64 form
accepted by the REST gateway.

Supply exactly one of --from, --account, --private-key, or
--mnemonic so the builder can sign. Pair --dry-run with send if you
want a round-trip that stops just before broadcast instead.
`)
	case msgs.ModeSend:
		short = "Build, sign, and broadcast a transaction."
		long = strings.TrimSpace(`
Build a Heimdall v2 transaction for the chosen message type, sign
it, and POST it to the REST gateway. The default mode is
BROADCAST_MODE_SYNC: polycli waits for CheckTx to return, prints
the tx hash, and then polls CometBFT for inclusion. --async skips
both waits. --confirmations N waits for N blocks past inclusion.
--dry-run stops after building (useful for CI).
`)
	case msgs.ModeEstimate:
		short = "Simulate a transaction and report gas usage."
		long = strings.TrimSpace(`
Build a transaction for the chosen message type and call
/cosmos/tx/v1beta1/simulate to estimate gas without broadcasting.
Pair with --gas-price to see the implied fee for the simulated gas
amount.
`)
	}
	cmd := &cobra.Command{
		Use:     mode.String() + " <MSG>",
		Short:   short,
		Long:    long,
		Args:    cobra.NoArgs, // children enforce their own args
		Aliases: nil,
	}
	// Helpful hint when no msg subcommand is provided.
	cmd.SilenceUsage = true
	cmd.RunE = func(c *cobra.Command, _ []string) error {
		names := msgs.Names()
		return fmt.Errorf("%s requires a message subcommand (one of: %s)", mode.String(), strings.Join(names, ", "))
	}
	for _, child := range msgs.BuildChildren(mode, globalFlags) {
		cmd.AddCommand(child)
	}
	return cmd
}
