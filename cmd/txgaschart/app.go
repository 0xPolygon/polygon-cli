package txgaschart

import (
	_ "embed"
	"math"

	"github.com/spf13/cobra"
)

type args struct {
	rpcURL      string
	rateLimit   float64
	concurrency uint64

	scale string

	startBlock uint64
	endBlock   uint64

	targetAddr string

	output string
}

var inputArgs = args{}

//go:embed usage.md
var usage string
var Cmd = &cobra.Command{
	Use:   "tx-gas-chart",
	Short: "plot a chart of transaction gas prices and limits",
	Long:  usage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildChart(cmd)
	},
}

func init() {
	f := Cmd.PersistentFlags()
	f.StringVar(&inputArgs.rpcURL, "rpc-url", "http://localhost:8545", "RPC URL of network")
	f.Float64Var(&inputArgs.rateLimit, "rate-limit", 4, "requests per second limit (use negative value to remove limit)")
	f.Uint64VarP(&inputArgs.concurrency, "concurrency", "c", 1, "number of tasks to perform concurrently (default: one at a time)")

	f.StringVar(&inputArgs.scale, "scale", "log", "scale for gas price axis (options: log, linear)")

	f.Uint64Var(&inputArgs.startBlock, "start-block", 0, "starting block number (inclusive)")
	f.Uint64Var(&inputArgs.endBlock, "end-block", math.MaxUint64, "ending block number (inclusive)")
	f.StringVar(&inputArgs.targetAddr, "target-address", "", "address that will have tx sent from or to highlighted in the chart")
	f.StringVarP(&inputArgs.output, "output", "o", "tx_gasprice_chart.png", "where to save the chart image (default: tx_gasprice_chart.png)")
}
