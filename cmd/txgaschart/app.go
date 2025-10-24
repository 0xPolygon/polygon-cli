package txgaschart

import (
	_ "embed"
	"math"

	"github.com/spf13/cobra"
)

// const (
// 	rpcURL     = "https://sepolia.infura.io/v3/f694519bed4a476bbe8905b8c2e00ace"
// 	startBlock = uint64(9356826)
// 	endBlock   = uint64(9358826)
// 	targetAddr = "0xeE76bECaF80fFe451c8B8AFEec0c21518Def02f9"
// )

type args struct {
	rpcURL    string
	rateLimit float64

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
	f.Uint64VarP(&inputArgs.startBlock, "start-block", "s", 0, "starting block number (inclusive)")
	f.Uint64VarP(&inputArgs.endBlock, "end-block", "e", math.MaxUint64, "ending block number (inclusive)")
	f.StringVar(&inputArgs.targetAddr, "target-address", "a", "address that will have tx sent from or to highlighted in the chart")
	f.StringVarP(&inputArgs.output, "output", "o", "tx_gasprice_chart.png", "where to save the chart image (default: tx_gasprice_chart.png)")
}
