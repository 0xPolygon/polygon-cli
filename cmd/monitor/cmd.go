package monitor

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/maticnetwork/polygon-cli/util"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	// flags
	rpcUrl         string
	batchSizeValue string
	intervalStr    string
)

// MonitorCmd represents the monitor command
var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor blocks using a JSON-RPC endpoint.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// By default, remove logs from `polycli monitor`.
		verbosityFlag := cmd.Flags().Lookup("verbosity")
		if verbosityFlag != nil && !verbosityFlag.Changed {
			return util.SetLogLevel(int(util.Silent))
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return monitor(cmd.Context())
	},
}

func init() {
	MonitorCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	MonitorCmd.PersistentFlags().StringVarP(&batchSizeValue, "batch-size", "b", "auto", "Number of requests per batch")
	MonitorCmd.PersistentFlags().StringVarP(&intervalStr, "interval", "i", "5s", "Amount of time between batch block rpc calls")
}

func checkFlags() (err error) {
	// Check rpc-url flag.
	if err = util.ValidateUrl(rpcUrl); err != nil {
		return
	}

	// Check interval duration flag.
	interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}

	// Check batch-size flag.
	if batchSizeValue == "auto" {
		batchSize = -1
	} else {
		batchSize, err = strconv.Atoi(batchSizeValue)
		if batchSize == 0 {
			return fmt.Errorf("batch-size can't be equal to zero")
		}
		if err != nil {
			// Failed to convert to int, handle the error
			return fmt.Errorf("batch-size needs to be an integer")
		}
	}

	return nil
}
